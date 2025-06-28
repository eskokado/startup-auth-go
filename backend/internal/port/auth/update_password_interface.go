package port

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type UpdatePasswordInterface interface {
	Execute(
		ctx context.Context,
		userID vo.ID,
		currentPassword string,
		newPassword string,
	) error
}
