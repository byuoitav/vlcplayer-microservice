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
	vlc.GET("/play/:streamURL", v.PlayStream)
	vlc.GET("/stop", v.StopStream)
	vlc.GET("/status", v.GetStatus)
	vlc.GET("/volume", v.SetVolume)
	vlc.GET("/stream", v.GetStream)

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	v.Log.Info("running http server")
	router.Run(server.Addr)

	return fmt.Errorf("http server stopped")

}
