package usecase

import (
	"context"
	"errors"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	service "github.com/eskokado/startup-auth-go/backend/pkg/domain/services"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
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
	email vo.Email, // Corrigido para vo.Email
) error {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	// Corrigido: tratar apenas o erro específico de não encontrado
	if err != nil {
		if errors.Is(err, msgerror.AnErrNotFound) {
			// Não revelar que o usuário não existe
			return nil
		}
		return msgerror.Wrap("failed to get user", err)
	}

	if err := user.GeneratePasswordResetToken(); err != nil {
		return msgerror.Wrap("failed to generate reset token", err)
	}

	if _, err := uc.userRepo.Save(ctx, user); err != nil {
		return msgerror.Wrap("failed to save user", err)
	}

	err = uc.emailSender.SendResetPasswordEmail(user.Email, user.PasswordResetToken)
	if err != nil {
		return msgerror.Wrap("failed to send reset email", err)
	}
	return nil
}
