package entity

type User struct {
	ID              uint64 `json:"id"`
	StaffID         string `json:"staff_id"`
	Username        string `json:"username"`
	Name            string `json:"name"`
	Contact         string `json:"contact"`
	Email           string `json:"email"`
	EmailVerifiedAt string `json:"email_verified_at"`
	Password        string `json:"password"`
	RememberToken   string `json:"remember_token"`
	Role            int    `json:"role"`
	Photo           string `json:"photo"`
	Address         string `json:"address"`
	Office          uint64 `json:"office"`
}
