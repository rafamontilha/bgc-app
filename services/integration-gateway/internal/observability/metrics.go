package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ConnectorRequests total de requisições por connector
	ConnectorRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_connector_requests_total",
			Help: "Total number of requests per connector",
		},
		[]string{"connector", "endpoint", "status"},
	)

	// ConnectorDuration duração das requisições
	ConnectorDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bgc_connector_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"connector", "endpoint"},
	)

	// ConnectorCircuitBreakerState estado do circuit breaker
	ConnectorCircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bgc_connector_circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
		},
		[]string{"connector"},
	)

	// ConnectorRateLimitRemaining requisições restantes no rate limit
	ConnectorRateLimitRemaining = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bgc_connector_rate_limit_remaining",
			Help: "Number of requests remaining in rate limit window",
		},
		[]string{"connector"},
	)

	// ConnectorCacheHits cache hits
	ConnectorCacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_connector_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"connector", "endpoint"},
	)

	// ConnectorCacheMisses cache misses
	ConnectorCacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_connector_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"connector", "endpoint"},
	)

	// ConnectorRetries número de retries
	ConnectorRetries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_connector_retries_total",
			Help: "Total number of retries",
		},
		[]string{"connector", "endpoint"},
	)

	// ConnectorErrors erros por tipo
	ConnectorErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_connector_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"connector", "endpoint", "error_type"},
	)

	// CertificateExpiryDays dias até expiração do certificado
	CertificateExpiryDays = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bgc_certificate_expiry_days",
			Help: "Days until certificate expiration",
		},
		[]string{"certificate_ref"},
	)

	// TransformPluginDuration duração de plugins de transformação
	TransformPluginDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bgc_transform_plugin_duration_seconds",
			Help:    "Transform plugin execution duration",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"plugin"},
	)
)

// RecordRequest registra uma requisição
func RecordRequest(connector, endpoint, status string, duration float64) {
	ConnectorRequests.WithLabelValues(connector, endpoint, status).Inc()
	ConnectorDuration.WithLabelValues(connector, endpoint).Observe(duration)
}

// RecordCacheHit registra cache hit
func RecordCacheHit(connector, endpoint string) {
	ConnectorCacheHits.WithLabelValues(connector, endpoint).Inc()
}

// RecordCacheMiss registra cache miss
func RecordCacheMiss(connector, endpoint string) {
	ConnectorCacheMisses.WithLabelValues(connector, endpoint).Inc()
}

// RecordRetry registra retry
func RecordRetry(connector, endpoint string) {
	ConnectorRetries.WithLabelValues(connector, endpoint).Inc()
}

// RecordError registra erro
func RecordError(connector, endpoint, errorType string) {
	ConnectorErrors.WithLabelValues(connector, endpoint, errorType).Inc()
}

// SetCircuitBreakerState atualiza estado do circuit breaker
// 0=closed, 1=half-open, 2=open
func SetCircuitBreakerState(connector string, state float64) {
	ConnectorCircuitBreakerState.WithLabelValues(connector).Set(state)
}

// SetRateLimitRemaining atualiza requisições restantes
func SetRateLimitRemaining(connector string, remaining float64) {
	ConnectorRateLimitRemaining.WithLabelValues(connector).Set(remaining)
}

// SetCertificateExpiryDays atualiza dias até expiração
func SetCertificateExpiryDays(certificateRef string, days float64) {
	CertificateExpiryDays.WithLabelValues(certificateRef).Set(days)
}

// RecordTransformPlugin registra execução de plugin
func RecordTransformPlugin(plugin string, duration float64) {
	TransformPluginDuration.WithLabelValues(plugin).Observe(duration)
}
