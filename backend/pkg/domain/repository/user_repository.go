package repository

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByEmail(ctx context.Context, email vo.Email) (*entity.User, error)
	GetByID(ctx context.Context, userID vo.ID) (*entity.User, error)
	GetByResetToken(ctx context.Context, token string) (*entity.User, error)
}
