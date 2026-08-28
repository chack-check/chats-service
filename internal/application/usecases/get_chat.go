package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type GetChatUseCase struct {
	chatsRepository       ports.ChatsRepository
	usersRepository       ports.UsersRepository
	userActionsRepository ports.UserActionsRepository
}

func NewGetChatUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
	userActionsRepository ports.UserActionsRepository,
) *GetChatUseCase {
	return &GetChatUseCase{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
	}
}

func (useCase *GetChatUseCase) Execute(userID int, chatID int) (*entities.Chat, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	if chat.GetType() != "user" {
		setupSavedMessagesChatAvatar(chat)
		return chat, nil
	}

	anotherUserID := getAnotherUserIDForUserChat(*chat, userID)
	if anotherUserID == 0 {
		return chat, nil
	}

	anotherUser, err := useCase.usersRepository.GetById(anotherUserID)
	if err != nil {
		return chat, nil
	}

	chatActions := useCase.userActionsRepository.GetAllChatActionsUsers(*chat)
	chat.SetupActions(chatActions)
	chat.SetupUserData(anotherUser)
	return chat, nil
}
