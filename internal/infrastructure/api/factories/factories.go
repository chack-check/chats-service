package factories

import (
	"time"

	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/dtos"
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/api/graph/model"
)

func UploadingFileMetaToModel(meta model.UploadingFileMeta) dtos.UploadingFileMeta {
	return dtos.NewUploadingFileMeta(
		meta.URL,
		meta.Filename,
		meta.Signature,
		constants.SystemFiletype(meta.SystemFiletype.String()),
	)
}

func UploadingFileToModel(file model.UploadingFile) dtos.UploadingFile {
	original := UploadingFileMetaToModel(*file.Original)
	var converted *dtos.UploadingFileMeta
	if convertedMeta := file.Converted; convertedMeta != nil {
		convertedFile := UploadingFileMetaToModel(*convertedMeta)
		converted = &convertedFile
	}

	return dtos.NewUploadingFile(
		original,
		converted,
	)
}

func CreateMessageRequestToModel(request model.CreateMessageRequest) dtos.CreateMessageData {
	var voice *dtos.UploadingFile
	if request.Voice != nil {
		file := UploadingFileToModel(*request.Voice)
		voice = &file
	}

	var circle *dtos.UploadingFile
	if request.Circle != nil {
		file := UploadingFileToModel(*request.Circle)
		circle = &file
	}

	attachments := make([]dtos.UploadingFile, 0, len(request.Attachments))
	for _, attachment := range request.Attachments {
		file := UploadingFileToModel(*attachment)
		attachments = append(attachments, file)
	}

	return dtos.NewCreateMessageData(
		request.ChatID,
		constants.MessageTypes(request.Type),
		request.Content,
		voice,
		attachments,
		request.ReplyToID,
		request.Mentioned,
		circle,
	)
}

func UpdateMessageRequestToModel(request model.ChangeMessageRequest) dtos.UpdateMessageData {
	mentioned := make([]int, 0, len(request.Mentioned))
	for _, user := range request.Mentioned {
		mentioned = append(mentioned, *user)
	}

	attachments := make([]dtos.UploadingFile, 0, len(request.Attachments))
	for _, attachment := range request.Attachments {
		file := UploadingFileToModel(*attachment)
		attachments = append(attachments, file)
	}

	return dtos.NewUpdateMessageData(
		request.Content,
		attachments,
		mentioned,
	)
}

func SavedFileToResponse(file entities.SavedFile) model.SavedFile {
	return model.SavedFile{
		OriginalURL:       file.GetOriginalURL(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedURL:      file.GetConvertedURL(),
		ConvertedFilename: file.GetConvertedFilename(),
	}
}

func ReactionModelToResponse(reaction entities.MessageReaction) model.Reaction {
	return model.Reaction{
		UserID:  reaction.GetUserID(),
		Content: reaction.GetContent(),
	}
}

func MessageModelToResponse(message entities.Message) model.Message {
	chat := message.GetChat()
	var voice *model.SavedFile
	if message.GetVoice() != nil {
		file := SavedFileToResponse(*message.GetVoice())
		voice = &file
	}

	var circle *model.SavedFile
	if message.GetCircle() != nil {
		file := SavedFileToResponse(*message.GetCircle())
		circle = &file
	}

	attachments := make([]*model.SavedFile, 0, len(message.GetAttachments()))
	for _, attachment := range message.GetAttachments() {
		file := SavedFileToResponse(attachment)
		attachments = append(attachments, &file)
	}

	reactions := make([]*model.Reaction, 0, len(message.GetReactions()))
	for _, reaction := range message.GetReactions() {
		reactionResponse := ReactionModelToResponse(reaction)
		reactions = append(reactions, &reactionResponse)
	}

	return model.Message{
		ID:          message.GetID(),
		Type:        model.MessageType(string(message.GetType())),
		SenderID:    message.GetSenderID(),
		ChatID:      chat.GetID(),
		Content:     message.GetContent(),
		Voice:       voice,
		Circle:      circle,
		ReplyToID:   message.GetReplyToID(),
		ReadedBy:    message.GetReadedBy(),
		Reactions:   reactions,
		Attachments: attachments,
		Mentioned:   message.GetMentioned(),
		CreatedAt:   message.GetCreatedAt().Format(time.RFC3339),
	}
}

func OffsetMessagesToResponse(messages dtos.OffsetResponse[entities.Message], chatID int) model.PaginatedMessages {
	data := messages.GetData()
	messagesResponse := make([]*model.Message, 0, len(data))
	for _, message := range data {
		messageResponse := MessageModelToResponse(message)
		messagesResponse = append(messagesResponse, &messageResponse)
	}

	return model.PaginatedMessages{
		Offset: messages.GetOffset(),
		Limit:  messages.GetLimit(),
		Total:  messages.GetTotal(),
		Data:   messagesResponse,
		ID:     chatID,
	}
}

func CreateChatRequestToModel(request model.CreateChatRequest, chatType constants.ChatTypes) dtos.CreateChatData {
	var avatar *dtos.UploadingFile
	if request.Avatar != nil {
		file := UploadingFileToModel(*request.Avatar)
		avatar = &file
	}

	return dtos.NewCreateChatData(
		chatType,
		avatar,
		request.Title,
		request.Members,
		request.User,
	)
}

func ActionUserModelToResponse(user entities.ActionUser) model.ChatActionUser {
	return model.ChatActionUser{
		FullName: user.GetFullName(),
		ID:       user.GetID(),
	}
}

func ChatModelToResponse(chat entities.Chat) model.Chat {
	var avatar *model.SavedFile
	if chatAvatar := chat.GetAvatar(); chatAvatar != nil {
		file := SavedFileToResponse(*chatAvatar)
		avatar = &file
	}

	actions := make([]*model.ChatAction, 0, len(chat.GetActions()))
	for key, value := range chat.GetActions() {
		actionUsers := make([]*model.ChatActionUser, 0, len(value))
		for _, user := range value {
			userResponse := ActionUserModelToResponse(user)
			actionUsers = append(actionUsers, &userResponse)
		}

		action := model.ChatAction{
			Action:      model.ActionTypes(key),
			ActionUsers: actionUsers,
		}
		actions = append(actions, &action)
	}

	return model.Chat{
		ID:         chat.GetID(),
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       model.ChatType(string(chat.GetType())),
		Members:    chat.GetMembers(),
		IsArchived: chat.GetIsArchived(),
		OwnerID:    chat.GetOwnerID(),
		Admins:     chat.GetAdmins(),
		Actions:    actions,
	}
}

func PaginatedChatsToResponse(chats dtos.PaginatedResponse[entities.Chat]) model.PaginatedChats {
	data := chats.GetData()
	chatsResponse := make([]*model.Chat, 0, len(data))
	for _, chat := range data {
		chatResponse := ChatModelToResponse(chat)
		chatsResponse = append(chatsResponse, &chatResponse)
	}

	return model.PaginatedChats{
		Page:     chats.GetPage(),
		NumPages: chats.GetPagesCount(),
		PerPage:  chats.GetPerPage(),
		Total:    chats.GetTotal(),
		Data:     chatsResponse,
	}
}
