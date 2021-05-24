package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
	"strings"
)

type CommandInterface interface {
	GetHandler() handlerFunc
	GetName() string
	GetOptions() []*Option
	GetOperands() []*Operand
	GetDescription() string
	GetShortDescription() string
	SetName(string) (*Command, error)
	SetHandler(handlerFunc) *Command
	SetDescription(string) *Command
}

type Command struct {
	WithOptions
	WithOperands

	name string
	shortDescription string
	longDescription string

	handler handlerFunc
}

func NewCommand(name string, handler handlerFunc, options ...*Option) (*Command, error) {
	c := new(Command)
	return c.Build(name, handler, options...)
}

func (c *Command) Build(name string, handler handlerFunc, options ...*Option) (*Command, error) {
	c.options = make([]*Option, 0)
	c.optionMapping = make(map[string]*Option)
	c.operands = make([]*Operand, 0)

	_, err := c.SetName(name)

	if err != nil {
		return nil, err
	}

	c.handler = handler

	if len(options) > 0 {
		_, err = c.AddOptions(options...)

		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Command) SetName(name string) (*Command, error) {
	if len(name) == 0 || name[0] == '-' || strings.Contains(name, " ") {
		return nil, errors.New(fmt.Sprintf("Command name has to be an alphanumeric string not starting with dash, found '%s'", name))
	}
	c.name = name
	return c, nil
}

func (c *Command) SetHandler(handler handlerFunc) *Command {
	c.handler = handler
	return c
}

func (c *Command) SetDescription(longDescription string) *Command {
	c.longDescription = longDescription
	if c.shortDescription == "" {
		c.shortDescription = longDescription
	}
	return c
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) GetHandler() handlerFunc {
	return c.handler
}

func (c *Command) GetDescription() string {
	return c.longDescription
}

func (c *Command) Description() string {
	return c.longDescription
}

func (c *Command) GetShortDescription() string {
	return c.shortDescription
}

func (c *Command) ShortDescription() string {
	return c.shortDescription
}