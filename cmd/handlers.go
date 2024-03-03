package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "messanger/pkg/auth"
	. "messanger/pkg/models"
	"net/http"
)

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "OK")
	if err != nil {
		return
	}
}
func Register(writer http.ResponseWriter, request *http.Request) {

	login := request.FormValue("login")
	password := request.FormValue("password")
	user := User{}
	err := user.RegisterUser(login, password, writer)
	if err != nil {
		return
	}

}
func Login(writer http.ResponseWriter, request *http.Request) {

	login := request.FormValue("login")
	password := request.FormValue("password")

	user := User{}
	err := user.LoginUser(login, password, writer)
	if err != nil {
		return
	}
}
func Update(writer http.ResponseWriter, request *http.Request) {
	payload, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
	updateType := mux.Vars(request)["type"]
	login := payload["sub"].(string)
	user := User{}
	user.UpdateUser(login, updateType, writer, request)
}

func DeleteUserHandler(writer http.ResponseWriter, request *http.Request) {
	payload, check := JwtPayloadFromRequest(writer, request)
	if !check {
		return
	}
	login := payload["sub"].(string)
	user := User{}
	user.DeleteUser(login, writer)
}

func GetAllUsersHandler(writer http.ResponseWriter, request *http.Request) {
	payload, check := JwtPayloadFromRequest(writer, request)
	fmt.Println(payload["sub"])
	if !check {
		return
	}
	user := User{}
	user.GetAllUsers(writer)
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
