package cache

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CacheHits total de cache hits por nível
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "integration_gateway_cache_hits_total",
			Help: "Total number of cache hits by level",
		},
		[]string{"level", "connector", "endpoint"},
	)

	// CacheMisses total de cache misses por nível
	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "integration_gateway_cache_misses_total",
			Help: "Total number of cache misses by level",
		},
		[]string{"level", "connector", "endpoint"},
	)

	// CacheLatency latência de operações de cache
	CacheLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "integration_gateway_cache_latency_seconds",
			Help:    "Cache operation latency in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1, .25, .5},
		},
		[]string{"level", "operation"}, // operation: get, set, delete
	)

	// CacheSize tamanho atual do cache em bytes
	CacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "integration_gateway_cache_size_bytes",
			Help: "Current cache size in bytes by level",
		},
		[]string{"level"},
	)

	// CacheEvictions total de evictions por nível
	CacheEvictions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "integration_gateway_cache_evictions_total",
			Help: "Total number of cache evictions by level",
		},
		[]string{"level"},
	)

	// CacheHitRate taxa de hit do cache (0.0 a 1.0)
	CacheHitRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "integration_gateway_cache_hit_rate",
			Help: "Cache hit rate (0.0 to 1.0) by level",
		},
		[]string{"level"},
	)

	// CacheSets total de operações set por nível
	CacheSets = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "integration_gateway_cache_sets_total",
			Help: "Total number of cache set operations by level",
		},
		[]string{"level", "connector", "endpoint"},
	)

	// CachePromotions total de promoções entre níveis de cache
	CachePromotions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "integration_gateway_cache_promotions_total",
			Help: "Total number of cache promotions between levels",
		},
		[]string{"from_level", "to_level"},
	)

	// CacheErrors erros de operações de cache
	CacheErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "integration_gateway_cache_errors_total",
			Help: "Total number of cache operation errors",
		},
		[]string{"level", "operation", "error_type"},
	)
)

// RecordCacheHit registra um cache hit
func RecordCacheHit(level, connector, endpoint string) {
	CacheHits.WithLabelValues(level, connector, endpoint).Inc()
}

// RecordCacheMiss registra um cache miss
func RecordCacheMiss(level, connector, endpoint string) {
	CacheMisses.WithLabelValues(level, connector, endpoint).Inc()
}

// RecordCacheLatency registra latência de operação
func RecordCacheLatency(level, operation string, durationSeconds float64) {
	CacheLatency.WithLabelValues(level, operation).Observe(durationSeconds)
}

// SetCacheSize atualiza tamanho do cache
func SetCacheSize(level string, sizeBytes uint64) {
	CacheSize.WithLabelValues(level).Set(float64(sizeBytes))
}

// RecordCacheEviction registra uma eviction
func RecordCacheEviction(level string) {
	CacheEvictions.WithLabelValues(level).Inc()
}

// SetCacheHitRate atualiza taxa de hit
func SetCacheHitRate(level string, hitRate float64) {
	CacheHitRate.WithLabelValues(level).Set(hitRate)
}

// RecordCacheSet registra uma operação set
func RecordCacheSet(level, connector, endpoint string) {
	CacheSets.WithLabelValues(level, connector, endpoint).Inc()
}

// RecordCachePromotion registra promoção entre níveis
func RecordCachePromotion(fromLevel, toLevel string) {
	CachePromotions.WithLabelValues(fromLevel, toLevel).Inc()
}

// RecordCacheError registra erro de operação
func RecordCacheError(level, operation, errorType string) {
	CacheErrors.WithLabelValues(level, operation, errorType).Inc()
}
