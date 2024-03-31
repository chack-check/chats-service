package services

import (
	"fmt"
	"log"
	"slices"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/repositories"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/api/v1/utils"
	"github.com/chack-check/chats-service/generic_factories"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/getsentry/sentry-go"
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
	ChatsRepository *repositories.ChatsRepository
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
	messages := manager.MessagesQueries.GetAllInChat(models.GetAllInChatParams{ChatId: chatId, Page: page, PerPage: perPage, UserId: uint(tokenSubject.UserId)})
	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(count), *messages)
	return &paginatedResponse, nil
}

func (manager *MessagesManager) createTextMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	log.Printf("Constructing text message: %+v. Message data = %+v", message, messageData)
	if err := utils.ValidateTextMessage(messageData); err != nil {
		return err
	}

	messageDataContent := *messageData.Content
	if len(messageDataContent) != 0 {
		message.Content = messageDataContent
	}

	var message_attachments []models.SavedFile
	if len(messageData.Attachments) != 0 {
		log.Printf("Saving message attachments: %+v", messageData.Attachments)
		for _, attachment := range messageData.Attachments {
			utils.ValidateUploadingFile(*attachment, "file_in_chat")
			saved_file := generic_factories.UploadingFileToDbFile(*attachment)
			message_attachments = append(message_attachments, saved_file)
		}
	}
	message.Attachments = message_attachments

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

	fmt.Printf("Constructed text message: %+v", message)
	return nil
}

func (manager *MessagesManager) createVoiceMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	log.Printf("Constructing voice message %+v. Message data = %+v", message, messageData)
	if err := utils.ValidateVoiceMessage(messageData); err != nil {
		return err
	}

	if messageData.Voice == nil {
		log.Printf("Trying to create voice message without voice")
		sentry.CaptureException(fmt.Errorf("trying to create voice message without voice"))
		return fmt.Errorf("you can't create voice message without voice field")
	}

	message.Voice = generic_factories.UploadingFileToDbFile(*messageData.Voice)

	if messageData.ReplyToID != nil && *messageData.ReplyToID != 0 {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	log.Printf("Constructed voice message: %+v", message)
	return nil
}

func (manager *MessagesManager) createCircleMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	log.Printf("Constructing circle message %+v. Message data = %+v", message, messageData)
	if err := utils.ValidateCircleMessage(messageData); err != nil {
		return err
	}

	if messageData.Circle == nil {
		sentry.CaptureException(fmt.Errorf("trying to create circle message without circle"))
		return fmt.Errorf("you can't create circle message without circle field")
	}

	message.Circle = generic_factories.UploadingFileToDbFile(*messageData.Circle)

	if messageData.ReplyToID != nil && *messageData.ReplyToID != 0 {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	log.Printf("Constructed circle message: %+v", message)
	return nil
}

func (manager *MessagesManager) sendMessageEvent(message *models.Message, chat *dtos.ChatDto, eventType string, includedUsers *[]int) error {
	log.Printf("Sending message event to rabbitmq. Message = %+v. Chat = %+v. Event type = %s. IncludedUsers = %v", message, chat, eventType, includedUsers)
	if includedUsers == nil {
		includedUsers = &([]int{})
		for _, member := range chat.Members {
			new_included_users := append(*includedUsers, int(member))
			includedUsers = &new_included_users
		}
	}

	messageEvent, err := rabbit.NewSystemEvent(eventType, *includedUsers, message)
	if err != nil {
		log.Printf("Error constructing message event: %v. Message event: %+v", err, messageEvent)
		sentry.CaptureException(err)
		return err
	}

	log.Printf("Sending rabbitmq event %+v", messageEvent)
	err = rabbit.EventsRabbitConnection.SendEvent(messageEvent)
	if err != nil {
		log.Printf("Error sending event with type %s to rabbitmq", eventType)
		sentry.CaptureException(err)
	}

	return nil
}

func (manager *MessagesManager) CreateMessage(messageData *model.CreateMessageRequest, chat *dtos.ChatDto, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	for _, attachment := range messageData.Attachments {
		if err := utils.ValidateUploadingFile(*attachment, "file_in_chat"); err != nil {
			return nil, err
		}
	}
	if messageData.Voice != nil {
		if err := utils.ValidateUploadingFile(*messageData.Voice, "voice"); err != nil {
			return nil, err
		}
	}
	if messageData.Circle != nil {
		if err := utils.ValidateUploadingFile(*messageData.Circle, "circle"); err != nil {
			return nil, err
		}
	}

	message := &models.Message{
		ChatId:   uint(chat.Id),
		Type:     messageData.Type.String(),
		SenderId: uint(tokenSubject.UserId),
	}

	if messageData.Type.String() == "text" {
		if err = manager.createTextMessage(message, messageData); err != nil {
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

	manager.sendMessageEvent(message, chat, "message_created", nil)
	return message, nil
}

func (manager *MessagesManager) Read(chat *dtos.ChatDto, messageId uint, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: uint(chat.Id), MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return nil, err
	}

	manager.MessagesQueries.Read(message, uint(tokenSubject.UserId))

	manager.sendMessageEvent(message, chat, "message_readed", nil)

	return message, nil
}

func (manager *MessagesManager) ReactMessage(token *jwt.Token, chatId uint, messageId uint, content string) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	chat, err := manager.ChatsRepository.GetWithMember(uint(chatId), uint(tokenSubject.UserId))
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: chatId, MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return &models.Message{}, err
	}

	created := manager.MessagesQueries.AddOrGetReaction(uint(tokenSubject.UserId), content, message)
	if !created {
		return message, nil
	}

	manager.sendMessageEvent(message, chat, "message_reacted", nil)

	return message, nil
}

func (manager *MessagesManager) validateCanUserDeleteMessage(userId uint, chat *dtos.ChatDto, message *models.Message) bool {
	return slices.Contains(chat.Members, int(userId))
}

func (manager *MessagesManager) DeleteMessage(token *jwt.Token, chatId uint, messageId uint, deleteFor DeleteForOptions) error {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return err
	}

	chat, err := manager.ChatsRepository.GetWithMember(chatId, uint(tokenSubject.UserId))
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

	var deleteForArray []int32
	var includedUsersArray []int
	log.Printf("Deleting message for: %v", deleteFor)
	if deleteFor == DeleteForMe {
		deleteForArray = append(deleteForArray, int32(tokenSubject.UserId))
		includedUsersArray = append(includedUsersArray, tokenSubject.UserId)
	} else {
		for _, member := range chat.Members {
			deleteForArray = append(deleteForArray, int32(member))
			includedUsersArray = append(includedUsersArray, int(member))
		}
	}

	log.Printf("Deleting message %+v. Delete for array: %v", message, deleteForArray)
	manager.MessagesQueries.DeleteMessage(message, deleteForArray)

	manager.sendMessageEvent(message, chat, "message_deleted", &includedUsersArray)

	return err
}

func (manager *MessagesManager) DeleteReaction(token *jwt.Token, chatId int, messageId int) error {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return err
	}

	chat, err := manager.ChatsRepository.GetWithMember(uint(chatId), uint(tokenSubject.UserId))
	if err != nil {
		return err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: uint(chat.Id), MessageId: uint(messageId), UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return err
	}

	manager.MessagesQueries.DeleteReaction(tokenSubject.UserId, message)
	updatedMessage, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: uint(chat.Id), MessageId: uint(messageId), UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return err
	}

	manager.sendMessageEvent(updatedMessage, chat, "message_reaction_deleted", nil)
	return nil
}

func (manager *MessagesManager) Update(chat *dtos.ChatDto, messageId uint, updateData model.ChangeMessageRequest, token *jwt.Token) (*models.Message, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	message, err := manager.MessagesQueries.GetConcrete(models.GetConcreteMessageParams{ChatId: uint(chat.Id), MessageId: messageId, UserId: uint(tokenSubject.UserId)})
	if err != nil {
		return nil, err
	}

	log.Printf("Updating message %+v. Update data = %+v", message, updateData)
	message.Content = *updateData.Content

	var attachments []models.SavedFile
	for _, attachment := range updateData.Attachments {
		if err := utils.ValidateUploadingFile(*attachment, "file_in_chat"); err != nil {
			return nil, err
		}

		saved_file := generic_factories.UploadingFileToDbFile(*attachment)
		attachments = append(attachments, saved_file)
	}
	message.Attachments = attachments

	mentioned := pq.Int32Array{}
	for _, ment := range updateData.Mentioned {
		mentioned = append(mentioned, int32(*ment))
	}
	message.Mentioned = mentioned

	manager.MessagesQueries.Update(message)

	log.Printf("Updated message: %+v", message)

	manager.sendMessageEvent(message, chat, "message_updated", nil)

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
		ChatsRepository: &repositories.ChatsRepository{},
	}
}
