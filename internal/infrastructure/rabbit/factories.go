package rabbit

import (
	"chats-service/internal/domain/chats"
	"chats-service/internal/domain/files"
	"chats-service/internal/domain/messages"
	"chats-service/internal/domain/users"
)

func SavedFileToEventSavedFile(file files.SavedFile) EventSavedFile {
	return EventSavedFile{
		OriginalUrl:       file.GetOriginalUrl(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedUrl:      file.GetConvertedUrl(),
		ConvertedFilename: file.GetConvertedFilename(),
	}
}

func ActionUserToEventActionUser(user users.ActionUser) EventActionUser {
	return EventActionUser{
		Id:         user.GetId(),
		LastName:   user.GetLastName(),
		FirstName:  user.GetFirstName(),
		MiddleName: user.GetMiddleName(),
		Username:   user.GetUsername(),
	}
}

func ChatToChatEvent(chat chats.Chat) ChatEvent {
	var avatar *EventSavedFile
	if file := chat.GetAvatar(); file != nil {
		eventFile := SavedFileToEventSavedFile(*file)
		avatar = &eventFile
	}

	actions := make(map[string][]EventActionUser)
	for action, users := range chat.GetActions() {
		eventUsers := make([]EventActionUser, 0, len(users))
		for _, user := range users {
			eventUsers = append(eventUsers, ActionUserToEventActionUser(user))
		}

		actions[string(action)] = eventUsers
	}

	return ChatEvent{
		Id:         chat.GetId(),
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       string(chat.GetType()),
		Members:    chat.GetMembers(),
		IsArchived: chat.GetIsArchived(),
		OwnerId:    chat.GetOwnerId(),
		Admins:     chat.GetAdmins(),
		Actions:    actions,
	}
}

func MessageReactionToEventReaction(reaction messages.MessageReaction) EventMessageReaction {
	return EventMessageReaction{
		UserId:  reaction.GetUserId(),
		Content: reaction.GetContent(),
	}
}

func MessageToMessageEvent(message messages.Message) MessageEvent {
	chat := message.GetChat()
	var eventVoice *EventSavedFile
	if voice := message.GetVoice(); voice != nil {
		eventFile := SavedFileToEventSavedFile(*voice)
		eventVoice = &eventFile
	}

	var eventCircle *EventSavedFile
	if circle := message.GetCircle(); circle != nil {
		eventFile := SavedFileToEventSavedFile(*circle)
		eventCircle = &eventFile
	}

	attachments := make([]EventSavedFile, 0, len(message.GetAttachments()))
	for _, attachment := range message.GetAttachments() {
		eventFile := SavedFileToEventSavedFile(attachment)
		attachments = append(attachments, eventFile)
	}

	reactions := make([]EventMessageReaction, 0, len(message.GetReactions()))
	for _, reaction := range message.GetReactions() {
		eventReaction := MessageReactionToEventReaction(reaction)
		reactions = append(reactions, eventReaction)
	}

	return MessageEvent{
		Id:          message.GetId(),
		SenderId:    message.GetSenderId(),
		ChatId:      chat.GetId(),
		Type:        string(message.GetType()),
		Content:     message.GetContent(),
		Voice:       eventVoice,
		Circle:      eventCircle,
		Attachments: attachments,
		ReplyToId:   message.GetReplyToId(),
		Mentioned:   message.GetMentioned(),
		ReadedBy:    message.GetReadedBy(),
		Reactions:   reactions,
		CreatedAt:   message.GetCreatedAt(),
	}
}
