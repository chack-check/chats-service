package ports

import (
	"chats-service/internal/domain/entities"
)

type UsersRepositoryPort interface {
	GetById(id int) (*entities.User, error)
	GetByIds(ids []int) []entities.User
}
