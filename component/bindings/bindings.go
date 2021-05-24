package bindings

import (
	"github.com/daforester/go-di-container/di"
	"github.com/daforester/go-sky-streamer/component/controllers/http"
	"github.com/daforester/go-sky-streamer/component/services/sockets"
	"github.com/daforester/go-sky-streamer/component/services/sockets/engine"
	"github.com/daforester/go-sky-streamer/component/websockets/engines"
	"github.com/daforester/go-sky-streamer/component/websockets/handshake"
)

type flag int64

const (
	FLAG_NONE    flag = 0
)

/*
	The following sets up all the Dependency Injections
 */
func Register(app *di.App, flags flag) {
	app.Bind((*sockets.CollectionInterface)(nil), (*sockets.Collection)(nil))
	app.When((*engines.ClientEngine)(nil)).Needs((*handshake.Interface)(nil)).Give((*handshake.ClientHandshake)(nil))
	app.Singleton((*engines.ClientEngine)(nil), app.Make((*engines.ClientEngine)(nil)))
	app.When((*http.EstablishSocketController)(nil)).Needs((*engine.Interface)(nil)).Give(app.Make((*engines.ClientEngine)(nil))).Singleton()
}
