package models

import (
	"gorm.io/gorm"
)

func Paginate(page *int, perPage *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageValue := *page
		perPageValue := *perPage
		if pageValue <= 0 {
			pageValue = 1
		}

		if perPageValue < 20 {
			perPageValue = 20
		} else {
			perPageValue = 100
		}

		offset := (pageValue - 1) * perPageValue
		return db.Offset(offset).Limit(perPageValue)
	}
}
