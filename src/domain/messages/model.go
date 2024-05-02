package messages

import (
	"slices"
	"time"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
)

type MessageTypes string

const (
	TextMessageType   MessageTypes = "text"
	EventMessageType  MessageTypes = "event"
	CallMessageType   MessageTypes = "call"
	VoiceMessageType  MessageTypes = "voice"
	CircleMessageType MessageTypes = "circle"
)

type CreateMessageData struct {
	chatId      int
	type_       MessageTypes
	content     *string
	voice       *files.UploadingFile
	attachments []files.UploadingFile
	replyToId   *int
	mentioned   []int
	circle      *files.UploadingFile
}

func (model *CreateMessageData) GetChatId() int {
	return model.chatId
}

func (model *CreateMessageData) GetType() MessageTypes {
	return model.type_
}

func (model *CreateMessageData) GetContent() *string {
	return model.content
}

func (model *CreateMessageData) GetVoice() *files.UploadingFile {
	return model.voice
}

func (model *CreateMessageData) GetAttachments() []files.UploadingFile {
	return model.attachments
}

func (model *CreateMessageData) GetReplyToId() *int {
	return model.replyToId
}

func (model *CreateMessageData) GetMentioned() []int {
	return model.mentioned
}

func (model *CreateMessageData) GetCircle() *files.UploadingFile {
	return model.circle
}

type UpdateMessageData struct {
	content     *string
	attachments []files.UploadingFile
	mentioned   []int
}

func (model *UpdateMessageData) GetContent() *string {
	return model.content
}

func (model *UpdateMessageData) GetAttachments() []files.UploadingFile {
	return model.attachments
}

func (model *UpdateMessageData) GetMentioned() []int {
	return model.mentioned
}

type MessageReaction struct {
	userId  int
	content string
}

func (model *MessageReaction) GetUserId() int {
	return model.userId
}

func (model *MessageReaction) GetContent() string {
	return model.content
}

type Message struct {
	id            int
	senderId      int
	chat          chats.Chat
	type_         MessageTypes
	content       *string
	voice         *files.SavedFile
	circle        *files.SavedFile
	attachments   []files.SavedFile
	replyToId     *int
	mentioned     []int
	readedBy      []int
	reactions     []MessageReaction
	deletedForIds []int
	createdAt     *time.Time
}

func (model *Message) GetId() int {
	return model.id
}

func (model *Message) GetSenderId() int {
	return model.senderId
}

func (model *Message) GetChat() chats.Chat {
	return model.chat
}

func (model *Message) GetType() MessageTypes {
	return model.type_
}

func (model *Message) GetContent() *string {
	return model.content
}

func (model *Message) SetContent(newContent *string) {
	model.content = newContent
}

func (model *Message) GetVoice() *files.SavedFile {
	return model.voice
}

func (model *Message) GetCircle() *files.SavedFile {
	return model.circle
}

func (model *Message) GetAttachments() []files.SavedFile {
	return model.attachments
}

func (model *Message) SetAttachments(attachments []files.SavedFile) {
	model.attachments = attachments
}

func (model *Message) GetReplyToId() *int {
	return model.replyToId
}

func (model *Message) GetMentioned() []int {
	return model.mentioned
}

func (model *Message) SetMentioned(mentioned []int) {
	model.mentioned = mentioned
}

func (model *Message) GetReadedBy() []int {
	return model.readedBy
}

func (model *Message) Read(userId int) {
	if slices.Contains(model.readedBy, userId) {
		return
	}

	model.readedBy = append(model.readedBy, userId)
}

func (model *Message) Unread(userId int) {
	var newReadedBy []int
	for _, user := range model.readedBy {
		if user == userId {
			continue
		}

		newReadedBy = append(newReadedBy, user)
	}

	model.readedBy = newReadedBy
}

func (model *Message) GetReactions() []MessageReaction {
	return model.reactions
}

func (model *Message) AddReaction(newReaction MessageReaction) {
	for _, reaction := range model.reactions {
		if reaction.GetUserId() == newReaction.GetUserId() {
			return
		}
	}

	model.reactions = append(model.reactions, newReaction)
}

func (model *Message) RemoveReaction(reaction MessageReaction) {
	var newReactions []MessageReaction
	for _, react := range model.reactions {
		if react.userId == reaction.userId {
			continue
		}

		newReactions = append(newReactions, react)
	}

	model.reactions = newReactions
}

func (model *Message) GetDeletedForIds() []int {
	return model.deletedForIds
}

func (model *Message) DeleteFor(users []int) {
	for _, user := range users {
		if slices.Contains(model.deletedForIds, user) {
			continue
		}

		model.deletedForIds = append(model.deletedForIds, user)
	}
}

func (model *Message) GetCreatedAt() *time.Time {
	return model.createdAt
}

func NewMessageReaction(userId int, content string) MessageReaction {
	return MessageReaction{
		userId:  userId,
		content: content,
	}
}

func NewMessage(
	id int,
	senderId int,
	chat chats.Chat,
	type_ MessageTypes,
	content *string,
	voice *files.SavedFile,
	circle *files.SavedFile,
	attachments []files.SavedFile,
	replyToId *int,
	mentioned []int,
	readedBy []int,
	reactions []MessageReaction,
	deletedForIds []int,
	createdAt *time.Time,
) Message {
	return Message{
		id:            id,
		senderId:      senderId,
		chat:          chat,
		type_:         type_,
		content:       content,
		voice:         voice,
		circle:        circle,
		attachments:   attachments,
		replyToId:     replyToId,
		mentioned:     mentioned,
		readedBy:      readedBy,
		reactions:     reactions,
		deletedForIds: deletedForIds,
		createdAt:     createdAt,
	}
}

func NewCreateMessageData(
	chatId int,
	type_ MessageTypes,
	content *string,
	voice *files.UploadingFile,
	attachments []files.UploadingFile,
	replyToId *int,
	mentioned []int,
	circle *files.UploadingFile,
) CreateMessageData {
	return CreateMessageData{
		chatId:      chatId,
		type_:       type_,
		content:     content,
		voice:       voice,
		attachments: attachments,
		replyToId:   replyToId,
		mentioned:   mentioned,
		circle:      circle,
	}
}

func NewUpdateMessageData(
	content *string,
	attachments []files.UploadingFile,
	mentioned []int,
) UpdateMessageData {
	return UpdateMessageData{
		content:     content,
		attachments: attachments,
		mentioned:   mentioned,
	}
}
