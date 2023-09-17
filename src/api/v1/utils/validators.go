package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/golang-jwt/jwt/v5"
)

func UserRequired(token *jwt.Token) error {
	if token == nil {
		log.Fatal("No token")
		return fmt.Errorf("Incorrect token")
	}

	exp, err := token.Claims.GetExpirationTime()
	if err == nil && token.Valid && exp.Unix() > time.Now().Unix() {
		log.Printf("Token not expired: %v", exp)
		return nil
	}

	log.Printf("Token expired: %v. Is valid: %v. Exp: %v. Now: %v", err, token.Valid, exp.Unix(), time.Now().Unix())
	return fmt.Errorf("Incorrect token")
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
