package entities

import (
	"chats-service/internal/domain/constants"
	"slices"
	"time"
)

type MessageReaction struct {
	userID  int
	content string
}

func NewMessageReaction(userID int, content string) MessageReaction {
	return MessageReaction{
		userID:  userID,
		content: content,
	}
}

func (reaction *MessageReaction) GetUserID() int {
	return reaction.userID
}

func (reaction *MessageReaction) GetContent() string {
	return reaction.content
}

type Message struct {
	id            int
	senderID      int
	chat          Chat
	type_         constants.MessageTypes
	content       *string
	voice         *SavedFile
	circle        *SavedFile
	attachments   []SavedFile
	replyToID     *int
	mentioned     []int
	readedBy      []int
	reactions     []MessageReaction
	deletedForIds []int
	createdAt     *time.Time
}

func NewMessage(
	id int,
	senderID int,
	chat Chat,
	type_ constants.MessageTypes,
	content *string,
	voice *SavedFile,
	circle *SavedFile,
	attachments []SavedFile,
	replyToID *int,
	mentioned []int,
	readedBy []int,
	reactions []MessageReaction,
	deletedForIds []int,
	createdAt *time.Time,
) Message {
	return Message{
		id:            id,
		senderID:      senderID,
		chat:          chat,
		type_:         type_,
		content:       content,
		voice:         voice,
		circle:        circle,
		attachments:   attachments,
		replyToID:     replyToID,
		mentioned:     mentioned,
		readedBy:      readedBy,
		reactions:     reactions,
		deletedForIds: deletedForIds,
		createdAt:     createdAt,
	}
}

func (message *Message) GetID() int {
	return message.id
}

func (message *Message) GetSenderID() int {
	return message.senderID
}

func (message *Message) GetChat() Chat {
	return message.chat
}

func (message *Message) GetType() constants.MessageTypes {
	return message.type_
}

func (message *Message) GetContent() *string {
	return message.content
}

func (message *Message) SetContent(newContent *string) {
	message.content = newContent
}

func (message *Message) GetVoice() *SavedFile {
	return message.voice
}

func (message *Message) GetCircle() *SavedFile {
	return message.circle
}

func (message *Message) GetAttachments() []SavedFile {
	return message.attachments
}

func (message *Message) SetAttachments(attachments []SavedFile) {
	message.attachments = attachments
}

func (message *Message) GetReplyToID() *int {
	return message.replyToID
}

func (message *Message) GetMentioned() []int {
	return message.mentioned
}

func (message *Message) SetMentioned(mentioned []int) {
	message.mentioned = mentioned
}

func (message *Message) GetReadedBy() []int {
	return message.readedBy
}

func (message *Message) Read(userID int) {
	if slices.Contains(message.readedBy, userID) {
		return
	}

	message.readedBy = append(message.readedBy, userID)
}

func (message *Message) Unread(userID int) {
	var newReadedBy []int
	for _, user := range message.readedBy {
		if user == userID {
			continue
		}

		newReadedBy = append(newReadedBy, user)
	}

	message.readedBy = newReadedBy
}

func (message *Message) GetReactions() []MessageReaction {
	return message.reactions
}

func (message *Message) AddReaction(newReaction MessageReaction) {
	for _, reaction := range message.reactions {
		if reaction.userID == newReaction.userID {
			return
		}
	}

	message.reactions = append(message.reactions, newReaction)
}

func (message *Message) RemoveReaction(reaction MessageReaction) {
	var newReactions []MessageReaction
	for _, react := range message.reactions {
		if react.userID == reaction.userID {
			continue
		}

		newReactions = append(newReactions, react)
	}

	message.reactions = newReactions
}

func (message *Message) GetDeletedForIds() []int {
	return message.deletedForIds
}

func (message *Message) DeleteFor(users []int) {
	for _, user := range users {
		if slices.Contains(message.deletedForIds, user) {
			continue
		}

		message.deletedForIds = append(message.deletedForIds, user)
	}
}

func (message *Message) GetCreatedAt() *time.Time {
	return message.createdAt
}
