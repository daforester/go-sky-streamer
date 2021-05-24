package handshake

import (
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ClientHandshake struct {
}

func (P ClientHandshake) New() *ClientHandshake {
	p := new(ClientHandshake)

	return p
}

func (P *ClientHandshake) HandshakeSocket(ws *sockets.Connection, data map[string]string) error {
	ws.SetUUID(uuid.New().URN())
	err := ws.WriteMessage(websocket.TextMessage, []byte("UUID " + ws.GetUUID()))
	if err != nil {
		logrus.Error("Failed to send websocket UUID")
		return err
	}

	return nil
}
