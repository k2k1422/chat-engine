package Channel

import (
	"messaging/Model"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients         = make(map[string]*websocket.Conn) // Connected clients mapped by username
	ClientsMutex    sync.Mutex                         // Mutex for thread-safe access to clients map
	Broadcast       = make(chan Model.Message)         // Broadcast channel
	ConsumerUnicast = make(chan Model.Message)         // Broadcast channel
)
