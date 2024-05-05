package chats

import (
	"slices"

	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/users"
)

type ActionTypes string

const (
	WritingActionType         = "writing"
	AudioRecordingActionType  = "audio_recording"
	AudioSendingActionType    = "audio_sending"
	CircleRecordingActionType = "circle_recording"
	CircleSendingActionType   = "circle_sending"
	FilesSendingActionType    = "files_sending"
)

type ChatTypes string

var (
	UserChatType          ChatTypes = "user"
	GroupChatType         ChatTypes = "group"
	SavedMessagesChatType ChatTypes = "saved_messages"
)

type Chat struct {
	id         int
	avatar     *files.SavedFile
	title      string
	type_      ChatTypes
	members    []int
	isArchived bool
	ownerId    int
	admins     []int
	actions    map[ActionTypes][]users.ActionUser
}

func (model *Chat) GetId() int {
	return model.id
}

func (model *Chat) GetAvatar() *files.SavedFile {
	return model.avatar
}

func (model *Chat) SetAvatar(avatar files.SavedFile) {
	model.avatar = &avatar
}

func (model *Chat) GetTitle() string {
	return model.title
}

func (model *Chat) SetTitle(title string) {
	model.title = title
}

func (model *Chat) GetType() ChatTypes {
	return model.type_
}

func (model *Chat) SetType(type_ ChatTypes) {
	model.type_ = type_
}

func (model *Chat) GetMembers() []int {
	return model.members
}

func (model *Chat) SetMembers(members []int) {
	model.members = members
}

func (model *Chat) GetIsArchived() bool {
	return model.isArchived
}

func (model *Chat) Archive() {
	model.isArchived = true
}

func (model *Chat) Unarchive() {
	model.isArchived = false
}

func (model *Chat) GetOwnerId() int {
	return model.ownerId
}

func (model *Chat) SetOwnerId(ownerId int) {
	var newMembers []int
	for _, member := range model.members {
		if member == model.ownerId {
			continue
		}

		newMembers = append(newMembers, member)
	}

	var newAdmins []int
	for _, admin := range model.admins {
		if admin == model.ownerId {
			continue
		}

		newAdmins = append(newAdmins, admin)
	}

	model.ownerId = ownerId
	model.members = newMembers
	model.admins = newAdmins
}

func (model *Chat) GetAdmins() []int {
	return model.admins
}

func (model *Chat) SetAdmins(admins []int) {
	if !slices.Contains(admins, model.ownerId) {
		admins = append(admins, model.ownerId)
	}

	model.admins = admins
}

func (model *Chat) SetupUserData(anotherUser *users.User) {
	if model.type_ != "user" || anotherUser == nil {
		return
	}

	model.SetTitle(anotherUser.GetFullName())
	if anotherUser.GetAvatar() != nil {
		model.SetAvatar(*anotherUser.GetAvatar())
	}
}

func (model *Chat) GetActions() map[ActionTypes][]users.ActionUser {
	return model.actions
}

func (model *Chat) SetupActions(actions map[ActionTypes][]users.ActionUser) {
	model.actions = actions
}

type CreateChatData struct {
	avatar     *files.UploadingFile
	title      *string
	membersIds []int
	userId     *int
	type_      ChatTypes
}

func (data *CreateChatData) GetAvatar() *files.UploadingFile {
	return data.avatar
}

func (data *CreateChatData) GetTitle() *string {
	return data.title
}

func (data *CreateChatData) GetMembersIds() []int {
	return data.membersIds
}

func (data *CreateChatData) GetUserId() *int {
	return data.userId
}

func (data *CreateChatData) GetType() ChatTypes {
	return data.type_
}

func NewChat(id int, avatar *files.SavedFile, title string, type_ ChatTypes, members []int, isArchived bool, ownerId int, admins []int) Chat {
	return Chat{
		id:         id,
		avatar:     avatar,
		title:      title,
		type_:      type_,
		members:    members,
		isArchived: isArchived,
		ownerId:    ownerId,
		admins:     admins,
	}
}

func NewCreateChatData(
	chatType ChatTypes,
	avatar *files.UploadingFile,
	title *string,
	membersIds []int,
	userId *int,
) CreateChatData {
	return CreateChatData{
		type_:      chatType,
		avatar:     avatar,
		title:      title,
		membersIds: membersIds,
		userId:     userId,
	}
}
