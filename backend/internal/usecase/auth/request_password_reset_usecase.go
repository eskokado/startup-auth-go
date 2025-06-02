package usecase

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	service "github.com/eskokado/startup-auth-go/backend/pkg/domain/services"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type RequestPasswordResetUsecase struct {
	userRepo    repository.UserRepository
	emailSender service.EmailServiceInterface
}

func NewRequestPasswordReset(
	repo repository.UserRepository,
	emailSender service.EmailServiceInterface,
) *RequestPasswordResetUsecase {
	return &RequestPasswordResetUsecase{userRepo: repo, emailSender: emailSender}
}

func (uc *RequestPasswordResetUsecase) Execute(
	ctx context.Context,
	email vo.Email,
) error {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil
	}

	if err := user.GeneratePasswordResetToken(); err != nil {
		return err
	}

	if _, err := uc.userRepo.Save(ctx, user); err != nil {
		return err
	}

	err = uc.emailSender.SendResetPasswordEmail(user.Email, user.PasswordResetToken)
	if err != nil {
		return err
	}
	return nil
}
