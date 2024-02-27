package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

var db *sql.DB
var envFile, err = godotenv.Read(".env")
var jwtSecretKey = envFile["secretKey"]

func init() {

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envFile["host"], envFile["port"], envFile["user"], envFile["password"], envFile["dbname"], envFile["sslmode"])
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
}

type User struct {
	id       int
	username string
	password string
}

func jwtPayloadFromRequest(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return nil, false
	}

	// Проверяем формат токена
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return nil, false
	}

	tokenString := authHeader[len(bearerPrefix):]

	// Парсинг и валидация токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Замените 'yourSecretKey' на ваш реальный ключ.
		return []byte(envFile["secretKey"]), nil
	})

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Токен валиден, и claims успешно извлечены
		return claims, true
	} else {
		// Токен невалиден или claims не могут быть приведены к типу MapClaims
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return nil, false
	}
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

	login := request.FormValue("login")
	password := request.FormValue("password")

	existingUser := User{}
	err := db.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&existingUser.username)
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
	_ = rows.Scan(&user.id, &user.username, &user.password)
	if user.username == "" {
		fmt.Fprintf(writer, "User not found")
		return
	}
	if CheckPasswordHash(password, user.password) {
		payload := jwt.MapClaims{
			"sub": user.username,
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
	updateType := mux.Vars(request)["type"]
	login := request.FormValue("login")
	password := request.FormValue("password")
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.id, &user.username, &user.password)
	if CheckPasswordHash(password, user.password) {
		if rows != nil {
			if updateType == "password" {
				newPassword := request.FormValue("new-password")
				hashPass, _ := HashPassword(newPassword)
				db.Exec("UPDATE users SET password = $1 WHERE username = $2", hashPass, login)

			}
			if updateType == "login" {
				newLogin := request.FormValue("new-login")
				db.Exec("UPDATE users SET username = $1 WHERE username = $2", newLogin, login)
			}
		} else {
			fmt.Fprintf(writer, "User not found write credentials")
		}
	}
}
func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	login := request.FormValue("login")
	password := request.FormValue("password")
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.id, &user.username, &user.password)
	if CheckPasswordHash(password, user.password) {
		if rows != nil {
			db.Exec("DELETE FROM users WHERE username = $1", login)
			fmt.Fprintf(writer, "User deleted")
		} else {
			fmt.Fprintf(writer, "User not found write credentials")
		}
	} else {
		fmt.Fprintf(writer, "Wrong password")
	}
}
func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	payload, check := jwtPayloadFromRequest(writer, request)
	fmt.Println(payload)
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
		err = rows.Scan(&user.id, &user.username, &user.password)
		user.password = ""
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
