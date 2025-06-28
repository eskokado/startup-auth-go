package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	mocks "github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestResetPasswordUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	validToken := "valid-token"
	expiredToken := "expired-token"
	validPassword := "valid-password123"
	shortPassword := "short"
	emptyPassword := ""

	t.Run("should return invalid token error when user not found", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		userRepo.On("GetByResetToken", ctx, "invalid-token").Return(nil, nil)

		err := uc.Execute(ctx, "invalid-token", validPassword)

		assert.ErrorIs(t, err, msgerror.AnErrInvalidToken)
	})

	t.Run("should return wrapped error when repository returns an error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		expectedErr := errors.New("database error")
		userRepo.On("GetByResetToken", ctx, "invalid-token").Return(nil, expectedErr)

		err := uc.Execute(ctx, "invalid-token", validPassword)

		assert.ErrorContains(t, err, "falha ao buscar usuário pelo token")
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should return expired token error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(-1 * time.Hour),
		}

		userRepo.On("GetByResetToken", ctx, expiredToken).Return(user, nil)

		err := uc.Execute(ctx, expiredToken, validPassword)

		assert.ErrorIs(t, err, msgerror.AnErrExpiredToken)
	})

	t.Run("should return error for invalid new password", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(1 * time.Hour),
		}

		userRepo.On("GetByResetToken", ctx, validToken).Return(user, nil)

		err := uc.Execute(ctx, validToken, shortPassword)
		assert.ErrorIs(t, err, msgerror.AnErrPasswordInvalid)

		err = uc.Execute(ctx, validToken, emptyPassword)
		assert.ErrorIs(t, err, msgerror.AnErrPasswordInvalid)
	})

	t.Run("should reset password successfully", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetToken:   "original-token",
			PasswordResetExpires: time.Now().Add(1 * time.Hour),
		}

		userRepo.On("GetByResetToken", ctx, validToken).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, nil)

		err := uc.Execute(ctx, validToken, validPassword)

		assert.NoError(t, err)
		assert.Empty(t, user.PasswordResetToken)
		assert.True(t, user.PasswordResetExpires.IsZero())
		assert.NotEmpty(t, user.PasswordHash)
		userRepo.AssertExpectations(t)
	})

	t.Run("should handle save error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		uc := usecase.NewResetPassword(userRepo)

		user := &entity.User{
			PasswordResetExpires: time.Now().Add(1 * time.Hour),
		}

		expectedErr := errors.New("save failed")
		userRepo.On("GetByResetToken", ctx, validToken).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, expectedErr)

		err := uc.Execute(ctx, validToken, validPassword)

		assert.ErrorContains(t, err, "falha ao salvar usuário")
		assert.ErrorIs(t, err, expectedErr)
	})
}

// Helpers para simular falhas no vo.NewPasswordHash
var voNewPasswordHash = func(password string) (string, error) {
	return "hashed-password", nil
}
