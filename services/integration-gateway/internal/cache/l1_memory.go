package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
)

// L1MemoryCache implementa cache em memória com Ristretto (LFU eviction)
type L1MemoryCache struct {
	cache      *ristretto.Cache
	defaultTTL time.Duration
	stats      *L1Stats
}

// L1Stats estatísticas do cache L1
type L1Stats struct {
	Hits      uint64
	Misses    uint64
	Sets      uint64
	Evictions uint64
	Size      uint64
}

// L1Config configuração do cache L1
type L1Config struct {
	MaxSize    int64         // Tamanho máximo em bytes (default: 100MB)
	DefaultTTL time.Duration // TTL padrão (default: 5min)
	NumCounters int64        // Número de contadores para LFU (default: 10x MaxSize)
	BufferItems int64        // Buffer de items (default: 64)
}

// DefaultL1Config retorna configuração padrão do L1
func DefaultL1Config() L1Config {
	maxSize := int64(100 * 1024 * 1024) // 100MB
	return L1Config{
		MaxSize:     maxSize,
		DefaultTTL:  5 * time.Minute,
		NumCounters: maxSize / 10, // ~10x items esperados
		BufferItems: 64,
	}
}

// NewL1MemoryCache cria um novo cache L1 em memória
func NewL1MemoryCache(config L1Config) (*L1MemoryCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: config.NumCounters,
		MaxCost:     config.MaxSize,
		BufferItems: config.BufferItems,
		Metrics:     true, // Habilita métricas internas
		OnEvict: func(item *ristretto.Item) {
			// Callback quando item é removido (para métricas)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ristretto cache: %w", err)
	}

	return &L1MemoryCache{
		cache:      cache,
		defaultTTL: config.DefaultTTL,
		stats:      &L1Stats{},
	}, nil
}

// Get busca valor do cache L1
func (l *L1MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	value, found := l.cache.Get(key)

	if found {
		l.stats.Hits++
		return value, true
	}

	l.stats.Misses++
	return nil, false
}

// Set armazena valor no cache L1 com TTL padrão
func (l *L1MemoryCache) Set(ctx context.Context, key string, value interface{}, cost int64) error {
	return l.SetWithTTL(ctx, key, value, cost, l.defaultTTL)
}

// SetWithTTL armazena valor no cache L1 com TTL customizado
func (l *L1MemoryCache) SetWithTTL(ctx context.Context, key string, value interface{}, cost int64, ttl time.Duration) error {
	// Estima cost se não fornecido (1KB por item por padrão)
	if cost == 0 {
		cost = 1024
	}

	success := l.cache.SetWithTTL(key, value, cost, ttl)
	if !success {
		return fmt.Errorf("failed to set key %s in L1 cache", key)
	}

	// Aguarda escrita assíncrona completar (Ristretto é assíncrono)
	l.cache.Wait()

	l.stats.Sets++
	return nil
}

// Delete remove valor do cache L1
func (l *L1MemoryCache) Delete(ctx context.Context, key string) {
	l.cache.Del(key)
}

// Clear limpa todo o cache L1
func (l *L1MemoryCache) Clear(ctx context.Context) {
	l.cache.Clear()
	l.stats = &L1Stats{} // Reseta estatísticas
}

// GetStats retorna estatísticas atuais do cache
func (l *L1MemoryCache) GetStats() L1Stats {
	metrics := l.cache.Metrics

	return L1Stats{
		Hits:      metrics.Hits(),
		Misses:    metrics.Misses(),
		Sets:      l.stats.Sets,
		Evictions: metrics.KeysEvicted(),
		Size:      l.cache.Metrics.CostAdded() - l.cache.Metrics.CostEvicted(),
	}
}

// GetHitRate calcula a taxa de hit do cache (0.0 a 1.0)
func (l *L1MemoryCache) GetHitRate() float64 {
	stats := l.GetStats()
	total := stats.Hits + stats.Misses
	if total == 0 {
		return 0.0
	}
	return float64(stats.Hits) / float64(total)
}

// Close fecha o cache e libera recursos
func (l *L1MemoryCache) Close() {
	l.cache.Close()
}
