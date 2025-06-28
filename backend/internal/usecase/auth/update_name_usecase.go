package usecase

import (
	"context"
	"errors"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type UpdateNameUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateNameUseCase(userRepo repository.UserRepository) *UpdateNameUseCase {
	return &UpdateNameUseCase{userRepo: userRepo}
}

func (uc *UpdateNameUseCase) Execute(ctx context.Context, userID vo.ID, newName string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if errors.Is(err, msgerror.AnErrNotFound) {
		return msgerror.AnErrUserNotFound
	}
	if err != nil {
		return msgerror.Wrap("failed to get user", err)
	}

	// Adicionar verificação de user nil
	if user == nil {
		return msgerror.AnErrUserNotFound
	}

	validName, err := vo.NewName(newName, 3, 50)
	if err != nil {
		return msgerror.AnErrInvalidName
	}

	updatedUser, err := user.WithName(validName)
	if err != nil {
		if errors.Is(err, msgerror.AnErrNameDifferent) {
			return err
		}
		return msgerror.Wrap("failed to update name", err)
	}

	_, err = uc.userRepo.Save(ctx, updatedUser)
	if err != nil {
		return msgerror.Wrap("failed to save user", err)
	}

	return nil
}
