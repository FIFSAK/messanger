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
	_, err := m.DB.Exec("INSERT INTO messages (sender_id, receiver_id, message_text) VALUES ($1, $2, $3)", senderId, receiverId, messageText)

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdateMessage(senderId int, receiverId int, messageText string) error {
	_, err := m.DB.Exec("UPDATE messages SET message_text = $1 WHERE sender_id = $2 AND receiver_id = $3 AND readed = $4", messageText, senderId, receiverId, false)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteMessage(senderId int, receiverId int, message_id int) error {
	_, err := m.DB.Exec("DELETE FROM messages WHERE sender_id = $1 AND receiver_id = $2 AND id = $3", senderId, receiverId, message_id)
	if err != nil {
		return err
	}
	return nil
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
	rows, err := m.DB.Query("SELECT * FROM messages WHERE receiver_id = $1 AND readed = true", receiverId)
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
	rows, err := m.DB.Query("SELECT * FROM messages WHERE receiver_id = $1 AND readed = false", receiverId)
	if err != nil {
		return nil, err
	}
	_, err = m.DB.Exec("UPDATE  messages SET readed = true WHERE receiver_id = $1 and readed = false", receiverId)
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
