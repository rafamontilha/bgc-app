package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// L2RedisCache implementa cache distribuído com Redis
type L2RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
	prefix     string // Prefixo para todas as chaves (ex: "bgc:cache:")
}

// L2Config configuração do cache L2
type L2Config struct {
	Addr       string        // Endereço do Redis (host:port)
	Password   string        // Senha (via KubernetesSecretStore)
	DB         int           // Database number (default: 0)
	DefaultTTL time.Duration // TTL padrão (default: 7 dias)
	Prefix     string        // Prefixo de chaves (default: "bgc:cache:")
	MaxRetries int           // Tentativas de retry (default: 3)
	PoolSize   int           // Tamanho do connection pool (default: 10)
}

// DefaultL2Config retorna configuração padrão do L2
func DefaultL2Config() L2Config {
	return L2Config{
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 7 * 24 * time.Hour, // 7 dias
		Prefix:     "bgc:cache:",
		MaxRetries: 3,
		PoolSize:   10,
	}
}

// NewL2RedisCache cria um novo cache L2 com Redis
func NewL2RedisCache(config L2Config) (*L2RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       config.Addr,
		Password:   config.Password,
		DB:         config.DB,
		MaxRetries: config.MaxRetries,
		PoolSize:   config.PoolSize,

		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		// Connection pool
		MinIdleConns: 2,
		MaxIdleConns: 5,
		ConnMaxLifetime: 5 * time.Minute,
	})

	// Testa conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", config.Addr, err)
	}

	return &L2RedisCache{
		client:     client,
		defaultTTL: config.DefaultTTL,
		prefix:     config.Prefix,
	}, nil
}

// Get busca valor do cache L2
func (l *L2RedisCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	fullKey := l.prefix + key

	val, err := l.client.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		// Key não existe
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("redis get error for key %s: %w", fullKey, err)
	}

	// Desserializa JSON
	var data interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal cached value for key %s: %w", fullKey, err)
	}

	return data, true, nil
}

// Set armazena valor no cache L2 com TTL padrão
func (l *L2RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return l.SetWithTTL(ctx, key, value, l.defaultTTL)
}

// SetWithTTL armazena valor no cache L2 com TTL customizado
func (l *L2RedisCache) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := l.prefix + key

	// Serializa para JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", fullKey, err)
	}

	// Armazena no Redis
	if err := l.client.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set error for key %s: %w", fullKey, err)
	}

	return nil
}

// Delete remove valor do cache L2
func (l *L2RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := l.prefix + key

	if err := l.client.Del(ctx, fullKey).Err(); err != nil {
		return fmt.Errorf("redis delete error for key %s: %w", fullKey, err)
	}

	return nil
}

// Clear limpa todo o cache L2 (apenas chaves com o prefixo)
func (l *L2RedisCache) Clear(ctx context.Context) error {
	pattern := l.prefix + "*"

	// Usa SCAN para não bloquear o Redis
	iter := l.client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := l.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("redis clear error: %w", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis scan error during clear: %w", err)
	}

	return nil
}

// Exists verifica se uma chave existe
func (l *L2RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := l.prefix + key

	count, err := l.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error for key %s: %w", fullKey, err)
	}

	return count > 0, nil
}

// GetTTL retorna o TTL restante de uma chave
func (l *L2RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := l.prefix + key

	ttl, err := l.client.TTL(ctx, fullKey).Result()
	if err != nil {
		return 0, fmt.Errorf("redis ttl error for key %s: %w", fullKey, err)
	}

	return ttl, nil
}

// GetStats retorna estatísticas do Redis
func (l *L2RedisCache) GetStats(ctx context.Context) (*redis.PoolStats, error) {
	stats := l.client.PoolStats()
	return stats, nil
}

// Ping verifica conectividade com Redis
func (l *L2RedisCache) Ping(ctx context.Context) error {
	return l.client.Ping(ctx).Err()
}

// Close fecha a conexão com Redis
func (l *L2RedisCache) Close() error {
	return l.client.Close()
}
