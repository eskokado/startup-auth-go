package port

import "context"

type LogoutInterface interface {
	Execute(ctx context.Context, token string) error
}
