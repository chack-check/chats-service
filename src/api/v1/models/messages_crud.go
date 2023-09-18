package models

import "github.com/chack-check/chats-service/database"

type MessagesQueries struct{}

func (queries *MessagesQueries) GetAllInChat(page int, perPage int, chatId uint) *[]Message {
	var messages []Message
	database.DB.Scopes(Paginate(page, perPage)).Preload(
		"Reactions",
	).Where(
		"chat_id = ?", chatId,
	).Order(
		"created_at DESC",
	).Find(&messages)
	return &messages
}

func (queries *MessagesQueries) GetAllInChatCount(chatId uint) int64 {
	var count int64
	database.DB.Model(&Message{}).Where("chat_id = ?", chatId).Count(&count)
	return count
}

func (queries *MessagesQueries) Create(message *Message) error {
	result := database.DB.Create(message)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
