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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
}

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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	mockBlacklist.AssertNotCalled(t, "SetWithKey")
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
	createdAt := time.Now()

	user := &entity.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
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

	// Expectativas para salvamento no Redis
	prefix := "startup-auth-go"
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":UserID", userID.String(), 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":Name", name.String(), 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":Email", email.String(), 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":Token", "generated_token", 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":CreatedAt", createdAt.Format(time.RFC3339), 24*time.Hour).Return(nil)

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
	mockBlacklist.AssertExpectations(t)
}

func TestLoginBlacklistAddFailure_UserID(t *testing.T) {
	testRedisSaveFailure(t, "UserID", "failed to save UserID")
}

func TestLoginBlacklistAddFailure_Name(t *testing.T) {
	testRedisSaveFailure(t, "Name", "failed to save Name")
}

func TestLoginBlacklistAddFailure_Email(t *testing.T) {
	testRedisSaveFailure(t, "Email", "failed to save Email")
}

func TestLoginBlacklistAddFailure_Token(t *testing.T) {
	testRedisSaveFailure(t, "Token", "failed to save Token")
}

func TestLoginBlacklistAddFailure_CreatedAt(t *testing.T) {
	testRedisSaveFailure(t, "CreatedAt", "failed to save CreatedAt")
}

func testRedisSaveFailure(t *testing.T, field string, expectedError string) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	name, _ := vo.NewName("Test User", 0, 0)
	email, _ := vo.NewEmail("user@test.com")
	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	userID := vo.NewID()
	createdAt := time.Now()

	user := &entity.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
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

	// Configurar as expectativas para o Redis
	prefix := "startup-auth-go"
	redisErr := errors.New("redis error")

	// Ordem dos campos conforme execução
	fieldsOrder := []string{"UserID", "Name", "Email", "Token", "CreatedAt"}
	calls := []*mock.Call{}

	for _, f := range fieldsOrder {
		key := prefix + ":" + f
		if f == field {
			call := mockBlacklist.On("SetWithKey", mock.Anything, key, mock.Anything, 24*time.Hour)
			call.Return(redisErr)
			calls = append(calls, call)
			break // Para após o campo que causa erro
		} else {
			call := mockBlacklist.On("SetWithKey", mock.Anything, key, mock.Anything, 24*time.Hour)
			call.Return(nil)
			calls = append(calls, call)
		}
	}

	// Garantir que as chamadas são verificadas na ordem
	for _, call := range calls {
		call.Once()
	}

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	_, err := handler.Execute(context.Background(), "user@test.com", "valid-password")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError)
	assert.Contains(t, err.Error(), "redis error")

	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
}

func TestLoginWithCreatedAtFormat(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	mockCrypto := new(mocks.MockCrypto)
	mockToken := new(mocks.MockTokenProvider)
	mockBlacklist := new(mocks.MockBlacklist)

	name, _ := vo.NewName("Test User", 0, 0)
	email, _ := vo.NewEmail("user@test.com")
	validHash := "$2a$10$0MwrQkGO0Bw6dYpVfiX4mefEVgTdgtCYCJ7LxltXfzj5qscr4sive"
	passwordHash, _ := vo.NewPasswordHash(validHash)
	userID := vo.NewID()

	// Criar um tempo específico para testar a formatação
	testTime := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	expectedFormat := testTime.Format(time.RFC3339)

	user := &entity.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    testTime,
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

	// Expectativas para salvamento no Redis
	prefix := "startup-auth-go"
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":UserID", userID.String(), 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":Name", name.String(), 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":Email", email.String(), 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":Token", "generated_token", 24*time.Hour).Return(nil)
	mockBlacklist.On("SetWithKey", mock.Anything, prefix+":CreatedAt", expectedFormat, 24*time.Hour).Return(nil)

	handler := usecase.NewLoginUsecase(mockRepo, mockCrypto, mockToken, mockBlacklist)
	result, err := handler.Execute(context.Background(), "user@test.com", "valid-password")

	assert.NoError(t, err)
	assert.Equal(t, testTime, result.CreatedAt)
	mockBlacklist.AssertCalled(t, "SetWithKey", mock.Anything, prefix+":CreatedAt", expectedFormat, 24*time.Hour)
}
