package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginWithInvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto)
	_, err := handler.Execute(context.Background(), "invalid-email", "any")

	assert.ErrorIs(t, err, msgerror.AnErrInvalidEmail)
	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithNonExistentUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("nonexistent@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return(nil, msgerror.AnErrNotFound)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto)
	_, err := handler.Execute(context.Background(), "nonexistent@test.com", "any")

	assert.ErrorIs(t, err, msgerror.AnErrUserNotFound)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithUnexpectedErrorOnGetByEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	email, _ := vo.NewEmail("test@test.com")
	expectedErr := errors.New("unexpected error")
	mockRepo.On("GetByEmail", mock.Anything, email).Return(nil, expectedErr)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto)
	_, err := handler.Execute(context.Background(), "test@test.com", "any")

	assert.ErrorIs(t, err, expectedErr)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithInvalidPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	email, _ := vo.NewEmail("user@test.com")

	user := &entity.User{
		PasswordHash: passwordHash,
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

	mockCrypto.On("Compare", "wrong-password", mock.Anything).Return(false, nil)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto)
	_, err := handler.Execute(context.Background(), "user@test.com", "wrong-password")

	assert.ErrorIs(t, err, msgerror.AnErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestLoginWithCompareError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	email, _ := vo.NewEmail("user@test.com")

	user := &entity.User{
		PasswordHash: passwordHash,
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	compareErr := errors.New("comparison failed")

	mockCrypto.On("Compare", "any-password", mock.Anything).Return(false, compareErr)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto)
	_, err := handler.Execute(context.Background(), "user@test.com", "any-password")

	assert.ErrorContains(t, err, "failed to verify password")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestLoginSuccessfully(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)

	name, _ := vo.NewName("Test User", 0, 0)
	email, _ := vo.NewEmail("user@test.com")

	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)

	user := &entity.User{
		ID:           vo.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

	mockCrypto.On("Compare", "valid-password", validHash).Return(true, nil)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto)
	result, err := handler.Execute(context.Background(), "user@test.com", "valid-password")

	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.CreatedAt, result.CreatedAt)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}
