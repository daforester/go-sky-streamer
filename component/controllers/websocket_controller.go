package controllers

import (
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
)

type WebSocketController struct {
	BaseController
	c *engine.Context
	context *engine.Context
}

func (C *WebSocketController) GetContext() interface{} {
	return C.context
}

func (C *WebSocketController) SetContext(context interface{}) Controller {
	C.c = context.(*engine.Context)
	C.context = context.(*engine.Context)
	return C
}

func (C *WebSocketController) Run() error {
	return nil
}
