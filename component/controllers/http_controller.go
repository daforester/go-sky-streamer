package controllers

import (
	"github.com/gin-gonic/gin"
)

type HTTPController struct {
	BaseController
	c *gin.Context
	context *gin.Context
}

func (C *HTTPController) GetContext() interface{} {
	return C.context
}

func (C *HTTPController) SetContext(context interface{}) Controller {
	C.c = context.(*gin.Context)
	C.context = context.(*gin.Context)
	return C
}

func (C *HTTPController) Run() error {
	return nil
}
