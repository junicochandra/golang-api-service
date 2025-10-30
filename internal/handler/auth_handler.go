package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junicochandra/golang-api-service/internal/dto"
	usecase "github.com/junicochandra/golang-api-service/internal/usecase/auth"
)

type AuthHandler struct {
	usecase usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

// @Tags Auth
// @Summary Register a new user
// @Description Add a new user to the database
// @Router /auth/register [post]
// @Accept json
// @Produce json
// @Param user body dto.UserCreateRequest true "Register data"
// @Success 201 {object} dto.UserCreateResponse
// @Failure 400
// @Failure 500
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.Register(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}
