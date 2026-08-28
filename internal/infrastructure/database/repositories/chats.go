package repositories

import (
	"fmt"
	"log"
	"math"

	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/database"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func getOrCreateFile(file *entities.SavedFile, db gorm.DB) database.SavedFile {
	if file == nil {
		return database.SavedFile{}
	}

	var convertedURL string
	var convertedFilename string
	if url := file.GetConvertedURL(); url != nil {
		convertedURL = *url
		filename := file.GetConvertedFilename()
		convertedFilename = *filename
	}

	var foundedFile database.SavedFile
	db.Where("original_url = ?", file.GetOriginalURL()).First(&foundedFile)

	if foundedFile.ID != 0 {
		return foundedFile
	}

	foundedFile.OriginalUrl = file.GetOriginalURL()
	foundedFile.OriginalFilename = file.GetOriginalFilename()
	foundedFile.ConvertedUrl = convertedURL
	foundedFile.ConvertedFilename = convertedFilename
	db.Save(&foundedFile)

	return foundedFile
}

type ChatsLoggingRepository struct {
	repository ports.ChatsRepository
}

func (adapter ChatsLoggingRepository) GetById(id int) (*entities.Chat, error) {
	log.Printf("fetching chat by id: %d", id)
	chat, err := adapter.repository.GetById(id)
	if err != nil {
		log.Printf("error fetching chat by id: %v", err)
		return chat, err
	}

	log.Printf("fetched chat by id: %+v", chat)
	return chat, err
}

func (adapter ChatsLoggingRepository) GetByIdForUser(id int, userID int) (*entities.Chat, error) {
	log.Printf("fetching chat by id for user: id=%d, userId=%d", id, userID)
	chat, err := adapter.repository.GetByIdForUser(id, userID)
	if err != nil {
		log.Printf("error fetching chat by id for user: %v", err)
		return chat, err
	}

	log.Printf("fetched chat by id for user: %+v", chat)
	return chat, err
}

func (adapter ChatsLoggingRepository) GetByIdsForUser(ids []int, userID int) []entities.Chat {
	log.Printf("fetching chats by ids for user: ids=%+v, userId=%d", ids, userID)
	chats := adapter.repository.GetByIdsForUser(ids, userID)
	log.Printf("fetched chats by ids for user: %+v", chats)
	return chats
}

func (adapter ChatsLoggingRepository) GetUserAll(userID int, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	log.Printf("fetching all chats for user: userId=%d, page=%d, perPage=%d", userID, page, perPage)
	chats := adapter.repository.GetUserAll(userID, page, perPage)
	log.Printf("fetched all chats for user: %+v", chats)
	return chats
}

func (adapter ChatsLoggingRepository) Save(chat entities.Chat) (*entities.Chat, error) {
	log.Printf("saving chat: %+v", chat)
	savedChat, err := adapter.repository.Save(chat)
	if err != nil {
		log.Printf("error saving chat: %v", err)
		return savedChat, err
	}

	log.Printf("saved chat: %+v", savedChat)
	return savedChat, err
}

func (adapter ChatsLoggingRepository) HasDeletedUserChat(chat entities.Chat) bool {
	log.Printf("checking has deleted user chat: %+v", chat)
	has := adapter.repository.HasDeletedUserChat(chat)
	log.Printf("has deleted user chat: %v", has)
	return has
}

func (adapter ChatsLoggingRepository) RestoreChat(chat entities.Chat) (*entities.Chat, error) {
	log.Printf("restoring chat: %+v", chat)
	restoredChat, err := adapter.repository.RestoreChat(chat)
	if err != nil {
		log.Printf("error restoring chat: %v", err)
		return restoredChat, err
	}

	log.Printf("restored chat: %+v", restoredChat)
	return restoredChat, err
}

func (adapter ChatsLoggingRepository) CheckChatExists(chat entities.Chat) bool {
	log.Printf("checking chat existing: %+v", chat)
	chatExists := adapter.repository.CheckChatExists(chat)
	log.Printf("chat exists: %v", chatExists)
	return chatExists
}

func (adapter ChatsLoggingRepository) Delete(chat entities.Chat) {
	log.Printf("deleting chat: %+v", chat)
	adapter.repository.Delete(chat)
	log.Printf("deleted chat")
}

func (adapter ChatsLoggingRepository) SearchChats(userID int, query string, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	log.Printf("searching chats: query=%s, page=%d, perPage=%d", query, page, perPage)
	chats := adapter.repository.SearchChats(userID, query, page, perPage)
	log.Printf("founded chats count: %d", len(chats.GetData()))
	return chats
}

type ChatsRepository struct {
	db gorm.DB
}

func (adapter ChatsRepository) GetById(id int) (*entities.Chat, error) {
	var chat database.Chat
	result := adapter.db.Preload("Avatar").Where("id = ?", id).First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := database.DbChatToModel(chat)
	return &chatModel, nil
}

func (adapter ChatsRepository) GetByIdForUser(id int, userID int) (*entities.Chat, error) {
	var chat database.Chat
	result := adapter.db.Preload("Avatar").Where("id = ? AND ? = ANY(members)", id, userID).First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := database.DbChatToModel(chat)
	return &chatModel, nil
}

func (adapter ChatsRepository) GetByIdsForUser(ids []int, userID int) []entities.Chat {
	var foundedChats []database.Chat
	result := adapter.db.Preload("Avatar").Where("id IN ? AND ? = ANY(members)", ids, userID).Find(&foundedChats)
	if result.Error != nil {
		return []entities.Chat{}
	}

	chatModels := make([]entities.Chat, 0, len(foundedChats))
	for _, chat := range foundedChats {
		chatModels = append(chatModels, database.DbChatToModel(chat))
	}

	return chatModels
}

func (adapter ChatsRepository) getUserAllCount(userID int) int {
	var count int64
	adapter.db.Model(&database.Chat{}).Where("? = ANY(members)", userID).Count(&count)
	return int(count)
}

func (adapter ChatsRepository) GetUserAll(userID int, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	totalCount := adapter.getUserAllCount(userID)
	if totalCount == 0 {
		return dtos.NewPaginatedResponse(
			1, 1, 1, 0, []entities.Chat{},
		)
	}

	var foundedChats []*database.Chat
	result := adapter.db.Scopes(database.Paginate(page, perPage)).Preload("Avatar").Where(
		"? = ANY(members)", userID,
	).Order(
		"(SELECT created_at FROM messages WHERE chat_id = chats.id ORDER BY created_at DESC LIMIT 1) DESC NULLS LAST",
	).Find(&foundedChats)

	if result.Error != nil {
		return dtos.NewPaginatedResponse(
			1, 1, 1, 0, []entities.Chat{},
		)
	}

	pagesCount := math.Ceil(float64(totalCount) / float64(perPage))
	if pagesCount == 0 {
		pagesCount = 1
	}

	chats := make([]entities.Chat, 0, len(foundedChats))
	for _, chat := range foundedChats {
		chats = append(chats, database.DbChatToModel(*chat))
	}

	return dtos.NewPaginatedResponse(
		page,
		perPage,
		int(pagesCount),
		totalCount,
		chats,
	)
}

func (adapter ChatsRepository) Save(chat entities.Chat) (*entities.Chat, error) {
	avatarFile := getOrCreateFile(chat.GetAvatar(), adapter.db)
	dbChat := database.ModelToDbChat(chat, avatarFile)
	result := adapter.db.Save(&dbChat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := database.DbChatToModel(dbChat)
	return &chatModel, nil
}

func (adapter ChatsRepository) HasDeletedUserChat(chat entities.Chat) bool {
	var count int64
	adapter.db.Unscoped().Model(&database.Chat{}).Where("deleted_at IS NOT NULL AND members = ? AND type = ?", chat.GetMembers(), "user").Count(&count)
	return count > 0
}

func (adapter ChatsRepository) RestoreChat(chat entities.Chat) (*entities.Chat, error) {
	var dbChat database.Chat
	result := adapter.db.Unscoped().Model(&database.Chat{}).Where("id = ?", chat.GetID()).Update("deleted_at", nil)
	if result.Error != nil {
		return nil, result.Error
	}

	adapter.db.Where("id = ?", chat.GetID()).First(&dbChat)
	chatModel := database.DbChatToModel(dbChat)
	return &chatModel, nil
}

func (adapter ChatsRepository) CheckChatExists(chat entities.Chat) bool {
	var count int64
	membersIds := make(pq.Int32Array, 0, len(chat.GetMembers()))
	for _, member := range chat.GetMembers() {
		membersIds = append(membersIds, int32(member))
	}

	adapter.db.Model(&database.Chat{}).Where("members @> ? AND ? @> members AND type = ?", membersIds, membersIds, "user").Count(&count)
	return count > 0
}

func (adapter ChatsRepository) Delete(chat entities.Chat) {
	adapter.db.Delete(&database.Chat{ID: uint(chat.GetID())})
}

func (adapter ChatsRepository) SearchChats(userID int, query string, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	stmt := adapter.db.Model(&database.Chat{}).Where("(lower(title) LIKE lower(?) OR title = '') AND ? = ANY(members)", fmt.Sprintf("%%%s%%", query), userID)
	var totalCount int64
	stmt.Count(&totalCount)

	var foundedChats []*database.Chat
	result := stmt.Scopes(database.Paginate(page, perPage)).Preload("Avatar").Order(
		"(SELECT created_at FROM messages WHERE chat_id = chats.id ORDER BY created_at DESC LIMIT 1) DESC NULLS LAST",
	).Find(&foundedChats)

	if result.Error != nil {
		return dtos.NewPaginatedResponse(
			1, 1, 1, 0, []entities.Chat{},
		)
	}

	pagesCount := math.Ceil(float64(totalCount) / float64(perPage))
	if pagesCount == 0 {
		pagesCount = 1
	}

	chats := make([]entities.Chat, 0, len(foundedChats))
	for _, chat := range foundedChats {
		chats = append(chats, database.DbChatToModel(*chat))
	}

	return dtos.NewPaginatedResponse(
		page,
		perPage,
		int(pagesCount),
		int(totalCount),
		chats,
	)
}

func NewChatsRepository(db gorm.DB) ports.ChatsRepository {
	return ChatsLoggingRepository{repository: ChatsRepository{db: db}}
}
