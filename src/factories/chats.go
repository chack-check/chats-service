package factories

import (
	"github.com/chack-check/chats-service/api/v1/models"
	pb "github.com/chack-check/chats-service/protochats"
)

type ChatsFactory struct{}

func (factory *ChatsFactory) SchemaToProto(schema *models.Chat) *pb.ChatResponse {
	return &pb.ChatResponse{
		Id:        int32(schema.ID),
		AvatarUrl: schema.AvatarURL,
		Title:     schema.Title,
		Type:      schema.Type,
	}
}
