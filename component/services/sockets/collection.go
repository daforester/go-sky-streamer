package sockets

import (
	"bitbucket.org/daforester/go-backend-repository/component/libs/strgen"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type CollectionInterface interface {
	LockUp()
	Unlock()
	AddConnection(connection *Connection)
	RemoveConnection(connection *Connection)
	RemoveConnectionByUUID(uuid string)
	UUIDExists(uuid string) bool
	ConnectionExists(connection *Connection) bool
	Count() int
	String() string
	GetConnections() []*Connection
	GetName() string
	All() []*Connection
	HasTag(tag string) []*Connection
	HasRelated(key string, values ...interface{}) []*Connection
	HasUUID(uuid string) []*Connection
	SendCommand(command *JSONRequest, connections ...*Connection)
	SendData(data []byte, connections ...*Connection)
	SendText(data string, connections ...*Connection)
	Send(messageType int, data []byte, connections ...*Connection)
}

type Collection struct {
	Name string
	connections []*Connection
	masterLock *sync.Mutex
	modifyLock *sync.Mutex
}

var WebSocketCollections = make(map[string]*Collection)
var lockCollections sync.Mutex

func NewCollection(name string) *Collection {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	ec, exists := WebSocketCollections[name]
	if exists {
		return ec
	}

	c := new(Collection)
	c.Name = name
	c.masterLock = new(sync.Mutex)
	c.modifyLock = new(sync.Mutex)
	WebSocketCollections[c.Name] = c
	return c
}

func LockCollections() {
	lockCollections.Lock()
}

func UnlockCollections() {
	lockCollections.Unlock()
}

func RemoveCollection(name string) {
	LockCollections()
	defer UnlockCollections()

	_, exists := WebSocketCollections[name]
	if exists {
		delete(WebSocketCollections, name)
	}
}

func RemoveCollectionIfEmpty(name string) {
	c, exists := WebSocketCollections[name]
	if exists && len(c.connections) == 0 {
		RemoveCollection(name)
	}
}

func (C Collection) New() *Collection {
	lockCollections.Lock()
	defer lockCollections.Unlock()

	c := new(Collection)
	c.masterLock = new(sync.Mutex)
	c.modifyLock = new(sync.Mutex)
	c.Name, _ = strgen.RandString(32)
	_, exists := WebSocketCollections[c.Name]
	for exists {
		c.Name, _ = strgen.RandString(32)
		_, exists = WebSocketCollections[c.Name]
	}
	WebSocketCollections[c.Name] = c
	return c
}

func (C *Collection) LockUp() {
	if C.masterLock != nil {
		C.masterLock.Lock()
	}
}

func (C *Collection) Unlock() {
	if C.masterLock != nil {
		C.masterLock.Unlock()
	}
}

func (C *Collection) AddConnection(connection *Connection) {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	connection.SetCollection(C)
	C.connections = append(C.connections, connection)
	lines := strings.Split(C.String(), "\n")
	for _, l := range lines {
		logrus.Debug(l)
	}
}

func (C *Collection) RemoveConnection(connection *Connection) {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	for n, c := range C.connections {
		if c == connection {
			connection.SetCollection(nil)
			C.connections = append(C.connections[:n], C.connections[n+1:]...)
			lines := strings.Split(C.String(), "\n")
			for _, l := range lines {
				logrus.Debug(l)
			}
			break
		}
	}
}

func (C *Collection) RemoveConnectionByUUID(uuid string) {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	for n, c := range C.connections {
		if c.uuid == uuid {
			C.connections = append(C.connections[:n], C.connections[n+1:]...)
			break
		}
	}
}

func (C *Collection) UUIDExists(uuid string) bool {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	for _, c := range C.connections {
		if c.uuid == uuid {
			return true
		}
	}
	return false
}

func (C *Collection) ConnectionExists(connection *Connection) bool {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	for _, c := range C.connections {
		if c == connection {
			return true
		}
	}
	return false
}

func (C *Collection) Count() int {
	return len(C.connections)
}

func (C *Collection) String() string {
	output := strconv.Itoa(len(C.connections)) + " Connections\n"
	for _, conn := range C.connections {
		output += conn.String() + "\n"
	}
	return output
}

func (C *Collection) GetConnections() []*Connection {
	return C.connections
}

func (C *Collection) GetName() string {
	return C.Name
}

func (C *Collection) All() []*Connection {
	return C.connections
}

func (C *Collection) HasTag(tag string) []*Connection {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	a := make([]*Connection, 0)
	for _, c := range C.connections {
		for _, t := range c.Tags {
			if t == tag {
				a = append(a, c)
			}
		}
	}
	return a
}

func (C *Collection) HasRelated(key string, values ...interface{}) []*Connection {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	a := make([]*Connection, 0)
	for _, c := range C.connections {
		value, exists := c.Related[key]
		if exists {
			for _, v := range values {
				if value == v {
					a = append(a, c)
					break
				}
			}
		}
	}
	return a
}

func (C *Collection) HasUUID(uuid string) []*Connection {
	C.modifyLock.Lock()
	defer C.modifyLock.Unlock()
	a := make([]*Connection, 0)
	for _, c := range C.connections {
		if c.uuid == uuid {
			a = append(a, c)
		}
	}
	return a
}

func (C Collection) SendCommand(command *JSONRequest, connections ...*Connection) {
	data, err := json.Marshal(command)
	if err != nil {
		return
	}
	C.Send(websocket.TextMessage, data, connections...)
}

func (C Collection) SendData(data []byte, connections ...*Connection) {
	C.Send(websocket.BinaryMessage, data, connections...)
}

func (C Collection) SendText(data string, connections ...*Connection) {
	C.Send(websocket.TextMessage, []byte(data), connections...)
}

func (C Collection) Send(messageType int, data []byte, connections ...*Connection) {
	if connections == nil {
		connections = C.GetConnections()
	}
	if len(connections) == 0 {
		return
	}

	var wg sync.WaitGroup
	active := 0
	for _, c := range connections {
		if active >= int(math.Ceil(float64(runtime.NumCPU()) / 2)) {
			wg.Wait()
			active = 0
		}
		active += 1
		wg.Add(1)

		go func(connection *Connection) {
			defer wg.Done()
			if connection == nil {
				return
			}
			err := connection.WriteMessage(messageType, data)
			if err != nil {
				logrus.Error(err.Error())
				_ = connection.websocket.Close()
			}
		}(c)
	}

	wg.Wait()
}
