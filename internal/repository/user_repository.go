package repository

import (
	"errors"

	"github.com/junicochandra/golang-api-service/internal/config/database"
	"github.com/junicochandra/golang-api-service/internal/entity"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
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

func (repo *userRepository) FindByEmail(email string) (bool, error) {
	var exists bool
	err := repo.db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repo *userRepository) Update(user *entity.User) error {
	return repo.db.Save(user).Error
}

func (repo *userRepository) Delete(id uint64) error {
	return repo.db.Delete(&entity.User{}, id).Error
}
