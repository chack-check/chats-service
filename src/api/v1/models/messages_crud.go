package models

import (
	"fmt"
	"log"
	"slices"

	"github.com/chack-check/chats-service/database"
	"github.com/lib/pq"
)

type GetConcreteMessageParams struct {
	ChatId      uint
	MessageId   uint
	WithDeleted bool `default:"false"`
	UserId      uint `default:"0"`
}

type GetAllInChatParams struct {
	Page        int
	PerPage     int
	ChatId      uint
	WithDeleted bool `default:"false"`
	UserId      uint `default:"0"`
}

type GetAllInChatCountParams struct {
	ChatId      uint
	WithDeleted bool `default:"false"`
	UserId      uint `default:"0"`
}

type MessagesQueries struct{}

func (queries *MessagesQueries) GetConcreteById(messageId int, userId int) (*Message, error) {
	var message Message

	database.DB.Where(
		"(deleted_for IS NULL OR NOT ? = ANY(deleted_for)) AND id = ?", userId, messageId,
	).First(&message)

	if message.ID == 0 {
		return nil, fmt.Errorf("Message with this id doesn't exist")
	}

	return &message, nil
}

func (queries *MessagesQueries) GetByIds(messageIds []int, userId int) ([]*Message, error) {
	var messages []*Message

	log.Printf("Getting messages with ids %v and for user id %d", messageIds, userId)

	database.DB.Where(
		"(deleted_for IS NULL OR NOT ? = ANY(deleted_for)) AND id IN ?", userId, messageIds,
	).Find(&messages)

	return messages, nil
}

func (queries *MessagesQueries) GetConcrete(params GetConcreteMessageParams) (*Message, error) {
	var message Message

	if params.WithDeleted {
		database.DB.Preload("Reactions").Where("chat_id = ? AND id = ?", params.ChatId, params.MessageId).First(&message)
	} else {
		database.DB.Preload("Reactions").Where(
			"chat_id = ? AND (deleted_for IS NULL OR NOT ? = ANY(deleted_for)) AND id = ?", params.ChatId, params.UserId, params.MessageId,
		).First(&message)
	}

	if message.ID == 0 {
		return &Message{}, fmt.Errorf("Message with this ID doesn't exist")
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

func (queries *MessagesQueries) GetAllInChat(params GetAllInChatParams) *[]Message {
	var messages []Message

	if params.WithDeleted {
		log.Printf("Fetching chat messages with deleted messages. Chat id: %d", params.ChatId)
		database.DB.Scopes(Paginate(params.Page, params.PerPage)).Preload(
			"Reactions",
		).Where(
			"chat_id = ?", params.ChatId,
		).Order(
			"created_at DESC",
		).Find(&messages)
	} else {
		log.Printf("Fetching chat messages without deleted messages. Chat id: %d. User id: %d", params.ChatId, params.UserId)
		database.DB.Scopes(Paginate(params.Page, params.PerPage)).Preload(
			"Reactions",
		).Where(
			"chat_id = ? AND (deleted_for IS NULL OR NOT ? = ANY(deleted_for))", params.ChatId, params.UserId,
		).Order(
			"created_at DESC",
		).Find(&messages)
	}

	return &messages
}

func (queries *MessagesQueries) GetAllInChatCount(params GetAllInChatCountParams) int64 {
	var count int64

	if params.WithDeleted {
		database.DB.Model(&Message{}).Where("chat_id = ?", params.ChatId).Count(&count)
	} else {
		database.DB.Model(&Message{}).Where("chat_id = ? AND (deleted_for IS NULL OR NOT ? = ANY(deleted_for))", params.ChatId, params.UserId).Count(&count)
	}

	return count
}

func (queries *MessagesQueries) Create(message *Message) error {
	result := database.DB.Create(message)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (queris *MessagesQueries) DeleteMessage(message *Message, deleteFor []int32) error {
	deleteForArray := pq.Int32Array{}
	for _, deleteForId := range deleteFor {
		deleteForArray = append(deleteForArray, deleteForId)
	}

	message.DeletedFor = deleteForArray
	database.DB.Save(message)
	return nil
}

func (queries *MessagesQueries) AddOrGetReaction(userId uint, content string, message *Message) bool {
	if slices.IndexFunc(message.Reactions, func(r Reaction) bool { return r.UserId == userId }) >= 0 {
		return false
	}

	newReaction := Reaction{
		MessageId: message.ID,
		UserId:    userId,
		Content:   content,
	}

	database.DB.Create(&newReaction)
	message.Reactions = append(message.Reactions, newReaction)
	return true
}

func (queries *MessagesQueries) DeleteReaction(userId int, message *Message) {
	database.DB.Delete(&Reaction{}, "message_id = ? AND user_id = ?", message.ID, userId)
}

func (queries *MessagesQueries) GetLastForChatId(chatId int, userId int) *Message {
	var message Message
	database.DB.Model(
		&Message{},
	).Joins(
		"JOIN chats ON chats.id = messages.chat_id",
	).Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatId, userId,
	).Order(
		"messages.created_at DESC",
	).Limit(1).Find(&message)

	log.Printf("Received last messages for chat id: %d; and user id: %d. Last message: %v", chatId, userId, message)

	if message.ID != 0 {
		return &message
	}

	return nil
}

func (queries *MessagesQueries) Update(message *Message) {
	database.DB.Save(message)
}
