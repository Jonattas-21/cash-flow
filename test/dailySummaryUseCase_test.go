package dailySummary

import (
	"github.com/Jonattas-21/cash-flow/internal/domain/entities"
	"github.com/Jonattas-21/cash-flow/internal/usecases"

	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Interface que define os métodos do Redis que vamos usar
type RedisClient interface {
	Get(key string) *redis.StringCmd
}

// Mocks
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	args := m.Called(key)
	return args.Get(0).(*redis.StringCmd)
}

// Test cases
func TestGetDailySummary(t *testing.T) {
	// Setup
	mockRedisClient := new(MockRedisClient)
	//useCase := &usecases.DailySummaryUseCase{Rdb: mockRedisClient}
	useCase := &usecases.DailySummaryUseCase{}


	testDate := time.Now()
	cacheKey := fmt.Sprintf("daily_summary:%s", testDate.Format("2006-01-02"))
	expectedSummary := &entities.DailySummary{Total: 100}

	t.Run("Cache hit", func(t *testing.T) {
		// Mocking Redis to return a valid JSON
		jsonSummary, _ := json.Marshal(expectedSummary)
		mockRedisClient.On("Get", cacheKey).Return(redis.NewStringResult(string(jsonSummary), nil))

		// Call the function
		summary, err := useCase.GetDailySummary(testDate)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedSummary, summary)
	})

	t.Run("Cache miss", func(t *testing.T) {
		// Mocking Redis to return redis.Nil (cache miss)
		mockRedisClient.On("Get", cacheKey).Return(redis.NewStringResult("", redis.Nil))

		// Call the function
		summary, err := useCase.GetDailySummary(testDate)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedSummary, summary)
	})

	t.Run("Redis error", func(t *testing.T) {
		// Mocking Redis to return an error
		mockRedisClient.On("Get", cacheKey).Return(redis.NewStringResult("", errors.New("redis error")))

		// Call the function
		_, err := useCase.GetDailySummary(testDate)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "redis error", err.Error())
	})

	t.Run("Invalid JSON in cache", func(t *testing.T) {
		// Mocking Redis to return invalid JSON
		mockRedisClient.On("Get", cacheKey).Return(redis.NewStringResult("invalid-json", nil))

		// Call the function
		_, err := useCase.GetDailySummary(testDate)

		// Assertions
		assert.Error(t, err)
	})
}
