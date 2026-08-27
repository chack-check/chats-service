// TODO: Make file for every controller
// TODO: Rename controllers to use cases
package controllers

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/dtos"
	"chats-service/internal/domain/entities"
	"chats-service/internal/domain/services"
)

// TODO: To constants or configuration
const SavedMessagesChatAvatarURL string = "https://storage.yandexcloud.net/diffaction/saved-messages-logo.svg"

// TODO: To errors package
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
	ErrChatWithSelf            = fmt.Errorf("you can't create chat with self user")
)

func setupSavedMessagesChatAvatar(chat *entities.Chat) {
	if chat == nil || chat.GetType() != constants.SavedMessagesChatType {
		return
	}

	chat.SetAvatar(entities.NewSavedFile(
		SavedMessagesChatAvatarURL,
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

type CreateChatController struct {
	chatsRepository      ports.ChatsRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
	usersRepository      ports.UsersRepositoryPort
	filesRepository      ports.FilesRepositoryPort
}

func NewCreateChatController(
	chatsRepository ports.ChatsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	filesRepository ports.FilesRepositoryPort,
) *CreateChatController {
	return &CreateChatController{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
		usersRepository:      usersRepository,
		filesRepository:      filesRepository,
	}
}

func (controller *CreateChatController) Execute(data dtos.CreateChatData, currentUserID int) (*entities.Chat, error) {
	if err := services.ValidateUploadingFile(controller.filesRepository, data.GetAvatar(), constants.AvatarFiletype, false); err != nil {
		return nil, err
	}

	currentUser, err := controller.usersRepository.GetById(currentUserID)
	if err != nil {
		return nil, ErrFindingUser
	}

	var savedChat *entities.Chat
	var savingError error
	switch data.GetType() {
	case constants.GroupChatType:
		savedChat, savingError = controller.createGroupChat(data, currentUser)
	case constants.UserChatType:
		savedChat, savingError = controller.createUserChat(data, currentUser)
	default:
		savingError = ErrInvalidCreatingChatType
	}

	if savingError != nil {
		return nil, savingError
	}

	controller.chatEventsRepository.SendChatCreated(*savedChat)
	return savedChat, nil
}

func (controller *CreateChatController) createGroupChat(data dtos.CreateChatData, currentUser *entities.User) (*entities.Chat, error) {
	chat := services.CreateChatDataToChat(data, 0)
	chat.SetOwnerID(currentUser.GetID())
	if !ValidateUserChatMember(chat, currentUser.GetID()) {
		newMembers := chat.GetMembers()
		newMembers = append(newMembers, currentUser.GetID())
		chat.SetMembers(newMembers)
	}
	if !ValidateUserChatAdmin(chat, currentUser.GetID()) {
		newAdmins := chat.GetAdmins()
		newAdmins = append(newAdmins, currentUser.GetID())
		chat.SetAdmins(newAdmins)
	}

	chat.SetType("group")
	savedChat, err := controller.chatsRepository.Save(chat)
	if err != nil {
		return nil, errors.Join(ErrSavingChat, err)
	}

	return savedChat, nil
}

func (controller *CreateChatController) createUserChat(data dtos.CreateChatData, currentUser *entities.User) (*entities.Chat, error) {
	if data.GetUserID() == nil {
		return nil, ErrCreatingNotUserChat
	}
	if *data.GetUserID() == currentUser.GetID() {
		return nil, ErrChatWithSelf
	}

	chatUser, err := controller.usersRepository.GetById(*data.GetUserID())
	if err != nil {
		return nil, ErrFindingUser
	}

	chat := services.CreateChatDataToChat(data, currentUser.GetID())
	if controller.chatsRepository.HasDeletedUserChat(chat) {
		chat, err := controller.chatsRepository.RestoreChat(chat)
		if err != nil {
			return nil, errors.Join(ErrRestoringChat, err)
		}

		return chat, nil
	}

	if controller.chatsRepository.CheckChatExists(chat) {
		return nil, ErrChatAlreadyExists
	}

	savedChat, err := controller.chatsRepository.Save(chat)
	if err != nil {
		return nil, errors.Join(ErrSavingChat, err)
	}

	savedChat.SetupUserData(chatUser)
	return savedChat, nil
}

type CreateSavedMessagesChatController struct {
	chatsRepository ports.ChatsRepositoryPort
}

func NewCreateSavedMessagesChatController(chatsRepository ports.ChatsRepositoryPort) *CreateSavedMessagesChatController {
	return &CreateSavedMessagesChatController{
		chatsRepository: chatsRepository,
	}
}

func (controller *CreateSavedMessagesChatController) Execute(data dtos.CreateChatData, currentUserID int) (*entities.Chat, error) {
	chat := services.CreateChatDataToChat(data, currentUserID)
	chat.SetOwnerID(currentUserID)
	chat.SetMembers([]int{currentUserID})
	chat.SetTitle("Saved messages")
	savedChat, err := controller.chatsRepository.Save(chat)
	if err != nil {
		return nil, errors.Join(ErrSavingChat, err)
	}

	setupSavedMessagesChatAvatar(savedChat)
	return savedChat, nil
}

type DeleteChatController struct {
	chatsRepository      ports.ChatsRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewDeleteChatController(
	chatsRepository ports.ChatsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *DeleteChatController {
	return &DeleteChatController{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *DeleteChatController) Execute(chatID, userID int) error {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return ErrChatNotFound
	}

	controller.chatsRepository.Delete(*chat)
	controller.chatEventsRepository.SendChatDeleted(*chat)
	return nil
}

type GetChatsController struct {
	chatsRepository       ports.ChatsRepositoryPort
	usersRepository       ports.UsersRepositoryPort
	userActionsRepository ports.UserActionsRepositoryPort
}

func NewGetChatsController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	userActionsRepository ports.UserActionsRepositoryPort,
) *GetChatsController {
	return &GetChatsController{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
	}
}

func (controller *GetChatsController) Execute(userID int, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	paginatedChats := controller.chatsRepository.GetUserAll(userID, page, perPage)
	fetchingUsers := GetUserChatsUsersIds(paginatedChats.GetData(), userID)
	fetchedUsers := controller.usersRepository.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(paginatedChats.GetData(), fetchedUsers, userID)
	completeChats := make([]entities.Chat, 0, len(chatsWithUsersData))
	for _, chat := range chatsWithUsersData {
		setupSavedMessagesChatAvatar(&chat)
		chatActions := controller.userActionsRepository.GetAllChatActionsUsers(chat)
		chat.SetupActions(chatActions)
		completeChats = append(completeChats, chat)
	}

	paginatedChats.SetData(completeChats)
	return paginatedChats
}

type GetChatsByIdsController struct {
	chatsRepository       ports.ChatsRepositoryPort
	usersRepository       ports.UsersRepositoryPort
	userActionsRepository ports.UserActionsRepositoryPort
}

func NewGetChatsByIdsController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	userActionsRepository ports.UserActionsRepositoryPort,
) *GetChatsByIdsController {
	return &GetChatsByIdsController{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
	}
}

func (controller *GetChatsByIdsController) Execute(chatIds []int, userID int) []entities.Chat {
	requestChats := controller.chatsRepository.GetByIdsForUser(chatIds, userID)
	fetchingUsers := GetUserChatsUsersIds(requestChats, userID)
	fetchedUsers := controller.usersRepository.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(requestChats, fetchedUsers, userID)
	completeChats := make([]entities.Chat, 0, len(chatsWithUsersData))
	for _, chat := range chatsWithUsersData {
		setupSavedMessagesChatAvatar(&chat)
		chatActions := controller.userActionsRepository.GetAllChatActionsUsers(chat)
		chat.SetupActions(chatActions)
		completeChats = append(completeChats, chat)
	}

	return completeChats
}

type GetChatController struct {
	chatsRepository       ports.ChatsRepositoryPort
	usersRepository       ports.UsersRepositoryPort
	userActionsRepository ports.UserActionsRepositoryPort
}

func NewGetChatController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	userActionsRepository ports.UserActionsRepositoryPort,
) *GetChatController {
	return &GetChatController{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
	}
}

func (controller *GetChatController) Execute(userID int, chatID int) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if chat.GetType() != "user" {
		setupSavedMessagesChatAvatar(chat)
		return chat, nil
	}

	anotherUserID := getAnotherUserIDForUserChat(*chat, userID)
	if anotherUserID == 0 {
		return chat, nil
	}

	anotherUser, err := controller.usersRepository.GetById(anotherUserID)
	if err != nil {
		return chat, nil
	}

	chatActions := controller.userActionsRepository.GetAllChatActionsUsers(*chat)
	chat.SetupActions(chatActions)
	chat.SetupUserData(anotherUser)
	return chat, nil
}

type UserActionController struct {
	chatsRepository       ports.ChatsRepositoryPort
	usersRepository       ports.UsersRepositoryPort
	userActionsRepository ports.UserActionsRepositoryPort
	chatEventsRepository  ports.ChatEventsRepositoryPort
}

func NewUserActionController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	userActionsRepository ports.UserActionsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *UserActionController {
	return &UserActionController{
		chatsRepository:       chatsRepository,
		usersRepository:       usersRepository,
		userActionsRepository: userActionsRepository,
		chatEventsRepository:  chatEventsRepository,
	}
}

func (controller *UserActionController) Execute(chatID int, userID int, actionType constants.ActionTypes) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	user, err := controller.usersRepository.GetById(userID)
	if err != nil {
		return nil, ErrFindingUser
	}

	newChatActions := controller.userActionsRepository.AddChatActionUser(*chat, *user, actionType)
	chat.SetupActions(newChatActions)
	controller.chatEventsRepository.SendChatUserAction(*chat)
	return chat, nil
}

type StopUserActionController struct {
	chatsRepository       ports.ChatsRepositoryPort
	userActionsRepository ports.UserActionsRepositoryPort
	chatEventsRepository  ports.ChatEventsRepositoryPort
}

func NewStopUserActionController(
	chatsRepository ports.ChatsRepositoryPort,
	userActionsRepository ports.UserActionsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *StopUserActionController {
	return &StopUserActionController{
		chatsRepository:       chatsRepository,
		userActionsRepository: userActionsRepository,
		chatEventsRepository:  chatEventsRepository,
	}
}

func (controller *StopUserActionController) Execute(chatID int, userID int, actionType constants.ActionTypes) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	newChatActions := controller.userActionsRepository.RemoveChatActionUser(*chat, userID, actionType)
	chat.SetupActions(newChatActions)
	controller.chatEventsRepository.SendChatUserAction(*chat)
	return chat, nil
}

type AddChatMembersController struct {
	chatsRepository      ports.ChatsRepositoryPort
	usersRepository      ports.UsersRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewAddChatMembersController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *AddChatMembersController {
	return &AddChatMembersController{
		chatsRepository:      chatsRepository,
		usersRepository:      usersRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *AddChatMembersController) Execute(chatID int, userID int, members []int) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userID) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != constants.GroupChatType {
		return nil, ErrChatNotGroup
	}

	newMembers := chat.GetMembers()
	users := controller.usersRepository.GetByIds(members)
	for _, member := range users {
		if !slices.Contains(newMembers, member.GetID()) {
			newMembers = append(newMembers, member.GetID())
		}
	}

	chat.SetMembers(newMembers)
	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type AddChatAdminsController struct {
	chatsRepository      ports.ChatsRepositoryPort
	usersRepository      ports.UsersRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewAddChatAdminsController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *AddChatAdminsController {
	return &AddChatAdminsController{
		chatsRepository:      chatsRepository,
		usersRepository:      usersRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *AddChatAdminsController) Execute(chatID int, userID int, admins []int) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userID) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != constants.GroupChatType {
		return nil, ErrChatNotGroup
	}

	newAdmins := chat.GetAdmins()
	users := controller.usersRepository.GetByIds(admins)
	for _, admin := range users {
		if !slices.Contains(newAdmins, admin.GetID()) {
			newAdmins = append(newAdmins, admin.GetID())
		}
	}

	chat.SetAdmins(newAdmins)
	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type RemoveChatMembersController struct {
	chatsRepository      ports.ChatsRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewRemoveChatMembersController(
	chatsRepository ports.ChatsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *RemoveChatMembersController {
	return &RemoveChatMembersController{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *RemoveChatMembersController) Execute(chatID int, userID int, members []int) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userID) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != constants.GroupChatType {
		return nil, ErrChatNotGroup
	}

	var newMembers []int
	for _, member := range chat.GetMembers() {
		if !slices.Contains(members, member) || member == userID {
			newMembers = append(newMembers, member)
		}
	}

	chat.SetMembers(newMembers)
	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type RemoveChatAdminsController struct {
	chatsRepository      ports.ChatsRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewRemoveChatAdminsController(
	chatsRepository ports.ChatsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *RemoveChatAdminsController {
	return &RemoveChatAdminsController{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *RemoveChatAdminsController) Execute(chatID int, userID int, admins []int) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if !ValidateUserChatAdmin(*chat, userID) {
		return nil, ErrChatNotAdmin
	}
	if chat.GetType() != constants.GroupChatType {
		return nil, ErrChatNotGroup
	}

	var newAdmins []int
	for _, admin := range chat.GetAdmins() {
		if !slices.Contains(admins, admin) || admin == userID {
			newAdmins = append(newAdmins, admin)
		}
	}

	chat.SetAdmins(newAdmins)
	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type QuitChatController struct {
	chatsRepository      ports.ChatsRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewQuitChatController(
	chatsRepository ports.ChatsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *QuitChatController {
	return &QuitChatController{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *QuitChatController) Execute(chatID int, userID int) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	var newMembers []int
	for _, member := range chat.GetMembers() {
		if member != userID {
			newMembers = append(newMembers, member)
		}
	}

	var newAdmins []int
	for _, admin := range chat.GetAdmins() {
		if admin != userID {
			newAdmins = append(newAdmins, admin)
		}
	}

	chat.SetMembers(newMembers)
	chat.SetAdmins(newAdmins)
	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type ChangeGroupChatController struct {
	chatsRepository      ports.ChatsRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewChangeGroupChatController(
	chatsRepository ports.ChatsRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *ChangeGroupChatController {
	return &ChangeGroupChatController{
		chatsRepository:      chatsRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *ChangeGroupChatController) Execute(chatID int, userID int, chatData dtos.ChangeGroupChatData) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if userID != chat.GetOwnerID() && !ValidateUserChatAdmin(*chat, userID) {
		return nil, ErrChatNotAdmin
	}

	if chat.GetType() != constants.GroupChatType {
		return nil, ErrChatNotGroup
	}

	if chatData.GetTitle() != nil {
		chat.SetTitle(*chatData.GetTitle())
	} else {
		chat.SetTitle("")
	}

	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type UpdateGroupChatAvatarController struct {
	chatsRepository      ports.ChatsRepositoryPort
	filesRepository      ports.FilesRepositoryPort
	chatEventsRepository ports.ChatEventsRepositoryPort
}

func NewUpdateGroupChatAvatarController(
	chatsRepository ports.ChatsRepositoryPort,
	filesRepository ports.FilesRepositoryPort,
	chatEventsRepository ports.ChatEventsRepositoryPort,
) *UpdateGroupChatAvatarController {
	return &UpdateGroupChatAvatarController{
		chatsRepository:      chatsRepository,
		filesRepository:      filesRepository,
		chatEventsRepository: chatEventsRepository,
	}
}

func (controller *UpdateGroupChatAvatarController) Execute(chatID int, userID int, newAvatar dtos.UploadingFile) (*entities.Chat, error) {
	chat, err := controller.chatsRepository.GetByIdForUser(chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	if userID != chat.GetOwnerID() && !ValidateUserChatAdmin(*chat, userID) {
		return nil, ErrChatNotAdmin
	}

	if chat.GetType() != constants.GroupChatType {
		return nil, ErrChatNotGroup
	}

	err = services.ValidateUploadingFile(controller.filesRepository, &newAvatar, constants.AvatarFiletype, true)
	if err != nil {
		return nil, err
	}

	savedFile := services.UploadingFileToSavedFile(newAvatar)
	chat.SetAvatar(savedFile)
	savedChat, err := controller.chatsRepository.Save(*chat)
	if err != nil {
		return nil, ErrSavingChat
	}

	controller.chatEventsRepository.SendChatChanged(*savedChat)
	return savedChat, nil
}

type SearchChatsController struct {
	chatsRepository ports.ChatsRepositoryPort
	usersRepository ports.UsersRepositoryPort
}

func NewSearchChatsController(
	chatsRepository ports.ChatsRepositoryPort,
	usersRepository ports.UsersRepositoryPort,
) *SearchChatsController {
	return &SearchChatsController{
		chatsRepository: chatsRepository,
		usersRepository: usersRepository,
	}
}

func (controller *SearchChatsController) Execute(userID int, query string, page int, perPage int) dtos.PaginatedResponse[entities.Chat] {
	requestChats := controller.chatsRepository.SearchChats(userID, query, page, perPage)

	fetchingUsers := GetUserChatsUsersIds(requestChats.GetData(), userID)
	fetchedUsers := controller.usersRepository.GetByIds(fetchingUsers)
	chatsWithUsersData := SetupUserChatsData(requestChats.GetData(), fetchedUsers, userID)

	var resultChats []entities.Chat

	for _, chat := range chatsWithUsersData {
		if chat.GetType() == constants.SavedMessagesChatType {
			setupSavedMessagesChatAvatar(&chat)
		}

		if chat.GetType() == constants.UserChatType && !strings.Contains(strings.ToLower(chat.GetTitle()), strings.ToLower(query)) {
			continue
		}

		resultChats = append(resultChats, chat)
	}

	requestChats.SetData(resultChats)
	return requestChats
}
