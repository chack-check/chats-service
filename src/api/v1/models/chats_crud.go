package models

import (
	"fmt"
	"log"

	"github.com/chack-check/chats-service/database"
)

type ChatsQueries struct{}

func (queries *ChatsQueries) Create(chat *Chat) error {
	result := database.DB.Create(chat)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (queries *ChatsQueries) GetConcrete(userId uint, id uint) (*Chat, error) {
	var chat Chat
	database.DB.Where("owner_id = ?", userId).Where("id = ?", id).First(&chat)
	if chat.ID == 0 {
		return &Chat{}, fmt.Errorf("Chat with this ID doesn't exist")
	}
	return &chat, nil
}

func (queries *ChatsQueries) GetWithMember(chatId uint, userId uint) (*Chat, error) {
	var chat Chat
	database.DB.Where("? = ANY(members)", userId).Where("id = ?", chatId).First(&chat)
	if chat.ID == 0 {
		return &Chat{}, fmt.Errorf("Chat with this ID doesn't exist")
	}
	return &chat, nil
}

func (queries *ChatsQueries) GetAllWithMember(userId uint, page int, perPage int) *[]Chat {
	var chats []Chat

	database.DB.Scopes(Paginate(page, perPage)).Where("? = ANY(members)", userId).Find(&chats)
	log.Printf("User chats count: %v", len(chats))
	return &chats
}

func (queries *ChatsQueries) GetAllWithMemberCount(userId uint) *int64 {
	var count int64
	database.DB.Model(&Chat{}).Where("? = ANY(members)", userId).Count(&count)
	return &count
}

func (queries *ChatsQueries) GetAll(userId uint) *[]Chat {
	var chats []Chat
	database.DB.Where(&Chat{OwnerId: userId}).Find(&chats)
	return &chats
}

func (queries *ChatsQueries) GetExistingWithUser(userId uint, anotherUserId uint) bool {
	var count int64
	database.DB.Model(&Chat{}).Where("? = ANY(members) AND ? = ANY(members)", userId, anotherUserId).Count(&count)
	return count > 0
}
