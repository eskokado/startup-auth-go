package mocks

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/stretchr/testify/mock"
)

type MockLoginUseCase struct {
	mock.Mock
}

func (m *MockLoginUseCase) Execute(ctx context.Context, email string, password string) (dto.LoginResult, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(dto.LoginResult), args.Error(1)
}
