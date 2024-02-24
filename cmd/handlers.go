package main

import (
	"fmt"
	"net/http"
)

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "OK")
	if err != nil {
		return
	}
}
func Register(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Register")
	if err != nil {
		return
	}
}
func Login(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Login")
	if err != nil {
		return
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
