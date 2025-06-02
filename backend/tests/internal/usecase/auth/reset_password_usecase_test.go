package usecase_test

import (
	"context"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	mocks "github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResetUserRepository struct {
	mock.Mock
}

func (m *MockResetUserRepository) GetByResetToken(ctx context.Context, token string) (*entity.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockResetUserRepository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestResetPasswordUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	validToken := "valid-token"
	expiredToken := "expired-token"
	newPassword := "new-password"

	t.Run("should return invalid token error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		userRepo.On("GetByResetToken", ctx, "invalid-token").Return(nil, msgerror.AnErrNotFound)

		err := uc.Execute(ctx, "invalid-token", newPassword)

		assert.Equal(t, msgerror.AnErrInvalidToken, err)
	})

	t.Run("should return expired token error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(-1 * time.Hour),
		}

		userRepo.On("GetByResetToken", ctx, expiredToken).Return(user, nil)

		err := uc.Execute(ctx, expiredToken, newPassword)

		assert.Equal(t, msgerror.AnErrExpiredToken, err)
	})

	t.Run("should reset password successfully", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(1 * time.Hour),
		}

		userRepo.On("GetByResetToken", ctx, validToken).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, nil)

		err := uc.Execute(ctx, validToken, newPassword)

		assert.NoError(t, err)
		assert.Empty(t, user.PasswordResetToken)
		assert.True(t, user.PasswordResetExpires.IsZero())
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("should handle save error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(1 * time.Hour),
		}

		userRepo.On("GetByResetToken", ctx, validToken).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, assert.AnError)

		err := uc.Execute(ctx, validToken, newPassword)

		assert.Error(t, err)
	})

	t.Run("should return error for invalid new password", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(1 * time.Hour),
		}

		invalidPassword := ""

		userRepo.On("GetByResetToken", ctx, validToken).Return(user, nil)

		err := uc.Execute(ctx, validToken, invalidPassword)

		assert.ErrorIs(t, err, msgerror.AnErrPasswordInvalid)
		userRepo.AssertNotCalled(t, "Save")
	})
}
