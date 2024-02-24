package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("GET")
	router.HandleFunc("/login", UpdateUser).Methods("UPDATE")
	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/message/{id}", SendMessage).Methods("POST")
	router.HandleFunc("/message/{id}", UpdateMessage).Methods("UPDATE")
	router.HandleFunc("/message/{id}", DeleteMessage).Methods("DELETE")
	router.HandleFunc("/notifications", Notifications).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
