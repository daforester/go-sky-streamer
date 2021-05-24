package routers

import (
	"github.com/daforester/go-di-container/di"
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
)

type WebSocketRouter struct {
	BaseRouter
}

func NewWebSocketRouter(app di.AppInterface, engine... *engine.Engine) Router {
	r := (&WebSocketRouter{}).New(app)
	if len(engine) > 0 {
		r.RegisterRoutes(engine[0])
	}
	return r
}

func (P WebSocketRouter) New(app di.AppInterface) Router {
	p := new(WebSocketRouter)
	p.SetContainer(app)
	return p
}
