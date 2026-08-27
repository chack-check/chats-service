package usersproto

import (
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/grpc_service/usersproto/usersprotobuf"
)

func ProtoSavedFileToModel(file *usersprotobuf.SavedFile) entities.SavedFile {
	return entities.NewSavedFile(
		file.OriginalUrl,
		file.OriginalFilename,
		file.ConvertedUrl,
		file.ConvertedFilename,
	)
}

func ProtoUserToModel(user *usersprotobuf.UserResponse) entities.User {
	var avatar *entities.SavedFile
	if user.Avatar != nil {
		file := ProtoSavedFileToModel(user.Avatar)
		avatar = &file
	}

	return entities.NewUser(
		int(user.Id),
		avatar,
		user.LastName,
		user.FirstName,
		user.MiddleName,
		user.Username,
	)
}
