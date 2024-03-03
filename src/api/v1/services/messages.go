package services

import (
	"fmt"
	"log"
	"slices"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/api/v1/utils"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type DeleteForOptions string

const (
	DeleteForMe  DeleteForOptions = "me"
	DeleteForAll DeleteForOptions = "all"
)

type MessagesManager struct {
	MessagesQueries *models.MessagesQueries
	ChatsQueries    *models.ChatsQueries
}

func (manager *MessagesManager) GetConcrete(messageId int, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcreteById(messageId, tokenSubject.UserId)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (manager *MessagesManager) GetByIds(messageIds []int, token *jwt.Token) ([]*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	messages, err := manager.MessagesQueries.GetByIds(messageIds, tokenSubject.UserId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (manager *MessagesManager) GetChatAll(token *jwt.Token, chatId uint, page int, perPage int) (*schemas.PaginatedResponse[models.Message], error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	count := manager.MessagesQueries.GetAllInChatCount(models.GetAllInChatCountParams{ChatId: chatId, UserId: uint(tokenSubject.UserId)})
	messages := manager.MessagesQueries.GetAllInChat(models.GetAllInChatParams{ChatId: chatId, Page: page, PerPage: perPage})
	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(count), *messages)
	return &paginatedResponse, nil
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

	var includedUsers []int
	for _, member := range chat.Members {
		includedUsers = append(includedUsers, int(member))
	}

	messageEvent, err := rabbit.NewSystemEvent("message_created", includedUsers, message)
	if err != nil {
		return nil, err
	}

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

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: chat.ID, MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return nil, err
	}

	manager.MessagesQueries.Read(message, uint(tokenSubject.UserId))

	var included_users []int
	for _, user := range chat.Members {
		included_users = append(included_users, int(user))
	}

	readMessageEvent, err := rabbit.NewSystemEvent("message_readed", included_users, message)
	if err != nil {
		return nil, err
	}

	err = rabbit.EventsRabbitConnection.SendEvent(readMessageEvent)

	if err != nil {
		log.Printf("Error when publishing message readed event in queue: %v", err)
	}

	log.Printf("Sended message readed event to rabbitmq")

	return message, nil
}

func (manager *MessagesManager) ReactMessage(token *jwt.Token, chatId uint, messageId uint, content string) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	chat, err := manager.ChatsQueries.GetWithMember(uint(chatId), uint(tokenSubject.UserId))
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: chatId, MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return &models.Message{}, err
	}

	manager.MessagesQueries.AddReaction(uint(tokenSubject.UserId), content, message)

	var included_users []int
	for _, user := range chat.Members {
		included_users = append(included_users, int(user))
	}

	reactMessageEvent, err := rabbit.NewSystemEvent("message_reacted", included_users, message)
	if err != nil {
		return nil, err
	}

	err = rabbit.EventsRabbitConnection.SendEvent(reactMessageEvent)

	if err != nil {
		log.Printf("Error when publishing message reacted event in queue: %v", err)
	}

	log.Printf("Sended message reacted event to rabbitmq")

	return message, nil
}

func (manager *MessagesManager) validateCanUserDeleteMessage(userId uint, chat *models.Chat, message *models.Message) bool {
	return slices.Contains(chat.Members, int64(userId))
}

func (manager *MessagesManager) DeleteMessage(token *jwt.Token, chatId uint, messageId uint, deleteFor DeleteForOptions) error {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return err
	}

	chat, err := manager.ChatsQueries.GetConcrete(uint(tokenSubject.UserId), chatId)
	if err != nil {
		return err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: chatId, MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return err
	}

	if !manager.validateCanUserDeleteMessage(uint(tokenSubject.UserId), chat, message) {
		return fmt.Errorf("the user with id %d can't delete message with id %d", tokenSubject.UserId, messageId)
	}

	deleteForArray := []int32{}
	if deleteFor == DeleteForMe {
		deleteForArray = append(deleteForArray, int32(tokenSubject.UserId))
	} else {
		for _, member := range chat.Members {
			deleteForArray = append(deleteForArray, int32(member))
		}
	}

	err = manager.MessagesQueries.DeleteMessage(message, deleteForArray)
	return err
}

func (manager *MessagesManager) Update(chat *models.Chat, messageId uint, updateData model.ChangeMessageRequest, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: chat.ID, MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return nil, err
	}

	message.Content = *updateData.Content
	attachments := pq.StringArray{}
	for _, attachment := range updateData.Attachments {
		attachments = append(attachments, attachment)
	}

	mentioned := pq.Int32Array{}
	for _, ment := range updateData.Mentioned {
		mentioned = append(mentioned, int32(*ment))
	}

	message.Attachments = attachments
	message.Mentioned = mentioned
	manager.MessagesQueries.Update(message)
	return message, nil
}

func (manager *MessagesManager) GetLastByChatIds(chatIds []int, userId int) []*models.Message {
	var last_messages []*models.Message

	for _, chat_id := range chatIds {
		message := manager.MessagesQueries.GetLastForChatId(chat_id, userId)
		if message != nil {
			last_messages = append(last_messages, message)
		}
	}

	return last_messages
}

func NewMessagesManager() *MessagesManager {
	return &MessagesManager{
		MessagesQueries: &models.MessagesQueries{},
		ChatsQueries:    &models.ChatsQueries{},
	}
}
