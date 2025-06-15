package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockResetPasswordUseCase struct {
	mock.Mock
}

func (m *MockResetPasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}
