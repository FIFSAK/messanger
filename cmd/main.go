package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("GET")
	router.HandleFunc("/login/{type}", UpdateUser).Methods("PATCH")
	router.HandleFunc("/login", DeleteUser).Methods("DELETE")
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
