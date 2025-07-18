package dto

type UpdateNameInput struct {
	Name string `json:"name"`
}

type UpdatePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token    string `json:"reset_password_token" binding:"required"`
	Password string `json:"password" binding:"required"`
}
