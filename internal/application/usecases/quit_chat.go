package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type QuitChatUseCase struct {
	chatsRepository      ports.ChatsRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewQuitChatUseCase(
	chatsRepository ports.ChatsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *QuitChatUseCase {
	return &QuitChatUseCase{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *QuitChatUseCase) Execute(chatID int, userID int) (*entities.Chat, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	var newMembers []int
	for _, member := range chat.GetMembers() {
		if member != userID {
			newMembers = append(newMembers, member)
		}
	}

	var newAdmins []int
	for _, admin := range chat.GetAdmins() {
		if admin != userID {
			newAdmins = append(newAdmins, admin)
		}
	}

	chat.SetMembers(newMembers)
	chat.SetAdmins(newAdmins)
	savedChat, err := useCase.chatsRepository.Save(*chat)
	if err != nil {
		return nil, domainerrors.ErrSavingChat
	}

	useCase.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}
