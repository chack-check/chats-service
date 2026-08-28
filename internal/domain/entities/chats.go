package entities

import (
	"chats-service/internal/domain/constants"
	"slices"
)

type Chat struct {
	id         int
	avatar     *SavedFile
	title      string
	type_      constants.ChatTypes
	members    []int
	isArchived bool
	ownerID    int
	admins     []int
	actions    map[constants.ActionTypes][]ActionUser
}

func NewChat(id int, avatar *SavedFile, title string, type_ constants.ChatTypes, members []int, isArchived bool, ownerID int, admins []int) Chat {
	return Chat{
		id:         id,
		avatar:     avatar,
		title:      title,
		type_:      type_,
		members:    members,
		isArchived: isArchived,
		ownerID:    ownerID,
		admins:     admins,
	}
}

func (chat *Chat) GetID() int {
	return chat.id
}

func (chat *Chat) GetAvatar() *SavedFile {
	return chat.avatar
}

func (chat *Chat) SetAvatar(avatar SavedFile) {
	chat.avatar = &avatar
}

func (chat *Chat) GetTitle() string {
	return chat.title
}

func (chat *Chat) SetTitle(title string) {
	chat.title = title
}

func (chat *Chat) GetType() constants.ChatTypes {
	return chat.type_
}

func (chat *Chat) SetType(type_ constants.ChatTypes) {
	chat.type_ = type_
}

func (chat *Chat) GetMembers() []int {
	return chat.members
}

func (chat *Chat) SetMembers(members []int) {
	chat.members = members
}

func (chat *Chat) GetIsArchived() bool {
	return chat.isArchived
}

func (chat *Chat) Archive() {
	chat.isArchived = true
}

func (chat *Chat) Unarchive() {
	chat.isArchived = false
}

func (chat *Chat) GetOwnerID() int {
	return chat.ownerID
}

func (chat *Chat) SetOwnerID(ownerID int) {
	var newMembers []int
	for _, member := range chat.members {
		if member == chat.ownerID {
			continue
		}

		newMembers = append(newMembers, member)
	}

	var newAdmins []int
	for _, admin := range chat.admins {
		if admin == chat.ownerID {
			continue
		}

		newAdmins = append(newAdmins, admin)
	}

	chat.ownerID = ownerID
	chat.members = newMembers
	chat.admins = newAdmins
}

func (chat *Chat) GetAdmins() []int {
	return chat.admins
}

func (chat *Chat) SetAdmins(admins []int) {
	if !slices.Contains(admins, chat.ownerID) {
		admins = append(admins, chat.ownerID)
	}

	chat.admins = admins
}

func (chat *Chat) SetupUserData(anotherUser *User) {
	if chat.type_ != "user" || anotherUser == nil {
		return
	}

	chat.SetTitle(anotherUser.GetFullName())
	if anotherUser.GetAvatar() != nil {
		chat.SetAvatar(*anotherUser.GetAvatar())
	}
}

func (chat *Chat) GetActions() map[constants.ActionTypes][]ActionUser {
	return chat.actions
}

func (chat *Chat) SetupActions(actions map[constants.ActionTypes][]ActionUser) {
	chat.actions = actions
}
