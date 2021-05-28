package sockets

import (
	"github.com/gorilla/websocket"
	"time"
)

type ReadFunc = func(*Connection, []byte)

func ReadSocket(connection *Connection, f ...ReadFunc) {
	done := make(chan bool)
	running := true
	go pingSocket(connection, done, &running)
	for {
		t, d, err := connection.websocket.ReadMessage()
		if err != nil {
			break
		}
		if t == websocket.TextMessage {
			text := string(d)
			if len(text) == 4 && text[:4] == "PING" {
				_ = connection.Pong(nil)
			}
			if text[:5] == "PING " {
				_ = connection.Pong(text[5:])
			}

			for n := 0; n < len(f); n ++ {
				if f[n] != nil {
					f[n](connection, d)
				}
			}
		}
	}
	if running {
		done <- true
	}
}

func pingSocket(connection *Connection, done chan bool, running *bool) {
	defer close(done)
	pingTicker := time.NewTicker((time.Second * time.Duration(30)))
	for {
		select {
			case <- done:
				*running = false
				return
			case <- pingTicker.C:
				err := connection.Ping()
				if err != nil {
					*running = false
					return
				}
		}
	}
}
