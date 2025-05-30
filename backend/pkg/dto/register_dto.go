package dto

type RegisterParams struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImageURL string `json:"image_url"`
}

type RegisterUserOutput struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	ImageURL             string `json:"image_url"`
	PasswordResetToken   string `json:"password_reset_token"`
	PasswordResetExpires string `json:"password_reset_expires"`
}
