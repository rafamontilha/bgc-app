package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockL3Cache implementa L3Cache para testes
type MockL3Cache struct {
	data  map[string]interface{}
	mu    sync.RWMutex
	calls map[string]int // Rastreia chamadas para verificar promoção
}

func NewMockL3Cache() *MockL3Cache {
	return &MockL3Cache{
		data:  make(map[string]interface{}),
		calls: make(map[string]int),
	}
}

func (m *MockL3Cache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.calls["get"]++
	value, exists := m.data[key]
	return value, exists, nil
}

func (m *MockL3Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls["set"]++
	m.data[key] = value
	return nil
}

func (m *MockL3Cache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls["delete"]++
	delete(m.data, key)
	return nil
}

func (m *MockL3Cache) GetCalls(operation string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[operation]
}

func TestNewMultiLevelCacheManager(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false, // Redis não disponível em teste unitário
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	require.NotNil(t, manager)
	assert.NotNil(t, manager.l1)
	assert.Nil(t, manager.l2)
	assert.True(t, manager.enableL1)
	assert.False(t, manager.enableL2)

	manager.Close()
}

func TestMultiLevelCacheManager_L1Only(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	ctx := context.Background()

	// Set valor
	err = manager.Set(ctx, "test-key", "test-value", 5*time.Minute)
	require.NoError(t, err)

	// Get deve vir do L1
	value, level, err := manager.Get(ctx, "test-key")
	require.NoError(t, err)
	assert.Equal(t, "test-value", value)
	assert.Equal(t, LevelL1, level)
}

func TestMultiLevelCacheManager_L1_L3_Cascade(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     true,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	// Injeta mock L3
	mockL3 := NewMockL3Cache()
	manager.SetL3Cache(mockL3)

	ctx := context.Background()

	// Pre-popula L3 (simula dado já em cache L3)
	mockL3.Set(ctx, "l3-key", "l3-value", 10*time.Minute)

	// Get deve buscar em cascata: L1 miss → L3 hit
	value, level, err := manager.Get(ctx, "l3-key")
	require.NoError(t, err)
	assert.Equal(t, "l3-value", value)
	assert.Equal(t, LevelL3, level)

	// Verifica que L3 foi chamado
	assert.Equal(t, 1, mockL3.GetCalls("get"))

	// Próxima chamada deve vir de L1 (promoção automática)
	value, level, err = manager.Get(ctx, "l3-key")
	require.NoError(t, err)
	assert.Equal(t, "l3-value", value)
	assert.Equal(t, LevelL1, level) // Agora está em L1

	// L3 não deve ser chamado novamente
	assert.Equal(t, 1, mockL3.GetCalls("get"))
}

func TestMultiLevelCacheManager_CacheMiss(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     true,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	mockL3 := NewMockL3Cache()
	manager.SetL3Cache(mockL3)

	ctx := context.Background()

	// Get em chave inexistente (miss em todos os níveis)
	value, level, err := manager.Get(ctx, "nonexistent-key")
	require.NoError(t, err)
	assert.Nil(t, value)
	assert.Equal(t, LevelExternal, level) // Indica que deve buscar API externa

	// L3 deve ter sido consultado
	assert.Equal(t, 1, mockL3.GetCalls("get"))
}

func TestMultiLevelCacheManager_SetPropagation(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     true,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	mockL3 := NewMockL3Cache()
	manager.SetL3Cache(mockL3)

	ctx := context.Background()

	// Set deve propagar para todos os níveis
	err = manager.Set(ctx, "propagate-key", "propagate-value", 5*time.Minute)
	require.NoError(t, err)

	// Verifica L1
	value, level, err := manager.Get(ctx, "propagate-key")
	require.NoError(t, err)
	assert.Equal(t, "propagate-value", value)
	assert.Equal(t, LevelL1, level)

	// Verifica L3
	assert.Equal(t, 1, mockL3.GetCalls("set"))
	value, found, err := mockL3.Get(ctx, "propagate-key")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "propagate-value", value)
}

func TestMultiLevelCacheManager_Delete(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     true,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	mockL3 := NewMockL3Cache()
	manager.SetL3Cache(mockL3)

	ctx := context.Background()

	// Set valor em todos os níveis
	manager.Set(ctx, "delete-key", "delete-value", 5*time.Minute)

	// Verifica que existe
	value, level, _ := manager.Get(ctx, "delete-key")
	assert.Equal(t, "delete-value", value)
	assert.Equal(t, LevelL1, level)

	// Delete
	err = manager.Delete(ctx, "delete-key")
	require.NoError(t, err)

	// Verifica remoção de L1
	value, level, err = manager.Get(ctx, "delete-key")
	require.NoError(t, err)
	assert.Nil(t, value)
	assert.Equal(t, LevelExternal, level)

	// Verifica que L3 foi chamado para delete
	assert.Equal(t, 1, mockL3.GetCalls("delete"))
}

func TestMultiLevelCacheManager_Clear(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	ctx := context.Background()

	// Set múltiplos valores
	manager.Set(ctx, "key1", "value1", 5*time.Minute)
	manager.Set(ctx, "key2", "value2", 5*time.Minute)
	manager.Set(ctx, "key3", "value3", 5*time.Minute)

	// Clear
	err = manager.Clear(ctx)
	require.NoError(t, err)

	// Verifica que todos foram removidos
	value, level, _ := manager.Get(ctx, "key1")
	assert.Nil(t, value)
	assert.Equal(t, LevelExternal, level)

	value, level, _ = manager.Get(ctx, "key2")
	assert.Nil(t, value)
	assert.Equal(t, LevelExternal, level)

	value, level, _ = manager.Get(ctx, "key3")
	assert.Nil(t, value)
	assert.Equal(t, LevelExternal, level)
}

func TestMultiLevelCacheManager_GetStats(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	ctx := context.Background()

	// Gera algumas operações
	manager.Set(ctx, "key1", "value1", 5*time.Minute)
	manager.Get(ctx, "key1")          // Hit
	manager.Get(ctx, "nonexistent")   // Miss

	// Get stats
	stats := manager.GetStats(ctx)
	require.NotNil(t, stats)

	l1Stats, exists := stats["l1"]
	assert.True(t, exists)
	assert.NotNil(t, l1Stats)

	l1Map := l1Stats.(map[string]interface{})
	assert.Greater(t, l1Map["hits"].(uint64), uint64(0))
	assert.Greater(t, l1Map["misses"].(uint64), uint64(0))
	assert.Greater(t, l1Map["sets"].(uint64), uint64(0))
}

func TestMultiLevelCacheManager_ConcurrentAccess(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	ctx := context.Background()

	// Acesso concorrente
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := "concurrent-key"
			manager.Set(ctx, key, id, 5*time.Minute)
			manager.Get(ctx, key)
		}(i)
	}

	wg.Wait()

	// Manager deve estar funcional
	stats := manager.GetStats(ctx)
	assert.NotNil(t, stats)
}

func TestMultiLevelCacheManager_TTLRespected(t *testing.T) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, err := NewMultiLevelCacheManager(config)
	require.NoError(t, err)
	defer manager.Close()

	ctx := context.Background()

	// Set com TTL curto (100ms)
	err = manager.Set(ctx, "ttl-key", "ttl-value", 100*time.Millisecond)
	require.NoError(t, err)

	// Get imediatamente (deve existir)
	value, level, err := manager.Get(ctx, "ttl-key")
	require.NoError(t, err)
	assert.Equal(t, "ttl-value", value)
	assert.Equal(t, LevelL1, level)

	// Aguarda TTL expirar
	time.Sleep(150 * time.Millisecond)

	// Get após expiração (deve ter sumido)
	value, level, err = manager.Get(ctx, "ttl-key")
	require.NoError(t, err)
	assert.Nil(t, value)
	assert.Equal(t, LevelExternal, level)
}

func BenchmarkMultiLevelCacheManager_Get(b *testing.B) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, _ := NewMultiLevelCacheManager(config)
	defer manager.Close()

	ctx := context.Background()
	manager.Set(ctx, "bench-key", "bench-value", 5*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Get(ctx, "bench-key")
	}
}

func BenchmarkMultiLevelCacheManager_Set(b *testing.B) {
	config := ManagerConfig{
		L1Config:     DefaultL1Config(),
		L2Config:     DefaultL2Config(),
		EnableL1:     true,
		EnableL2:     false,
		EnableL3:     false,
		ConnectorID:  "test-connector",
		EndpointName: "test-endpoint",
	}

	manager, _ := NewMultiLevelCacheManager(config)
	defer manager.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Set(ctx, "bench-key", "bench-value", 5*time.Minute)
	}
}
