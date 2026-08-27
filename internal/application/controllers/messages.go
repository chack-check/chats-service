// TODO: Make files for every controller
package controllers

import (
	"fmt"
	"slices"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/dtos"
	"chats-service/internal/domain/entities"
	"chats-service/internal/domain/services"
)

// TODO: To errors package
var (
	ErrMessageNotFound        = fmt.Errorf("message not found")
	ErrCantDeleteMessage      = fmt.Errorf("you can't delete message")
	ErrIncorrectCircleMessage = fmt.Errorf("you need to specify circle for circle message")
	ErrIncorrectVoiceMessage  = fmt.Errorf("you need to specify voice for voice message")
	ErrIncorrectTextMessage   = fmt.Errorf("you need to specify content or attachments for text message")
	ErrSavingMessage          = fmt.Errorf("error saving message")
)

type CreateMessageController struct {
	chatsRepository         ports.ChatsRepositoryPort
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
	filesRepository         ports.FilesRepositoryPort
}

func NewCreateMessageController(
	chatsRepository ports.ChatsRepositoryPort,
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
	filesRepository ports.FilesRepositoryPort,
) *CreateMessageController {
	return &CreateMessageController{
		chatsRepository:         chatsRepository,
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
		filesRepository:         filesRepository,
	}
}

func (controller *CreateMessageController) Execute(data dtos.CreateMessageData, userID int) (*entities.Message, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(data.GetChatID(), userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	var savedAttachments []entities.SavedFile
	for _, attachment := range data.GetAttachments() {
		if err := services.ValidateUploadingFile(controller.filesRepository, &attachment, constants.FileInChatFiletype, true); err != nil {
			return nil, err
		}

		savedAttachments = append(savedAttachments, services.UploadingFileToSavedFile(attachment))
	}

	var voiceSavedFile *entities.SavedFile
	if voice := data.GetVoice(); voice != nil {
		if err := services.ValidateUploadingFile(controller.filesRepository, voice, constants.VoiceFiletype, true); err != nil {
			return nil, err
		}

		savedFile := services.UploadingFileToSavedFile(*voice)
		voiceSavedFile = &savedFile
	}

	var circleSavedFile *entities.SavedFile
	if circle := data.GetCircle(); circle != nil {
		if err := services.ValidateUploadingFile(controller.filesRepository, circle, constants.CircleFiletype, true); err != nil {
			return nil, err
		}

		savedFile := services.UploadingFileToSavedFile(*circle)
		circleSavedFile = &savedFile
	}

	if data.GetType() == constants.VoiceMessageType && data.GetVoice() == nil {
		return nil, ErrIncorrectVoiceMessage
	}

	if data.GetType() == constants.CircleMessageType && data.GetCircle() == nil {
		return nil, ErrIncorrectCircleMessage
	}

	if data.GetType() == constants.TextMessageType && data.GetContent() == nil && len(data.GetAttachments()) == 0 {
		return nil, ErrIncorrectTextMessage
	}

	message := entities.NewMessage(
		0,
		userID,
		*chat,
		data.GetType(),
		data.GetContent(),
		voiceSavedFile,
		circleSavedFile,
		savedAttachments,
		data.GetReplyToID(),
		data.GetMentioned(),
		[]int{},
		[]entities.MessageReaction{},
		[]int{},
		nil,
	)

	savedMessage, err := controller.messagesRepository.Save(message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	controller.messageEventsRepository.SendMessageCreated(*savedMessage)
	return savedMessage, nil
}

type GetConcreteMessageController struct {
	messagesRepository ports.MessagesRepositoryPort
}

func NewGetConcreteMessageController(
	messagesRepository ports.MessagesRepositoryPort,
) *GetConcreteMessageController {
	return &GetConcreteMessageController{messagesRepository: messagesRepository}
}

func (controller *GetConcreteMessageController) Execute(messageID int, userID int) (*entities.Message, error) {
	message, err := controller.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	return message, nil
}

type GetMessagesByIdsController struct {
	messagesRepository ports.MessagesRepositoryPort
}

func NewGetMessagesByIdsController(
	messagesRepository ports.MessagesRepositoryPort,
) *GetMessagesByIdsController {
	return &GetMessagesByIdsController{messagesRepository: messagesRepository}
}

func (controller *GetMessagesByIdsController) Execute(messageIds []int, userID int) []entities.Message {
	messages := controller.messagesRepository.GetByIdsForUser(messageIds, userID)
	return messages
}

type GetChatMessagesController struct {
	messagesRepository ports.MessagesRepositoryPort
	chatsRepository    ports.ChatsRepositoryPort
}

func NewGetChatMessagesController(
	messagesRepository ports.MessagesRepositoryPort,
	chatsRepository ports.ChatsRepositoryPort,
) *GetChatMessagesController {
	return &GetChatMessagesController{
		messagesRepository: messagesRepository,
		chatsRepository:    chatsRepository,
	}
}

func (controller *GetChatMessagesController) Execute(chatID int, userID int, offset int, limit int) (*dtos.OffsetResponse[entities.Message], error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	messages := controller.messagesRepository.GetChatAllForUser(chat.GetID(), userID, offset, limit)
	return &messages, nil
}

type GetChatMessagesByCursorController struct {
	messagesRepository ports.MessagesRepositoryPort
	chatsRepository    ports.ChatsRepositoryPort
}

func NewGetChatMessagesByCursorController(
	messagesRepository ports.MessagesRepositoryPort,
	chatsRepository ports.ChatsRepositoryPort,
) *GetChatMessagesByCursorController {
	return &GetChatMessagesByCursorController{
		messagesRepository: messagesRepository,
		chatsRepository:    chatsRepository,
	}
}

func (controller *GetChatMessagesByCursorController) Execute(chatID int, userID int, messageID int, aroundOffset int) (*dtos.OffsetResponse[entities.Message], error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	messages := controller.messagesRepository.GetChatCursorAllForUser(chat.GetID(), userID, messageID, aroundOffset)
	return &messages, nil
}

type GetChatsLastMessagesController struct {
	messagesRepository ports.MessagesRepositoryPort
	chatsRepository    ports.ChatsRepositoryPort
}

func NewGetChatsLastMessagesController(
	messagesRepository ports.MessagesRepositoryPort,
	chatsRepository ports.ChatsRepositoryPort,
) *GetChatsLastMessagesController {
	return &GetChatsLastMessagesController{
		messagesRepository: messagesRepository,
		chatsRepository:    chatsRepository,
	}
}

func (controller *GetChatsLastMessagesController) Execute(chatIds []int, userID int) []entities.Message {
	chats := controller.chatsRepository.GetByIdsForUser(chatIds, userID)
	fetchedChatIds := make([]int, 0, len(chats))
	for _, chat := range chats {
		fetchedChatIds = append(fetchedChatIds, chat.GetID())
	}

	messages := controller.messagesRepository.GetChatsLast(fetchedChatIds, userID)
	return messages
}

type ReadMessageController struct {
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
}

func NewReadMessageController(
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
) *ReadMessageController {
	return &ReadMessageController{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (controller *ReadMessageController) Execute(messageID int, userID int) (*entities.Message, error) {
	message, err := controller.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	if slices.Contains(message.GetReadedBy(), userID) {
		return message, nil
	}

	message.Read(userID)
	savedMessage, err := controller.messagesRepository.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	controller.messageEventsRepository.SendMessageReaded(*savedMessage)
	return savedMessage, nil
}

type ReactMessageController struct {
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
}

func NewReactMessageController(
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
) *ReactMessageController {
	return &ReactMessageController{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (controller *ReactMessageController) Execute(messageID int, userID int, content string) (*entities.Message, error) {
	message, err := controller.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	reaction := entities.NewMessageReaction(userID, content)

	if slices.Contains(message.GetReactions(), reaction) {
		return message, nil
	}

	message.AddReaction(reaction)
	savedMessage, err := controller.messagesRepository.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	controller.messageEventsRepository.SendMessageReacted(*savedMessage)
	return savedMessage, nil
}

type DeleteMessageReactionController struct {
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
}

func NewDeleteMessageReactionController(
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
) *DeleteMessageReactionController {
	return &DeleteMessageReactionController{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (controller *DeleteMessageReactionController) Execute(messageID int, userID int) (*entities.Message, error) {
	message, err := controller.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	var userReaction *entities.MessageReaction
	for _, reaction := range message.GetReactions() {
		if reaction.GetUserID() == userID {
			userReaction = &reaction
		}
	}

	if userReaction == nil {
		return message, nil
	}

	message.RemoveReaction(*userReaction)
	savedMessage, err := controller.messagesRepository.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	controller.messageEventsRepository.SendReactionDeleted(*savedMessage)
	return savedMessage, nil
}

type UpdateMessageController struct {
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
	filesRepository         ports.FilesRepositoryPort
}

func NewUpdateMessageController(
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
	filesRepository ports.FilesRepositoryPort,
) *UpdateMessageController {
	return &UpdateMessageController{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
		filesRepository:         filesRepository,
	}
}

func (controller *UpdateMessageController) Execute(messageID int, userID int, data dtos.UpdateMessageData) (*entities.Message, error) {
	message, err := controller.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, ErrMessageNotFound
	}

	if content := data.GetContent(); content != nil {
		message.SetContent(content)
	}
	if attachments := data.GetAttachments(); len(attachments) > 0 {
		var savedFiles []entities.SavedFile
		for _, attachment := range attachments {
			if err := services.ValidateUploadingFile(controller.filesRepository, &attachment, constants.FileInChatFiletype, true); err != nil {
				return nil, err
			}

			savedFiles = append(savedFiles, services.UploadingFileToSavedFile(attachment))
		}

		message.SetAttachments(savedFiles)
	}
	if mentioned := data.GetMentioned(); len(mentioned) > 0 {
		message.SetMentioned(mentioned)
	}

	savedMessage, err := controller.messagesRepository.Save(*message)
	if err != nil {
		return nil, ErrSavingMessage
	}

	controller.messageEventsRepository.SendMessageUpdated(*savedMessage)
	return savedMessage, nil
}

type DeleteMessageController struct {
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
}

func NewDeleteMessageController(
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
) *DeleteMessageController {
	return &DeleteMessageController{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (controller *DeleteMessageController) Execute(messageID int, userID int) error {
	message, err := controller.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return ErrMessageNotFound
	}

	chat := message.GetChat()
	if !slices.Contains(chat.GetMembers(), userID) {
		return ErrCantDeleteMessage
	}

	controller.messagesRepository.Delete(*message)
	controller.messageEventsRepository.SendMessageDeleted(*message)
	return nil
}

type RecognizeMessageController struct {
	messagesRepository      ports.MessagesRepositoryPort
	messageEventsRepository ports.MessageEventsRepositoryPort
}

func NewRecognizeMessageController(
	messagesRepository ports.MessagesRepositoryPort,
	messageEventsRepository ports.MessageEventsRepositoryPort,
) *RecognizeMessageController {
	return &RecognizeMessageController{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
	}
}

func (controller *RecognizeMessageController) Execute(messageID int, content string) error {
	message, err := controller.messagesRepository.GetById(messageID)
	if err != nil {
		return ErrMessageNotFound
	}

	message.SetContent(&content)
	_, err = controller.messagesRepository.Save(*message)
	if err != nil {
		return ErrSavingMessage
	}

	controller.messageEventsRepository.SendMessageUpdated(*message)
	return nil
}
