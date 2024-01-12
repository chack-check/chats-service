package services

import (
	"log"
	"time"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/api/v1/utils"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/golang-jwt/jwt/v5"
)

func getMessageEventFromMessage(message *models.Message, chat *models.Chat) *rabbit.MessageEvent {
	members := []int{}
	attachments := []string{}
	mentioned := []int{}
	readedBy := []int{}

	for _, member := range chat.Members {
		members = append(members, int(member))
	}

	for _, attachment := range message.Attachments {
		attachments = append(attachments, string(attachment))
	}

	for _, ment := range message.Mentioned {
		mentioned = append(mentioned, int(ment))
	}

	for _, reader := range message.ReadedBy {
		readedBy = append(readedBy, int(reader))
	}

	return &rabbit.MessageEvent{
		Type:          "message",
		MessageId:     int(message.ID),
		IncludedUsers: members,
		ChatID:        int(message.ChatId),
		SenderID:      int(message.SenderId),
		MessageType:   message.Type,
		Content:       message.Content,
		VoiceURL:      message.VoiceURL,
		CircleURL:     message.CircleURL,
		Attachments:   attachments,
		ReplyToID:     int(message.ReplyToID),
		Mentioned:     mentioned,
		ReadedBy:      readedBy,
		Datetime:      message.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func getReadMessageEventFromMessage(message *models.Message, chat *models.Chat) *rabbit.ReadMessageEvent {
	members := []int{}
	readedBy := []int{}

	for _, member := range chat.Members {
		members = append(members, int(member))
	}

	for _, reader := range message.ReadedBy {
		readedBy = append(readedBy, int(reader))
	}

	return &rabbit.ReadMessageEvent{
		Type:          "message_readed",
		MessageId:     int(message.ID),
		ChatID:        int(message.ChatId),
		ReadedBy:      readedBy,
		IncludedUsers: members,
	}
}

type MessagesManager struct {
	MessagesQueries *models.MessagesQueries
}

func (manager *MessagesManager) GetChatAll(chatId uint, page int, perPage int) *schemas.PaginatedResponse[models.Message] {
	count := manager.MessagesQueries.GetAllInChatCount(chatId)
	messages := manager.MessagesQueries.GetAllInChat(page, perPage, chatId)
	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(count), *messages)
	return &paginatedResponse
}

func (manager *MessagesManager) getTextMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	log.Print("Creating text message")
	if err := utils.ValidateTextMessage(messageData); err != nil {
		return err
	}

	messageDataContent := *messageData.Content
	if len(messageDataContent) != 0 {
		message.Content = messageDataContent
	}

	if len(messageData.Attachments) != 0 {
		message.Attachments = messageData.Attachments
	}

	if messageData.ReplyToID != nil {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	if len(messageData.Mentioned) != 0 {
		var mentioned []int32
		for _, v := range messageData.Mentioned {
			mentioned = append(mentioned, int32(v))
		}
		message.Mentioned = mentioned
	}

	return nil
}

func (manager *MessagesManager) createVoiceMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	log.Print("Creating voice message")
	if err := utils.ValidateVoiceMessage(messageData); err != nil {
		return err
	}

	message.VoiceURL = *messageData.Voice

	if messageData.ReplyToID != nil && *messageData.ReplyToID != 0 {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	return nil
}

func (manager *MessagesManager) createCircleMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	log.Print("Creating circle message")
	if err := utils.ValidateCircleMessage(messageData); err != nil {
		return err
	}

	message.CircleURL = *messageData.Circle

	if *messageData.ReplyToID != 0 {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	return nil
}

func (manager *MessagesManager) CreateMessage(messageData *model.CreateMessageRequest, chat *models.Chat, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		ChatId:   chat.ID,
		Type:     messageData.Type.String(),
		SenderId: uint(tokenSubject.UserId),
	}

	if messageData.Type.String() == "text" {
		if err = manager.getTextMessage(message, messageData); err != nil {
			return nil, err
		}
	} else if messageData.Type.String() == "voice" {
		if err = manager.createVoiceMessage(message, messageData); err != nil {
			return nil, err
		}
	} else if messageData.Type.String() == "circle" {
		if err = manager.createCircleMessage(message, messageData); err != nil {
			return nil, err
		}
	}

	if err = manager.MessagesQueries.Create(message); err != nil {
		return nil, err
	}

	sendingIds32 := []int32{}
	for _, v := range chat.Members {
		sendingIds32 = append(sendingIds32, int32(v))
	}

	messageEvent := getMessageEventFromMessage(message, chat)
	err = rabbit.EventsRabbitConnection.SendEvent(messageEvent)
	log.Printf("Sended message to rabbitmq")

	if err != nil {
		log.Printf("Error when publishing message event in queue: %v", err)
	}

	return message, nil
}

func (manager *MessagesManager) Read(chat *models.Chat, messageId uint, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcrete(chat.ID, messageId)
	if err != nil {
		return nil, err
	}

	manager.MessagesQueries.Read(message, uint(tokenSubject.UserId))

	readMessageEvent := getReadMessageEventFromMessage(message, chat)
	err = rabbit.EventsRabbitConnection.SendEvent(readMessageEvent)

	if err != nil {
		log.Printf("Error when publishing message readed event in queue: %v", err)
	}

	log.Printf("Sended message readed event to rabbitmq")

	return message, nil
}

func (manager *MessagesManager) ReactMessage(userId uint, chatId uint, messageId uint, content string) (*models.Message, error) {
	message, err := manager.MessagesQueries.GetConcrete(chatId, messageId)
	if err != nil {
		return &models.Message{}, err
	}

	manager.MessagesQueries.AddReaction(userId, content, message)

	return message, nil
}

func NewMessagesManager() *MessagesManager {
	return &MessagesManager{
		MessagesQueries: &models.MessagesQueries{},
	}
}
