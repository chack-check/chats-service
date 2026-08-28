package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/application/services"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type UpdateGroupChatAvatarUseCase struct {
	chatsRepository      ports.ChatsRepository
	filesRepository      ports.FilesProvider
	chatEventsRepository ports.ChatEventsPublisher
}

func NewUpdateGroupChatAvatarUseCase(
	chatsRepository ports.ChatsRepository,
	filesRepository ports.FilesProvider,
	chatEventsRepository ports.ChatEventsPublisher,
) *UpdateGroupChatAvatarUseCase {
	return &UpdateGroupChatAvatarUseCase{
		chatsRepository:      chatsRepository,
		filesRepository:      filesRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (useCase *UpdateGroupChatAvatarUseCase) Execute(chatID int, userID int, newAvatar dtos.UploadingFile) (*entities.Chat, error) {
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

	err = services.ValidateUploadingFile(useCase.filesRepository, &newAvatar, constants.AvatarFiletype, true)
	if err != nil {
		return nil, err
	}

	savedFile := services.UploadingFileToSavedFile(newAvatar)
	chat.SetAvatar(savedFile)
	savedChat, err := useCase.chatsRepository.Save(*chat)
	if err != nil {
		return nil, domainerrors.ErrSavingChat
	}

	useCase.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}
