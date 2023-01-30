package vlc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VlcManager struct {
	Log *zap.Logger
}

func (v *VlcManager) RunHTTPServer(router *gin.Engine, port string) error {
	v.Log.Info("Starting VLC")

	// endpoints
	vlc := router.Group("/api/v1")
	vlc.GET("/play/:streamURL", v.playStream)
	vlc.GET("/stop", v.stopStream)
	vlc.GET("/status", v.getStatus)
	vlc.GET("/stream", v.getStream)
	vlc.GET("/volume", v.getVolume)
	vlc.GET("/mute", v.muteStream)
	vlc.GET("/unmute", v.unmuteStream)

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	v.Log.Info("running http server")
	router.Run(server.Addr)

	return fmt.Errorf("http server stopped")

}
