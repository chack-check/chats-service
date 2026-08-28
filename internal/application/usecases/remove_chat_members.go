package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type RemoveChatMembersUseCase struct {
	chatsRepository      ports.ChatsRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewRemoveChatMembersUseCase(
	chatsRepository ports.ChatsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *RemoveChatMembersUseCase {
	return &RemoveChatMembersUseCase{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *RemoveChatMembersUseCase) Execute(chatID int, userID int, members []int) (*entities.Chat, error) {
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

	var newMembers []int
	for _, member := range chat.GetMembers() {
		if !slices.Contains(members, member) || member == userID {
			newMembers = append(newMembers, member)
		}
	}

	chat.SetMembers(newMembers)
	savedChat, err := useCase.chatsRepository.Save(*chat)
	if err != nil {
		return nil, domainerrors.ErrSavingChat
	}

	useCase.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}
