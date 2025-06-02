package usecase

import (
	"context"
	"fmt"

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
	user, err := uc.userRepo.GetByID(ctx, userID) // Adicione GetByID ao UserRepository
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	validName, err := vo.NewName(newName, 3, 50)
	if err != nil {
		return msgerror.AnErrInvalidUser
	}

	updatedUser, err := user.WithName(validName)
	if err != nil {
		return err
	}

	_, err = uc.userRepo.Save(ctx, updatedUser)
	return err
}
