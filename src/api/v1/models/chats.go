package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Chat struct {
	*gorm.Model
	ID         uint `gorm:"primaryKey"`
	AvatarURL  string
	Title      string `gorm:"unique"`
	Type       string
	Members    pq.Int64Array `gorm:"type:integer[]"`
	IsArchived bool          `gorm:"default:false"`
	OwnerId    uint
	Admins     pq.Int64Array `gorm:"type:integer[]"`
}
