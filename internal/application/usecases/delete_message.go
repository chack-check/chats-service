package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
)

type DeleteMessageUseCase struct {
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
}

func NewDeleteMessageUseCase(
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
) *DeleteMessageUseCase {
	return &DeleteMessageUseCase{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (useCase *DeleteMessageUseCase) Execute(messageID int, userID int) error {
	message, err := useCase.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return domainerrors.ErrMessageNotFound
	}

	chat := message.GetChat()
	if !slices.Contains(chat.GetMembers(), userID) {
		return domainerrors.ErrCantDeleteMessage
	}

	useCase.messagesRepository.Delete(*message)
	useCase.messageEventsRepository.SendMessageDeleted(*message)
	return nil
}
