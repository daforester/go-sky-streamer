package client

import (
	"github.com/daforester/go-sky-streamer/component/controllers"
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)


type GetICEController struct {
	controllers.WebSocketController
}

func (C GetICEController) New() *GetICEController {
	c := new(GetICEController)

	return c
}

func (C *GetICEController) Run() error {
	c := C.GetContext().(*engine.Context)
	C.GetICE(c)

	return nil
}

func (C *GetICEController) GetICE(c *engine.Context) {
	data := []byte("")
	
	err := c.Connection.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logrus.Debug(err)
		return
	}

	return
}
