package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	ControlConfigPath string
}

func (h *Handlers) PlayStream(c *gin.Context) {
	stramURL := c.Param("streamURL")
	var args []string

	stramURL, err := url.QueryUnescape(stramURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		fmt.Errorf("error getting stream URL: %s", err)
		return
	}
}

func (h *Handlers) switchStram() error {
	return nil
}

func (h *Handlers) checkStream() bool {
	return true
}

func (h *Handlers) StopStream(c *gin.Context) {
	var args []string
}

func (h *Handlers) GetStream(c *gin.Context) {
	var args []string
}
