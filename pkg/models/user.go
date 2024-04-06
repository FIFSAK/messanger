package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	. "messanger/pkg/auth"
	"net/http"
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

func (m *UserModel) RegisterUser(login string, password string, writer http.ResponseWriter) error {
	var username string
	err := m.DB.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			hashPass, err := HashPassword(password)
			if err != nil {
				http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
				return err
			}
			_, err = m.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", login, hashPass)
			if err != nil {
				http.Error(writer, "Failed to insert user into database", http.StatusInternalServerError)
				return err
			}
			_, err = fmt.Fprintf(writer, "User registered successfully")
			if err != nil {
				return err
			}
			return nil
		}
		http.Error(writer, "Error checking existing user", http.StatusInternalServerError)
		return err
	}

	http.Error(writer, "User already exists", http.StatusBadRequest)
	return nil
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
		err := CreateToken(user.Username, user.Id, writer)
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

func (m *UserModel) GetAllUsers(writer http.ResponseWriter, ordering string, page int, direction string, search string) {
	limit := 2
	offset := limit * (page - 1)
	var totalUsers int
	countQuery := "SELECT COUNT(*) FROM users"
	err := m.DB.QueryRow(countQuery).Scan(&totalUsers)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Failed to fetch user count", http.StatusInternalServerError)
		return
	}

	// Расчет максимального количества страниц
	maxPages := int(math.Ceil(float64(totalUsers) / float64(limit)))

	if page < 0 && maxPages < page {
		http.Error(writer, "Page parameter out of range", http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("SELECT * FROM users WHERE username LIKE $3 ORDER BY %s %s LIMIT $1 OFFSET $2", ordering, direction)

	searchPattern := "%" + search + "%"
	rows, err := m.DB.Query(query, limit, offset, searchPattern)

	if err != nil {
		fmt.Println(err)
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
	type UserResponse struct {
		Users    []User
		Total    int
		MaxPages int
		Page     int
	}
	response := UserResponse{
		Users:    users,
		Total:    totalUsers,
		MaxPages: maxPages,
		Page:     page,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
	}
}

//func (m *UserModel) GetUser(login string, writer http.ResponseWriter) (id int, username string) {
//	rows := m.DB.QueryRow("SELECT * FROM users WHERE username = $1", login)
//	user := User{}
//	_ = rows.Scan(&user.Id, &user.Username)
//	return user.Id, user.
//}
