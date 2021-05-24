package handshake

import (
	"github.com/daforester/go-sky-streamer/component/services/sockets"
)

type Interface interface {
	HandshakeSocket(*sockets.Connection, map[string]string) error
}
