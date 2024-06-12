package chats

import (
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/users"
	"github.com/chack-check/chats-service/domain/utils"
)

type ChatsPort interface {
	GetById(id int) (*Chat, error)
	GetByIdForUser(id int, userId int) (*Chat, error)
	GetByIdsForUser(ids []int, userId int) []Chat
	GetUserAll(userId int, page int, perPage int) utils.PaginatedResponse[Chat]
	Save(chat Chat) (*Chat, error)
	HasDeletedUserChat(chat Chat) bool
	RestoreChat(chat Chat) (*Chat, error)
	CheckChatExists(chat Chat) bool
	Delete(chat Chat)
	SearchChats(userId int, query string, page int, perPage int) utils.PaginatedResponse[Chat]
}

type ChatEventsPort interface {
	SendChatCreated(chat Chat)
	SendChatDeleted(chat Chat)
	SendChatUserAction(chat Chat)
	SendChatChanged(chat Chat)
}

type UserActionsPort interface {
	AddChatActionUser(chat Chat, user users.User, actionType ActionTypes) map[ActionTypes][]users.ActionUser
	RemoveChatActionUser(chat Chat, userId int, actionType ActionTypes) map[ActionTypes][]users.ActionUser
	GetAllChatActionsUsers(chat Chat) map[ActionTypes][]users.ActionUser
}

func NewCreateChatHandler(
	chatsPort ChatsPort,
	chatEventsPort ChatEventsPort,
	usersPort users.UsersPort,
	filesPort files.FilesPort,
) CreateChatHandler {
	return CreateChatHandler{
		chatsPort:      chatsPort,
		chatEventsPort: chatEventsPort,
		usersPort:      usersPort,
		filesPort:      filesPort,
	}
}

func NewCreateSavedMessagesChatHandler(
	chatsPort ChatsPort,
) CreateSavedMessagesChat {
	return CreateSavedMessagesChat{chatsPort: chatsPort}
}

func NewDeleteChatHandler(
	chatsPort ChatsPort,
	chatEventsPort ChatEventsPort,
) DeleteChatHandler {
	return DeleteChatHandler{
		chatsPort:      chatsPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewUserActionHandler(
	chatsPort ChatsPort,
	chatEventsPort ChatEventsPort,
	usersPort users.UsersPort,
	userActionsPort UserActionsPort,
) UserActionHandler {
	return UserActionHandler{
		chatsPort:       chatsPort,
		chatEventsPort:  chatEventsPort,
		usersPort:       usersPort,
		userActionsPort: userActionsPort,
	}
}

func NewStopUserActionHandler(
	chatsPort ChatsPort,
	chatEventsPort ChatEventsPort,
	userActionsPort UserActionsPort,
) StopUserActionHandler {
	return StopUserActionHandler{
		chatsPort:       chatsPort,
		chatEventsPort:  chatEventsPort,
		userActionsPort: userActionsPort,
	}
}

func NewGetChatsHandler(
	chatsPort ChatsPort,
	usersPort users.UsersPort,
	userActionsPort UserActionsPort,
) GetChatsHandler {
	return GetChatsHandler{
		chatsPort:       chatsPort,
		usersPort:       usersPort,
		userActionsPort: userActionsPort,
	}
}

func NewGetChatHandler(
	chatsPort ChatsPort,
	usersPort users.UsersPort,
	userActionsPort UserActionsPort,
) GetChatHandler {
	return GetChatHandler{
		chatsPort:       chatsPort,
		usersPort:       usersPort,
		userActionsPort: userActionsPort,
	}
}

func NewGetChatsByIdsHandler(
	chatsPort ChatsPort,
	usersPort users.UsersPort,
	userActionsPort UserActionsPort,
) GetChatsByIdsHandler {
	return GetChatsByIdsHandler{
		chatsPort:       chatsPort,
		usersPort:       usersPort,
		userActionsPort: userActionsPort,
	}
}

func NewAddChatsMembersHandler(
	chatsPort ChatsPort,
	usersPort users.UsersPort,
	chatEventsPort ChatEventsPort,
) AddChatMembersHandler {
	return AddChatMembersHandler{
		chatsPort:      chatsPort,
		usersPort:      usersPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewAddChatsAdminsHandler(
	chatsPort ChatsPort,
	usersPort users.UsersPort,
	chatEventsPort ChatEventsPort,
) AddChatAdminsHandler {
	return AddChatAdminsHandler{
		chatsPort:      chatsPort,
		usersPort:      usersPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewRemoveChatMembersHandler(
	chatsPort ChatsPort,
	chatEventsPort ChatEventsPort,
) RemoveChatMembersHandler {
	return RemoveChatMembersHandler{
		chatsPort:      chatsPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewRemoveChatAdminsHandler(
	chatsPort ChatsPort,
	chatEventsPort ChatEventsPort,
) RemoveChatAdminsHandler {
	return RemoveChatAdminsHandler{
		chatsPort:      chatsPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewQuitChatHandler(chatsPort ChatsPort, chatEventsPort ChatEventsPort) QuitChatHandler {
	return QuitChatHandler{
		chatsPort:      chatsPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewChangeGroupChatHandler(chatsPort ChatsPort, chatEventsPort ChatEventsPort) ChangeGroupChatHandler {
	return ChangeGroupChatHandler{
		chatsPort:      chatsPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewUpdateGroupChatAvatar(chatsPort ChatsPort, filesPort files.FilesPort, chatEventsPort ChatEventsPort) UpdateGroupChatAvatar {
	return UpdateGroupChatAvatar{
		chatsPort:      chatsPort,
		filesPort:      filesPort,
		chatEventsPort: chatEventsPort,
	}
}

func NewSearchChatsHandler(chatsPort ChatsPort, usersPort users.UsersPort, userActionsPort UserActionsPort) SearchChatsHandler {
	return SearchChatsHandler{
		chatsPort:       chatsPort,
		usersPort:       usersPort,
		userActionsPort: userActionsPort,
	}
}
