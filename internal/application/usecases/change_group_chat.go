package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type ChangeGroupChatUseCase struct {
	chatsRepository      ports.ChatsRepository
	chatEventsRepository ports.ChatEventsPublisher
}

func NewChangeGroupChatUseCase(
	chatsRepository ports.ChatsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
) *ChangeGroupChatUseCase {
	return &ChangeGroupChatUseCase{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *ChangeGroupChatUseCase) Execute(chatID int, userID int, chatData dtos.ChangeGroupChatData) (*entities.Chat, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	if userID != chat.GetOwnerID() && !ValidateUserChatAdmin(*chat, userID) {
		return nil, domainerrors.ErrChatNotAdmin
	}

	if chat.GetType() != constants.GroupChatType {
		return nil, domainerrors.ErrChatNotGroup
	}

	if chatData.GetTitle() != nil {
		chat.SetTitle(*chatData.GetTitle())
	} else {
		chat.SetTitle("")
	}

	savedChat, err := useCase.chatsRepository.Save(*chat)
	if err != nil {
		return nil, domainerrors.ErrSavingChat
	}

	useCase.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}
