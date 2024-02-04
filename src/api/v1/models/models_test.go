package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/chack-check/chats-service/database"
)

func setup() error {
	database.DB.AutoMigrate(&Chat{})

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	baseDir, err := filepath.Abs(filepath.Dir(filepath.Dir(filepath.Dir(cwd))))
	if err != nil {
		panic(err)
	}

	chatsFilename := filepath.Join(baseDir, "test_data/chats.json")
	f, err := os.Open(chatsFilename)
	if err != nil {
		panic(err)
	}

	data := make([]byte, 1024*10)
	n, err := f.Read(data)
	if err != nil {
		panic(err)
	}

	chats := make([]Chat, 20)
	err = json.Unmarshal(data[:n], &chats)
	if err != nil {
		panic(err)
	}

	database.DB.Create(&chats)
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
