package service

import (
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type EmailServiceInterface interface {
	SendResetPasswordEmail(email vo.Email, token string) error
}
