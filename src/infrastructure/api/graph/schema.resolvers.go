package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/messages"
	"github.com/chack-check/chats-service/infrastructure/api/factories"
	"github.com/chack-check/chats-service/infrastructure/api/graph/model"
	"github.com/chack-check/chats-service/infrastructure/api/middlewares"
	"github.com/chack-check/chats-service/infrastructure/api/utils"
	"github.com/chack-check/chats-service/infrastructure/database"
	"github.com/chack-check/chats-service/infrastructure/filesservice"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/usersproto"
	"github.com/chack-check/chats-service/infrastructure/rabbit"
	"github.com/chack-check/chats-service/infrastructure/redisdb"
	jwt "github.com/golang-jwt/jwt/v5"
)

// CreateMessage is the resolver for the createMessage field.
func (r *mutationResolver) CreateMessage(ctx context.Context, request model.CreateMessageRequest) (model.MessageErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewCreateMessageHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		database.NewMessagesAdapter(*database.DatabaseConnection),
		rabbit.NewMessageEventsAdapter(*rabbit.EventsRabbitConnection),
		filesservice.NewFilesAdapter(),
	)

	data := factories.CreateMessageRequestToModel(request)
	message, err := messagesHandler.Execute(data, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	messageResponse := factories.MessageModelToResponse(*message)
	return &messageResponse, nil
}

// EditMessage is the resolver for the editMessage field.
func (r *mutationResolver) EditMessage(ctx context.Context, messageID int, request model.ChangeMessageRequest) (model.MessageErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewUpdateMessageHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
		rabbit.NewMessageEventsAdapter(*rabbit.EventsRabbitConnection),
		filesservice.NewFilesAdapter(),
	)

	data := factories.UpdateMessageRequestToModel(request)
	message, err := messagesHandler.Execute(messageID, tokenSubject.UserId, data)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	messageResponse := factories.MessageModelToResponse(*message)
	return &messageResponse, nil
}

// CreateChat is the resolver for the createChat field.
func (r *mutationResolver) CreateChat(ctx context.Context, request model.CreateChatRequest) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewCreateChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		filesservice.NewFilesAdapter(),
	)

	var chatType chats.ChatTypes
	if request.User != nil {
		chatType = chats.UserChatType
	} else {
		chatType = chats.GroupChatType
	}

	data := factories.CreateChatRequestToModel(request, chatType)
	chat, err := chatsHandler.Execute(data, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	chatResponse := factories.ChatModelToResponse(*chat)
	return &chatResponse, nil
}

// ReadMessage is the resolver for the readMessage field.
func (r *mutationResolver) ReadMessage(ctx context.Context, messageID int) (model.MessageErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewReadMessageHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
		rabbit.NewMessageEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	message, err := messagesHandler.Execute(messageID, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	messageResponse := factories.MessageModelToResponse(*message)
	return &messageResponse, nil
}

// ReactMessage is the resolver for the reactMessage field.
func (r *mutationResolver) ReactMessage(ctx context.Context, messageID int, content string) (model.MessageErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewReactMessageHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
		rabbit.NewMessageEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	message, err := messagesHandler.Execute(messageID, tokenSubject.UserId, content)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	messageResponse := factories.MessageModelToResponse(*message)
	return &messageResponse, nil
}

// DeleteMessageReaction is the resolver for the deleteMessageReaction field.
func (r *mutationResolver) DeleteMessageReaction(ctx context.Context, messageID int) (model.MessageErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewDeleteMessageReactionHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
		rabbit.NewMessageEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	message, err := messagesHandler.Execute(messageID, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	messageResponse := factories.MessageModelToResponse(*message)
	return &messageResponse, nil
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(ctx context.Context, messageID int) (model.BooleanResultErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewDeleteMessageHandler(
		database.NewMessagesAdapter(*database.DatabaseConnection),
		rabbit.NewMessageEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	err = messagesHandler.Execute(messageID, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return model.BooleanResult{Result: true}, nil
}

// DeleteChat is the resolver for the deleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, chatID int) (model.BooleanResultErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewDeleteChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	err = chatsHandler.Execute(chatID, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return model.BooleanResult{Result: true}, nil
}

// SendUserAction is the resolver for the sendUserAction field.
func (r *mutationResolver) SendUserAction(ctx context.Context, chatID int, actionType model.ActionTypes) (model.BooleanResultErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewUserActionHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	_, err = chatsHandler.Execute(chatID, tokenSubject.UserId, chats.ActionTypes(actionType.String()))
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return model.BooleanResult{Result: true}, nil
}

// StopUserAction is the resolver for the stopUserAction field.
func (r *mutationResolver) StopUserAction(ctx context.Context, chatID int, actionType model.ActionTypes) (model.BooleanResultErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewStopUserActionHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	_, err = chatsHandler.Execute(chatID, tokenSubject.UserId, chats.ActionTypes(actionType.String()))
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return model.BooleanResult{Result: true}, nil
}

// AddMembers is the resolver for the addMembers field.
func (r *mutationResolver) AddMembers(ctx context.Context, chatID int, members []int) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewAddChatsMembersHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId, members)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// AddAdmins is the resolver for the addAdmins field.
func (r *mutationResolver) AddAdmins(ctx context.Context, chatID int, admins []int) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewAddChatsAdminsHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId, admins)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// RemoveMembers is the resolver for the removeMembers field.
func (r *mutationResolver) RemoveMembers(ctx context.Context, chatID int, members []int) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewRemoveChatMembersHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId, members)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// RemoveAdmins is the resolver for the removeAdmins field.
func (r *mutationResolver) RemoveAdmins(ctx context.Context, chatID int, admins []int) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewRemoveChatAdminsHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId, admins)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// QuitChat is the resolver for the quitChat field.
func (r *mutationResolver) QuitChat(ctx context.Context, chatID int) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewQuitChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// ChangeGroupChat is the resolver for the changeGroupChat field.
func (r *mutationResolver) ChangeGroupChat(ctx context.Context, chatID int, chatData model.ChangeGroupChatData) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewChangeGroupChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId, chats.NewChangeGroupChatData(chatData.Title))
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// UpdateGroupChatAvatar is the resolver for the updateGroupChatAvatar field.
func (r *mutationResolver) UpdateGroupChatAvatar(ctx context.Context, chatID int, avatar model.UploadingFile) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewUpdateGroupChatAvatar(
		database.NewChatsAdapter(*database.DatabaseConnection),
		filesservice.NewFilesAdapter(),
		rabbit.NewChatEventsAdapter(*rabbit.EventsRabbitConnection),
	)

	chat, err := chatsHandler.Execute(chatID, tokenSubject.UserId, factories.UploadingFileToModel(avatar))
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	return factories.ChatModelToResponse(*chat), nil
}

// GetChatMessages is the resolver for the getChatMessages field.
func (r *queryResolver) GetChatMessages(ctx context.Context, chatID int, offset *int, limit *int) (model.PaginatedMessagesErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewGetChatMessagesHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		database.NewMessagesAdapter(*database.DatabaseConnection),
	)

	var offsetValue int
	if offset != nil && *offset > 0 {
		offsetValue = *offset
	} else {
		offsetValue = 0
	}

	var limitValue int
	if limit != nil && *limit > 0 {
		limitValue = *limit
	} else {
		limitValue = 100
	}

	messages, err := messagesHandler.Execute(chatID, tokenSubject.UserId, offsetValue, limitValue)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	response := factories.OffsetMessagesToResponse(*messages, chatID)
	return &response, nil
}

// GetChatMessagesByCursor is the resolver for the getChatMessagesByCursor field.
func (r *queryResolver) GetChatMessagesByCursor(ctx context.Context, chatID int, messageID int, aroundOffset *int) (model.PaginatedMessagesErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewGetChatMessagesByCursorHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		database.NewMessagesAdapter(*database.DatabaseConnection),
	)

	var aroundOffsetValue int
	if aroundOffset != nil && *aroundOffset > 0 {
		aroundOffsetValue = *aroundOffset
	} else {
		aroundOffsetValue = 50
	}

	messages, err := messagesHandler.Execute(chatID, tokenSubject.UserId, messageID, aroundOffsetValue)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	response := factories.OffsetMessagesToResponse(*messages, chatID)
	return &response, nil
}

// GetChats is the resolver for the getChats field.
func (r *queryResolver) GetChats(ctx context.Context, page *int, perPage *int) (model.PaginatedChatsErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewGetChatsHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	var pageValue int
	var perPageValue int
	if page != nil && *page > 0 {
		pageValue = *page
	} else {
		pageValue = 1
	}

	if perPage != nil && *perPage > 0 {
		perPageValue = *perPage
	} else {
		perPageValue = 20
	}

	chats := chatsHandler.Execute(tokenSubject.UserId, pageValue, perPageValue)
	chatsResponse := factories.PaginatedChatsToResponse(chats)
	return &chatsResponse, nil
}

// GetChat is the resolver for the getChat field.
func (r *queryResolver) GetChat(ctx context.Context, chatID int) (model.ChatErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	chatsHandler := chats.NewGetChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	chat, err := chatsHandler.Execute(tokenSubject.UserId, chatID)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	chatResponse := factories.ChatModelToResponse(*chat)
	return &chatResponse, nil
}

// GetLastMessagesForChats is the resolver for the getLastMessagesForChats field.
func (r *queryResolver) GetLastMessagesForChats(ctx context.Context, chatIds []int) (model.MessagesArrayErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	messagesHandler := messages.NewGetChatsLastMessagesHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		database.NewMessagesAdapter(*database.DatabaseConnection),
	)

	messages := messagesHandler.Execute(chatIds, tokenSubject.UserId)
	if err != nil {
		return model.ErrorResponse{Message: err.Error()}, nil
	}

	var response []*model.Message
	for _, message := range messages {
		messageResponse := factories.MessageModelToResponse(message)
		response = append(response, &messageResponse)
	}
	return model.MessagesArray{Messages: response}, nil
}

// SearchChats is the resolver for the searchChats field.
func (r *queryResolver) SearchChats(ctx context.Context, query string, page *int, perPage *int) (model.PaginatedChatsErrorResponse, error) {
	token, _ := ctx.Value("token").(*jwt.Token)
	if err := utils.UserRequired(token); err != nil {
		return model.ErrorResponse{Message: "Token required"}, nil
	}

	tokenSubject, err := middlewares.GetTokenSubject(token)
	if err != nil {
		return model.ErrorResponse{Message: "Incorrect token"}, nil
	}

	var pageValue int = 1
	var perPageValue int = 100
	if page != nil && *page > 0 {
		pageValue = *page
	}
	if perPage != nil && *perPage > 0 {
		perPageValue = *perPage
	}

	searchHandler := chats.NewSearchChatsHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
		usersproto.NewUsersAdapter(usersproto.UsersClientConnect()),
		redisdb.NewUserActionsAdapter(redisdb.RedisConnection),
	)

	chats := searchHandler.Execute(tokenSubject.UserId, query, pageValue, perPageValue)
	var response []*model.Chat
	for _, chat := range chats.GetData() {
		chatResponse := factories.ChatModelToResponse(chat)
		response = append(response, &chatResponse)
	}
	return model.PaginatedChats{Page: chats.GetPage(), NumPages: chats.GetPagesCount(), PerPage: chats.GetPerPage(), Total: chats.GetTotal(), Data: response}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
