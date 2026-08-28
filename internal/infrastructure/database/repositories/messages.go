package repositories

import (
	"log"

	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/database"

	"gorm.io/gorm"
)

type MessagesLoggingRepository struct {
	repository ports.MessagesRepository
}

func (adapter MessagesLoggingRepository) GetChatAllForUser(chatID int, userID int, offset int, limit int) dtos.OffsetResponse[entities.Message] {
	log.Printf("fetching chat all messages for user: chatId=%d, userId=%d, offset=%d, limit=%d", chatID, userID, offset, limit)
	messages := adapter.repository.GetChatAllForUser(chatID, userID, offset, limit)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingRepository) GetChatCursorAllForUser(chatID int, userID int, messageID int, aroundOffset int) dtos.OffsetResponse[entities.Message] {
	log.Printf("fetching chat all messages for user by cursor: chatId=%d, userId=%d, messageId=%d, aroundOffset=%d", chatID, userID, messageID, aroundOffset)
	messages := adapter.repository.GetChatCursorAllForUser(chatID, userID, messageID, aroundOffset)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingRepository) GetChatsLast(chatIds []int, userID int) []entities.Message {
	log.Printf("fetching last messages for chats: chatIds=%v, userId=%d", chatIds, userID)
	messages := adapter.repository.GetChatsLast(chatIds, userID)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingRepository) GetByIdForUser(messageID int, userID int) (*entities.Message, error) {
	log.Printf("fetching message by id for user: messageId=%d, userId=%d", messageID, userID)
	message, err := adapter.repository.GetByIdForUser(messageID, userID)
	if err != nil {
		log.Printf("error fetching message by id for user: %v", err)
		return message, err
	}

	log.Printf("fetchec message by id for user: %+v", message)
	return message, err
}

func (adapter MessagesLoggingRepository) GetByIdsForUser(messageIds []int, userID int) []entities.Message {
	log.Printf("fetching messages by ids for user: messageIds=%v, userId=%d", messageIds, userID)
	messages := adapter.repository.GetByIdsForUser(messageIds, userID)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingRepository) GetById(messageID int) (*entities.Message, error) {
	log.Printf("fetching message by id messageId=%v", messageID)
	message, err := adapter.repository.GetById(messageID)
	log.Printf("fetched message: %+v", message)
	return message, err
}

func (adapter MessagesLoggingRepository) Save(message entities.Message) (*entities.Message, error) {
	log.Printf("saving message: %+v", message)
	savedMessage, err := adapter.repository.Save(message)
	if err != nil {
		log.Printf("error saving message: %v", err)
		return savedMessage, err
	}

	log.Printf("saved message: %+v", savedMessage)
	return savedMessage, err
}

func (adapter MessagesLoggingRepository) Delete(message entities.Message) {
	log.Printf("deleting message: %+v", message)
	adapter.repository.Delete(message)
	log.Printf("message deleted")
}

type MessagesRepository struct {
	db gorm.DB
}

func (adapter MessagesRepository) getChatAllForUserTotal(chatID int, userID int) int {
	var count int64

	adapter.db.Model(&database.Message{}).Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatID, userID,
	).Count(&count)

	return int(count)
}

func (adapter MessagesRepository) GetChatAllForUser(chatID int, userID int, offset int, limit int) dtos.OffsetResponse[entities.Message] {
	var dbMessages []database.Message

	total := adapter.getChatAllForUserTotal(chatID, userID)

	adapter.db.Preload("database.Chat").Preload("Reactions").Preload("Voice").Preload("Circle").Preload("Attachments").Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatID, userID,
	).Order(
		"messages.created_at DESC NULLS LAST",
	).Offset(offset).Limit(limit).Find(&dbMessages)

	messagesModels := make([]entities.Message, 0, len(dbMessages))
	for _, dbMessage := range dbMessages {
		messagesModels = append(messagesModels, database.DbMessageToModel(dbMessage))
	}

	return dtos.NewOffsetResponse(
		offset,
		limit,
		total,
		messagesModels,
	)
}

func (adapter MessagesRepository) getMessageOffsetById(chatID int, userID int, messageID int) int {
	var offset int64

	adapter.db.Model(&database.Message{}).Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members) AND messages.id >= ?", chatID, userID, messageID,
	).Order("messages.created_at DESC NULLS LAST").Count(&offset)

	return int(offset)
}

func (adapter MessagesRepository) GetChatCursorAllForUser(chatID int, userID int, messageID int, aroundOffset int) dtos.OffsetResponse[entities.Message] {
	offset := adapter.getMessageOffsetById(chatID, userID, messageID)
	startOffset := max(offset-aroundOffset, 0)
	return adapter.GetChatAllForUser(chatID, userID, startOffset, aroundOffset*2)
}

func (adapter MessagesRepository) GetChatsLast(chatIds []int, userID int) []entities.Message {
	messages := make([]entities.Message, 0, len(chatIds))

	for _, chatID := range chatIds {
		var message database.Message

		adapter.db.Preload("database.Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
			"messages.chat_id = ? AND ? = ANY(chats.members)", chatID, userID,
		).Order("messages.created_at DESC NULLS LAST").Limit(1).First(&message)

		messageModel := database.DbMessageToModel(message)
		messages = append(messages, messageModel)
	}

	return messages
}

func (adapter MessagesRepository) GetById(messageID int) (*entities.Message, error) {
	var dbMessage database.Message

	result := adapter.db.Preload("database.Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id = ?", messageID,
	).First(&dbMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	messageModel := database.DbMessageToModel(dbMessage)
	return &messageModel, nil
}

func (adapter MessagesRepository) GetByIdForUser(messageID int, userID int) (*entities.Message, error) {
	var dbMessage database.Message

	result := adapter.db.Preload("database.Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id = ? AND ? = ANY(chats.members)", messageID, userID,
	).First(&dbMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	messageModel := database.DbMessageToModel(dbMessage)
	return &messageModel, nil
}

func (adapter MessagesRepository) GetByIdsForUser(messageIds []int, userID int) []entities.Message {
	var dbMessages []database.Message

	adapter.db.Preload("database.Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id IN ? AND ? = ANY(chats.members)", messageIds, userID,
	).Find(&dbMessages)

	modelMessages := make([]entities.Message, 0, len(dbMessages))
	for _, dbMessage := range dbMessages {
		modelMessages = append(modelMessages, database.DbMessageToModel(dbMessage))
	}

	return modelMessages
}

func (adapter MessagesRepository) getOrCreateReaction(reaction entities.MessageReaction) database.Reaction {
	var foundedReaction database.Reaction

	adapter.db.Where("user_id = ?", reaction.GetUserID()).First(&foundedReaction)

	if foundedReaction.UserId == uint(reaction.GetUserID()) {
		foundedReaction.Content = reaction.GetContent()
		adapter.db.Save(&foundedReaction)
		return foundedReaction
	}

	dbReaction := database.Reaction{
		UserId:  uint(reaction.GetUserID()),
		Content: reaction.GetContent(),
	}
	adapter.db.Save(&dbReaction)
	return dbReaction
}

func (adapter MessagesRepository) Save(message entities.Message) (*entities.Message, error) {
	circle := getOrCreateFile(message.GetCircle(), adapter.db)
	var circlePointer *database.SavedFile
	if circle.ID != 0 {
		circlePointer = &circle
	}

	voice := getOrCreateFile(message.GetVoice(), adapter.db)
	var voicePointer *database.SavedFile
	if voice.ID != 0 {
		voicePointer = &voice
	}

	attachments := make([]database.SavedFile, 0, len(message.GetAttachments()))
	for _, attachment := range message.GetAttachments() {
		attachments = append(attachments, getOrCreateFile(&attachment, adapter.db))
	}

	reactions := make([]database.Reaction, 0, len(message.GetReactions()))
	for _, reaction := range message.GetReactions() {
		reactions = append(reactions, adapter.getOrCreateReaction(reaction))
	}

	dbMessage := database.ModelToDbMessage(message, voicePointer, circlePointer, attachments, reactions)
	result := adapter.db.Save(&dbMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	savedMessage, err := adapter.GetById(int(dbMessage.ID))
	if err != nil {
		return nil, err
	}

	return savedMessage, nil
}

func (adapter MessagesRepository) Delete(message entities.Message) {
	adapter.db.Delete(&database.Message{ID: uint(message.GetID())})
}

func NewMessagesRepository(db gorm.DB) ports.MessagesRepository {
	return MessagesLoggingRepository{repository: MessagesRepository{db: db}}
}
