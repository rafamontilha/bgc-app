# Observability & Monitoring Guide

**Version:** 1.0
**Last Updated:** 2025-11-06
**Epic:** 2 - Observabilidade & Padrões

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Prometheus Metrics](#prometheus-metrics)
4. [Distributed Tracing (OpenTelemetry)](#distributed-tracing-opentelemetry)
5. [Dashboards & Visualization](#dashboards--visualization)
6. [Alerting](#alerting)
7. [Local Development Setup](#local-development-setup)
8. [Kubernetes Deployment](#kubernetes-deployment)
9. [Metrics Reference](#metrics-reference)
10. [Troubleshooting](#troubleshooting)

---

## Overview

BGC App implements a comprehensive observability stack based on industry best practices:

- **Metrics**: Prometheus for time-series metrics collection
- **Tracing**: OpenTelemetry for distributed tracing
- **Visualization**: Grafana for dashboards and analytics
- **Tracing UI**: Jaeger for trace visualization

### Key Features

✅ **Automatic Instrumentation**
- HTTP requests (latency, throughput, errors)
- Database queries (duration, count, connection pool)
- Business metrics (idempotency cache, custom counters)

✅ **Distributed Tracing**
- End-to-end request tracing across services
- Database query spans with parameters
- Automatic trace context propagation

✅ **Pre-built Dashboards**
- API performance overview
- Database health monitoring
- Error rate tracking
- Real-time metrics

✅ **Production-ready Alerts**
- High error rates (> 5%)
- High latency (P95 > 2s)
- Database connection exhaustion
- Service health checks

---

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│              BGC API (Go + Gin)                 │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │   OpenTelemetry Middleware               │  │
│  │   - Trace ID propagation                 │  │
│  │   - Span creation for handlers           │  │
│  └──────────────────────────────────────────┘  │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │   Prometheus Middleware                  │  │
│  │   - Request counter                      │  │
│  │   - Duration histogram                   │  │
│  │   - In-flight gauge                      │  │
│  └──────────────────────────────────────────┘  │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │   Repository Layer                       │  │
│  │   - DB query spans                       │  │
│  │   - Query duration metrics               │  │
│  └──────────────────────────────────────────┘  │
└─────────────┬───────────────────┬───────────────┘
              │                   │
              ▼                   ▼
    ┌─────────────────┐   ┌──────────────┐
    │   PostgreSQL    │   │    Jaeger    │
    │   (Metrics)     │   │   (Traces)   │
    └─────────────────┘   └──────────────┘
              │                   │
              ▼                   │
    ┌─────────────────┐           │
    │   Prometheus    │◄──────────┘
    │  (Scraper)      │
    └────────┬────────┘
             │
             ▼
    ┌─────────────────┐
    │    Grafana      │
    │  (Dashboards)   │
    └─────────────────┘
```

---

## Prometheus Metrics

### Available Metrics

#### HTTP Request Metrics

| Metric Name | Type | Labels | Description |
|-------------|------|--------|-------------|
| `bgc_http_requests_total` | Counter | method, path, status | Total HTTP requests |
| `bgc_http_request_duration_seconds` | Histogram | method, path | Request duration |
| `bgc_http_requests_in_flight` | Gauge | - | Current requests being processed |

#### Database Metrics

| Metric Name | Type | Labels | Description |
|-------------|------|--------|-------------|
| `bgc_db_queries_total` | Counter | operation, table | Total DB queries executed |
| `bgc_db_query_duration_seconds` | Histogram | operation, table | DB query duration |
| `bgc_db_connections_open` | Gauge | - | Number of open DB connections |
| `bgc_db_connections_in_use` | Gauge | - | Number of DB connections in use |
| `bgc_db_connections_idle` | Gauge | - | Number of idle DB connections |

#### Application Metrics

| Metric Name | Type | Labels | Description |
|-------------|------|--------|-------------|
| `bgc_errors_total` | Counter | type, severity | Application errors |
| `bgc_idempotency_cache_hits_total` | Counter | - | Idempotency cache hits |
| `bgc_idempotency_cache_misses_total` | Counter | - | Idempotency cache misses |
| `bgc_idempotency_cache_size` | Gauge | - | Current cache size |

### Metric Collection

Metrics are collected automatically via middleware:

```go
// api/internal/app/server.go
r.Use(metrics.PrometheusMiddleware()) // Automatic HTTP metrics
```

Database metrics are collected during query execution:

```go
// api/internal/repository/postgres/market.go
start := time.Now()
rows, err := r.db.QueryContext(ctx, query, args...)
duration := time.Since(start)
metrics.RecordDBQuery("SELECT", "v_tam_by_year_chapter", duration)
```

### Accessing Metrics

**Local (Docker Compose):**
- API Metrics: http://localhost:8080/metrics
- Prometheus UI: http://localhost:9090

**Kubernetes:**
```bash
kubectl port-forward -n observability svc/prometheus 9090:9090
```

---

## Distributed Tracing (OpenTelemetry)

### Configuration

The tracer is initialized automatically on application startup:

```go
// api/cmd/api/main.go
tracerShutdown, err := tracing.InitTracer("bgc-api", "production")
defer tracerShutdown(context.Background())
```

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint | `jaeger:4317` |
| `ENVIRONMENT` | Deployment environment | `production` |

### Automatic Instrumentation

#### HTTP Handlers

All HTTP requests are automatically traced:

```go
// api/internal/app/server.go
r.Use(otelgin.Middleware("bgc-api"))
```

Each request creates a span with:
- HTTP method, path, status code
- Request duration
- Error information (if any)

#### Database Queries

All database queries are instrumented with spans:

```go
// api/internal/repository/postgres/market.go
ctx, span := tracing.StartSpan(context.Background(), "db.GetMarketDataByYearRange")
defer span.End()

span.SetAttributes(
    attribute.Int("year.from", yearFrom),
    attribute.Int("year.to", yearTo),
)

rows, err := r.db.QueryContext(ctx, query, args...)
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return nil, err
}

span.SetAttributes(attribute.Int("result.count", len(items)))
span.SetStatus(codes.Ok, "query successful")
```

### Trace Context Propagation

Trace context is automatically propagated via HTTP headers:
- `traceparent`: W3C Trace Context
- `tracestate`: Additional vendor-specific data

### Viewing Traces

**Local (Docker Compose):**
- Jaeger UI: http://localhost:16686

**Kubernetes:**
```bash
kubectl port-forward -n observability svc/jaeger-query 16686:16686
```

**Example Trace:**
```
Span: GET /v1/market/size
├─ Span: db.GetMarketDataByYearRange
│  └─ Attributes:
│     - year.from: 2020
│     - year.to: 2023
│     - result.count: 42
│     - db.system: postgresql
└─ Duration: 125ms
```

---

## Dashboards & Visualization

### Grafana Setup

**Access Grafana:**
- Docker Compose: http://localhost:3001
- Kubernetes: http://localhost:3000 (via port-forward)

**Default Credentials:**
- Username: `admin`
- Password: `admin` (change in production!)

### Pre-configured Dashboards

#### BGC API Overview

Location: `bgcstack/observability/dashboards/bgc-api.json`

**Panels:**
1. **Request Rate** - Requests per second by endpoint
2. **Error Rate** - 4xx/5xx errors as percentage
3. **Request Duration** - P50, P95, P99 latency
4. **Database Connections** - Open, in-use, idle connections

**Queries:**
```promql
# Request rate
sum(rate(bgc_http_requests_total[5m])) by (path)

# Error rate
sum(rate(bgc_http_requests_total{status="5xx"}[5m]))
  / sum(rate(bgc_http_requests_total[5m])) * 100

# P95 latency
histogram_quantile(0.95,
  sum(rate(bgc_http_request_duration_seconds_bucket[5m])) by (le)
)

# DB connections
bgc_db_connections_in_use
```

### Creating Custom Dashboards

1. Open Grafana UI
2. Click **Dashboards** → **New Dashboard**
3. Add panel → Select **Prometheus** datasource
4. Enter PromQL query
5. Configure visualization (Graph, Stat, Gauge, etc.)
6. Save dashboard

---

## Alerting

### Alert Rules

Location: `k8s/observability/prometheus-alert-rules.yaml`

**Critical Alerts:**

| Alert Name | Condition | Threshold | Duration | Action |
|------------|-----------|-----------|----------|--------|
| `HighErrorRate` | Error rate > 5% | 5% | 5 minutes | Page on-call |
| `APIDown` | Service unavailable | - | 1 minute | Page immediately |
| `DatabaseConnectionPoolExhaustion` | Connections > 90% | 90% | 5 minutes | Investigate |

**Warning Alerts:**

| Alert Name | Condition | Threshold | Duration | Action |
|------------|-----------|-----------|----------|--------|
| `HighLatencyP95` | P95 > 2 seconds | 2s | 10 minutes | Investigate |
| `HighDatabaseLatency` | DB P95 > 500ms | 500ms | 10 minutes | Review queries |
| `HighRequestRate` | Requests > 1000/s | 1000 req/s | 5 minutes | Check capacity |

### Example Alert Rule

```yaml
- alert: HighErrorRate
  expr: |
    (
      sum(rate(bgc_http_requests_total{status="5xx"}[5m]))
      /
      sum(rate(bgc_http_requests_total[5m]))
    ) > 0.05
  for: 5m
  labels:
    severity: critical
    component: api
  annotations:
    summary: "High error rate detected"
    description: "API error rate is {{ $value | humanizePercentage }}"
```

### Alert Notifications (Kubernetes)

To configure alert notifications, deploy Alertmanager:

```bash
kubectl apply -f k8s/observability/alertmanager.yaml
```

Configure notification channels (Slack, PagerDuty, email) in:
```yaml
# k8s/observability/alertmanager-config.yaml
receivers:
  - name: 'slack'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
        channel: '#alerts'
```

---

## Local Development Setup

### Option 1: Docker Compose (Recommended)

1. **Start all services:**
```bash
cd bgcstack
docker-compose up -d
```

2. **Access observability UIs:**
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001 (admin/admin)
- Jaeger: http://localhost:16686
- API Metrics: http://localhost:8080/metrics

3. **Generate traffic:**
```bash
# Make some API requests
curl http://localhost:8080/v1/market/size?year_from=2020&year_to=2023

# Check metrics
curl http://localhost:8080/metrics | grep bgc_http_requests_total
```

4. **View in Grafana:**
- Open http://localhost:3001
- Navigate to **Dashboards** → **BGC** → **BGC API Overview**

5. **View traces in Jaeger:**
- Open http://localhost:16686
- Select service: `bgc-api`
- Click **Find Traces**

### Option 2: Standalone API (Development)

If running the API locally without Docker:

```bash
cd api
export ENVIRONMENT=development
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
go run cmd/api/main.go
```

Metrics available at: http://localhost:8080/metrics

---

## Kubernetes Deployment

### Deploy Observability Stack

```bash
# Create namespace and deploy all services
kubectl apply -f k8s/observability/

# Verify deployments
kubectl get pods -n observability

# Expected output:
# NAME                          READY   STATUS    RESTARTS   AGE
# prometheus-xxx                1/1     Running   0          1m
# grafana-xxx                   1/1     Running   0          1m
# jaeger-xxx                    1/1     Running   0          1m
```

### Access Services (Port Forwarding)

```bash
# Prometheus
kubectl port-forward -n observability svc/prometheus 9090:9090

# Grafana
kubectl port-forward -n observability svc/grafana 3000:3000

# Jaeger UI
kubectl port-forward -n observability svc/jaeger-query 16686:16686
```

### Access via NodePort

Services are exposed via NodePort:
- Grafana: http://<node-ip>:30030
- Jaeger: http://<node-ip>:30016

Get node IP:
```bash
kubectl get nodes -o wide
```

### Configure API to Send Metrics

The API automatically discovers Prometheus via Kubernetes service discovery. Ensure your API deployment has the correct labels:

```yaml
# k8s/api/deployment.yaml
metadata:
  labels:
    app: bgc-api
```

---

## Metrics Reference

### Custom Metrics

To add custom application metrics:

```go
import "bgc-app/internal/observability/metrics"

// Record an error
metrics.RecordError("validation_error", "warning")

// Record idempotency cache hit
metrics.RecordIdempotencyCacheHit()

// Update cache size
metrics.UpdateIdempotencyCacheSize(cacheSize)
```

### PromQL Query Examples

**Request rate by status code:**
```promql
sum(rate(bgc_http_requests_total[5m])) by (status)
```

**Average request duration:**
```promql
avg(rate(bgc_http_request_duration_seconds_sum[5m]) / rate(bgc_http_request_duration_seconds_count[5m]))
```

**Database query rate:**
```promql
sum(rate(bgc_db_queries_total[5m])) by (table)
```

**Cache hit ratio:**
```promql
rate(bgc_idempotency_cache_hits_total[5m]) /
  (rate(bgc_idempotency_cache_hits_total[5m]) + rate(bgc_idempotency_cache_misses_total[5m]))
```

---

## Troubleshooting

### Metrics not appearing in Prometheus

1. **Check API /metrics endpoint:**
```bash
curl http://localhost:8080/metrics
```

2. **Verify Prometheus scrape config:**
```bash
# Docker Compose
cat bgcstack/observability/prometheus.yml

# Kubernetes
kubectl get cm -n observability prometheus-config -o yaml
```

3. **Check Prometheus targets:**
- Open http://localhost:9090
- Go to **Status** → **Targets**
- Ensure `bgc-api` target is **UP**

### Traces not appearing in Jaeger

1. **Verify OTLP endpoint:**
```bash
# Docker Compose
echo $OTEL_EXPORTER_OTLP_ENDPOINT  # Should be: jaeger:4317

# Check Jaeger collector
docker logs bgc_jaeger
```

2. **Check API logs for tracer initialization:**
```bash
docker logs bgc_api | grep -i "tracer initialized"
```

3. **Verify Jaeger is receiving spans:**
```bash
# Check Jaeger collector metrics
curl http://localhost:14269/metrics
```

### Grafana dashboards empty

1. **Verify datasource configuration:**
- Open Grafana UI → **Configuration** → **Data Sources**
- Test connection to Prometheus

2. **Check if metrics exist:**
```bash
# Open Prometheus expression browser
http://localhost:9090/graph

# Run query:
bgc_http_requests_total
```

3. **Verify time range:**
- Grafana default is "Last 15 minutes"
- Ensure you have traffic in that window

### High memory usage in Prometheus

1. **Reduce retention period:**
```bash
# Edit prometheus.yml or deployment
--storage.tsdb.retention.time=7d  # Instead of 30d
```

2. **Reduce scrape frequency:**
```yaml
# prometheus.yml
global:
  scrape_interval: 30s  # Instead of 15s
```

---

## Best Practices

### Metrics

✅ **DO:**
- Use consistent label names across metrics
- Keep cardinality low (avoid user IDs in labels)
- Use histograms for latency measurements
- Document custom metrics

❌ **DON'T:**
- Create high-cardinality metrics (millions of series)
- Use metrics for logging (use structured logs instead)
- Over-instrument (measure what matters)

### Tracing

✅ **DO:**
- Trace critical paths (HTTP handlers, DB queries)
- Add meaningful attributes to spans
- Record errors in spans
- Propagate trace context

❌ **DON'T:**
- Trace every function (high overhead)
- Store sensitive data in span attributes
- Create too many child spans (> 100 per trace)

### Dashboards

✅ **DO:**
- Focus on actionable metrics (RED: Rate, Errors, Duration)
- Use percentiles (P50, P95, P99) for latency
- Include context (comparisons, thresholds)
- Version control dashboard JSON files

❌ **DON'T:**
- Create dashboards without alerts
- Use only averages (hides outliers)
- Make dashboards too busy (> 12 panels)

---

## References

- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/instrumentation/go/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)

---

**Maintained by:** BGC Development Team
**Questions:** [GitHub Issues](https://github.com/rafamontilha/bgc-app/issues)
