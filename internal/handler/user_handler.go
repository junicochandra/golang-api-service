package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/junicochandra/golang-api-service/internal/dto"
	usecase "github.com/junicochandra/golang-api-service/internal/usecase/user"
)

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{usecase: uc}
}

// @Tags Users
// @Summary Get all users
// @Description Get all users from database
// @Router /users [get]
// @Accept json
// @Produce json
// @Success 200 {array} dto.UserListResponse
// @Failure 400
func (h *UserHandler) GetUsers(c *gin.Context) {
	list, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Tags Users
// @Summary Get user detail
// @Description Get detail information of a specific user by ID
// @Router /users/{id} [get]
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDetailResponse
// @Failure 400
// @Failure 404
// @Failure 500
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.usecase.GetUserByID(uint64(id64))
	if err != nil {
		if err == usecase.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": usecase.ErrNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Tags Users
// @Summary Create a new user
// @Description Add a new user to the database
// @Router /users [post]
// @Accept json
// @Produce json
// @Param user body dto.UserCreateRequest true "User data"
// @Success 201 {object} dto.UserCreateResponse
// @Failure 400
// @Failure 500
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// @Tags Users
// @Summary Update user
// @Description Update user data by ID
// @Router /users/{id} [put]
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UserUpdateRequest true "Updated user data"
// @Success 200 {object} dto.UserUpdateResponse
// @Failure 400
// @Failure 404
// @Failure 500
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.Update(id, &req)
	if err != nil {
		if err == usecase.ErrEmailExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": usecase.ErrEmailExists.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags Users
// @Summary Delete a user
// @Description Delete a user from the database by ID
// @Router /users/{id} [delete]
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	err := h.usecase.Delete(id)
	if err != nil {
		if err.Error() == "User not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
