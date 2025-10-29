package entity

import "time"

type User struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StaffID         *string    `gorm:"size:32" json:"staffId,omitempty"`
	Username        *string    `gorm:"size:255;unique" json:"username,omitempty"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	Contact         *string    `gorm:"size:32" json:"contact,omitempty"`
	Email           string     `gorm:"size:255;unique;not null" json:"email"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt,omitempty"`
	Password        string     `gorm:"size:255;not null" json:"password"`
	RememberToken   *string    `gorm:"size:100" json:"rememberToken,omitempty"`
	Role            *int8      `json:"role,omitempty"`
	Photo           *string    `gorm:"size:255" json:"photo,omitempty"`
	Address         *string    `gorm:"type:text" json:"address,omitempty"`
	Office          *int       `json:"office,omitempty"`
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
	DeletedBy       *string    `gorm:"size:32" json:"deletedBy,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}
