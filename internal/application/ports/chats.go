package ports

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type ChatsRepository interface {
	GetById(id int) (*entities.Chat, error)
	GetByIdForUser(id int, userID int) (*entities.Chat, error)
	GetByIdsForUser(ids []int, userID int) []entities.Chat
	GetUserAll(userID int, page int, perPage int) dtos.PaginatedResponse[entities.Chat]
	Save(chat entities.Chat) (*entities.Chat, error)
	HasDeletedUserChat(chat entities.Chat) bool
	RestoreChat(chat entities.Chat) (*entities.Chat, error)
	CheckChatExists(chat entities.Chat) bool
	Delete(chat entities.Chat)
	SearchChats(userID int, query string, page int, perPage int) dtos.PaginatedResponse[entities.Chat]
}

type ChatEventsPublisher interface {
	SendChatCreated(chat entities.Chat)
	SendChatDeleted(chat entities.Chat)
	SendChatUserAction(chat entities.Chat)
	SendChatChanged(chat entities.Chat)
}

type UserActionsRepository interface {
	AddChatActionUser(chat entities.Chat, user entities.User, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser
	RemoveChatActionUser(chat entities.Chat, userID int, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser
	GetAllChatActionsUsers(chat entities.Chat) map[constants.ActionTypes][]entities.ActionUser
}
