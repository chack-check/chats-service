package usecases

import (
	domainerrors "chats-service/internal/domain/errors"
	"errors"

	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/application/services"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type CreateChatUseCase struct {
	chatsRepository      ports.ChatsRepository
	chatEventsRepository ports.ChatEventsPublisher
	usersRepository      ports.UsersRepository
	filesRepository      ports.FilesProvider
}

func NewCreateChatUseCase(
	chatsRepository ports.ChatsRepository,
	chatEventsRepository ports.ChatEventsPublisher,
	usersRepository ports.UsersRepository,
	filesRepository ports.FilesProvider,
) *CreateChatUseCase {
	return &CreateChatUseCase{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
		usersRepository:      usersRepository,
		filesRepository:      filesRepository,
	}
}

func (useCase *CreateChatUseCase) Execute(data dtos.CreateChatData, currentUserID int) (*entities.Chat, error) {
	if err := services.ValidateUploadingFile(useCase.filesRepository, data.GetAvatar(), constants.AvatarFiletype, false); err != nil {
		return nil, err
	}

	currentUser, err := useCase.usersRepository.GetById(currentUserID)
	if err != nil {
		return nil, domainerrors.ErrFindingUser
	}

	var savedChat *entities.Chat
	var savingError error
	switch data.GetType() {
	case constants.GroupChatType:
		savedChat, savingError = useCase.createGroupChat(data, currentUser)
	case constants.UserChatType:
		savedChat, savingError = useCase.createUserChat(data, currentUser)
	default:
		savingError = domainerrors.ErrInvalidCreatingChatType
	}

	if savingError != nil {
		return nil, savingError
	}

	useCase.chatEventsRepository.SendChatCreated(*savedChat)
	return savedChat, nil
}

func (useCase *CreateChatUseCase) createGroupChat(data dtos.CreateChatData, currentUser *entities.User) (*entities.Chat, error) {
	chat := services.CreateChatDataToChat(data, 0)
	chat.SetOwnerID(currentUser.GetID())
	if !ValidateUserChatMember(chat, currentUser.GetID()) {
		newMembers := chat.GetMembers()
		newMembers = append(newMembers, currentUser.GetID())
		chat.SetMembers(newMembers)
	}
	if !ValidateUserChatAdmin(chat, currentUser.GetID()) {
		newAdmins := chat.GetAdmins()
		newAdmins = append(newAdmins, currentUser.GetID())
		chat.SetAdmins(newAdmins)
	}

	chat.SetType("group")
	savedChat, err := useCase.chatsRepository.Save(chat)
	if err != nil {
		return nil, errors.Join(domainerrors.ErrSavingChat, err)
	}

	return savedChat, nil
}

func (useCase *CreateChatUseCase) createUserChat(data dtos.CreateChatData, currentUser *entities.User) (*entities.Chat, error) {
	if data.GetUserID() == nil {
		return nil, domainerrors.ErrCreatingNotUserChat
	}
	if *data.GetUserID() == currentUser.GetID() {
		return nil, domainerrors.ErrChatWithSelf
	}

	chatUser, err := useCase.usersRepository.GetById(*data.GetUserID())
	if err != nil {
		return nil, domainerrors.ErrFindingUser
	}

	chat := services.CreateChatDataToChat(data, currentUser.GetID())
	if useCase.chatsRepository.HasDeletedUserChat(chat) {
		chat, err := useCase.chatsRepository.RestoreChat(chat)
		if err != nil {
			return nil, errors.Join(domainerrors.ErrRestoringChat, err)
		}

		return chat, nil
	}

	if useCase.chatsRepository.CheckChatExists(chat) {
		return nil, domainerrors.ErrChatAlreadyExists
	}

	savedChat, err := useCase.chatsRepository.Save(chat)
	if err != nil {
		return nil, errors.Join(domainerrors.ErrSavingChat, err)
	}

	savedChat.SetupUserData(chatUser)
	return savedChat, nil
}
