package models

type User struct {
	Id       int
	Username string
	Password string
}

type Message struct {
	Id       int
	Sender   int
	Receiver int
	Text     string
	SentAt   string
}
