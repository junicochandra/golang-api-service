package payment

import "github.com/junicochandra/golang-api-service/internal/app/payment/dto"

type TopUpUseCase interface {
	CreateTopUp(account *dto.TopUpRequest) (*dto.TopUpResponse, error)
}
