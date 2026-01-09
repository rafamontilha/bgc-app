package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewL1MemoryCache(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)

	require.NoError(t, err)
	require.NotNil(t, cache)
	assert.NotNil(t, cache.cache)
	assert.Equal(t, config.DefaultTTL, cache.defaultTTL)

	cache.Close()
}

func TestL1MemoryCache_SetAndGet(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set valor
	err = cache.Set(ctx, "test-key", "test-value", 1024)
	require.NoError(t, err)

	// Get valor
	value, found := cache.Get(ctx, "test-key")
	assert.True(t, found)
	assert.Equal(t, "test-value", value)
}

func TestL1MemoryCache_GetMiss(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Get chave inexistente
	value, found := cache.Get(ctx, "nonexistent-key")
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestL1MemoryCache_SetWithTTL(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set com TTL curto (100ms)
	err = cache.SetWithTTL(ctx, "ttl-key", "ttl-value", 1024, 100*time.Millisecond)
	require.NoError(t, err)

	// Get imediatamente (deve existir)
	value, found := cache.Get(ctx, "ttl-key")
	assert.True(t, found)
	assert.Equal(t, "ttl-value", value)

	// Aguarda TTL expirar
	time.Sleep(150 * time.Millisecond)

	// Get após expiração (deve ter sumido)
	value, found = cache.Get(ctx, "ttl-key")
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestL1MemoryCache_Delete(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set e confirma existência
	err = cache.Set(ctx, "delete-key", "delete-value", 1024)
	require.NoError(t, err)

	value, found := cache.Get(ctx, "delete-key")
	assert.True(t, found)
	assert.Equal(t, "delete-value", value)

	// Delete
	cache.Delete(ctx, "delete-key")

	// Confirma remoção
	value, found = cache.Get(ctx, "delete-key")
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestL1MemoryCache_Clear(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set múltiplas chaves
	cache.Set(ctx, "key1", "value1", 1024)
	cache.Set(ctx, "key2", "value2", 1024)
	cache.Set(ctx, "key3", "value3", 1024)

	// Confirma existência
	_, found1 := cache.Get(ctx, "key1")
	_, found2 := cache.Get(ctx, "key2")
	_, found3 := cache.Get(ctx, "key3")
	assert.True(t, found1)
	assert.True(t, found2)
	assert.True(t, found3)

	// Clear
	cache.Clear(ctx)

	// Confirma remoção de todas
	_, found1 = cache.Get(ctx, "key1")
	_, found2 = cache.Get(ctx, "key2")
	_, found3 = cache.Get(ctx, "key3")
	assert.False(t, found1)
	assert.False(t, found2)
	assert.False(t, found3)
}

func TestL1MemoryCache_GetStats(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Operações para gerar estatísticas
	cache.Set(ctx, "key1", "value1", 1024)
	cache.Set(ctx, "key2", "value2", 1024)

	cache.Get(ctx, "key1")          // Hit
	cache.Get(ctx, "key2")          // Hit
	cache.Get(ctx, "nonexistent")   // Miss

	// Verifica stats
	stats := cache.GetStats()
	assert.Equal(t, uint64(2), stats.Hits)
	assert.Equal(t, uint64(1), stats.Misses)
	assert.Equal(t, uint64(2), stats.Sets)
	assert.True(t, stats.Size > 0)
}

func TestL1MemoryCache_GetHitRate(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Nenhuma operação = hit rate 0
	hitRate := cache.GetHitRate()
	assert.Equal(t, 0.0, hitRate)

	// Set e Get (2 hits, 0 misses)
	cache.Set(ctx, "key1", "value1", 1024)
	cache.Get(ctx, "key1") // Hit
	cache.Get(ctx, "key1") // Hit

	hitRate = cache.GetHitRate()
	assert.Equal(t, 1.0, hitRate) // 100% hit rate

	// 1 Miss
	cache.Get(ctx, "nonexistent") // Miss

	hitRate = cache.GetHitRate()
	assert.InDelta(t, 0.666, hitRate, 0.01) // ~66.6% hit rate (2 hits, 1 miss)
}

func TestL1MemoryCache_MultipleTypes(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// String
	cache.Set(ctx, "string", "hello", 1024)
	value, found := cache.Get(ctx, "string")
	assert.True(t, found)
	assert.Equal(t, "hello", value)

	// Int
	cache.Set(ctx, "int", 42, 1024)
	value, found = cache.Get(ctx, "int")
	assert.True(t, found)
	assert.Equal(t, 42, value)

	// Struct
	type TestStruct struct {
		Name  string
		Value int
	}
	testStruct := TestStruct{Name: "test", Value: 123}
	cache.Set(ctx, "struct", testStruct, 1024)
	value, found = cache.Get(ctx, "struct")
	assert.True(t, found)
	assert.Equal(t, testStruct, value)
}

func TestL1MemoryCache_ConcurrentAccess(t *testing.T) {
	config := DefaultL1Config()
	cache, err := NewL1MemoryCache(config)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Acesso concorrente (simula múltiplas goroutines)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			key := "concurrent-key"
			cache.Set(ctx, key, id, 1024)
			cache.Get(ctx, key)
			done <- true
		}(i)
	}

	// Aguarda todas as goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Cache deve estar funcional
	stats := cache.GetStats()
	assert.True(t, stats.Hits+stats.Misses > 0)
}

func BenchmarkL1MemoryCache_Set(b *testing.B) {
	config := DefaultL1Config()
	cache, _ := NewL1MemoryCache(config)
	defer cache.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(ctx, "bench-key", "bench-value", 1024)
	}
}

func BenchmarkL1MemoryCache_Get(b *testing.B) {
	config := DefaultL1Config()
	cache, _ := NewL1MemoryCache(config)
	defer cache.Close()

	ctx := context.Background()
	cache.Set(ctx, "bench-key", "bench-value", 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(ctx, "bench-key")
	}
}
