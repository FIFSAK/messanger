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

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

type User struct {
	Id       int
	Username string
	Password string
}

// RegisterUser регистрирует нового пользователя в системе.
func (m *UserModel) RegisterUser(login string, password string, writer http.ResponseWriter) error {
	// Проверяем, существует ли пользователь
	var username string
	err := m.DB.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Имя пользователя не существует, продолжаем регистрацию
			hashPass, err := HashPassword(password)
			if err != nil {
				http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
				return err // Возвращаем ошибку наверх
			}
			_, err = m.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", login, hashPass)
			if err != nil {
				http.Error(writer, "Failed to insert user into database", http.StatusInternalServerError)
				return err // Возвращаем ошибку наверх
			}
			_, err = fmt.Fprintf(writer, "User registered successfully")
			if err != nil {
				return err
			} // Убираем _, так как мы игнорируем возвращаемое значение
			return nil // Возвращаем nil, указывая на успешное выполнение
		}
		http.Error(writer, "Error checking existing user", http.StatusInternalServerError)
		return err // Возвращаем ошибку наверх
	}

	http.Error(writer, "User already exists", http.StatusBadRequest)
	return nil // Здесь нет ошибки, но пользователь уже существует
}
func (m *UserModel) LoginUser(login string, password string, writer http.ResponseWriter) error {
	rows := m.DB.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
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
		t, err := token.SignedString([]byte(os.Getenv("secretKey")))
		fmt.Println(os.Getenv("secretKey"))

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

func (m *UserModel) UpdateUser(login string, updateType string, writer http.ResponseWriter, request *http.Request) {
	rows := m.DB.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		if updateType == "password" {
			newPassword := request.FormValue("new-password")
			hashPass, _ := HashPassword(newPassword)
			_, err := m.DB.Exec("UPDATE users SET password = $1 WHERE username = $2", hashPass, login)
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
			err := m.DB.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&existingUser.Username)
			if err == nil {
				_, err := m.DB.Exec("UPDATE users SET username = $1 WHERE username = $2", newLogin, login)
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

func (m *UserModel) DeleteUser(login string, writer http.ResponseWriter) {
	rows := m.DB.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		_, err := m.DB.Exec("DELETE FROM users WHERE username = $1", login)
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

func (m *UserModel) GetAllUsers(writer http.ResponseWriter) {
	rows, err := m.DB.Query("SELECT * FROM users")
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
