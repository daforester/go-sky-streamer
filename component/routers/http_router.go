package routers

import (
	"github.com/daforester/go-sky-streamer/component/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPRouter struct {
	BaseRouter
}

func (B *HTTPRouter) RunController(context interface{}, controller controllers.Controller) {
	// Builds Controller Object on Demand with Dependency Injection Container
	c := B.container.Make(controller)
	c.(controllers.Controller).SetContext(context)

	err := c.(controllers.Controller).Run()
	if err != nil {
		context.(*gin.Context).JSON(http.StatusInternalServerError, err.Error())
	}
}
