package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type GetMessagesByIdsUseCase struct {
	messagesRepository ports.MessagesRepository
}

func NewGetMessagesByIdsUseCase(
	messagesRepository ports.MessagesRepository,
) *GetMessagesByIdsUseCase {
	return &GetMessagesByIdsUseCase{messagesRepository: messagesRepository}
}

func (useCase *GetMessagesByIdsUseCase) Execute(messageIds []int, userID int) []entities.Message {
	messages := useCase.messagesRepository.GetByIdsForUser(messageIds, userID)
	return messages
}
