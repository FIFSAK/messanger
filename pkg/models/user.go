package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	. "messanger/pkg/auth"
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

type User struct {
	Id       int
	Username string
	Password string
}

func (user *User) RegisterUser(login string, password string, writer http.ResponseWriter) error {
	err := db.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Username doesn't exist, proceed with registration
			hashPass, err := HashPassword(password)
			if err != nil {
				http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
				return nil
			}
			_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", login, hashPass)
			if err != nil {
				http.Error(writer, "Failed to insert user into database", http.StatusInternalServerError)
				return nil
			}
			_, err = fmt.Fprintf(writer, "User registered successfully")
			return nil
		}
		// Other error occurred during query
		http.Error(writer, "Error checking existing user", http.StatusInternalServerError)
		return nil
	}

	http.Error(writer, "User already exists", http.StatusBadRequest)
	return nil
}
func (user *User) LoginUser(login string, password string, writer http.ResponseWriter) error {
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)
	if user.Username == "" {
		_, err := fmt.Fprintf(writer, "User not found")
		if err != nil {
			return err
		}
		return nil
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

		err = json.NewEncoder(writer).Encode(t)
		if err != nil {
			return err
		}

	} else {
		_, err := fmt.Fprintf(writer, "Wrong password")
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) UpdateUser(login string, updateType string, writer http.ResponseWriter, request *http.Request) {
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		if updateType == "password" {
			newPassword := request.FormValue("new-password")
			hashPass, _ := HashPassword(newPassword)
			_, err := db.Exec("UPDATE users SET password = $1 WHERE username = $2", hashPass, login)
			if err != nil {
				log.Println("Error updating password:", err)
				_, err := fmt.Fprintf(writer, "Failed to update password")
				if err != nil {
					return
				}
				return
			}
			_, err = fmt.Fprintf(writer, "updated password successfully")
			if err != nil {
				return
			}
		}
		if updateType == "login" {
			newLogin := request.FormValue("new-login")
			existingUser := User{}
			err := db.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&existingUser.Username)
			if err == nil {
				_, err := db.Exec("UPDATE users SET username = $1 WHERE username = $2", newLogin, login)
				if err != nil {
					return
				}
				_, err = fmt.Fprintf(writer, "updated login successfully")
				if err != nil {
					return
				}
			} else {
				_, err := fmt.Fprintf(writer, "User with this login already exists")
				if err != nil {
					return
				}
			}
		}
	} else {
		_, err := fmt.Fprintf(writer, "User not found write credentials")
		if err != nil {
			return
		}
	}
}

func (user *User) DeleteUser(login string, writer http.ResponseWriter) {
	rows := db.QueryRow("SELECT * FROM users WHERE username = $1", login)
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		_, err := db.Exec("DELETE FROM users WHERE username = $1", login)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(writer, "User deleted")
		if err != nil {
			return
		}
	} else {
		_, err := fmt.Fprintf(writer, "User not found write credentials")
		if err != nil {
			return
		}
	}
}

func (user *User) GetAllUsers(writer http.ResponseWriter) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(writer, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	var users []User
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
