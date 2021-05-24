package sockets

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func GetUpgrader() *websocket.Upgrader {
	return &wsUpgrade
}
