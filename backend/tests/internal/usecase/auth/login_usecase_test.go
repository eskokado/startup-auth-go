package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/entity"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/providers"
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginWithInvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)

	// Teste com e-mail inválido
	_, err := handler.Execute(context.Background(), "invalid-email", "any")

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Contains(t, valErr.FieldErrors["email"], "invalid email format")

	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithEmptyEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)

	_, err := handler.Execute(context.Background(), "", "any")

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "cannot be empty", valErr.FieldErrors["email"])

	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithShortPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)

	_, err := handler.Execute(context.Background(), "valid@test.com", "short")

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "must be at least 8 characters", valErr.FieldErrors["password"])

	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithNonExistentUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	email, _ := vo.NewEmail("nonexistent@test.com")
	mockRepo.On("GetByEmail", mock.Anything, email).Return(nil, msgerror.AnErrNotFound)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "nonexistent@test.com", "valid-password")

	assert.ErrorIs(t, err, msgerror.AnErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithUnexpectedErrorOnGetByEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	email, _ := vo.NewEmail("test@test.com")
	expectedErr := errors.New("unexpected error")
	mockRepo.On("GetByEmail", mock.Anything, email).Return(nil, expectedErr)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "test@test.com", "valid-password")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user")
	assert.Contains(t, err.Error(), "unexpected error")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertNotCalled(t, "Compare")
}

func TestLoginWithInvalidPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	email, _ := vo.NewEmail("user@test.com")

	user := &entity.User{
		PasswordHash: passwordHash,
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	mockCrypto.On("Compare", "wrong-password", mock.Anything).Return(false, nil)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "user@test.com", "wrong-password")

	assert.ErrorIs(t, err, msgerror.AnErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestLoginWithCompareError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	email, _ := vo.NewEmail("user@test.com")

	user := &entity.User{
		PasswordHash: passwordHash,
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	compareErr := errors.New("comparison failed")
	mockCrypto.On("Compare", "any-password", mock.Anything).Return(false, compareErr)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "user@test.com", "any-password")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to verify password")
	assert.Contains(t, err.Error(), "comparison failed")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestLoginSuccessfully(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	name, _ := vo.NewName("Test User", 0, 0)
	email, _ := vo.NewEmail("user@test.com")
	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	userID := vo.NewID()

	user := &entity.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	expectedClaims := providers.Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	mockCrypto.On("Compare", "valid-password", validHash).Return(true, nil)
	mockToken.On("Generate", expectedClaims).Return("generated_token", nil)
	mockBlacklist.On("Add", mock.Anything, "generated_token", 24*time.Hour).Return(nil)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	result, err := handler.Execute(context.Background(), "user@test.com", "valid-password")

	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.CreatedAt, result.CreatedAt)
	assert.Equal(t, "generated_token", result.Token)

	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	// mockBlacklist.AssertExpectations(t)
}

func TestLoginTokenGenerationFailure(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	name, _ := vo.NewName("Test User", 0, 0)
	email, _ := vo.NewEmail("user@test.com")
	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	userID := vo.NewID()

	user := &entity.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	expectedClaims := providers.Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	mockCrypto.On("Compare", "valid-password", validHash).Return(true, nil)
	mockToken.On("Generate", expectedClaims).Return("", errors.New("token generation error"))

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "user@test.com", "valid-password")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate token")
	assert.Contains(t, err.Error(), "token generation error")

	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	// mockBlacklist.AssertNotCalled(t, "Add")
}

// func TestLoginBlacklistAddFailure(t *testing.T) {
// 	mockRepo := new(mocks.MockUserRepo)
// 	mockCrypto := new(mocks.MockCrypto)
// 	mockToken := new(mocks.MockTokenProvider)
// 	mockBlacklist := new(mocks.MockBlacklist)

// 	name, _ := vo.NewName("Test User", 0, 0)
// 	email, _ := vo.NewEmail("user@test.com")
// 	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
// 	passwordHash, _ := vo.NewPasswordHash(validHash)
// 	userID := vo.NewID()

// 	user := &entity.User{
// 		ID:           userID,
// 		Name:         name,
// 		Email:        email,
// 		PasswordHash: passwordHash,
// 		CreatedAt:    time.Now(),
// 	}

// 	expectedClaims := providers.Claims{
// 		UserID: userID.String(),
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Subject:   user.Email.String(),
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
// 		},
// 	}

// 	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
// 	mockCrypto.On("Compare", "valid-password", validHash).Return(true, nil)
// 	mockToken.On("Generate", expectedClaims).Return("generated_token", nil)
// 	// mockBlacklist.On("Add", mock.Anything, "generated_token", 24*time.Hour).Return(errors.New("blacklist error"))

// 	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
// 	_, err := handler.Execute(context.Background(), "user@test.com", "valid-password")

// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "failed to secure session")
// 	assert.Contains(t, err.Error(), "blacklist error")

// 	mockRepo.AssertExpectations(t)
// 	mockCrypto.AssertExpectations(t)
// 	mockToken.AssertExpectations(t)
// 	// mockBlacklist.AssertExpectations(t)
// }

func TestLoginWithEmptyPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "valid@test.com", "")

	var valErr *msgerror.ValidationErrors
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "cannot be empty", valErr.FieldErrors["password"])

	mockRepo.AssertNotCalled(t, "GetByEmail")
	mockCrypto.AssertNotCalled(t, "Compare")
}
