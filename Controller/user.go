package Controller

import (
	"encoding/json"
	"messaging/Database"
	"messaging/Model"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoute(serverMux *mux.Router) {

	serverMux.HandleFunc("/v1/create", createUser)

}

// Handler for creating a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser Model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the user in the database
	newUser, err = Database.UserRepo().CreateUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}
