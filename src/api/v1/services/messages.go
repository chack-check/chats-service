package services

import (
	"log"
	"slices"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/api/v1/utils"
	"github.com/chack-check/chats-service/protousers"
	"github.com/chack-check/chats-service/rabbit"
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
        Type: "message",
        IncludedUsers: members,
        ChatID: int(message.ChatId),
        SenderID: int(message.SenderId),
        MessageType: message.Type,
        Content: message.Content,
        VoiceURL: message.VoiceURL,
        CircleURL: message.CircleURL,
        Attachments: attachments,
        ReplyToID: int(message.ReplyToID),
        Mentioned: mentioned,
        ReadedBy: readedBy,
    }
}

type MessagesManager struct {
	MessagesQueries *models.MessagesQueries
}

func (manager *MessagesManager) GetChatAll(chatId uint, page *int, perPage *int) *schemas.PaginatedResponse[models.Message] {
	count := manager.MessagesQueries.GetAllInChatCount(chatId)
	countValue := *count
	messages := manager.MessagesQueries.GetAllInChat(page, perPage, chatId)
	paginatedResponse := schemas.NewPaginatedResponse[models.Message](*page, *perPage, int(countValue), *messages)
	return &paginatedResponse
}

func (manager *MessagesManager) getTextMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
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
	if err := utils.ValidateVoiceMessage(messageData); err != nil {
		return err
	}

	message.VoiceURL = *messageData.Voice

	if *messageData.ReplyToID != 0 {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	return nil
}

func (manager *MessagesManager) createCircleMessage(message *models.Message, messageData *model.CreateMessageRequest) error {
	if err := utils.ValidateCircleMessage(messageData); err != nil {
		return err
	}

	message.CircleURL = *messageData.Circle

	if *messageData.ReplyToID != 0 {
		message.ReplyToID = uint(*messageData.ReplyToID)
	}

	return nil
}

func (manager *MessagesManager) CreateMessage(messageData *model.CreateMessageRequest, chat *models.Chat, user *protousers.UserResponse) (*models.Message, error) {
	getMessage := map[string]func(message *models.Message, messageData *model.CreateMessageRequest) error{
		"text":   manager.getTextMessage,
		"voice":  manager.createVoiceMessage,
		"circle": manager.createCircleMessage,
	}[messageData.Type.String()]

	message := &models.Message{
		ChatId:   chat.ID,
		Type:     messageData.Type.String(),
		SenderId: uint(user.Id),
	}

	if err := getMessage(message, messageData); err != nil {
		return nil, err
	}

	if err := manager.MessagesQueries.Create(message); err != nil {
		return nil, err
	}

	sendingIds := slices.DeleteFunc(
		chat.Members,
		func(id int64) bool { return id != int64(user.Id) },
	)
	sendingIds32 := []int32{}
	for _, v := range sendingIds {
		sendingIds32 = append(sendingIds32, int32(v))
	}

    messageEvent := getMessageEventFromMessage(message, chat)
    err := rabbit.EventsRabbitConnection.SendMessageEvent(messageEvent)

    if err != nil {
        log.Printf("Error when publishing message event in queue: %v", err)
    }

	return message, nil
}

func NewMessagesManager() *MessagesManager {
	return &MessagesManager{
		MessagesQueries: &models.MessagesQueries{},
	}
}
