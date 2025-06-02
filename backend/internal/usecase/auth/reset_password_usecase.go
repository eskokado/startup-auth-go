package usecase

import (
	"context"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type ResetPasswordUsecase struct {
	userRepo repository.UserRepository
}

func NewResetPassword(repo repository.UserRepository) *ResetPasswordUsecase {
	return &ResetPasswordUsecase{userRepo: repo}
}

func (uc *ResetPasswordUsecase) Execute(
	ctx context.Context,
	token, newPassword string,
) error {
	user, err := uc.userRepo.GetByResetToken(ctx, token)
	if err != nil || user == nil {
		return msgerror.AnErrInvalidToken
	}

	if user.PasswordResetExpires.Before(time.Now()) {
		return msgerror.AnErrExpiredToken
	}

	newHash, err := vo.NewPasswordHash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	user.ClearResetToken()

	_, err = uc.userRepo.Save(ctx, user)

	return err
}
