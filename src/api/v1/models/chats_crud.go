package models

import (
	"fmt"
	"log"

	"github.com/chack-check/chats-service/database"
	"github.com/lib/pq"
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
	database.DB.Where("owner_id = ? AND id = ?", userId, id).First(&chat)
	if chat.ID == 0 {
		return &Chat{}, fmt.Errorf("Chat with this ID doesn't exist")
	}
	return &chat, nil
}

func (queries *ChatsQueries) GetWithMember(chatId uint, userId uint) (*Chat, error) {
	var chat Chat
	database.DB.Where("? = ANY(members) AND id = ?", userId, chatId).First(&chat)
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

func (queries *ChatsQueries) SearchCount(userId uint, query string, page int, perPage int) int64 {
	var count int64

	database.DB.Scopes(Paginate(page, perPage)).Model(&Chat{}).Joins(
		"JOIN json_each_text(title) d ON true",
	).Where(
		"owner_id = ? AND d.key ILIKE '%?%' AND d.value = '?'", userId, query, userId,
	).Count(&count)

	return count
}

func (queries *ChatsQueries) Search(userId uint, query string, page int, perPage int) *[]Chat {
	var chats []Chat

	database.DB.Scopes(Paginate(page, perPage)).Model(&Chat{}).Joins(
		"JOIN json_each_text(to_json(title)) d ON true",
	).Where(
		"? = ANY(members) AND ((d.key ILIKE ? AND d.value = ?) OR title ILIKE ?)", userId, fmt.Sprintf("%%%s%%", query), fmt.Sprintf("%d", userId), fmt.Sprintf("%%%s%%", query),
	).Find(&chats)

	return &chats
}

func (queries *ChatsQueries) GetExistingWithUser(userId uint, anotherUserId uint) bool {
	var count int64
	database.DB.Model(&Chat{}).Where("? = ANY(members) AND ? = ANY(members)", userId, anotherUserId).Count(&count)
	return count > 0
}

func (queries *ChatsQueries) Delete(chat *Chat) {
	database.DB.Delete(chat)
}

func (queries *ChatsQueries) GetDeletedChatId(members pq.Int32Array) uint {
	chats := new([]Chat)
	database.DB.Unscoped().Model(&Chat{}).Where("deleted_at IS NOT NULL AND members = ? AND type = ?", members, "user").Find(chats)
	if len(*chats) == 0 {
		return 0
	}

	return (*chats)[0].ID
}

func (queries *ChatsQueries) RestoreChat(chatId uint) {
	database.DB.Unscoped().Model(&Chat{}).Where("id  = ?", chatId).Update("deleted_at", nil)
}
