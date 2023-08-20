package schemas

import (
	"math"
)

type PaginatedResponse[T interface{}] struct {
	Page       int
	PerPage    int
	PagesCount int
	Data       *[]T
}

func NewPaginatedResponse[T interface{}](page int, perPage int, count int, data []T) PaginatedResponse[T] {
	if perPage > count {
		perPage = count
	}
	if perPage < 20 {
		perPage = 20
	}

	pagesCount := math.Ceil(float64(count) / float64(perPage))

	if page <= 0 {
		page = 1
	}

	return PaginatedResponse[T]{
		Page:       page,
		PerPage:    perPage,
		PagesCount: int(pagesCount),
		Data:       &data,
	}
}
