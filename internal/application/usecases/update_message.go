package usecases

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/application/services"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

type UpdateMessageUseCase struct {
	messagesRepository      ports.MessagesRepository
	messageEventsRepository ports.MessageEventsPublisher
	filesRepository         ports.FilesProvider
}

func NewUpdateMessageUseCase(
	messagesRepository ports.MessagesRepository,
	messageEventsRepository ports.MessageEventsPublisher,
	filesRepository ports.FilesProvider,
) *UpdateMessageUseCase {
	return &UpdateMessageUseCase{
		messagesRepository:      messagesRepository,
		messageEventsRepository: messageEventsRepository,
		filesRepository:         filesRepository,
	}
}

func (useCase *UpdateMessageUseCase) Execute(messageID int, userID int, data dtos.UpdateMessageData) (*entities.Message, error) {
	message, err := useCase.messagesRepository.GetByIdForUser(messageID, userID)
	if err != nil {
		return nil, domainerrors.ErrMessageNotFound
	}

	if content := data.GetContent(); content != nil {
		message.SetContent(content)
	}
	if attachments := data.GetAttachments(); len(attachments) > 0 {
		var savedFiles []entities.SavedFile
		for _, attachment := range attachments {
			if err := services.ValidateUploadingFile(useCase.filesRepository, &attachment, constants.FileInChatFiletype, true); err != nil {
				return nil, err
			}

			savedFiles = append(savedFiles, services.UploadingFileToSavedFile(attachment))
		}

		message.SetAttachments(savedFiles)
	}
	if mentioned := data.GetMentioned(); len(mentioned) > 0 {
		message.SetMentioned(mentioned)
	}

	savedMessage, err := useCase.messagesRepository.Save(*message)
	if err != nil {
		return nil, domainerrors.ErrSavingMessage
	}

	useCase.messageEventsRepository.SendMessageUpdated(*savedMessage)
	return savedMessage, nil
}
