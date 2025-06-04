package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockBlacklist struct {
	mock.Mock
}

func (m *MockBlacklist) Add(ctx context.Context, token string, ttl time.Duration) error {
	args := m.Called(ctx, token, ttl)
	return args.Error(0)
}

func (m *MockBlacklist) Exists(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}
