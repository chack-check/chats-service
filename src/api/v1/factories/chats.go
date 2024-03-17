package factories

import (
	"fmt"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/lib/pq"
)

func DbChatToChatDto(chat models.Chat, actions *[]dtos.ChatActionDto) dtos.ChatDto {
	var chat_members []int
	for _, member := range chat.Members {
		chat_members = append(chat_members, int(member))
	}

	var chat_admins []int
	for _, admin := range chat.Admins {
		chat_admins = append(chat_admins, int(admin))
	}

	return dtos.ChatDto{
		Id:         int(chat.ID),
		Avatar:     SavedFileToFileDto(chat.Avatar),
		Title:      chat.Title,
		Type:       chat.Type,
		Members:    chat_members,
		IsArchived: chat.IsArchived,
		OwnerId:    int(chat.OwnerId),
		Admins:     chat_admins,
		Actions:    actions,
	}
}

func ChatDtoToDbChat(chat dtos.ChatDto) models.Chat {
	var chat_members pq.Int64Array
	var chat_admins pq.Int64Array
	for _, member := range chat.Members {
		chat_members = append(chat_members, int64(member))
	}

	for _, admin := range chat.Admins {
		chat_admins = append(chat_admins, int64(admin))
	}

	if chat.Id != 0 {
		return models.Chat{
			ID:         uint(chat.Id),
			Avatar:     FileDtoToSavedFile(chat.Avatar),
			Title:      chat.Title,
			Type:       chat.Type,
			Members:    chat_members,
			IsArchived: chat.IsArchived,
			OwnerId:    uint(chat.OwnerId),
			Admins:     chat_admins,
		}
	}

	return models.Chat{
		Avatar:     FileDtoToSavedFile(chat.Avatar),
		Title:      chat.Title,
		Type:       chat.Type,
		Members:    chat_members,
		IsArchived: chat.IsArchived,
		OwnerId:    uint(chat.OwnerId),
		Admins:     chat_admins,
	}
}

func ChatRequestToChatDto(request model.CreateChatRequest) (dtos.ChatDto, error) {
	if request.User != nil {
		return dtos.ChatDto{}, nil
	}

	if request.Title != nil && request.Avatar != nil && request.Members != nil {
		var converted_url *string
		var converted_filename *string
		if request.Avatar.Converted != nil {
			converted_url = &request.Avatar.Converted.URL
			converted_filename = &request.Avatar.Converted.Filename
		} else {
			converted_url = nil
			converted_filename = nil
		}

		avatar := dtos.FileDto{
			OriginalUrl:       request.Avatar.Original.URL,
			OriginalFilename:  request.Avatar.Original.Filename,
			ConvertedUrl:      converted_url,
			ConvertedFilename: converted_filename,
		}

		return dtos.ChatDto{
			Title:   *request.Title,
			Members: request.Members,
			Avatar:  avatar,
		}, nil
	}

	return dtos.ChatDto{}, fmt.Errorf("not enough data for chat creation")
}

func ChatActionDtoToResponse(chatAction dtos.ChatActionDto) model.ChatAction {
	var action_users []*model.ChatActionUser
	for _, action_user := range chatAction.ActionUsers {
		action_users = append(action_users, &model.ChatActionUser{
			Name: action_user.Name,
			ID:   action_user.Id,
		})
	}

	return model.ChatAction{
		Action:      model.ActionTypes(chatAction.Action),
		ActionUsers: action_users,
	}
}

func ChatDtoToResponse(chat dtos.ChatDto) model.Chat {
	avatar := FileDtoToSchema(chat.Avatar)
	var actions []*model.ChatAction
	if chat.Actions != nil {
		for _, action := range *chat.Actions {
			action_response := ChatActionDtoToResponse(action)
			actions = append(actions, &action_response)
		}
	}

	return model.Chat{
		ID:         chat.Id,
		Avatar:     &avatar,
		Title:      chat.Title,
		Type:       model.ChatType(chat.Type),
		Members:    chat.Members,
		IsArchived: chat.IsArchived,
		OwnerID:    chat.OwnerId,
		Admins:     chat.Admins,
		Actions:    actions,
	}
}
