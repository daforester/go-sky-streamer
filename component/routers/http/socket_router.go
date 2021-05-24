package http

import (
	"github.com/daforester/go-di-container/di"
	"github.com/daforester/go-sky-streamer/component/controllers/http"
	"github.com/daforester/go-sky-streamer/component/routers"
	"github.com/gin-gonic/gin"
)

type SocketRouter struct {
	routers.HTTPRouter
}

func NewSocketRouter(app di.AppInterface, engine *gin.Engine) routers.Router {
	r := (&SocketRouter{}).New(app)
	r.RegisterRoutes(engine)
	return r
}

func (P SocketRouter) New(app di.AppInterface) *SocketRouter {
	a := new(SocketRouter)
	a.SetContainer(app)
	return a
}

func (P *SocketRouter) RegisterRoutes(engine interface{}) {
	r := engine.(*gin.Engine)

	r.GET("/ws", func(c *gin.Context) {
		P.RunController(c, &http.EstablishSocketController{})
	})
}
