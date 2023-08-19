package utils

import (
	"fmt"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
)

func ChatRequestToDbChat(request *model.CreateChatRequest) (*models.Chat, error) {
	if request.User != nil {
		return &models.Chat{}, nil
	}

	if request.Title == nil || request.Avatar == nil || request.Members == nil || len(request.Members) == 0 {
		var members []int64
		for _, v := range request.Members {
			members = append(members, int64(v))
		}

		return &models.Chat{Title: *request.Title, Members: members}, nil
	}

	return nil, fmt.Errorf("Not enough data for chat creation")
}

func DbChatToSchema(chat *models.Chat) *model.Chat {
	var members []int
	var admins []int

	for _, v := range chat.Members {
		members = append(members, int(v))
	}

	for _, v := range chat.Admins {
		admins = append(admins, int(v))
	}

	return &model.Chat{
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
