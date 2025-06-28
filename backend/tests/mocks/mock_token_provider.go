package mocks

import "github.com/stretchr/testify/mock"

type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) Generate(claims interface{}) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) Validate(token string) (interface{}, error) {
	args := m.Called(token)
	return args.Get(0), args.Error(1)
}
