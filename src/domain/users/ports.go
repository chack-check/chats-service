package users

type UsersPort interface {
	GetById(id int) (*User, error)
	GetByIds(ids []int) []User
}
