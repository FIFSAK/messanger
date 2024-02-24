package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", HealthCheck).Methods("GET") // Теперь HealthCheck доступен без указания пакета
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
