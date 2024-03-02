package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	. "messanger/pkg/auth"
	. "messanger/pkg/models"
	"net/http"
	"os"
	"time"
)

var db *sql.DB
var jwtSecretKey = os.Getenv("secretKey")

func init() {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"), os.Getenv("sslmode"))
	db, _ = sql.Open("postgres", connStr)
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "OK")
	if err != nil {
		return
	}
}
func Register(writer http.ResponseWriter, request *http.Request) {

	login := request.FormValue("login")
	password := request.FormValue("password")

	existingUser := User{}
	err := db.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&existingUser.Username)
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

	login := request.FormValue("login")
	password := request.FormValue("password")

	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)
	if user.Username == "" {
		fmt.Fprintf(writer, "User not found")
		return
	}
	if CheckPasswordHash(password, user.Password) {
		payload := jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		t, err := token.SignedString([]byte(jwtSecretKey))

		if err != nil {
			http.Error(writer, "jwt token signing", http.StatusBadRequest)
		}

		json.NewEncoder(writer).Encode(t)

	} else {
		fmt.Fprintf(writer, "Wrong password")
	}
}
func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	payload, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
	updateType := mux.Vars(request)["type"]
	login := payload["sub"].(string)
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		if updateType == "password" {
			newPassword := request.FormValue("new-password")
			hashPass, _ := HashPassword(newPassword)
			_, err := db.Exec("UPDATE users SET password = $1 WHERE username = $2", hashPass, login)
			if err != nil {
				log.Println("Error updating password:", err)
				fmt.Fprintf(writer, "Failed to update password")
				return
			}
			fmt.Fprintf(writer, "updated password successfully")
		}
		if updateType == "login" {
			newLogin := request.FormValue("new-login")
			existingUser := User{}
			err := db.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&existingUser.Username)
			if err == nil {
				db.Exec("UPDATE users SET username = $1 WHERE username = $2", newLogin, login)
				fmt.Fprintf(writer, "updated login successfully")
			} else {
				fmt.Fprintf(writer, "User with this login already exists")
			}
		}
	} else {
		fmt.Fprintf(writer, "User not found write credentials")
	}
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	payload, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
	login := payload["sub"].(string)
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		db.Exec("DELETE FROM users WHERE username = $1", login)
		fmt.Fprintf(writer, "User deleted")
	} else {
		fmt.Fprintf(writer, "User not found write credentials")
	}
}

func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	payload, check := JwtPayloadFromRequest(writer, request)
	fmt.Println(payload["sub"])
	if !check {
		return
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(writer, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	users := []User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		user.Password = ""
		fmt.Println(user)
		if err != nil {
			http.Error(writer, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	_, err = fmt.Fprintf(writer, "%v", users)
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
