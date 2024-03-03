package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Chat struct {
	*gorm.Model
	ID         uint          `gorm:"primaryKey" json:"id"`
	AvatarURL  string        `json:"avatar_url"`
	Title      string        `json:"title"`
	Type       string        `json:"type"`
	Members    pq.Int64Array `gorm:"type:integer[]" json:"members"`
	IsArchived bool          `gorm:"default:false" json:"is_archived"`
	OwnerId    uint          `json:"owner_id"`
	Admins     pq.Int64Array `gorm:"type:integer[]" json:"admins"`
}
