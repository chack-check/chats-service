package services

import (
	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
	domainerrors "chats-service/internal/domain/errors"
)

func ValidateUploadingFile(repository ports.FilesProvider, file *dtos.UploadingFile, useFor constants.SystemFiletype, required bool) error {
	if file == nil && !required {
		return nil
	}

	if file == nil {
		return domainerrors.ErrFileRequired
	}

	if file.GetOriginal().GetSystemFiletype() != useFor {
		return domainerrors.ErrIncorrectUsing
	}
	if file.GetConverted() != nil && file.GetConverted().GetSystemFiletype() != useFor {
		return domainerrors.ErrIncorrectUsing
	}

	originalHexdigest := repository.GetSignatureForFile(file.GetOriginal().GetFilename(), file.GetOriginal().GetSystemFiletype())
	if originalHexdigest != file.GetOriginal().GetSignature() {
		return domainerrors.ErrIncorrectSignature
	}

	if file.GetConverted() != nil {
		convertedHexdigest := repository.GetSignatureForFile(file.GetConverted().GetFilename(), file.GetConverted().GetSystemFiletype())
		if convertedHexdigest != file.GetConverted().GetSignature() {
			return domainerrors.ErrIncorrectSignature
		}
	}

	return nil
}

func UploadingFileToSavedFile(file dtos.UploadingFile) entities.SavedFile {
	var convertedURL *string
	var convertedFilename *string
	if file.GetConverted() != nil {
		url := file.GetConverted().GetURL()
		filename := file.GetConverted().GetFilename()
		convertedURL = &url
		convertedFilename = &filename
	} else {
		convertedURL = nil
		convertedFilename = nil
	}

	return entities.NewSavedFile(
		file.GetOriginal().GetURL(),
		file.GetOriginal().GetFilename(),
		convertedURL,
		convertedFilename,
	)
}
