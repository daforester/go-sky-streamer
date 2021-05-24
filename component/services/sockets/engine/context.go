package engine

import (
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Context struct {
	Connection  *sockets.Connection
	Engine      Interface
	Data        interface{}
	Params      interface{}
	RequestData []byte
}

func (C Context) New(connection *sockets.Connection, engine Interface) *Context {
	c := new(Context)
	c.Connection = connection
	c.Engine = engine

	return c
}

func (C *Context) GetDataMap() map[string]interface{} {
	switch C.Data.(type) {
	case map[string]interface{}:
		return C.Data.(map[string]interface{})
	default:
		logrus.Debug("Missing Data Map")
		return nil
	}
}

func (C *Context) JSON(data []byte) {
	err := C.Connection.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logrus.Error(err.Error())
		_ = C.Connection.Close()
	}
}
