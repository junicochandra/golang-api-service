package repository

import (
	"errors"

	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	userRepo "github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userRepo.UserRepository {
	return &userRepository{db: database.DB}
}

func (repo *userRepository) GetAll() ([]entity.User, error) {
	var users []entity.User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) GetUserByID(id uint64) (*entity.User, error) {
	var user entity.User
	if err := repo.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) Create(user *entity.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) Update(user *entity.User) error {
	return repo.db.Save(user).Error
}

func (repo *userRepository) Delete(id uint64) error {
	return repo.db.Delete(&entity.User{}, id).Error
}
