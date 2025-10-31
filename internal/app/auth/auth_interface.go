package auth

import "github.com/junicochandra/golang-api-service/internal/app/auth/dto"

type AuthUseCase interface {
	Register(user *dto.RegisterRequest) error
	Login(user *dto.UserAuthRequest) (string, error)
	Logout(token string) error
}
