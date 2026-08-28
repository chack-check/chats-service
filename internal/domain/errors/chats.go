package errors

import "errors"

var (
	ErrFindingUser             = errors.New("error finding user")
	ErrCreatingNotUserChat     = errors.New("trying to create user chat with not specified user id")
	ErrSavingChat              = errors.New("error saving chat")
	ErrRestoringChat           = errors.New("error restoring chat")
	ErrChatAlreadyExists       = errors.New("you already have chat with this user")
	ErrChatNotFound            = errors.New("there is no such chat")
	ErrNotGroupAdmin           = errors.New("you are not a group chat admin")
	ErrChatNotGroup            = errors.New("the editing chat is not group")
	ErrInvalidCreatingChatType = errors.New("invalid creating chat type. Valid values: group, user, saved_messages")
	ErrChatNotAdmin            = errors.New("user is not admin in chat")
	ErrChatWithSelf            = errors.New("you can't create chat with self user")
)
