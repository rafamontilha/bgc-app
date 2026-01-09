package cache

import (
	"context"
	"fmt"
	"time"
)

// CacheLevel representa um nível de cache
type CacheLevel string

const (
	LevelL1       CacheLevel = "l1"
	LevelL2       CacheLevel = "l2"
	LevelL3       CacheLevel = "l3"
	LevelExternal CacheLevel = "external"
)

// MultiLevelCacheManager gerencia cache em múltiplos níveis (L1 → L2 → L3 → External API)
type MultiLevelCacheManager struct {
	l1              *L1MemoryCache
	l2              *L2RedisCache
	l3              L3Cache // Interface para PostgreSQL (implementar depois)
	enableL1        bool
	enableL2        bool
	enableL3        bool
	connectorID     string
	endpointName    string
}

// L3Cache interface para cache L3 (PostgreSQL materialized views)
type L3Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// ManagerConfig configuração do manager
type ManagerConfig struct {
	L1Config     L1Config
	L2Config     L2Config
	EnableL1     bool
	EnableL2     bool
	EnableL3     bool
	ConnectorID  string
	EndpointName string
}

// NewMultiLevelCacheManager cria um novo gerenciador de cache multinível
func NewMultiLevelCacheManager(config ManagerConfig) (*MultiLevelCacheManager, error) {
	manager := &MultiLevelCacheManager{
		enableL1:     config.EnableL1,
		enableL2:     config.EnableL2,
		enableL3:     config.EnableL3,
		connectorID:  config.ConnectorID,
		endpointName: config.EndpointName,
	}

	// Inicializa L1 (in-memory)
	if config.EnableL1 {
		l1, err := NewL1MemoryCache(config.L1Config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize L1 cache: %w", err)
		}
		manager.l1 = l1
	}

	// Inicializa L2 (Redis)
	if config.EnableL2 {
		l2, err := NewL2RedisCache(config.L2Config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize L2 cache: %w", err)
		}
		manager.l2 = l2
	}

	// L3 será inicializado externamente e injetado via SetL3Cache()

	return manager, nil
}

// SetL3Cache injeta implementação de L3 cache (PostgreSQL)
func (m *MultiLevelCacheManager) SetL3Cache(l3 L3Cache) {
	m.l3 = l3
}

// Get busca valor em cascata (L1 → L2 → L3 → retorna nil)
// Retorna: (value, cacheLevel, error)
func (m *MultiLevelCacheManager) Get(ctx context.Context, key string) (interface{}, CacheLevel, error) {
	startTime := time.Now()

	// Tenta L1 (in-memory)
	if m.enableL1 && m.l1 != nil {
		if value, found := m.l1.Get(ctx, key); found {
			RecordCacheHit(string(LevelL1), m.connectorID, m.endpointName)
			RecordCacheLatency(string(LevelL1), "get", time.Since(startTime).Seconds())
			return value, LevelL1, nil
		}
		RecordCacheMiss(string(LevelL1), m.connectorID, m.endpointName)
	}

	// Tenta L2 (Redis)
	if m.enableL2 && m.l2 != nil {
		l2StartTime := time.Now()
		value, found, err := m.l2.Get(ctx, key)
		RecordCacheLatency(string(LevelL2), "get", time.Since(l2StartTime).Seconds())

		if err != nil {
			RecordCacheError(string(LevelL2), "get", "redis_error")
			// Continua para L3 em caso de erro no Redis
		} else if found {
			RecordCacheHit(string(LevelL2), m.connectorID, m.endpointName)

			// Promove para L1 (warm up)
			if m.enableL1 && m.l1 != nil {
				if err := m.l1.Set(ctx, key, value, 1024); err == nil {
					RecordCachePromotion(string(LevelL2), string(LevelL1))
				}
			}

			return value, LevelL2, nil
		}
		RecordCacheMiss(string(LevelL2), m.connectorID, m.endpointName)
	}

	// Tenta L3 (PostgreSQL)
	if m.enableL3 && m.l3 != nil {
		l3StartTime := time.Now()
		value, found, err := m.l3.Get(ctx, key)
		RecordCacheLatency(string(LevelL3), "get", time.Since(l3StartTime).Seconds())

		if err != nil {
			RecordCacheError(string(LevelL3), "get", "postgres_error")
		} else if found {
			RecordCacheHit(string(LevelL3), m.connectorID, m.endpointName)

			// Promove para L2 e L1
			if m.enableL2 && m.l2 != nil {
				if err := m.l2.Set(ctx, key, value); err == nil {
					RecordCachePromotion(string(LevelL3), string(LevelL2))
				}
			}
			if m.enableL1 && m.l1 != nil {
				if err := m.l1.Set(ctx, key, value, 1024); err == nil {
					RecordCachePromotion(string(LevelL3), string(LevelL1))
				}
			}

			return value, LevelL3, nil
		}
		RecordCacheMiss(string(LevelL3), m.connectorID, m.endpointName)
	}

	// Nenhum cache hit
	return nil, LevelExternal, nil
}

// Set armazena valor em todos os níveis habilitados (cascata reversa)
func (m *MultiLevelCacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var errs []error

	// Armazena em L1
	if m.enableL1 && m.l1 != nil {
		startTime := time.Now()
		if err := m.l1.SetWithTTL(ctx, key, value, 1024, ttl); err != nil {
			errs = append(errs, fmt.Errorf("L1 set error: %w", err))
			RecordCacheError(string(LevelL1), "set", "set_failed")
		} else {
			RecordCacheSet(string(LevelL1), m.connectorID, m.endpointName)
			RecordCacheLatency(string(LevelL1), "set", time.Since(startTime).Seconds())
		}
	}

	// Armazena em L2
	if m.enableL2 && m.l2 != nil {
		startTime := time.Now()
		if err := m.l2.SetWithTTL(ctx, key, value, ttl); err != nil {
			errs = append(errs, fmt.Errorf("L2 set error: %w", err))
			RecordCacheError(string(LevelL2), "set", "set_failed")
		} else {
			RecordCacheSet(string(LevelL2), m.connectorID, m.endpointName)
			RecordCacheLatency(string(LevelL2), "set", time.Since(startTime).Seconds())
		}
	}

	// Armazena em L3
	if m.enableL3 && m.l3 != nil {
		startTime := time.Now()
		if err := m.l3.Set(ctx, key, value, ttl); err != nil {
			errs = append(errs, fmt.Errorf("L3 set error: %w", err))
			RecordCacheError(string(LevelL3), "set", "set_failed")
		} else {
			RecordCacheSet(string(LevelL3), m.connectorID, m.endpointName)
			RecordCacheLatency(string(LevelL3), "set", time.Since(startTime).Seconds())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache set errors: %v", errs)
	}

	return nil
}

// Delete remove valor de todos os níveis
func (m *MultiLevelCacheManager) Delete(ctx context.Context, key string) error {
	var errs []error

	if m.enableL1 && m.l1 != nil {
		m.l1.Delete(ctx, key)
	}

	if m.enableL2 && m.l2 != nil {
		if err := m.l2.Delete(ctx, key); err != nil {
			errs = append(errs, fmt.Errorf("L2 delete error: %w", err))
		}
	}

	if m.enableL3 && m.l3 != nil {
		if err := m.l3.Delete(ctx, key); err != nil {
			errs = append(errs, fmt.Errorf("L3 delete error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache delete errors: %v", errs)
	}

	return nil
}

// Clear limpa todos os níveis de cache
func (m *MultiLevelCacheManager) Clear(ctx context.Context) error {
	var errs []error

	if m.enableL1 && m.l1 != nil {
		m.l1.Clear(ctx)
	}

	if m.enableL2 && m.l2 != nil {
		if err := m.l2.Clear(ctx); err != nil {
			errs = append(errs, fmt.Errorf("L2 clear error: %w", err))
		}
	}

	if m.enableL3 && m.l3 != nil {
		if err := m.l3.Delete(ctx, "*"); err != nil {
			errs = append(errs, fmt.Errorf("L3 clear error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache clear errors: %v", errs)
	}

	return nil
}

// GetStats retorna estatísticas agregadas de todos os níveis
func (m *MultiLevelCacheManager) GetStats(ctx context.Context) map[string]interface{} {
	stats := make(map[string]interface{})

	if m.enableL1 && m.l1 != nil {
		l1Stats := m.l1.GetStats()
		stats["l1"] = map[string]interface{}{
			"hits":       l1Stats.Hits,
			"misses":     l1Stats.Misses,
			"sets":       l1Stats.Sets,
			"evictions":  l1Stats.Evictions,
			"size_bytes": l1Stats.Size,
			"hit_rate":   m.l1.GetHitRate(),
		}

		// Atualiza gauge Prometheus
		SetCacheHitRate(string(LevelL1), m.l1.GetHitRate())
		SetCacheSize(string(LevelL1), l1Stats.Size)
	}

	if m.enableL2 && m.l2 != nil {
		l2Stats, err := m.l2.GetStats(ctx)
		if err == nil {
			stats["l2"] = map[string]interface{}{
				"pool_hits":   l2Stats.Hits,
				"pool_misses": l2Stats.Misses,
				"pool_total":  l2Stats.TotalConns,
				"pool_idle":   l2Stats.IdleConns,
			}
		}
	}

	return stats
}

// Close fecha todas as conexões de cache
func (m *MultiLevelCacheManager) Close() error {
	var errs []error

	if m.l1 != nil {
		m.l1.Close()
	}

	if m.l2 != nil {
		if err := m.l2.Close(); err != nil {
			errs = append(errs, fmt.Errorf("L2 close error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache close errors: %v", errs)
	}

	return nil
}
