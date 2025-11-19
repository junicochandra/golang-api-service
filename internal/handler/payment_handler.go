package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junicochandra/golang-api-service/internal/app/payment"
	"github.com/junicochandra/golang-api-service/internal/app/payment/dto"
)

type PaymentHandler struct {
	usecase payment.TopUpUseCase
}

func NewPaymentHandler(uc payment.TopUpUseCase) *PaymentHandler {
	return &PaymentHandler{usecase: uc}
}

// @Tags         Payment
// @Summary      Create a top-up transaction
// @Description  Create a new top-up transaction and return a pending transaction id
// @Router       /payments/topup [post]
// @Accept       json
// @Produce      json
// @Param        request body dto.TopUpRequest true "TopUp request payload"
// @Success      202 "transaction accepted"
// @Failure      400 "bad request"
// @Failure      500 "internal server error"
func (h *PaymentHandler) CreateTopUp(c *gin.Context) {
	var req dto.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txID, err := h.usecase.CreateTopUp(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"transaction_id": txID, "status": "pending"})
}
