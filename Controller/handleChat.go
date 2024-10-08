package Controller

import (
	"encoding/json"
	"log"
	"messaging/Cache"
	"messaging/Channel"
	"messaging/Database"
	"messaging/Message"
	"messaging/Model"
	"messaging/Utils"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {

	if !validateWebsocket(w, r) {
		return
	}

	username := r.Context().Value("username").(string)

	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Read username from query parameter

	// Register the WebSocket connection with the username
	Channel.AddWS(ws, username)

	log.Printf("User %s connected", username)

	Cache.LRemove("connection", os.Getenv("TOPIC_NAME"), username)
	log.Printf("Removed existing entry in redis for username:%s in server:%s", username, os.Getenv("TOPIC_NAME"))

	Cache.LPush("connection", os.Getenv("TOPIC_NAME"), username)
	log.Printf("Created a entry in redis for username:%s in server:%s", username, os.Getenv("TOPIC_NAME"))

	Message.FindAndSendTheUndelivedChat(ws, username)

	for {
		var msg Model.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err1 := Utils.Validate.Struct(msg); err != nil || err1 != nil {
			log.Printf("Validation failed for the web scoket message: %v", err)
			log.Printf("error: %v", err)
			Channel.RemoveWS(ws, username)

			// remove cache if it dosen't have the connection
			if !Channel.HasWS(username) {
				Cache.LRemove("connection", os.Getenv("TOPIC_NAME"), username)
				log.Printf("Removed the entry in redis for username:%s in server:%s", username, os.Getenv("TOPIC_NAME"))
			}
			break
		}
		// Send the newly received message to the broadcast channel
		msg.FromUsername = username

		Channel.Broadcast <- msg
	}
}

func validateWebsocket(w http.ResponseWriter, r *http.Request) bool {
	username := r.Context().Value("username").(string)
	if username == "" {
		log.Println("Username not provided")
		response := Model.Response{
			Message: "Hello, World!",
		}
		jsonData, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonData)
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return false
	}

	_, err := Database.UserRepo().GetUser(username)

	if err != nil {
		http.Error(w, "Username not found", http.StatusBadRequest)
		log.Println("Username not found")
		return false
	}

	return true
}
