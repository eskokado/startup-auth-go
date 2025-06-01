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

		userRepo.On("GetByEmail", ctx, validEmail).Return(&entity.User{}, msgerror.AnErrNotFound)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("should handle save error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{Email: validEmail}

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, assert.AnError)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("should send reset email successfully", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{Email: validEmail}

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, nil)
		emailService.On("SendResetPasswordEmail", validEmail, mock.Anything).Return(nil)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.NoError(t, err)
		emailService.AssertExpectations(t)
	})

	t.Run("should return error when email sending fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)
		user := &entity.User{Email: validEmail}

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)
		userRepo.On("Save", ctx, user).Return(user, nil)
		emailService.On("SendResetPasswordEmail", validEmail, mock.Anything).Return(assert.AnError)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		assert.EqualError(t, err, assert.AnError.Error())
		emailService.AssertExpectations(t)
	})

	t.Run("should return error when token generation fails", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepo)
		emailService := new(mocks.MockEmailService)

		// Mock para usuário com falha na geração de token
		user := &entity.User{
			Email: validEmail,
		}

		// Sobrescreva temporariamente a função de geração de token
		originalGenToken := entity.GenerateSecureToken
		entity.GenerateSecureToken = func() (string, error) {
			return "", errors.New("token generation failed")
		}
		defer func() { entity.GenerateSecureToken = originalGenToken }()

		userRepo.On("GetByEmail", ctx, validEmail).Return(user, nil)

		uc := usecase.NewRequestPasswordReset(userRepo, emailService)
		err := uc.Execute(ctx, validEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token generation failed")
		userRepo.AssertExpectations(t)
		userRepo.AssertNotCalled(t, "Save")
		emailService.AssertNotCalled(t, "SendResetPasswordEmail")
	})
}
