package user

import (
	"errors"
	"fmt"

	"github.com/junicochandra/golang-api-service/internal/app/user/dto"
	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/service"
)

var (
	ErrNotFound    = errors.New("User not found")
	ErrEmailExists = errors.New("Email already exists")
)

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) GetAll() ([]dto.UserListResponse, error) {
	users, err := u.userRepo.GetAll()
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
	user, err := u.userRepo.GetUserByID(id)
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
	exists, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}

	if exists != nil {
		return ErrEmailExists
	}

	// Hash password
	hashed, err := service.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	return u.userRepo.Create(user)
}

func (u *userUseCase) Update(id uint64, req *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error) {
	// Check if user exists
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := u.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}

		if exists != nil {
			return nil, ErrEmailExists
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserUpdateResponse{
		ID:    fmt.Sprintf("%d", user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUseCase) Delete(id uint64) error {
	find, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	if find == nil {
		return ErrNotFound
	}

	return u.userRepo.Delete(id)
}
