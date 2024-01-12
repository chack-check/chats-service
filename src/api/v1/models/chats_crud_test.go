package models

import (
	"os"
	"testing"

	"github.com/chack-check/chats-service/database"
	"github.com/lib/pq"
)

func setup() error {
	database.DB.AutoMigrate(&Chat{})
	return nil
}

func tearDown() error {
	database.DB.Migrator().DropTable(&Chat{})
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := tearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}

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

	if chat.ID != 1 || chat.OwnerId != 1 || chat.Title != "sometitle" {
		t.Fatalf("Getted chat is incorrect: %v", chat)
	}
}

func TestGetConcreteNotExisting(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetConcrete(1, 5)
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
	chat, err := chatsQueries.GetWithMember(1, 3)
	if err == nil {
		t.Fatal("Error when getting the chat with incorrect member: error is nil")
	}
	if chat.ID != 0 {
		t.Fatalf("Error when getting the chat with incorrect member: Chat id %d != 0", chat.ID)
	}
}

func TestGetWithMbmerIncorrectChatId(t *testing.T) {
	chatsQueries := ChatsQueries{}
	chat, err := chatsQueries.GetWithMember(2, 1)
	if err == nil {
		t.Fatal("Error when getting the chat with member with incorrect chat id: error is nil")
	}
	if chat.ID != 0 {
		t.Fatalf("Error when getting the chat with member with incorrect chat id: chat id %d != 0", chat.ID)
	}
}
