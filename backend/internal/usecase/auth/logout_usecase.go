package usecase

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
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
	// Validação básica do token
	if token == "" {
		return msgerror.AnErrTokenIsRequired
	}

	prefix := "startup-auth-go"
	keys := []string{
		prefix + ":" + token + ":UserID",
		prefix + ":" + token + ":Name",
		prefix + ":" + token + ":Email",
		prefix + ":" + token + ":Token",
		prefix + ":" + token + ":CreatedAt",
	}

	err := uc.blacklistProvider.Del(ctx, keys...)
	if err != nil {
		return msgerror.Wrap("failed to remove user session data", err)
	}

	return nil
}
