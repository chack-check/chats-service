package services

import (
	"fmt"

	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
	"github.com/golang-jwt/jwt/v5"
)

type ChatsManager struct {
	ChatsQueries *models.ChatsQueries
}

func (manager *ChatsManager) GetConcrete(chatID uint, token *jwt.Token) (*models.Chat, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	chat, err := manager.ChatsQueries.GetWithMember(chatID, uint(tokenSubject.UserId))

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (manager *ChatsManager) GetAll(token *jwt.Token, page *int, perPage *int) *schemas.PaginatedResponse[models.Chat] {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil
	}

	count := manager.ChatsQueries.GetAllWithMemberCount(uint(tokenSubject.UserId))
	countValue := *count
	chats := manager.ChatsQueries.GetAllWithMember(uint(tokenSubject.UserId), page, perPage)
	paginatedResponse := schemas.NewPaginatedResponse[models.Chat](*page, *perPage, int(countValue), *chats)
	return &paginatedResponse
}

func (manager *ChatsManager) createGroupChat(chat *models.Chat, userId int) error {
	chat.OwnerId = uint(userId)
	chat.Type = "group"

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	return nil
}

func (manager *ChatsManager) createUserChat(chat *models.Chat, userId int, chatUser *protousers.UserResponse) error {
	chat.Members = []int64{int64(userId), int64(chatUser.Id)}
	chat.OwnerId = 0
	chat.Title = fmt.Sprintf("%v %v %v", chatUser.LastName, chatUser.FirstName, chatUser.MiddleName)
	chat.Type = "user"
	chat.AvatarURL = "https://google.com"

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	return nil
}

func (manager *ChatsManager) Create(chat *models.Chat, token *jwt.Token, chatUserId uint) error {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return err
	}

	if chatUserId == 0 {
		return manager.createGroupChat(chat, tokenSubject.UserId)
	}

	if chatUserId == uint(tokenSubject.UserId) {
		return fmt.Errorf("You can't create chat for you")
	}

	chatUser, err := grpc_client.UsersGrpcClient.GetUserById(int(chatUserId))

	if err != nil {
		return fmt.Errorf("There is no user with id %d", chatUserId)
	}

	return manager.createUserChat(chat, tokenSubject.UserId, chatUser)
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		ChatsQueries: &models.ChatsQueries{},
	}
}
