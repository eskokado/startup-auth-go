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

func TestUpdatePasswordUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	validUserID := vo.NewID()
	currentHash, _ := vo.NewPasswordHash("valid_hash")
	validUser := &entity.User{
		ID:           validUserID,
		PasswordHash: currentHash,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "new_password", currentHash.String()).Return(false, nil)
		mockCrypto.On("Encrypt", "new_password").Return("new_hash", nil)
		mockRepo.On("Save", ctx, mock.Anything).Return(validUser, nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(nil, msgerror.AnErrUserNotFound)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorIs(t, err, msgerror.AnErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCurrentPassword", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "wrong_password", currentHash.String()).Return(false, nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "wrong_password", "new_password")

		assert.ErrorIs(t, err, msgerror.AnErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("SamePassword", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "current_password")

		assert.ErrorIs(t, err, usecase.ErrSamePassword)
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("EncryptError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "new_password", currentHash.String()).Return(false, nil)
		mockCrypto.On("Encrypt", "new_password").Return("", errors.New("encryption failed"))

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "failed to encrypt password")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("InvalidNewPassword", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "new_password", currentHash.String()).Return(false, nil)
		mockCrypto.On("Encrypt", "new_password").Return("", nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "invalid hash")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("CompareError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(false, errors.New("compare error"))

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "compare error")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("SaveError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "new_password", currentHash.String()).Return(false, nil)
		mockCrypto.On("Encrypt", "new_password").Return("new_hash", nil)
		mockRepo.On("Save", ctx, mock.Anything).Return(nil, errors.New("save error"))

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "save error")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("WithPasswordHashError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "new_password", currentHash.String()).Return(false, nil)
		mockCrypto.On("Encrypt", "new_password").Return("", nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "invalid hash")
		mockRepo.AssertNotCalled(t, "Save")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("NewPasswordComparisonError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
		mockCrypto.On("Compare", "new_password", currentHash.String()).Return(false, errors.New("comparison error"))

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "failed to verify password difference")
		mockRepo.AssertNotCalled(t, "Save")
		mockCrypto.AssertNotCalled(t, "Encrypt")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("NilUser", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(nil, nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorIs(t, err, msgerror.AnErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound_AnErrNotFound", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		// Return msgerror.AnErrNotFound specifically
		mockRepo.On("GetByID", ctx, validUserID).Return(nil, msgerror.AnErrNotFound)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorIs(t, err, msgerror.AnErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})
}
