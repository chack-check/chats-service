package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Message struct {
	*gorm.Model
	ID          uint           `gorm:"primaryKey" json:"id"`
	SenderId    uint           `json:"sender_id"`
	ChatId      uint           `json:"chat_id"`
	Type        string         `json:"type"`
	Content     string         `json:"content"`
	VoiceURL    string         `json:"voice_url"`
	CircleURL   string         `json:"circle_url"`
	Attachments pq.StringArray `gorm:"type:text[]" json:"attachments"`
	ReplyToID   uint           `json:"reply_to_id"`
	Mentioned   pq.Int32Array  `gorm:"type:integer[]" json:"mentioned"`
	ReadedBy    pq.Int32Array  `gorm:"type:integer[]" json:"readed_by"`
	Reactions   []Reaction     `gorm:"foreignKey:MessageId" json:"reactions"`
	DeletedFor  pq.Int32Array  `gorm:"type:integer[]" json:"deleted_for"`
}

type Reaction struct {
	*gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	MessageId uint   `json:"message_id"`
	UserId    uint   `json:"user_id"`
	Content   string `json:"content"`
}
