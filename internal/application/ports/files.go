package ports

import (
	"chats-service/internal/domain/constants"
)

type FilesRepositoryPort interface {
	GetSignatureForFile(filename string, systemFiletype constants.SystemFiletype) string
}
