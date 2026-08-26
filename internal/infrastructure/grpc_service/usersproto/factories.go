package usersproto

import (
	"chats-service/internal/domain/files"
	"chats-service/internal/domain/users"
	"chats-service/internal/infrastructure/grpc_service/usersproto/usersprotobuf"
)

func ProtoSavedFileToModel(file *usersprotobuf.SavedFile) files.SavedFile {
	return files.NewSavedFile(
		file.OriginalUrl,
		file.OriginalFilename,
		file.ConvertedUrl,
		file.ConvertedFilename,
	)
}

func ProtoUserToModel(user *usersprotobuf.UserResponse) users.User {
	var avatar *files.SavedFile
	if user.Avatar != nil {
		file := ProtoSavedFileToModel(user.Avatar)
		avatar = &file
	}

	return users.NewUser(
		int(user.Id),
		avatar,
		user.LastName,
		user.FirstName,
		user.MiddleName,
		user.Username,
	)
}
