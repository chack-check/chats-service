package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/repositories"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/api/v1/utils"
	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/chack-check/chats-service/redisdb"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func setupChatTitleAndAvatar(chat *dtos.ChatDto, chatUser *protousers.UserResponse) {
	log.Printf("Setting up chat title and avatar for chat = %+v and chat user = %+v", chat, chatUser)
	var chat_avatar_converted_url *string
	if chatUser.ConvertedAvatarUrl != "" {
		chat_avatar_converted_url = &chatUser.ConvertedAvatarUrl
	} else {
		chat_avatar_converted_url = nil
	}

	chat.Title = fmt.Sprintf("%s %s", chatUser.LastName, chatUser.FirstName)
	chat.Avatar.OriginalUrl = chatUser.OriginalAvatarUrl
	chat.Avatar.ConvertedUrl = chat_avatar_converted_url
	log.Printf("Chat with new title and avatar: %+v", chat)
}

func setupUserChatData(chat *dtos.ChatDto, currentUserId int) {
	log.Printf("Setting up user chat data for chat = %+v and current user id = %d", chat, currentUserId)
	var anotherUserId int
	for _, member := range chat.Members {
		if member != int(currentUserId) {
			anotherUserId = member
		}
	}

	log.Printf("Another user id for user chat = %d", anotherUserId)
	anotherUser, err := grpc_client.UsersGrpcClient.GetUserById(anotherUserId)
	log.Printf("Another user for user chat = %+v", anotherUser)
	if err != nil {
		sentry.CaptureException(err)
		chat.Title = "Untitled"
	} else {
		setupChatTitleAndAvatar(chat, anotherUser)
	}
}

func setupUserManyChatsData(chats []*dtos.ChatDto, currentUserId int) {
	log.Printf("Setting up many chats data for current user id = %d", currentUserId)
	var anotherUsersIds []int
	for _, chat := range chats {
		if chat.Type != "user" {
			log.Printf("Skipping chat setting up data because chat type != user: %s", chat.Type)
			continue
		}

		for _, member := range chat.Members {
			if member != currentUserId && !slices.Contains(anotherUsersIds, int(member)) {
				anotherUsersIds = append(anotherUsersIds, int(member))
			}
		}
	}

	log.Printf("Fetching another users for chats by ids: %v", anotherUsersIds)
	anotherUsers, err := grpc_client.UsersGrpcClient.GetUsersByIds(anotherUsersIds)
	if err != nil {
		log.Printf("Error fetching users by ids: %v", err)
		sentry.CaptureException(err)
		return
	}

	for _, chat := range chats {
		if chat.Type != "user" {
			continue
		}

		var anotherUser *protousers.UserResponse

		for _, member := range chat.Members {
			if member == currentUserId {
				continue
			}

			for _, user := range anotherUsers.Users {
				if user.Id == int32(member) {
					anotherUser = user
				}
			}
		}

		setupChatTitleAndAvatar(chat, anotherUser)
	}
}

type ChatsManager struct {
	ChatsRepository *repositories.ChatsRepository
}

func (manager *ChatsManager) GetConcrete(chatID uint, token *jwt.Token) (*dtos.ChatDto, error) {
	log.Printf("Fetching concrete chat id = %d", chatID)
	tokenSubject, err := GetTokenSubject(token)

	if err != nil {
		log.Printf("Error fetching token subject: %v", err)
		return nil, err
	}

	chat, err := manager.ChatsRepository.GetWithMember(chatID, uint(tokenSubject.UserId))
	log.Printf("Fetched chat: %+v", chat)
	if err != nil {
		return nil, err
	}

	if chat.Type == "user" {
		setupUserChatData(chat, tokenSubject.UserId)
	}

	manager.setupChatActions(chat)

	log.Printf("Fetched chat with actions: %+v", chat)
	return chat, nil
}

func (manager *ChatsManager) GetByIds(chatsIds []int, token *jwt.Token) ([]*dtos.ChatDto, error) {
	log.Printf("Fetching chats by ids: %v", chatsIds)
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		log.Printf("Error fetching token subject: %v", err)
		return nil, err
	}

	chats, err := manager.ChatsRepository.GetByIds(chatsIds, tokenSubject.UserId)
	if err != nil {
		return nil, err
	}

	setupUserManyChatsData(chats, tokenSubject.UserId)
	manager.setupManyChatsActions(chats)

	return chats, nil
}

func (manager *ChatsManager) GetAll(token *jwt.Token, page int, perPage int) *schemas.PaginatedResponse[*dtos.ChatDto] {
	log.Printf("Fetching all chats for user")
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		log.Printf("Error parsing token subject: %v", err)
		return nil
	}

	log.Printf("Fetching all chats for user id = %d", tokenSubject.UserId)
	count := manager.ChatsRepository.GetAllWithMemberCount(uint(tokenSubject.UserId))
	countValue := *count
	chats := manager.ChatsRepository.GetAllWithMember(uint(tokenSubject.UserId), page, perPage)

	setupUserManyChatsData(chats, tokenSubject.UserId)
	manager.setupManyChatsActions(chats)

	log.Printf("Fetched %d chats for user id = %d", countValue, tokenSubject.UserId)
	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(countValue), chats)
	return &paginatedResponse
}

func (manager *ChatsManager) createGroupChat(chat *dtos.ChatDto, userId int) error {
	log.Printf("Creating group chat: %+v for user id = %d", chat, userId)
	chat.OwnerId = userId
	chat.Type = "group"
	if !slices.Contains(chat.Members, userId) {
		chat.Members = append(chat.Members, userId)
	}

	if err := manager.ChatsRepository.Create(chat); err != nil {
		log.Printf("Error creating chat: %v", err)
		return err
	}

	log.Printf("Created group chat: %+v", chat)
	return nil
}

func (manager *ChatsManager) validateUserChatExists(userId int, anotherUserId int) error {
	alreadyExists := manager.ChatsRepository.GetExistingWithUser(uint(userId), uint(anotherUserId))
	if alreadyExists {
		log.Printf("chat with members = [%d, %d] already exists", userId, anotherUserId)
		return fmt.Errorf("you already have chat with this user")
	}

	return nil
}

func (manager *ChatsManager) createUserChat(chat *dtos.ChatDto, currentUser *protousers.UserResponse, chatUser *protousers.UserResponse) error {
	log.Printf("Creating user chat %+v for current user %+v and chat user %+v", chat, currentUser, chatUser)
	chat.Members = []int{int(currentUser.Id), int(chatUser.Id)}
	chat.OwnerId = 0
	chat.Title = ""
	chat.Type = "user"

	err := manager.validateUserChatExists(int(currentUser.Id), int(chatUser.Id))
	if err != nil {
		return err
	}

	log.Printf("Creating chat: %+v", chat)

	if err := manager.ChatsRepository.Create(chat); err != nil {
		return err
	}

	setupChatTitleAndAvatar(chat, chatUser)
	log.Printf("Created user chat: %+v", chat)
	return nil
}

func (manager *ChatsManager) setupChatActions(chat *dtos.ChatDto) error {
	log.Printf("Setting up chat actions")
	ctx := context.Background()
	chatAllActions, err := redisdb.RedisConnection.HGetAll(ctx, manager.getChatActionsKey(chat.Id)).Result()
	if err != nil {
		log.Printf("%v", err)
		sentry.CaptureException(err)
		return fmt.Errorf("error when getting chat actions")
	}
	log.Printf("Fetched redis chat actions: %+v", chatAllActions)
	var chat_actions []dtos.ChatActionDto
	for key, value := range chatAllActions {
		var action_users_struct []dtos.ActionUserDto
		err = json.Unmarshal([]byte(value), &action_users_struct)
		if err != nil {
			log.Printf("%v", err)
			sentry.CaptureException(err)
			return fmt.Errorf("error when parsing action from redis")
		}

		chat_actions = append(chat_actions, dtos.ChatActionDto{
			Action:      key,
			ActionUsers: action_users_struct,
		})
	}

	chat.Actions = &chat_actions
	log.Printf("Chat with actions: %+v", chat)
	return nil
}

func (manager *ChatsManager) setupManyChatsActions(chats []*dtos.ChatDto) error {
	for _, chat := range chats {
		if err := manager.setupChatActions(chat); err != nil {
			return err
		}
	}

	return nil
}

func (manager *ChatsManager) sendChatEvent(chat *dtos.ChatDto, eventType string) error {
	log.Printf("Sending chat event for chat %+v event type = %s", chat, eventType)
	var included_users []int
	for _, user := range chat.Members {
		included_users = append(included_users, int(user))
	}

	chatEvent, err := rabbit.NewSystemEvent(eventType, included_users, chat)
	if err != nil {
		log.Printf("%v", err)
		sentry.CaptureException(err)
		return err
	}

	log.Printf("Sending event %+v", chatEvent)
	err = rabbit.EventsRabbitConnection.SendEvent(chatEvent)

	if err != nil {
		log.Printf("Error when publishing event with type %s in queue: %v", eventType, err)
		sentry.CaptureException(err)
		return fmt.Errorf("error sending event with type %s", eventType)
	}
	return nil
}

func (manager *ChatsManager) Create(chat *dtos.ChatDto, token *jwt.Token, chatUserId uint) error {
	log.Printf("Creating chat %+v. Chat user id = %d", chat, chatUserId)
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		log.Printf("Error decoding token subject: %v", err)
		return err
	}

	if chatUserId == 0 {
		err = manager.createGroupChat(chat, tokenSubject.UserId)
		if err != nil {
			return err
		}

		manager.sendChatEvent(chat, "chat_created")
		return nil
	}

	if chatUserId == uint(tokenSubject.UserId) {
		return fmt.Errorf("you can't create chat for you")
	}

	members := pq.Int32Array{int32(tokenSubject.UserId), int32(chatUserId)}
	if chatId := manager.ChatsRepository.GetDeletedChatId(members); chatId > 0 {
		manager.ChatsRepository.RestoreChat(chatId)
		if err := manager.sendChatEvent(chat, "chat_created"); err != nil {
			return err
		}

		return nil
	}

	currentUser, err := grpc_client.UsersGrpcClient.GetUserById(tokenSubject.UserId)

	if err != nil || currentUser == nil {
		log.Printf("Error fetching current user by id: %v. User = %+v", err, currentUser)
		sentry.CaptureException(err)
		return fmt.Errorf("there is no user with id %d", tokenSubject.UserId)
	}

	chatUser, err := grpc_client.UsersGrpcClient.GetUserById(int(chatUserId))

	if err != nil || chatUser == nil {
		log.Printf("Error fetching user by id: %v. User = %+v", err, chatUser)
		sentry.CaptureException(err)
		return fmt.Errorf("there is no user with id %d", chatUserId)
	}

	log.Printf("Creating user chat: %v, user id: %v, chat user: %v", chat, tokenSubject.UserId, chatUser)
	err = manager.createUserChat(chat, currentUser, chatUser)
	if err != nil {
		return err
	}

	manager.sendChatEvent(chat, "chat_created")
	return nil
}

func (manager *ChatsManager) Delete(chat *dtos.ChatDto) {
	manager.ChatsRepository.Delete(chat)
	manager.sendChatEvent(chat, "chat_deleted")
}

func (manager *ChatsManager) getChatActionsKey(chatId int) string {
	return fmt.Sprintf("chat_actions:%d", chatId)
}

func (manager *ChatsManager) getChatActionUsers(chatId int, actionType string) ([]dtos.ActionUserDto, error) {
	log.Printf("Fetching chat action users for chat id = %d action type = %s", chatId, actionType)
	ctx := context.Background()
	action_users, err := redisdb.RedisConnection.HGet(ctx, manager.getChatActionsKey(int(chatId)), actionType).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching chat action users: %v", err)
		sentry.CaptureException(err)
		return nil, fmt.Errorf("error getting action from redis")
	} else if err == redis.Nil {
		return []dtos.ActionUserDto{}, nil
	}

	log.Printf("Fetched chat action users: %+v", action_users)
	var action_users_struct []dtos.ActionUserDto
	err = json.Unmarshal([]byte(action_users), &action_users_struct)
	if err != nil {
		log.Printf("Error decoding chat action user: %v", err)
		sentry.CaptureException(err)
		return nil, fmt.Errorf("error when parsing action from redis")
	}

	log.Printf("Chat action users for chat id = %d and action type = %s: %+v", chatId, actionType, action_users_struct)
	return action_users_struct, nil
}

func (manager *ChatsManager) setupNewChatActionUsers(actionUsers []dtos.ActionUserDto, user *protousers.UserResponse, start bool) []dtos.ActionUserDto {
	log.Printf("Setting up new chat action users: action users = %+v user = %+v start = %v", actionUsers, user, start)
	user_full_name := utils.GetUserFullName(user.FirstName, user.LastName, &user.MiddleName)
	var new_action_users_struct []dtos.ActionUserDto
	if len(actionUsers) == 0 && start {
		new_action_users_struct = append(new_action_users_struct, dtos.ActionUserDto{Id: int(user.Id), Name: user_full_name})
		log.Printf("New chat action users: %+v", new_action_users_struct)
		return new_action_users_struct
	}

	for _, action_user := range actionUsers {
		if action_user.Id != int(user.Id) {
			new_action_users_struct = append(new_action_users_struct, dtos.ActionUserDto{Id: int(user.Id), Name: user_full_name})
			continue
		}

		if start {
			new_action_users_struct = append(new_action_users_struct, dtos.ActionUserDto{Id: int(user.Id), Name: user_full_name})
		} else {
			continue
		}
	}

	log.Printf("New action users: %+v", new_action_users_struct)
	return new_action_users_struct
}

func (manager *ChatsManager) saveNewChatActionUsers(chatId int, actionType string, actionUsers []dtos.ActionUserDto) error {
	log.Printf("Saving new chat action users. Chat id = %d action type = %s action users = %+v", chatId, actionType, actionUsers)
	ctx := context.Background()
	if len(actionUsers) == 0 {
		redisdb.RedisConnection.HDel(ctx, manager.getChatActionsKey(chatId), actionType)
		log.Printf("Deleted user actions for chat id = %d", chatId)
		return nil
	}

	encoded_users, err := json.Marshal(actionUsers)
	if err != nil {
		log.Printf("Error marshalling action users: %v", err)
		sentry.CaptureException(err)
		return fmt.Errorf("error when saving new chat action users")
	}

	err = redisdb.RedisConnection.HSet(ctx, manager.getChatActionsKey(chatId), []string{actionType, string(encoded_users)}).Err()
	if err != nil {
		log.Printf("Error saving new chat actions: %v", err)
		sentry.CaptureException(err)
		return fmt.Errorf("error when saving new chat action users")
	}

	return nil
}

func (manager *ChatsManager) HandleChatUserAction(token *jwt.Token, chat *dtos.ChatDto, actionType string, start bool) error {
	log.Printf("Handling chat user action: chat = %+v, actionType = %s, start = %v", chat, actionType, start)
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		log.Printf("Error decoding token subject: %v", err)
		return err
	}

	user, err := grpc_client.UsersGrpcClient.GetUserById(tokenSubject.UserId)
	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("error when getting token user")
	}

	action_users_struct, err := manager.getChatActionUsers(int(chat.Id), actionType)
	if err != nil {
		return err
	}

	new_action_users_struct := manager.setupNewChatActionUsers(action_users_struct, user, start)
	err = manager.saveNewChatActionUsers(int(chat.Id), actionType, new_action_users_struct)
	if err != nil {
		return err
	}

	manager.setupChatActions(chat)
	manager.sendChatEvent(chat, "chat_user_action")
	log.Printf("Handled new chat user actions: chat = %+v", chat)
	return nil
}

func (manager *ChatsManager) SystemSave(chatDto *dtos.ChatDto) (*dtos.ChatDto, error) {
	if chatDto == nil {
		return nil, fmt.Errorf("Error saving chat: chatDto is nil pointer")
	}

	manager.ChatsRepository.Create(chatDto)
	return chatDto, nil
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		ChatsRepository: &repositories.ChatsRepository{},
	}
}
