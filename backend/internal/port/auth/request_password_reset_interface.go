package port

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type RequestPasswordResetInterface interface {
	Execute(ctx context.Context, email vo.Email) error
}
