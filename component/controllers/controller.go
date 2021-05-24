package controllers

import (
	"github.com/daforester/go-di-container/di"
)

type Controller interface {
	GetContext() interface{}
	Run() error
	SetContainer(di.AppInterface) Controller
	SetContext(interface{}) Controller
}

type BaseController struct {
	container di.AppInterface
	c interface{}
	context interface{}
}

func (C *BaseController) App() di.AppInterface {
	return C.container
}

func (C *BaseController) GetContext() interface{} {
	return C.context
}

func (C *BaseController) Run() error {
	return nil
}

func (C *BaseController) SetContainer(app di.AppInterface) Controller {
	C.container = app
	return C
}

func (C *BaseController) SetContext(context interface{}) Controller {
	C.c = context
	C.context = context
	return C
}
