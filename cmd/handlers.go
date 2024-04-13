package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "messanger/pkg/auth"
	"messanger/pkg/models"
	"net/http"
	"strconv"
	"strings"
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

		ordering := request.URL.Query().Get("ordering")
		if ordering == "" {
			ordering = "user_id" // Значение по умолчанию
		}
		page := request.URL.Query().Get("page")
		search := request.URL.Query().Get("search")
		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err != nil || pageInt < 1 {
			pageInt = 1
		}
		direction := "asc"
		if strings.Contains(ordering, "-") {
			direction = "desc"
			ordering = ordering[1:len(ordering)]
		}

		validOrderings := map[string]bool{"user_id": true, "username": true}
		if _, ok := validOrderings[ordering]; !ok {
			http.Error(writer, "Invalid ordering parameter", http.StatusBadRequest)
			return
		}

		fmt.Println("ordering "+ordering, "search "+search, "page "+page)
		userModel.GetAllUsers(writer, ordering, int(pageInt), direction, search)
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
		fmt.Fprintf(writer, "Message sent successfully")

	}
}
func UpdateMessageHandler(userModel *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		messageId, _ := strconv.ParseInt(request.FormValue("message_id"), 10, 64)
		messageText := request.FormValue("message_text")
		err := userModel.UpdateMessage(int(messageId), messageText)
		if err != nil {
			return
		}
		fmt.Fprintf(writer, "Message updated successfully")
	}
}
func DeleteMessageHandler(userModel *models.UserModel) http.HandlerFunc {
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
		messageId, _ := strconv.ParseInt(request.FormValue("message_id"), 10, 64)
		deleted, err := userModel.DeleteMessage(int(messageId), int(senderId))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if !deleted {
			http.Error(writer, "No message found to delete", http.StatusNotFound)
			return
		}
		fmt.Fprintf(writer, "Message deleted successfully")
	}
}

func GetSendMessageHandler(userModle *models.UserModel) http.HandlerFunc {
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
		messages, err := userModle.GetSendMessage(int(senderId))
		if err != nil {
			return
		}
		fmt.Println(messages)
		err = json.NewEncoder(writer).Encode(messages)
		if err != nil {
			return
		}
	}
}

func GetReceivedMessageHandler(userModle *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		receiverId, ok := payload["id"].(float64)
		if !ok {
			http.Error(writer, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		messages, err := userModle.GetReceivedMessage(int(receiverId))
		if err != nil {
			return
		}
		err = json.NewEncoder(writer).Encode(messages)
		if err != nil {
			return
		}
	}
}

func GetUnreadMessageHandler(userModle *models.UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		receiverId, ok := payload["id"].(float64)
		if !ok {
			http.Error(writer, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		messages, err := userModle.GetUnreadedMessage(int(receiverId))
		if err != nil {
			return
		}
		err = json.NewEncoder(writer).Encode(messages)
		if err != nil {
			return
		}
	}
}

func RefreshToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		receiverId, ok := payload["id"].(float64)
		if !ok {
			http.Error(writer, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		userName, _ := payload["sub"].(string)
		err := CreateToken(userName, int(receiverId), writer)

		fmt.Println(err)
	}
}
