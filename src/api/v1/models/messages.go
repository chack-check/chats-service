package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Message struct {
	*gorm.Model
	ID          uint           `gorm:"primaryKey" json:"id"`
	SenderId    uint           `json:"senderId"`
	ChatId      uint           `json:"chatId"`
	Type        string         `json:"type"`
	Content     string         `json:"content"`
	VoiceURL    string         `json:"voiceUrl"`
	CircleURL   string         `json:"circleUrl"`
	Attachments pq.StringArray `gorm:"type:text[]" json:"attachments"`
	ReplyToID   uint           `json:"replyToId"`
	Mentioned   pq.Int32Array  `gorm:"type:integer[]" json:"mentioned"`
	ReadedBy    pq.Int32Array  `gorm:"type:integer[]" json:"readedBy"`
	Reactions   []Reaction     `gorm:"foreignKey:MessageId" json:"reactions"`
	DeletedFor  pq.Int32Array  `gorm:"type:integer[]" json:"deletedFor"`
}

type Reaction struct {
	*gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	MessageId uint   `json:"messageId"`
	UserId    uint   `json:"userId"`
	Content   string `json:"content"`
}
