package rabbit

import (
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
)

type ChatEventsLoggingAdapter struct {
	repository ports.ChatEventsRepositoryPort
}

func (adapter ChatEventsLoggingAdapter) SendChatCreated(chat entities.Chat) {
	log.Printf("sending chat created event: %+v", chat)
	adapter.repository.SendChatCreated(chat)
}

func (adapter ChatEventsLoggingAdapter) SendChatDeleted(chat entities.Chat) {
	log.Printf("sending chat deleted event: %+v", chat)
	adapter.repository.SendChatDeleted(chat)
}

func (adapter ChatEventsLoggingAdapter) SendChatUserAction(chat entities.Chat) {
	log.Printf("sending chat user action event: %+v", chat)
	adapter.repository.SendChatUserAction(chat)
}

func (adapter ChatEventsLoggingAdapter) SendChatChanged(chat entities.Chat) {
	log.Printf("sending chat changed event: %+v", chat)
	adapter.repository.SendChatChanged(chat)
}

type ChatEventsAdapter struct {
	connection RabbitConnection
}

func (adapter ChatEventsAdapter) getSystemEventForChat(chat entities.Chat, eventType string) (*SystemEvent, error) {
	chatEvent := ChatToChatEvent(chat)
	systemEvent, err := NewSystemEvent(
		eventType,
		chat.GetMembers(),
		chatEvent,
	)
	if err != nil {
		return nil, err
	}

	return systemEvent, nil
}

func (adapter ChatEventsAdapter) sendChatEvent(chat entities.Chat, eventType string) {
	systemEvent, err := adapter.getSystemEventForChat(chat, eventType)
	if err != nil {
		return
	}

	adapter.connection.SendEvent(systemEvent)
}

func (adapter ChatEventsAdapter) SendChatCreated(chat entities.Chat) {
	adapter.sendChatEvent(chat, "chat_created")
}

func (adapter ChatEventsAdapter) SendChatDeleted(chat entities.Chat) {
	adapter.sendChatEvent(chat, "chat_deleted")
}

func (adapter ChatEventsAdapter) SendChatUserAction(chat entities.Chat) {
	adapter.sendChatEvent(chat, "chat_user_action")
}

func (adapter ChatEventsAdapter) SendChatChanged(chat entities.Chat) {
	adapter.sendChatEvent(chat, "chat_changed")
}

type MessageEventsLoggingAdapter struct {
	repository ports.MessageEventsRepositoryPort
}

func (adapter MessageEventsLoggingAdapter) SendMessageReacted(message entities.Message) {
	log.Printf("sending message reacted event: %+v", message)
	adapter.repository.SendMessageReacted(message)
}

func (adapter MessageEventsLoggingAdapter) SendReactionDeleted(message entities.Message) {
	log.Printf("sending message reaction deleted event: %+v", message)
	adapter.repository.SendReactionDeleted(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageReaded(message entities.Message) {
	log.Printf("sending message readed event: %+v", message)
	adapter.repository.SendMessageReaded(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageDeleted(message entities.Message) {
	log.Printf("sending message deleted event: %+v", message)
	adapter.repository.SendMessageDeleted(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageUpdated(message entities.Message) {
	log.Printf("sending message updated event: %+v", message)
	adapter.repository.SendMessageUpdated(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageCreated(message entities.Message) {
	log.Printf("sending message created event: %+v", message)
	adapter.repository.SendMessageCreated(message)
}

type MessageEventsAdapter struct {
	connection RabbitConnection
}

func (adapter MessageEventsAdapter) getSystemEventForMessage(message entities.Message, eventType string) (*SystemEvent, error) {
	messageEvent := MessageToMessageEvent(message)
	chat := message.GetChat()
	systemEvent, err := NewSystemEvent(
		eventType,
		chat.GetMembers(),
		messageEvent,
	)
	if err != nil {
		return nil, err
	}

	return systemEvent, nil
}

func (adapter MessageEventsAdapter) sendMessageEvent(message entities.Message, eventType string) {
	systemEvent, err := adapter.getSystemEventForMessage(message, eventType)
	if err != nil {
		return
	}

	adapter.connection.SendEvent(systemEvent)
}

func (adapter MessageEventsAdapter) SendMessageReacted(message entities.Message) {
	adapter.sendMessageEvent(message, "message_reacted")
}

func (adapter MessageEventsAdapter) SendReactionDeleted(message entities.Message) {
	adapter.sendMessageEvent(message, "message_reaction_deleted")
}

func (adapter MessageEventsAdapter) SendMessageReaded(message entities.Message) {
	adapter.sendMessageEvent(message, "message_readed")
}

func (adapter MessageEventsAdapter) SendMessageDeleted(message entities.Message) {
	adapter.sendMessageEvent(message, "message_deleted")
}

func (adapter MessageEventsAdapter) SendMessageUpdated(message entities.Message) {
	adapter.sendMessageEvent(message, "message_updated")
}

func (adapter MessageEventsAdapter) SendMessageCreated(message entities.Message) {
	adapter.sendMessageEvent(message, "message_created")
}

func NewChatEventsAdapter(connection RabbitConnection) ports.ChatEventsRepositoryPort {
	return ChatEventsLoggingAdapter{repository: ChatEventsAdapter{connection: connection}}
}

func NewMessageEventsAdapter(connection RabbitConnection) ports.MessageEventsRepositoryPort {
	return MessageEventsLoggingAdapter{repository: MessageEventsAdapter{connection: connection}}
}
