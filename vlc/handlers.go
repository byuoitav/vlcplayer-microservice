package vlc

import (
	"net/http"
	"net/url"

	"github.com/byuoitav/vlcplayer-microservice/vlc/helpers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	ControlConfigPath string
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

	//err := helpers.StopStream(vlcPlayer)

}

// GetStream returns the stream that is currently playing
func (v *VlcManager) GetStream(c *gin.Context) {
	v.Log.Debug("getting stream player")

	//stream, err := helpers.GetStream(vlcPlayer)

}

// GetStatus returns the status of the stream player
func (v *VlcManager) GetStatus(c *gin.Context) {
	v.Log.Debug("getting status of stream player")

	//status, err := helpers.GetPlaybackStatus(vlcPlayer)

}

// SetVolume sets the volume of the stream player
func (v *VlcManager) SetVolume(c *gin.Context) {
	v.Log.Debug("setting volume of stream player")

	//volume, err := helpers.SetVolume(vlcPlayer, c.Param("volume"))

}
