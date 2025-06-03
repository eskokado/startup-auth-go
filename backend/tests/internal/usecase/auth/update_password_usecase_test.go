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

		mockRepo.On("GetByID", ctx, validUserID).Return(&entity.User{}, msgerror.AnErrUserNotFound)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "failed to get user")
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

	t.Run("EncryptError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)
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
		// Simular erro na comparação
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
		mockCrypto.On("Encrypt", "new_password").Return("new_hash", nil)
		// Simular erro no save
		mockRepo.On("Save", ctx, mock.Anything).Return(&entity.User{}, errors.New("save error"))

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "save error")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("InvalidNewPasswordHash", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		validUser := &entity.User{
			ID:           validUserID,
			PasswordHash: currentHash,
		}

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)

		// Retornar string vazia que deve fazer vo.NewPasswordHash falhar
		mockCrypto.On("Encrypt", "new_password").Return("", nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid hash")

		// Verificar que Save não foi chamado pois o erro aconteceu antes
		mockRepo.AssertNotCalled(t, "Save")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("NoSaveOnError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		// Simular erro na comparação
		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(false, errors.New("compare error"))

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.ErrorContains(t, err, "compare error")
		mockRepo.AssertNotCalled(t, "Save")
		mockCrypto.AssertNotCalled(t, "Encrypt") // Corrigido: era mockRepo.AssertNotCalled
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("WithPasswordHashError", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		validUser := &entity.User{
			ID:           validUserID,
			PasswordHash: currentHash,
		}

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)

		// Retornar uma string vazia que deve causar erro no vo.NewPasswordHash
		mockCrypto.On("Encrypt", "new_password").Return("", nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid hash")

		// Verificar que Save não foi chamado
		mockRepo.AssertNotCalled(t, "Save")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("InvalidNewPasswordFromEmptyEncrypt", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepo)
		mockCrypto := new(mocks.MockCrypto)

		mockRepo.On("GetByID", ctx, validUserID).Return(validUser, nil)
		mockCrypto.On("Compare", "current_password", currentHash.String()).Return(true, nil)

		// Encrypt retorna string vazia - isso deve fazer vo.NewPasswordHash falhar
		mockCrypto.On("Encrypt", "new_password").Return("", nil)

		uc := usecase.NewUpdatePasswordUseCase(mockRepo, mockCrypto)
		err := uc.Execute(ctx, validUserID, "current_password", "new_password")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid hash")
		mockRepo.AssertNotCalled(t, "Save")
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})
}
