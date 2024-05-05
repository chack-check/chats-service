package database

import (
	"log"

	"gorm.io/gorm"
)

func Paginate(page int, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		log.Printf("Paginating query. page = %d perPage = %d", page, perPage)
		if page <= 0 {
			page = 1
		}

		if perPage <= 20 {
			perPage = 20
		} else {
			perPage = 100
		}

		log.Printf("Fixed page = %d perPage = %d", page, perPage)
		offset := (page - 1) * perPage
		log.Printf("Offset: %v", offset)
		return db.Offset(offset).Limit(perPage)
	}
}
