package repositories

import (
	"fmt"
	"log"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/factories"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/database"
	"github.com/lib/pq"
)

type ChatsRepository struct{}

func (queries *ChatsRepository) Create(chat *dtos.ChatDto) error {
	db_chat := factories.ChatDtoToDbChat(*chat)
	result := database.DB.Create(&db_chat)

	if result.Error != nil {
		return result.Error
	}

	chat.Id = int(db_chat.ID)
	return nil
}

func (queries *ChatsRepository) GetConcrete(userId uint, id uint) (dtos.ChatDto, error) {
	var chat models.Chat
	database.DB.Preload("Avatar").Where("owner_id = ? AND id = ?", userId, id).First(&chat)
	if chat.ID == 0 {
		return dtos.ChatDto{}, fmt.Errorf("chat with this ID doesn't exist")
	}
	return factories.DbChatToChatDto(chat, nil), nil
}

func (queries *ChatsRepository) GetByIds(chatIds []int, userId int) ([]*dtos.ChatDto, error) {
	var chats []*models.Chat

	log.Printf("Getting chats with ids %v and for user id %d", chatIds, userId)

	database.DB.Preload("Avatar").Where(
		"? = ANY(members) AND id IN ?", userId, chatIds,
	).Find(&chats)
	var chats_dtos []*dtos.ChatDto
	for _, chat := range chats {
		chat_dto := factories.DbChatToChatDto(*chat, nil)
		chats_dtos = append(chats_dtos, &chat_dto)
	}

	return chats_dtos, nil
}

func (queries *ChatsRepository) GetWithMember(chatId uint, userId uint) (*dtos.ChatDto, error) {
	var chat models.Chat
	database.DB.Preload("Avatar").Where("? = ANY(members) AND id = ?", userId, chatId).First(&chat)
	if chat.ID == 0 {
		return nil, fmt.Errorf("chat with this ID doesn't exist")
	}
	chat_dto := factories.DbChatToChatDto(chat, nil)
	return &chat_dto, nil
}

func (queries *ChatsRepository) GetAllWithMember(userId uint, page int, perPage int) []*dtos.ChatDto {
	var chats []*models.Chat

	database.DB.Scopes(
		models.Paginate(page, perPage),
	).Preload("Avatar").Where(
		"? = ANY(members)", userId,
	).Order(
		"(SELECT created_at FROM messages WHERE chat_id = chats.id ORDER BY created_at DESC LIMIT 1) DESC NULLS LAST",
	).Find(&chats)

	log.Printf("User chats count: %v", len(chats))
	var chats_dtos []*dtos.ChatDto
	for _, chat := range chats {
		chat_dto := factories.DbChatToChatDto(*chat, nil)
		chats_dtos = append(chats_dtos, &chat_dto)
	}

	return chats_dtos
}

func (queries *ChatsRepository) GetAllWithMemberCount(userId uint) *int64 {
	var count int64
	database.DB.Model(&models.Chat{}).Where("? = ANY(members)", userId).Count(&count)
	return &count
}

func (queries *ChatsRepository) GetAll(userId uint) []*dtos.ChatDto {
	var chats []models.Chat
	database.DB.Preload("Avatar").Where(&models.Chat{OwnerId: userId}).Find(&chats)
	var chats_dtos []*dtos.ChatDto
	for _, chat := range chats {
		chat_dto := factories.DbChatToChatDto(chat, nil)
		chats_dtos = append(chats_dtos, &chat_dto)
	}

	return chats_dtos
}

func (queries *ChatsRepository) SearchCount(userId uint, query string, page int, perPage int) int64 {
	var count int64

	database.DB.Scopes(models.Paginate(page, perPage)).Model(&models.Chat{}).Joins(
		"JOIN json_each_text(title) d ON true",
	).Where(
		"owner_id = ? AND d.key ILIKE '%?%' AND d.value = '?'", userId, query, userId,
	).Count(&count)

	return count
}

func (queries *ChatsRepository) GetExistingWithUser(userId uint, anotherUserId uint) bool {
	var count int64
	database.DB.Model(&models.Chat{}).Where("? = ANY(members) AND ? = ANY(members)", userId, anotherUserId).Count(&count)
	return count > 0
}

func (queries *ChatsRepository) Delete(chat *dtos.ChatDto) {
	db_chat := factories.ChatDtoToDbChat(*chat)
	database.DB.Delete(&db_chat)
}

func (queries *ChatsRepository) GetDeletedChatId(members pq.Int32Array) uint {
	chats := new([]models.Chat)
	database.DB.Unscoped().Model(&models.Chat{}).Where("deleted_at IS NOT NULL AND members = ? AND type = ?", members, "user").Find(chats)
	if len(*chats) == 0 {
		return 0
	}

	return (*chats)[0].ID
}

func (queries *ChatsRepository) RestoreChat(chatId uint) {
	database.DB.Unscoped().Model(&models.Chat{}).Where("id  = ?", chatId).Update("deleted_at", nil)
}
