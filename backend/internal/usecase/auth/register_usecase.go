package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type RegisterUsecase struct {
	userRepo       repository.UserRepository
	cryptoProvider providers.CryptoProvider
}

func NewRegisterUsecase(
	userRepo repository.UserRepository,
	cryptoProvider providers.CryptoProvider,
) *RegisterUsecase {
	return &RegisterUsecase{
		userRepo:       userRepo,
		cryptoProvider: cryptoProvider,
	}
}

func (h *RegisterUsecase) Execute(ctx context.Context, input dto.RegisterParams) error {
	name, err := vo.NewName(input.Name, 3, 100)
	if err != nil {
		return errors.New("invalid name: " + err.Error())
	}

	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return errors.New("invalid email: " + err.Error())
	}

	existingUser, err := h.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, msgerror.AnErrNotFound) {
		return err
	}
	if existingUser != nil {
		return msgerror.AnErrUserExists
	}

	hashedPassword, err := h.cryptoProvider.Encrypt(input.Password)
	if err != nil {
		return errors.New("failed to secure password")
	}

	passwordHashed, err := vo.NewPasswordHash(hashedPassword)
	if err != nil {
		return errors.New("failed to secure password")
	}

	newUser := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHashed,
		ImageURL:     vo.URL{},
		CreatedAt:    time.Now(),
	}

	savedUser, err := h.userRepo.Save(ctx, newUser)
	if err != nil {
		return errors.New("failed to create user: " + err.Error())
	}

	if savedUser == nil {
		return msgerror.AnErrNoSavedUser
	}

	return nil
}
