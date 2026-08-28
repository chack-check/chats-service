package usecases

import (
	"chats-service/internal/application/ports"
	domainerrors "chats-service/internal/domain/errors"
)

type DeleteChatUseCase struct {
	chatsRepository      ports.ChatsRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewDeleteChatUseCase(
	chatsRepository ports.ChatsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *DeleteChatUseCase {
	return &DeleteChatUseCase{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *DeleteChatUseCase) Execute(chatID, userID int) error {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return domainerrors.ErrChatNotFound
	}

	useCase.chatsRepository.Delete(*chat)
	useCase.chatEventsRepository.SendChatDeleted(*chat)
	return nil
}
