package messages

import (
	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/utils"
)

type MessagesPort interface {
	GetChatAllForUser(chatId int, userId int, offset int, limit int) utils.OffsetResponse[Message]
	GetChatCursorAllForUser(chatId int, userId int, messageId int, aroundOffset int) utils.OffsetResponse[Message]
	GetChatsLast(chatIds []int, userId int) []Message
	GetByIdForUser(messageId int, userId int) (*Message, error)
	GetByIdsForUser(messageIds []int, userId int) []Message
	Save(message Message) (*Message, error)
	Delete(message Message)
}

type MessageEventsPort interface {
	SendMessageReacted(message Message)
	SendReactionDeleted(message Message)
	SendMessageReaded(message Message)
	SendMessageDeleted(message Message)
	SendMessageUpdated(message Message)
	SendMessageCreated(message Message)
}

func NewCreateMessageHandler(
	chatsPort chats.ChatsPort,
	messagesPort MessagesPort,
	messageEventsPort MessageEventsPort,
	filesPort files.FilesPort,
) CreateMessageHandler {
	return CreateMessageHandler{
		chatsPort:         chatsPort,
		messagesPort:      messagesPort,
		messageEventsPort: messageEventsPort,
		filesPort:         filesPort,
	}
}

func NewUpdateMessageHandler(
	messagesPort MessagesPort,
	messageEventsPort MessageEventsPort,
	filesPort files.FilesPort,
) UpdateMessageHandler {
	return UpdateMessageHandler{
		messagesPort:      messagesPort,
		messageEventsPort: messageEventsPort,
		filesPort:         filesPort,
	}
}

func NewReadMessageHandler(
	messagesPort MessagesPort,
	messageEventsPort MessageEventsPort,
) ReadMessageHandler {
	return ReadMessageHandler{
		messagesPort:      messagesPort,
		messageEventsPort: messageEventsPort,
	}
}

func NewReactMessageHandler(
	messagesPort MessagesPort,
	messageEventsPort MessageEventsPort,
) ReactMessageHandler {
	return ReactMessageHandler{
		messagesPort:      messagesPort,
		messageEventsPort: messageEventsPort,
	}
}

func NewDeleteMessageReactionHandler(
	messagesPort MessagesPort,
	messageEventsPort MessageEventsPort,
) DeleteMessageReactionHandler {
	return DeleteMessageReactionHandler{
		messagesPort:      messagesPort,
		messageEventsPort: messageEventsPort,
	}
}

func NewDeleteMessageHandler(
	messagesPort MessagesPort,
	messageEventsPort MessageEventsPort,
) DeleteMessageHandler {
	return DeleteMessageHandler{
		messagesPort:      messagesPort,
		messageEventsPort: messageEventsPort,
	}
}

func NewGetChatMessagesHandler(
	chatsPort chats.ChatsPort,
	messagesPort MessagesPort,
) GetChatMessagesHandler {
	return GetChatMessagesHandler{
		messagesPort: messagesPort,
		chatsPort:    chatsPort,
	}
}

func NewGetChatMessagesByCursorHandler(
	chatsPort chats.ChatsPort,
	messagesPort MessagesPort,
) GetChatMessagesByCursorHandler {
	return GetChatMessagesByCursorHandler{
		messagesPort: messagesPort,
		chatsPort:    chatsPort,
	}
}

func NewGetChatsLastMessagesHandler(
	chatsPort chats.ChatsPort,
	messagesPort MessagesPort,
) GetChatsLastMessagesHandler {
	return GetChatsLastMessagesHandler{
		messagesPort: messagesPort,
		chatsPort:    chatsPort,
	}
}

func NewGetConcreteMessageHandler(
	messagesPort MessagesPort,
) GetConcreteMessageHandler {
	return GetConcreteMessageHandler{
		messagesPort: messagesPort,
	}
}

func NewGetMessagesByidsHandler(
	messagesPort MessagesPort,
) GetMessagesByIdsHandler {
	return GetMessagesByIdsHandler{
		messagesPort: messagesPort,
	}
}
