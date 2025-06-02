package mocks

import (
	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"github.com/stretchr/testify/mock"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendResetPasswordEmail(email vo.Email, token string) error {
	args := m.Called(email, token)
	return args.Error(0)
}
