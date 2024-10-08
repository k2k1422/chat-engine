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
	go Message.HandleUnicastProducerMessage()
	go Message.HandleUnicastConsumerMessage()
	go Message.DeleteKeyCacheIfNotConnected()

	// Start the server on localhost port 8000 and log any errors
	log.Println("Server started on :", os.Getenv("SERVER_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), serverMux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
