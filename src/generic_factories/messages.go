package generic_factories

import (
	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/models"
	pb "github.com/chack-check/chats-service/protochats"
)

type MessagesFactory struct{}

func (factory *MessagesFactory) SchemaToProto(message *models.Message, chat *dtos.ChatDto) *pb.MessageResponse {
	var attachments []string

	for _, attachment := range message.Attachments {
		var attachment_url string

		if attachment.ConvertedUrl != "" {
			attachment_url = attachment.ConvertedUrl
		} else {
			attachment_url = attachment.OriginalUrl
		}

		attachments = append(attachments, attachment_url)
	}

	var chat_avatar_url string
	if chat.Avatar.ConvertedUrl != nil {
		chat_avatar_url = *chat.Avatar.ConvertedUrl
	} else {
		chat_avatar_url = chat.Avatar.OriginalUrl
	}

	var voice_url string
	if message.Voice.ConvertedUrl != "" {
		voice_url = message.Voice.ConvertedUrl
	} else {
		voice_url = message.Voice.OriginalUrl
	}

	var circle_url string
	if message.Circle.ConvertedUrl != "" {
		circle_url = message.Circle.ConvertedUrl
	} else {
		circle_url = message.Circle.OriginalUrl
	}

	return &pb.MessageResponse{
		Id:            int32(message.ID),
		SenderId:      int32(message.SenderId),
		ChatId:        int32(chat.Id),
		ChatTitle:     chat.Title,
		ChatAvatarUrl: chat_avatar_url,
		ChatType:      chat.Type,
		Content:       message.Content,
		VoiceUrl:      voice_url,
		CircleUrl:     circle_url,
		Attachments:   attachments,
		ReadedBy:      message.ReadedBy,
	}
}
