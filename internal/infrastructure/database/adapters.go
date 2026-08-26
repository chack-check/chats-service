package database

import (
	"fmt"
	"log"
	"math"

	"chats-service/internal/domain/chats"
	"chats-service/internal/domain/files"
	"chats-service/internal/domain/messages"
	"chats-service/internal/domain/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetOrCreateFile(file *files.SavedFile, db gorm.DB) SavedFile {
	if file == nil {
		return SavedFile{}
	}

	var convertedURL string
	var convertedFilename string
	if url := file.GetConvertedUrl(); url != nil {
		convertedURL = *url
		filename := file.GetConvertedFilename()
		convertedFilename = *filename
	}

	var foundedFile SavedFile
	db.Where("original_url = ?", file.GetOriginalUrl()).First(&foundedFile)

	if foundedFile.ID != 0 {
		return foundedFile
	}

	foundedFile.OriginalUrl = file.GetOriginalUrl()
	foundedFile.OriginalFilename = file.GetOriginalFilename()
	foundedFile.ConvertedUrl = convertedURL
	foundedFile.ConvertedFilename = convertedFilename
	db.Save(&foundedFile)

	return foundedFile
}

type ChatsLoggingAdapter struct {
	adapter chats.ChatsPort
}

func (adapter ChatsLoggingAdapter) GetById(id int) (*chats.Chat, error) {
	log.Printf("fetching chat by id: %d", id)
	chat, err := adapter.adapter.GetById(id)
	if err != nil {
		log.Printf("error fetching chat by id: %v", err)
		return chat, err
	}

	log.Printf("fetched chat by id: %+v", chat)
	return chat, err
}

func (adapter ChatsLoggingAdapter) GetByIdForUser(id int, userID int) (*chats.Chat, error) {
	log.Printf("fetching chat by id for user: id=%d, userId=%d", id, userID)
	chat, err := adapter.adapter.GetByIdForUser(id, userID)
	if err != nil {
		log.Printf("error fetching chat by id for user: %v", err)
		return chat, err
	}

	log.Printf("fetched chat by id for user: %+v", chat)
	return chat, err
}

func (adapter ChatsLoggingAdapter) GetByIdsForUser(ids []int, userID int) []chats.Chat {
	log.Printf("fetching chats by ids for user: ids=%+v, userId=%d", ids, userID)
	chats := adapter.adapter.GetByIdsForUser(ids, userID)
	log.Printf("fetched chats by ids for user: %+v", chats)
	return chats
}

func (adapter ChatsLoggingAdapter) GetUserAll(userID int, page int, perPage int) utils.PaginatedResponse[chats.Chat] {
	log.Printf("fetching all chats for user: userId=%d, page=%d, perPage=%d", userID, page, perPage)
	chats := adapter.adapter.GetUserAll(userID, page, perPage)
	log.Printf("fetched all chats for user: %+v", chats)
	return chats
}

func (adapter ChatsLoggingAdapter) Save(chat chats.Chat) (*chats.Chat, error) {
	log.Printf("saving chat: %+v", chat)
	savedChat, err := adapter.adapter.Save(chat)
	if err != nil {
		log.Printf("error saving chat: %v", err)
		return savedChat, err
	}

	log.Printf("saved chat: %+v", savedChat)
	return savedChat, err
}

func (adapter ChatsLoggingAdapter) HasDeletedUserChat(chat chats.Chat) bool {
	log.Printf("checking has deleted user chat: %+v", chat)
	has := adapter.adapter.HasDeletedUserChat(chat)
	log.Printf("has deleted user chat: %v", has)
	return has
}

func (adapter ChatsLoggingAdapter) RestoreChat(chat chats.Chat) (*chats.Chat, error) {
	log.Printf("restoring chat: %+v", chat)
	restoredChat, err := adapter.adapter.RestoreChat(chat)
	if err != nil {
		log.Printf("error restoring chat: %v", err)
		return restoredChat, err
	}

	log.Printf("restored chat: %+v", restoredChat)
	return restoredChat, err
}

func (adapter ChatsLoggingAdapter) CheckChatExists(chat chats.Chat) bool {
	log.Printf("checking chat existing: %+v", chat)
	chatExists := adapter.adapter.CheckChatExists(chat)
	log.Printf("chat exists: %v", chatExists)
	return chatExists
}

func (adapter ChatsLoggingAdapter) Delete(chat chats.Chat) {
	log.Printf("deleting chat: %+v", chat)
	adapter.adapter.Delete(chat)
	log.Printf("deleted chat")
}

func (adapter ChatsLoggingAdapter) SearchChats(userID int, query string, page int, perPage int) utils.PaginatedResponse[chats.Chat] {
	log.Printf("searching chats: query=%s, page=%d, perPage=%d", query, page, perPage)
	chats := adapter.adapter.SearchChats(userID, query, page, perPage)
	log.Printf("founded chats count: %d", len(chats.GetData()))
	return chats
}

type ChatsAdapter struct {
	db gorm.DB
}

func (adapter ChatsAdapter) GetById(id int) (*chats.Chat, error) {
	var chat Chat
	result := adapter.db.Preload("Avatar").Where("id = ?", id).First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := DbChatToModel(chat)
	return &chatModel, nil
}

func (adapter ChatsAdapter) GetByIdForUser(id int, userID int) (*chats.Chat, error) {
	var chat Chat
	result := adapter.db.Preload("Avatar").Where("id = ? AND ? = ANY(members)", id, userID).First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := DbChatToModel(chat)
	return &chatModel, nil
}

func (adapter ChatsAdapter) GetByIdsForUser(ids []int, userID int) []chats.Chat {
	var foundedChats []Chat
	result := adapter.db.Preload("Avatar").Where("id IN ? AND ? = ANY(members)", ids, userID).Find(&foundedChats)
	if result.Error != nil {
		return []chats.Chat{}
	}

	chatModels := make([]chats.Chat, 0, len(foundedChats))
	for _, chat := range foundedChats {
		chatModels = append(chatModels, DbChatToModel(chat))
	}

	return chatModels
}

func (adapter ChatsAdapter) getUserAllCount(userID int) int {
	var count int64
	adapter.db.Model(&Chat{}).Where("? = ANY(members)", userID).Count(&count)
	return int(count)
}

func (adapter ChatsAdapter) GetUserAll(userID int, page int, perPage int) utils.PaginatedResponse[chats.Chat] {
	totalCount := adapter.getUserAllCount(userID)
	if totalCount == 0 {
		return utils.NewPaginatedResponse(
			1, 1, 1, 0, []chats.Chat{},
		)
	}

	var foundedChats []*Chat
	result := adapter.db.Scopes(Paginate(page, perPage)).Preload("Avatar").Where(
		"? = ANY(members)", userID,
	).Order(
		"(SELECT created_at FROM messages WHERE chat_id = chats.id ORDER BY created_at DESC LIMIT 1) DESC NULLS LAST",
	).Find(&foundedChats)

	if result.Error != nil {
		return utils.NewPaginatedResponse(
			1, 1, 1, 0, []chats.Chat{},
		)
	}

	pagesCount := math.Ceil(float64(totalCount) / float64(perPage))
	if pagesCount == 0 {
		pagesCount = 1
	}

	chats := make([]chats.Chat, 0, len(foundedChats))
	for _, chat := range foundedChats {
		chats = append(chats, DbChatToModel(*chat))
	}

	return utils.NewPaginatedResponse(
		page,
		perPage,
		int(pagesCount),
		totalCount,
		chats,
	)
}

func (adapter ChatsAdapter) Save(chat chats.Chat) (*chats.Chat, error) {
	avatarFile := GetOrCreateFile(chat.GetAvatar(), adapter.db)
	dbChat := ModelToDbChat(chat, avatarFile)
	result := adapter.db.Save(&dbChat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := DbChatToModel(dbChat)
	return &chatModel, nil
}

func (adapter ChatsAdapter) HasDeletedUserChat(chat chats.Chat) bool {
	var count int64
	adapter.db.Unscoped().Model(&Chat{}).Where("deleted_at IS NOT NULL AND members = ? AND type = ?", chat.GetMembers(), "user").Count(&count)
	return count > 0
}

func (adapter ChatsAdapter) RestoreChat(chat chats.Chat) (*chats.Chat, error) {
	var dbChat Chat
	result := adapter.db.Unscoped().Model(&Chat{}).Where("id = ?", chat.GetId()).Update("deleted_at", nil)
	if result.Error != nil {
		return nil, result.Error
	}

	adapter.db.Where("id = ?", chat.GetId()).First(&dbChat)
	chatModel := DbChatToModel(dbChat)
	return &chatModel, nil
}

func (adapter ChatsAdapter) CheckChatExists(chat chats.Chat) bool {
	var count int64
	membersIds := make(pq.Int32Array, 0, len(chat.GetMembers()))
	for _, member := range chat.GetMembers() {
		membersIds = append(membersIds, int32(member))
	}

	adapter.db.Model(&Chat{}).Where("members @> ? AND ? @> members AND type = ?", membersIds, membersIds, "user").Count(&count)
	return count > 0
}

func (adapter ChatsAdapter) Delete(chat chats.Chat) {
	adapter.db.Delete(&Chat{ID: uint(chat.GetId())})
}

func (adapter ChatsAdapter) SearchChats(userID int, query string, page int, perPage int) utils.PaginatedResponse[chats.Chat] {
	stmt := adapter.db.Model(&Chat{}).Where("(lower(title) LIKE lower(?) OR title = '') AND ? = ANY(members)", fmt.Sprintf("%%%s%%", query), userID)
	var totalCount int64
	stmt.Count(&totalCount)

	var foundedChats []*Chat
	result := stmt.Scopes(Paginate(page, perPage)).Preload("Avatar").Order(
		"(SELECT created_at FROM messages WHERE chat_id = chats.id ORDER BY created_at DESC LIMIT 1) DESC NULLS LAST",
	).Find(&foundedChats)

	if result.Error != nil {
		return utils.NewPaginatedResponse(
			1, 1, 1, 0, []chats.Chat{},
		)
	}

	pagesCount := math.Ceil(float64(totalCount) / float64(perPage))
	if pagesCount == 0 {
		pagesCount = 1
	}

	chats := make([]chats.Chat, 0, len(foundedChats))
	for _, chat := range foundedChats {
		chats = append(chats, DbChatToModel(*chat))
	}

	return utils.NewPaginatedResponse(
		page,
		perPage,
		int(pagesCount),
		int(totalCount),
		chats,
	)
}

type MessagesLoggingAdapter struct {
	adapter messages.MessagesPort
}

func (adapter MessagesLoggingAdapter) GetChatAllForUser(chatID int, userID int, offset int, limit int) utils.OffsetResponse[messages.Message] {
	log.Printf("fetching chat all messages for user: chatId=%d, userId=%d, offset=%d, limit=%d", chatID, userID, offset, limit)
	messages := adapter.adapter.GetChatAllForUser(chatID, userID, offset, limit)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetChatCursorAllForUser(chatID int, userID int, messageID int, aroundOffset int) utils.OffsetResponse[messages.Message] {
	log.Printf("fetching chat all messages for user by cursor: chatId=%d, userId=%d, messageId=%d, aroundOffset=%d", chatID, userID, messageID, aroundOffset)
	messages := adapter.adapter.GetChatCursorAllForUser(chatID, userID, messageID, aroundOffset)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetChatsLast(chatIds []int, userID int) []messages.Message {
	log.Printf("fetching last messages for chats: chatIds=%v, userId=%d", chatIds, userID)
	messages := adapter.adapter.GetChatsLast(chatIds, userID)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetByIdForUser(messageID int, userID int) (*messages.Message, error) {
	log.Printf("fetching message by id for user: messageId=%d, userId=%d", messageID, userID)
	message, err := adapter.adapter.GetByIdForUser(messageID, userID)
	if err != nil {
		log.Printf("error fetching message by id for user: %v", err)
		return message, err
	}

	log.Printf("fetchec message by id for user: %+v", message)
	return message, err
}

func (adapter MessagesLoggingAdapter) GetByIdsForUser(messageIds []int, userID int) []messages.Message {
	log.Printf("fetching messages by ids for user: messageIds=%v, userId=%d", messageIds, userID)
	messages := adapter.adapter.GetByIdsForUser(messageIds, userID)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetById(messageID int) (*messages.Message, error) {
	log.Printf("fetching message by id messageId=%v", messageID)
	message, err := adapter.adapter.GetById(messageID)
	log.Printf("fetched message: %+v", message)
	return message, err
}

func (adapter MessagesLoggingAdapter) Save(message messages.Message) (*messages.Message, error) {
	log.Printf("saving message: %+v", message)
	savedMessage, err := adapter.adapter.Save(message)
	if err != nil {
		log.Printf("error saving message: %v", err)
		return savedMessage, err
	}

	log.Printf("saved message: %+v", savedMessage)
	return savedMessage, err
}

func (adapter MessagesLoggingAdapter) Delete(message messages.Message) {
	log.Printf("deleting message: %+v", message)
	adapter.adapter.Delete(message)
	log.Printf("message deleted")
}

type MessagesAdapter struct {
	db gorm.DB
}

func (adapter MessagesAdapter) getChatAllForUserTotal(chatID int, userID int) int {
	var count int64

	adapter.db.Model(&Message{}).Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatID, userID,
	).Count(&count)

	return int(count)
}

func (adapter MessagesAdapter) GetChatAllForUser(chatID int, userID int, offset int, limit int) utils.OffsetResponse[messages.Message] {
	var dbMessages []Message

	total := adapter.getChatAllForUserTotal(chatID, userID)

	adapter.db.Preload("Chat").Preload("Reactions").Preload("Voice").Preload("Circle").Preload("Attachments").Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatID, userID,
	).Order(
		"messages.created_at DESC NULLS LAST",
	).Offset(offset).Limit(limit).Find(&dbMessages)

	messagesModels := make([]messages.Message, 0, len(dbMessages))
	for _, dbMessage := range dbMessages {
		messagesModels = append(messagesModels, DbMessageToModel(dbMessage))
	}

	return utils.NewOffsetResponse(
		offset,
		limit,
		total,
		messagesModels,
	)
}

func (adapter MessagesAdapter) getMessageOffsetById(chatID int, userID int, messageID int) int {
	var offset int64

	adapter.db.Model(&Message{}).Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members) AND messages.id >= ?", chatID, userID, messageID,
	).Order("messages.created_at DESC NULLS LAST").Count(&offset)

	return int(offset)
}

func (adapter MessagesAdapter) GetChatCursorAllForUser(chatID int, userID int, messageID int, aroundOffset int) utils.OffsetResponse[messages.Message] {
	offset := adapter.getMessageOffsetById(chatID, userID, messageID)
	startOffset := max(offset-aroundOffset, 0)
	return adapter.GetChatAllForUser(chatID, userID, startOffset, aroundOffset*2)
}

func (adapter MessagesAdapter) GetChatsLast(chatIds []int, userID int) []messages.Message {
	messages := make([]messages.Message, 0, len(chatIds))

	for _, chatID := range chatIds {
		var message Message

		adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
			"messages.chat_id = ? AND ? = ANY(chats.members)", chatID, userID,
		).Order("messages.created_at DESC NULLS LAST").Limit(1).First(&message)

		messageModel := DbMessageToModel(message)
		messages = append(messages, messageModel)
	}

	return messages
}

func (adapter MessagesAdapter) GetById(messageID int) (*messages.Message, error) {
	var dbMessage Message

	result := adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id = ?", messageID,
	).First(&dbMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	messageModel := DbMessageToModel(dbMessage)
	return &messageModel, nil
}

func (adapter MessagesAdapter) GetByIdForUser(messageID int, userID int) (*messages.Message, error) {
	var dbMessage Message

	result := adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id = ? AND ? = ANY(chats.members)", messageID, userID,
	).First(&dbMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	messageModel := DbMessageToModel(dbMessage)
	return &messageModel, nil
}

func (adapter MessagesAdapter) GetByIdsForUser(messageIds []int, userID int) []messages.Message {
	var dbMessages []Message

	adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id IN ? AND ? = ANY(chats.members)", messageIds, userID,
	).Find(&dbMessages)

	modelMessages := make([]messages.Message, 0, len(dbMessages))
	for _, dbMessage := range dbMessages {
		modelMessages = append(modelMessages, DbMessageToModel(dbMessage))
	}

	return modelMessages
}

func (adapter MessagesAdapter) getOrCreateReaction(reaction messages.MessageReaction) Reaction {
	var foundedReaction Reaction

	adapter.db.Where("user_id = ?", reaction.GetUserId()).First(&foundedReaction)

	if foundedReaction.UserId == uint(reaction.GetUserId()) {
		foundedReaction.Content = reaction.GetContent()
		adapter.db.Save(&foundedReaction)
		return foundedReaction
	}

	dbReaction := Reaction{
		UserId:  uint(reaction.GetUserId()),
		Content: reaction.GetContent(),
	}
	adapter.db.Save(&dbReaction)
	return dbReaction
}

func (adapter MessagesAdapter) Save(message messages.Message) (*messages.Message, error) {
	circle := GetOrCreateFile(message.GetCircle(), adapter.db)
	var circlePointer *SavedFile
	if circle.ID != 0 {
		circlePointer = &circle
	}

	voice := GetOrCreateFile(message.GetVoice(), adapter.db)
	var voicePointer *SavedFile
	if voice.ID != 0 {
		voicePointer = &voice
	}

	attachments := make([]SavedFile, 0, len(message.GetAttachments()))
	for _, attachment := range message.GetAttachments() {
		attachments = append(attachments, GetOrCreateFile(&attachment, adapter.db))
	}

	reactions := make([]Reaction, 0, len(message.GetReactions()))
	for _, reaction := range message.GetReactions() {
		reactions = append(reactions, adapter.getOrCreateReaction(reaction))
	}

	dbMessage := ModelToDbMessage(message, voicePointer, circlePointer, attachments, reactions)
	result := adapter.db.Save(&dbMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	savedMessage, err := adapter.GetById(int(dbMessage.ID))
	if err != nil {
		return nil, err
	}

	return savedMessage, nil
}

func (adapter MessagesAdapter) Delete(message messages.Message) {
	adapter.db.Delete(&Message{ID: uint(message.GetId())})
}

func NewChatsAdapter(db gorm.DB) chats.ChatsPort {
	return ChatsLoggingAdapter{adapter: ChatsAdapter{db: db}}
}

func NewMessagesAdapter(db gorm.DB) messages.MessagesPort {
	return MessagesLoggingAdapter{adapter: MessagesAdapter{db: db}}
}
