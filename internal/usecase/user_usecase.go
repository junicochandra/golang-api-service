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
	GetUserByID(id uint64) (*dto.UserDetailResponse, error)
	Create(user *dto.UserCreateRequest) error
	Update(id uint64, user *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error)
	Delete(id uint64) error
}

var (
	ErrNotFound    = errors.New("User not found")
	ErrEmailExists = errors.New("Email already exists")
)

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserRepository(r repository.UserRepository) UserUseCase {
	return &userUseCase{repo: r}
}

func (u *userUseCase) GetAll() ([]dto.UserListResponse, error) {
	users, err := u.repo.GetAll()
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

func (u *userUseCase) GetUserByID(id uint64) (*dto.UserDetailResponse, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	userResponse := &dto.UserDetailResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return userResponse, nil
}

func (u *userUseCase) Create(req *dto.UserCreateRequest) error {
	// Check if email already exists
	exists, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}

	if exists {
		return fmt.Errorf("email already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	// Insert user into repository
	return u.repo.Create(user)
}

func (u *userUseCase) Update(id uint64, req *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error) {
	// Check if user exists
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := u.repo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, ErrEmailExists
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := u.repo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserUpdateResponse{
		ID:    fmt.Sprintf("%d", user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUseCase) Delete(id uint64) error {
	find, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if find == nil {
		return ErrNotFound
	}

	return u.repo.Delete(id)
}
