package mocks

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/stretchr/testify/mock"
)

type MockForgotPasswordUseCase struct {
	mock.Mock
}

func (m *MockForgotPasswordUseCase) Execute(ctx context.Context, email vo.Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}
