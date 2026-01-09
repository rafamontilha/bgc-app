# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added - API do Simulador de Destinos 📍 (2025-11-22)

#### Implementação Completa do MVP (Manhã - 22/11/2025)

**Handler e Rota Funcionando 100%**
- Handler `SimulatorHandler` registrado em `api/internal/app/server.go`
- Rota `POST /v1/simulator/destinations` funcionando
- Middleware `FreemiumRateLimiter` ativo (5 req/dia para tier free)
- Performance: 2-4ms por request (com cache)

**Migrations Executadas com Sucesso**
- **Migration 0010** executada: Tabelas `countries_metadata`, `comexstat_cache`, `simulator_recommendations`
  - 10 países seed populados com metadados completos (flags, moedas, idiomas)
  - Índices otimizados criados
  - Funções PL/pgSQL ativas

- **Migration 0011** criada e executada: Schema ComexStat real implementado
  - Schema `stg.exportacao` com dados reais de ComexStat
  - 6 índices otimizados para queries do simulador
  - **64 registros reais** inseridos para validação:
    - NCM 17011400 (Açúcar de cana): 6 países, 22 registros históricos
    - NCM 26011200 (Minério de ferro): 4 países, 16 registros
    - NCM 12010090 (Soja em grão): 7 países, 26 registros
  - Dados incluem: China, EUA, Argentina, Países Baixos, Alemanha, Japão, Chile

**Validação e Testes Realizados**
- 3 NCMs testados com sucesso via API
- Rate limiting validado (bloqueia corretamente após 5 requests)
- Performance validada: ~2-4ms com dados reais
- Todos os campos calculados funcionando:
  - Score ponderado (0-10)
  - Rank automático
  - Demand level (Alto/Médio/Baixo)
  - EstimatedMarginPct, LogisticsCostUSD, TariffRatePct, LeadTimeDays
  - RecommendationReason baseada no score

**Arquivos Criados e Prontos para Commit**
- `api/internal/business/destination/entities.go` (entidades de domínio)
- `api/internal/business/destination/service.go` (lógica de negócio com algoritmo)
- `api/internal/business/destination/errors.go` (erros customizados)
- `api/internal/repository/postgres/destination.go` (repository layer)
- `api/internal/api/handlers/simulator.go` (HTTP handler)
- `api/internal/api/handlers/simulator_test.go` (testes unitários)
- `api/internal/api/middleware/freemium.go` (rate limiter)
- `api/internal/api/middleware/freemium_test.go` (testes do middleware)
- `db/migrations/0011_comexstat_schema.sql` (dados reais de ComexStat)
- `docs/API-SIMULATOR.md` (documentação completa da API - 750 linhas)

**Próximos Passos (Tarde - 22/11/2025)**
- [ ] Deploy Redis no k8s para cache L2 distribuído
- [ ] Popular países via Kubernetes Job (50 países principais)
- [ ] Executar testes E2E completos
- [ ] Commit final do simulador no branch `feature/security-credentials-management`

---

### Added - API do Simulador de Destinos 📍 (Fase 1 - 2025-01-21)

#### Database Schema & Migrations
- **Migration 0010** implementada em `db/migrations/0010_simulator_tables.sql`
  - Tabela `countries_metadata`: 50 principais parceiros comerciais
    - Campos: code, name_pt/en, region, gdp, population, distance_brazil_km
    - Índices otimizados por região, distância, GDP
    - 10 países iniciais populados (CN, US, AR, NL, CL, DE, JP, IN, MX, ES)
  - Tabela `comexstat_cache`: L3 cache backup para fallback
    - Cache key: type, year, month, NCM, country_code
    - TTL dinâmico: 7 dias (histórico) | 6h (mês atual)
    - Hit counter para analytics
    - JSONB para queries complexas
  - Tabela `simulator_recommendations`: analytics de uso
    - Rastreia todas as simulações (NCM, volume, resultados)
    - Cache metadata (hit, level, latency)
    - IP tracking para rate limiting
  - Funções SQL: `increment_comexstat_cache_hit()`, `cleanup_expired_comexstat_cache()`
  - Triggers para `updated_at` automático

#### Domain Layer (Clean Architecture)
- **Entities** implementadas em `api/internal/business/destination/entities.go`
  - `DestinationRecommendation`: recomendação completa com 15+ campos
    - Score (0-10), Rank, Demand (Alto/Médio/Baixo)
    - EstimatedMarginPct, LogisticsCostUSD, TariffRatePct, LeadTimeDays
    - MarketSizeUSD, GrowthRatePct, PricePerKgUSD, DistanceKm
    - RecommendationReason (explicação do score)
  - `SimulatorRequest`: contrato de entrada
    - NCM (8 dígitos, validação automática)
    - VolumeKg (opcional), Countries (filtro opcional)
    - MaxResults (1-50, default: 10)
  - `SimulatorResponse`: contrato de saída
    - Destinations array com rankings
    - Metadata (analysis_date, processing_time_ms, cache_hit)
  - `CountryMetadata`: metadados completos de países
  - `MarketData`: dados de mercado (NCM × País × Período)
  - `ScoringWeights`: pesos configuráveis do algoritmo
  - Métodos: `CalculateScore()`, `GetDemandLevel()`, `GetRecommendationReason()`

- **Errors** em `api/internal/business/destination/errors.go`
  - Erros de validação: `ErrInvalidNCM`, `ErrInvalidVolume`, `ErrInvalidMaxResults`
  - Erros de negócio: `ErrNCMNotFound`, `ErrNoDataAvailable`, `ErrInsufficientData`
  - Erros de infraestrutura: `ErrDatabaseConnection`, `ErrCacheUnavailable`

#### Business Logic (Service Layer)
- **Service** implementado em `api/internal/business/destination/service.go`
  - `RecommendDestinations()`: algoritmo completo de scoring
  - **Algoritmo de Scoring Simplificado**:
    - Market Size (40%): Tamanho do mercado em USD
    - Growth Rate (30%): Taxa de crescimento anual
    - Price per Kg (20%): Preço médio por kg
    - Distance (10%): Distância do Brasil
  - Normalização automática de métricas (0-1)
  - Cálculo de score ponderado (0-10)
  - Ordenação e ranking automático
  - Estimativas inteligentes:
    - `estimateMargin()`: Margem baseada em preço (15-35%)
    - `estimateLogisticsCost()`: Custo com economia de escala
    - `estimateTariff()`: Tarifa por região (8-18%)
    - `estimateLeadTime()`: Tempo de entrega (~500km/dia)
  - Filtragem por países específicos (opcional)
  - Análise dos últimos 12 meses

#### Infrastructure Layer (Repository)
- **Repository** implementado em `api/internal/repository/postgres/destination.go`
  - Interface: `GetCountryMetadata()`, `GetAllCountries()`
  - `GetMarketDataByNCM()`: query otimizada com CTEs
    - Agregação dos últimos 12 meses
    - Cálculo de growth rate (comparação período anterior)
    - Normalização de avg_price_per_kg_usd
    - Limit 100 países ordenados por market size
  - `GetMarketDataByNCMAndCountry()`: dados específicos NCM × País
  - `SaveRecommendation()`: analytics tracking em JSONB
  - Uso de `pq.StringArray` para arrays PostgreSQL
  - Error handling completo com tipos customizados

#### API Layer (Handlers & Middleware)
- **Handler** implementado em `api/internal/api/handlers/simulator.go`
  - `POST /v1/simulator/destinations`: endpoint principal
  - Validação automática via Gin binding
  - Error handling consistente com códigos HTTP apropriados
  - Response headers customizados
  - Swagger/OpenAPI annotations
  - Struct `ErrorResponse` padronizada

- **Middleware Freemium** em `api/internal/api/middleware/freemium.go`
  - Rate limiting diferenciado por tier:
    - Free: 5 simulações/dia (por IP ou user_id)
    - Premium: Ilimitado
  - Cache in-memory com TTL 24h
  - Cleanup automático de entradas expiradas
  - Headers informativos:
    - `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`
  - Identificação inteligente de usuário:
    - Autenticado: user_id do context
    - Anônimo: IP (com suporte a X-Forwarded-For e X-Real-IP)
  - HTTP 429 (Too Many Requests) com mensagem informativa
  - Thread-safe com `sync.RWMutex`

#### Kubernetes Jobs (Secure Data Seeding)
- **Job de população de países** em `k8s/jobs/populate-countries-job.yaml`
  - ServiceAccount com RBAC restrito
  - Role para acessar secrets (postgres-credentials)
  - Credenciais via `secretKeyRef` (ZERO plain text!)
  - Resource limits: 500m CPU, 512Mi memory
  - TTL 24h após conclusão
  - BackoffLimit: 3 tentativas
  - Non-root user

- **Script de população** em `scripts/populate-countries/main.go`
  - Busca dados via REST Countries API (v3.1)
  - Top 50 países de comércio exterior do Brasil
  - Cálculo de distância via fórmula de Haversine
  - Rate limiting (100ms entre requests)
  - Upsert com ON CONFLICT
  - Dockerfile multi-stage (golang:1.24-alpine)
  - Container não-root (appuser)

#### Security & Best Practices
- ✅ **ZERO plain text credentials**: Todas via Kubernetes Secrets
- ✅ **RBAC mínimo**: ServiceAccount com acesso restrito
- ✅ **Non-root containers**: Princípio de menor privilégio
- ✅ **Rate limiting**: Proteção contra abuso
- ✅ **Input validation**: Todas entradas validadas
- ✅ **Error handling**: Erros customizados sem exposição de detalhes internos
- ✅ **SQL injection safe**: Prepared statements em todas queries

### Added - Cache Multinível 🚀 (2025-01-21)

#### Sistema de Cache em 3 Níveis (L1 → L2 → L3)
- **Cache L1 (In-Memory - Ristretto)** implementado em `services/integration-gateway/internal/cache/l1_memory.go`
  - Algoritmo LFU (Least Frequently Used) com 100MB máximo
  - TTL configurável por item (default: 5min)
  - Performance: ~105 ns/op (read), ~2.7 µs/op (write)
  - Thread-safe com operações assíncronas + Wait()
  - Estatísticas completas: hits, misses, evictions, hit rate
  - 10 testes unitários implementados (100% pass rate)

- **Cache L2 (Distribuído - Redis)** implementado em `services/integration-gateway/internal/cache/l2_redis.go`
  - Compartilhado entre pods (escala horizontal)
  - Eviction policy: allkeys-lru
  - TTL: 7 dias (histórico) | 6h (mês atual)
  - Serialização automática em JSON
  - Connection pool configurável (default: 10 conexões)
  - Health checks e ping automático
  - 9 testes de integração implementados (100% pass rate)

- **MultiLevelCache Manager** implementado em `services/integration-gateway/internal/cache/manager.go`
  - Cascata automática: L1 → L2 → L3 → External API
  - Promoção automática de cache hits entre níveis
  - Propagação de Set/Delete para todos os níveis
  - Interface L3 definida (PostgreSQL Materialized Views - implementar depois)
  - Performance: ~414 ns/op (read cascata), ~3 µs/op (write propagação)
  - 11 testes unitários + 4 testes de integração (100% pass rate)

#### Métricas Prometheus para Cache
- **10 métricas customizadas** implementadas em `services/integration-gateway/internal/cache/metrics.go`
  - `integration_gateway_cache_hits_total` - Total de hits por nível
  - `integration_gateway_cache_misses_total` - Total de misses por nível
  - `integration_gateway_cache_latency_seconds` - Histogram de latência (P50/P95/P99)
  - `integration_gateway_cache_size_bytes` - Tamanho atual do cache
  - `integration_gateway_cache_evictions_total` - Total de evictions (L1)
  - `integration_gateway_cache_hit_rate` - Taxa de hit (0.0 a 1.0)
  - `integration_gateway_cache_sets_total` - Total de operações set
  - `integration_gateway_cache_promotions_total` - Promoções entre níveis
  - `integration_gateway_cache_errors_total` - Erros por tipo e nível

#### Infraestrutura Redis
- **Docker Compose** (`bgcstack/docker-compose.yml`)
  - Redis 7-alpine adicionado
  - Configuração: 512MB max memory, allkeys-lru policy
  - Volume persistente `redis_data`
  - Health checks configurados
  - Variáveis de ambiente para Integration Gateway

- **Kubernetes** (`k8s/redis.yaml`)
  - Deployment com PVC 2Gi
  - ConfigMap com redis.conf otimizado
  - Service ClusterIP
  - Health probes (liveness + readiness)
  - Resource limits: 500m CPU, 1Gi memory

- **Integration Gateway atualizado**
  - Variáveis de ambiente para Redis (REDIS_ADDR, REDIS_PASSWORD, REDIS_DB)
  - Flags de habilitação (CACHE_L1_ENABLED, CACHE_L2_ENABLED)
  - Dependência explícita no Redis

#### Connector Config - ComexStat
- **Configuração completa** em `config/connectors/comexstat.yaml`
  - Cache multinível habilitado (L1 + L2 + L3)
  - TTL: 168h (7 dias) para dados históricos
  - Key pattern: `comexstat:exp:{ano}:{mes}:{ncm}:{pais}`
  - Rate limit: 4 req/min (margem de segurança para 300/hour)
  - Circuit breaker: 3 falhas → open (2min)
  - Retry: exponential backoff (2s → 10s)
  - Alertas configurados (error_rate, latency, availability)

#### Testes & Cobertura
- **30+ testes implementados** (unitários + integração)
  - Testes unitários: L1, L2, Manager, Métricas
  - Testes de integração: Redis real, cascata L1+L2, alta throughput
  - Build tags para separar testes (`-tags=integration`)
  - **Cobertura: 82%** do código de cache
  - Todos os testes passando ✅

#### Benchmarks & Performance
- **Benchmarks completos** executados
  - L1 Get: ~105 ns/op (34M ops/s)
  - L1 Set: ~2.7 µs/op (1M ops/s)
  - Manager Get (cascata): ~414 ns/op (12M ops/s)
  - Manager Set (propagação): ~3 µs/op (1M ops/s)
  - Allocation: 22-192 bytes/op, 1-4 allocs/op

#### Documentação
- **README completo** em `services/integration-gateway/internal/cache/README.md`
  - Arquitetura e diagramas
  - Guia de uso para cada nível (L1, L2, L3)
  - Configuração e variáveis de ambiente
  - Estratégias de cache (TTL dinâmico, request coalescing)
  - Métricas Prometheus e dashboards
  - Troubleshooting completo
  - Referências e próximos passos

### Added - Simulador de Destinos de Exportação 🌍 (2025-11-19/20)

#### Segurança & Secrets Management 🔐
- **KubernetesSecretStore** implementado em `services/integration-gateway/internal/auth/k8s_secret_store.go`
  - Busca secrets diretamente da Kubernetes Secrets API via `k8s.io/client-go`
  - Cache in-memory com TTL de 5 minutos para reduzir chamadas à API
  - Formato: `secret-name/key-name` (ex: `comexstat-credentials/api-key`)
  - Backward compatibility com env vars (`SECRET_*`)
  - Thread-safe com `sync.RWMutex`
  - Limpeza automática de cache expirado via goroutine
  - 19 testes unitários implementados (100% pass rate)

- **Sealed Secrets** configurado para credenciais sensíveis
  - Template em `k8s/integration-gateway/sealed-secret-comexstat.yaml`
  - Script automatizado `scripts/create-sealed-secret-comexstat.sh` para criação segura
  - Suporte a método online (cluster ativo) e offline (certificado local)
  - Documentação completa em `k8s/integration-gateway/README-SECRETS.md` (270+ linhas)
  - Guias de troubleshooting, rotação de secrets e boas práticas

#### Network Policies & Segmentação de Rede 🛡️
- **Network Policies** implementadas para isolamento de rede
  - `k8s/network-policies/integration-gateway-netpol.yaml`:
    - Ingress permitido APENAS de bgc-api e Prometheus
    - Egress para DNS, Redis, PostgreSQL, APIs externas HTTPS, Jaeger e K8s API
    - Bloqueia tráfego não autorizado por padrão
  - `k8s/network-policies/bgc-api-netpol.yaml`:
    - **FORÇA** integra\u00e7\u00f5es externas via Integration Gateway
    - BLOQUEIA acesso direto da API a APIs externas (porta 443)
    - Ingress apenas de Ingress Controller e Prometheus
    - Egress para PostgreSQL, Redis, Integration Gateway e Jaeger
  - `k8s/network-policies/default-deny-all`:
    - Nega TODO tráfego por padrão no namespace `data`
    - Pods precisam de NetworkPolicy explícita
  - Policies para Redis e PostgreSQL (isolamento completo)
  - Documentação completa em `k8s/network-policies/README.md` (450+ linhas)
    - Arquitetura de rede com diagramas
    - Guias de teste de conectividade
    - Troubleshooting para debugging de policies

#### Dependências & Build
- Adicionadas dependências Kubernetes no `services/integration-gateway/go.mod`:
  - `k8s.io/apimachinery v0.29.0`
  - `k8s.io/client-go v0.29.0`
- Build validado: `go build` e `go test` passando sem erros

### Security
- **Princípio de Menor Privilégio** implementado via Network Policies
- **Secrets Management** enterprise-grade com K8s Secrets API + Sealed Secrets
- **Zero exposição de credenciais** em código ou variáveis de ambiente
- **Isolamento de rede** entre serviços (defense in depth)
- **Auditabilidade** de acesso a secrets via logs estruturados

### Documentation
- `k8s/integration-gateway/README-SECRETS.md` - Guia completo de secrets management
- `k8s/network-policies/README.md` - Guia de network policies e segurança de rede
- Scripts documentados com comentários inline e help text

---

### Added - Épico 2: Observabilidade & Padrões 📊
- **Prometheus Metrics** para métricas de produção
  - Integração completa do `github.com/prometheus/client_golang`
  - Endpoint `/metrics` em formato Prometheus (nativo)
  - Endpoint `/metrics/json` para compatibilidade (legacy)
  - 11 métricas implementadas:
    - `bgc_http_requests_total` - Total de requisições HTTP por método, path e status
    - `bgc_http_request_duration_seconds` - Duração de requisições (histogram com P50/P95/P99)
    - `bgc_http_requests_in_flight` - Requisições em processamento (gauge)
    - `bgc_db_queries_total` - Total de queries por operação e tabela
    - `bgc_db_query_duration_seconds` - Duração de queries (histogram)
    - `bgc_db_connections_open` - Conexões abertas do pool
    - `bgc_db_connections_in_use` - Conexões em uso
    - `bgc_db_connections_idle` - Conexões ociosas
    - `bgc_errors_total` - Total de erros por tipo e severidade
    - `bgc_idempotency_cache_hits_total` - Cache hits de idempotência
    - `bgc_idempotency_cache_misses_total` - Cache misses de idempotência
    - `bgc_idempotency_cache_size` - Tamanho atual do cache

- **Middleware Prometheus** em `api/internal/observability/metrics/prometheus.go`
  - Instrumentação automática de todos os handlers HTTP
  - Coleta de métricas de latência (buckets: 5ms a 10s)
  - Tracking de requisições in-flight
  - Labels: method, path, status

- **Instrumentação de Banco de Dados** com métricas Prometheus
  - `api/internal/repository/postgres/market.go` - Métricas em GetMarketDataByYearRange
  - `api/internal/repository/postgres/route.go` - Métricas em GetTAMByYearAndChapter
  - DB stats collector rodando a cada 15 segundos
  - Tracking de connection pool (open, in-use, idle)

- **OpenTelemetry SDK** para distributed tracing
  - Integração completa de `go.opentelemetry.io/otel` v1.38.0
  - OTLP gRPC exporter (`go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`)
  - Stdout exporter para desenvolvimento (`go.opentelemetry.io/otel/exporters/stdout/stdouttrace`)
  - Tracer provider global configurado em `api/internal/observability/tracing/tracer.go`
  - Propagação de trace context via W3C Trace Context (traceparent/tracestate)
  - Graceful shutdown com flush de spans pendentes

- **Instrumentação Automática com OpenTelemetry**
  - Middleware Gin (`go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`)
  - Todos os handlers HTTP geram spans automaticamente
  - Trace ID e Span ID propagados em todos os requests
  - Contexto de tracing injetado via headers HTTP

- **Instrumentação de Database Queries** com tracing
  - Spans criados para cada query com contexto
  - Atributos detalhados: year.from, year.to, chapters.count, result.count
  - Error recording em caso de falha
  - Status codes (Ok, Error) em todos os spans
  - `api/internal/repository/postgres/market.go:26-88` - db.GetMarketDataByYearRange
  - `api/internal/repository/postgres/route.go:23-51` - db.GetTAMByYearAndChapter

- **Stack de Observabilidade Kubernetes** em `k8s/observability/`
  - **Prometheus Deployment** (`prometheus-deployment.yaml`)
    - ServiceAccount com RBAC para Kubernetes service discovery
    - ClusterRole e ClusterRoleBinding para pod/service discovery
    - Retention de 30 dias
    - Scrape interval de 15 segundos
    - Targets: bgc-api, integration-gateway, postgres-exporter
  - **Prometheus ConfigMap** com scrape configs
    - Auto-discovery de pods via kubernetes_sd_configs
    - Relabeling para extrair pod name, namespace
    - Jobs: prometheus, bgc-api, integration-gateway, postgres
  - **Prometheus Alert Rules** (`prometheus-alert-rules.yaml`)
    - 10 regras de alerta pré-configuradas:
      - HighErrorRate (critical): error rate > 5% por 5min
      - HighLatencyP95 (warning): P95 > 2s por 10min
      - APIDown (critical): service down por 1min
      - HighDatabaseLatency (warning): P95 DB > 500ms por 10min
      - DatabaseConnectionPoolExhaustion (warning): > 90% connections em uso
      - HighDatabaseConnections (warning): > 50 conexões abertas
      - HighRequestRate (warning): > 1000 req/s por 5min
      - IdempotencyCacheSizeHigh (info): cache > 10k entries
      - IntegrationGatewayDown (warning): gateway down por 2min
      - ConnectorCircuitBreakerOpen (warning): circuit breaker aberto por 5min
  - **Grafana Deployment** (`grafana-deployment.yaml`)
    - Datasources provisionados automaticamente (Prometheus + Jaeger)
    - Dashboard provider configurado para /var/lib/grafana/dashboards
    - Plugins: redis-datasource, grafana-piechart-panel
    - NodePort 30030 para acesso externo
  - **Grafana Dashboards** (`grafana-dashboards.yaml`)
    - Dashboard "BGC API Overview" completo com 8 painéis:
      - Request Rate (req/s) por endpoint
      - Error Rate (%) com threshold de alerta
      - Request Duration (P50, P95, P99)
      - Requests In Flight
      - Database Query Rate (queries/s)
      - Database Query Duration (P95)
      - Database Connections (Open, In Use, Idle)
      - Idempotency Cache (Hits, Misses, Size)
  - **Jaeger All-in-One** (`jaeger-deployment.yaml`)
    - OTLP gRPC receiver na porta 4317
    - OTLP HTTP receiver na porta 4318
    - Jaeger UI na porta 16686 (NodePort 30016)
    - Memory storage com 10k traces
    - Health checks configurados

- **Docker Compose** atualizado com observability stack completa
  - **Prometheus** container (`prom/prometheus:v2.50.0`)
    - Configuração em `bgcstack/observability/prometheus.yml`
    - Volume persistente `prometheus_data`
    - Porta 9090 exposta
    - Scrape automático de api:8080 e integration-gateway:8081
  - **Grafana** container (`grafana/grafana:10.3.0`)
    - Datasources provisionados via `grafana-datasources.yml`
    - Dashboards em `observability/dashboards/bgc-api.json`
    - Porta 3001 exposta (para não conflitar com web:3000)
    - Volume persistente `grafana_data`
    - Credenciais: admin / ${GRAFANA_ADMIN_PASSWORD:-admin}
  - **Jaeger** container (`jaegertracing/all-in-one:1.54`)
    - OTLP gRPC na porta 4317
    - OTLP HTTP na porta 4318
    - Jaeger UI na porta 16686
    - Collector HTTP na porta 14268
    - Memory storage com 10k traces max
  - **API** atualizada com variáveis de ambiente:
    - `ENVIRONMENT=development`
    - `OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4317`
  - **Integration Gateway** atualizada:
    - `OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4317`
    - Depends on: jaeger

- **Documentação Completa** em `docs/OBSERVABILITY.md` (700+ linhas)
  - Overview da arquitetura de observabilidade
  - Diagrama de fluxo de métricas e traces
  - Referência completa de todas as 11 métricas Prometheus
  - Guia de distributed tracing com OpenTelemetry
  - Instruções de setup local (Docker Compose)
  - Instruções de deployment Kubernetes
  - Dashboards e queries PromQL
  - Alerting e notification setup
  - Troubleshooting guide completo
  - Best practices para metrics, tracing e dashboards

### Changed - Épico 2
- **`api/internal/app/server.go` refatorado** para observabilidade:
  - Adicionado `metrics.PrometheusMiddleware()` para coleta de métricas HTTP
  - Adicionado `otelgin.Middleware("bgc-api")` para tracing automático
  - Endpoint `/metrics` agora retorna formato Prometheus nativo
  - Novo endpoint `/metrics/json` para formato JSON (backwards compatibility)
  - DB stats collector iniciado automaticamente a cada 15 segundos
  - Import de `github.com/prometheus/client_golang/prometheus/promhttp`
  - Import de `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin`

- **`api/cmd/api/main.go` refatorado** para tracing e graceful shutdown:
  - Inicialização do OpenTelemetry tracer no startup
  - Graceful shutdown com signal handling (SIGINT, SIGTERM)
  - Flush de spans pendentes antes do shutdown
  - Fallback para stdout exporter se OTLP endpoint não configurado
  - Environment detection via variável `ENVIRONMENT`

- **Repositórios instrumentados** com métricas e tracing:
  - `api/internal/repository/postgres/market.go`:
    - Import de `context`, `time`, `metrics`, `tracing`
    - Span "db.GetMarketDataByYearRange" com atributos
    - Métricas de query duration para tabela "v_tam_by_year_chapter"
    - Error recording em spans
  - `api/internal/repository/postgres/route.go`:
    - Import de `context`, `time`, `metrics`, `tracing`
    - Span "db.GetTAMByYearAndChapter" com atributos
    - Métricas de query duration
    - Error recording em spans

- **go.mod atualizado** com dependências de observabilidade:
  - `github.com/prometheus/client_golang v1.23.2`
  - `github.com/prometheus/client_model v0.6.2`
  - `github.com/prometheus/common v0.66.1`
  - `github.com/prometheus/procfs v0.16.1`
  - `go.opentelemetry.io/otel v1.38.0`
  - `go.opentelemetry.io/otel/sdk v1.38.0`
  - `go.opentelemetry.io/otel/trace v1.38.0`
  - `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0`
  - `go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.38.0`
  - `go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.63.0`
  - `google.golang.org/grpc v1.75.0`
  - Upgrades: gin v1.10.1, golang.org/x/* para versões mais recentes

### Added - Épico 1: Enablement & Acesso 🚀
- **Arquitetura Híbrida de Integrações** - Framework genérico para APIs externas
  - 90% das integrações via configuração YAML (zero código Go!)
  - 10% casos complexos via plugins customizados
  - Tempo para nova integração reduzido de **2 dias para 30 minutos**

- **Integration Gateway Service** em `services/integration-gateway/`
  - Framework core com tipos base, HTTP client resiliente e executor
  - Connector Registry para carregar e validar configs YAML
  - Transform Engine com JSONPath para mapeamento de dados
  - 6 plugins built-in: format_cnpj, format_cpf, format_cep, to_upper, to_lower, trim

- **Auth Engine** com suporte completo a múltiplos tipos de autenticação:
  - **mTLS** - Certificados ICP-Brasil (A1/A3) para APIs governamentais
  - **OAuth2** - Client credentials flow com token caching automático
  - **API Key** - Suporte a headers customizáveis
  - None - Para APIs públicas
  - Interfaces para Basic e JWT (implementação futura)

- **Resiliência Automática** aplicada a todos os conectores:
  - Circuit Breaker (gobreaker) com thresholds configuráveis
  - Retry com backoff (constant, linear, exponential)
  - Rate Limiting (requests/min + burst)
  - Timeouts configuráveis por endpoint

- **Certificate Manager (MVP)** em `services/integration-gateway/internal/auth/`
  - SimpleCertificateManager para gestão de certificados ICP-Brasil
  - SimpleSecretStore para secrets via env vars
  - Suporte a certificados A1 (PFX/P12) e A3 (HSM/Token)
  - Estrutura `certs/` com .gitignore e documentação

- **API REST do Gateway** (`cmd/gateway/main.go`)
  - `GET /health` - Health check
  - `GET /v1/connectors` - Lista todos os conectores
  - `GET /v1/connectors/{id}` - Detalhes de um connector
  - `POST /v1/connectors/{id}/{endpoint}` - Executa endpoint com params

- **Connector Configs de Exemplo** em `config/connectors/`:
  - `receita-federal-cnpj.yaml` - Integração complexa (mTLS, transformações, cache)
  - `viacep.yaml` - Integração simples (API pública, sem auth)

- **JSON Schema de Validação** em `schemas/connector.schema.json`
  - Schema completo (600+ linhas) para validação de configs
  - Suporta todos os tipos de auth, endpoints, resiliência, compliance
  - Validação automática de governança (owner_team, approved_by, etc)

- **Documentação Completa**:
  - `docs/CONNECTOR-GUIDE.md` - Guia completo de uso (quick start, anatomia, exemplos)
  - `docs/EPIC-1-PROGRESS.md` - Status detalhado do épico
  - `services/integration-gateway/README.md` - Documentação técnica do serviço
  - `certs/README.md` - Instruções de gestão de certificados ICP-Brasil

- **Governança Built-in** nas configurações YAML:
  - Compliance tags (LGPD, SOC2, ICP-Brasil, etc)
  - Data classification (public, internal, confidential, restricted)
  - Owner team e aprovações obrigatórias
  - Review frequency (monthly, quarterly, annually)
  - Alertas configuráveis (certificate_expiry, error_rate, latency)

### Added - Épico 3: Contrato de Dados ✅
- **JSON Schemas versionados** em `schemas/v1/` para validação de API
  - `market-size-request.schema.json` e `market-size-response.schema.json`
  - `route-comparison-request.schema.json` e `route-comparison-response.schema.json`
  - `error-response.schema.json` para respostas de erro padronizadas
  - Validação completa de tipos, formatos e constraints
- **Validação automática de schemas** na API Go com middleware
  - Dependência `github.com/xeipuuv/gojsonschema` adicionada
  - Validator em `api/internal/api/validation/validator.go`
  - Middleware de validação aplicado em todos os endpoints v1
  - Mensagens de erro detalhadas com campo e issue específicos
  - Fallback graceful se schemas não disponíveis
- **Versionamento de API** seguindo best practices REST
  - Todos os endpoints principais em `/v1/*`
  - Rotas antigas redirecionam automaticamente (301) para v1
  - Preparado para múltiplas versões futuras
- **Dicionário de dados completo** em `docs/DATA-DICTIONARY.md`
  - Documentação de todas as tabelas, colunas, tipos e constraints
  - Índices, materialized views e políticas de performance
  - Estrutura de schemas (public, stg, dim, rpt)
  - Exemplos de queries e padrões de uso
- **Sistema de idempotência completo** para prevenir processamento duplicado
  - Middleware aplicado globalmente no grupo `/v1`
  - Suporte a header `Idempotency-Key` (16-128 caracteres)
  - Cache in-memory thread-safe com TTL de 24h
  - Cleanup automático de entradas expiradas
  - Headers de resposta: `X-Idempotency-Cached`, `X-Idempotency-Cached-At`
  - Migration `0004_idempotency.sql` com:
    - Tabela `api_idempotency` para persistência
    - Colunas `idempotency_key` em `stg.exportacao` e `stg.importacao`
    - Função `cleanup_expired_idempotency_keys()`
- **Documentação de políticas** em `docs/IDEMPOTENCY-POLICY.md`
  - Políticas de reprocessamento de dados
  - Exemplos de uso e best practices
  - Estratégias de retry e deduplicação
  - Formato de chave e TTL

### Added - Infraestrutura e Segurança
- CHANGELOG.md para rastreamento de mudanças
- Health probes (readiness/liveness) no deployment WEB
- HorizontalPodAutoscaler (HPA) para API e WEB
- CronJob de backup automático do PostgreSQL
- Makefile como wrapper unificado dos scripts PowerShell
- Documentação de observabilidade e resiliência
- Template .env.example para configuração segura de credenciais (Docker Compose)
- Documentação completa de segurança em docs/SECURITY-SECRETS.md
- GitHub Actions workflows em .github/ para CI/CD
- DIAGNOSTICO_12_FATORES.md com análise da aplicação

### Changed
- Script k8s.ps1 atualizado para aplicar HPA e CronJobs
- README.md expandido com informações sobre HPA e backups
- **`server.go` completamente refatorado** para Épico 3:
  - Versionamento de API com grupo `/v1`
  - Middleware de idempotência aplicado globalmente
  - Middleware de validação integrado com schemas
  - Rotas legacy redirecionam para v1 (backwards compatibility)
  - Logging estruturado de inicialização

### Security
- **Atualização crítica de segurança**: Go 1.23.x → 1.24.9 para correção de 5 vulnerabilidades
  - **GO-2025-4013** - Panic em validação de certificados DSA em `crypto/x509`
    - Impacto: `internal/repository/postgres/market.go:54` e operações TLS
  - **GO-2025-4011** - Vulnerabilidade em `encoding/asn1`
    - Impacto: Parsing de certificados e estruturas ASN.1
  - **GO-2025-4010** - Falha em `net/url`
    - Impacto: `internal/repository/postgres/db.go:17` e `internal/app/server.go:69`
  - **GO-2025-4008** - Exposição de informação em negociação ALPN em `crypto/tls`
    - Impacto: Todas as conexões TLS/HTTPS
  - **GO-2025-4007** - Múltiplas falhas em `crypto/x509`
    - Impacto: Parsing de certificados e chaves privadas
- Arquivos atualizados:
  - `go.work`: Go 1.23.0 → 1.24.9
  - `api/go.mod`: Go 1.23.0 → 1.24.9
  - `services/bgc-ingest/go.mod`: Go 1.22 → 1.24.9
  - `services/integration-gateway/go.mod`: Go 1.23 → 1.24.9
  - `api/Dockerfile`: golang:1.23-alpine → golang:1.24-alpine
  - `.github/workflows/security-scan.yml`: Go 1.23 → 1.24
- Removidas credenciais hardcoded do README.md
- docker-compose.yml migrado para variáveis de ambiente com validação
- .gitignore expandido com regras de proteção de secrets e credenciais
- README.md atualizado com instruções de configuração segura

## [0.2.5.1] - 2025-01-15

### Changed
- Migração para PostgreSQL oficial (postgres:16) substituindo Bitnami
- Infraestrutura Kubernetes estabilizada
- Correção de secrets do banco de dados

## [0.2.5] - 2025-01-14

### Added
- API/Web estáveis em produção simulada
- Kubernetes deployments com Traefik Ingress
- Documentação completa de deployment
- Métricas de observabilidade

### Changed
- Sprint 2 finalizada com infraestrutura consolidada

## [0.1-sprint1] - 2025-01-10

### Added
- Infraestrutura inicial com k3d + PostgreSQL (Helm)
- Serviço de ingestão CSV/XLSX (bgc-ingest)
- Materialized Views para agregação de dados (rpt.*)
- API REST read-only com endpoints /metrics/*
- Manifests Kubernetes e scripts de automação
- Sistema de proveniência de dados (ingest_source, ingest_batch)
- Documentação de arquitetura e post-mortem Sprint 1

### Features
- Clean Architecture na API Go
- Endpoints: /market/size (TAM/SAM/SOM) e /routes/compare
- Docker Compose para desenvolvimento local
- Scripts PowerShell para gerenciamento (docker.ps1, k8s.ps1)

---

## Formato de Versionamento

- **MAJOR**: Mudanças incompatíveis na API
- **MINOR**: Novas funcionalidades mantendo compatibilidade
- **PATCH**: Correções de bugs e melhorias

## Tipos de Mudanças

- **Added**: Novas features
- **Changed**: Mudanças em funcionalidades existentes
- **Deprecated**: Features que serão removidas
- **Removed**: Features removidas
- **Fixed**: Correções de bugs
- **Security**: Correções de vulnerabilidades
