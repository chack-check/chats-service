package dtos

import (
	"chats-service/internal/domain/constants"
)

type ChangeGroupChatData struct {
	title *string
}

func NewChangeGroupChatData(title *string) ChangeGroupChatData {
	return ChangeGroupChatData{
		title: title,
	}
}

func (dto *ChangeGroupChatData) GetTitle() *string {
	return dto.title
}

type CreateChatData struct {
	avatar     *UploadingFile
	title      *string
	membersIds []int
	userID     *int
	type_      constants.ChatTypes
}

func NewCreateChatData(
	chatType constants.ChatTypes,
	avatar *UploadingFile,
	title *string,
	membersIds []int,
	userID *int,
) CreateChatData {
	return CreateChatData{
		type_:      chatType,
		avatar:     avatar,
		title:      title,
		membersIds: membersIds,
		userID:     userID,
	}
}

func (dto *CreateChatData) GetAvatar() *UploadingFile {
	return dto.avatar
}

func (dto *CreateChatData) GetTitle() *string {
	return dto.title
}

func (dto *CreateChatData) GetMembersIds() []int {
	return dto.membersIds
}

func (dto *CreateChatData) SetMembersIds(newMembers []int) {
	dto.membersIds = newMembers
}

func (dto *CreateChatData) GetUserID() *int {
	return dto.userID
}

func (dto *CreateChatData) GetType() constants.ChatTypes {
	return dto.type_
}
