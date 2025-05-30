package port

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
)

type RegisterInterface interface {
	Execute(ctx context.Context, input dto.RegisterParams) error
}
