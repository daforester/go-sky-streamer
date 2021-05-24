package engines

import (
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
	"github.com/daforester/go-sky-streamer/component/websockets/handshake"
)

type ClientEngine struct {
	engine.Engine
}

func (E ClientEngine) New(collection sockets.CollectionInterface, handshakeProcessor handshake.Interface) *ClientEngine {
	e := new(ClientEngine)
	e.Handlers = make(map[string][]*engine.Handler)
	e.Collection = collection
	e.HandshakeProcessor = handshakeProcessor

	return e
}

func (E *ClientEngine) AddConnection(connection *sockets.Connection) {
	E.Collection.AddConnection(connection)
}

func (E *ClientEngine) RemoveConnection(connection *sockets.Connection) {
	E.Collection.RemoveConnection(connection)
}

func (E *ClientEngine) ReadMessage(connection *sockets.Connection, msg []byte) {
	E.ReadMessageWithEngine(connection, msg, E)
}
