package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"slices"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/services"
	"github.com/chack-check/chats-service/generic_factories"
	pb "github.com/chack-check/chats-service/protochats"
	"google.golang.org/grpc"
)

type ChatsServerImplementation struct {
	pb.ChatsServer
}

func (s ChatsServerImplementation) GetChatById(ctx context.Context, request *pb.GetChatByIdRequest) (*pb.ChatResponse, error) {
	token, err := services.GetTokenFromString(request.Token)
	if err != nil {
		return nil, err
	}

	chats_service := services.NewChatsManager()
	chat, err := chats_service.GetConcrete(uint(request.Id), token)

	if err != nil {
		return nil, err
	}

	chats_factory := generic_factories.ChatsFactory{}
	chat_response := chats_factory.SchemaToProto(chat)
	return chat_response, nil
}

func (s ChatsServerImplementation) GetMessageById(ctx context.Context, request *pb.GetMessageByIdRequest) (*pb.MessageResponse, error) {
	token, err := services.GetTokenFromString(request.Token)
	if err != nil {
		return nil, err
	}

	messages_service := services.NewMessagesManager()
	chats_service := services.NewChatsManager()

	message, err := messages_service.GetConcrete(int(request.Id), token)
	if err != nil {
		return nil, err
	}

	chat, err := chats_service.GetConcrete(message.ChatId, token)
	if err != nil {
		return nil, err
	}

	messages_factory := generic_factories.MessagesFactory{}
	message_response := messages_factory.SchemaToProto(message, chat)
	return message_response, nil
}

func (s ChatsServerImplementation) GetMessagesByIds(ctx context.Context, request *pb.GetMessagesByIdsRequest) (*pb.MessagesArrayResponse, error) {
	token, err := services.GetTokenFromString(request.Token)
	if err != nil {
		return nil, err
	}

	messages_service := services.NewMessagesManager()
	chats_service := services.NewChatsManager()

	var messages_ids []int
	for _, id := range request.Ids {
		if slices.Contains(messages_ids, int(id)) {
			continue
		}

		messages_ids = append(messages_ids, int(id))
	}

	messages, err := messages_service.GetByIds(messages_ids, token)
	if err != nil {
		return nil, err
	} else if messages == nil {
		return &pb.MessagesArrayResponse{Messages: []*pb.MessageResponse{}}, nil
	}

	var chat_ids []int
	for _, message := range messages {
		if slices.Contains(chat_ids, int(message.ChatId)) {
			continue
		}

		chat_ids = append(chat_ids, int(message.ChatId))
	}

	chats, err := chats_service.GetByIds(chat_ids, token)
	if err != nil {
		return nil, err
	}

	messages_factory := generic_factories.MessagesFactory{}
	var messages_response []*pb.MessageResponse
	for _, message := range messages {
		var message_chat *dtos.ChatDto
		for _, chat := range chats {
			if chat.Id == int(message.ChatId) {
				message_chat = chat
			}
		}

		messages_response = append(messages_response, messages_factory.SchemaToProto(message, message_chat))
	}

	return &pb.MessagesArrayResponse{Messages: messages_response}, nil
}

func (s ChatsServerImplementation) GetChatsByIds(ctx context.Context, request *pb.GetChatsByIdsRequest) (*pb.ChatsArrayResponse, error) {
	token, err := services.GetTokenFromString(request.Token)
	if err != nil {
		return nil, err
	}

	chats_service := services.NewChatsManager()

	log.Printf("getting chats for request %v", request)
	var chats_ids []int
	for _, id := range request.Ids {
		chats_ids = append(chats_ids, int(id))
	}

	chats, err := chats_service.GetByIds(chats_ids, token)
	if err != nil {
		return nil, err
	} else if chats == nil {
		return &pb.ChatsArrayResponse{Chats: []*pb.ChatResponse{}}, nil
	}

	chats_factory := generic_factories.ChatsFactory{}
	var chats_response []*pb.ChatResponse
	for _, chat := range chats {
		chats_response = append(chats_response, chats_factory.SchemaToProto(chat))
	}

	return &pb.ChatsArrayResponse{Chats: chats_response}, nil
}

func (ChatsServerImplementation) mustEmbedUnimplementedChatsServer() {}

func StartServer(host string, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpc_server := grpc.NewServer(opts...)
	chats_server_implementation := ChatsServerImplementation{}
	pb.RegisterChatsServer(grpc_server, chats_server_implementation)
	grpc_server.Serve(lis)
}
