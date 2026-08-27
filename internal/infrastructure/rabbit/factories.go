package rabbit

import (
	"chats-service/internal/domain/entities"
)

func SavedFileToEventSavedFile(file entities.SavedFile) EventSavedFile {
	return EventSavedFile{
		OriginalUrl:       file.GetOriginalURL(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedUrl:      file.GetConvertedURL(),
		ConvertedFilename: file.GetConvertedFilename(),
	}
}

func ActionUserToEventActionUser(user entities.ActionUser) EventActionUser {
	return EventActionUser{
		Id:         user.GetID(),
		LastName:   user.GetLastName(),
		FirstName:  user.GetFirstName(),
		MiddleName: user.GetMiddleName(),
		Username:   user.GetUsername(),
	}
}

func ChatToChatEvent(chat entities.Chat) ChatEvent {
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
		Id:         chat.GetID(),
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       string(chat.GetType()),
		Members:    chat.GetMembers(),
		IsArchived: chat.GetIsArchived(),
		OwnerId:    chat.GetOwnerID(),
		Admins:     chat.GetAdmins(),
		Actions:    actions,
	}
}

func MessageReactionToEventReaction(reaction entities.MessageReaction) EventMessageReaction {
	return EventMessageReaction{
		UserId:  reaction.GetUserID(),
		Content: reaction.GetContent(),
	}
}

func MessageToMessageEvent(message entities.Message) MessageEvent {
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
		Id:          message.GetID(),
		SenderId:    message.GetSenderID(),
		ChatId:      chat.GetID(),
		Type:        string(message.GetType()),
		Content:     message.GetContent(),
		Voice:       eventVoice,
		Circle:      eventCircle,
		Attachments: attachments,
		ReplyToId:   message.GetReplyToID(),
		Mentioned:   message.GetMentioned(),
		ReadedBy:    message.GetReadedBy(),
		Reactions:   reactions,
		CreatedAt:   message.GetCreatedAt(),
	}
}
