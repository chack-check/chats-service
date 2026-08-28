package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type GetConcreteMessageUseCase struct {
	messagesRepository ports.MessagesRepository
}

func NewGetConcreteMessageUseCase(
	messagesRepository ports.MessagesRepository,
) *GetConcreteMessageUseCase {
	return &GetConcreteMessageUseCase{messagesRepository: messagesRepository}
}

func (useCase *GetConcreteMessageUseCase) Execute(messageID int, userID int) (*entities.Message, error) {
	message, err := useCase.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, domainerrors.ErrMessageNotFound
	}

	return message, nil
}
