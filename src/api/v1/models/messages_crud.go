package models

import (
	"fmt"
	"slices"

	"github.com/chack-check/chats-service/database"
	"github.com/lib/pq"
)

type MessagesQueries struct{}

func (queries *MessagesQueries) GetConcrete(chatId uint, messageId uint) (*Message, error) {
	var message Message
	database.DB.Where("chat_id = ?", chatId).Where("id = ?", messageId).First(&message)
	if message.ID == 0 {
		return &Message{}, fmt.Errorf("Chat with this ID doesn't exist")
	}

	return &message, nil
}

func (queries *MessagesQueries) Read(message *Message, userId uint) error {
	var readedBy []int32
	for _, v := range message.ReadedBy {
		readedBy = append(readedBy, int32(v))
	}

	if slices.Contains(readedBy, int32(userId)) {
		return nil
	}

	readedBy = append(readedBy, int32(userId))
	database.DB.Model(message).Update("readed_by", pq.Int32Array(readedBy))
	return nil
}

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

func (queries *MessagesQueries) AddReaction(userId uint, content string, message *Message) {
	if slices.IndexFunc(message.Reactions, func(r Reaction) bool { return r.UserId == userId }) >= 0 {
		return
	}

	newReaction := Reaction{
		MessageId: message.ID,
		UserId:    userId,
		Content:   content,
	}

	database.DB.Create(&newReaction)
	message.Reactions = append(message.Reactions, newReaction)
}
