package usecases

import (
	"chats-service/internal/application/ports"
	domainerrors "chats-service/internal/domain/errors"
)

type RecognizeMessageUseCase struct {
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
}

func NewRecognizeMessageUseCase(
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
) *RecognizeMessageUseCase {
	return &RecognizeMessageUseCase{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (useCase *RecognizeMessageUseCase) Execute(messageID int, content string) error {
	message, err := useCase.messagesRepository.GetById(messageID)
	if err != nil {
		return domainerrors.ErrMessageNotFound
	}

	message.SetContent(&content)
	_, err = useCase.messagesRepository.Save(*message)
	if err != nil {
		return domainerrors.ErrSavingMessage
	}

	useCase.messageEventsRepository.SendMessageUpdated(*message)
	return nil
}
