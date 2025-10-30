package user

import "github.com/junicochandra/golang-api-service/internal/dto"

type UserUseCase interface {
	GetAll() ([]dto.UserListResponse, error)
	GetUserByID(id uint64) (*dto.UserDetailResponse, error)
	Create(user *dto.UserCreateRequest) error
	Update(id uint64, user *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error)
	Delete(id uint64) error
}
