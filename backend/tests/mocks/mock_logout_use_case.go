package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLogoutUseCase struct {
	mock.Mock
}

func (m *MockLogoutUseCase) Execute(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
