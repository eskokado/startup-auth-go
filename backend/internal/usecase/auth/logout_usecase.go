package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
)

type LogoutUsecase struct {
	blacklistProvider providers.BlacklistProvider
}

func NewLogoutUsecase(blacklistProvider providers.BlacklistProvider) *LogoutUsecase {
	return &LogoutUsecase{
		blacklistProvider: blacklistProvider,
	}
}

func (uc *LogoutUsecase) Execute(ctx context.Context, token string) error {
	exists, err := uc.blacklistProvider.Exists(ctx, token)
	if err != nil {
		return errors.New("failed to verify token status")
	}

	if exists {
		return errors.New("token already revoked")
	}

	return uc.blacklistProvider.Add(ctx, token, 24*time.Hour)
}
