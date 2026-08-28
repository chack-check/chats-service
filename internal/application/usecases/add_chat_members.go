package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type AddChatMembersUseCase struct {
	chatsRepository      ports.ChatsRepository
	usersRepository      ports.UsersRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewAddChatMembersUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *AddChatMembersUseCase {
	return &AddChatMembersUseCase{
		chatsRepository:      chatsRepository,
		usersRepository:      usersRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *AddChatMembersUseCase) Execute(chatID int, userID int, members []int) (*entities.Chat, error) {
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

	newMembers := chat.GetMembers()
	users := useCase.usersRepository.GetByIds(members)
	for _, member := range users {
		if !slices.Contains(newMembers, member.GetID()) {
			newMembers = append(newMembers, member.GetID())
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
