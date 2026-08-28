package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type AddChatAdminsUseCase struct {
	chatsRepository      ports.ChatsRepository
	usersRepository      ports.UsersRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewAddChatAdminsUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *AddChatAdminsUseCase {
	return &AddChatAdminsUseCase{
		chatsRepository:      chatsRepository,
		usersRepository:      usersRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *AddChatAdminsUseCase) Execute(chatID int, userID int, admins []int) (*entities.Chat, error) {
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

	newAdmins := chat.GetAdmins()
	users := useCase.usersRepository.GetByIds(admins)
	for _, admin := range users {
		if !slices.Contains(newAdmins, admin.GetID()) {
			newAdmins = append(newAdmins, admin.GetID())
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
