package port

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type UpdateNameInterface interface {
	Execute(ctx context.Context, userID vo.ID, newName string) error
}
