package mocks

import (
	"github.com/stretchr/testify/mock"
	"gopkg.in/gomail.v2"
)

type MockSenderService struct {
	mock.Mock
}

func (m *MockSenderService) DialAndSend(msg ...*gomail.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}
