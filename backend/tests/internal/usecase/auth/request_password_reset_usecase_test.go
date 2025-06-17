package usecase_test

import (
	"context"
	"errors"
	"testing"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestPasswordResetUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	validEmail, _ := vo.NewEmail("user@example.com")

	t.Run("should return nil when user not found", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)

		userRepo.On("GetByEmail", ctx, validEmail).Return((*entity.User)(nil), msgerror.AnErrNotFound)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.NoError(t, err)
	})

	t.Run("should handle save error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{}

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, assert.AnError)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save user")
	})

	t.Run("should send reset email successfully", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{}

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, nil)
		emailService.On("SendResetPasswordEmail", mock.Anything, mock.Anything).Return(nil)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.NoError(t, err)
	})

	t.Run("should return error when email sending fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{}

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, nil)
		emailService.On("SendResetPasswordEmail", mock.Anything, mock.Anything).Return(assert.AnError)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send reset email")
	})

	t.Run("should return error when token generation fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{}

		// Mock falha na geração de token
		originalGenToken := entity.GenerateSecureToken
		entity.GenerateSecureToken = func() (string, error) {
			return "", errors.New("token generation failed")
		}
		defer func() { entity.GenerateSecureToken = originalGenToken }()

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate reset token")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)

		// Simular um erro de repositório (diferente de AnErrNotFound)
		expectedErr := errors.New("database connection failed")
		userRepo.On("GetByEmail", ctx, validEmail).Return((*entity.User)(nil), expectedErr)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user")
		assert.Contains(t, err.Error(), expectedErr.Error())
	})
}
