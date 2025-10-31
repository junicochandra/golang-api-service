package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecase "github.com/junicochandra/golang-api-service/internal/app/auth"
	"github.com/junicochandra/golang-api-service/internal/app/auth/dto"
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
// @Param user body dto.RegisterRequest true "Register data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400
// @Failure 500
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
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

// @Tags Auth
// @Summary Login user
// @Description Login user
// @Router /auth/login [post]
// @Accept json
// @Produce json
// @Param user body dto.UserAuthRequest true "Login data"
// @Success 200 {object} dto.UserAuthResponse
// @Failure 400
// @Failure 500
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.UserAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.Login(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags Auth
// @Summary Logout user
// @Description Logout user by invalidating JWT (client-side)
// @Security BearerAuth
// @Router /auth/logout [post]
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
func (h *AuthHandler) Logout(c *gin.Context) {
	autHeader := c.GetHeader("Authorization")
	if autHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Authorization header"})
		return
	}

	token := autHeader[len("Bearer "):]
	if err := h.usecase.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}
