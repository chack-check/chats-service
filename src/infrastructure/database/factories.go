package database

import (
	"time"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/messages"
	"github.com/lib/pq"
)

func DbSavedFileToModel(file SavedFile) files.SavedFile {
	var convertedUrl *string
	var convertedFilename *string
	if file.ConvertedUrl == "" {
		convertedUrl = nil
		convertedFilename = nil
	} else {
		convertedUrl = &file.ConvertedUrl
		convertedFilename = &file.ConvertedFilename
	}

	return files.NewSavedFile(
		file.OriginalUrl,
		file.OriginalFilename,
		convertedUrl,
		convertedFilename,
	)
}

func ModelToDbSavedFile(file files.SavedFile) SavedFile {
	var convertedUrl string
	var convertedFilename string
	if url := file.GetConvertedUrl(); url != nil {
		convertedUrl = *url
		filename := file.GetConvertedFilename()
		convertedFilename = *filename
	} else {
		convertedUrl = ""
		convertedFilename = ""
	}

	return SavedFile{
		OriginalUrl:       file.GetOriginalUrl(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedUrl:      convertedUrl,
		ConvertedFilename: convertedFilename,
	}
}

func DbChatToModel(chat Chat) chats.Chat {
	var avatar *files.SavedFile
	if chat.AvatarId != nil {
		savedFile := DbSavedFileToModel(chat.Avatar)
		avatar = &savedFile
	}

	var members []int
	for _, member := range chat.Members {
		members = append(members, int(member))
	}

	var admins []int
	for _, admin := range chat.Admins {
		admins = append(admins, int(admin))
	}

	return chats.NewChat(
		int(chat.ID),
		avatar,
		chat.Title,
		chats.ChatTypes(chat.Type),
		members,
		chat.IsArchived,
		int(chat.OwnerId),
		admins,
	)
}

func ModelToDbChat(chat chats.Chat, avatar SavedFile) Chat {
	var avatarId *uint
	if avatar.OriginalUrl != "" {
		id := avatar.ID
		avatarId = &id
	}

	var members pq.Int64Array
	for _, member := range chat.GetMembers() {
		members = append(members, int64(member))
	}

	var admins pq.Int64Array
	for _, admin := range chat.GetAdmins() {
		admins = append(admins, int64(admin))
	}

	return Chat{
		ID:         uint(chat.GetId()),
		AvatarId:   avatarId,
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       string(chat.GetType()),
		Members:    members,
		IsArchived: chat.GetIsArchived(),
		OwnerId:    uint(chat.GetOwnerId()),
		Admins:     admins,
	}
}

func DbMessageReactionToModel(reaction Reaction) messages.MessageReaction {
	return messages.NewMessageReaction(
		int(reaction.UserId),
		reaction.Content,
	)
}

func DbMessageToModel(message Message) messages.Message {
	var voice *files.SavedFile
	if message.Voice == nil {
		voice = nil
	} else {
		savedFile := DbSavedFileToModel(*message.Voice)
		voice = &savedFile
	}

	var circle *files.SavedFile
	if message.Circle == nil {
		circle = nil
	} else {
		savedFile := DbSavedFileToModel(*message.Circle)
		circle = &savedFile
	}

	var attachments []files.SavedFile
	for _, attachment := range message.Attachments {
		attachments = append(attachments, DbSavedFileToModel(attachment))
	}

	var replyToId *int
	if message.ReplyToID != 0 {
		replyToIdInt := int(message.ReplyToID)
		replyToId = &replyToIdInt
	}

	var mentioned []int
	for _, ment := range message.Mentioned {
		mentioned = append(mentioned, int(ment))
	}

	var readedBy []int
	for _, user := range message.ReadedBy {
		readedBy = append(readedBy, int(user))
	}

	var deletedFor []int
	for _, user := range message.DeletedFor {
		deletedFor = append(deletedFor, int(user))
	}

	var reactions []messages.MessageReaction
	for _, reaction := range message.Reactions {
		reactionModel := DbMessageReactionToModel(reaction)
		reactions = append(reactions, reactionModel)
	}

	return messages.NewMessage(
		int(message.ID),
		int(message.SenderId),
		DbChatToModel(message.Chat),
		messages.MessageTypes(message.Type),
		&message.Content,
		voice,
		circle,
		attachments,
		replyToId,
		mentioned,
		readedBy,
		reactions,
		deletedFor,
		&message.CreatedAt,
	)
}

func ModelToDbMessage(message messages.Message, voiceId *int, circleId *int, attachments []SavedFile, reactions []Reaction) Message {
	chat := message.GetChat()

	var content string
	if messageContent := message.GetContent(); messageContent != nil {
		content = *messageContent
	}

	var replyToId int
	if messageReplyToId := message.GetReplyToId(); messageReplyToId != nil {
		replyToId = *messageReplyToId
	}

	var mentioned pq.Int32Array
	for _, ment := range message.GetMentioned() {
		mentioned = append(mentioned, int32(ment))
	}

	var readedBy pq.Int32Array
	for _, reader := range message.GetReadedBy() {
		readedBy = append(readedBy, int32(reader))
	}

	var createdAt time.Time
	if dt := message.GetCreatedAt(); dt != nil {
		createdAt = *dt
	}

	return Message{
		ID:          uint(message.GetId()),
		SenderId:    uint(message.GetSenderId()),
		ChatId:      uint(chat.GetId()),
		Type:        string(message.GetType()),
		Content:     content,
		VoiceId:     voiceId,
		CircleId:    circleId,
		Attachments: attachments,
		ReplyToID:   uint(replyToId),
		Mentioned:   mentioned,
		ReadedBy:    readedBy,
		Reactions:   reactions,
		CreatedAt:   createdAt,
	}
}
