package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"errors"

	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/application/services"
	"chats-service/internal/domain/entities"
)

type CreateSavedMessagesChatUseCase struct {
	chatsRepository ports.ChatsRepository
}

func NewCreateSavedMessagesChatUseCase(chatsRepository ports.ChatsRepository) *CreateSavedMessagesChatUseCase {
	return &CreateSavedMessagesChatUseCase{
		chatsRepository: chatsRepository,
	}
}

func (useCase *CreateSavedMessagesChatUseCase) Execute(data dtos.CreateChatData, currentUserID int) (*entities.Chat, error) {
	chat := services.CreateChatDataToChat(data, currentUserID)
	chat.SetOwnerID(currentUserID)
	chat.SetMembers([]int{currentUserID})
	chat.SetTitle("Saved messages")
	savedChat, err := useCase.chatsRepository.Save(chat)
	if err != nil {
		return nil, errors.Join(domainerrors.ErrSavingChat, err)
	}

	setupSavedMessagesChatAvatar(savedChat)
	return savedChat, nil
}
