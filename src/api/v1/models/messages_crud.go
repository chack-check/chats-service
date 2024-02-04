package models

import (
	"fmt"
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

func (queries *MessagesQueries) GetConcrete(params GetConcreteMessageParams) (*Message, error) {
	var message Message

	if params.WithDeleted {
		database.DB.Where("chat_id = ? AND id = ?", params.ChatId, params.MessageId).First(&message)
	} else {
		database.DB.Where(
			"chat_id = ? AND NOT ? = ANY(deleted_for) AND id = ?", params.ChatId, params.UserId, params.MessageId,
		).First(&message)
	}

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

func (queries *MessagesQueries) GetAllInChat(params GetAllInChatParams) *[]Message {
	var messages []Message

	if params.WithDeleted {
		database.DB.Scopes(Paginate(params.Page, params.PerPage)).Preload(
			"Reactions",
		).Where(
			"chat_id = ?", params.ChatId,
		).Order(
			"created_at DESC",
		).Find(&messages)
	} else {
		database.DB.Scopes(Paginate(params.Page, params.PerPage)).Preload(
			"Reactions",
		).Where(
			"chat_id = ? AND NOT ? = ANY(deleted_for)", params.ChatId, params.UserId,
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
		database.DB.Model(&Message{}).Where("chat_id = ? AND NOT ? = ANY(deleted_for)", params.ChatId, params.UserId).Count(&count)
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

func (queries *MessagesQueries) Update(message *Message) {
	database.DB.Save(message)
}
