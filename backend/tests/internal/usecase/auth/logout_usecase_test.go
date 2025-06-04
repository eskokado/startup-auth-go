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

	// Token existe na blacklist
	mockBlacklist.On("Exists", mock.Anything, token).Return(true, nil)

	// Atualiza com TTL=0
	mockBlacklist.On("Add", mock.Anything, token, time.Duration(0)).Return(nil)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.NoError(t, err)
	mockBlacklist.AssertExpectations(t)
}

func TestLogoutTokenNotFound(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "invalid_token"

	// Token não existe na blacklist
	mockBlacklist.On("Exists", mock.Anything, token).Return(false, nil)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.Equal(t, "token not found", err.Error())
	mockBlacklist.AssertExpectations(t)
	mockBlacklist.AssertNotCalled(t, "Add")
}

func TestLogoutExistsError(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "any_token"
	expectedErr := errors.New("database error")

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

	// Token existe na blacklist
	mockBlacklist.On("Exists", mock.Anything, token).Return(true, nil)

	// Erro ao atualizar com TTL=0
	mockBlacklist.On("Add", mock.Anything, token, time.Duration(0)).Return(expectedErr)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockBlacklist.AssertExpectations(t)
}
