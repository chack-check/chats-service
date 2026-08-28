package dtos

import (
	"chats-service/internal/domain/constants"
)

type CreateMessageData struct {
	chatID      int
	type_       constants.MessageTypes
	content     *string
	voice       *UploadingFile
	attachments []UploadingFile
	replyToID   *int
	mentioned   []int
	circle      *UploadingFile
}

func NewCreateMessageData(
	chatID int,
	type_ constants.MessageTypes,
	content *string,
	voice *UploadingFile,
	attachments []UploadingFile,
	replyToID *int,
	mentioned []int,
	circle *UploadingFile,
) CreateMessageData {
	return CreateMessageData{
		chatID:      chatID,
		type_:       type_,
		content:     content,
		voice:       voice,
		attachments: attachments,
		replyToID:   replyToID,
		mentioned:   mentioned,
		circle:      circle,
	}
}

func (dto *CreateMessageData) GetChatID() int {
	return dto.chatID
}

func (dto *CreateMessageData) GetType() constants.MessageTypes {
	return dto.type_
}

func (dto *CreateMessageData) GetContent() *string {
	return dto.content
}

func (dto *CreateMessageData) GetVoice() *UploadingFile {
	return dto.voice
}

func (dto *CreateMessageData) GetAttachments() []UploadingFile {
	return dto.attachments
}

func (dto *CreateMessageData) GetReplyToID() *int {
	return dto.replyToID
}

func (dto *CreateMessageData) GetMentioned() []int {
	return dto.mentioned
}

func (dto *CreateMessageData) GetCircle() *UploadingFile {
	return dto.circle
}

type UpdateMessageData struct {
	content     *string
	attachments []UploadingFile
	mentioned   []int
}

func NewUpdateMessageData(
	content *string,
	attachments []UploadingFile,
	mentioned []int,
) UpdateMessageData {
	return UpdateMessageData{
		content:     content,
		attachments: attachments,
		mentioned:   mentioned,
	}
}

func (dto *UpdateMessageData) GetContent() *string {
	return dto.content
}

func (dto *UpdateMessageData) GetAttachments() []UploadingFile {
	return dto.attachments
}

func (dto *UpdateMessageData) GetMentioned() []int {
	return dto.mentioned
}
