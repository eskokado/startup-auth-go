package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterWithInvalidName(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name: "A",
	})

	assert.ErrorContains(t, err, "invalid name")
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithInvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:  "Valid Name",
		Email: "invalid-email",
	})

	assert.ErrorContains(t, err, "invalid email")
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithExistingUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("existing@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return(&entity.User{}, nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "Existing User",
		Email:    "existing@test.com",
		Password: "any",
	})

	assert.ErrorIs(t, err, msgerror.AnErrUserExists)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterSuccessfully(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	name, _ := vo.NewName("New User", 0, 0)
	passwordHash, _ := vo.NewPasswordHash("hashed-password")
	email, _ := vo.NewEmail("new@test.com")

	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)

	newUser := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		ImageURL:     vo.URL{},
		CreatedAt:    time.Now(),
	}
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(newUser, nil)
	mockCrypto.On("Encrypt", "valid-password").Return("hashed-password", nil)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "New User",
		Email:    "new@test.com",
		Password: "valid-password",
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithUnexpectedErrorOnGetByEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	expectedErr := errors.New("unexpected error")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), expectedErr)

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "Valid Name",
		Email:    "test@test.com",
		Password: "any",
	})

	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Encrypt")
}

func TestRegisterWithEncryptError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("", errors.New("encryption failed"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "Valid Name",
		Email:    "test@test.com",
		Password: "any",
	})

	assert.ErrorContains(t, err, "failed to secure password")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithPasswordHashError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("", nil) // Retorna string vazia

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "Valid Name",
		Email:    "test@test.com",
		Password: "any",
	})

	assert.ErrorContains(t, err, "failed to secure password")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWithSaveError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("hashed-password", nil)

	// Configurar para retornar nil explicitamente
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "Valid Name",
		Email:    "test@test.com",
		Password: "any",
	})

	assert.ErrorContains(t, err, "failed to create user")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestRegisterWhenSaveReturnsNilUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return((*entity.User)(nil), msgerror.AnErrNotFound)
	mockCrypto.On("Encrypt", "any").Return("hashed-password", nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil, nil) // Retorna nil sem erro

	handler := usecase.NewRegisterUsecase(mockRepo, mockCrypto)
	err := handler.Execute(context.Background(), dto.RegisterParams{
		Name:     "Valid Name",
		Email:    "test@test.com",
		Password: "any",
	})

	assert.ErrorIs(t, err, msgerror.AnErrNoSavedUser)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}
