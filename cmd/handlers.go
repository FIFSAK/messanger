package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "messanger/pkg/auth"
	"messanger/pkg/models"
	"net/http"
	"strconv"
)

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "OK")
	if err != nil {
		return
	}
}
func RegisterHandler(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		login := request.FormValue("login")
		password := request.FormValue("password")

		err := userModel.RegisterUser(login, password, writer)
		if err != nil {
			// здесь вы уже обработали ошибку внутри RegisterUser
			return
		}
	}
}

func Login(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		login := request.FormValue("login")
		password := request.FormValue("password")

		err := userModel.LoginUser(login, password, writer)
		if err != nil {
			return
		}
	}
}
func Update(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		updateType := mux.Vars(request)["type"]
		login := payload["sub"].(string)
		userModel.UpdateUser(login, updateType, writer, request)
	}
}

func DeleteUserHandler(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		login := payload["sub"].(string)
		userModel.DeleteUser(login, writer)
	}
}

func GetAllUsersHandler(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		fmt.Println(payload["sub"])
		if !check {
			return
		}
		userModel.GetAllUsers(writer)
	}
}
func SendMessageHandler(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		senderId, ok := payload["id"].(float64)
		if !ok {
			http.Error(writer, "Invalid sender ID", http.StatusBadRequest)
			return
		}
		receiverId, _ := strconv.ParseInt(request.FormValue("receiver_id"), 10, 64)
		messageText := request.FormValue("message")
		err := userModel.SendMessage(int(senderId), int(receiverId), messageText)
		if err != nil {
			return
		}

	}
}
func UpdateMessageHandler(writer http.ResponseWriter, request *http.Request) {
	_, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
}
func DeleteMessageHandler(writer http.ResponseWriter, request *http.Request) {
	_, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
}
func NotificationsHandler(writer http.ResponseWriter, request *http.Request) {
	_, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
}
