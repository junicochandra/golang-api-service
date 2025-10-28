package usecase

import (
	"errors"
	"fmt"

	"github.com/junicochandra/golang-api-service/internal/dto"
	"github.com/junicochandra/golang-api-service/internal/entity"
	"github.com/junicochandra/golang-api-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	GetAll() ([]dto.UserListResponse, error)
	GetById(id uint64) (*dto.UserDetailResponse, error)
	Create(user *dto.UserCreateRequest) (*dto.UserCreateResponse, error)
	Update(id uint64, user *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error)
	Delete(id uint64) error
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserRepository(r repository.UserRepository) UserUseCase {
	return &userUseCase{repo: r}
}

func (u *userUseCase) GetAll() ([]dto.UserListResponse, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserListResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserListResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return userResponses, nil
}

func (u *userUseCase) GetById(id uint64) (*dto.UserDetailResponse, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	userResponse := &dto.UserDetailResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return userResponse, nil
}

func (u *userUseCase) Create(req *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {
	// Check if email already exists
	exists, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	// Insert user into repository
	id, err := u.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserCreateResponse{ID: id, Name: user.Name, Email: user.Email}, nil
}

func (u *userUseCase) Update(id uint64, req *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error) {
	// Check if user exists
	user, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := u.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, fmt.Errorf("email already exists")
		}
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	_, err = u.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserUpdateResponse{
		ID:    fmt.Sprintf("%d", user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUseCase) Delete(id uint64) error {
	find, err := u.repo.FindById(id)
	if err != nil {
		return err
	}

	if find == nil {
		return errors.New("User not found")
	}

	_, err = u.repo.Delete(id)
	return err
}
