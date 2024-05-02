package chats

import "github.com/chack-check/chats-service/domain/files"

func CreateChatDataToChat(data CreateChatData, currentUserId int) Chat {
	var avatar *files.SavedFile
	if uploadingFile := data.GetAvatar(); uploadingFile != nil {
		savedFile := files.UploadingFileToSavedFile(*uploadingFile)
		avatar = &savedFile
	} else {
		avatar = nil
	}

	var title string
	if data.title != nil {
		title = *data.GetTitle()
	} else {
		title = ""
	}

	if data.GetType() == "user" && currentUserId != 0 && data.userId != nil {
		data.membersIds = []int{currentUserId, *data.userId}
	}

	return NewChat(
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
