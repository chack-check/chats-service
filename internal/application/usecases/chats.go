package usecases

import (
	"slices"

	"chats-service/configs"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"
)

func setupSavedMessagesChatAvatar(chat *entities.Chat) {
	configuration := configs.GetAPIConfiguration()
	if chat == nil || chat.GetType() != constants.SavedMessagesChatType {
		return
	}

	chat.SetAvatar(entities.NewSavedFile(
		configuration.SavedMessagesChatAvatarURL.String(),
		"saved-messages-logo.svg",
		nil,
		nil,
	))
}

func ValidateUserChatMember(chat entities.Chat, userID int) bool {
	return slices.Contains(chat.GetMembers(), userID)
}

func ValidateUserChatAdmin(chat entities.Chat, userID int) bool {
	return chat.GetOwnerID() == userID || slices.Contains(chat.GetAdmins(), userID)
}

func getAnotherUserIDForUserChat(chat entities.Chat, currentUserID int) int {
	if chat.GetType() != "user" {
		return 0
	}

	var anotherUserID int
	for _, member := range chat.GetMembers() {
		if member != currentUserID {
			anotherUserID = member
		}
	}

	return anotherUserID
}

func GetUserChatsUsersIds(chats []entities.Chat, currentUserID int) []int {
	var fetchingUsers []int
	for _, chat := range chats {
		if chat.GetType() != "user" {
			continue
		}

		anotherUserID := getAnotherUserIDForUserChat(chat, currentUserID)
		if anotherUserID == 0 {
			continue
		}

		fetchingUsers = append(fetchingUsers, anotherUserID)
	}

	return fetchingUsers
}

func SetupUserChatsData(requestChats []entities.Chat, fetchedUsers []entities.User, currentUserID int) []entities.Chat {
	var newChats []entities.Chat
	for _, chat := range requestChats {
		if chat.GetType() != "user" {
			newChats = append(newChats, chat)
			continue
		}

		anotherUserID := getAnotherUserIDForUserChat(chat, currentUserID)
		if anotherUserID == 0 {
			newChats = append(newChats, chat)
			continue
		}

		var chatUser entities.User
		for _, user := range fetchedUsers {
			if user.GetID() == anotherUserID {
				chatUser = user
			}
		}

		chat.SetupUserData(&chatUser)
		newChats = append(newChats, chat)
	}

	return newChats
}
