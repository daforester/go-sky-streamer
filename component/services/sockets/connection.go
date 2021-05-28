package sockets

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"strconv"
	"sync"
	"time"
)

type Connection struct {
	uuid string
	websocket *websocket.Conn
	Related map[string]interface{}
	Tags []string
	mux *sync.Mutex
	collection *Collection
}

func (C Connection) New(socket *websocket.Conn) *Connection {
	c := new(Connection)
	c.websocket = socket
	c.Related = make(map[string]interface{})
	c.mux = new(sync.Mutex)
	return c
}

func (C *Connection) GetCollection() *Collection {
	return C.collection
}

func (C *Connection) SetCollection(collection *Collection) {
	C.collection = collection
}

func (C *Connection) GetUUID() string {
	return C.uuid
}

// Closes the WebSocket and Removes from collection
func (C *Connection) Close() error {
	if C.collection != nil {
		C.collection.RemoveConnection(C)
		C.collection = nil
	}
	return C.websocket.Close()
}

// Returns the remote network address.
func (C *Connection) PeerAddress() net.Addr {
	return C.websocket.RemoteAddr()
}

func (C *Connection) Ping() error {
	return C.WriteMessage(websocket.TextMessage, []byte("PING " + strconv.FormatInt(time.Now().Unix(), 10)))
}

func (C *Connection) Pong(v interface{}) error {
	switch v.(type) {
	case nil:
		return C.pongNil()
	case int:
		return C.pongInt(int64(v.(int)))
	case int8:
		return C.pongInt(int64(v.(int8)))
	case int16:
		return C.pongInt(int64(v.(int16)))
	case int32:
		return C.pongInt(int64(v.(int32)))
	case int64:
		return C.pongInt(v.(int64))
	case uint:
		return C.pongInt(int64(v.(uint)))
	case uint8:
		return C.pongInt(int64(v.(uint8)))
	case uint16:
		return C.pongInt(int64(v.(uint16)))
	case uint32:
		return C.pongInt(int64(v.(uint32)))
	case string:
		return C.pongString(v.(string))
	}
	return errors.New(fmt.Sprintf("unsupported type: %T", v))
}

func (C *Connection) pongInt(v int64) error {
	return C.WriteMessage(websocket.TextMessage, []byte("PONG " + strconv.FormatInt(v, 10)))
}

func (C *Connection) pongNil() error {
	return C.WriteMessage(websocket.TextMessage, []byte("PONG"))
}

func (C *Connection) pongString(v string) error {
	return C.WriteMessage(websocket.TextMessage, []byte("PONG " + v))
}

func (C *Connection) ReadMessage() (messageType int, p []byte, err error) {
	return C.websocket.ReadMessage()
}

func (C *Connection) SetUUID(uuid string) {
	C.uuid = uuid
}

func (C *Connection) String() string {
	return C.uuid + ":" + C.websocket.RemoteAddr().String() + ":" + C.websocket.LocalAddr().String()
}

func (C *Connection) SwapTags(from []string, to[]string) ([]string, error) {
	var err error
	if len(from) == len(to) {
		// 1 to 1 replacement
		for i, k := range from {
			for j, t := range C.Tags {
				if t == k {
					C.Tags[j] = to[i]
				}
			}
		}
	} else if len(to) == 1 {
		// Remove all From, Add Tag
		for _, k := range from {
			for j, t := range C.Tags {
				if t == k {
					C.Tags[j] = C.Tags[len(C.Tags)-1]
					C.Tags = C.Tags[:len(C.Tags)-1]
					break
				}
			}
		}
		C.Tags = append(C.Tags, to[0])
	} else if len(from) == 0 {
		// Add Tags
		C.Tags = append(C.Tags, to...)
	} else {
		err = errors.New("invalid replacement, from and to should be equal length or only one to length")
	}

	C.Tags = C.uniqueStrings(C.Tags)

	return C.Tags, err
}

func (C *Connection) WriteMessage(messageType int, data []byte) error {
	C.mux.Lock()
	defer C.mux.Unlock()
	return C.websocket.WriteMessage(messageType, data)
}

func (C *Connection) uniqueStrings(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}