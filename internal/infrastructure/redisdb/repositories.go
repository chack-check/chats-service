package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"

	"github.com/redis/go-redis/v9"
)

type RedisActionUser struct {
	ID         int
	LastName   string
	FirstName  string
	MiddleName *string
	Username   string
}

type UserActionsLoggingRepository struct {
	repository ports.UserActionsRepository
}

func (repository UserActionsLoggingRepository) AddChatActionUser(chat entities.Chat, user entities.User, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	log.Printf("adding chat action user: chat=%+v, user=%+v, actionType=%v", chat, user, actionType)
	actions := repository.repository.AddChatActionUser(chat, user, actionType)
	log.Printf("chat actions: %+v", actions)
	return actions
}

func (repository UserActionsLoggingRepository) RemoveChatActionUser(chat entities.Chat, userID int, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	log.Printf("removing chat action user: chat=%+v, userId=%d, actionType=%v", chat, userID, actionType)
	actions := repository.repository.RemoveChatActionUser(chat, userID, actionType)
	log.Printf("chat actions: %+v", actions)
	return actions
}

func (repository UserActionsLoggingRepository) GetAllChatActionsUsers(chat entities.Chat) map[constants.ActionTypes][]entities.ActionUser {
	log.Printf("fetching all chat actions users: chat=%+v", chat)
	actions := repository.repository.GetAllChatActionsUsers(chat)
	log.Printf("chat actions: %+v", actions)
	return actions
}

type UserActionsRepository struct {
	db *redis.Client
}

func (repository UserActionsRepository) getChatActionsKey(chatID int) string {
	return fmt.Sprintf("chat:%d:actions", chatID)
}

func (repository UserActionsRepository) AddChatActionUser(chat entities.Chat, user entities.User, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	actionUsers, err := repository.db.HGet(context.Background(), repository.getChatActionsKey(chat.GetID()), string(actionType)).Result()
	log.Printf("action users: %+v, err: %v", actionUsers, err)
	if err != nil && err != redis.Nil {
		return repository.GetAllChatActionsUsers(chat)
	}
	if actionUsers == "" || actionUsers == "null" {
		actionUsers = "[]"
	}

	users := make([]RedisActionUser, 0, 1)
	err = json.Unmarshal([]byte(actionUsers), &users)
	log.Printf("users: %+v, err: %v", users, err)
	if err != nil {
		return repository.GetAllChatActionsUsers(chat)
	}

	users = append(users, RedisActionUser{
		ID:         user.GetID(),
		LastName:   user.GetLastName(),
		FirstName:  user.GetFirstName(),
		MiddleName: user.GetMiddleName(),
		Username:   user.GetUsername(),
	})
	usersJSON, err := json.Marshal(users)
	log.Printf("users: %+v, usersJson: %s, err: %v", users, usersJSON, err)
	if err != nil {
		return repository.GetAllChatActionsUsers(chat)
	}

	repository.db.HSet(context.Background(), repository.getChatActionsKey(chat.GetID()), string(actionType), string(usersJSON)).Result()
	return repository.GetAllChatActionsUsers(chat)
}

func (repository UserActionsRepository) RemoveChatActionUser(chat entities.Chat, userID int, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	actionUsers, err := repository.db.HGet(context.Background(), repository.getChatActionsKey(chat.GetID()), string(actionType)).Result()
	if err != nil || actionUsers == "" {
		return repository.GetAllChatActionsUsers(chat)
	}

	var users []RedisActionUser
	err = json.Unmarshal([]byte(actionUsers), &users)
	if err != nil {
		return repository.GetAllChatActionsUsers(chat)
	}

	var resultUsers []RedisActionUser
	for _, user := range users {
		if user.ID == userID {
			continue
		}

		resultUsers = append(resultUsers, user)
	}

	usersJSON, err := json.Marshal(resultUsers)
	if err != nil {
		return repository.GetAllChatActionsUsers(chat)
	}

	repository.db.HSet(context.Background(), repository.getChatActionsKey(chat.GetID()), string(actionType), string(usersJSON))
	return repository.GetAllChatActionsUsers(chat)
}

func (repository UserActionsRepository) GetAllChatActionsUsers(chat entities.Chat) map[constants.ActionTypes][]entities.ActionUser {
	chatAllActions, err := repository.db.HGetAll(context.Background(), repository.getChatActionsKey(chat.GetID())).Result()
	if err != nil {
		return map[constants.ActionTypes][]entities.ActionUser{}
	}

	actions := make(map[constants.ActionTypes][]entities.ActionUser)
	for key, value := range chatAllActions {
		actionType := constants.ActionTypes(key)
		var actionUsers []RedisActionUser
		err = json.Unmarshal([]byte(value), &actionUsers)
		if err != nil {
			continue
		}

		var actionUsersSchemas []entities.ActionUser
		for _, user := range actionUsers {
			actionUsersSchemas = append(actionUsersSchemas, entities.NewActionUser(user.ID, user.LastName, user.FirstName, user.MiddleName, user.Username))
		}

		actions[actionType] = actionUsersSchemas
	}

	return actions
}

func NewUserActionsRepository(db *redis.Client) ports.UserActionsRepository {
	return UserActionsLoggingRepository{repository: UserActionsRepository{db: db}}
}
