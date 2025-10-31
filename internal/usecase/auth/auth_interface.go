package auth

import "github.com/junicochandra/golang-api-service/internal/dto"

type AuthUseCase interface {
	Register(user *dto.UserCreateRequest) error
	Login(user *dto.UserAuthRequest) (string, error)
	Logout(token string) error
}
