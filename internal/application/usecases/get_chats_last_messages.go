package usecases

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type GetChatsLastMessagesUseCase struct {
	messagesRepository ports.MessagesRepository
	chatsRepository    ports.ChatsRepository
}

func NewGetChatsLastMessagesUseCase(
	messagesRepository ports.MessagesRepository,
	chatsRepository ports.ChatsRepository,
) *GetChatsLastMessagesUseCase {
	return &GetChatsLastMessagesUseCase{
		messagesRepository: messagesRepository,
		chatsRepository:    chatsRepository,
	}
}

func (useCase *GetChatsLastMessagesUseCase) Execute(chatIds []int, userID int) []entities.Message {
	chats := useCase.chatsRepository.GetByIdsForUser(chatIds, userID)
	fetchedChatIds := make([]int, 0, len(chats))
	for _, chat := range chats {
		fetchedChatIds = append(fetchedChatIds, chat.GetID())
	}

	messages := useCase.messagesRepository.GetChatsLast(fetchedChatIds, userID)
	return messages
}
