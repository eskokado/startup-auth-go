package usecase_test

import (
	"context"
	"errors"
	"testing"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogoutSuccess(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	logoutUsecase := usecase.NewLogoutUsecase(mockBlacklist)

	ctx := context.Background()
	token := "valid_token"

	// Definir as chaves que serão deletadas
	expectedKeys := []string{
		"startup-auth-go:valid_token:UserID",
		"startup-auth-go:valid_token:Name",
		"startup-auth-go:valid_token:Email",
		"startup-auth-go:valid_token:Token",
		"startup-auth-go:valid_token:CreatedAt",
	}

	// Configurar o mock para retornar sucesso
	mockBlacklist.On("Del", ctx, expectedKeys).Return(nil)

	err := logoutUsecase.Execute(ctx, token)

	assert.NoError(t, err)
	mockBlacklist.AssertExpectations(t)
}

func TestLogoutError(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	logoutUsecase := usecase.NewLogoutUsecase(mockBlacklist)

	ctx := context.Background()
	token := "valid_token"
	expectedErr := errors.New("redis error")

	expectedKeys := []string{
		"startup-auth-go:valid_token:UserID",
		"startup-auth-go:valid_token:Name",
		"startup-auth-go:valid_token:Email",
		"startup-auth-go:valid_token:Token",
		"startup-auth-go:valid_token:CreatedAt",
	}

	mockBlacklist.On("Del", ctx, expectedKeys).Return(expectedErr)

	err := logoutUsecase.Execute(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to remove user session data")
	assert.Contains(t, err.Error(), "redis error")
	mockBlacklist.AssertExpectations(t)
}

func TestLogoutEmptyToken(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	logoutUsecase := usecase.NewLogoutUsecase(mockBlacklist)

	ctx := context.Background()
	err := logoutUsecase.Execute(ctx, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is required")
	mockBlacklist.AssertNotCalled(t, "Del")
}

func TestLogoutKeysDefinition(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	logoutUsecase := usecase.NewLogoutUsecase(mockBlacklist)

	ctx := context.Background()
	token := "token123"

	expectedKeys := []string{
		"startup-auth-go:token123:UserID",
		"startup-auth-go:token123:Name",
		"startup-auth-go:token123:Email",
		"startup-auth-go:token123:Token",
		"startup-auth-go:token123:CreatedAt",
	}

	mockBlacklist.On("Del", ctx, expectedKeys).Return(nil)

	_ = logoutUsecase.Execute(ctx, token)

	// Verifica se as chaves passadas são exatamente as esperadas
	mockBlacklist.AssertCalled(t, "Del", ctx, expectedKeys)
}
