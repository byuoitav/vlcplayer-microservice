package vlc

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/byuoitav/vlcplayer-microservice/cache"
	"github.com/byuoitav/vlcplayer-microservice/couch"
	"github.com/byuoitav/vlcplayer-microservice/data"
	"github.com/gin-gonic/gin"
	kivik "github.com/go-kivik/kivik/v3"
	"go.uber.org/zap"
)

type VlcManager struct {
	Log               *zap.Logger
	ConfigService     data.ConfigService
	ControlConfigPath string
}

func (v *VlcManager) RunHTTPServer(router *gin.Engine, port string) error {
	v.Log.Info("Starting VLC")

	var configService data.ConfigService
	dbAddress := os.Getenv("DB_ADDRESS")

	if len(dbAddress) > 0 {
		couchURL, err := url.Parse(dbAddress)
		if err != nil {
			v.Log.Fatal("invalid db address", zap.Error(err))
		}

		couchURL.User = url.UserPassword(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))

		client, err := kivik.New("couch", couchURL.String())
		if err != nil {
			v.Log.Fatal("failed to create couch client", zap.Error(err))
		}

		couch := &couch.ConfigService{
			Client:         client,
			StreamConfigDB: "stream-configs",
		}

		configService, err = cache.New(couch, os.Getenv("CACHE_DATABASE_LOCATION"))
		if err != nil {
			v.Log.Fatal("failed to build cache", zap.Error(err))
		}
		v.Log.Info("succesfully connected to config service")
	}

	v.ConfigService = configService
	v.ControlConfigPath = os.Getenv("CONTROL_CONFIG_PATH")

	// endpoints
	vlc := router.Group("/api/v1")
	vlc.PUT("/play", v.playStream)
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
