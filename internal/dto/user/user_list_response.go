package user

import "time"

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
