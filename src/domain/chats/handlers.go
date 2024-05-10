package chats

import (
	"errors"
	"fmt"
	"slices"

	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/users"
	"github.com/chack-check/chats-service/domain/utils"
)

var (
	ErrFindingUser             = fmt.Errorf("error finding user")
	ErrCreatingNotUserChat     = fmt.Errorf("trying to create user chat with not specified user id")
	ErrSavingChat              = fmt.Errorf("error saving chat")
	ErrRestoringChat           = fmt.Errorf("error restoring chat")
	ErrChatAlreadyExists       = fmt.Errorf("you already have chat with this user")
	ErrChatNotFound            = fmt.Errorf("there is no such chat")
	ErrNotGroupAdmin           = fmt.Errorf("you are not a group chat admin")
	ErrChatNotGroup            = fmt.Errorf("the editing chat is not group")
	ErrInvalidCreatingChatType = fmt.Errorf("invalid creating chat type. Valid values: group, user, saved_messages")
	ErrChatNotAdmin            = fmt.Errorf("user is not admin in chat")
)

func ValidateUserChatMember(chat Chat, userId int) bool {
	return slices.Contains(chat.GetMembers(), userId)
}

func ValidateUserChatAdmin(chat Chat, userId int) bool {
	return chat.GetOwnerId() == userId || slices.Contains(chat.GetAdmins(), userId)
}

func GetAnotherUserIdForUserChat(chat Chat, currentUserId int) int {
	if chat.GetType() != "user" {
		return 0
	}

	var anotherUserId int
	for _, member := range chat.GetMembers() {
		if member != currentUserId {
			anotherUserId = member
		}
	}

	return anotherUserId
}

func GetUserChatsUsersIds(chats []Chat, currentUserId int) []int {
	var fetchingUsers []int
	for _, chat := range chats {
		if chat.GetType() != "user" {
			continue
		}

		anotherUserId := GetAnotherUserIdForUserChat(chat, currentUserId)
		if anotherUserId == 0 {
			continue
		}

		fetchingUsers = append(fetchingUsers, anotherUserId)
	}

	return fetchingUsers
}

func SetupUserChatsData(chats []Chat, fetchedUsers []users.User, currentUserId int) []Chat {
	var newChats []Chat
	for _, chat := range chats {
		if chat.GetType() != "user" {
			newChats = append(newChats, chat)
			continue
		}

		anotherUserId := GetAnotherUserIdForUserChat(chat, currentUserId)
		if anotherUserId == 0 {
			newChats = append(newChats, chat)
			continue
		}

		var chatUser users.User
		for _, user := range fetchedUsers {
			if user.GetId() == anotherUserId {
				chatUser = user
			}
		}

		chat.SetupUserData(&chatUser)
		newChats = append(newChats, chat)
	}

	return newChats
}

type CreateChatHandler struct {
	chatsPort      ChatsPort
	chatEventsPort ChatEventsPort
	usersPort      users.UsersPort
	filesPort      files.FilesPort
}

func (handler *CreateChatHandler) createGroupChat(data CreateChatData, currentUser *users.User) (*Chat, error) {
	chat := CreateChatDataToChat(data, 0)
	chat.SetOwnerId(currentUser.GetId())
	if !ValidateUserChatMember(chat, currentUser.GetId()) {
		newMembers := chat.GetMembers()
		newMembers = append(newMembers, currentUser.GetId())
		chat.SetMembers(newMembers)
	}
	if !ValidateUserChatAdmin(chat, currentUser.GetId()) {
		newAdmins := chat.GetAdmins()
		newAdmins = append(newAdmins, currentUser.GetId())
		chat.SetAdmins(newAdmins)
	}

	chat.SetType("group")
	savedChat, err := handler.chatsPort.Save(chat)
	if err != nil {
		return nil, errors.Join(ErrSavingChat, err)
	}

	return savedChat, nil
}

func (handler *CreateChatHandler) createUserChat(data CreateChatData, currentUser *users.User) (*Chat, error) {
	if data.userId == nil {
		return nil, ErrCreatingNotUserChat
	}

	chatUser, err := handler.usersPort.GetById(*data.userId)
	if err != nil {
		return nil, ErrFindingUser
	}

	chat := CreateChatDataToChat(data, currentUser.GetId())
	if handler.chatsPort.HasDeletedUserChat(chat) {
		chat, err := handler.chatsPort.RestoreChat(chat)
		if err != nil {
			return nil, errors.Join(ErrRestoringChat, err)
		}

		return chat, nil
	}

	if handler.chatsPort.CheckChatExists(chat) {
		return nil, ErrChatAlreadyExists
	}

	savedChat, err := handler.chatsPort.Save(chat)
	if err != nil {
		return nil, errors.Join(ErrSavingChat, err)
	}

	savedChat.SetupUserData(chatUser)
	return savedChat, nil
}

func (handler *CreateChatHandler) Execute(data CreateChatData, currentUserId int) (*Chat, error) {
	if err := files.ValidateUploadingFile(handler.filesPort, data.avatar, files.AvatarFiletype, false); err != nil {
		return nil, err
	}

	currentUser, err := handler.usersPort.GetById(currentUserId)
	if err != nil {
		return nil, ErrFindingUser
	}

	var savedChat *Chat
	var savingError error
	switch data.GetType() {
	case GroupChatType:
		savedChat, savingError = handler.createGroupChat(data, currentUser)
	case UserChatType:
		savedChat, savingError = handler.createUserChat(data, currentUser)
	default:
		savingError = ErrInvalidCreatingChatType
	}

	if savingError != nil {
		return nil, savingError
	}

	handler.chatEventsPort.SendChatCreated(*savedChat)
	return savedChat, nil
}

type CreateSavedMessagesChat struct {
	chatsPort ChatsPort
}

func (handler *CreateSavedMessagesChat) Execute(data CreateChatData, currentUserId int) (*Chat, error) {
	chat := CreateChatDataToChat(data, currentUserId)
	chat.SetOwnerId(currentUserId)
	chat.SetMembers([]int{currentUserId})
	chat.SetTitle("Saved messages")
	savedChat, err := handler.chatsPort.Save(chat)
	if err != nil {
		return nil, errors.Join(ErrSavingChat, err)
	}

	return savedChat, nil
}

type DeleteChatHandler struct {
	chatsPort      ChatsPort
	chatEventsPort ChatEventsPort
}

func (handler *DeleteChatHandler) Execute(chatId, userId int) error {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return ErrChatNotFound
	}

	handler.chatsPort.Delete(*chat)
	handler.chatEventsPort.SendChatDeleted(*chat)
	return nil
}

type GetChatsHandler struct {
	chatsPort       ChatsPort
	usersPort       users.UsersPort
	userActionsPort UserActionsPort
}

func (handler *GetChatsHandler) Execute(userId int, page int, perPage int) utils.PaginatedResponse[Chat] {
	paginatedChats := handler.chatsPort.GetUserAll(userId, page, perPage)
	fetchingUsers := GetUserChatsUsersIds(paginatedChats.GetData(), userId)
	fetchedUsers := handler.usersPort.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(paginatedChats.GetData(), fetchedUsers, userId)
	var completeChats []Chat
	for _, chat := range chatsWithUsersData {
		chatActions := handler.userActionsPort.GetAllChatActionsUsers(chat)
		chat.SetupActions(chatActions)
		completeChats = append(completeChats, chat)
	}

	paginatedChats.SetData(completeChats)
	return paginatedChats
}

type GetChatsByIdsHandler struct {
	chatsPort       ChatsPort
	usersPort       users.UsersPort
	userActionsPort UserActionsPort
}

func (handler *GetChatsByIdsHandler) Execute(chatIds []int, userId int) []Chat {
	chats := handler.chatsPort.GetByIdsForUser(chatIds, userId)
	fetchingUsers := GetUserChatsUsersIds(chats, userId)
	fetchedUsers := handler.usersPort.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(chats, fetchedUsers, userId)
	var completeChats []Chat
	for _, chat := range chatsWithUsersData {
		chatActions := handler.userActionsPort.GetAllChatActionsUsers(chat)
		chat.SetupActions(chatActions)
		completeChats = append(completeChats, chat)
	}

	return completeChats
}

type GetChatHandler struct {
	chatsPort       ChatsPort
	usersPort       users.UsersPort
	userActionsPort UserActionsPort
}

func (handler *GetChatHandler) Execute(userId int, chatId int) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if chat.GetType() != "user" {
		return chat, nil
	}

	anotherUserId := GetAnotherUserIdForUserChat(*chat, userId)
	if anotherUserId == 0 {
		return chat, nil
	}

	anotherUser, err := handler.usersPort.GetById(anotherUserId)
	if err != nil {
		return chat, nil
	}

	chatActions := handler.userActionsPort.GetAllChatActionsUsers(*chat)
	chat.SetupActions(chatActions)
	chat.SetupUserData(anotherUser)
	return chat, nil
}

type UserActionHandler struct {
	chatsPort       ChatsPort
	usersPort       users.UsersPort
	userActionsPort UserActionsPort
	chatEventsPort  ChatEventsPort
}

func (handler *UserActionHandler) Execute(chatId int, userId int, actionType ActionTypes) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	user, err := handler.usersPort.GetById(userId)
	if err != nil {
		return nil, ErrFindingUser
	}

	newChatActions := handler.userActionsPort.AddChatActionUser(*chat, *user, actionType)
	chat.SetupActions(newChatActions)
	handler.chatEventsPort.SendChatUserAction(*chat)
	return chat, nil
}

type StopUserActionHandler struct {
	chatsPort       ChatsPort
	userActionsPort UserActionsPort
	chatEventsPort  ChatEventsPort
}

func (handler *StopUserActionHandler) Execute(chatId int, userId int, actionType ActionTypes) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	newChatActions := handler.userActionsPort.RemoveChatActionUser(*chat, userId, actionType)
	chat.SetupActions(newChatActions)
	handler.chatEventsPort.SendChatUserAction(*chat)
	return chat, nil
}

type AddChatMembersHandler struct {
	chatsPort      ChatsPort
	usersPort      users.UsersPort
	chatEventsPort ChatEventsPort
}

func (handler *AddChatMembersHandler) Execute(chatId int, userId int, members []int) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userId) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != GroupChatType {
		return nil, ErrChatNotGroup
	}

	newMembers := chat.GetMembers()
	users := handler.usersPort.GetByIds(members)
	for _, member := range users {
		if !slices.Contains(newMembers, member.GetId()) {
			newMembers = append(newMembers, member.GetId())
		}
	}

	chat.SetMembers(newMembers)
	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}

type AddChatAdminsHandler struct {
	chatsPort      ChatsPort
	usersPort      users.UsersPort
	chatEventsPort ChatEventsPort
}

func (handler *AddChatAdminsHandler) Execute(chatId int, userId int, admins []int) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userId) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != GroupChatType {
		return nil, ErrChatNotGroup
	}

	newAdmins := chat.GetAdmins()
	users := handler.usersPort.GetByIds(admins)
	for _, admin := range users {
		if !slices.Contains(newAdmins, admin.GetId()) {
			newAdmins = append(newAdmins, admin.GetId())
		}
	}

	chat.SetAdmins(newAdmins)
	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}

type RemoveChatMembersHandler struct {
	chatsPort      ChatsPort
	chatEventsPort ChatEventsPort
}

func (handler *RemoveChatMembersHandler) Execute(chatId int, userId int, members []int) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userId) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != GroupChatType {
		return nil, ErrChatNotGroup
	}

	var newMembers []int
	for _, member := range chat.GetMembers() {
		if !slices.Contains(members, member) || member == userId {
			newMembers = append(newMembers, member)
		}
	}

	chat.SetMembers(newMembers)
	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}

type RemoveChatAdminsHandler struct {
	chatsPort      ChatsPort
	chatEventsPort ChatEventsPort
}

func (handler *RemoveChatAdminsHandler) Execute(chatId int, userId int, admins []int) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userId) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != GroupChatType {
		return nil, ErrChatNotGroup
	}

	var newAdmins []int
	for _, admin := range chat.GetAdmins() {
		if !slices.Contains(admins, admin) || admin == userId {
			newAdmins = append(newAdmins, admin)
		}
	}

	chat.SetAdmins(newAdmins)
	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}

type QuitChatHandler struct {
	chatsPort      ChatsPort
	chatEventsPort ChatEventsPort
}

func (handler *QuitChatHandler) Execute(chatId int, userId int) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	var newMembers []int
	for _, member := range chat.GetMembers() {
		if member != userId {
			newMembers = append(newMembers, member)
		}
	}

	var newAdmins []int
	for _, admin := range chat.GetAdmins() {
		if admin != userId {
			newAdmins = append(newAdmins, admin)
		}
	}

	chat.SetMembers(newMembers)
	chat.SetAdmins(newAdmins)
	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}

type ChangeGroupChatHandler struct {
	chatsPort      ChatsPort
	chatEventsPort ChatEventsPort
}

func (handler *ChangeGroupChatHandler) Execute(chatId int, userId int, chatData ChangeGroupChatData) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if userId != chat.GetOwnerId() && !ValidateUserChatAdmin(*chat, userId) {
		return nil, ErrChatNotAdmin
	}

	if chat.GetType() != GroupChatType {
		return nil, ErrChatNotGroup
	}

	if chatData.GetTitle() != nil {
		chat.SetTitle(*chatData.GetTitle())
	} else {
		chat.SetTitle("")
	}

	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}

type UpdateGroupChatAvatar struct {
	chatsPort      ChatsPort
	filesPort      files.FilesPort
	chatEventsPort ChatEventsPort
}

func (handler *UpdateGroupChatAvatar) Execute(chatId int, userId int, newAvatar files.UploadingFile) (*Chat, error) {
	chat, err := handler.chatsPort.GetByIdForUser(chatId, userId)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if userId != chat.GetOwnerId() && !ValidateUserChatAdmin(*chat, userId) {
		return nil, ErrChatNotAdmin
	}

	if chat.GetType() != GroupChatType {
		return nil, ErrChatNotGroup
	}

	err = files.ValidateUploadingFile(handler.filesPort, &newAvatar, files.AvatarFiletype, true)
	if err != nil {
		return nil, err
	}

	savedFile := files.UploadingFileToSavedFile(newAvatar)
	chat.SetAvatar(savedFile)
	savedChat, err := handler.chatsPort.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	handler.chatEventsPort.SendChatChanged(*savedChat)
	return savedChat, nil
}
