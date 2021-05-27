package routers

import (
	"github.com/daforester/go-di-container/di"
	"github.com/daforester/go-sky-streamer/component/controllers"
)

type Router interface {
	RegisterRoutes(interface{})
	SetContainer(container di.AppInterface)
	RunController(interface{}, controllers.Controller)
}

type BaseRouter struct {
	container di.AppInterface
}

func (B *BaseRouter) RegisterRoutes(engine interface{}) {

}

func (B *BaseRouter) SetContainer(container di.AppInterface) {
	B.container = container
}

func (B *BaseRouter) RunController(context interface{}, controller controllers.Controller) {
	c := B.container.Make(controller)
	c.(controllers.Controller).SetContainer(B.container)
	c.(controllers.Controller).SetContext(context)

	_ = c.(controllers.Controller).Run()
}
