package database

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SavedFile struct {
	*gorm.Model
	ID                uint   `gorm:"primaryKey" json:"id"`
	OriginalUrl       string `json:"original_url"`
	OriginalFilename  string `json:"original_filename"`
	ConvertedUrl      string `json:"converted_url"`
	ConvertedFilename string `json:"converted_filename"`
}

type Chat struct {
	*gorm.Model
	ID         uint          `gorm:"primaryKey" json:"id"`
	AvatarId   *uint         `json:"avatar_id"`
	Avatar     SavedFile     `gorm:"foreignKey:AvatarId" json:"avatar"`
	Title      string        `json:"title"`
	Type       string        `json:"type"`
	Members    pq.Int64Array `gorm:"type:integer[]" json:"members"`
	IsArchived bool          `gorm:"default:false" json:"is_archived"`
	OwnerId    uint          `json:"owner_id"`
	Admins     pq.Int64Array `gorm:"type:integer[]" json:"admins"`
}

type Message struct {
	*gorm.Model
	ID          uint          `gorm:"primaryKey" json:"id"`
	SenderId    uint          `json:"sender_id"`
	ChatId      uint          `json:"chat_id"`
	Chat        Chat          `gorm:"foreignKey:ChatId"`
	Type        string        `json:"type"`
	Content     string        `json:"content"`
	VoiceId     *int          `json:"voice_id"`
	Voice       *SavedFile    `gorm:"foreignKey:VoiceId" json:"voice"`
	CircleId    *int          `json:"circle_id"`
	Circle      *SavedFile    `gorm:"foreignKey:CircleId" json:"circle"`
	Attachments []SavedFile   `gorm:"many2many:message_attachments" json:"attachments"`
	ReplyToID   uint          `json:"reply_to_id"`
	Mentioned   pq.Int32Array `gorm:"type:integer[]" json:"mentioned"`
	ReadedBy    pq.Int32Array `gorm:"type:integer[]" json:"readed_by"`
	Reactions   []Reaction    `gorm:"foreignKey:MessageId" json:"reactions"`
	DeletedFor  pq.Int32Array `gorm:"type:integer[]" json:"deleted_for"`
	CreatedAt   time.Time
}

type Reaction struct {
	*gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	MessageId uint   `json:"message_id"`
	UserId    uint   `json:"user_id"`
	Content   string `json:"content"`
}
