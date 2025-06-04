package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutSuccess(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "valid_token"

	// Mock: token não está na blacklist
	mockBlacklist.On("Exists", mock.Anything, token).Return(false, nil)
	// Mock: adição bem sucedida
	mockBlacklist.On("Add", mock.Anything, token, 24*time.Hour).Return(nil)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.NoError(t, err)
	mockBlacklist.AssertExpectations(t)
}

func TestLogoutTokenAlreadyRevoked(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "revoked_token"

	// Mock: token já está na blacklist
	mockBlacklist.On("Exists", mock.Anything, token).Return(true, nil)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.Equal(t, "token already revoked", err.Error())
	mockBlacklist.AssertExpectations(t)
	mockBlacklist.AssertNotCalled(t, "Add")
}

func TestLogoutExistsError(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "any_token"
	expectedErr := errors.New("database error")

	// Mock: erro na verificação
	mockBlacklist.On("Exists", mock.Anything, token).Return(false, expectedErr)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.Equal(t, "failed to verify token status", err.Error())
	mockBlacklist.AssertExpectations(t)
	mockBlacklist.AssertNotCalled(t, "Add")
}

func TestLogoutAddError(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "valid_token"
	expectedErr := errors.New("redis error")

	// Mock: token não está na blacklist
	mockBlacklist.On("Exists", mock.Anything, token).Return(false, nil)
	// Mock: erro ao adicionar
	mockBlacklist.On("Add", mock.Anything, token, 24*time.Hour).Return(expectedErr)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err) // Erro original é propagado
	mockBlacklist.AssertExpectations(t)
}
