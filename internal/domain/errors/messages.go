package errors

import "errors"

var (
	ErrMessageNotFound        = errors.New("message not found")
	ErrCantDeleteMessage      = errors.New("you can't delete message")
	ErrIncorrectCircleMessage = errors.New("you need to specify circle for circle message")
	ErrIncorrectVoiceMessage  = errors.New("you need to specify voice for voice message")
	ErrIncorrectTextMessage   = errors.New("you need to specify content or attachments for text message")
	ErrSavingMessage          = errors.New("error saving message")
)
