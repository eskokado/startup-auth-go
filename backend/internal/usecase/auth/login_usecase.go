package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/repository"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
)

type LoginUsecase struct {
	userRepo          repository.UserRepository
	cryptoProvider    providers.CryptoProvider
	tokenProvider     providers.TokenProvider
	blacklistProvider providers.BlacklistProvider
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	cryptoProvider providers.CryptoProvider,
	tokenProvider providers.TokenProvider,
	blacklistProvider providers.BlacklistProvider,
) *LoginUsecase {
	return &LoginUsecase{
		userRepo:          userRepo,
		cryptoProvider:    cryptoProvider,
		tokenProvider:     tokenProvider,
		blacklistProvider: blacklistProvider,
	}
}

func (h *LoginUsecase) Execute(ctx context.Context, email string, password string) (dto.LoginResult, error) {
	validEmail, err := vo.NewEmail(email)
	if err != nil {
		return dto.LoginResult{}, msgerror.AnErrInvalidEmail
	}

	user, err := h.userRepo.GetByEmail(ctx, validEmail)
	if errors.Is(err, msgerror.AnErrNotFound) {
		return dto.LoginResult{}, msgerror.AnErrUserNotFound
	}
	if err != nil {
		return dto.LoginResult{}, err
	}

	// Corrigido: converter PasswordHash para string diretamente
	match, err := h.cryptoProvider.Compare(password, string(user.PasswordHash.String()))
	if err != nil {
		return dto.LoginResult{}, errors.New("failed to verify password")
	}
	if !match {
		return dto.LoginResult{}, msgerror.AnErrInvalidCredentials
	}

	token, err := h.tokenProvider.Generate(user.ID)
	if err != nil {
		return dto.LoginResult{}, errors.New("failed to generate token")
	}

	ttl := 24 * time.Hour

	if err := h.blacklistProvider.Add(ctx, token, ttl); err != nil {
		return dto.LoginResult{}, errors.New("failed to secure session")
	}

	return dto.LoginResult{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Token:     token,
	}, nil
}
