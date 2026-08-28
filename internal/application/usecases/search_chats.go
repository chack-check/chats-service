package usecases

import (
	"strings"

	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

type SearchChatsUseCase struct {
	chatsRepository ports.ChatsRepository
	usersRepository ports.UsersRepository
}

func NewSearchChatsUseCase(
	chatsRepository ports.ChatsRepository,
	usersRepository ports.UsersRepository,
) *SearchChatsUseCase {
	return &SearchChatsUseCase{
		chatsRepository: chatsRepository,
		usersRepository: usersRepository,
	}
}

func (useCase *SearchChatsUseCase) Execute(userID int, query string, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	requestChats := useCase.chatsRepository.SearchChats(userID, query, page, perPage)

	fetchingUsers := GetUserChatsUsersIds(requestChats.GetData(), userID)
	fetchedUsers := useCase.usersRepository.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(requestChats.GetData(), fetchedUsers, userID)

	var resultChats []entities.Chat

	for _, chat := range chatsWithUsersData {
		if chat.GetType() == constants.SavedMessagesChatType {
			setupSavedMessagesChatAvatar(&chat)
		}

		if chat.GetType() == constants.UserChatType && !strings.Contains(strings.ToLower(chat.GetTitle()), strings.ToLower(query)) {
			continue
		}

		resultChats = append(resultChats, chat)
	}

	requestChats.SetData(resultChats)
	return requestChats
}
