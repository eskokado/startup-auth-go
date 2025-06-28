package usecase

import (
	"context"

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
	// exists, err := uc.blacklistProvider.Exists(ctx, token)
	// if err != nil {
	// 	return msgerror.Wrap("failed to verify token status", err)
	// }

	// if !exists {
	// 	return msgerror.AnErrInvalidToken
	// }

	// if err := uc.blacklistProvider.Add(ctx, token, 0); err != nil {
	// 	return msgerror.Wrap("failed to revoke token", err)
	// }

	return nil
}
