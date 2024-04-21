package main

import (
	"log"
	"messaging/Controller"
	"messaging/Message"
	"messaging/Middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	serverMux := mux.NewRouter()

	userRouter := serverMux.PathPrefix("/api/user").Subrouter()

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("templates"))

	serverMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	// serverMux.HandleFunc("/", fs)

	Controller.UserRoute(userRouter)

	// Handle WebSocket connections
	serverMux.Handle("/ws", Middleware.JwtMiddleware(http.HandlerFunc((Controller.HandleConnections))))

	// Start listening for incoming chat messages
	go Message.HandleMessages1()
	go Message.HandleUnicastConsumerMessage()

	// Start the server on localhost port 8000 and log any errors
	log.Println("Server started on :8000")
	err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), serverMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
