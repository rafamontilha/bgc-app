# Cache Multinível - Integration Gateway

Sistema de cache em 3 níveis (L1 → L2 → L3) para o Integration Gateway do BGC, com objetivo de reduzir chamadas a APIs externas de 1000/dia para < 20/dia.

## 🎯 Objetivos

- **Hit Rate**: 98% (L1 + L2 + L3)
- **Latência P95**: < 50ms (com cache)
- **Redução de Chamadas**: 1000/dia → < 20/dia ao ComexStat
- **Throughput**: > 1000 req/s

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Request                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  MultiLevelCache Manager      │
         └───────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
   ┌────────┐      ┌─────────┐      ┌─────────┐
   │ L1 (🔥)│      │ L2 (💾) │      │ L3 (🗄️) │
   │Ristretto│─────▶│  Redis  │─────▶│PostgreSQL│
   │100MB   │      │ 512MB   │      │   2GB   │
   │LFU     │      │ LRU     │      │  MView  │
   │5min TTL│      │ 7d TTL  │      │ 30d TTL │
   └────────┘      └─────────┘      └─────────┘
        │                │                │
        │    ❌ Miss     │    ❌ Miss     │    ❌ Miss
        └────────────────┴────────────────┴────────────┐
                                                        │
                                                        ▼
                                            ┌─────────────────┐
                                            │ External API    │
                                            │ (ComexStat)     │
                                            └─────────────────┘
```

## 📦 Componentes

### L1 - In-Memory Cache (Ristretto)

**Características:**
- **Algoritmo**: LFU (Least Frequently Used)
- **Tamanho**: 100MB máximo
- **TTL**: Configurável (default: 5min)
- **Latência**: ~105 ns/op (read)
- **Escopo**: Por pod (não compartilhado)

**Arquivo**: `l1_memory.go`

```go
config := DefaultL1Config()
cache, err := NewL1MemoryCache(config)
defer cache.Close()

// Set com TTL customizado
cache.SetWithTTL(ctx, "key", "value", 1024, 10*time.Minute)

// Get
value, found := cache.Get(ctx, "key")
```

### L2 - Distributed Cache (Redis)

**Características:**
- **Algoritmo**: LRU (allkeys-lru)
- **Tamanho**: 512MB máximo
- **TTL**: 7 dias (histórico) | 6h (mês atual)
- **Latência**: ~3-5ms (local network)
- **Escopo**: Compartilhado entre pods

**Arquivo**: `l2_redis.go`

```go
config := L2Config{
    Addr:       "redis:6379",
    Password:   "",
    DB:         0,
    DefaultTTL: 7 * 24 * time.Hour,
    Prefix:     "bgc:cache:",
}

cache, err := NewL2RedisCache(config)
defer cache.Close()

// Set
cache.Set(ctx, "key", map[string]interface{}{"data": "value"})

// Get (retorna interface{}, serialização automática em JSON)
value, found, err := cache.Get(ctx, "key")
```

### L3 - Persistent Cache (PostgreSQL)

**Características:**
- **Implementação**: Materialized Views
- **Refresh**: Diário às 3:00 AM (CronJob)
- **TTL**: 30 dias
- **Latência**: ~10-50ms
- **Escopo**: Dados agregados (top NCMs × países)

**Status**: Interface definida, implementação pendente

```go
type L3Cache interface {
    Get(ctx context.Context, key string) (interface{}, bool, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

### MultiLevelCache Manager

Gerencia cascata e promoção automática entre níveis.

**Arquivo**: `manager.go`

```go
config := ManagerConfig{
    L1Config:     DefaultL1Config(),
    L2Config:     DefaultL2Config(),
    EnableL1:     true,
    EnableL2:     true,
    EnableL3:     false,
    ConnectorID:  "comexstat",
    EndpointName: "exportacao_mes",
}

manager, err := NewMultiLevelCacheManager(config)
defer manager.Close()

// Get com cascata automática (L1 → L2 → L3 → External)
value, level, err := manager.Get(ctx, "cache-key")
switch level {
case LevelL1:
    // Hit em L1 (~105ns)
case LevelL2:
    // Hit em L2 (~3-5ms), promovido para L1
case LevelL3:
    // Hit em L3 (~10-50ms), promovido para L2 e L1
case LevelExternal:
    // Miss em todos, buscar API externa
}

// Set propaga para todos os níveis
manager.Set(ctx, "cache-key", data, 7*24*time.Hour)
```

## 📊 Métricas Prometheus

Todas as operações de cache geram métricas Prometheus:

```promql
# Cache hits por nível
integration_gateway_cache_hits_total{level="l1|l2|l3", connector="comexstat", endpoint="exportacao_mes"}

# Cache misses por nível
integration_gateway_cache_misses_total{level="l1|l2|l3", connector="comexstat", endpoint="exportacao_mes"}

# Latência de operações (histogram)
integration_gateway_cache_latency_seconds{level="l1|l2|l3", operation="get|set|delete"}

# Tamanho do cache (bytes)
integration_gateway_cache_size_bytes{level="l1|l2"}

# Taxa de hit (0.0 a 1.0)
integration_gateway_cache_hit_rate{level="l1|l2"}

# Promoções entre níveis
integration_gateway_cache_promotions_total{from_level="l2|l3", to_level="l1|l2"}

# Evictions (L1 apenas)
integration_gateway_cache_evictions_total{level="l1"}

# Erros
integration_gateway_cache_errors_total{level="l1|l2|l3", operation="get|set|delete", error_type="redis_error|..."}
```

## 🧪 Testes

### Testes Unitários

```bash
# L1 (Ristretto)
go test -v ./internal/cache/... -run TestL1

# Manager
go test -v ./internal/cache/... -run TestMultiLevelCacheManager

# Cobertura
go test -cover ./internal/cache/...
# coverage: 82.0% of statements
```

### Testes de Integração (Redis)

```bash
# Iniciar Redis
cd bgcstack && docker-compose up -d redis

# Executar testes de integração
cd services/integration-gateway
go test -tags=integration -v ./internal/cache/... -run Integration

# Todos os testes (unitário + integração)
go test -tags=integration -v ./internal/cache/... -cover
```

### Benchmarks

```bash
go test -bench=. -benchmem ./internal/cache/... -run=^$ -benchtime=3s

# Resultados (AMD Ryzen 5 5600H):
# BenchmarkL1MemoryCache_Get-12         34M ops    105.0 ns/op    23 B/op
# BenchmarkL1MemoryCache_Set-12          1M ops   2777 ns/op     192 B/op
# BenchmarkMultiLevelCacheManager_Get   12M ops    414.0 ns/op    22 B/op
# BenchmarkMultiLevelCacheManager_Set    1M ops   2998 ns/op     192 B/op
```

## 🔧 Configuração

### Variáveis de Ambiente

```bash
# Cache L1 (in-memory)
CACHE_L1_ENABLED=true
CACHE_L1_MAX_SIZE_MB=100
CACHE_L1_DEFAULT_TTL=5m

# Cache L2 (Redis)
CACHE_L2_ENABLED=true
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10

# Cache L3 (PostgreSQL)
CACHE_L3_ENABLED=false  # Implementar depois
```

### Connector Config (YAML)

```yaml
# config/connectors/comexstat.yaml
cache:
  enabled: true
  ttl: 168h  # 7 dias (histórico)
  key_pattern: "comexstat:exp:{ano}:{mes}:{ncm}:{pais}"
```

## 🚀 Uso em Produção

### 1. Docker Compose

```yaml
services:
  redis:
    image: redis:7-alpine
    command: redis-server --maxmemory 512mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]

  integration-gateway:
    environment:
      - REDIS_ADDR=redis:6379
      - CACHE_L1_ENABLED=true
      - CACHE_L2_ENABLED=true
    depends_on:
      - redis
```

### 2. Kubernetes

```bash
# Deploy Redis
kubectl apply -f k8s/redis.yaml

# Deploy Integration Gateway (já configurado)
kubectl apply -f k8s/integration-gateway/deployment.yaml
```

### 3. Monitoramento (Grafana)

Dashboard recomendado:

```json
{
  "title": "BGC Cache Performance",
  "panels": [
    {
      "title": "Cache Hit Rate",
      "targets": ["integration_gateway_cache_hit_rate"]
    },
    {
      "title": "Cache Latency (P95)",
      "targets": ["histogram_quantile(0.95, integration_gateway_cache_latency_seconds)"]
    },
    {
      "title": "Cache Size",
      "targets": ["integration_gateway_cache_size_bytes"]
    }
  ]
}
```

## 📝 Estratégias de Cache

### TTL Dinâmico (ComexStat)

```go
func GetTTLForComexStat(ano, mes int) time.Duration {
    now := time.Now()
    currentYear := now.Year()
    currentMonth := int(now.Month())

    if ano < currentYear {
        // Dados históricos (imutáveis)
        return 7 * 24 * time.Hour // 7 dias
    }
    if ano == currentYear && mes < currentMonth {
        // Mês fechado no ano atual
        return 24 * time.Hour // 1 dia
    }
    // Mês atual (dados podem mudar)
    return 6 * time.Hour // 6 horas
}
```

### Request Coalescing

Agrupa requests idênticas em uma janela de 10s:

```go
// TODO: Implementar em framework/coalescing.go
type RequestCoalescer struct {
    inflight map[string]*sync.WaitGroup
}
```

### Cache Warming

CronJob semanal para pré-popular cache:

```bash
# scripts/warm-cache-comexstat.sh
# Top 50 NCMs × Top 10 países × Últimos 12 meses
```

## 🐛 Troubleshooting

### L1 não está cacheando

```go
// Verificar se Wait() está sendo chamado após Set
cache.SetWithTTL(ctx, key, value, cost, ttl)
// Ristretto é assíncrono, precisa de Wait()
```

### L2 (Redis) com erros de conexão

```bash
# Verificar conectividade
kubectl exec -it deployment/integration-gateway -n data -- redis-cli -h redis ping
# PONG

# Verificar logs
kubectl logs -f deployment/integration-gateway -n data | grep redis
```

### Hit rate baixo (< 90%)

1. **TTL muito curto**: Aumentar TTL no connector config
2. **Keys não padronizadas**: Verificar `key_pattern` no YAML
3. **L1 muito pequeno**: Aumentar `CACHE_L1_MAX_SIZE_MB`
4. **Dados não repetidos**: Revisar padrão de acesso

### Métricas não aparecem no Prometheus

```bash
# Verificar endpoint /metrics
curl http://localhost:8081/metrics | grep integration_gateway_cache

# Verificar scrape config
kubectl get configmap prometheus-config -n observability -o yaml
```

## 📚 Referências

- **Ristretto**: https://github.com/dgraph-io/ristretto
- **Redis Go**: https://github.com/redis/go-redis
- **Prometheus Client**: https://github.com/prometheus/client_golang

## 🔗 Próximos Passos

- [ ] Implementar L3 (PostgreSQL Materialized Views)
- [ ] Request Coalescing (janela de 10s)
- [ ] Cache Warming (CronJob semanal)
- [ ] TTL dinâmico por tipo de dados
- [ ] Dashboard Grafana customizado

---

**Última Atualização**: 2025-01-21
**Cobertura de Testes**: 82%
**Status**: ✅ Produção Ready (L1 + L2)
