package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type GetChatMessagesUseCase struct {
	messagesRepository ports.MessagesRepository
	chatsRepository    ports.ChatsRepository
}

func NewGetChatMessagesUseCase(
	messagesRepository ports.MessagesRepository,
	chatsRepository ports.ChatsRepository,
) *GetChatMessagesUseCase {
	return &GetChatMessagesUseCase{
		messagesRepository: messagesRepository,
		chatsRepository:    chatsRepository,
	}
}

func (useCase *GetChatMessagesUseCase) Execute(chatID int, userID int, offset int, limit int) (*dtos.OffsetResponse[entities.Message], error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	messages := useCase.messagesRepository.GetChatAllForUser(chat.GetID(), userID, offset, limit)
	return &messages, nil
}
