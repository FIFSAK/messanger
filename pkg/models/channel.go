package models

import "errors"

type Channel struct {
	ChanelId   int    `json:"chanelId"`
	OwnerId    int    `json:"ownerId"`
	ChanelName string `json:"chanelName"`
}

type ChannelUsers struct {
	ChanelUserId int `json:"chanelUserId"`
	ChanelId     int `json:"chanelId"`
	UserId       int `json:"userId"`
}

type ChannelMessages struct {
	ChanelMessageId int    `json:"chanelMessageId"`
	ChanelId        int    `json:"chanelId"`
	MessageText     string `json:"messageText"`
	SentAt          string `json:"sentAt"`
}

func (m *UserModel) CreateChannel(ownerId int, chanelName string) error {
	_, err := m.DB.Exec("INSERT INTO channel (owner_id, chanel_name) VALUES ($1, $2) returning *", ownerId, chanelName)

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdateChannel(chanelId int, chanelName string, ownerId int) error {
	_, err := m.DB.Exec("UPDATE channel SET chanel_name = $1 WHERE chanel_id = $2 and owner_id=$3", chanelName, chanelId, ownerId)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteChannel(chanelId int, ownerId int) (bool, error) {
	result, err := m.DB.Exec("DELETE FROM channel WHERE chanel_id = $1 and owner_id = $2", chanelId, ownerId)
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

func (m *UserModel) GetAllChannels() ([]Channel, error) {
	rows, err := m.DB.Query("SELECT * FROM channel")
	if err != nil {
		return nil, err
	}
	var channels []Channel
	for rows.Next() {
		var channel Channel
		err := rows.Scan(&channel.ChanelId, &channel.OwnerId, &channel.ChanelName)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (m *UserModel) FollowChannel(userId int, chanelId int) error {
	_, err := m.DB.Exec("INSERT INTO channelusers (chanel_id, user_id) VALUES ($1, $2) returning *", chanelId, userId)

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UnFollowChannel(userId int, chanelId int) error {
	_, err := m.DB.Exec("DELETE FROM channelusers WHERE chanel_id = $1 and user_id = $2", chanelId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) SendMessageToChannel(chanelId int, messageText string, ownerId int) error {
	rows, err := m.DB.Query("SELECT * FROM channel WHERE chanel_id = $1 and owner_id = $2", chanelId, ownerId)

	if err != nil {
		return err
	}
	if rows.Next() {
		_, err := m.DB.Exec("INSERT INTO channelmessages (chanel_id, message_text) VALUES ($1, $2) returning *", chanelId, messageText)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("You are not the owner of this channel")
}

func (m *UserModel) GetFollowedChannelsMessages(userId int) ([]ChannelMessages, error) {
	rows, err := m.DB.Query("select message_text, user_id, chanel_name from channelmessages as cm join channelusers as cu on cm.chanel_id = cu.chanel_id join channel as c on cu.chanel_id = c.chanel_id where cu.user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	var messages []ChannelMessages
	for rows.Next() {
		var message ChannelMessages
		err := rows.Scan(&message.ChanelMessageId, &message.ChanelId, &message.MessageText, &message.SentAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
