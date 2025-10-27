package entity

import "time"

type User struct {
	ID              uint64
	StaffID         string
	Username        string
	Name            string
	Contact         string
	Email           string
	Password        string
	RememberToken   string
	Role            int
	Photo           string
	Address         string
	Office          uint64
	LastLoginAt     *time.Time
	DeletedAt       *time.Time
	DeletedBy       string
	EmailVerifiedAt *time.Time
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
