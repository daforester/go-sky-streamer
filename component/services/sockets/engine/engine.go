package engine

import (
	"encoding/json"
	"errors"
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/daforester/go-sky-streamer/component/websockets/handshake"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type HandlerMethod int

const (
	COMMAND_HANDLER HandlerMethod = 0
	JSON_HANDLER    HandlerMethod = 1
)

type HandlerFunc func(context *Context)

type Handler struct {
	Method      HandlerMethod
	HandlerFunc HandlerFunc
}

type Interface interface {
	AddCommand(command string, handlers ...HandlerFunc)
	AddHandler(command string, method HandlerMethod, handlers ...HandlerFunc) error
	AddJSON(command string, handlers ...HandlerFunc)
	GetCollection() sockets.CollectionInterface
	ReadMessage(connection *sockets.Connection, msg []byte)
	Run()
}

type Engine struct {
	Handlers           map[string][]*Handler
	Collection         sockets.CollectionInterface
	HandshakeProcessor handshake.Interface
}

func (E Engine) New(collection sockets.CollectionInterface, handshakeProcessor handshake.Interface) *Engine {
	e := new(Engine)
	e.Handlers = make(map[string][]*Handler)
	e.Collection = collection
	e.HandshakeProcessor = handshakeProcessor

	return e
}

func (E *Engine) AddCommand(command string, handlers ...HandlerFunc) {
	err := E.AddHandler(command, COMMAND_HANDLER, handlers...)
	if err != nil {
		logrus.Error(err)
	}
}

func (E *Engine) AddHandler(command string, method HandlerMethod, handlers ...HandlerFunc) error {
	uCommand := strings.ToUpper(command)

	re := regexp.MustCompile("^[A-Z-_]+$")
	if !re.Match([]byte(uCommand)) {
		return errors.New("command may only contain A-Z, - or _")
	}

	for _, hF := range handlers {
		h := new(Handler)
		h.Method = method
		h.HandlerFunc = hF
		E.Handlers[uCommand] = append(E.Handlers[uCommand], h)
	}

	return nil
}

func (E *Engine) AddJSON(command string, handlers ...HandlerFunc) {
	err := E.AddHandler(command, JSON_HANDLER, handlers...)
	if err != nil {
		logrus.Error(err)
	}
}

func (E *Engine) Run() {
	logrus.Debug("Running Engine")
}

func (E *Engine) GetCollection() sockets.CollectionInterface {
	return E.Collection
}

func (E *Engine) ReadMessage(connection *sockets.Connection, msg []byte) {
	E.ReadMessageWithEngine(connection, msg, E)
}

func (E *Engine) ReadMessageWithEngine(connection *sockets.Connection, msg []byte, engine Interface) {
	rType, r := E.ParseData(msg)

	handlers := E.Handlers[r.GetCommand()]

	if handlers == nil || len(handlers) == 0 {
		logrus.Debug("No handler registered for " + r.GetCommand())
		return
	}

	for _, h := range handlers {
		if h.Method != rType {
			continue
		}
		go func(handler *Handler, request sockets.Request) {
			context := Context{}.New(connection, engine)
			context.Data = request.GetData()
			context.Params = request.GetParams()
			context.RequestData = msg
			handler.HandlerFunc(context)
		}(h, r)
	}
}

func (E *Engine) ParseData(input []byte) (HandlerMethod, sockets.Request) {
	logrus.Debug("Parsing input")
	logrus.Debug(string(input))
	if input[0] == '{' {
		// Assume JSON
		logrus.Debug("Assuming JSON based on first byte")
		r, err := E.parseJSONData(input)
		logrus.Debug(err)
		if err == nil {
			logrus.Debug("Parsing JSON complete")
			logrus.Debug(r)
			return JSON_HANDLER, r
		}
	}
	r := E.parseCommandData(string(input))
	logrus.Debug("Parsing Command complete")
	logrus.Debug(r)
	return COMMAND_HANDLER, r
}

func (E *Engine) parseJSONData(input []byte) (*sockets.JSONRequest, error) {
	r := new(sockets.JSONRequest)
	err := json.Unmarshal(input, r)
	return r, err
}

func (E *Engine) parseCommandData(input string) *sockets.CommandRequest {
	r := new(sockets.CommandRequest)

	trimStr := strings.TrimLeft(input, " ")
	i := strings.IndexRune(trimStr, ' ')
	if i == -1 {
		r.Command = trimStr
		return r
	}

	r.Command = trimStr[:i]
	if len(trimStr) > i+1 {
		r.Data = trimStr[i+1:]
	}
	return r
}
