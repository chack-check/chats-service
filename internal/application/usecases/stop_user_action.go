package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type StopUserActionUseCase struct {
	chatsRepository       ports.ChatsRepository
	userActionsRepository ports.UserActionsRepository
	chatEventsRepository  ports.ChatEventsPublisher
}

func NewStopUserActionUseCase(
	chatsRepository ports.ChatsRepository,
	userActionsRepository ports.UserActionsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *StopUserActionUseCase {
	return &StopUserActionUseCase{
		chatsRepository:       chatsRepository,
		userActionsRepository: userActionsRepository,
		chatEventsRepository:  chatEventsRepository,
	}
}

func (useCase *StopUserActionUseCase) Execute(chatID int, userID int, actionType constants.ActionTypes) (*entities.Chat, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	newChatActions := useCase.userActionsRepository.RemoveChatActionUser(*chat, userID, actionType)
	chat.SetupActions(newChatActions)
	useCase.chatEventsRepository.SendChatUserAction(*chat)
	return chat, nil
}
