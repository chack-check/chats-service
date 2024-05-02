package utils

type PaginatedResponse[T any] struct {
	page       int
	perPage    int
	pagesCount int
	total      int
	data       []T
}

func (model *PaginatedResponse[T]) GetPage() int {
	return model.page
}

func (model *PaginatedResponse[T]) GetPerPage() int {
	return model.perPage
}

func (model *PaginatedResponse[T]) GetPagesCount() int {
	return model.pagesCount
}

func (model *PaginatedResponse[T]) GetTotal() int {
	return model.total
}

func (model *PaginatedResponse[T]) GetData() []T {
	return model.data
}

func (model *PaginatedResponse[T]) SetData(data []T) {
	model.data = data
}

type OffsetResponse[T any] struct {
	offset int
	limit  int
	total  int
	data   []T
}

func (model *OffsetResponse[T]) GetOffset() int {
	return model.offset
}

func (model *OffsetResponse[T]) SetOffset(offset int) {
	model.offset = offset
}

func (model *OffsetResponse[T]) GetLimit() int {
	return model.limit
}

func (model *OffsetResponse[T]) SetLimit(limit int) {
	model.limit = limit
}

func (model *OffsetResponse[T]) GetTotal() int {
	return model.total
}

func (model *OffsetResponse[T]) SetTotal(total int) {
	model.total = total
}

func (model *OffsetResponse[T]) GetData() []T {
	return model.data
}

func (model *OffsetResponse[T]) SetData(data []T) {
	model.data = data
}

func NewPaginatedResponse[T any](page, perPage, pagesCount, total int, data []T) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		page:       page,
		perPage:    perPage,
		pagesCount: pagesCount,
		total:      total,
		data:       data,
	}
}

func NewOffsetResponse[T any](offset, limit, total int, data []T) OffsetResponse[T] {
	return OffsetResponse[T]{
		offset: offset,
		limit:  limit,
		total:  total,
		data:   data,
	}
}
