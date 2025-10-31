package repository

import "github.com/junicochandra/golang-api-service/internal/domain/entity"

type UserRepository interface {
	GetAll() ([]entity.User, error)
	GetUserByID(id uint64) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint64) error
}
