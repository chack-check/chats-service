package models

import (
	"slices"
	"testing"

	"github.com/lib/pq"
)

func TestCreateChat(t *testing.T) {
	members := pq.Int64Array{}
	members = append(members, 1)
	members = append(members, 2)

	admins := pq.Int64Array{}
	admins = append(admins, 1)
	admins = append(admins, 2)

	chatsQueries := ChatsQueries{}

	creatingChat := Chat{
		AvatarURL:  "someurl",
		Title:      "sometitle",
		Type:       "sometype",
		Members:    members,
		IsArchived: false,
		OwnerId:    1,
		Admins:     admins,
	}

	err := chatsQueries.Create(&creatingChat)
	if err != nil {
		t.Fatalf("Error when creating chat: %s", err)
	}
	if creatingChat.ID < 1 {
		t.Fatalf("Error when creating chat: id = %d", creatingChat.ID)
	}
}

func TestGetConcrete(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetConcrete(1, 1)
	if err != nil {
		t.Fatalf("Error when getting an existing chat: %s", err)
	}

	if chat.ID != 1 {
		t.Fatalf("Getted chat is incorrect: %v", chat)
	}
}

func TestGetConcreteNotExisting(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetConcrete(1, 525)
	if err == nil {
		t.Fatalf("Error when getting not existing chat: error is nil. Getted chat ID: %d", chat.ID)
	}

	if chat.ID != 0 {
		t.Fatalf("Error when getting not existing chat. Getted: %v", chat)
	}
}

func TestGetWithMember(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetWithMember(1, 2)
	if err != nil {
		t.Fatalf("Error when getting chat with member: %s", err)
	}
	if chat.ID != 1 {
		t.Fatalf("Error when getting chat with member: Chat id %d != 1", chat.ID)
	}
}

func TestGetWithAnotherMember(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetWithMember(1, 1)
	if err != nil {
		t.Fatalf("Error when getting chat with member: %s", err)
	}
	if chat.ID != 1 {
		t.Fatalf("Error when getting chat with member: Chat id %d != 1", chat.ID)
	}
}

func TestGetWithIncorrectMember(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetWithMember(1, 525)
	if err == nil {
		t.Fatal("Error when getting the chat with incorrect member: error is nil")
	}
	if chat.ID != 0 {
		t.Fatalf("Error when getting the chat with incorrect member: Chat id %d != 0", chat.ID)
	}
}

func TestGetWithMbmerIncorrectChatId(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetWithMember(525, 1)
	if err == nil {
		t.Fatal("Error when getting the chat with member with incorrect chat id: error is nil")
	}
	if chat.ID != 0 {
		t.Fatalf("Error when getting the chat with member with incorrect chat id: chat id %d != 0", chat.ID)
	}
}

func TestGetAllWithMember(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chats := chatsQueries.GetAllWithMember(1, 1, 20)
	if len(*chats) < 1 {
		t.Fatalf("Error when getting all chats with member 1: Getted %d chats != %d", len(*chats), 1)
	}
	for _, chat := range *chats {
		if !slices.Contains(chat.Members, 1) {
			t.Fatalf("Error when getting all chats with member 1: chat owner id is %d != %d", (*chats)[0].OwnerId, 1)
		}
	}
}

func TestGetAllWithMemberAnotherMember(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chats := chatsQueries.GetAllWithMember(2, 1, 20)
	if len(*chats) < 1 {
		t.Fatalf("Error when getting all chats with member 2: Getted %d chats != %d", len(*chats), 1)
	}
	for _, chat := range *chats {
		if !slices.Contains(chat.Members, 2) {
			t.Fatalf("Error when getting all chats with member 1: chat owner id is %d != %d", (*chats)[0].OwnerId, 1)
		}
	}
}

func TestGetAllWithMemberIncorrectMember(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chats := chatsQueries.GetAllWithMember(525, 1, 20)
	if len(*chats) != 0 {
		t.Fatalf("Error when getting all chats with incorrect member: founded %d chats: %v", len(*chats), *chats)
	}
}

func TestGetAll(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chats := chatsQueries.GetAll(1)
	if len(*chats) < 1 {
		t.Fatalf("Error when getting all chats: getted %d chats: %v", len(*chats), *chats)
	}
	for _, chat := range *chats {
		if chat.OwnerId != 1 {
			t.Fatalf("Error when getting all chats: getted chat owner id is %d, not %d", (*chats)[0].OwnerId, 1)
		}
	}
}

func TestGetAllWithIncorrectOwner(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chats := chatsQueries.GetAll(525)
	if len(*chats) != 0 {
		t.Fatalf("Error when getting all chats: getted %d chats != %d. Value: %v", len(*chats), 0, *chats)
	}
}

// func TestSearchChats(t *testing.T) {
// 	chatsQueries := ChatsQueries{}
// 	chats := chatsQueries.Search(1, "sometitle", 1, 20)
// 	if len(*chats) != 1 {
// 		t.Fatalf("Error when searching chats: searched %d != %d. Value: %v", len(*chats), 1, *chats)
// 	}
// 	if (*chats)[0].Title != "sometitle" {
// 		t.Fatalf("Error when searching chats: searched chat title %s != %s", (*chats)[0].Title, "sometitle")
// 	}
// }

func TestGetExistingWithUser(t *testing.T) {
	chatsQueries := ChatsQueries{}
	existing := chatsQueries.GetExistingWithUser(1, 2)
	if !existing {
		t.Fatalf("Error when getting is chat already existing between two users: 1 and 2: Chat doesn't exist")
	}
}

func TestGetExistingWithUserIncorrectUsers(t *testing.T) {
	chatsQueries := ChatsQueries{}
	existing := chatsQueries.GetExistingWithUser(1, 525)
	if existing {
		t.Fatalf("Error when getting is chat already existing between two users: 1 and 5: Chat exists")
	}
}
