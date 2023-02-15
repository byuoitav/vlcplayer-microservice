package vlc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/byuoitav/vlcplayer-microservice/data"
	"github.com/byuoitav/vlcplayer-microservice/vlc/helpers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Player *vlc.Player

type StreamURL struct {
	URL string `json:"url"`
}

func (v *VlcManager) startVLC() (*vlc.Player, error) {
	v.Log.Debug("starting vlc player")

	player, err := helpers.StartVLC()
	if err != nil {
		v.Log.Warn("failed to start vlc player", zap.Error(err))
		return nil, err
	}
	Player = player
	return player, nil
}

func (v *VlcManager) playStream(c *gin.Context) {
	v.Log.Debug("playing stream")

	streamURL := c.Param("streamURL")
	fmt.Println("streamURL: ", streamURL)

	if v.ConfigService != nil {
		stream, err := v.ConfigService.GetStreamConfig(c.Request.Context(), streamURL)
		if err == nil && stream.Secret != "" {
			// the token is everything after the base URL
			token, err := v.generateToken(stream)
			if err != nil {
				v.Log.Error("error generating secure token: %s", zap.Error(err))
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			v.Log.Info("generated secure token", zap.String("token", token))
			streamURL += token
		}

	}
	streamURL, err := url.QueryUnescape(streamURL)
	if err != nil {
		v.Log.Warn("failed to unescape stream url", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	vlcPlayer, err := helpers.StartVLC()
	if err != nil {
		v.Log.Warn("failed to start vlc player", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return

	}

	Player = vlcPlayer

	err = helpers.SwitchStream(Player, streamURL)
	if err != nil {
		v.Log.Warn("failed to switch stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "stream playing")
}

// StopStream stops the stream that is currently playing
func (v *VlcManager) stopStream(c *gin.Context) {
	v.Log.Debug("stopping stream player")
	if Player == nil {
		v.Log.Warn("no stream to stop")
		c.JSON(http.StatusInternalServerError, "no stream to stop")
		return
	}
	err := helpers.StopStream(Player)
	if err != nil {
		v.Log.Warn("failed to stop stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	Player = nil
	c.JSON(http.StatusOK, "stream stopped")
}

// GetStream returns the stream that is currently playing
func (v *VlcManager) getStream(c *gin.Context) {
	v.Log.Debug("getting stream player")
	if Player == nil {
		v.Log.Warn("no stream to get")
		c.JSON(http.StatusInternalServerError, "no stream to get")
		return
	}
	stream, err := helpers.GetStream(Player)
	if err != nil {
		v.Log.Warn("failed to get stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, stream)
}

// GetStatus returns the status of the stream player
func (v *VlcManager) getStatus(c *gin.Context) {
	v.Log.Debug("getting status of stream player")
	if Player == nil {
		v.Log.Warn("no stream to get status of")
		c.JSON(http.StatusInternalServerError, "no stream to get status of")
		return
	}
	status, err := helpers.GetPlaybackStatus(Player)
	if err != nil {
		v.Log.Warn("failed to get status of stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

// SetVolume sets the volume of the stream player
func (v *VlcManager) setVolume(c *gin.Context) {
	v.Log.Debug("setting volume of stream player")
	if Player == nil {
		v.Log.Warn("no stream to set volume of")
		c.JSON(http.StatusInternalServerError, "no stream to set volume of")
		return
	}
	volume, err := strconv.Atoi(c.Param("volume"))
	if err != nil {
		v.Log.Warn("failed to convert volume to int", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = helpers.SetVolume(Player, volume)
	if err != nil {
		v.Log.Warn("failed to set volume of stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "volume set to "+c.Param("volume"))
}

// GetVolume returns the volume of the stream player
func (v *VlcManager) getVolume(c *gin.Context) {
	v.Log.Debug("getting volume of stream player")
	if Player == nil {
		v.Log.Warn("no stream to get volume of")
		c.JSON(http.StatusInternalServerError, "no stream to get volume of")
		return
	}
	volume, err := helpers.Volume(Player)
	if err != nil {
		v.Log.Warn("failed to get volume of stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, volume)
}

// MuteStream mutes the stream that is currently playing
func (v *VlcManager) muteStream(c *gin.Context) {
	v.Log.Debug("muting stream player")
	if Player == nil {
		v.Log.Warn("no stream to mute")
		c.JSON(http.StatusInternalServerError, "no stream to mute")
		return
	}
	err := helpers.Mute(Player)
	if err != nil {
		v.Log.Warn("failed to mute stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "stream muted")
}

// UnmuteStream unmutes the stream that is currently playing
func (v *VlcManager) unmuteStream(c *gin.Context) {
	v.Log.Debug("unmuting stream player")
	if Player == nil {
		v.Log.Warn("no stream to unmute")
		c.JSON(http.StatusInternalServerError, "no stream to unmute")
		return
	}
	err := helpers.Unmute(Player)
	if err != nil {
		v.Log.Warn("failed to unmute stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "stream unmuted")
}

/*
Example of a stream URL:
https://streaming-stg.byu.edu:8443/live-stream01/stream01/playlist.m3u8?zbyutokenstarttime=1675881066&zbyutokenendtime=1675917066&zbyutokenhash=_Y5WqFc8VwT0ahWYWO9Trj-4Wz8Ap7NXjP_gfyENW_k=
*/

// generateToken generates a token for the encrypted streams
func (v *VlcManager) generateToken(stream data.Stream) (string, error) {
	duration, err := time.ParseDuration(stream.Duration)
	if err != nil {
		v.Log.Warn("failed to parse duration", zap.Error(err))
		return "", err
	}

	start := time.Now().UTC()
	end := start.Add(duration)

	url := fmt.Sprintf("%s?%s&%sendtime=%d&%sstarttime=%d", stream.URL, stream.Secret, stream.QueryPrefix, end.Unix(), stream.QueryPrefix, start.Unix())
	input := strings.NewReader(url)
	hash := sha256.New()

	if _, err := io.Copy(hash, input); err != nil {
		v.Log.Warn("failed to copy hash", zap.Error(err))
		return "", err
	}

	finalHash := string(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	finalHash = strings.ReplaceAll(finalHash, "+", "-")
	finalHash = strings.ReplaceAll(finalHash, "/", "_")

	return fmt.Sprintf("?%sstarttime=%d&%sendtime=%d&%shash=%s", stream.QueryPrefix, start.Unix(), stream.QueryPrefix, end.Unix(), stream.QueryPrefix, finalHash), nil
}
