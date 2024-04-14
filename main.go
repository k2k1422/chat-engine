package main

import (
	"encoding/json"
	"log"
	"messaging/Channel"
	"messaging/Controller"
	"messaging/Database"
	"messaging/Middleware"
	"messaging/Model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	serverMux := mux.NewRouter()

	userRouter := serverMux.PathPrefix("/api/user").Subrouter()

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("templates"))

	serverMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	// serverMux.HandleFunc("/", fs)

	Controller.UserRoute(userRouter)

	// Handle WebSocket connections
	serverMux.Handle("/ws", Middleware.JwtMiddleware(http.HandlerFunc((handleConnections))))

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", serverMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

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
	Channel.ClientsMutex.Lock()
	Channel.Clients[username] = ws
	Channel.ClientsMutex.Unlock()

	log.Printf("User %s connected", username)

	for {
		var msg Model.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			Channel.ClientsMutex.Lock()
			delete(Channel.Clients, username)
			Channel.ClientsMutex.Unlock()
			break
		}
		// Send the newly received message to the broadcast channel
		msg.FromUsername = username
		Channel.Broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-Channel.Broadcast

		// Send it out to every client that is currently connected
		Channel.ClientsMutex.Lock()
		toUser, isKey := Channel.Clients[msg.ToUser]
		if isKey {
			err := toUser.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				log.Println("Connection has been closed for user's web scoket:", msg.ToUser)

				// If an error occurs while writing, close the connection
				toUser.Close()
			}
		} else {
			log.Println("user not found")
			msg.Message = "User not found: " + msg.ToUser
			FromUser := Channel.Clients[msg.FromUsername]
			FromUser.WriteJSON(msg)
		}
		Channel.ClientsMutex.Unlock()
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
