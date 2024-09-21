package Channel

import (
	"errors"
	"log"
	"messaging/Model"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients         = make(map[string][]*websocket.Conn) // Connected clients mapped by username
	ClientsMutex    sync.Mutex                           // Mutex for thread-safe access to clients map
	Broadcast       = make(chan Model.Message)           // Broadcast channel
	ConsumerUnicast = make(chan Model.Message)           // Broadcast channel
)

func AddWS(ws *websocket.Conn, username string) error {
	ClientsMutex.Lock()
	if connList, ok := Clients[username]; ok {
		Clients[username] = append(connList, ws)
	} else {
		Clients[username] = []*websocket.Conn{ws}
	}
	ClientsMutex.Unlock()
	return nil
}

func HasWS(username string) bool {
	if _, ok := Clients[username]; ok {
		return true
	}
	return false
}

func WriteJSONWS(msg Model.Message, username string) error {
	log.Printf("For user: %s, number of connection: %d", username, len(Clients[username]))
	sentOneMsgSucess := false
	if connList, ok := Clients[username]; ok {
		for _, conn := range connList {
			err := conn.WriteJSON(msg)
			if err == nil {
				sentOneMsgSucess = true
			} else {
				RemoveWS(conn, username)
			}
		}
	}
	if !sentOneMsgSucess {
		return errors.New("not able to send a sucessfull msg")
	}
	return nil
}

func RemoveWS(ws *websocket.Conn, username string) error {
	ClientsMutex.Lock()
	newConnList := []*websocket.Conn{}
	if connList, ok := Clients[username]; ok {
		for _, conn := range connList {
			if conn != ws {
				newConnList = append(newConnList, conn)
			}
		}
	}

	if len(newConnList) != 0 {
		Clients[username] = newConnList
	} else {
		delete(Clients, username)
	}
	ClientsMutex.Unlock()
	return nil
}
