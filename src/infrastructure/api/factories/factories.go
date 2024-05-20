package factories

import (
	"time"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/messages"
	"github.com/chack-check/chats-service/domain/users"
	"github.com/chack-check/chats-service/domain/utils"
	"github.com/chack-check/chats-service/infrastructure/api/graph/model"
)

func UploadingFileMetaToModel(meta model.UploadingFileMeta) files.UploadingFileMeta {
	return files.NewUploadingFileMeta(
		meta.URL,
		meta.Filename,
		meta.Signature,
		files.SystemFiletype(meta.SystemFiletype.String()),
	)
}

func UploadingFileToModel(file model.UploadingFile) files.UploadingFile {
	original := UploadingFileMetaToModel(*file.Original)
	var converted *files.UploadingFileMeta
	if convertedMeta := file.Converted; convertedMeta != nil {
		convertedFile := UploadingFileMetaToModel(*convertedMeta)
		converted = &convertedFile
	}

	return files.NewUploadingFile(
		original,
		converted,
	)
}

func CreateMessageRequestToModel(request model.CreateMessageRequest) messages.CreateMessageData {
	var voice *files.UploadingFile
	if request.Voice != nil {
		file := UploadingFileToModel(*request.Voice)
		voice = &file
	}

	var circle *files.UploadingFile
	if request.Circle != nil {
		file := UploadingFileToModel(*request.Circle)
		circle = &file
	}

	var attachments []files.UploadingFile
	for _, attachment := range request.Attachments {
		file := UploadingFileToModel(*attachment)
		attachments = append(attachments, file)
	}

	return messages.NewCreateMessageData(
		request.ChatID,
		messages.MessageTypes(request.Type),
		request.Content,
		voice,
		attachments,
		request.ReplyToID,
		request.Mentioned,
		circle,
	)
}

func UpdateMessageRequestToModel(request model.ChangeMessageRequest) messages.UpdateMessageData {
	var mentioned []int
	for _, user := range request.Mentioned {
		mentioned = append(mentioned, *user)
	}

	var attachments []files.UploadingFile
	for _, attachment := range request.Attachments {
		file := UploadingFileToModel(*attachment)
		attachments = append(attachments, file)
	}

	return messages.NewUpdateMessageData(
		request.Content,
		attachments,
		mentioned,
	)
}

func SavedFileToResponse(file files.SavedFile) model.SavedFile {
	return model.SavedFile{
		OriginalURL:       file.GetOriginalUrl(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedURL:      file.GetConvertedUrl(),
		ConvertedFilename: file.GetConvertedFilename(),
	}
}

func ReactionModelToResponse(reaction messages.MessageReaction) model.Reaction {
	return model.Reaction{
		UserID:  reaction.GetUserId(),
		Content: reaction.GetContent(),
	}
}

func MessageModelToResponse(message messages.Message) model.Message {
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

	var attachments []*model.SavedFile
	for _, attachment := range message.GetAttachments() {
		file := SavedFileToResponse(attachment)
		attachments = append(attachments, &file)
	}

	var reactions []*model.Reaction
	for _, reaction := range message.GetReactions() {
		reactionResponse := ReactionModelToResponse(reaction)
		reactions = append(reactions, &reactionResponse)
	}

	return model.Message{
		ID:          message.GetId(),
		Type:        model.MessageType(string(message.GetType())),
		SenderID:    message.GetSenderId(),
		ChatID:      chat.GetId(),
		Content:     message.GetContent(),
		Voice:       voice,
		Circle:      circle,
		ReplyToID:   message.GetReplyToId(),
		ReadedBy:    message.GetReadedBy(),
		Reactions:   reactions,
		Attachments: attachments,
		Mentioned:   message.GetMentioned(),
		CreatedAt:   message.GetCreatedAt().Format(time.RFC3339),
	}
}

func OffsetMessagesToResponse(messages utils.OffsetResponse[messages.Message], chatId int) model.PaginatedMessages {
	data := messages.GetData()
	var messagesResponse []*model.Message
	for _, message := range data {
		messageResponse := MessageModelToResponse(message)
		messagesResponse = append(messagesResponse, &messageResponse)
	}

	return model.PaginatedMessages{
		Offset: messages.GetOffset(),
		Limit:  messages.GetLimit(),
		Total:  messages.GetTotal(),
		Data:   messagesResponse,
		ID:     chatId,
	}
}

func CreateChatRequestToModel(request model.CreateChatRequest, chatType chats.ChatTypes) chats.CreateChatData {
	var avatar *files.UploadingFile
	if request.Avatar != nil {
		file := UploadingFileToModel(*request.Avatar)
		avatar = &file
	}

	return chats.NewCreateChatData(
		chatType,
		avatar,
		request.Title,
		request.Members,
		request.User,
	)
}

func ActionUserModelToResponse(user users.ActionUser) model.ChatActionUser {
	return model.ChatActionUser{
		FullName: user.GetFullName(),
		ID:       user.GetId(),
	}
}

func ChatModelToResponse(chat chats.Chat) model.Chat {
	var avatar *model.SavedFile
	if chatAvatar := chat.GetAvatar(); chatAvatar != nil {
		file := SavedFileToResponse(*chatAvatar)
		avatar = &file
	}

	var actions []*model.ChatAction
	for key, value := range chat.GetActions() {
		var actionUsers []*model.ChatActionUser
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
		ID:         chat.GetId(),
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       model.ChatType(string(chat.GetType())),
		Members:    chat.GetMembers(),
		IsArchived: chat.GetIsArchived(),
		OwnerID:    chat.GetOwnerId(),
		Admins:     chat.GetAdmins(),
		Actions:    actions,
	}
}

func PaginatedChatsToResponse(chats utils.PaginatedResponse[chats.Chat]) model.PaginatedChats {
	data := chats.GetData()
	var chatsResponse []*model.Chat
	for _, chat := range data {
		chatResponse := ChatModelToResponse(chat)
		chatsResponse = append(chatsResponse, &chatResponse)
	}

	return model.PaginatedChats{
		Page:       chats.GetPage(),
		NumPages:   chats.GetPagesCount(),
		PagesCount: chats.GetPagesCount(),
		Total:      chats.GetTotal(),
		Data:       chatsResponse,
	}
}
