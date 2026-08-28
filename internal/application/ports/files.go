package ports

import (
	"chats-service/internal/domain/constants"
)

type FilesProvider interface {
	GetSignatureForFile(filename string, systemFiletype constants.SystemFiletype) string
}
