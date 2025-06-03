package usecase

import (
	"context"
	"fmt"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type UpdatePasswordUseCase struct {
	userRepo       repository.UserRepository
	cryptoProvider providers.CryptoProvider
}

func NewUpdatePasswordUseCase(
	userRepo repository.UserRepository,
	cryptoProvider providers.CryptoProvider,
) *UpdatePasswordUseCase {
	return &UpdatePasswordUseCase{
		userRepo:       userRepo,
		cryptoProvider: cryptoProvider,
	}
}

func (uc *UpdatePasswordUseCase) Execute(
	ctx context.Context,
	userID vo.ID,
	currentPassword string,
	newPassword string,
) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	match, err := uc.cryptoProvider.Compare(currentPassword, user.PasswordHash.String())
	if err != nil {
		return fmt.Errorf("failed to compare passwords: %w", err)
	}
	if !match {
		return msgerror.AnErrInvalidCredentials
	}

	newHash, err := uc.cryptoProvider.Encrypt(newPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	newPasswordHash, err := vo.NewPasswordHash(newHash)
	if err != nil {
		return fmt.Errorf("invalid hash: %w", err)
	}

	updatedUser, _ := user.WithPasswordHash(newPasswordHash)
	_, err = uc.userRepo.Save(ctx, updatedUser)
	return err
}
