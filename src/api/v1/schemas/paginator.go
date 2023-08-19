package schemas

type PaginatedResponse[T interface{}] struct {
	Page       int
	PerPage    int
	PagesCount int
	Data       *[]T
}
