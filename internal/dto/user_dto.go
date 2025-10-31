package dto

// Request / Response DTOs for User API

import "time"

// Create
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserCreateResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Detail
type UserDetailResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// List
type UserListResponse struct {
	ID              uint64     `json:"id"`
	StaffID         string     `json:"staff_id,omitempty"`
	Username        string     `json:"username,omitempty"`
	Name            string     `json:"name,omitempty"`
	Contact         string     `json:"contact,omitempty"`
	Email           string     `json:"email,omitempty"`
	Role            int        `json:"role,omitempty"`
	Photo           string     `json:"photo,omitempty"`
	Address         string     `json:"address,omitempty"`
	Office          uint64     `json:"office,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// Update
type UserUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,email"`
}

type UserUpdateResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Auth
type UserAuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
