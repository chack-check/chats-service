package utils

import (
	"fmt"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/protousers"
)

func UserRequired(user *protousers.UserResponse) error {
	if user == nil || user.Id == 0 {
		return fmt.Errorf("Incorrect token")
	}

	return nil
}

func ValidateTextMessage(message *model.CreateMessageRequest) error {
	if len(*message.Content) == 0 && len(message.Attachments) == 0 {
		return fmt.Errorf("You need to specify content or attachments for text message")
	}

	return nil
}

func ValidateVoiceMessage(message *model.CreateMessageRequest) error {
	if len(*message.Voice) == 0 {
		return fmt.Errorf("You need to specify voice url for voice message")
	}

	return nil
}

func ValidateCircleMessage(message *model.CreateMessageRequest) error {
	if len(*message.Circle) == 0 {
		return fmt.Errorf("You need to specify circle url for circle message")
	}

	return nil
}
