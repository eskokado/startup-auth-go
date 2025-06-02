package mocks

import "github.com/stretchr/testify/mock"

type MockCrypto struct {
	mock.Mock
}

func (m *MockCrypto) Encrypt(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockCrypto) Compare(password, hash string) (bool, error) {
	args := m.Called(password, hash)
	return args.Bool(0), args.Error(1)
}
