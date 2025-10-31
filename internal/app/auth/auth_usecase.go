package auth

import (
	"errors"
	"fmt"

	"github.com/junicochandra/golang-api-service/internal/app/auth/dto"
	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/domain/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/service"
)

var (
	ErrEmailExists        = errors.New("Email already exists")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepo: userRepo}
}

func (u *authUseCase) Register(req *dto.RegisterRequest) error {
	exists, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return err
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

func (u *authUseCase) Login(req *dto.UserAuthRequest) (string, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", err
	}

	// Compare password
	err = service.CheckPassword(req.Password, user.Password)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT
	token, err := service.GenerateToken(user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUseCase) Logout(token string) error {
	fmt.Println("Logout token:", token)

	return nil
}
