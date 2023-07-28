// main.go

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var testVar int32 = 123
var testVar2 = "zort"
var foo = "zort"

// int mapping
const (
	Summer int = 0
	Autumn int = 1
	Winter int = 2
	Spring int = 3
)

// string mapping
const (
	Sommer string = "summer"
	Autum  string = "autumn"
	Wintir string = "winter"
	Sprig  string = "spring"
)

// User represents a user in the system.
type User struct {
	ID       int
	Username string
	Email    string
}

// usersDB simulates a database of users.
var usersDB = []User{
	{ID: 1, Username: "user1", Email: "user1@example.com"},
	{ID: 2, Username: "user2", Email: "user2@example.com"},
	{ID: 3, Username: "user3", Email: "user3@example.com"},
}

// handler for /users endpoint.
func usersHandler(w http.ResponseWriter, r *http.Request) {
	usersJSON, err := json.Marshal(usersDB)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
}

// handler for /status endpoint.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "Server is up and running!",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/status", statusHandler)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
