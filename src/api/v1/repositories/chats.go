package repositories

import (
	"fmt"
	"log"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/factories"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/database"
	"github.com/getsentry/sentry-go"
	"github.com/lib/pq"
)

type ChatsRepository struct{}

func (queries *ChatsRepository) Create(chat *dtos.ChatDto) error {
	log.Printf("Creating chat: %+v", chat)
	db_chat := factories.ChatDtoToDbChat(*chat)
	log.Printf("Creating db chat: %+v", db_chat)
	result := database.DB.Create(&db_chat)

	if result.Error != nil {
		log.Printf("Error saving chat to database: %s", result.Error)
		sentry.CaptureException(fmt.Errorf("error saving chat: %s", result.Error))
		return result.Error
	}

	chat.Id = int(db_chat.ID)
	log.Printf("Saved chat id: %d", chat.Id)
	return nil
}

func (queries *ChatsRepository) GetConcrete(userId uint, id uint) (dtos.ChatDto, error) {
	log.Printf("Fetching concrete chat by user id = %d and chat id = %d", userId, id)
	var chat models.Chat
	database.DB.Preload("Avatar").Where("owner_id = ? AND id = ?", userId, id).First(&chat)
	log.Printf("Fetched chat: %+v", chat)
	if chat.ID == 0 {
		log.Printf("Error fetching chat: %+v", chat)
		sentry.CaptureException(fmt.Errorf("error fetching chat: %+v", chat))
		return dtos.ChatDto{}, fmt.Errorf("chat with this ID doesn't exist")
	}
	return factories.DbChatToChatDto(chat, nil), nil
}

func (queries *ChatsRepository) GetByIds(chatIds []int, userId int) ([]*dtos.ChatDto, error) {
	log.Printf("Fetching chats by ids = %v for user id = %d", chatIds, userId)
	var chats []*models.Chat

	database.DB.Preload("Avatar").Where(
		"? = ANY(members) AND id IN ?", userId, chatIds,
	).Find(&chats)
	var chats_dtos []*dtos.ChatDto
	log.Printf("Fetched chats count = %d", len(chats))
	for _, chat := range chats {
		chat_dto := factories.DbChatToChatDto(*chat, nil)
		chats_dtos = append(chats_dtos, &chat_dto)
	}

	return chats_dtos, nil
}

func (queries *ChatsRepository) GetWithMember(chatId uint, userId uint) (*dtos.ChatDto, error) {
	log.Printf("Fetching chat with id = %d and member = %d", chatId, userId)
	var chat models.Chat
	database.DB.Preload("Avatar").Where("? = ANY(members) AND id = ?", userId, chatId).First(&chat)
	log.Printf("Fetched chat: %+v", chat)
	if chat.ID == 0 {
		log.Printf("Error fetching chat with id = %d member id = %d: %+v", chatId, userId, chat)
		sentry.CaptureException(fmt.Errorf("chat with id = %d and member id = %d doesn't exist", chatId, userId))
		return nil, fmt.Errorf("chat with this ID doesn't exist")
	}
	chat_dto := factories.DbChatToChatDto(chat, nil)
	return &chat_dto, nil
}

func (queries *ChatsRepository) GetAllWithMember(userId uint, page int, perPage int) []*dtos.ChatDto {
	log.Printf("Fetching all chats with member id = %d page = %d perPage = %d", userId, page, perPage)
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
	log.Printf("Fetching all user chats count for user id = %d", userId)
	var count int64
	database.DB.Model(&models.Chat{}).Where("? = ANY(members)", userId).Count(&count)
	log.Printf("User chats count: %d", count)
	return &count
}

func (queries *ChatsRepository) GetAll(userId uint) []*dtos.ChatDto {
	log.Printf("Fetching all chats with owner id = %d", userId)
	var chats []models.Chat
	database.DB.Preload("Avatar").Where(&models.Chat{OwnerId: userId}).Find(&chats)
	var chats_dtos []*dtos.ChatDto
	for _, chat := range chats {
		chat_dto := factories.DbChatToChatDto(chat, nil)
		chats_dtos = append(chats_dtos, &chat_dto)
	}

	log.Printf("Fetched all chats with owner id = %d count = %d", userId, len(chats_dtos))
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
	log.Printf("Check is chat with members [%d, %d] exists", userId, anotherUserId)
	var count int64
	database.DB.Model(&models.Chat{}).Where("? = ANY(members) AND ? = ANY(members)", userId, anotherUserId).Count(&count)
	log.Printf("Chat with members [%d, %d] count: %d", userId, anotherUserId, count)
	return count > 0
}

func (queries *ChatsRepository) Delete(chat *dtos.ChatDto) {
	log.Printf("Deleting chat %+v", chat)
	db_chat := factories.ChatDtoToDbChat(*chat)
	database.DB.Delete(&db_chat)
}

func (queries *ChatsRepository) GetDeletedChatId(members pq.Int32Array) uint {
	log.Printf("Fetching deleted chat id with members = %v", members)
	chats := new([]models.Chat)
	database.DB.Unscoped().Model(&models.Chat{}).Where("deleted_at IS NOT NULL AND members = ? AND type = ?", members, "user").Find(chats)
	log.Printf("Fetched deleted chats with members %v count = %d", members, len(*chats))
	if len(*chats) == 0 {
		log.Printf("There is no deleted chat with members %v", members)
		return 0
	}

	log.Printf("Deleted chat id with members = %v: %d", members, (*chats)[0].ID)
	return (*chats)[0].ID
}

func (queries *ChatsRepository) RestoreChat(chatId uint) {
	log.Printf("Restoring chat with id %d", chatId)
	database.DB.Unscoped().Model(&models.Chat{}).Where("id  = ?", chatId).Update("deleted_at", nil)
}
