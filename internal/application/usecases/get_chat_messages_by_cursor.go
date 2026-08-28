package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type GetChatMessagesByCursorUseCase struct {
	messagesRepository ports.MessagesRepository
	chatsRepository    ports.ChatsRepository
}

func NewGetChatMessagesByCursorUseCase(
	messagesRepository ports.MessagesRepository,
	chatsRepository ports.ChatsRepository,
) *GetChatMessagesByCursorUseCase {
	return &GetChatMessagesByCursorUseCase{
		messagesRepository: messagesRepository,
		chatsRepository:    chatsRepository,
	}
}

func (useCase *GetChatMessagesByCursorUseCase) Execute(chatID int, userID int, messageID int, aroundOffset int) (*dtos.OffsetResponse[entities.Message], error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	messages := useCase.messagesRepository.GetChatCursorAllForUser(chat.GetID(), userID, messageID, aroundOffset)
	return &messages, nil
}
