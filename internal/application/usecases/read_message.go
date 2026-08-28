package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type ReadMessageUseCase struct {
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
}

func NewReadMessageUseCase(
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
) *ReadMessageUseCase {
	return &ReadMessageUseCase{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (useCase *ReadMessageUseCase) Execute(messageID int, userID int) (*entities.Message, error) {
	message, err := useCase.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, domainerrors.ErrMessageNotFound
	}

	if slices.Contains(message.GetReadedBy(), userID) {
		return message, nil
	}

	message.Read(userID)
	savedMessage, err := useCase.messagesRepository.Save(*message)
	if err != nil {
		return nil, domainerrors.ErrSavingMessage
	}

	useCase.messageEventsRepository.SendMessageReaded(*savedMessage)
	return savedMessage, nil
}
