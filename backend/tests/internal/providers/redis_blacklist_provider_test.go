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

func (m *MockRedisCmdable) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisCmdable) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.SliceCmd)
}

func (m *MockRedisCmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisCmdable) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisCmdable) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

// ===== Testes existentes para Add e Exists =====

func TestRedisBlacklist_Add_Success(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"
	ttl := 5 * time.Minute

	cmd := redis.NewStatusResult("", nil)
	mockClient.On("Set", ctx, "startup-auth-go:"+token, token, ttl).Return(cmd)

	err := provider.Add(ctx, token, ttl)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Add_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"
	ttl := time.Duration(0)
	expectedErr := errors.New("redis error")

	cmd := redis.NewStatusResult("", expectedErr)
	mockClient.On("Set", ctx, "startup-auth-go:"+token, "", ttl).Return(cmd)

	err := provider.Add(ctx, token, 0)
	assert.ErrorIs(t, err, expectedErr)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Exists_True(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	token := "test_token"

	cmd := redis.NewIntResult(1, nil)
	mockClient.On("Exists", ctx, []string{"startup-auth-go:" + token}).Return(cmd)

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
	mockClient.On("Exists", ctx, []string{"startup-auth-go:" + token}).Return(cmd)

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
	mockClient.On("Exists", ctx, []string{"startup-auth-go:" + token}).Return(cmd)

	exists, err := provider.Exists(ctx, token)
	assert.ErrorIs(t, err, expectedErr)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}

// ===== Novos testes para ExistsKey =====

func TestRedisBlacklist_ExistsKey_True(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "some_key"

	cmd := redis.NewIntResult(1, nil)
	mockClient.On("Exists", ctx, []string{key}).Return(cmd)

	exists, err := provider.ExistsKey(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_ExistsKey_False(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "some_key"

	cmd := redis.NewIntResult(0, nil)
	mockClient.On("Exists", ctx, []string{key}).Return(cmd)

	exists, err := provider.ExistsKey(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_ExistsKey_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "some_key"
	expectedErr := errors.New("redis error")

	cmd := redis.NewIntResult(0, expectedErr)
	mockClient.On("Exists", ctx, []string{key}).Return(cmd)

	exists, err := provider.ExistsKey(ctx, key)
	assert.ErrorIs(t, err, expectedErr)
	assert.False(t, exists)
	mockClient.AssertExpectations(t)
}

// ===== Testes para SetWithKey =====

func TestRedisBlacklist_SetWithKey_Success(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "test_key"
	value := "test_value"
	ttl := 5 * time.Minute

	cmd := redis.NewStatusResult("", nil)
	mockClient.On("Set", ctx, key, value, ttl).Return(cmd)

	err := provider.SetWithKey(ctx, key, value, ttl)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_SetWithKey_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "test_key"
	value := "test_value"
	ttl := 5 * time.Minute
	expectedErr := errors.New("redis error")

	cmd := redis.NewStatusResult("", expectedErr)
	mockClient.On("Set", ctx, key, value, ttl).Return(cmd)

	err := provider.SetWithKey(ctx, key, value, ttl)
	assert.ErrorIs(t, err, expectedErr)
	mockClient.AssertExpectations(t)
}

// ===== Testes para Get =====

func TestRedisBlacklist_Get_Success(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "test_key"
	expectedValue := "test_value"

	cmd := redis.NewStringResult(expectedValue, nil)
	mockClient.On("Get", ctx, key).Return(cmd)

	value, err := provider.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Get_NotFound(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "test_key"

	cmd := redis.NewStringResult("", redis.Nil)
	mockClient.On("Get", ctx, key).Return(cmd)

	value, err := provider.Get(ctx, key)
	assert.NoError(t, err)
	assert.Empty(t, value)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Get_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	key := "test_key"
	expectedErr := errors.New("redis error")

	cmd := redis.NewStringResult("", expectedErr)
	mockClient.On("Get", ctx, key).Return(cmd)

	value, err := provider.Get(ctx, key)
	assert.ErrorIs(t, err, expectedErr)
	assert.Empty(t, value)
	mockClient.AssertExpectations(t)
}

// ===== Testes para MGet =====

func TestRedisBlacklist_MGet_Success(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	keys := []string{"key1", "key2"}
	expectedValues := []interface{}{"value1", "value2"}

	cmd := redis.NewSliceResult(expectedValues, nil)
	mockClient.On("MGet", ctx, keys).Return(cmd)

	values, err := provider.MGet(ctx, keys...)
	assert.NoError(t, err)
	assert.Equal(t, expectedValues, values)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_MGet_WithSomeNil(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	keys := []string{"key1", "key2"}
	expectedValues := []interface{}{"value1", nil}

	cmd := redis.NewSliceResult(expectedValues, nil)
	mockClient.On("MGet", ctx, keys).Return(cmd)

	values, err := provider.MGet(ctx, keys...)
	assert.NoError(t, err)
	assert.Equal(t, expectedValues, values)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_MGet_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	keys := []string{"key1", "key2"}
	expectedErr := errors.New("redis error")

	cmd := redis.NewSliceResult(nil, expectedErr)
	mockClient.On("MGet", ctx, keys).Return(cmd)

	values, err := provider.MGet(ctx, keys...)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, values)
	mockClient.AssertExpectations(t)
}

// ===== Testes para Del =====

func TestRedisBlacklist_Del_Success(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	keys := []string{"key1", "key2"}

	// Configurar para retornar 2 chaves deletadas
	cmd := redis.NewIntResult(2, nil)
	mockClient.On("Del", ctx, keys).Return(cmd)

	err := provider.Del(ctx, keys...)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Del_Error(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	keys := []string{"key1", "key2"}
	expectedErr := errors.New("redis error")

	cmd := redis.NewIntResult(0, expectedErr)
	mockClient.On("Del", ctx, keys).Return(cmd)

	err := provider.Del(ctx, keys...)
	assert.ErrorIs(t, err, expectedErr)
	mockClient.AssertExpectations(t)
}

func TestRedisBlacklist_Del_NoKeys(t *testing.T) {
	mockClient := new(MockRedisCmdable)
	provider := providers.NewRedisBlacklist(mockClient)

	ctx := context.Background()
	err := provider.Del(ctx) // Sem chaves

	assert.NoError(t, err)
	mockClient.AssertNotCalled(t, "Del")
}
