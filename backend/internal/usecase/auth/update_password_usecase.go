package usecase

import (
	"context"
	"errors"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

var (
	ErrSamePassword = errors.New("new password must be different")
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
	if errors.Is(err, msgerror.AnErrNotFound) {
		return msgerror.AnErrUserNotFound
	}
	if err != nil {
		return msgerror.Wrap("failed to get user", err)
	}
	if user == nil {
		return msgerror.AnErrUserNotFound
	}

	// Verify current password
	match, err := uc.cryptoProvider.Compare(currentPassword, user.PasswordHash.String())
	if err != nil {
		return msgerror.Wrap("failed to compare passwords", err)
	}
	if !match {
		return msgerror.AnErrInvalidCredentials
	}

	// Check if new password is different
	same, err := uc.cryptoProvider.Compare(newPassword, user.PasswordHash.String())
	if err != nil {
		return msgerror.Wrap("failed to verify password difference", err)
	}
	if same {
		return ErrSamePassword
	}

	// Encrypt new password
	newHash, err := uc.cryptoProvider.Encrypt(newPassword)
	if err != nil {
		return msgerror.Wrap("failed to encrypt password", err)
	}

	newPasswordHash, err := vo.NewPasswordHash(newHash)
	if err != nil {
		return msgerror.Wrap("invalid hash", err)
	}

	updatedUser, err := user.WithPasswordHash(newPasswordHash)
	if err != nil {
		return msgerror.Wrap("failed to update password", err)
	}

	_, err = uc.userRepo.Save(ctx, updatedUser)
	if err != nil {
		return msgerror.Wrap("failed to save user", err)
	}

	return nil
}
