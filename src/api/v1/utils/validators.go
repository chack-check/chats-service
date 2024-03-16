package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/settings"
	"github.com/golang-jwt/jwt/v5"
)

type ValidatingFileMeta struct {
	SystemFiletype string
	Filename       string
	Signature      string
}

type ValidatingFile struct {
	Original  ValidatingFileMeta
	Converted *ValidatingFileMeta
}

func UserRequired(token *jwt.Token) error {
	if token == nil {
		log.Print("No token")
		return fmt.Errorf("incorrect token")
	}

	exp, err := token.Claims.GetExpirationTime()
	if err == nil && token.Valid && exp.Unix() > time.Now().Unix() {
		log.Printf("Token not expired: %v", exp)
		return nil
	}

	log.Printf("Token expired: %v. Is valid: %v. Exp: %v. Now: %v", err, token.Valid, exp.Unix(), time.Now().Unix())
	return fmt.Errorf("incorrect token")
}

func ValidateTextMessage(message *model.CreateMessageRequest) error {
	if len(*message.Content) == 0 && len(message.Attachments) == 0 {
		return fmt.Errorf("you need to specify content or attachments for text message")
	}

	return nil
}

func ValidateVoiceMessage(message *model.CreateMessageRequest) error {
	log.Printf("Message voice: %v", *message.Voice)
	if message.Voice == nil {
		return fmt.Errorf("you need to specify voice file for voice message")
	}

	return nil
}

func ValidateCircleMessage(message *model.CreateMessageRequest) error {
	if message.Circle == nil {
		return fmt.Errorf("you need to specify circle file for circle message")
	}

	return nil
}

func GetSignatureForFile(filename string, systemFiletype string) string {
	file_hmac := hmac.New(sha256.New, []byte(settings.Settings.FILES_SIGNATURE_KEY))
	file_hmac.Write([]byte(fmt.Sprintf("%s:%s", filename, systemFiletype)))
	hashsum := file_hmac.Sum(nil)
	hexdigest := make([]byte, hex.EncodedLen(len(hashsum)))
	hex.Encode(hexdigest, hashsum)
	return string(hexdigest)
}

func validateFile(file ValidatingFile, useFor string) error {
	if file.Original.SystemFiletype != useFor {
		return fmt.Errorf("you can't use file with system type %s for %s", file.Original.SystemFiletype, useFor)
	}
	if file.Converted != nil && file.Converted.SystemFiletype != useFor {
		return fmt.Errorf("you can't use file with system type %s for %s", file.Original.SystemFiletype, useFor)
	}

	original_hexdigest := GetSignatureForFile(file.Original.Filename, file.Original.SystemFiletype)
	if original_hexdigest != file.Original.Signature {
		return fmt.Errorf("incorrect signature for file")
	}

	if file.Converted != nil {
		converted_hexdigest := GetSignatureForFile(file.Converted.Filename, file.Converted.SystemFiletype)
		if converted_hexdigest != file.Converted.Signature {
			return fmt.Errorf("incorrect signature for file")
		}
	}

	return nil
}

func uploadingFileToValidatingFile(uploadingFile model.UploadingFile) ValidatingFile {
	var converted_file *ValidatingFileMeta
	if uploadingFile.Converted != nil {
		converted_file = &ValidatingFileMeta{
			Filename:       uploadingFile.Converted.Filename,
			Signature:      uploadingFile.Converted.Signature,
			SystemFiletype: uploadingFile.Converted.SystemFiletype.String(),
		}
	} else {
		converted_file = nil
	}

	return ValidatingFile{
		Original: ValidatingFileMeta{
			Filename:       uploadingFile.Original.Filename,
			Signature:      uploadingFile.Original.Signature,
			SystemFiletype: uploadingFile.Original.SystemFiletype.String(),
		},
		Converted: converted_file,
	}
}

func ValidateUploadingFile(uploadingFile model.UploadingFile, useFor string) error {
	validating_file := uploadingFileToValidatingFile(uploadingFile)
	return validateFile(validating_file, useFor)
}
