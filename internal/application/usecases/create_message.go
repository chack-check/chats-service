package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/application/services"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type CreateMessageUseCase struct {
	chatsRepository         ports.ChatsRepository
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
	filesRepository         ports.FilesProvider
}

func NewCreateMessageUseCase(
	chatsRepository ports.ChatsRepository,
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
	filesRepository ports.FilesProvider,
) *CreateMessageUseCase {
	return &CreateMessageUseCase{
		chatsRepository:         chatsRepository,
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
		filesRepository:         filesRepository,
	}
}

func (useCase *CreateMessageUseCase) Execute(data dtos.CreateMessageData, userID int) (*entities.Message, error) {
	chat, err := useCase.chatsRepository.GetByIdForUser(data.GetChatID(), userID)
	if err != nil {
		return nil, domainerrors.ErrChatNotFound
	}

	var savedAttachments []entities.SavedFile
	for _, attachment := range data.GetAttachments() {
		if err := services.ValidateUploadingFile(useCase.filesRepository, &attachment, constants.FileInChatFiletype, true); err != nil {
			return nil, err
		}

		savedAttachments = append(savedAttachments, services.UploadingFileToSavedFile(attachment))
	}

	var voiceSavedFile *entities.SavedFile
	if voice := data.GetVoice(); voice != nil {
		if err := services.ValidateUploadingFile(useCase.filesRepository, voice, constants.VoiceFiletype, true); err != nil {
			return nil, err
		}

		savedFile := services.UploadingFileToSavedFile(*voice)
		voiceSavedFile = &savedFile
	}

	var circleSavedFile *entities.SavedFile
	if circle := data.GetCircle(); circle != nil {
		if err := services.ValidateUploadingFile(useCase.filesRepository, circle, constants.CircleFiletype, true); err != nil {
			return nil, err
		}

		savedFile := services.UploadingFileToSavedFile(*circle)
		circleSavedFile = &savedFile
	}

	if data.GetType() == constants.VoiceMessageType && data.GetVoice() == nil {
		return nil, domainerrors.ErrIncorrectVoiceMessage
	}

	if data.GetType() == constants.CircleMessageType && data.GetCircle() == nil {
		return nil, domainerrors.ErrIncorrectCircleMessage
	}

	if data.GetType() == constants.TextMessageType && data.GetContent() == nil && len(data.GetAttachments()) == 0 {
		return nil, domainerrors.ErrIncorrectTextMessage
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

	savedMessage, err := useCase.messagesRepository.Save(message)
	if err != nil {
		return nil, domainerrors.ErrSavingMessage
	}

	useCase.messageEventsRepository.SendMessageCreated(*savedMessage)
	return savedMessage, nil
}
