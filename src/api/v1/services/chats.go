package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/api/v1/utils"
	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
	"github.com/golang-jwt/jwt/v5"
)

type ChatsManager struct {
	ChatsQueries *models.ChatsQueries
}

func (manager *ChatsManager) decodeUserChatTitle(currentUserId int, title string) string {
    decodedTitles := map[string]int{}
    err := json.Unmarshal([]byte(title), &decodedTitles)

    if err != nil {
        return "Untitled"
    }

    for chatTitle, id := range decodedTitles {
        if id != currentUserId {
            return chatTitle
        }
    }

    return "Untitled"
}

func (manager *ChatsManager) decodeUserChatAvatar(currentUserId int, avatarUrl string) string {
    decodedAvatars := map[string]int{}
    err := json.Unmarshal([]byte(avatarUrl), &decodedAvatars)

    if err != nil {
        return "Untitled"
    }

    for chatAvatarUrl, id := range decodedAvatars {
        if id != currentUserId {
            return chatAvatarUrl
        }
    }

    return "Untitled"
}

func (manager *ChatsManager) decodeUsersChats(userId int, chats []models.Chat) []models.Chat {
    newChats := []models.Chat{}

    for _, chat := range chats {
        if chat.Type == "user" {
            chat.Title = manager.decodeUserChatTitle(userId, chat.Title)
            chat.AvatarURL = manager.decodeUserChatAvatar(userId, chat.AvatarURL)
            log.Printf("Chat title: %v", chat.Title)
        }
        newChats = append(newChats, chat)
    }

    return newChats
}

func (manager *ChatsManager) GetConcrete(chatID uint, token *jwt.Token) (*models.Chat, error) {
	tokenSubject, err := GetTokenSubject(token)

	if err != nil {
		return nil, err
	}

	chat, err := manager.ChatsQueries.GetWithMember(chatID, uint(tokenSubject.UserId))
    chat.Title = manager.decodeUserChatTitle(tokenSubject.UserId, chat.Title)
    chat.AvatarURL = manager.decodeUserChatAvatar(tokenSubject.UserId, chat.AvatarURL)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (manager *ChatsManager) GetAll(token *jwt.Token, page int, perPage int) *schemas.PaginatedResponse[models.Chat] {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil
	}

	count := manager.ChatsQueries.GetAllWithMemberCount(uint(tokenSubject.UserId))
	countValue := *count
	chats := manager.ChatsQueries.GetAllWithMember(uint(tokenSubject.UserId), page, perPage)
    decodedChats := manager.decodeUsersChats(tokenSubject.UserId, *chats)
    log.Printf("Decoded chats: %v", decodedChats)
	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(countValue), decodedChats)
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

func (manager *ChatsManager) validateUserChatExists(userId int, chat *models.Chat) error {
    alreadyCreatedChatsCount := manager.ChatsQueries.GetAllWithMemberCount(uint(userId))
    alreadyCreatedChats := manager.ChatsQueries.GetAllWithMember(uint(userId), 1, int(*alreadyCreatedChatsCount))

    for _, existingChat := range *alreadyCreatedChats {
        log.Printf("Comparing chats members: %v and %v", existingChat.Members, chat.Members)
        if existingChat.Type == "user" && utils.AreSlicesEqual(existingChat.Members, chat.Members) {
            log.Printf("Slices are the same")
            return fmt.Errorf("You already have chat with this user")
        }
        log.Printf("Slices are differrent")
    }

    return nil
}

func (manager *ChatsManager) createUserChat(chat *models.Chat, currentUser *protousers.UserResponse, chatUser *protousers.UserResponse) error {
    chatTitleForCurrentUser := fmt.Sprintf("%s %s", currentUser.LastName, currentUser.FirstName)
    chatTitleForChatUser := fmt.Sprintf("%s %s", chatUser.LastName, chatUser.FirstName)
    chatTitles := map[string]int{
        chatTitleForCurrentUser: int(currentUser.Id),
        chatTitleForChatUser: int(chatUser.Id),
    }
    chatAvatarUrls := map[string]int{
        currentUser.AvatarUrl: int(currentUser.Id),
        chatUser.AvatarUrl: int(chatUser.Id),
    }
    jsonTitles, err := json.Marshal(chatTitles)
    if err != nil {
        return err
    }

    jsonAvatars, err := json.Marshal(chatAvatarUrls)
    if err != nil {
        return err
    }

	chat.Members = []int64{int64(currentUser.Id), int64(chatUser.Id)}
	chat.OwnerId = 0
	chat.Title = string(jsonTitles)
	chat.Type = "user"
	chat.AvatarURL = string(jsonAvatars)

    err = manager.validateUserChatExists(int(currentUser.Id), chat)
    if err != nil {
        return err
    }

    log.Printf("Creating chat: %v", chat)

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

    chat.Title = chatTitleForChatUser
    chat.AvatarURL = chatUser.AvatarUrl

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

	currentUser, err := grpc_client.UsersGrpcClient.GetUserById(tokenSubject.UserId)

	if err != nil || currentUser == nil {
		return fmt.Errorf("There is no user with id %d", tokenSubject.UserId)
	}

	chatUser, err := grpc_client.UsersGrpcClient.GetUserById(int(chatUserId))

	if err != nil || chatUser == nil {
		return fmt.Errorf("There is no user with id %d", chatUserId)
	}

    log.Printf("Creating user chat: %v, user id: %v, chat user: %v", chat, tokenSubject.UserId, chatUser)
	return manager.createUserChat(chat, currentUser, chatUser)
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		ChatsQueries: &models.ChatsQueries{},
	}
}
