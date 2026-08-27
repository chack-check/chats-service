package database

import (
	"time"

	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"

	"github.com/lib/pq"
)

func DbSavedFileToModel(file SavedFile) entities.SavedFile {
	var convertedURL *string
	var convertedFilename *string
	if file.ConvertedUrl == "" {
		convertedURL = nil
		convertedFilename = nil
	} else {
		convertedURL = &file.ConvertedUrl
		convertedFilename = &file.ConvertedFilename
	}

	return entities.NewSavedFile(
		file.OriginalUrl,
		file.OriginalFilename,
		convertedURL,
		convertedFilename,
	)
}

func ModelToDbSavedFile(file entities.SavedFile) SavedFile {
	var convertedURL string
	var convertedFilename string
	if url := file.GetConvertedURL(); url != nil {
		convertedURL = *url
		filename := file.GetConvertedFilename()
		convertedFilename = *filename
	} else {
		convertedURL = ""
		convertedFilename = ""
	}

	return SavedFile{
		OriginalUrl:       file.GetOriginalURL(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedUrl:      convertedURL,
		ConvertedFilename: convertedFilename,
	}
}

func DbChatToModel(chat Chat) entities.Chat {
	var avatar *entities.SavedFile
	if chat.AvatarId != nil {
		savedFile := DbSavedFileToModel(chat.Avatar)
		avatar = &savedFile
	}

	members := make([]int, 0, len(chat.Members))
	for _, member := range chat.Members {
		members = append(members, int(member))
	}

	admins := make([]int, 0, len(chat.Admins))
	for _, admin := range chat.Admins {
		admins = append(admins, int(admin))
	}

	return entities.NewChat(
		int(chat.ID),
		avatar,
		chat.Title,
		constants.ChatTypes(chat.Type),
		members,
		chat.IsArchived,
		int(chat.OwnerId),
		admins,
	)
}

func ModelToDbChat(chat entities.Chat, avatar SavedFile) Chat {
	var avatarID *uint
	if avatar.OriginalUrl != "" {
		id := avatar.ID
		avatarID = &id
	}

	members := make(pq.Int64Array, 0, len(chat.GetMembers()))
	for _, member := range chat.GetMembers() {
		members = append(members, int64(member))
	}

	admins := make(pq.Int64Array, 0, len(chat.GetAdmins()))
	for _, admin := range chat.GetAdmins() {
		admins = append(admins, int64(admin))
	}

	return Chat{
		ID:         uint(chat.GetID()),
		AvatarId:   avatarID,
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       string(chat.GetType()),
		Members:    members,
		IsArchived: chat.GetIsArchived(),
		OwnerId:    uint(chat.GetOwnerID()),
		Admins:     admins,
	}
}

func DbMessageReactionToModel(reaction Reaction) entities.MessageReaction {
	return entities.NewMessageReaction(
		int(reaction.UserId),
		reaction.Content,
	)
}

func DbMessageToModel(message Message) entities.Message {
	var voice *entities.SavedFile
	if message.Voice == nil {
		voice = nil
	} else {
		savedFile := DbSavedFileToModel(*message.Voice)
		voice = &savedFile
	}

	var circle *entities.SavedFile
	if message.Circle == nil {
		circle = nil
	} else {
		savedFile := DbSavedFileToModel(*message.Circle)
		circle = &savedFile
	}

	attachments := make([]entities.SavedFile, 0, len(message.Attachments))
	for _, attachment := range message.Attachments {
		attachments = append(attachments, DbSavedFileToModel(attachment))
	}

	var replyToID *int
	if message.ReplyToID != 0 {
		replyToIDInt := int(message.ReplyToID)
		replyToID = &replyToIDInt
	}

	mentioned := make([]int, 0, len(message.Mentioned))
	for _, ment := range message.Mentioned {
		mentioned = append(mentioned, int(ment))
	}

	readedBy := make([]int, 0, len(message.ReadedBy))
	for _, user := range message.ReadedBy {
		readedBy = append(readedBy, int(user))
	}

	deletedFor := make([]int, 0, len(message.DeletedFor))
	for _, user := range message.DeletedFor {
		deletedFor = append(deletedFor, int(user))
	}

	reactions := make([]entities.MessageReaction, 0, len(message.Reactions))
	for _, reaction := range message.Reactions {
		reactionModel := DbMessageReactionToModel(reaction)
		reactions = append(reactions, reactionModel)
	}

	return entities.NewMessage(
		int(message.ID),
		int(message.SenderId),
		DbChatToModel(message.Chat),
		constants.MessageTypes(message.Type),
		&message.Content,
		voice,
		circle,
		attachments,
		replyToID,
		mentioned,
		readedBy,
		reactions,
		deletedFor,
		&message.CreatedAt,
	)
}

func ModelToDbMessage(message entities.Message, voice *SavedFile, circle *SavedFile, attachments []SavedFile, reactions []Reaction) Message {
	chat := message.GetChat()

	var content string
	if messageContent := message.GetContent(); messageContent != nil {
		content = *messageContent
	}

	var replyToID int
	if messageReplyToID := message.GetReplyToID(); messageReplyToID != nil {
		replyToID = *messageReplyToID
	}

	mentioned := make(pq.Int32Array, 0, len(message.GetMentioned()))
	for _, ment := range message.GetMentioned() {
		mentioned = append(mentioned, int32(ment))
	}

	readedBy := make(pq.Int32Array, 0, len(message.GetReadedBy()))
	for _, reader := range message.GetReadedBy() {
		readedBy = append(readedBy, int32(reader))
	}

	var createdAt time.Time
	if dt := message.GetCreatedAt(); dt != nil {
		createdAt = *dt
	}

	return Message{
		ID:          uint(message.GetID()),
		SenderId:    uint(message.GetSenderID()),
		ChatId:      uint(chat.GetID()),
		Type:        string(message.GetType()),
		Content:     content,
		Voice:       voice,
		Circle:      circle,
		Attachments: attachments,
		ReplyToID:   uint(replyToID),
		Mentioned:   mentioned,
		ReadedBy:    readedBy,
		Reactions:   reactions,
		CreatedAt:   createdAt,
	}
}
