package client

import (
	"github.com/daforester/go-sky-streamer/component/controllers"
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
	"github.com/daforester/go-sky-streamer/component/stream"
	"github.com/daforester/go-sky-streamer/component/websockets/actions"
	"github.com/sirupsen/logrus"
)


type GetICEController struct {
	controllers.WebSocketController
}

func (C GetICEController) New() *GetICEController {
	c := new(GetICEController)

	return c
}

func (C *GetICEController) Run() error {
	c := C.GetContext().(*engine.Context)
	r := C.GetICE(c)

	if r != nil {
		c.JSON(r.ResponseData)
	}

	return nil
}

func (C *GetICEController) GetICE(c *engine.Context)  *actions.JSONResponse {

	data := c.GetDataMap()
	offer, e := data["Offer"]

	if !e {
		logrus.Debug("No offer")
		return nil
	}
	s := C.App().Make((*stream.Stream)(nil)).(*stream.Stream)

	localDescription := s.AddOffer(offer.(string))

	msg := sockets.JSONRequest{}.New()
	msg.Command = "ICE_DATA"
	msg.Data["Offer"] = localDescription

	return actions.StandardResponse(msg)
}
