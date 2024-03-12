package models

import (
	"time"
)

type Message struct {
	id          int
	senderId    int
	receiverId  int
	messageText string
	sentAt      string
	readed      bool
}

func (m *UserModel) SendMessage(senderId int, receiverId int, messageText string) error {
	_, err := m.DB.Exec("INSERT INTO messages (sender_id, receiver_id, message_text) VALUES ($1, $2, $3)", senderId, receiverId, messageText)

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdateMessage(id int, messageText string) error {
	_, err := m.DB.Exec("UPDATE messages SET message_text = $1, readed=$2 WHERE id = $3", messageText, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteMessage(id int) error {
	_, err := m.DB.Exec("DELETE FROM messages WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) getSendMessage(userId int) ([]Message, error) {
	rows, err := m.DB.Query("SELECT * FROM messages WHERE sender_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.id, &message.senderId, &message.receiverId, &message.messageText, &message.sentAt, &message.readed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (m *UserModel) getReceivedMessage(userId int) ([]Message, error) {
	rows, err := m.DB.Query("SELECT * FROM messages WHERE receiver_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.id, &message.senderId, &message.receiverId, &message.messageText, &message.sentAt, &message.readed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (m *UserModel) getUnreadedMessage(userId int) ([]Message, error) {
	rows, err := m.DB.Query("SELECT * FROM messages WHERE receiver_id = $1 AND readed = false", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.id, &message.senderId, &message.receiverId, &message.messageText, &message.sentAt, &message.readed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
