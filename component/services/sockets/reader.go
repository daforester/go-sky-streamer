package sockets

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

type ReadFunc = func(*Connection, []byte)

func ReadSocket(connection *Connection, f ...ReadFunc) {
	done := make(chan bool)
	running := true
	go pingSocket(connection, done, &running)
	for {
		logrus.Debug("Reading Socket")
		t, d, err := connection.websocket.ReadMessage()
		if err != nil {
			logrus.Debug("Aborted Reading Socket")
			break
		}
		logrus.Debug("Read Socket")
		if t == websocket.TextMessage {
			text := string(d)
			logrus.Debug("WEBSOCKET DATA FROM: " + connection.GetUUID())
			logrus.Debug(text)
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
	defer logrus.Debug("Ping Loop Ended")
	defer close(done)
	logrus.Debug("Starting Ping Loop")
	pingTicker := time.NewTicker((time.Second * time.Duration(30)))
	for {
		select {
			case <- done:
				logrus.Debug("Ping Loop Got Termination Signal")
				*running = false
				return
			case <- pingTicker.C:
				logrus.Debug("Ping Loop Sending Ping")
				err := connection.Ping()
				if err != nil {
					*running = false
					return
				}
		}
	}
}
