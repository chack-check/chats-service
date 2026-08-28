package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type GetChatsUseCase struct {
	chatsRepository       ports.ChatsRepository
	usersRepository       ports.UsersRepository
	userActionsRepository ports.UserActionsRepository
}

func NewGetChatsUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
	userActionsRepository ports.UserActionsRepository,
) *GetChatsUseCase {
	return &GetChatsUseCase{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
	}
}

func (useCase *GetChatsUseCase) Execute(userID int, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	paginatedChats := useCase.chatsRepository.GetUserAll(userID, page, perPage)
	fetchingUsers := GetUserChatsUsersIds(paginatedChats.GetData(), userID)
	fetchedUsers := useCase.usersRepository.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(paginatedChats.GetData(), fetchedUsers, userID)
	completeChats := make([]entities.Chat, 0, len(chatsWithUsersData))
	for _, chat := range chatsWithUsersData {
		setupSavedMessagesChatAvatar(&chat)
		chatActions := useCase.userActionsRepository.GetAllChatActionsUsers(chat)
		chat.SetupActions(chatActions)
		completeChats = append(completeChats, chat)
	}

	paginatedChats.SetData(completeChats)
	return paginatedChats
}
