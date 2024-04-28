package models

type Message struct {
	Id          int    `json:"id"`
	SenderId    int    `json:"senderId"`
	ReceiverId  int    `json:"receiverId"`
	MessageText string `json:"messageText"`
	SentAt      string `json:"sentAt"`
	Readed      bool   `json:"readed"`
}

func (m *UserModel) SendMessage(senderId int, receiverId int, messageText string) error {
	_, err := m.DB.Exec("INSERT INTO messages (sender_id, receiver_id, message_text) VALUES ($1, $2, $3) returning *", senderId, receiverId, messageText)

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdateMessage(messageId int, messageText string) error {
	_, err := m.DB.Exec("UPDATE messages SET message_text = $1, read = $2 WHERE message_id = $3", messageText, false, messageId)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteMessage(message_id int, sender_id int) (bool, error) {
	result, err := m.DB.Exec("DELETE FROM messages WHERE message_id = $1 and sender_id = $2", message_id, sender_id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil // No rows affected, might return a custom error instead
	}
	return true, nil
}

func (m *UserModel) GetSendMessage(senderId int) ([]Message, error) {
	rows, err := m.DB.Query("SELECT * FROM messages WHERE sender_id = $1", senderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.MessageText, &message.SentAt, &message.Readed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (m *UserModel) GetReceivedMessage(receiverId int) ([]Message, error) {
	rows, err := m.DB.Query("SELECT * FROM messages WHERE receiver_id = $1 AND read = true", receiverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.MessageText, &message.SentAt, &message.Readed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (m *UserModel) GetUnreadedMessage(receiverId int) ([]Message, error) {
	rows, err := m.DB.Query("SELECT * FROM messages WHERE receiver_id = $1 AND read = false", receiverId)
	if err != nil {
		return nil, err
	}
	_, err = m.DB.Exec("UPDATE  messages SET read = true WHERE receiver_id = $1 and read = false", receiverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.MessageText, &message.SentAt, &message.Readed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
