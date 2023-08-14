package services

import (
	"github.com/chack-check/chats-service/api/v1/models"
)

type ChatsManager struct {
	ChatsQueries *models.ChatsQueries
}

func (manager *ChatsManager) GetConcrete(chatID uint) (*models.Chat, error) {
	chat, err := manager.ChatsQueries.GetConcrete(1, chatID)

	if err != nil {
		return nil, err
	}

	return chat, nil
}
