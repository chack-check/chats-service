package chatsproto

import (
	"context"
	"fmt"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/messages"
	"github.com/chack-check/chats-service/infrastructure/database"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/chatsproto/chatsprotobuf"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/usersproto"
	"github.com/chack-check/chats-service/infrastructure/redisdb"
)

var (
	ErrIncorrectToken = fmt.Errorf("incorrect token")
)

type ChatsServer struct {
	chatsprotobuf.ChatsServer
}

func (ChatsServer) GetChatById(ctx context.Context, request *chatsprotobuf.GetChatByIdRequest) (*chatsprotobuf.ChatResponse, error) {
	token, err := GetTokenFromString(request.Token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	chatsHandler := chats.NewGetChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	chat, err := chatsHandler.Execute(tokenSubject.UserId, int(request.Id))
	if err != nil {
		return nil, err
	}

	chatResponse := ChatModelToProto(*chat)
	return chatResponse, nil
}

func (ChatsServer) GetMessageById(ctx context.Context, request *chatsprotobuf.GetMessageByIdRequest) (*chatsprotobuf.MessageResponse, error) {
	token, err := GetTokenFromString(request.Token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	messagesHandler := messages.NewGetConcreteMessageHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
	)

	message, err := messagesHandler.Execute(int(request.Id), tokenSubject.UserId)
	if err != nil {
		return nil, err
	}

	messageResponse := MessageToProto(*message)
	return messageResponse, nil
}

func (ChatsServer) GetChatsByIds(ctx context.Context, request *chatsprotobuf.GetChatsByIdsRequest) (*chatsprotobuf.ChatsArrayResponse, error) {
	token, err := GetTokenFromString(request.Token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	chatsHandler := chats.NewGetChatsByIdsHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	var ids []int
	for _, id := range request.Ids {
		ids = append(ids, int(id))
	}

	chats := chatsHandler.Execute(ids, tokenSubject.UserId)
	var chatsResponse []*chatsprotobuf.ChatResponse
	for _, chat := range chats {
		chatsResponse = append(chatsResponse, ChatModelToProto(chat))
	}

	response := &chatsprotobuf.ChatsArrayResponse{Chats: chatsResponse}
	return response, nil
}

func (ChatsServer) GetMessagesByIds(ctx context.Context, request *chatsprotobuf.GetMessagesByIdsRequest) (*chatsprotobuf.MessagesArrayResponse, error) {
	token, err := GetTokenFromString(request.Token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	messagesHandler := messages.NewGetMessagesByidsHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
	)

	var ids []int
	for _, id := range request.Ids {
		ids = append(ids, int(id))
	}

	messages := messagesHandler.Execute(ids, tokenSubject.UserId)

	var messagesResponse []*chatsprotobuf.MessageResponse
	for _, message := range messages {
		messagesResponse = append(messagesResponse, MessageToProto(message))
	}

	response := &chatsprotobuf.MessagesArrayResponse{Messages: messagesResponse}
	return response, nil
}

func (ChatsServer) GetMessagesByChatId(ctx context.Context, request *chatsprotobuf.GetMessagesByChatIdRequest) (*chatsprotobuf.PaginatedMessages, error) {
	token, err := GetTokenFromString(request.Token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, ErrIncorrectToken
	}

	messagesHandler := messages.NewGetChatMessagesHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		database.NewMessagesAdapter(*database.DatabaseConnection),
	)

	var offsetValue int
	if request.Offset != nil && *request.Offset > 0 {
		offsetValue = int(*request.Offset)
	} else {
		offsetValue = 0
	}

	var limitValue int
	if request.Limit != nil && *request.Limit > 0 {
		limitValue = int(*request.Limit)
	} else {
		limitValue = 0
	}

	messages, err := messagesHandler.Execute(int(request.ChatId), tokenSubject.UserId, offsetValue, limitValue)
	if err != nil {
		return nil, err
	}

	messagesResponse := OffsetMessagesToProto(*messages)
	return messagesResponse, nil
}
