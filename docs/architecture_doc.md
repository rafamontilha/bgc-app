# Arquitetura do Sistema BGC Analytics

**Versão:** 1.0  
**Última atualização:** Setembro 2025  
**Status:** Sprint 1 Implementada

## 📋 Visão Geral

O BGC Analytics é um sistema de analytics para dados de exportação brasileira, construído com arquitetura cloud-native para execução em ambiente Kubernetes local (k3d) durante desenvolvimento.

### Objetivos do Sistema
- **Performance:** Consultas analíticas rápidas via Materialized Views
- **Simplicidade:** Stack mínima e bem documentada 
- **Desenvolvimento ágil:** Ambiente local reproducível
- **Escalabilidade:** Preparado para migração cloud futura

---

## 🏗️ Arquitetura de Alto Nível

```
┌─────────────────────────────────────────────────────────────────┐
│                           k3d Cluster                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   bgc-api       │  │  bgc-postgres   │  │   bgc-ingest    │ │
│  │                 │  │                 │  │                 │ │
│  │ ┌─────────────┐ │  │ ┌─────────────┐ │  │ ┌─────────────┐ │ │
│  │ │ GET /metrics│ │  │ │ PostgreSQL  │ │  │ │ CSV/XLSX    │ │ │
│  │ │    /resumo  │ │  │ │   Database  │ │  │ │   Loader    │ │ │
│  │ │    /pais    │ │  │ │             │ │  │ │             │ │ │
│  │ └─────────────┘ │  │ │  ┌───────┐  │ │  │ └─────────────┘ │ │
│  │                 │  │ │  │  MVs  │  │ │  │                 │ │
│  │ Port: 3000      │  │ │  └───────┘  │ │  │ CronJob         │ │
│  └─────────────────┘  │ └─────────────┘ │  └─────────────────┘ │
│           │            │        ▲        │           │         │
│           │            │        │        │           │         │
│           └────────────┼────────┘        │           │         │
│                        │                 │           │         │
│                        └─────────────────┼───────────┘         │
└─────────────────────────────────────────────────────────────────┘
                         │
                    ┌────▼────┐
                    │ kubectl │
                    │port-fwd │
                    │   :3000 │
                    └─────────┘
                         │
                    ┌────▼────┐
                    │   Web   │
                    │ Browser │
                    │Postman  │
                    └─────────┘
```

---

## 🔧 Componentes Principais

### 1. BGC API (Go)
**Responsabilidade:** API REST para consultas analíticas

**Tecnologias:**
- **Runtime:** Go 1.23+
- **Framework:** HTTP nativo + gorilla/mux (planejado)
- **Database:** lib/pq (PostgreSQL driver)
- **Deploy:** Kubernetes Deployment

**Endpoints Atuais:**
```
GET /metrics/resumo[?ano=YYYY&setor=Nome]
GET /metrics/pais[?ano=YYYY&limit=N]
```

**Arquitetura Interna:**
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Handler   │───▶│  Service    │───▶│ Repository  │
│             │    │   Layer     │    │             │
│ HTTP Routes │    │ Business    │    │ SQL Queries │
│ JSON/Error  │    │ Logic       │    │ Connection  │
└─────────────┘    └─────────────┘    └─────────────┘
```

### 2. BGC Ingest (Go)
**Responsabilidade:** ETL de dados CSV/XLSX para PostgreSQL

**Tecnologias:**
- **Runtime:** Go 1.23+
- **CSV:** encoding/csv nativo
- **Excel:** github.com/xuri/excelize/v2
- **Database:** lib/pq
- **Deploy:** Kubernetes CronJob

**Comandos Disponíveis:**
```bash
bgc-ingest load-csv /path/to/file.csv
bgc-ingest load-xlsx /path/to/file.xlsx [--sheet=Nome]
bgc-ingest refresh-mv
```

**Processo de Ingest:**
```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│CSV/XLSX │───▶│ Parser  │───▶│Transform│───▶│  Load   │
│ Files   │    │         │    │ Validate│    │ INSERT  │
└─────────┘    └─────────┘    └─────────┘    └─────────┘
                    │              │              │
                    ▼              ▼              ▼
               ┌─────────┐    ┌─────────┐    ┌─────────┐
               │Row-by-  │    │Business │    │Batch    │
               │Row Read │    │Rules    │    │Upsert   │
               └─────────┘    └─────────┘    └─────────┘
```

### 3. PostgreSQL (Bitnami Helm)
**Responsabilidade:** Armazenamento e processamento de dados

**Configuração:**
- **Versão:** PostgreSQL 15+ (via Bitnami)
- **Storage:** PVC local (k3d)
- **Backup:** Manual (desenvolvimento)
- **Deploy:** Helm Chart

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
  ingest_batch VARCHAR
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
**Propósito:** Agregação por país para endpoint `/metrics/pais`

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
**Propósito:** Métricas gerais para endpoint `/metrics/resumo`

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| `total_usd` | DECIMAL(18,2) | Valor total geral |
| `total_kg` | DECIMAL(18,3) | Peso total geral |
| `paises_count` | INT | Número de países |
| `setores_count` | INT | Número de setores |
| `anos` | INT[] | Anos disponíveis |
| `updated_at` | TIMESTAMP | Último refresh |

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
        image: bgc/api:dev
        ports:
        - containerPort: 3000
        env:
        - name: DB_HOST
          value: bgc-postgres
        livenessProbe:
          httpGet:
            path: /health  # TODO: implementar
            port: 3000
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
  - port: 3000
    targetPort: 3000
  type: ClusterIP
```

#### CronJobs
```yaml
# Refresh MVs diariamente
apiVersion: batch/v1
kind: CronJob
metadata:
  name: refresh-mv
spec:
  schedule: "0 1 * * *"  # 01:00 daily
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: ingest
            image: bgc/ingest:dev
            command: ["refresh-mv"]
```

### Network Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Developer  │───▶│k3d-loadbala │───▶│   k3d       │
│  localhost  │    │   ncer      │    │  cluster    │
│    :3000    │    │    :8080    │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                                              │
                                              ▼
                                    ┌─────────────┐
                                    │ bgc-api     │
                                    │ service     │
                                    │ :3000       │
                                    └─────────────┘
```

---

## 🔒 Segurança

### Desenvolvimento (Sprint 1)
- **Database:** Credenciais via Kubernetes Secret (Helm generated)
- **Network:** Cluster interno, acesso via port-forward
- **Images:** Local build, sem registry externo
- **Data:** Dados de exemplo, não sensíveis

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

### Sprint 2 (Próxima)
- [ ] **Health Endpoints** - `/health`, `/ready`
- [ ] **OpenAPI Spec** - Documentação formal
- [ ] **Error Handling** - Padronização de erros HTTP
- [ ] **Logging** - Logs estruturados JSON
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