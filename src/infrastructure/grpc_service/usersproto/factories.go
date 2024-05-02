package usersproto

import (
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/users"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/usersproto/usersprotobuf"
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
