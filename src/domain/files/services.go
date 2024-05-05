package files

import (
	"fmt"
)

var (
	ErrFileRequired       error = fmt.Errorf("uploading file required")
	ErrIncorrectUsing     error = fmt.Errorf("incorrect using uploading file (incorrect system filetype)")
	ErrIncorrectSignature error = fmt.Errorf("incorrect file signature")
)

func ValidateUploadingFile(port FilesPort, file *UploadingFile, useFor SystemFiletype, required bool) error {
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

	originalHexdigest := port.GetSignatureForFile(file.GetOriginal().GetFilename(), file.GetOriginal().GetSystemFiletype())
	if originalHexdigest != file.GetOriginal().GetSignature() {
		return ErrIncorrectSignature
	}

	if file.GetConverted() != nil {
		convertedHexdigest := port.GetSignatureForFile(file.GetConverted().GetFilename(), file.GetConverted().GetSystemFiletype())
		if convertedHexdigest != file.GetConverted().GetSignature() {
			return ErrIncorrectSignature
		}
	}

	return nil
}
