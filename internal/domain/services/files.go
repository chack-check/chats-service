package services

import (
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/dtos"
	"chats-service/internal/domain/entities"
	"fmt"
)

var (
	ErrFileRequired       error = fmt.Errorf("uploading file required")
	ErrIncorrectUsing     error = fmt.Errorf("incorrect using uploading file (incorrect system filetype)")
	ErrIncorrectSignature error = fmt.Errorf("incorrect file signature")
)

// TODO: Move to application or refactor to not use application ports
func ValidateUploadingFile(repository ports.FilesRepositoryPort, file *dtos.UploadingFile, useFor constants.SystemFiletype, required bool) error {
	if file == nil && !required {
		return nil
	}

	if file == nil {
		return ErrFileRequired
	}

	if file.GetOriginal().GetSystemFiletype() != useFor {
		return ErrIncorrectUsing
	}
	if file.GetConverted() != nil && file.GetConverted().GetSystemFiletype() != useFor {
		return ErrIncorrectUsing
	}

	originalHexdigest := repository.GetSignatureForFile(file.GetOriginal().GetFilename(), file.GetOriginal().GetSystemFiletype())
	if originalHexdigest != file.GetOriginal().GetSignature() {
		return ErrIncorrectSignature
	}

	if file.GetConverted() != nil {
		convertedHexdigest := repository.GetSignatureForFile(file.GetConverted().GetFilename(), file.GetConverted().GetSystemFiletype())
		if convertedHexdigest != file.GetConverted().GetSignature() {
			return ErrIncorrectSignature
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
