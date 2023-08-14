package models

import (
	"fmt"

	"github.com/chack-check/chats-service/database"
)

type ChatsQueries struct{}

func (queries *ChatsQueries) GetConcrete(userId uint, id uint) (*Chat, error) {
	var chat Chat
	database.DB.Where("owner_id = ?", userId).Where("id = ?", id).First(&chat)
	if chat.ID == 0 {
		return &Chat{}, fmt.Errorf("User with this ID doesn't exist")
	}
	return &chat, nil
}

func (queries *ChatsQueries) GetAll(userId uint) *[]Chat {
	var chats []Chat
	database.DB.Where(&Chat{OwnerId: userId}).Find(&chats)
	return &chats
}
