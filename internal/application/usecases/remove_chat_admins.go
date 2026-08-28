package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type RemoveChatAdminsUseCase struct {
	chatsRepository      ports.ChatsRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewRemoveChatAdminsUseCase(
	chatsRepository ports.ChatsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *RemoveChatAdminsUseCase {
	return &RemoveChatAdminsUseCase{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *RemoveChatAdminsUseCase) Execute(chatID int, userID int, admins []int) (*entities.Chat, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userID) {
		return nil, domainerrors.ErrChatNotAdmin
	}
	if chat.GetType() != constants.GroupChatType {
		return nil, domainerrors.ErrChatNotGroup
	}

	var newAdmins []int
	for _, admin := range chat.GetAdmins() {
		if !slices.Contains(admins, admin) || admin == userID {
			newAdmins = append(newAdmins, admin)
		}
	}

	chat.SetAdmins(newAdmins)
	savedChat, err := useCase.chatsRepository.Save(*chat)
	if err != nil {
		return nil, domainerrors.ErrSavingChat
	}

	useCase.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}
