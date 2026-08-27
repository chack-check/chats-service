// TODO: Delete Port from names and think about Repository for all (maybe Provider for files signature, publisher for events, store for user actions)
package ports

import (
	"chats-service/internal/domain/dtos"
	"chats-service/internal/domain/entities"
)

type MessagesRepositoryPort interface {
	GetChatAllForUser(chatID int, userID int, offset int, limit int) dtos.OffsetResponse[entities.Message]
	GetChatCursorAllForUser(chatID int, userID int, messageID int, aroundOffset int) dtos.OffsetResponse[entities.Message]
	GetChatsLast(chatIds []int, userID int) []entities.Message
	GetByIdForUser(messageID int, userID int) (*entities.Message, error)
	GetByIdsForUser(messageIds []int, userID int) []entities.Message
	GetById(messageID int) (*entities.Message, error)
	Save(message entities.Message) (*entities.Message, error)
	Delete(message entities.Message)
}

type MessageEventsRepositoryPort interface {
	SendMessageReacted(message entities.Message)
	SendReactionDeleted(message entities.Message)
	SendMessageReaded(message entities.Message)
	SendMessageDeleted(message entities.Message)
	SendMessageUpdated(message entities.Message)
	SendMessageCreated(message entities.Message)
}
