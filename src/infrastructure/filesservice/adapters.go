package filesservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/chack-check/chats-service/domain/files"
)

type FilesLoggingAdapter struct {
	adapter files.FilesPort
}

func (adapter FilesLoggingAdapter) GetSignatureForFile(filename string, systemFiletype files.SystemFiletype) string {
	log.Printf("calculating signature for file: filename=%s, systemFiletype=%v", filename, systemFiletype)
	signature := adapter.adapter.GetSignatureForFile(filename, systemFiletype)
	log.Printf("calculated file signature: %s", signature)
	return signature
}

type FilesAdapter struct{}

func (adapter FilesAdapter) GetSignatureForFile(filename string, systemFiletype files.SystemFiletype) string {
	file_hmac := hmac.New(sha256.New, []byte(Settings.FILES_SIGNATURE_KEY))
	file_hmac.Write([]byte(fmt.Sprintf("%s:%s", filename, systemFiletype.String())))
	hashsum := file_hmac.Sum(nil)
	hexdigest := make([]byte, hex.EncodedLen(len(hashsum)))
	hex.Encode(hexdigest, hashsum)
	return string(hexdigest)
}

func NewFilesAdapter() files.FilesPort {
	return FilesLoggingAdapter{adapter: FilesAdapter{}}
}
