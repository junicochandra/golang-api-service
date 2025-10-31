package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags Middleware Test
// @Summary Get user profile with middleware
// @Description Get user profile
// @Router /profile [get]
// @Security BearerAuth
// @Success 200
// @Failure 400
// @Failure 500
func Profile(c *gin.Context) {
	// Get data from middleware (claimed email from token)
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your profile!",
		"email":   email,
	})
}
