package port

import "context"

type ResetPasswordInterface interface {
	Execute(ctx context.Context, token, newPassword string) error
}
