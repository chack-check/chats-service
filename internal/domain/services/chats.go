package services

import (
	"chats-service/internal/domain/dtos"
	"chats-service/internal/domain/entities"
)

func CreateChatDataToChat(data dtos.CreateChatData, currentUserID int) entities.Chat {
	var avatar *entities.SavedFile
	if uploadingFile := data.GetAvatar(); uploadingFile != nil {
		savedFile := UploadingFileToSavedFile(*uploadingFile)
		avatar = &savedFile
	} else {
		avatar = nil
	}

	var title string
	if data.GetTitle() != nil {
		title = *data.GetTitle()
	} else {
		title = ""
	}

	if data.GetType() == "user" && currentUserID != 0 && data.GetUserID() != nil {
		data.SetMembersIds([]int{currentUserID, *data.GetUserID()})
	}

	return entities.NewChat(
		0,
		avatar,
		title,
		data.GetType(),
		data.GetMembersIds(),
		false,
		0,
		[]int{},
	)
}
