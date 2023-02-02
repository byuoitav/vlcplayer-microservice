package vlc

import (
	"net/http"
	"strconv"

	vlc "github.com/adrg/libvlc-go/v3"
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

	var stream StreamURL
	err := c.ShouldBindJSON(&stream)
	if err != nil {
		v.Log.Warn("failed to bind json", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vlcPlayer, err := helpers.StartVLC()
	if err != nil {
		v.Log.Warn("failed to start vlc player", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)
		return

	}

	Player = vlcPlayer

	err = helpers.SwitchStream(Player, stream.URL)
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
