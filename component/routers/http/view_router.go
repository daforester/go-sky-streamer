package http

import (
	"github.com/daforester/go-di-container/di"
	"github.com/daforester/go-sky-streamer/component/routers"
	"github.com/daforester/go-sky-streamer/component/services/ginhelper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type ViewRouter struct {
	routers.HTTPRouter
}

func NewViewRouter(app di.AppInterface, engine... *gin.Engine) routers.Router {
	r := (&ViewRouter{}).New(app)
	if len(engine) > 0 {
		r.RegisterRoutes(engine[0])
	}
	return r
}

func (P ViewRouter) New(app di.AppInterface) routers.Router {
	a := new(ViewRouter)
	a.SetContainer(app)
	return a
}

func (P *ViewRouter) RegisterRoutes(engine interface{}) {
	r := engine.(*gin.Engine)
	r.GET("/", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=86400")
		c.HTML(http.StatusOK, "index.gohtml", nil)
	})

	r.GET("/service-worker.js", func(c *gin.Context) {
		modified, eTag := ginhelper.FileModified(c, viper.GetString("ROOT_PATH") + "/public/service-worker.js")
		if modified {
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "application/javascript")
			c.Header("Etag", "\"" + eTag + "\"")
			c.File(viper.GetString("ROOT_PATH") + "/public/service-worker.js")
		} else {
			c.Status(http.StatusNotModified)
		}
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=86400")
		c.Header("Content-Type", mime.TypeByExtension(".ico"))
		c.File(viper.GetString("ROOT_PATH") + "/public/favicon.ico")
	})

	r.GET("/js/client/:filename", func(c *gin.Context) {
		jsClientPath := viper.GetString("ROOT_PATH") + "/public/js/client/"
		mainJs := "main.[a-z0-9]+.js"
		fileName := c.Param("filename")
		if fileName == "main.js" {
			_ = filepath.Walk(jsClientPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatal(err)
					return err
				}
				match, err := regexp.Match(mainJs, []byte(filepath.Base(path)))
				if err != nil {
					log.Fatal(err)
					return err
				}
				if match {
					fileName = filepath.Base(path)
				}
				return nil
			})
		}

		modified, eTag := ginhelper.FileModified(c, jsClientPath + fileName)
		if modified {
			c.Header("Cache-Control", "no-cache")
			c.Header("Etag", "\"" + eTag + "\"")
			c.Header("Content-Type", mime.TypeByExtension(filepath.Ext(fileName)))
			c.File(jsClientPath + fileName)
		} else {
			c.Status(http.StatusNotModified)
		}
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "index.gohtml", map[string]string{"XErrorMessage": "Page not found"})
	})
}
