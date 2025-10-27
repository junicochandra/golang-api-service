package user

type UserDetailRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,email"`
}

type UserUpdateResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
