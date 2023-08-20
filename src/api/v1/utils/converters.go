package utils

import (
	"fmt"
	"time"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
)

func ChatRequestToDbChat(request *model.CreateChatRequest) (*models.Chat, error) {
	if request.User != nil {
		return &models.Chat{}, nil
	}

	if request.Title != nil && request.AvatarURL != nil && request.Members != nil {
		var members []int64
		for _, v := range request.Members {
			members = append(members, int64(v))
		}

		return &models.Chat{Title: *request.Title, Members: members, AvatarURL: *request.AvatarURL}, nil
	}

	return nil, fmt.Errorf("Not enough data for chat creation")
}

func DbChatToSchema(chat models.Chat) model.Chat {
	var members []int
	var admins []int

	for _, v := range chat.Members {
		members = append(members, int(v))
	}

	for _, v := range chat.Admins {
		admins = append(admins, int(v))
	}

	return model.Chat{
		ID:         int(chat.ID),
		AvatarURL:  chat.AvatarURL,
		Title:      chat.Title,
		Type:       model.ChatType(chat.Type),
		Members:    members,
		IsArchived: chat.IsArchived,
		OwnerID:    int(chat.OwnerId),
		Admins:     admins,
	}
}

func DbMessageToSchema(message models.Message) model.Message {
	senderId := int(message.SenderId)
	replyTo := int(message.ReplyToID)

	var readedBy []int
	for _, v := range message.ReadedBy {
		readedBy = append(readedBy, int(v))
	}

	var mentioned []int
	for _, v := range message.Mentioned {
		mentioned = append(mentioned, int(v))
	}

	var attachments []*model.FileObjectResponse

	var reactions []*model.Reaction
	for _, v := range message.Reactions {
		reactions = append(reactions, &model.Reaction{
			Content: v.Content,
			UserID:  int(v.UserId),
		})
	}

	return model.Message{
		ID:          int(message.ID),
		Type:        model.MessageType(message.Type),
		SenderID:    &senderId,
		ChatID:      int(message.ChatId),
		Content:     &message.Content,
		VoiceURL:    &message.VoiceURL,
		CircleURL:   &message.CircleURL,
		ReplyToID:   &replyTo,
		ReadedBy:    readedBy,
		Reactions:   reactions,
		Attachments: attachments,
		Mentioned:   mentioned,
		Datetime:    message.CreatedAt.UTC().Format(time.RFC3339),
	}
}
