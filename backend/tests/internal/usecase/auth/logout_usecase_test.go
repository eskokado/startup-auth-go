package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	"github.com/eskokado/startup-auth-go/backend/pkg/msgerror"
	"github.com/eskokado/startup-auth-go/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutSuccess(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "valid_token"

	mockBlacklist.On("Exists", mock.Anything, token).Return(true, nil)
	mockBlacklist.On("Add", mock.Anything, token, time.Duration(0)).Return(nil)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.NoError(t, err)
	mockBlacklist.AssertExpectations(t)
}

func TestLogoutTokenNotFound(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "invalid_token"

	mockBlacklist.On("Exists", mock.Anything, token).Return(false, nil)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.ErrorIs(t, err, msgerror.AnErrInvalidToken)
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
	assert.Contains(t, err.Error(), "failed to verify token status")
	assert.ErrorIs(t, err, expectedErr)
	mockBlacklist.AssertExpectations(t)
	mockBlacklist.AssertNotCalled(t, "Add")
}

func TestLogoutAddError(t *testing.T) {
	mockBlacklist := new(mocks.MockBlacklist)
	token := "valid_token"
	expectedErr := errors.New("redis error")

	mockBlacklist.On("Exists", mock.Anything, token).Return(true, nil)
	mockBlacklist.On("Add", mock.Anything, token, time.Duration(0)).Return(expectedErr)

	uc := usecase.NewLogoutUsecase(mockBlacklist)
	err := uc.Execute(context.Background(), token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to revoke token")
	assert.ErrorIs(t, err, expectedErr)
	mockBlacklist.AssertExpectations(t)
}
