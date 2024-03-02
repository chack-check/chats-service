package factories

import (
	"github.com/chack-check/chats-service/api/v1/models"
	pb "github.com/chack-check/chats-service/protochats"
)

type MessagesFactory struct{}

func (factory *MessagesFactory) SchemaToProto(message *models.Message, chat *models.Chat) *pb.MessageResponse {
	var attachments []string

	for _, attachment := range message.Attachments {
		attachments = append(attachments, attachment)
	}

	return &pb.MessageResponse{
		Id:            int32(message.ID),
		SenderId:      int32(message.SenderId),
		ChatId:        int32(chat.ID),
		ChatTitle:     chat.Title,
		ChatAvatarUrl: chat.AvatarURL,
		ChatType:      chat.Type,
		Content:       message.Content,
		VoiceUrl:      message.VoiceURL,
		CircleUrl:     message.CircleURL,
		Attachments:   attachments,
		ReadedBy:      message.ReadedBy,
	}
}
