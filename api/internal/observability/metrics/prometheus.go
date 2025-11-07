package metrics

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics for the BGC API
var (
	// HTTP Request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bgc_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bgc_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// Database metrics
	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bgc_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5},
		},
		[]string{"operation", "table"},
	)

	dbConnectionsOpen = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bgc_db_connections_open",
			Help: "Number of open database connections",
		},
	)

	dbConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bgc_db_connections_in_use",
			Help: "Number of database connections currently in use",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bgc_db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// Application errors
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bgc_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"type", "severity"},
	)

	// Idempotency cache metrics
	idempotencyCacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bgc_idempotency_cache_hits_total",
			Help: "Total number of idempotency cache hits",
		},
	)

	idempotencyCacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bgc_idempotency_cache_misses_total",
			Help: "Total number of idempotency cache misses",
		},
	)

	idempotencyCacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bgc_idempotency_cache_size",
			Help: "Current number of entries in idempotency cache",
		},
	)
)

// PrometheusMiddleware returns a Gin middleware that collects Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get route info
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		status := c.Writer.Status()

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, path, statusClass(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

// RecordDBQuery records a database query execution
func RecordDBQuery(operation, table string, duration time.Duration) {
	dbQueriesTotal.WithLabelValues(operation, table).Inc()
	dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordError records an application error
func RecordError(errorType, severity string) {
	errorsTotal.WithLabelValues(errorType, severity).Inc()
}

// RecordIdempotencyCacheHit records an idempotency cache hit
func RecordIdempotencyCacheHit() {
	idempotencyCacheHits.Inc()
}

// RecordIdempotencyCacheMiss records an idempotency cache miss
func RecordIdempotencyCacheMiss() {
	idempotencyCacheMisses.Inc()
}

// UpdateIdempotencyCacheSize updates the current idempotency cache size
func UpdateIdempotencyCacheSize(size int) {
	idempotencyCacheSize.Set(float64(size))
}

// UpdateDBConnectionStats updates database connection pool metrics
func UpdateDBConnectionStats(db *sql.DB) {
	stats := db.Stats()
	dbConnectionsOpen.Set(float64(stats.OpenConnections))
	dbConnectionsInUse.Set(float64(stats.InUse))
	dbConnectionsIdle.Set(float64(stats.Idle))
}

// StartDBStatsCollector starts a goroutine that periodically collects DB stats
func StartDBStatsCollector(db *sql.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			UpdateDBConnectionStats(db)
		}
	}()
}

// statusClass returns the HTTP status class (2xx, 4xx, 5xx, etc)
func statusClass(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
