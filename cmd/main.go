package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"messanger/pkg/models"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Postgres driver
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	db, err := initializeDB()
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}
	defer db.Close()

	// Create a new Router
	router := mux.NewRouter()

	// Initialize User model with DB connection
	userModel := models.NewUserModel(db)

	// Setup routes with handlers
	setupRoutes(router, userModel)

	// Start the server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func initializeDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("host"), os.Getenv("port"), os.Getenv("user"),
		os.Getenv("password"), os.Getenv("dbname"), os.Getenv("sslmode"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func setupRoutes(router *mux.Router, userModel *models.UserModel) {
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/register", RegisterHandler(userModel)).Methods("POST")
	router.HandleFunc("/login", Login(userModel)).Methods("GET")
	router.HandleFunc("/login/{type}", Update(userModel)).Methods("PATCH")
	router.HandleFunc("/login", DeleteUserHandler(userModel)).Methods("DELETE")
	router.HandleFunc("/users", GetAllUsersHandler(userModel)).Methods("GET")
	router.HandleFunc("/message/send", GetSendMessageHandler(userModel)).Methods("GET")
	router.HandleFunc("/message/received", GetReceivedMessageHandler(userModel)).Methods("GET")
	router.HandleFunc("/message", SendMessageHandler(userModel)).Methods("POST")
	router.HandleFunc("/message", UpdateMessageHandler(userModel)).Methods("PATCH")
	router.HandleFunc("/message/{id}", DeleteMessageHandler(userModel)).Methods("DELETE")
	router.HandleFunc("/message/notifications", GetUnreadMessageHandler(userModel)).Methods("GET")
}
