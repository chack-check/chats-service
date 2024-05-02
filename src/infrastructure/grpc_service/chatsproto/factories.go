package chatsproto

import (
	"time"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/domain/files"
	"github.com/chack-check/chats-service/domain/messages"
	"github.com/chack-check/chats-service/domain/utils"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/chatsproto/chatsprotobuf"
)

func SavedFileToProto(file files.SavedFile) *chatsprotobuf.SavedFile {
	return &chatsprotobuf.SavedFile{
		OriginalUrl:       file.GetOriginalUrl(),
		OriginalFilename:  file.GetOriginalFilename(),
		ConvertedUrl:      file.GetConvertedUrl(),
		ConvertedFilename: file.GetConvertedFilename(),
	}
}

func ChatModelToProto(chat chats.Chat) *chatsprotobuf.ChatResponse {
	var avatar *chatsprotobuf.SavedFile
	if file := chat.GetAvatar(); file != nil {
		avatar = SavedFileToProto(*file)
	}

	var members []int32
	for _, member := range chat.GetMembers() {
		members = append(members, int32(member))
	}

	var admins []int32
	for _, admin := range chat.GetAdmins() {
		admins = append(admins, int32(admin))
	}

	return &chatsprotobuf.ChatResponse{
		Id:         int32(chat.GetId()),
		Avatar:     avatar,
		Title:      chat.GetTitle(),
		Type:       string(chat.GetType()),
		MembersIds: members,
		IsArchived: chat.GetIsArchived(),
		OwnerId:    int32(chat.GetOwnerId()),
		AdminsIds:  admins,
	}
}

func ReactionToProto(reaction messages.MessageReaction) *chatsprotobuf.MessageReaction {
	return &chatsprotobuf.MessageReaction{
		UserId:  int32(reaction.GetUserId()),
		Content: reaction.GetContent(),
	}
}

func MessageToProto(message messages.Message) *chatsprotobuf.MessageResponse {
	var voice *chatsprotobuf.SavedFile
	if file := message.GetVoice(); file != nil {
		voice = SavedFileToProto(*file)
	}

	var circle *chatsprotobuf.SavedFile
	if file := message.GetCircle(); file != nil {
		circle = SavedFileToProto(*file)
	}

	var attachments []*chatsprotobuf.SavedFile
	for _, attachment := range message.GetAttachments() {
		attachments = append(attachments, SavedFileToProto(attachment))
	}

	var replyToId *int32
	if message.GetReplyToId() != nil {
		replyToId32 := int32(*message.GetReplyToId())
		replyToId = &replyToId32
	}

	var mentioned []int32
	for _, user := range message.GetMentioned() {
		mentioned = append(mentioned, int32(user))
	}

	var readedBy []int32
	for _, user := range message.GetReadedBy() {
		readedBy = append(readedBy, int32(user))
	}

	var reactions []*chatsprotobuf.MessageReaction
	for _, reaction := range message.GetReactions() {
		reactions = append(reactions, ReactionToProto(reaction))
	}

	var createdAt *string
	if dt := message.GetCreatedAt(); dt != nil {
		isodt := dt.Format(time.RFC3339)
		createdAt = &isodt
	}

	chat := message.GetChat()
	return &chatsprotobuf.MessageResponse{
		Id:          int32(message.GetId()),
		SenderId:    int32(message.GetSenderId()),
		ChatId:      int32(chat.GetId()),
		Type:        string(message.GetType()),
		Content:     message.GetContent(),
		Voice:       voice,
		Circle:      circle,
		Attachments: attachments,
		ReplyToId:   replyToId,
		Mentioned:   mentioned,
		ReadedBy:    readedBy,
		Reactions:   reactions,
		CreatedAt:   createdAt,
	}
}

func OffsetMessagesToProto(offsetMessages utils.OffsetResponse[messages.Message]) *chatsprotobuf.PaginatedMessages {
	var data []*chatsprotobuf.MessageResponse
	for _, message := range offsetMessages.GetData() {
		data = append(data, MessageToProto(message))
	}

	return &chatsprotobuf.PaginatedMessages{
		Offset: int32(offsetMessages.GetOffset()),
		Limit:  int32(offsetMessages.GetLimit()),
		Total:  int32(offsetMessages.GetTotal()),
		Data:   data,
	}
}
