package errors

import "errors"

var (
	ErrFileRequired       = errors.New("uploading file required")
	ErrIncorrectUsing     = errors.New("incorrect using uploading file (incorrect system filetype)")
	ErrIncorrectSignature = errors.New("incorrect file signature")
)
