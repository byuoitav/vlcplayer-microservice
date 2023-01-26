package vlc

import (
	"net/http"
	"net/url"
	"strconv"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/byuoitav/vlcplayer-microservice/vlc/helpers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Player *vlc.Player

func (v *VlcManager) StartVLC() (*vlc.Player, error) {
	v.Log.Debug("starting vlc player")

	player, err := helpers.StartVLC()
	if err != nil {
		v.Log.Warn("failed to start vlc player", zap.Error(err))
		return nil, err
	}
	Player = player
	return player, nil
}

func (v *VlcManager) PlayStream(c *gin.Context) {
	v.Log.Debug("playing stream", zap.String("streamURL", c.Param("streamURL")))
	streamURL := c.Param("streamURL")

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

	err = helpers.SwitchStream(vlcPlayer, streamURL)
	if err != nil {
		v.Log.Warn("failed to switch stream", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "stream playing")
}

// StopStream stops the stream that is currently playing
func (v *VlcManager) StopStream(c *gin.Context) {
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
func (v *VlcManager) GetStream(c *gin.Context) {
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
func (v *VlcManager) GetStatus(c *gin.Context) {
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
func (v *VlcManager) SetVolume(c *gin.Context) {
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
func (v *VlcManager) GetVolume(c *gin.Context) {
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
func (v *VlcManager) MuteStream(c *gin.Context) {
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
func (v *VlcManager) UnmuteStream(c *gin.Context) {
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
