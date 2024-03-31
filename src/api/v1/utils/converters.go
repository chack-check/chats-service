package utils

import (
	"fmt"
	"time"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/generic_factories"
)

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
		schema := generic_factories.DbFileToSchema(message.Voice)
		voice_schema = &schema
	}

	var circle_schema *model.SavedFile
	if message.CircleId != nil {
		schema := generic_factories.DbFileToSchema(message.Circle)
		circle_schema = &schema
	}

	var attachments []*model.SavedFile
	for _, attachment := range message.Attachments {
		schema := generic_factories.DbFileToSchema(attachment)
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

func GetUserFullName(firstName string, lastName string, middleName *string) string {
	if middleName != nil && *middleName != "" {
		return fmt.Sprintf("%s %s %s", lastName, firstName, *middleName)
	}

	return fmt.Sprintf("%s %s", lastName, firstName)
}
