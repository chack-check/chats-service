package messages

import (
	"fmt"
	"slices"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/utils"
)

var (
	ErrMessageNotFound        = fmt.Errorf("message not found")
	ErrCantDeleteMessage      = fmt.Errorf("you can't delete message")
	ErrIncorrectCircleMessage = fmt.Errorf("you need to specify circle for circle message")
	ErrIncorrectVoiceMessage  = fmt.Errorf("you need to specify voice for voice message")
	ErrIncorrectTextMessage   = fmt.Errorf("you need to specify content or attachments for text message")
	ErrSavingMessage          = fmt.Errorf("error saving message")
)

type CreateMessageHandler struct {
	chatsPort         chats.ChatsPort
	messagesPort      MessagesPort
	messageEventsPort MessageEventsPort
	filesPort         files.FilesPort
}

func (handler *CreateMessageHandler) Execute(data CreateMessageData, userId int) (*Message, error) {
	chat, err := handler.chatsPort.GetByIdForUser(data.chatId, userId)
	if err != nil {
		return nil, chats.ErrChatNotFound
	}

	var savedAttachments []files.SavedFile
	for _, attachment := range data.GetAttachments() {
		if err := files.ValidateUploadingFile(handler.filesPort, &attachment, files.FileInChatFiletype, true); err != nil {
			return nil, err
		}

		savedAttachments = append(savedAttachments, files.UploadingFileToSavedFile(attachment))
	}

	var voiceSavedFile *files.SavedFile
	if voice := data.GetVoice(); voice != nil {
		if err := files.ValidateUploadingFile(handler.filesPort, voice, files.VoiceFiletype, true); err != nil {
			return nil, err
		}

		savedFile := files.UploadingFileToSavedFile(*voice)
		voiceSavedFile = &savedFile
	}

	var circleSavedFile *files.SavedFile
	if circle := data.GetCircle(); circle != nil {
		if err := files.ValidateUploadingFile(handler.filesPort, circle, files.CircleFiletype, true); err != nil {
			return nil, err
		}

		savedFile := files.UploadingFileToSavedFile(*circle)
		circleSavedFile = &savedFile
	}

	if data.GetType() == VoiceMessageType && data.GetVoice() == nil {
		return nil, ErrIncorrectVoiceMessage
	}

	if data.GetType() == CircleMessageType && data.GetCircle() == nil {
		return nil, ErrIncorrectCircleMessage
	}

	if data.GetType() == TextMessageType && data.GetContent() == nil && len(data.GetAttachments()) == 0 {
		return nil, ErrIncorrectTextMessage
	}

	message := NewMessage(
		0,
		userId,
		*chat,
		data.GetType(),
		data.GetContent(),
		voiceSavedFile,
		circleSavedFile,
		savedAttachments,
		data.GetReplyToId(),
		data.GetMentioned(),
		[]int{},
		[]MessageReaction{},
		[]int{},
		nil,
	)

	savedMessage, err := handler.messagesPort.Save(message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	handler.messageEventsPort.SendMessageCreated(*savedMessage)
	return savedMessage, nil
}

type GetConcreteMessageHandler struct {
	messagesPort MessagesPort
}

func (handler *GetConcreteMessageHandler) Execute(messageId int, userId int) (*Message, error) {
	message, err := handler.messagesPort.GetByIdForUser(messageId, userId)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	return message, nil
}

type GetMessagesByIdsHandler struct {
	messagesPort MessagesPort
}

func (handler *GetMessagesByIdsHandler) Execute(messageIds []int, userId int) []Message {
	messages := handler.messagesPort.GetByIdsForUser(messageIds, userId)
	return messages
}

type GetChatMessagesHandler struct {
	messagesPort MessagesPort
	chatsPort    chats.ChatsPort
}

func (handler *GetChatMessagesHandler) Execute(chatId int, userId int, offset int, limit int) (*utils.OffsetResponse[Message], error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, chats.ErrChatNotFound
	}

	messages := handler.messagesPort.GetChatAllForUser(chat.GetId(), userId, offset, limit)
	return &messages, nil
}

type GetChatMessagesByCursorHandler struct {
	messagesPort MessagesPort
	chatsPort    chats.ChatsPort
}

func (handler *GetChatMessagesByCursorHandler) Execute(chatId int, userId int, messageId int, aroundOffset int) (*utils.OffsetResponse[Message], error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, chats.ErrChatNotFound
	}

	messages := handler.messagesPort.GetChatCursorAllForUser(chat.GetId(), userId, messageId, aroundOffset)
	return &messages, nil
}

type GetChatsLastMessagesHandler struct {
	messagesPort MessagesPort
	chatsPort    chats.ChatsPort
}

func (handler *GetChatsLastMessagesHandler) Execute(chatIds []int, userId int) []Message {
	chats := handler.chatsPort.GetByIdsForUser(chatIds, userId)
	var fetchedChatIds []int
	for _, chat := range chats {
		fetchedChatIds = append(fetchedChatIds, chat.GetId())
	}

	messages := handler.messagesPort.GetChatsLast(fetchedChatIds, userId)
	return messages
}

type ReadMessageHandler struct {
	messagesPort      MessagesPort
	messageEventsPort MessageEventsPort
}

func (handler *ReadMessageHandler) Execute(messageId int, userId int) (*Message, error) {
	message, err := handler.messagesPort.GetByIdForUser(messageId, userId)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	if slices.Contains(message.GetReadedBy(), userId) {
		return message, nil
	}

	message.Read(userId)
	savedMessage, err := handler.messagesPort.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	handler.messageEventsPort.SendMessageReaded(*savedMessage)
	return savedMessage, nil
}

type ReactMessageHandler struct {
	messagesPort      MessagesPort
	messageEventsPort MessageEventsPort
}

func (handler *ReactMessageHandler) Execute(messageId int, userId int, content string) (*Message, error) {
	message, err := handler.messagesPort.GetByIdForUser(messageId, userId)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	reaction := NewMessageReaction(userId, content)

	if slices.Contains(message.GetReactions(), reaction) {
		return message, nil
	}

	message.AddReaction(reaction)
	savedMessage, err := handler.messagesPort.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	handler.messageEventsPort.SendMessageReacted(*savedMessage)
	return savedMessage, nil
}

type DeleteMessageReactionHandler struct {
	messagesPort      MessagesPort
	messageEventsPort MessageEventsPort
}

func (handler *DeleteMessageReactionHandler) Execute(messageId int, userId int) (*Message, error) {
	message, err := handler.messagesPort.GetByIdForUser(messageId, userId)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	var userReaction *MessageReaction
	for _, reaction := range message.GetReactions() {
		if reaction.GetUserId() == userId {
			userReaction = &reaction
		}
	}

	if userReaction == nil {
		return message, nil
	}

	message.RemoveReaction(*userReaction)
	savedMessage, err := handler.messagesPort.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	handler.messageEventsPort.SendReactionDeleted(*savedMessage)
	return savedMessage, nil
}

type UpdateMessageHandler struct {
	messagesPort      MessagesPort
	messageEventsPort MessageEventsPort
	filesPort         files.FilesPort
}

func (handler *UpdateMessageHandler) Execute(messageId int, userId int, data UpdateMessageData) (*Message, error) {
	message, err := handler.messagesPort.GetByIdForUser(messageId, userId)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	if content := data.GetContent(); content != nil {
		message.SetContent(content)
	}
	if attachments := data.GetAttachments(); len(attachments) > 0 {
		var savedFiles []files.SavedFile
		for _, attachment := range attachments {
			if err := files.ValidateUploadingFile(handler.filesPort, &attachment, files.FileInChatFiletype, true); err != nil {
				return nil, err
			}

			savedFiles = append(savedFiles, files.UploadingFileToSavedFile(attachment))
		}

		message.SetAttachments(savedFiles)
	}
	if mentioned := data.GetMentioned(); len(mentioned) > 0 {
		message.SetMentioned(mentioned)
	}

	savedMessage, err := handler.messagesPort.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	handler.messageEventsPort.SendMessageUpdated(*savedMessage)
	return savedMessage, nil
}

type DeleteMessageHandler struct {
	messagesPort      MessagesPort
	messageEventsPort MessageEventsPort
}

func (handler *DeleteMessageHandler) Execute(messageId int, userId int) error {
	message, err := handler.messagesPort.GetByIdForUser(messageId, userId)
	if err != nil {
		return ErrMessageNotFound
	}

	chat := message.GetChat()
	if !slices.Contains(chat.GetMembers(), userId) {
		return ErrCantDeleteMessage
	}

	handler.messagesPort.Delete(*message)
	handler.messageEventsPort.SendMessageDeleted(*message)
	return nil
}
