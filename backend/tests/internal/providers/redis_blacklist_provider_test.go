package providers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eskokado/startup-auth-go/backend/internal/providers"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisCmdable struct {
	mock.Mock
}

func (m *MockRedisCmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisCmdable) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func TestRedisBlacklist_Add_Success(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"
	ttl := 5 * time.Minute

	cmd := redis.NewStatusResult("", nil)
	mockClient.On("Set", ctx, "blacklist:"+token, true, ttl).Return(cmd)

	err := provider.Add(ctx, token, ttl)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Add_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"
	ttl := 5 * time.Minute
	expectedErr := errors.New("redis error")

	cmd := redis.NewStatusResult("", expectedErr)
	mockClient.On("Set", ctx, "blacklist:"+token, true, ttl).Return(cmd)

	err := provider.Add(ctx, token, ttl)
	assert.ErrorIs(t, err, expectedErr)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Exists_True(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"

	cmd := redis.NewIntResult(1, nil)
	mockClient.On("Exists", ctx, []string{"blacklist:" + token}).Return(cmd)

	exists, err := provider.Exists(ctx, token)
	assert.NoError(t, err)
	assert.True(t, exists)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Exists_False(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"

	cmd := redis.NewIntResult(0, nil)
	mockClient.On("Exists", ctx, []string{"blacklist:" + token}).Return(cmd)

	exists, err := provider.Exists(ctx, token)
	assert.NoError(t, err)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Exists_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"
	expectedErr := errors.New("redis error")

	cmd := redis.NewIntResult(0, expectedErr)
	mockClient.On("Exists", ctx, []string{"blacklist:" + token}).Return(cmd)

	exists, err := provider.Exists(ctx, token)
	assert.ErrorIs(t, err, expectedErr)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}
