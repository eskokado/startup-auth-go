package dto

import (
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type LoginResult struct {
	UserID    vo.ID
	Name      vo.Name
	Email     vo.Email
	ImageURL  vo.URL
	CreatedAt time.Time
	Token     string
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken string     `json:"access_token"`
	User        UserOutput `json:"user"`
}

type UserOutput struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
