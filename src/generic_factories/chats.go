package generic_factories

import (
	"github.com/chack-check/chats-service/api/v1/dtos"
	pb "github.com/chack-check/chats-service/protochats"
)

type ChatsFactory struct{}

func (factory *ChatsFactory) SchemaToProto(schema *dtos.ChatDto) *pb.ChatResponse {
	var avatar_url string

	if schema.Avatar.ConvertedUrl != nil {
		avatar_url = *schema.Avatar.ConvertedUrl
	} else {
		avatar_url = schema.Avatar.OriginalUrl
	}

	return &pb.ChatResponse{
		Id:        int32(schema.Id),
		AvatarUrl: avatar_url,
		Title:     schema.Title,
		Type:      schema.Type,
	}
}
