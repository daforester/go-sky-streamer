package controllers

import (
	"github.com/daforester/getopt-golang/getopt"
)

type CLIController struct {
	BaseController
	c *getopt.GetOpt
	context *getopt.GetOpt
}

func (C *CLIController) GetContext() interface{} {
	return C.context
}

func (C *CLIController) SetContext(context interface{}) Controller {
	C.c = context.(*getopt.GetOpt)
	C.context = context.(*getopt.GetOpt)
	return C
}

func (C *CLIController) Run() error {
	return nil
}
