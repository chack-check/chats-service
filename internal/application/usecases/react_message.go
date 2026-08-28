package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type ReactMessageUseCase struct {
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
}

func NewReactMessageUseCase(
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
) *ReactMessageUseCase {
	return &ReactMessageUseCase{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (useCase *ReactMessageUseCase) Execute(messageID int, userID int, content string) (*entities.Message, error) {
	message, err := useCase.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, domainerrors.ErrMessageNotFound
	}

	reaction := entities.NewMessageReaction(userID, content)

	if slices.Contains(message.GetReactions(), reaction) {
		return message, nil
	}

	message.AddReaction(reaction)
	savedMessage, err := useCase.messagesRepository.Save(*message)
	if err != nil {
		return nil, domainerrors.ErrSavingMessage
	}

	useCase.messageEventsRepository.SendMessageReacted(*savedMessage)
	return savedMessage, nil
}
