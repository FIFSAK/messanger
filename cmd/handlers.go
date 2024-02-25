package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var envFile, _ = godotenv.Read(".env")

type User struct {
	id       int
	username string
	password string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "OK")
	if err != nil {
		return
	}
}
func Register(writer http.ResponseWriter, request *http.Request) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", envFile["host"], envFile["port"], envFile["user"], envFile["password"], envFile["dbname"], envFile["sslmode"])
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	login := request.FormValue("login")
	password := request.FormValue("password")

	existingUser := User{}
	err = db.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&existingUser.username)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username doesn't exist, proceed with registration
			hashPass, err := HashPassword(password)
			if err != nil {
				http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", login, hashPass)
			if err != nil {
				http.Error(writer, "Failed to insert user into database", http.StatusInternalServerError)
				return
			}
			_, err = fmt.Fprintf(writer, "User registered successfully")
			return
		}
		// Other error occurred during query
		http.Error(writer, "Error checking existing user", http.StatusInternalServerError)
		return
	}

	// Username already exists
	http.Error(writer, "User already exists", http.StatusBadRequest)
}
func Login(writer http.ResponseWriter, request *http.Request) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", envFile["host"], envFile["port"], envFile["user"], envFile["password"], envFile["dbname"], envFile["sslmode"])
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	login := request.FormValue("login")
	password := request.FormValue("password")

	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	err = rows.Scan(&user.id, &user.username, &user.password)
	if user.username == "" {
		fmt.Fprintf(writer, "User not found")
		return
	}
	if CheckPasswordHash(password, user.password) {
		_, err = fmt.Fprintf(writer, "Logged in")
	} else {
		_, err = fmt.Fprintf(writer, "Wrong password")
	}
}
func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "UpdateUser")
	if err != nil {
		return
	}

}
func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "GetAllUsers")
	if err != nil {
		return
	}

}
func SendMessage(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "SendMessage")
	if err != nil {
		return
	}
}
func UpdateMessage(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "UpdateMessage")
	if err != nil {
		return
	}
}
func DeleteMessage(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "DeleteMessage")
	if err != nil {
		return
	}
}
func Notifications(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Notifications")
	if err != nil {
		return
	}
}
