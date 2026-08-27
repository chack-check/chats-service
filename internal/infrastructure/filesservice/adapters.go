package filesservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
)

type FilesLoggingAdapter struct {
	repository ports.FilesRepositoryPort
}

func (adapter FilesLoggingAdapter) GetSignatureForFile(filename string, systemFiletype constants.SystemFiletype) string {
	log.Printf("calculating signature for file: filename=%s, systemFiletype=%v", filename, systemFiletype)
	signature := adapter.repository.GetSignatureForFile(filename, systemFiletype)
	log.Printf("calculated file signature: %s", signature)
	return signature
}

type FilesAdapter struct{}

func (adapter FilesAdapter) GetSignatureForFile(filename string, systemFiletype constants.SystemFiletype) string {
	fileHMAC := hmac.New(sha256.New, []byte(Settings.FILES_SIGNATURE_KEY))
	fmt.Fprintf(fileHMAC, "%s:%s", filename, systemFiletype.String())
	hashsum := fileHMAC.Sum(nil)
	hexdigest := make([]byte, hex.EncodedLen(len(hashsum)))
	hex.Encode(hexdigest, hashsum)
	return string(hexdigest)
}

func NewFilesAdapter() ports.FilesRepositoryPort {
	return FilesLoggingAdapter{repository: FilesAdapter{}}
}
