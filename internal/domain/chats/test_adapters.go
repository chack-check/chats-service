package chats

import (
	"slices"

	"chats-service/internal/domain/files"
	"chats-service/internal/domain/users"
	"chats-service/internal/domain/utils"
)

var existingAvatars = []files.SavedFile{
	files.NewSavedFile(
		"original_url_1",
		"original_filename_1",
		nil,
		nil,
	),
}

var deletedChats = []Chat{
	{
		id:         3,
		avatar:     nil,
		title:      "",
		type_:      UserChatType,
		members:    []int{1, 2},
		isArchived: false,
		ownerId:    0,
		admins:     []int{},
		actions:    map[ActionTypes][]users.ActionUser{},
	},
}

var existingChats = []Chat{
	{
		id:         1,
		avatar:     nil,
		title:      "",
		type_:      UserChatType,
		members:    []int{1, 2},
		isArchived: false,
		ownerId:    0,
		admins:     []int{},
		actions:    map[ActionTypes][]users.ActionUser{},
	},
	{
		id:         2,
		avatar:     &existingAvatars[0],
		title:      "group chat 1",
		type_:      GroupChatType,
		members:    []int{1, 2},
		isArchived: false,
		ownerId:    1,
		admins:     []int{1},
		actions:    map[ActionTypes][]users.ActionUser{},
	},
}

type TestChatsAdapter struct{}

func (adapter *TestChatsAdapter) GetById(id int) (*Chat, error) {
	var chat *Chat
	for _, dbChat := range existingChats {
		if dbChat.GetId() == id {
			chat = &dbChat
		}
	}

	return chat, nil
}

func (adapter *TestChatsAdapter) GetByIdForUser(id int, userID int) (*Chat, error) {
	var chat *Chat
	for _, dbChat := range existingChats {
		if dbChat.GetId() == id && slices.Contains(dbChat.GetMembers(), userID) {
			chat = &dbChat
		}
	}

	return chat, nil
}

func (adapter *TestChatsAdapter) GetByIdsForUser(ids []int, userID int) []Chat {
	var chats []Chat
	for _, dbChat := range existingChats {
		if slices.Contains(ids, dbChat.GetId()) && slices.Contains(dbChat.GetMembers(), userID) {
			chats = append(chats, dbChat)
		}
	}

	return chats
}

func (adapter *TestChatsAdapter) GetUserAll(userID int, page int, perPage int) utils.PaginatedResponse[Chat] {
	chats := make([]Chat, 0, len(existingChats))
	for _, dbChat := range existingChats {
		if slices.Contains(dbChat.GetMembers(), userID) {
			chats = append(chats, dbChat)
		}
	}

	return utils.NewPaginatedResponse(page, perPage, 1, len(chats), chats)
}

func (adapter *TestChatsAdapter) Save(chat Chat) (*Chat, error) {
	chatIds := make([]int, 0, len(existingChats))
	for _, dbChat := range existingChats {
		chatIds = append(chatIds, dbChat.GetId())
	}

	chat.id = slices.Max(chatIds) + 1
	existingChats = append(existingChats, chat)
	savedChat := &chat
	return savedChat, nil
}

func (adapter *TestChatsAdapter) HasDeletedUserChat(chat Chat) bool {
	for _, deletedChat := range deletedChats {
		if slices.Compare(deletedChat.GetMembers(), chat.GetMembers()) == 0 && deletedChat.GetType() == chat.GetType() {
			return true
		}
	}

	return false
}

func (adapter *TestChatsAdapter) RestoreChat(chat Chat) (*Chat, error) {
	restoredChat := &chat
	return restoredChat, nil
}

func (adapter *TestChatsAdapter) CheckChatExists(chat Chat) bool {
	chatIds := make([]int, 0, len(existingChats))
	for _, dbChat := range existingChats {
		chatIds = append(chatIds, dbChat.GetId())
	}

	return slices.Contains(chatIds, chat.GetId())
}

func (adapter *TestChatsAdapter) Delete(chat Chat) {
	var newExistingChats []Chat
	for _, dbChat := range existingChats {
		if dbChat.GetId() != chat.GetId() {
			newExistingChats = append(newExistingChats, dbChat)
		}
	}

	existingChats = newExistingChats
}

type TestChatEventsAdapter struct{}

func (adapter *TestChatEventsAdapter) SendChatCreated(chat Chat) {}

func (adapter *TestChatEventsAdapter) SendChatDeleted(chat Chat) {}

func (adapter *TestChatEventsAdapter) SendChatUserAction(chat Chat) {}
