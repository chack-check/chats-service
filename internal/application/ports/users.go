package ports

import (
	"chats-service/internal/domain/entities"
)

type UsersRepository interface {
	GetById(id int) (*entities.User, error)
	GetByIds(ids []int) []entities.User
}
