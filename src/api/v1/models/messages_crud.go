package models

import (
	"fmt"
	"log"
	"slices"

	"github.com/chack-check/chats-service/database"
	"github.com/getsentry/sentry-go"
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
	log.Printf("Fetching concrete message by id: messageId=%d userId=%d", messageId, userId)

	var message Message

	database.DB.Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Where(
		"(deleted_for IS NULL OR NOT ? = ANY(deleted_for)) AND id = ?", userId, messageId,
	).First(&message)

	log.Printf("Fetched message = %+v", message)
	if message.ID == 0 {
		log.Printf("Fetched message id == 0: %+v", message)
		sentry.CaptureException(fmt.Errorf("message with id %d for user id %d doesn't exist", messageId, userId))
		return nil, fmt.Errorf("Message with this id doesn't exist")
	}

	return &message, nil
}

func (queries *MessagesQueries) GetByIds(messageIds []int, userId int) ([]*Message, error) {
	log.Printf("Fetching messages by ids: %v for user id: %d", messageIds, userId)
	var messages []*Message

	database.DB.Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Where(
		"(deleted_for IS NULL OR NOT ? = ANY(deleted_for)) AND id IN ?", userId, messageIds,
	).Find(&messages)

	log.Printf("Fetched messages by ids %v: %v", messageIds, messages)

	return messages, nil
}

func (queries *MessagesQueries) GetConcrete(params GetConcreteMessageParams) (*Message, error) {
	log.Printf("Fetching concrete message with params: %+v", params)

	var message Message

	if params.WithDeleted {
		log.Printf("Fetching concrete message include deleted")
		database.DB.Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Where("chat_id = ? AND id = ?", params.ChatId, params.MessageId).First(&message)
	} else {
		log.Printf("Fetching concrete message exclude deleted")
		database.DB.Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Where(
			"chat_id = ? AND (deleted_for IS NULL OR NOT ? = ANY(deleted_for)) AND id = ?", params.ChatId, params.UserId, params.MessageId,
		).First(&message)
	}

	log.Printf("Fetched message: %+v", message)
	if message.ID == 0 {
		sentry.CaptureException(fmt.Errorf("message with params %+v doesn't exist", params))
		log.Printf("Not found message with params %+v: %+v", params, message)
		return &Message{}, fmt.Errorf("Message with this ID doesn't exist")
	}

	return &message, nil
}

func (queries *MessagesQueries) Read(message *Message, userId uint) error {
	log.Printf("Read message by user id = %d: %+v", userId, message)
	var readedBy []int32
	for _, v := range message.ReadedBy {
		readedBy = append(readedBy, int32(v))
	}

	log.Printf("Reading message readed by: %v", readedBy)
	if slices.Contains(readedBy, int32(userId)) {
		log.Printf("Reading message id = %d already readed by user with id %d", message.ID, userId)
		return nil
	}

	readedBy = append(readedBy, int32(userId))
	log.Printf("Updating message id = %d new readed by: %v", message.ID, readedBy)
	database.DB.Model(message).Update("readed_by", pq.Int32Array(readedBy))
	return nil
}

func (queries *MessagesQueries) GetAllInChat(params GetAllInChatParams) *[]Message {
	log.Printf("Fetching all messages in chat with params %+v", params)
	var messages []Message

	if params.WithDeleted {
		log.Printf("Fetching all messages in chat include deleted messages")
		database.DB.Scopes(Paginate(params.Page, params.PerPage)).Preload(
			"Reactions",
		).Preload("Voice").Preload("Circle").Preload("Attachments").Where(
			"chat_id = ?", params.ChatId,
		).Order(
			"created_at DESC",
		).Find(&messages)
	} else {
		log.Printf("Fetching all messages in chat exclude deleted messages")
		database.DB.Scopes(Paginate(params.Page, params.PerPage)).Preload(
			"Reactions",
		).Preload("Voice").Preload("Circle").Preload("Attachments").Where(
			"chat_id = ? AND (deleted_for IS NULL OR NOT ? = ANY(deleted_for))", params.ChatId, params.UserId,
		).Order(
			"created_at DESC",
		).Find(&messages)
	}

	return &messages
}

func (queries *MessagesQueries) GetAllInChatCount(params GetAllInChatCountParams) int64 {
	log.Printf("Fetching all messages in chat count with params %+v", params)
	var count int64

	if params.WithDeleted {
		database.DB.Model(&Message{}).Where("chat_id = ?", params.ChatId).Count(&count)
	} else {
		database.DB.Model(&Message{}).Where("chat_id = ? AND (deleted_for IS NULL OR NOT ? = ANY(deleted_for))", params.ChatId, params.UserId).Count(&count)
	}

	log.Printf("Fetched all messages in chat count: %d", count)
	return count
}

func (queries *MessagesQueries) Create(message *Message) error {
	log.Printf("Creating new message: %+v", message)

	result := database.DB.Create(message)

	if result.Error != nil {
		log.Printf("Creating message error: %s", result.Error)
		sentry.CaptureException(result.Error)
		return result.Error
	}

	return nil
}

func (queris *MessagesQueries) DeleteMessage(message *Message, deleteFor []int32) error {
	log.Printf("Deleting message %+v for %v", message, deleteFor)
	deleteForArray := pq.Int32Array{}
	for _, deleteForId := range deleteFor {
		deleteForArray = append(deleteForArray, deleteForId)
	}

	message.DeletedFor = deleteForArray
	database.DB.Save(message)
	return nil
}

func (queries *MessagesQueries) AddOrGetReaction(userId uint, content string, message *Message) bool {
	log.Printf("Adding reaction by user id = %d with content = %s for message = %+v", userId, content, message)
	if slices.IndexFunc(message.Reactions, func(r Reaction) bool { return r.UserId == userId }) >= 0 {
		log.Printf("Message already have reaction by user id = %d", userId)
		return false
	}

	newReaction := Reaction{
		MessageId: message.ID,
		UserId:    userId,
		Content:   content,
	}

	log.Printf("Creating new reaction for message: %+v", newReaction)
	database.DB.Create(&newReaction)
	message.Reactions = append(message.Reactions, newReaction)
	log.Printf("Message with new reactions: %+v", message)
	return true
}

func (queries *MessagesQueries) DeleteReaction(userId int, message *Message) {
	log.Printf("Deleting reaction by user id = %d for message = %+v", userId, message)
	database.DB.Delete(&Reaction{}, "message_id = ? AND user_id = ?", message.ID, userId)
}

func (queries *MessagesQueries) GetLastForChatId(chatId int, userId int) *Message {
	log.Printf("Fetching last message for chat id = %d and user id = %d", chatId, userId)

	var message Message
	database.DB.Model(
		&Message{},
	).Joins(
		"JOIN chats ON chats.id = messages.chat_id",
	).Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatId, userId,
	).Order(
		"messages.created_at DESC",
	).Limit(1).Find(&message)

	log.Printf("Received last messages for chat id: %d; and user id: %d. Last message: %+v", chatId, userId, message)

	if message.ID != 0 {
		return &message
	}

	log.Printf("Error fetching last message for chat id = %d, user id = %d: %+v", chatId, userId, message)
	sentry.CaptureException(fmt.Errorf("fetching last chat message for chat id = %d user id = %d: %+v", chatId, userId, message))
	return nil
}

func (queries *MessagesQueries) Update(message *Message) {
	log.Printf("Saving message: %+v", message)
	log.Printf("Clearing attachments for message")
	database.DB.Model(&Message{ID: message.ID}).Association("Attachments").Clear()
	database.DB.Save(message)
}
