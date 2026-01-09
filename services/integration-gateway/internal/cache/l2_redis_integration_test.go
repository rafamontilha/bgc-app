// +build integration

package cache

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestL2RedisCache_Integration testa Redis em ambiente real
// Para executar: go test -tags=integration -v ./internal/cache/... -run TestL2RedisCache_Integration
func TestL2RedisCache_Integration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	config := L2Config{
		Addr:       redisAddr,
		Password:   "",
		DB:         0,
		DefaultTTL: 10 * time.Minute,
		Prefix:     "test:cache:",
		MaxRetries: 3,
		PoolSize:   10,
	}

	cache, err := NewL2RedisCache(config)
	if err != nil {
		t.Skipf("Redis not available at %s: %v", redisAddr, err)
		return
	}
	defer cache.Close()

	ctx := context.Background()

	// Limpa antes de testar
	cache.Clear(ctx)

	t.Run("SetAndGet", func(t *testing.T) {
		err := cache.Set(ctx, "test-key", "test-value")
		require.NoError(t, err)

		value, found, err := cache.Get(ctx, "test-key")
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "test-value", value)
	})

	t.Run("GetMiss", func(t *testing.T) {
		value, found, err := cache.Get(ctx, "nonexistent-key")
		require.NoError(t, err)
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("SetWithTTL", func(t *testing.T) {
		err := cache.SetWithTTL(ctx, "ttl-key", "ttl-value", 2*time.Second)
		require.NoError(t, err)

		// Verifica existência
		value, found, err := cache.Get(ctx, "ttl-key")
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "ttl-value", value)

		// Verifica TTL
		ttl, err := cache.GetTTL(ctx, "ttl-key")
		require.NoError(t, err)
		assert.Greater(t, ttl, time.Duration(0))
		assert.LessOrEqual(t, ttl, 2*time.Second)

		// Aguarda expiração
		time.Sleep(3 * time.Second)

		// Verifica que expirou
		value, found, err = cache.Get(ctx, "ttl-key")
		require.NoError(t, err)
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("Delete", func(t *testing.T) {
		cache.Set(ctx, "delete-key", "delete-value")

		err := cache.Delete(ctx, "delete-key")
		require.NoError(t, err)

		value, found, err := cache.Get(ctx, "delete-key")
		require.NoError(t, err)
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("Exists", func(t *testing.T) {
		cache.Set(ctx, "exists-key", "exists-value")

		exists, err := cache.Exists(ctx, "exists-key")
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = cache.Exists(ctx, "nonexistent-key")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Clear", func(t *testing.T) {
		// Set múltiplos valores
		cache.Set(ctx, "clear1", "value1")
		cache.Set(ctx, "clear2", "value2")
		cache.Set(ctx, "clear3", "value3")

		// Clear
		err := cache.Clear(ctx)
		require.NoError(t, err)

		// Verifica remoção
		value, found, _ := cache.Get(ctx, "clear1")
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("ComplexDataTypes", func(t *testing.T) {
		// Map
		testMap := map[string]interface{}{
			"name":  "John",
			"age":   30,
			"email": "john@example.com",
		}
		cache.Set(ctx, "map-key", testMap)

		value, found, err := cache.Get(ctx, "map-key")
		require.NoError(t, err)
		assert.True(t, found)

		resultMap, ok := value.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "John", resultMap["name"])
		assert.Equal(t, float64(30), resultMap["age"]) // JSON deserializa números como float64
	})

	t.Run("Ping", func(t *testing.T) {
		err := cache.Ping(ctx)
		require.NoError(t, err)
	})

	t.Run("GetStats", func(t *testing.T) {
		stats, err := cache.GetStats(ctx)
		require.NoError(t, err)
		require.NotNil(t, stats)

		assert.GreaterOrEqual(t, stats.TotalConns, uint32(0))
	})
}

func TestMultiLevelCacheManager_FullIntegration(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Testa conectividade Redis
	testConfig := L2Config{
		Addr:       redisAddr,
		Password:   "",
		DB:         0,
		DefaultTTL: 10 * time.Minute,
		Prefix:     "test:full:",
		MaxRetries: 3,
		PoolSize:   10,
	}

	testRedis, err := NewL2RedisCache(testConfig)
	if err != nil {
		t.Skipf("Redis not available at %s: %v", redisAddr, err)
		return
	}
	testRedis.Close()

	// Configura manager com L1 e L2
	config := ManagerConfig{
		L1Config: DefaultL1Config(),
		L2Config: L2Config{
			Addr:       redisAddr,
			Password:   "",
			DB:         0,
			DefaultTTL: 10 * time.Minute,
			Prefix:     "test:manager:",
			MaxRetries: 3,
			PoolSize:   10,
		},
		EnableL1:     true,
		EnableL2:     true,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	ctx := context.Background()

	// Limpa antes de testar
	manager.Clear(ctx)

	t.Run("CascadeL1toL2", func(t *testing.T) {
		// Pre-popula L2 diretamente
		manager.l2.Set(ctx, "l2-only-key", "l2-only-value")

		// Get deve buscar em cascata: L1 miss → L2 hit
		value, level, err := manager.Get(ctx, "l2-only-key")
		require.NoError(t, err)
		assert.Equal(t, "l2-only-value", value)
		assert.Equal(t, LevelL2, level)

		// Próxima chamada deve vir de L1 (promoção automática)
		value, level, err = manager.Get(ctx, "l2-only-key")
		require.NoError(t, err)
		assert.Equal(t, "l2-only-value", value)
		assert.Equal(t, LevelL1, level)
	})

	t.Run("SetPropagation", func(t *testing.T) {
		// Set deve propagar para L1 e L2
		err := manager.Set(ctx, "propagate-key", "propagate-value", 5*time.Minute)
		require.NoError(t, err)

		// Verifica L1 (limpa L1 primeiro para forçar busca em L2)
		manager.l1.Delete(ctx, "propagate-key")

		value, level, err := manager.Get(ctx, "propagate-key")
		require.NoError(t, err)
		assert.Equal(t, "propagate-value", value)
		assert.Equal(t, LevelL2, level) // Vem de L2
	})

	t.Run("DeletePropagation", func(t *testing.T) {
		// Set em todos os níveis
		manager.Set(ctx, "delete-key", "delete-value", 5*time.Minute)

		// Verifica que existe
		value, level, _ := manager.Get(ctx, "delete-key")
		assert.Equal(t, "delete-value", value)
		assert.Equal(t, LevelL1, level)

		// Delete
		err := manager.Delete(ctx, "delete-key")
		require.NoError(t, err)

		// Verifica remoção de L2 diretamente
		value, found, err := manager.l2.Get(ctx, "delete-key")
		require.NoError(t, err)
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("HighThroughput", func(t *testing.T) {
		// Testa alta throughput (100 operações)
		for i := 0; i < 100; i++ {
			key := string(rune('a' + i%26))
			manager.Set(ctx, key, i, 5*time.Minute)
		}

		// Verifica hit rate
		stats := manager.GetStats(ctx)
		require.NotNil(t, stats)

		l1Stats := stats["l1"].(map[string]interface{})
		hits := l1Stats["hits"].(uint64)
		misses := l1Stats["misses"].(uint64)
		total := hits + misses

		if total > 0 {
			hitRate := float64(hits) / float64(total)
			t.Logf("L1 Hit Rate: %.2f%% (hits: %d, misses: %d)", hitRate*100, hits, misses)
		}
	})
}
