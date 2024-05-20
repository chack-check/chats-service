package database

import (
	"log"
	"math"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/messages"
	"github.com/chack-check/chats-service/domain/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetOrCreateFile(file *files.SavedFile, db gorm.DB) SavedFile {
	if file == nil {
		return SavedFile{}
	}

	var convertedUrl string
	var convertedFilename string
	if url := file.GetConvertedUrl(); url != nil {
		convertedUrl = *url
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
	foundedFile.ConvertedUrl = convertedUrl
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

func (adapter ChatsLoggingAdapter) GetByIdForUser(id int, userId int) (*chats.Chat, error) {
	log.Printf("fetching chat by id for user: id=%d, userId=%d", id, userId)
	chat, err := adapter.adapter.GetByIdForUser(id, userId)
	if err != nil {
		log.Printf("error fetching chat by id for user: %v", err)
		return chat, err
	}

	log.Printf("fetched chat by id for user: %+v", chat)
	return chat, err
}

func (adapter ChatsLoggingAdapter) GetByIdsForUser(ids []int, userId int) []chats.Chat {
	log.Printf("fetching chats by ids for user: ids=%+v, userId=%d", ids, userId)
	chats := adapter.adapter.GetByIdsForUser(ids, userId)
	log.Printf("fetched chats by ids for user: %+v", chats)
	return chats
}

func (adapter ChatsLoggingAdapter) GetUserAll(userId int, page int, perPage int) utils.PaginatedResponse[chats.Chat] {
	log.Printf("fetching all chats for user: userId=%d, page=%d, perPage=%d", userId, page, perPage)
	chats := adapter.adapter.GetUserAll(userId, page, perPage)
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

func (adapter ChatsAdapter) GetByIdForUser(id int, userId int) (*chats.Chat, error) {
	var chat Chat
	result := adapter.db.Preload("Avatar").Where("id = ? AND ? = ANY(members)", id, userId).First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

	chatModel := DbChatToModel(chat)
	return &chatModel, nil
}

func (adapter ChatsAdapter) GetByIdsForUser(ids []int, userId int) []chats.Chat {
	var foundedChats []Chat
	result := adapter.db.Preload("Avatar").Where("id IN ? AND ? = ANY(members)", ids, userId).Find(&foundedChats)
	if result.Error != nil {
		return []chats.Chat{}
	}

	var chatModels []chats.Chat
	for _, chat := range foundedChats {
		chatModels = append(chatModels, DbChatToModel(chat))
	}

	return chatModels
}

func (adapter ChatsAdapter) getUserAllCount(userId int, page int, perPage int) int {
	var count int64
	adapter.db.Model(&Chat{}).Where("? = ANY(members)", userId).Count(&count)
	return int(count)
}

func (adapter ChatsAdapter) GetUserAll(userId int, page int, perPage int) utils.PaginatedResponse[chats.Chat] {
	totalCount := adapter.getUserAllCount(userId, page, perPage)
	if totalCount == 0 {
		return utils.NewPaginatedResponse(
			1, 1, 1, 0, []chats.Chat{},
		)
	}

	var foundedChats []*Chat
	result := adapter.db.Scopes(Paginate(page, perPage)).Preload("Avatar").Where(
		"? = ANY(members)", userId,
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

	var chats []chats.Chat
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
	var membersIds pq.Int32Array
	for _, member := range chat.GetMembers() {
		membersIds = append(membersIds, int32(member))
	}

	adapter.db.Model(&Chat{}).Where("sort(members) = sort(?) AND type = ?", membersIds, "user").Count(&count)
	return count > 0
}

func (adapter ChatsAdapter) Delete(chat chats.Chat) {
	adapter.db.Delete(&Chat{ID: uint(chat.GetId())})
}

type MessagesLoggingAdapter struct {
	adapter messages.MessagesPort
}

func (adapter MessagesLoggingAdapter) GetChatAllForUser(chatId int, userId int, offset int, limit int) utils.OffsetResponse[messages.Message] {
	log.Printf("fetching chat all messages for user: chatId=%d, userId=%d, offset=%d, limit=%d", chatId, userId, offset, limit)
	messages := adapter.adapter.GetChatAllForUser(chatId, userId, offset, limit)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetChatCursorAllForUser(chatId int, userId int, messageId int, aroundOffset int) utils.OffsetResponse[messages.Message] {
	log.Printf("fetching chat all messages for user by cursor: chatId=%d, userId=%d, messageId=%d, aroundOffset=%d", chatId, userId, messageId, aroundOffset)
	messages := adapter.adapter.GetChatCursorAllForUser(chatId, userId, messageId, aroundOffset)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetChatsLast(chatIds []int, userId int) []messages.Message {
	log.Printf("fetching last messages for chats: chatIds=%v, userId=%d", chatIds, userId)
	messages := adapter.adapter.GetChatsLast(chatIds, userId)
	log.Printf("fetched messages: %+v", messages)
	return messages
}

func (adapter MessagesLoggingAdapter) GetByIdForUser(messageId int, userId int) (*messages.Message, error) {
	log.Printf("fetching message by id for user: messageId=%d, userId=%d", messageId, userId)
	message, err := adapter.adapter.GetByIdForUser(messageId, userId)
	if err != nil {
		log.Printf("error fetching message by id for user: %v", err)
		return message, err
	}

	log.Printf("fetchec message by id for user: %+v", message)
	return message, err
}

func (adapter MessagesLoggingAdapter) GetByIdsForUser(messageIds []int, userId int) []messages.Message {
	log.Printf("fetching messages by ids for user: messageIds=%v, userId=%d", messageIds, userId)
	messages := adapter.adapter.GetByIdsForUser(messageIds, userId)
	log.Printf("fetched messages: %+v", messages)
	return messages
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

func (adapter MessagesAdapter) getChatAllForUserTotal(chatId int, userId int) int {
	var count int64

	adapter.db.Model(&Message{}).Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatId, userId,
	).Count(&count)

	return int(count)
}

func (adapter MessagesAdapter) GetChatAllForUser(chatId int, userId int, offset int, limit int) utils.OffsetResponse[messages.Message] {
	var dbMessages []Message

	total := adapter.getChatAllForUserTotal(chatId, userId)

	adapter.db.Preload("Chat").Preload("Reactions").Preload("Voice").Preload("Circle").Preload("Attachments").Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members)", chatId, userId,
	).Order(
		"messages.created_at DESC NULLS LAST",
	).Offset(offset).Limit(limit).Find(&dbMessages)

	var messagesModels []messages.Message
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

func (adapter MessagesAdapter) getMessageOffsetById(chatId int, userId int, messageId int) int {
	var offset int64

	adapter.db.Model(&Message{}).Joins("JOIN chats ON messages.chat_id = chats.id").Where(
		"messages.chat_id = ? AND ? = ANY(chats.members) AND messages.id >= ?", chatId, userId, messageId,
	).Order("messages.created_at DESC NULLS LAST").Count(&offset)

	return int(offset)
}

func (adapter MessagesAdapter) GetChatCursorAllForUser(chatId int, userId int, messageId int, aroundOffset int) utils.OffsetResponse[messages.Message] {
	offset := adapter.getMessageOffsetById(chatId, userId, messageId)

	startOffset := offset - aroundOffset
	if startOffset < 0 {
		startOffset = 0
	}

	return adapter.GetChatAllForUser(chatId, userId, startOffset, aroundOffset*2)
}

func (adapter MessagesAdapter) GetChatsLast(chatIds []int, userId int) []messages.Message {
	var messages []messages.Message

	for _, chatId := range chatIds {
		var message Message

		adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
			"messages.chat_id = ? AND ? = ANY(chats.members)", chatId, userId,
		).Order("messages.created_at DESC NULLS LAST").Limit(1).First(&message)

		messageModel := DbMessageToModel(message)
		messages = append(messages, messageModel)
	}

	return messages
}

func (adapter MessagesAdapter) getByid(messageId int) (*messages.Message, error) {
	var dbMessage Message

	result := adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id = ?", messageId,
	).First(&dbMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	messageModel := DbMessageToModel(dbMessage)
	return &messageModel, nil
}

func (adapter MessagesAdapter) GetByIdForUser(messageId int, userId int) (*messages.Message, error) {
	var dbMessage Message

	result := adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id = ? AND ? = ANY(chats.members)", messageId, userId,
	).First(&dbMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	messageModel := DbMessageToModel(dbMessage)
	return &messageModel, nil
}

func (adapter MessagesAdapter) GetByIdsForUser(messageIds []int, userId int) []messages.Message {
	var dbMessages []Message

	adapter.db.Preload("Chat").Preload("Voice").Preload("Circle").Preload("Attachments").Preload("Reactions").Joins("JOIN chats ON messages.chat_id = chats.id").Preload("Circle").Preload("Voice").Preload("Attachments").Where(
		"messages.id IN ? AND ? = ANY(chats.members)", messageIds, userId,
	).Find(&dbMessages)

	var modelMessages []messages.Message
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

	var attachments []SavedFile
	for _, attachment := range message.GetAttachments() {
		attachments = append(attachments, GetOrCreateFile(&attachment, adapter.db))
	}

	var reactions []Reaction
	for _, reaction := range message.GetReactions() {
		reactions = append(reactions, adapter.getOrCreateReaction(reaction))
	}

	dbMessage := ModelToDbMessage(message, voicePointer, circlePointer, attachments, reactions)
	result := adapter.db.Save(&dbMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	savedMessage, err := adapter.getByid(int(dbMessage.ID))
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
