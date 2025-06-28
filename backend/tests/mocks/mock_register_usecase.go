package mocks

import (
	"context"

	"github.com/eskokado/startup-auth-go/backend/pkg/dto"
	"github.com/stretchr/testify/mock"
)

type MockRegisterUseCase struct {
	mock.Mock
}

func (m *MockRegisterUseCase) Execute(ctx context.Context, input dto.RegisterParams) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
