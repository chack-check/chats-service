package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Message struct {
	*gorm.Model
	ID          uint `gorm:"primaryKey"`
	SenderId    uint
	ChatId      uint
	Type        string
	Content     string
	VoiceURL    string
	CircleURL   string
	Attachments pq.StringArray `gorm:"type:text[]"`
	ReplyToID   uint
	Mentioned   pq.Int32Array `gorm:"type:integer[]"`
	ReadedBy    pq.Int32Array `gorm:"type:integer[]"`
	Reactions   []Reaction    `gorm:"foreignKey:MessageId"`
}

type Reaction struct {
	*gorm.Model
	ID        uint `gorm:"primaryKey"`
	MessageId uint
	UserId    uint
	Content   string
}
