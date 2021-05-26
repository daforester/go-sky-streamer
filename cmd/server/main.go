package main

// https://github.com/dwoja22/ffmpeg-webrtc

import (
	"fmt"
	"github.com/daforester/go-di-container/di"
	. "github.com/daforester/go-sky-streamer/component/bindings"
	"github.com/daforester/go-sky-streamer/component/capture"
	"github.com/daforester/go-sky-streamer/component/routers/http"
	"github.com/daforester/go-sky-streamer/component/routers/websocket"
	"github.com/daforester/go-sky-streamer/component/websockets/engines"
	"github.com/daforester/go-sky-streamer/setup"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	setup.Config()
	setup.MimeTypes()

	// Create Dependency Injection App
	app := di.New()
	Register(app, FLAG_NONE)
	logrus.SetLevel(logrus.DebugLevel)

	r := gin.Default()
	r.LoadHTMLGlob(viper.GetString("ROOT_PATH") + "/public/*.gohtml")

	// HTTP Default settings
	r.Use(location.New(location.Config{
		Scheme: viper.GetString("HTTP_SCHEMA"),
		Host:   viper.GetString("HTTP_HOST"),
		Base:   viper.GetString("HTTP_BASE"),
	}))

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"DELETE", "GET", "POST", "PUT"},
		AllowHeaders:    []string{"Origin", "Content-Type", "X-Auth-Token"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	// Client engine handles websocket connection requests from Users
	wsClientEngine := app.Make((*engines.ClientEngine)(nil)).(*engines.ClientEngine)
	wsClientEngine.Run()

	// Register HTTP request routes
	http.NewViewRouter(app, r)
	http.NewSocketRouter(app, r)

	// Register Websocket request routes
	websocket.NewClientRouter(app, wsClientEngine)

	// Start Capture
	capture := app.Make((*capture.Capture)(nil)).(*capture.Capture)
	capture.Start()

	// Run HTTP server
	_ = r.Run(":" + viper.GetString("HTTP_PORT"))

	fmt.Println("Pre-exit")
	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
		fmt.Println("Exit by User")
	case <-time.After(time.Second * 120):
		// Exit by timeout
		fmt.Println("Exit by timeout")
	}
	fmt.Println("Post-exit")
}
