package websocket

import (
	"github.com/daforester/go-di-container/di"
	"github.com/daforester/go-sky-streamer/component/controllers/websocket/client"
	"github.com/daforester/go-sky-streamer/component/routers"
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
	"github.com/gorilla/websocket"
)

type ClientRouter struct {
	routers.WebSocketRouter
}

func NewClientRouter(app di.AppInterface, engine... engine.Interface) routers.Router {
	r := (&ClientRouter{}).New(app)
	if len(engine) > 0 {
		r.RegisterRoutes(engine[0])
	}
	return r
}

func (P ClientRouter) New(app di.AppInterface) routers.Router {
	p := new(ClientRouter)
	p.SetContainer(app)
	return p
}

func (P *ClientRouter) RegisterRoutes(e interface{}) {
	P.RegisterPeerRoutes(e.(engine.Interface))
}

func (P *ClientRouter) RegisterPeerRoutes(r engine.Interface) {
	r.AddJSON("GET_ICE", func(c *engine.Context) {
		P.RunController(c, &client.GetICEController{})
	})
	r.AddJSON("GET_STATUS", func(c *engine.Context) {
		P.RunController(c, &client.GetICEController{})
	})
	r.AddCommand("PING", func(c *engine.Context) {
		_ = c.Connection.WriteMessage(websocket.TextMessage, []byte(c.Data.(string)))
	})
	r.AddCommand("PONG", func(c *engine.Context) {

	})
}
