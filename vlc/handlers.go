package vlc

import (
	"log"
	"net/http"
	"net/url"

	"github.com/byuoitav/vlcplayer-microservice/vlc/helpers"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	ControlConfigPath string
}

func (v *VlcManager) PlayStream(c *gin.Context) {
	streamURL := c.Param("streamURL")

	streamURL, err := url.QueryUnescape(streamURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		//fmt.Errorf("error getting stream URL: %v", err)
		return
	}
	vlcPlayer, err := helpers.StartVLC()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		//fmt.Errorf("error starting VLC: %s", err)
		return
	}
	err = helpers.SwitchStream(vlcPlayer, streamURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		//fmt.Errorf("error switching stream: %s", err)
		return
	}
	c.JSON(http.StatusOK, "stream playing")
}

// StopStream stops the stream that is currently playing
func (v *VlcManager) StopStream(c *gin.Context) {
	log.Println("Stopping stream player...")

}

// GetStream returns the stream that is currently playing
func (v *VlcManager) GetStream(c *gin.Context) {
	log.Println("Getting stream...")

}

// GetStatus returns the status of the stream player
func (v *VlcManager) GetStatus(c *gin.Context) {
	log.Println("Getting status...")

}

// SetVolume sets the volume of the stream player
func (v *VlcManager) SetVolume(c *gin.Context) {
	log.Println("Setting volume...")

}
