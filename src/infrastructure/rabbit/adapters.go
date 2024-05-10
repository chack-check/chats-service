package rabbit

import (
	"log"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/messages"
)

type ChatEventsLoggingAdapter struct {
	adapter chats.ChatEventsPort
}

func (adapter ChatEventsLoggingAdapter) SendChatCreated(chat chats.Chat) {
	log.Printf("sending chat created event: %+v", chat)
	adapter.adapter.SendChatCreated(chat)
}

func (adapter ChatEventsLoggingAdapter) SendChatDeleted(chat chats.Chat) {
	log.Printf("sending chat deleted event: %+v", chat)
	adapter.adapter.SendChatDeleted(chat)
}

func (adapter ChatEventsLoggingAdapter) SendChatUserAction(chat chats.Chat) {
	log.Printf("sending chat user action event: %+v", chat)
	adapter.adapter.SendChatUserAction(chat)
}

func (adapter ChatEventsLoggingAdapter) SendChatChanged(chat chats.Chat) {
	log.Printf("sending chat changed event: %+v", chat)
	adapter.adapter.SendChatChanged(chat)
}

type ChatEventsAdapter struct {
	connection RabbitConnection
}

func (adapter ChatEventsAdapter) getSystemEventForChat(chat chats.Chat, eventType string) (*SystemEvent, error) {
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

func (adapter ChatEventsAdapter) sendChatEvent(chat chats.Chat, eventType string) {
	systemEvent, err := adapter.getSystemEventForChat(chat, eventType)
	if err != nil {
		return
	}

	adapter.connection.SendEvent(systemEvent)
}

func (adapter ChatEventsAdapter) SendChatCreated(chat chats.Chat) {
	adapter.sendChatEvent(chat, "chat_created")
}

func (adapter ChatEventsAdapter) SendChatDeleted(chat chats.Chat) {
	adapter.sendChatEvent(chat, "chat_deleted")
}

func (adapter ChatEventsAdapter) SendChatUserAction(chat chats.Chat) {
	adapter.sendChatEvent(chat, "chat_user_action")
}

func (adapter ChatEventsAdapter) SendChatChanged(chat chats.Chat) {
	adapter.sendChatEvent(chat, "chat_changed")
}

type MessageEventsLoggingAdapter struct {
	adapter messages.MessageEventsPort
}

func (adapter MessageEventsLoggingAdapter) SendMessageReacted(message messages.Message) {
	log.Printf("sending message reacted event: %+v", message)
	adapter.adapter.SendMessageReacted(message)
}

func (adapter MessageEventsLoggingAdapter) SendReactionDeleted(message messages.Message) {
	log.Printf("sending message reaction deleted event: %+v", message)
	adapter.adapter.SendReactionDeleted(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageReaded(message messages.Message) {
	log.Printf("sending message readed event: %+v", message)
	adapter.adapter.SendMessageReaded(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageDeleted(message messages.Message) {
	log.Printf("sending message deleted event: %+v", message)
	adapter.adapter.SendMessageDeleted(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageUpdated(message messages.Message) {
	log.Printf("sending message updated event: %+v", message)
	adapter.adapter.SendMessageUpdated(message)
}

func (adapter MessageEventsLoggingAdapter) SendMessageCreated(message messages.Message) {
	log.Printf("sending message created event: %+v", message)
	adapter.adapter.SendMessageCreated(message)
}

type MessageEventsAdapter struct {
	connection RabbitConnection
}

func (adapter MessageEventsAdapter) getSystemEventForMessage(message messages.Message, eventType string) (*SystemEvent, error) {
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

func (adapter MessageEventsAdapter) sendMessageEvent(message messages.Message, eventType string) {
	systemEvent, err := adapter.getSystemEventForMessage(message, eventType)
	if err != nil {
		return
	}

	adapter.connection.SendEvent(systemEvent)
}

func (adapter MessageEventsAdapter) SendMessageReacted(message messages.Message) {
	adapter.sendMessageEvent(message, "message_reacted")
}

func (adapter MessageEventsAdapter) SendReactionDeleted(message messages.Message) {
	adapter.sendMessageEvent(message, "message_reaction_deleted")
}

func (adapter MessageEventsAdapter) SendMessageReaded(message messages.Message) {
	adapter.sendMessageEvent(message, "message_readed")
}

func (adapter MessageEventsAdapter) SendMessageDeleted(message messages.Message) {
	adapter.sendMessageEvent(message, "message_deleted")
}

func (adapter MessageEventsAdapter) SendMessageUpdated(message messages.Message) {
	adapter.sendMessageEvent(message, "message_updated")
}

func (adapter MessageEventsAdapter) SendMessageCreated(message messages.Message) {
	adapter.sendMessageEvent(message, "message_created")
}

func NewChatEventsAdapter(connection RabbitConnection) chats.ChatEventsPort {
	return ChatEventsLoggingAdapter{adapter: ChatEventsAdapter{connection: connection}}
}

func NewMessageEventsAdapter(connection RabbitConnection) messages.MessageEventsPort {
	return MessageEventsLoggingAdapter{adapter: MessageEventsAdapter{connection: connection}}
}
