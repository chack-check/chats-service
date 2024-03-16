package utils

import (
	"fmt"
	"time"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/factories"
)

func ChatRequestToDbChat(request *model.CreateChatRequest) (*models.Chat, error) {
	if request.User != nil {
		return &models.Chat{}, nil
	}

	if request.Title != nil && request.Avatar != nil && request.Members != nil {
		var members []int64
		for _, v := range request.Members {
			members = append(members, int64(v))
		}

		var converted_url string
		var converted_filename string
		if request.Avatar.Converted != nil {
			converted_url = request.Avatar.Converted.URL
			converted_filename = request.Avatar.Converted.Filename
		} else {
			converted_url = ""
			converted_filename = ""
		}

		avatar := models.SavedFile{
			OriginalUrl:       request.Avatar.Original.URL,
			OriginalFilename:  request.Avatar.Original.Filename,
			ConvertedUrl:      converted_url,
			ConvertedFilename: converted_filename,
		}

		return &models.Chat{Title: *request.Title, Members: members, Avatar: avatar}, nil
	}

	return nil, fmt.Errorf("not enough data for chat creation")
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

	avatar := factories.DbFileToSchema(chat.Avatar)

	return model.Chat{
		ID:         int(chat.ID),
		Avatar:     &avatar,
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

	var voice_schema *model.SavedFile
	if message.VoiceId != nil {
		schema := factories.DbFileToSchema(message.Voice)
		voice_schema = &schema
	}

	var circle_schema *model.SavedFile
	if message.CircleId != nil {
		schema := factories.DbFileToSchema(message.Circle)
		circle_schema = &schema
	}

	var attachments []*model.SavedFile
	for _, attachment := range message.Attachments {
		schema := factories.DbFileToSchema(attachment)
		attachments = append(attachments, &schema)
	}

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
		Voice:       voice_schema,
		Circle:      circle_schema,
		ReplyToID:   &replyTo,
		ReadedBy:    readedBy,
		Reactions:   reactions,
		Attachments: attachments,
		Mentioned:   mentioned,
		Datetime:    message.CreatedAt.UTC().Format(time.RFC3339),
	}
}
