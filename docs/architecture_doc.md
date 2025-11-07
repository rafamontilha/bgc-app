# Arquitetura do Sistema BGC Analytics

**Versão:** 4.0
**Última atualização:** Novembro 2025
**Status:** Production-Ready with Full Observability & Integration Framework

## 📋 Visão Geral

O BGC Analytics é um sistema de analytics para dados de exportação brasileira, construído com **Clean Architecture (Hexagonal Architecture)** e arquitetura cloud-native para execução em ambiente Kubernetes local (k3d) durante desenvolvimento e Docker Compose para desenvolvimento rápido.

### Objetivos do Sistema
- **Performance:** Consultas analíticas rápidas via Materialized Views (<500ms p95)
- **Manutenibilidade:** Código modular seguindo Clean Architecture
- **Testabilidade:** Separação clara de camadas com dependency injection
- **Observabilidade:** Métricas, tracing e logs completos (Prometheus, Jaeger, OTel)
- **Integração:** Framework para 30-min integration time (vs 2 dias antes)
- **Resiliência:** Circuit breaker, retry, rate limiting automáticos
- **Governança:** Validação de contratos, idempotency, versionamento
- **Segurança:** Go 1.24.9, mTLS, OAuth2, vulnerabilidades corrigidas
- **Desenvolvimento ágil:** Ambiente local reproducível
- **Escalabilidade:** Preparado para migração cloud futura

---

## 🏗️ Arquitetura de Alto Nível

```
┌────────────────────────────────────────────────────────────────────────────────┐
│                         k3d Cluster / Docker Compose                           │
│  ┌───────────────┐  ┌───────────────┐  ┌─────────────────┐ ┌──────────────┐  │
│  │   bgc-web     │  │   bgc-api     │  │ bgc-postgres    │ │ integration- │  │
│  │  (Next.js)    │  │   (Go API)    │  │  (PostgreSQL)   │ │   gateway    │  │
│  │               │  │               │  │                 │ │   (Go)       │  │
│  │ ┌───────────┐ │  │ ┌───────────┐ │  │ ┌─────────┐   │ │ ┌──────────┐ │  │
│  │ │ React UI  │ │  │ │/v1/market │ │  │ │PostgreSQL│  │ │ │Connectors│ │  │
│  │ │ SSR/SSG   │ │  │ │/v1/routes │ │  │ │ Database │  │ │ │ mTLS     │ │  │
│  │ │ Rewrites  │─┼──┼▶│ /healthz  │ │  │ │          │  │ │ │ OAuth2   │ │  │
│  │ └───────────┘ │  │ │ /metrics  │ │  │ │  ┌────┐  │  │ │ │ API Key  │ │  │
│  │               │  │ │OTel Traces│ │  │ │  │MVs │  │  │ │ └──────────┘ │  │
│  │ Port: 3000    │  │ └───────────┘ │  │ │  └────┘  │  │ │ Port: 8081   │  │
│  └───────────────┘  └───────────────┘  │ └─────────┘   │ └──────────────┘  │
│         │                   │           │        ▲       │        │          │
│         │                   └───────────┼────────┘       │        │          │
│         │                               │                │        │          │
│         │                           ┌───┴────────────────┴────────┘          │
│    ┌────▼──────┐                   │                                         │
│    │  Ingress  │              ┌────▼─────────────────────────────────┐       │
│    │  Traefik  │              │     Observability Stack              │       │
│    └───────────┘              │  ┌────────────┐  ┌────────────┐     │       │
│         │                      │  │ Prometheus │  │  Grafana   │     │       │
│         │                      │  │  Metrics   │  │ Dashboards │     │       │
│    ┌────▼──────┐              │  │  :9090     │  │   :3001    │     │       │
│    │ CronJobs  │              │  └────────────┘  └────────────┘     │       │
│    │ - Backup  │              │  ┌──────────────────────────┐       │       │
│    │ - Refresh │              │  │         Jaeger           │       │       │
│    └───────────┘              │  │   Distributed Tracing    │       │       │
│                                │  │        :16686            │       │       │
│                                │  └──────────────────────────┘       │       │
│                                └─────────────────────────────────────┘       │
└────────────────────────────────────────────────────────────────────────────────┘
              │
         ┌────▼────────┐
         │   Browser   │
         │ web.bgc.local (K8s) / localhost:3000 (Docker)               │
         └─────────────┘
```

---

## 🔧 Componentes Principais

### 1. BGC Web (Next.js 15) - Frontend Moderno

**Responsabilidade:** Interface de usuário para visualização e análise de dados de mercado

**Tecnologias:**
- **Framework:** Next.js 15.1 (App Router)
- **Runtime:** React 19, TypeScript 5
- **Styling:** Tailwind CSS
- **Deploy:** Docker (production build), Node.js standalone

**Páginas Principais:**
```
GET /                    # Dashboard TAM/SAM/SOM
GET /routes              # Comparação de rotas comerciais
GET /api/health          # Health check interno
```

**Arquitetura Next.js:**
- **Server-Side Rendering (SSR):** Páginas renderizadas no servidor
- **API Routes:** `/api/health` para health check
- **Rewrites Internos:** Proxy transparente para Go API
  ```typescript
  // next.config.ts
  rewrites() {
    return [
      { source: '/market/:path*', destination: 'http://bgc-api:8080/market/:path*' },
      { source: '/routes/:path*', destination: 'http://bgc-api:8080/routes/:path*' },
      { source: '/healthz', destination: 'http://bgc-api:8080/healthz' }
    ];
  }
  ```

**Features:**
- Dashboard interativo com gráficos e tabelas
- Export de dados para CSV
- Validação de formulários
- Estado compartilhado via React hooks
- Responsive design com Tailwind

**Estrutura de Diretórios:**
```
web-next/
├── app/
│   ├── page.tsx              # Dashboard TAM/SAM/SOM
│   ├── routes/
│   │   └── page.tsx          # Comparação de rotas
│   ├── api/
│   │   └── health/
│   │       └── route.ts      # Health check endpoint
│   └── layout.tsx            # Layout principal
├── components/               # Componentes React reutilizáveis
├── public/                   # Assets estáticos
├── next.config.ts            # Configuração (rewrites)
├── tailwind.config.ts        # Configuração Tailwind
├── tsconfig.json             # TypeScript config
└── Dockerfile                # Build para produção
```

### 2. BGC API (Go) - Clean Architecture

**Responsabilidade:** API REST para consultas analíticas com arquitetura hexagonal

**Tecnologias:**
- **Runtime:** Go 1.24.9+
- **Framework HTTP:** Gin (gin-gonic/gin)
- **Database Driver:** lib/pq (PostgreSQL)
- **Configuration:** gopkg.in/yaml.v3
- **Observability:** Prometheus client_golang, OpenTelemetry SDK
- **Validation:** JSON Schema (gojsonschema)
- **Deploy:** Kubernetes Deployment ou Docker Compose
- **Port:** 8080 (internal and external)

**Endpoints Atuais:**
```
# Health & Documentation
GET /health, /healthz           # Health check com status de config
GET /docs                       # Documentação Redoc
GET /openapi.yaml              # Especificação OpenAPI

# Observability
GET /metrics                    # Prometheus metrics (formato nativo)
GET /metrics/json              # Métricas JSON (legacy)

# API v1 (com JSON Schema validation)
GET /v1/market/size            # Cálculo de TAM/SAM/SOM
  ?metric=TAM|SAM|SOM
  &year_from=YYYY
  &year_to=YYYY
  &ncm_chapter=XX
  &scenario=base|aggressive
  Header: Idempotency-Key (optional)

GET /v1/routes/compare         # Comparação de rotas comerciais
  ?from=USA
  &alts=CHN,ARE,IND
  &ncm_chapter=XX
  &year=YYYY
  &tariff_scenario=base|tarifa10
  Header: Idempotency-Key (optional)

# Legacy endpoints (redirect 301 to /v1/*)
GET /market/size               # → /v1/market/size
GET /routes/compare            # → /v1/routes/compare
```

**Arquitetura Hexagonal (Camadas):**
```
┌─────────────────────────────────────────────────────────────┐
│                    cmd/api/main.go                          │
│                    (Entry Point)                            │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              internal/app/server.go                         │
│          (Dependency Injection & Wiring)                    │
└─────────────────────────────────────────────────────────────┘
           │                │                │
           ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  Handlers    │  │  Middleware  │  │   Config     │
│  (HTTP)      │  │  (CORS, Log) │  │  (YAML)      │
└──────────────┘  └──────────────┘  └──────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────┐
│              internal/business/                             │
│         (Domain Layer - Business Logic)                     │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                 │
│  │  Market  │  │  Route   │  │  Health  │                 │
│  │ Service  │  │ Service  │  │ Service  │                 │
│  └──────────┘  └──────────┘  └──────────┘                 │
│       │              │              │                       │
│       ▼              ▼              ▼                       │
│  Repository    Repository      (no repo)                   │
│  Interface     Interface                                    │
└─────────────────────────────────────────────────────────────┘
           │              │
           ▼              ▼
┌─────────────────────────────────────────────────────────────┐
│         internal/repository/postgres/                       │
│      (Infrastructure - Database Access)                     │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │   Market     │  │    Route     │                        │
│  │ Repository   │  │ Repository   │                        │
│  └──────────────┘  └──────────────┘                        │
│          │                 │                                │
│          └─────────────────┘                                │
│                    │                                        │
│                    ▼                                        │
│            PostgreSQL Database                              │
└─────────────────────────────────────────────────────────────┘
```

**Estrutura de Diretórios:**
```
api/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── config/                  # Configuração
│   │   └── config.go           # LoadConfig, LoadPartnerWeights, LoadTariffScenarios
│   ├── business/                # Domínios (Lógica de Negócio)
│   │   ├── market/
│   │   │   ├── entities.go    # MarketItem, MarketSizeRequest/Response
│   │   │   ├── repository.go  # Interface Repository
│   │   │   └── service.go     # CalculateMarketSize (TAM/SAM/SOM)
│   │   ├── route/
│   │   │   ├── entities.go    # RouteCompareRequest/Response, RouteItem
│   │   │   ├── repository.go  # Interface Repository
│   │   │   └── service.go     # CompareRoutes (partner weights + tariffs)
│   │   └── health/
│   │       └── service.go     # GetHealthStatus
│   ├── repository/              # Implementações de Persistência
│   │   └── postgres/
│   │       ├── db.go           # MustConnect (connection setup)
│   │       ├── market.go       # GetMarketDataByYearRange
│   │       └── route.go        # GetTAMByYearAndChapter
│   ├── api/                     # Camada HTTP
│   │   ├── handlers/
│   │   │   ├── health.go       # GET /health, /healthz
│   │   │   ├── market.go       # GET /market/size
│   │   │   └── route.go        # GET /routes/compare
│   │   └── middleware/
│   │       ├── cors.go         # CORS middleware
│   │       └── metrics.go      # Request ID, logging, metrics
│   └── app/
│       └── server.go           # NewServer (wiring), Run
├── config/                      # Arquivos de configuração YAML
│   ├── partners_stub.yaml
│   ├── tariff_scenarios.yaml
│   ├── scope.yaml
│   └── som.yaml
├── Dockerfile
├── go.mod
└── openapi.yaml
```

**Princípios Aplicados:**

1. **Separation of Concerns**: Cada camada tem responsabilidade única
   - Handlers: HTTP request/response
   - Services: Business logic
   - Repositories: Data access

2. **Dependency Inversion**: 
   - Services dependem de interfaces de Repository
   - Implementações concretas injetadas em runtime

3. **Testability**:
   - Services testáveis via mock repositories
   - Business logic isolada de HTTP e DB

4. **Framework Independence**:
   - Lógica de negócio não depende de Gin
   - Fácil migração para outro framework HTTP

5. **Clean Code**:
   - Packages pequenos e focados
   - Código autodocumentado
   - Sem comentários desnecessários

### 3. PostgreSQL Database
**Responsabilidade:** Armazenamento e processamento de dados

**Configuração:**
- **Versão:** PostgreSQL 16
- **Storage:** PVC local (k3d) ou Docker volumes
- **Backup:** Automático via CronJob (diário 02:00)
- **Deploy:** Direct deployment ou Helm Chart

**Schema Overview:**
```sql
-- Staging: dados raw
stg.exportacao (
  ano INT,
  mes INT,
  pais VARCHAR,
  setor VARCHAR,
  ncm VARCHAR,
  valor_usd DECIMAL,
  peso_kg DECIMAL,
  ingest_at TIMESTAMP,
  ingest_batch VARCHAR,
  idempotency_key VARCHAR(128)  -- NEW: Epic 3
)

-- Reports: views materializadas
rpt.mv_resumo_pais (
  pais VARCHAR,
  total_usd DECIMAL,
  total_kg DECIMAL,
  participacao_pct DECIMAL
)

rpt.mv_resumo_setor (
  setor VARCHAR,
  total_usd DECIMAL,
  anos INT[]
)

-- Idempotency tracking (Epic 3)
public.api_idempotency (
  id SERIAL PRIMARY KEY,
  idempotency_key VARCHAR(128) UNIQUE,
  request_hash TEXT,
  response_body TEXT,
  response_status INT,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP
)
```

### 4. Integration Gateway (Go) - Epic 1

**Responsabilidade:** Framework genérico para integrações com APIs externas

**Tecnologias:**
- **Runtime:** Go 1.24.9
- **Framework HTTP:** Gin
- **Auth:** mTLS (ICP-Brasil), OAuth2, API Key
- **Resilience:** gobreaker (Circuit Breaker), backoff (Retry)
- **Transform:** JSONPath engine
- **Deploy:** Kubernetes Deployment ou Docker Compose
- **Port:** 8081

**Capacidades:**
- **90% config-driven** - Integrações via YAML (zero código Go)
- **Multi-auth support** - mTLS com certificados ICP-Brasil A1/A3, OAuth2, API Key
- **Automatic resilience** - Circuit breaker, retry with backoff, rate limiting
- **Data transformation** - JSONPath para mapeamento request/response
- **Schema validation** - JSON Schema automático para connector configs
- **Certificate management** - SimpleCertificateManager para ICP-Brasil

**Endpoints:**
```
GET /health                          # Gateway health check
GET /v1/connectors                   # Listar conectores disponíveis
GET /v1/connectors/{id}              # Detalhes de connector
POST /v1/connectors/{id}/{endpoint}  # Executar endpoint com params
```

**Connector Examples:**
- `receita-federal-cnpj.yaml` - Consulta CNPJ com mTLS
- `viacep.yaml` - Consulta CEP (API pública)

**Estrutura de Diretórios:**
```
services/integration-gateway/
├── cmd/gateway/main.go         # Entry point
├── internal/
│   ├── auth/                   # Auth engine (mTLS, OAuth2, API Key)
│   │   ├── engine.go          # Auth factory
│   │   ├── mtls.go            # mTLS handler
│   │   ├── oauth2.go          # OAuth2 handler
│   │   └── certmanager.go     # Certificate manager
│   ├── framework/              # Core framework
│   │   ├── httpclient.go      # Resilient HTTP client
│   │   └── executor.go        # Request executor
│   ├── registry/               # Connector registry
│   │   ├── registry.go        # Registry manager
│   │   └── loader.go          # YAML loader & validator
│   ├── transform/              # Transform engine
│   │   └── engine.go          # JSONPath transformer
│   └── observability/          # Logging & metrics
└── go.mod
```

### 5. Observability Stack - Epic 2

**Responsabilidade:** Métricas, tracing e monitoramento de produção

#### 5.1 Prometheus (Metrics)

**Métricas Implementadas (11 customizadas):**
```go
// HTTP Metrics
bgc_http_requests_total          // Counter: requests por endpoint/method/status
bgc_http_request_duration_seconds // Histogram: latência (P50/P95/P99)
bgc_http_requests_in_flight       // Gauge: requests em processamento

// Database Metrics
bgc_db_queries_total              // Counter: queries por operação/tabela
bgc_db_query_duration_seconds     // Histogram: latência de queries
bgc_db_connections_open           // Gauge: conexões abertas
bgc_db_connections_in_use         // Gauge: conexões em uso
bgc_db_connections_idle           // Gauge: conexões ociosas

// Application Metrics
bgc_errors_total                  // Counter: erros por tipo/severidade
bgc_idempotency_cache_hits_total  // Counter: cache hits
bgc_idempotency_cache_misses_total // Counter: cache misses
bgc_idempotency_cache_size        // Gauge: tamanho do cache
```

**Deployment:**
- Docker Compose: `prom/prometheus:v2.50.0` na porta 9090
- Kubernetes: Deployment com ServiceAccount + RBAC

#### 5.2 Grafana (Dashboards)

**Dashboards Pré-configurados:**
- **BGC API Overview** - 8 painéis:
  - Request Rate (req/s)
  - Error Rate (%)
  - Request Duration (P50/P95/P99)
  - Requests In Flight
  - DB Query Rate
  - DB Query Duration (P95)
  - DB Connections (Open/In Use/Idle)
  - Idempotency Cache (Hits/Misses/Size)

**Deployment:**
- Docker Compose: `grafana/grafana:10.3.0` na porta 3001
- Kubernetes: Deployment com datasources provisionados

#### 5.3 Jaeger (Distributed Tracing)

**Capacidades:**
- OTLP gRPC receiver (porta 4317)
- OTLP HTTP receiver (porta 4318)
- W3C Trace Context propagation
- Automatic span creation para HTTP handlers
- Database query tracing com atributos

**Deployment:**
- Docker Compose: `jaegertracing/all-in-one:1.54` na porta 16686
- Kubernetes: All-in-one deployment com memory storage (10k traces)

**Instrumentação Automática:**
```go
// Middleware OpenTelemetry
otelgin.Middleware("bgc-api")

// Span de database query
ctx, span := tracing.StartSpan(ctx, "db.GetMarketDataByYearRange")
defer span.End()
span.SetAttributes(
  attribute.Int("year.from", yearFrom),
  attribute.Int("year.to", yearTo),
)
```

---

## 💾 Modelo de Dados

### Schema Staging (`stg`)

#### stg.exportacao
**Propósito:** Dados raw de exportação carregados via ingest

| Coluna | Tipo | Descrição | Exemplo |
|--------|------|-----------|---------|
| `ano` | INT | Ano da exportação | 2023 |
| `mes` | INT | Mês (1-12) | 3 |
| `pais` | VARCHAR(100) | País de destino | "China" |
| `setor` | VARCHAR(100) | Setor econômico | "Agricultura" |
| `ncm` | VARCHAR(20) | Código NCM | "17011100" |
| `valor_usd` | DECIMAL(15,2) | Valor em USD | 1250000.50 |
| `peso_kg` | DECIMAL(15,3) | Peso em kg | 850000.125 |
| `ingest_at` | TIMESTAMP | Quando foi carregado | 2025-09-07 14:30:00 |
| `ingest_batch` | VARCHAR(50) | ID do batch de carga | "batch_20250907_143000" |
| `ingest_source` | VARCHAR(100) | Arquivo fonte | "dados_jan_2023.xlsx" |

**Índices:**
```sql
-- Performance de queries analíticas
CREATE INDEX idx_exportacao_ano_pais ON stg.exportacao(ano, pais);
CREATE INDEX idx_exportacao_setor ON stg.exportacao(setor);
CREATE INDEX idx_exportacao_ingest ON stg.exportacao(ingest_batch);
```

### Schema Reports (`rpt`)

#### rpt.mv_resumo_pais
**Propósito:** Agregação por país (dados base para análises)

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| `pais` | VARCHAR(100) | País (PK) |
| `total_usd` | DECIMAL(18,2) | Soma valor USD |
| `total_kg` | DECIMAL(18,3) | Soma peso kg |
| `participacao_pct` | DECIMAL(5,2) | % do total |
| `anos` | INT[] | Array de anos com dados |
| `updated_at` | TIMESTAMP | Último refresh |

```sql
-- Definição da MV
CREATE MATERIALIZED VIEW rpt.mv_resumo_pais AS
SELECT 
  pais,
  SUM(valor_usd) as total_usd,
  SUM(peso_kg) as total_kg,
  ROUND(SUM(valor_usd) * 100.0 / SUM(SUM(valor_usd)) OVER(), 2) as participacao_pct,
  ARRAY_AGG(DISTINCT ano ORDER BY ano) as anos,
  NOW() as updated_at
FROM stg.exportacao 
GROUP BY pais;

-- Índice UNIQUE para REFRESH CONCURRENTLY
CREATE UNIQUE INDEX ON rpt.mv_resumo_pais (pais);
```

#### rpt.mv_resumo_geral
**Propósito:** Métricas gerais agregadas

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| `total_usd` | DECIMAL(18,2) | Valor total geral |
| `total_kg` | DECIMAL(18,3) | Peso total geral |
| `paises_count` | INT | Número de países |
| `setores_count` | INT | Número de setores |
| `anos` | INT[] | Anos disponíveis |
| `updated_at` | TIMESTAMP | Último refresh |

---

## 🏗️ Clean Architecture - Fluxo de Requisição

### Exemplo: GET /market/size?metric=TAM&year_from=2023&year_to=2024

```
1. HTTP Request
   └─▶ Gin Router (app/server.go)

2. Middleware Pipeline
   ├─▶ CORS Middleware (api/middleware/cors.go)
   ├─▶ Request ID Middleware (api/middleware/metrics.go)
   └─▶ Logging & Metrics Middleware (api/middleware/metrics.go)

3. Handler Layer (api/handlers/market.go)
   ├─ Parse query parameters
   ├─ Validate input
   └─▶ Call MarketService.CalculateMarketSize(req)

4. Service Layer (business/market/service.go)
   ├─ Apply business rules (TAM/SAM/SOM logic)
   ├─▶ Call Repository.GetMarketDataByYearRange(...)
   ├─ Receive data from repository
   ├─ Calculate SOM percentages (base: 1.5%, aggressive: 3%)
   └─▶ Return MarketSizeResponse

5. Repository Layer (repository/postgres/market.go)
   ├─ Build SQL query with filters
   ├─ Execute query on PostgreSQL
   ├─ Scan rows into MarketItem structs
   └─▶ Return []MarketItem

6. Handler Layer (continued)
   ├─ Receive response from service
   └─▶ Return JSON response with status 200

7. Middleware (continued)
   ├─ Log request (structured JSON)
   ├─ Update metrics counters
   └─▶ Send response to client
```

### Dependency Injection (app/server.go)

```go
func NewServer(cfg *config.AppConfig, db *sql.DB) *Server {
    // Load external configs
    weights := config.LoadPartnerWeights(cfg.PartnerWeightsFile)
    tariffs := config.LoadTariffScenarios(cfg.TariffScenariosFile)
    
    // Create repositories (infrastructure)
    marketRepo := postgres.NewMarketRepository(db)
    routeRepo := postgres.NewRouteRepository(db)
    
    // Create services (business logic) - inject repositories
    marketService := market.NewService(marketRepo, cfg)
    routeService := route.NewService(routeRepo, weights, tariffs)
    healthService := health.NewService(cfg, weights, tariffs)
    
    // Create handlers (presentation) - inject services
    marketHandler := handlers.NewMarketHandler(marketService)
    routeHandler := handlers.NewRouteHandler(routeService)
    healthHandler := handlers.NewHealthHandler(healthService)
    
    // Setup router with middleware
    r := gin.Default()
    r.Use(middleware.CORS())
    r.Use(middleware.RequestID())
    r.Use(middleware.MetricsAndLog())
    
    // Register routes
    r.GET("/health", healthHandler.GetHealth)
    r.GET("/market/size", marketHandler.GetMarketSize)
    r.GET("/routes/compare", routeHandler.CompareRoutes)
    
    return &Server{router: r, config: cfg}
}
```

### Benefícios da Arquitetura

✅ **Testabilidade**: Cada camada pode ser testada isoladamente com mocks  
✅ **Manutenibilidade**: Mudanças em uma camada não afetam outras  
✅ **Legibilidade**: Código organizado por domínio de negócio  
✅ **Reusabilidade**: Services podem ser usados por diferentes handlers  
✅ **Escalabilidade**: Fácil adicionar novos domínios sem afetar existentes  

---

## 🔄 Fluxo de Dados

### 1. Ingest Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Arquivo   │───▶│ bgc-ingest  │───▶│stg.exportac │
│  CSV/XLSX   │    │   Job       │    │     ao      │
└─────────────┘    └─────────────┘    └─────────────┘
                          │                   │
                          ▼                   ▼
                   ┌─────────────┐    ┌─────────────┐
                   │   Logs      │    │ Audit Trail │
                   │  kubectl    │    │ingest_batch │
                   └─────────────┘    └─────────────┘
```

### 2. Refresh Flow  
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   CronJob   │───▶│ REFRESH MV  │───▶│ rpt.mv_*    │
│  refresh-mv │    │CONCURRENTLY │    │  updated    │
└─────────────┘    └─────────────┘    └─────────────┘
       │                                     │
       ▼                                     ▼
┌─────────────┐                     ┌─────────────┐
│    Daily    │                     │     API     │
│   01:00     │                     │ Performance │
└─────────────┘                     └─────────────┘
```

### 3. Query Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│   bgc-api   │───▶│ rpt.mv_*    │
│Browser/Post │    │  Handler    │    │ PostgreSQL  │
└─────────────┘    └─────────────┘    └─────────────┘
       ▲                   │                   │
       │                   ▼                   ▼
       └───────────┌─────────────┐    ┌─────────────┐
                   │    JSON     │    │ SQL Query   │
                   │  Response   │    │ Execution   │
                   └─────────────┘    └─────────────┘
```

---

## 🐳 Deployment Architecture

### Kubernetes Resources

#### Deployments
```yaml
# bgc-api deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bgc-api
spec:
  replicas: 1
  selector:
    matchLabels:
      app: bgc-api
  template:
    spec:
      containers:
      - name: api
        image: bgc/bgc-api:dev
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres
        - name: DB_USER
          value: bgc
        - name: DB_PASS
          value: bgc
        - name: DB_NAME
          value: bgc
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8080
```

#### Services
```yaml
# bgc-api service
apiVersion: v1
kind: Service
metadata:
  name: bgc-api
spec:
  selector:
    app: bgc-api
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP

# bgc-web service
apiVersion: v1
kind: Service
metadata:
  name: bgc-web
spec:
  selector:
    app: bgc-web
  ports:
  - port: 3000
    targetPort: 3000
  type: ClusterIP
```

#### CronJobs
```yaml
# Backup PostgreSQL diariamente
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
spec:
  schedule: "0 2 * * *"  # 02:00 daily
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:16
            command: ["/bin/sh", "-c", "pg_dump ..."]

# Refresh Materialized Views
apiVersion: batch/v1
kind: CronJob
metadata:
  name: mview-refresh
spec:
  schedule: "0 3 * * *"  # 03:00 daily
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: refresh
            image: postgres:16
            command: ["/bin/sh", "-c", "psql -c 'REFRESH MATERIALIZED VIEW...'"]
```

### Network Flow (Kubernetes)
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Developer  │───▶│ Traefik     │───▶│   k3d       │
│  Browser    │    │ LoadBalancer│    │  cluster    │
│web.bgc.local│    │   :80       │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                                              │
                          ┌───────────────────┴────────────────┐
                          ▼                                    ▼
                  ┌─────────────┐                    ┌─────────────┐
                  │ bgc-web     │───rewrites────────▶│ bgc-api     │
                  │ service     │                    │ service     │
                  │ :3000       │                    │ :8080       │
                  └─────────────┘                    └─────────────┘
```

---

## 🔒 Segurança

### Desenvolvimento (Ambiente Local)
- **Database:** Credenciais consistentes (bgc/bgc/bgc) em todos ambientes
- **Network:** Cluster interno (K8s) ou localhost (Docker Compose)
- **Images:** Local build, sem registry externo
- **Data:** Dados de exemplo, não sensíveis
- **Frontend:** Next.js com rewrites para proxy API (sem exposição direta)

### Planos Futuros
- [ ] **RBAC:** Kubernetes role-based access control
- [ ] **TLS:** Certificados internos para comunicação
- [ ] **Secrets:** Vault ou External Secrets Operator
- [ ] **Network Policies:** Isolamento entre namespaces
- [ ] **Image Security:** Registry privado + scanning
- [ ] **Audit:** Logs de acesso e modificações

---

## 📊 Performance

### Objetivos de Performance (Sprint 1)
- **API Response Time:** < 500ms (p95)
- **Concurrent Users:** 10 (desenvolvimento)
- **Data Volume:** ~1M registros de exemplo
- **MV Refresh:** < 30 segundos

### Estratégias de Otimização
1. **Materialized Views:** Pre-computação de agregações
2. **Índices:** Cobertura para queries principais
3. **Connection Pooling:** Planejado para Sprint 2
4. **Caching:** Redis planejado para futuro

### Monitoring (Planejado)
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Application │───▶│ Prometheus  │───▶│  Grafana    │
│   Metrics   │    │   TSDB      │    │ Dashboard   │
└─────────────┘    └─────────────┘    └─────────────┘
```

---

## 🔄 CI/CD Pipeline (Futuro)

### Planejamento Sprint 2+
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   GitHub    │───▶│   Actions   │───▶│   k3d       │
│    Push     │    │   Build     │    │   Deploy    │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Tests     │    │   Images    │    │   Health    │
│   Lint      │    │   Registry  │    │   Check     │
└─────────────┘    └─────────────┘    └─────────────┘
```

---

## 🎯 Roadmap Técnico

### Sprint 2 ✅ (Completa - Outubro 2025)
- [x] **Clean Architecture** - Refatoração completa para hexagonal architecture
- [x] **Health Endpoints** - `/health`, `/healthz` com status de configuração
- [x] **OpenAPI Spec** - Documentação formal em `/openapi.yaml` e `/docs`
- [x] **Error Handling** - Padronização de erros HTTP por domínio
- [x] **Logging** - Logs estruturados JSON com request ID
- [x] **Metrics** - Endpoint `/metrics` com contadores e latências
- [x] **Market Analytics** - Endpoint `/market/size` com TAM/SAM/SOM
- [x] **Route Comparison** - Endpoint `/routes/compare` com tariff scenarios
- [x] **Configuration Management** - YAMLs para partners, tariffs, scope, SOM
- [x] **Docker Compose** - Ambiente de desenvolvimento simplificado

### Sprint 3 (Próxima)
- [ ] **Unit Tests** - Testes para services e repositories
- [ ] **Integration Tests** - Testes E2E automatizados
- [ ] **Basic Metrics** - Counters e histogramas

### Sprint 3-4 (Médio Prazo)
- [ ] **Connection Pooling** - pgxpool
- [ ] **Caching Layer** - Redis
- [ ] **Background Jobs** - Async processing
- [ ] **Multi-environment** - dev/staging/prod configs
- [ ] **Integration Tests** - API + Database

### Sprint 5+ (Longo Prazo)
- [ ] **Cloud Migration** - EKS/GKE
- [ ] **High Availability** - Multi-replica + LoadBalancer
- [ ] **Data Pipeline** - Stream processing
- [ ] **ML Integration** - Análises preditivas

---

## 🧪 Testing Strategy

### Níveis de Teste
```
┌─────────────────────────────────────────────────────────────┐
│                    Testing Pyramid                         │
│                                                             │
│                    ┌─────────────┐                         │
│                    │     E2E     │ ← Full system tests     │
│                    └─────────────┘                         │
│                 ┌─────────────────────┐                    │
│                 │    Integration      │ ← API + DB tests   │
│                 └─────────────────────┘                    │
│              ┌─────────────────────────────┐               │
│              │        Unit Tests           │ ← Business    │
│              └─────────────────────────────┘   logic       │
└─────────────────────────────────────────────────────────────┘
```

### Implementação Atual (Sprint 1)
- **Manual Testing:** Postman collections
- **Ad-hoc:** curl commands
- **DB Testing:** Manual SQL queries

### Planos Sprint 2+
```go
// Unit tests - business logic
func TestCalculateParticipacao(t *testing.T) {
    // Test percentage calculation logic
}

// Integration tests - API + DB
func TestMetricsResumoEndpoint(t *testing.T) {
    // Test full request flow
}

// E2E tests - full system
func TestIngestAndQuery(t *testing.T) {
    // Load data + query via API
}
```

---

## 📋 ADRs (Architecture Decision Records)

### ADR-001: k3d como Runtime Local
**Status:** ✅ Aceito  
**Contexto:** Necessidade de ambiente Kubernetes local  
**Decisão:** k3d em vez de minikube/kind  
**Razão:** Mais leve, Docker nativo, fácil setup  
**Consequências:** Limitado a desenvolvimento local

### ADR-002: PostgreSQL via Helm Bitnami
**Status:** ✅ Aceito  
**Contexto:** Banco de dados para desenvolvimento  
**Decisão:** Bitnami PostgreSQL chart  
**Razão:** Produção-ready, bem documentado  
**Consequências:** Dependência de registry externo

### ADR-003: Go como Linguagem Principal
**Status:** ✅ Aceito  
**Contexto:** Backend API e ingest  
**Decisão:** Go 1.23+ para ambos serviços  
**Razão:** Performance, simplicidade, ecosystem  
**Consequências:** Curva de aprendizado para equipe

### ADR-004: Materialized Views para Performance
**Status:** ✅ Aceito  
**Contexto:** Queries analíticas complexas  
**Decisão:** MVs em vez de views normais  
**Razão:** Performance previsível  
**Consequências:** Complexidade de refresh

### ADR-005: Monorepo Structure
**Status:** ✅ Aceito  
**Contexto:** Organização de código  
**Decisão:** Monorepo com services/api + services/ingest  
**Razão:** Simplicidade para equipe pequena  
**Consequências:** Deploy acoplado

---

## 🔧 Troubleshooting Guide

### Problemas Comuns

#### 1. API não responde
```bash
# Verificar pod status
kubectl get pods | grep bgc-api

# Ver logs
kubectl logs deployment/bgc-api

# Verificar service
kubectl get svc bgc-api

# Port-forward manual
kubectl port-forward svc/bgc-api 3000:3000
```

#### 2. Banco de dados inacessível
```bash
# Status do PostgreSQL
kubectl get pods | grep postgres

# Conectar ao banco
kubectl run psql-client --rm -it --image bitnami/postgresql:latest -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql -h bgc-postgres -U postgres

# Verificar dados
SELECT COUNT(*) FROM stg.exportacao;
```

#### 3. MVs não atualizadas
```sql
-- Verificar última atualização
SELECT updated_at FROM rpt.mv_resumo_pais LIMIT 1;

-- Refresh manual
REFRESH MATERIALIZED VIEW CONCURRENTLY rpt.mv_resumo_pais;
```

#### 4. Imagens não encontradas
```bash
# Listar imagens no k3d
k3d image list -c bgc

# Re-importar
docker build -t bgc/api:dev services/api/
k3d image import bgc/api:dev -c bgc
kubectl rollout restart deployment/bgc-api
```

### Health Checks
```bash
# Script de verificação rápida
#!/bin/bash
echo "🔍 BGC Health Check"

# 1. Cluster
kubectl get nodes | grep Ready && echo "✅ Cluster OK" || echo "❌ Cluster FAIL"

# 2. Pods
kubectl get pods | grep Running | wc -l | xargs echo "✅ Running pods:"

# 3. API
curl -s http://localhost:3000/metrics/resumo > /dev/null && echo "✅ API OK" || echo "❌ API FAIL"

# 4. Database
kubectl exec deployment/bgc-postgres -- psql -U postgres -c "SELECT 1" && echo "✅ DB OK" || echo "❌ DB FAIL"
```

---

## 📚 References & Standards

### Coding Standards
- **Go:** [Effective Go](https://golang.org/doc/effective_go.html)
- **SQL:** Snake_case, explicit naming
- **K8s:** [Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
- **Git:** Conventional Commits

### API Standards
- **REST:** Richardson Maturity Model Level 2
- **JSON:** camelCase para responses
- **HTTP:** Status codes padronizados
- **Versioning:** URL path (/v1/, /v2/)

### Database Standards
```sql
-- Naming conventions
schemas: stg (staging), rpt (reports), cfg (config)
tables: snake_case, plural
columns: snake_case, descriptive
indexes: idx_<table>_<columns>
mvs: mv_<purpose>_<grain>

-- Data types
timestamps: TIMESTAMP WITH TIME ZONE
money: DECIMAL(precision, scale)
text: VARCHAR with explicit limits
arrays: type[] for multiple values
```

### Container Standards
```dockerfile
# Multi-stage builds
FROM golang:1.23 AS builder
FROM gcr.io/distroless/base-debian11 AS runtime

# Non-root user
USER 65534

# Health checks
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:3000/health || exit 1
```

---

## 📊 Metrics & Observability

### Application Metrics (Planejado)
```go
// Prometheus metrics
var (
    httpRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "bgc_http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"endpoint", "method", "status"},
    )
    
    queryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "bgc_query_duration_seconds",
            Help: "Database query duration",
        },
        []string{"query_type"},
    )
)
```

### Infrastructure Metrics
- **Kubernetes:** CPU, Memory, Network via metrics-server
- **PostgreSQL:** Connections, queries/sec, cache hit ratio
- **Storage:** Disk usage, I/O patterns

### Logging Strategy
```json
{
  "timestamp": "2025-09-07T14:30:00Z",
  "level": "info",
  "service": "bgc-api",
  "endpoint": "/metrics/resumo",
  "method": "GET",
  "duration_ms": 234,
  "status": 200,
  "user_agent": "PostmanRuntime/7.32.3",
  "query_params": {"ano": "2023"},
  "trace_id": "abc123"
}
```

---

## 🚀 Deployment Environments

### Development (Current)
- **Runtime:** k3d local cluster
- **Database:** PostgreSQL via Helm
- **Storage:** Local Docker volumes
- **Networking:** Port-forward para acesso
- **Data:** Samples sintéticos
- **Monitoring:** Logs básicos

### Staging (Futuro)
- **Runtime:** Cloud Kubernetes (EKS/GKE)
- **Database:** Managed PostgreSQL
- **Storage:** Cloud persistent volumes
- **Networking:** Internal load balancer
- **Data:** Subset produção anonimizada
- **Monitoring:** Prometheus + Grafana

### Production (Futuro)
- **Runtime:** Multi-AZ Kubernetes
- **Database:** HA PostgreSQL cluster
- **Storage:** Replicated storage
- **Networking:** Public load balancer + CDN
- **Data:** Dados reais de exportação
- **Monitoring:** Full observability stack

---

## 🎯 Success Metrics

### Technical KPIs
- **Uptime:** 99.9% (objetivo futuro)
- **Response Time:** p95 < 500ms
- **Error Rate:** < 0.1%
- **Data Freshness:** MVs updated < 1h lag
- **Build Time:** < 5 minutes
- **Deploy Time:** < 2 minutes

### Business KPIs
- **Query Performance:** Complex analytics < 1s
- **Data Accuracy:** 100% consistency
- **User Experience:** Self-service analytics
- **Development Velocity:** Features/sprint
- **Operational Overhead:** Minimal manual intervention

---

**Documento mantido por:** [Time de Arquitetura]  
**Próxima revisão:** Sprint 2 Planning  
**Feedback:** [Link para issues do GitHub]
