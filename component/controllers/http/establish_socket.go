package http

import (
	"github.com/daforester/go-sky-streamer/component/controllers"
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/daforester/go-sky-streamer/component/websockets/engines"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


type EstablishSocketController struct {
	controllers.HTTPController
	wsClientEngine *engines.ClientEngine
	uuid string
}

func (C EstablishSocketController) New(
	wsClientEngine *engines.ClientEngine) *EstablishSocketController {
	c := new(EstablishSocketController)
	c.wsClientEngine = wsClientEngine
	c.uuid = viper.GetString("ZEROCONF_UUID")
	return c
}

func (C *EstablishSocketController) Run() error {
	c := C.GetContext().(*gin.Context)
	return C.EstablishSocket(c)
}

func (C *EstablishSocketController) EstablishSocket(c *gin.Context) error {
	// fmt.Println(fmt.Sprintf("Client Engine: %p", C.wsClientEngine))
	// fmt.Println(fmt.Sprintf("Client Collection: %p", C.wsClientEngine.GetCollection()))

	w := c.Writer
	r := c.Request

	socket, err := sockets.GetUpgrader().Upgrade(w, r, nil)
	defer func() {
		if socket != nil {
			_ = socket.Close()
		}
	}()

	if err != nil {
		logrus.Errorf("Failed to set websocket upgrade: %+v", err)
		return err
	}

	ws := sockets.Connection{}.New(socket)

	err = C.wsClientEngine.HandshakeProcessor.HandshakeSocket(ws, nil)

	if err != nil {
		logrus.Debug("Handshake Error")
		logrus.Debug(err)
		return nil
	}

	C.wsClientEngine.AddConnection(ws)

	sockets.ReadSocket(ws, C.wsClientEngine.ReadMessage)

	C.wsClientEngine.RemoveConnection(ws)

	return nil
}

