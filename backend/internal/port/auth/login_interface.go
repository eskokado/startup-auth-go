package port

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
)

type LoginInterface interface {
	Execute(ctx context.Context, email string, password string) (dto.LoginResult, error)
}
