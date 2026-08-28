package filesservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"chats-service/configs"
	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
)

type FilesLoggingProvider struct {
	provider ports.FilesProvider
}

func (provider FilesLoggingProvider) GetSignatureForFile(filename string, systemFiletype constants.SystemFiletype) string {
	log.Printf("calculating signature for file: filename=%s, systemFiletype=%v", filename, systemFiletype)
	signature := provider.provider.GetSignatureForFile(filename, systemFiletype)
	log.Printf("calculated file signature: %s", signature)
	return signature
}

type FilesProvider struct{}

func (provider FilesProvider) GetSignatureForFile(filename string, systemFiletype constants.SystemFiletype) string {
	configuration := configs.GetFilesServiceConfiguration()
	fileHMAC := hmac.New(sha256.New, []byte(configuration.SignatureKey))
	fmt.Fprintf(fileHMAC, "%s:%s", filename, systemFiletype.String())
	hashsum := fileHMAC.Sum(nil)
	hexdigest := make([]byte, hex.EncodedLen(len(hashsum)))
	hex.Encode(hexdigest, hashsum)
	return string(hexdigest)
}

func NewFilesProvider() ports.FilesProvider {
	return FilesLoggingProvider{provider: FilesProvider{}}
}
