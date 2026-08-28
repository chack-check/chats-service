package chatsproto

import (
	"context"
	"fmt"

	"chats-service/internal/application/usecases"
	"chats-service/internal/infrastructure/database"
	"chats-service/internal/infrastructure/database/repositories"
	"chats-service/internal/infrastructure/grpc_service/chatsproto/chatsprotobuf"
	"chats-service/internal/infrastructure/grpc_service/usersproto"
	"chats-service/internal/infrastructure/redisdb"
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

	useCase := usecases.NewGetChatUseCase(
		repositories.NewChatsRepository(*database.DatabaseConnection),
		usersproto.NewUsersRepository(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsRepository(redisdb.RedisConnection),
	)

	chat, err := useCase.Execute(tokenSubject.UserId, int(request.Id))
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

	useCase := usecases.NewGetConcreteMessageUseCase(
		repositories.NewMessagesRepository(*database.DatabaseConnection),
	)

	message, err := useCase.Execute(int(request.Id), tokenSubject.UserId)
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

	useCase := usecases.NewGetChatsByIdsUseCase(
		repositories.NewChatsRepository(*database.DatabaseConnection),
		usersproto.NewUsersRepository(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsRepository(redisdb.RedisConnection),
	)

	var ids []int
	for _, id := range request.Ids {
		ids = append(ids, int(id))
	}

	chats := useCase.Execute(ids, tokenSubject.UserId)
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

	useCase := usecases.NewGetMessagesByIdsUseCase(
		repositories.NewMessagesRepository(*database.DatabaseConnection),
	)

	var ids []int
	for _, id := range request.Ids {
		ids = append(ids, int(id))
	}

	messages := useCase.Execute(ids, tokenSubject.UserId)

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

	useCase := usecases.NewGetChatMessagesUseCase(
		repositories.NewMessagesRepository(*database.DatabaseConnection),
		repositories.NewChatsRepository(*database.DatabaseConnection),
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

	messages, err := useCase.Execute(int(request.ChatId), tokenSubject.UserId, offsetValue, limitValue)
	if err != nil {
		return nil, err
	}

	messagesResponse := OffsetMessagesToProto(*messages)
	return messagesResponse, nil
}
