package utils

import (
	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
)

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
