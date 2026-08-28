package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type GetChatsByIdsUseCase struct {
	chatsRepository       ports.ChatsRepository
	usersRepository       ports.UsersRepository
	userActionsRepository ports.UserActionsRepository
}

func NewGetChatsByIdsUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
	userActionsRepository ports.UserActionsRepository,
) *GetChatsByIdsUseCase {
	return &GetChatsByIdsUseCase{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
	}
}

func (useCase *GetChatsByIdsUseCase) Execute(chatIds []int, userID int) []entities.Chat {
	requestChats := useCase.chatsRepository.GetByIdsForUser(chatIds, userID)
	fetchingUsers := GetUserChatsUsersIds(requestChats, userID)
	fetchedUsers := useCase.usersRepository.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(requestChats, fetchedUsers, userID)
	completeChats := make([]entities.Chat, 0, len(chatsWithUsersData))
	for _, chat := range chatsWithUsersData {
		setupSavedMessagesChatAvatar(&chat)
		chatActions := useCase.userActionsRepository.GetAllChatActionsUsers(chat)
		chat.SetupActions(chatActions)
		completeChats = append(completeChats, chat)
	}

	return completeChats
}
