package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type UserActionUseCase struct {
	chatsRepository       ports.ChatsRepository
	usersRepository       ports.UsersRepository
	userActionsRepository ports.UserActionsRepository
	chatEventsRepository  ports.ChatEventsPublisher
}

func NewUserActionUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
	userActionsRepository ports.UserActionsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *UserActionUseCase {
	return &UserActionUseCase{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
		chatEventsRepository:  chatEventsRepository,
	}
}

func (useCase *UserActionUseCase) Execute(chatID int, userID int, actionType constants.ActionTypes) (*entities.Chat, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	user, err := useCase.usersRepository.GetById(userID)
	if err != nil {
		return nil, domainerrors.ErrFindingUser
	}

	newChatActions := useCase.userActionsRepository.AddChatActionUser(*chat, *user, actionType)
	chat.SetupActions(newChatActions)
	useCase.chatEventsRepository.SendChatUserAction(*chat)
	return chat, nil
}
