package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type DeleteMessageReactionUseCase struct {
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
}

func NewDeleteMessageReactionUseCase(
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
) *DeleteMessageReactionUseCase {
	return &DeleteMessageReactionUseCase{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (useCase *DeleteMessageReactionUseCase) Execute(messageID int, userID int) (*entities.Message, error) {
	message, err := useCase.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, domainerrors.ErrMessageNotFound
	}

	var userReaction *entities.MessageReaction
	for _, reaction := range message.GetReactions() {
		if reaction.GetUserID() == userID {
			userReaction = &reaction
		}
	}

	if userReaction == nil {
		return message, nil
	}

	message.RemoveReaction(*userReaction)
	savedMessage, err := useCase.messagesRepository.Save(*message)
	if err != nil {
		return nil, domainerrors.ErrSavingMessage
	}

	useCase.messageEventsRepository.SendReactionDeleted(*savedMessage)
	return savedMessage, nil
}
