# BGC App - Sistema de Analytics de Exportação

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.24.9+-00ADD8?logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?logo=postgresql)](https://www.postgresql.org)

Plataforma completa de analytics para dados de exportação brasileira com:
- **API REST** em Go 1.24.9 (Clean Architecture, Gin framework)
- **Frontend** em Next.js 15 (React 19, TypeScript, Tailwind CSS)
- **Banco de Dados** PostgreSQL 16 com Materialized Views
- **Integration Gateway** para APIs externas (mTLS, OAuth2, Circuit Breaker)
- **Observability Stack** completa (Prometheus, Grafana, Jaeger, OpenTelemetry)
- **API Contracts** com JSON Schemas e Idempotency

**Open Source** sob licença AGPL v3 - Garantindo que melhorias permaneçam livres e acessíveis à comunidade.

## 🚀 Quick Start

### Kubernetes com k3d (Recomendado)

```powershell
# 1. Setup inicial (primeira vez)
.\scripts\k8s.ps1 setup

# 2. Configurar hosts (executar como Administrador)
.\scripts\setup-hosts.ps1

# 3. Acessar aplicação
# Web UI:  http://web.bgc.local
# Routes:  http://web.bgc.local/routes
# API:     http://api.bgc.local/healthz
```

### Docker Compose (Desenvolvimento Local)

```powershell
# Iniciar ambiente
.\scripts\docker.ps1 up

# URLs disponíveis:
# Web:        http://localhost:3000
# API:        http://localhost:8080
# Prometheus: http://localhost:9090
# Grafana:    http://localhost:3001 (admin / admin)
# Jaeger UI:  http://localhost:16686
# PgAdmin:    http://localhost:5050 (admin@bgc.dev / admin)
```

---

## 📋 Pré-requisitos

### Para Docker Compose
- Docker Desktop
- PowerShell

### Para Kubernetes
- Docker Desktop
- k3d
- kubectl
- PowerShell

---

## 🎯 Scripts de Gerenciamento

### Makefile (Multiplataforma)

```bash
make help                # Mostrar todos os comandos disponíveis
make docker-up           # Iniciar Docker Compose
make k8s-setup           # Setup inicial Kubernetes
make k8s-status          # Status do cluster
make seed                # Carregar dados de exemplo
make restore-backup      # Restaurar backup do PostgreSQL
```

### Docker Compose

```powershell
.\scripts\docker.ps1 up          # Iniciar serviços
.\scripts\docker.ps1 down        # Parar serviços
.\scripts\docker.ps1 restart     # Reiniciar
.\scripts\docker.ps1 logs        # Ver logs
.\scripts\docker.ps1 ps          # Status dos containers
.\scripts\docker.ps1 build       # Rebuildar imagens
.\scripts\docker.ps1 clean       # Limpar tudo (remove volumes)
.\scripts\docker.ps1 help        # Ajuda
```

### Kubernetes

```powershell
.\scripts\k8s.ps1 setup          # Setup inicial (cluster + deploy)
.\scripts\k8s.ps1 up             # Deploy serviços
.\scripts\k8s.ps1 down           # Remover deployments
.\scripts\k8s.ps1 restart        # Reiniciar pods
.\scripts\k8s.ps1 logs           # Ver logs
.\scripts\k8s.ps1 status         # Status do cluster (inclui HPA e CronJobs)
.\scripts\k8s.ps1 build          # Rebuildar imagens
.\scripts\k8s.ps1 open           # Abrir no browser
.\scripts\k8s.ps1 clean          # Deletar cluster
.\scripts\k8s.ps1 help           # Ajuda
```

### Gerenciamento de Dados

```powershell
# Carregar dados de exemplo
.\scripts\seed.ps1

# Restaurar backup (Kubernetes)
.\scripts\restore-backup.ps1                    # Listar backups
.\scripts\restore-backup.ps1 -BackupFile <nome> # Restaurar específico
```

---

## 📁 Estrutura do Projeto

```
bgc-app/
├── api/                         # API Go (Clean Architecture)
│   ├── cmd/api/                # Entry point
│   ├── config/                 # Configurações YAML
│   ├── internal/               # Código interno
│   │   ├── business/          # Lógica de negócio (domain)
│   │   ├── repository/        # Persistência (postgres)
│   │   ├── api/               # Handlers HTTP, middleware, validation
│   │   ├── observability/     # Metrics (Prometheus) & Tracing (OTel)
│   │   └── app/               # Wiring & server
│   ├── Dockerfile
│   └── go.mod
│
├── web-next/                    # Frontend Next.js 15 (React + TypeScript)
│   ├── app/                    # App Router do Next.js
│   │   ├── v1/                # API routes v1 (proxies)
│   │   └── healthz/           # Health check route
│   ├── components/             # Componentes React
│   ├── lib/                    # Utilitários e API client
│   ├── Dockerfile
│   └── package.json
│
├── services/                    # Microserviços
│   ├── bgc-ingest/             # Serviço de ingestão (CSV/XLSX)
│   └── integration-gateway/    # Gateway de integrações externas
│       ├── cmd/gateway/        # Entry point
│       ├── internal/
│       │   ├── auth/           # Multi-auth (mTLS, OAuth2, API Key)
│       │   ├── framework/      # HTTP client resiliente
│       │   ├── registry/       # Connector registry
│       │   ├── transform/      # Transform engine (JSONPath)
│       │   └── observability/  # Logging & metrics
│       └── go.mod
│
├── config/                      # Configurações externas
│   └── connectors/             # YAML configs (Receita Federal, ViaCEP)
│
├── schemas/                     # JSON Schemas de validação
│   ├── connector.schema.json   # Schema para connectors
│   └── v1/                     # API v1 request/response schemas
│
├── certs/                       # Certificados ICP-Brasil (gitignored)
│
├── db/                          # Database
│   ├── init/                   # Schema inicial (Docker Compose)
│   └── migrations/             # Migrations SQL (inc. idempotency)
│
├── k8s/                         # Kubernetes Manifests
│   ├── api.yaml                # Deployment API com HPA
│   ├── web.yaml                # Deployment Web com HPA
│   ├── integration-gateway/    # Gateway deployment & configs
│   │   ├── deployment.yaml     # Deployment, Service, HPA
│   │   ├── configmap.yaml      # Connector configs
│   │   ├── sealed-secret-*.yaml # Sealed Secrets (Bitnami)
│   │   └── README-SECRETS.md   # Guia de secrets management
│   ├── network-policies/       # Network segmentation & isolation
│   │   ├── bgc-api-netpol.yaml # Policy da API (força uso do Gateway)
│   │   ├── integration-gateway-netpol.yaml
│   │   └── README.md           # Guia completo com testes
│   ├── observability/          # Prometheus, Grafana, Jaeger
│   ├── postgres-backup-cronjob.yaml
│   └── mview-refresh-cronjob.yaml
│
├── bgcstack/                    # Docker Compose stack
│   ├── docker-compose.yml      # Serviços principais
│   └── observability/          # Configs Prometheus, Grafana
│
├── tests/                       # Testes de integração
│
├── docs/                        # Documentação técnica
│   ├── OBSERVABILITY.md        # Guia completo de observabilidade
│   ├── CONNECTOR-GUIDE.md      # Guia de integrações externas
│   ├── DATA-DICTIONARY.md      # Dicionário de dados
│   ├── IDEMPOTENCY-POLICY.md   # Política de idempotência
│   ├── API-SIMULATOR.md        # Documentação da API do Simulador
│   ├── PRODUCT-ROADMAP.md      # Roadmap estratégico de produto
│   ├── PRODUCT-DECISIONS.md    # Registro de decisões de produto
│   ├── PRODUCT-METRICS.md      # Métricas e KPIs de produto
│   ├── NEXT-STEPS.md           # Próximos passos priorizados
│   └── EPIC-*.md               # Documentação dos épicos
│
├── Makefile                     # Wrapper multiplataforma
└── CHANGELOG.md                 # Histórico de mudanças
```

---

## 🏗️ Arquitetura

### Clean Architecture (Hexagonal)

A API segue os princípios de Clean Architecture com separação clara de responsabilidades:

- **Domain (business/)**: Lógica de negócio pura
  - Entities
  - Repository interfaces
  - Services

- **Infrastructure (repository/)**: Implementações de persistência
  - PostgreSQL repositories

- **Presentation (api/)**: Camada HTTP
  - Handlers
  - Middleware
  - Routing

- **Application (app/)**: Wiring e configuração
  - Dependency injection
  - Server initialization

### Stack Tecnológica

**Backend & Services:**
- **API**: Go 1.24.9 com Gin (Clean Architecture)
- **Integration Gateway**: Go 1.24.9 (Hybrid connector framework)
- **Frontend**: Next.js 15 (React 19, TypeScript, Tailwind CSS)

**Observability:**
- **Metrics**: Prometheus + Grafana
- **Tracing**: Jaeger + OpenTelemetry
- **Logging**: Structured JSON logs

**Data & Storage:**
- **Database**: PostgreSQL 16
- **Caching**: In-memory (idempotency)
- **Schemas**: JSON Schema validation

**Infrastructure:**
- **Container**: Docker + Docker Compose
- **Orchestration**: Kubernetes (k3d)
- **Ingress**: Traefik

**Security & Integration:**
- **Auth**: mTLS (ICP-Brasil), OAuth2, API Key
- **Secrets**: Kubernetes Secrets API + Sealed Secrets (Bitnami)
- **Network**: Network Policies (Zero Trust, Least Privilege)
- **Resilience**: Circuit Breaker, Retry, Rate Limiting

### Arquitetura de Segurança de Rede

O projeto implementa **Defense in Depth** com múltiplas camadas de segurança:

#### Segmentação de Rede (Network Policies)

```
Internet
   ↓ (HTTPS/TLS 1.3)
Ingress Controller
   ↓
bgc-api (namespace: data)
   ├─→ PostgreSQL ✅
   ├─→ Redis ✅
   ├─→ Integration Gateway ✅ (ÚNICO caminho para APIs externas)
   └─→ ❌ BLOQUEADO: Internet direta (porta 443)

Integration Gateway (namespace: data)
   ├─→ Kubernetes API ✅ (buscar Secrets)
   ├─→ Redis ✅ (cache L2)
   ├─→ PostgreSQL ✅ (cache L3)
   └─→ APIs Externas ✅ (ComexStat, ViaCEP, Receita Federal)
```

**Princípios Aplicados:**
- **Zero Trust**: Todo tráfego negado por padrão (`default-deny-all`)
- **Least Privilege**: Cada pod acessa APENAS o necessário
- **Network Isolation**: Isolamento completo entre serviços
- **Forced Gateway Pattern**: API principal OBRIGADA a usar Integration Gateway

**Arquivos:**
- `k8s/network-policies/` - Todas as policies + guia de testes
- Ver: `k8s/network-policies/README.md` para troubleshooting

#### Secrets Management

**Fluxo Seguro:**
```
Developer → Script (create-sealed-secret) → Sealed Secret (criptografado)
   → Git (safe to commit) → Kubernetes → Secret (runtime, in-memory)
```

**Componentes:**
- **KubernetesSecretStore** (`services/integration-gateway/internal/auth/k8s_secret_store.go`)
  - Busca secrets via Kubernetes API (`k8s.io/client-go`)
  - Cache in-memory com TTL de 5 minutos
  - Thread-safe, backward compatible com env vars

- **Sealed Secrets Controller** (Bitnami)
  - Secrets criptografados no Git (public-key cryptography)
  - Descriptografia automática no cluster
  - Rotação de secrets suportada

**Formato:**
```yaml
# Connector YAML
auth:
  type: api_key
  api_key:
    key_ref: comexstat-credentials/api-key  # secret-name/key-name
```

**Script de Criação:**
```bash
# Criar sealed secret de forma interativa
./scripts/create-sealed-secret-comexstat.sh
```

**Ver:** `k8s/integration-gateway/README-SECRETS.md` para guia completo

---

## 📊 Endpoints da API

### Base URLs

- **Docker Compose**: `http://localhost:8080`
- **Kubernetes**: `http://web.bgc.local` (via Ingress)

### Principais Endpoints

**Core API (v1):**
```
GET /healthz                          # Health check
GET /v1/market/size                   # TAM/SAM/SOM metrics (com JSON schema)
GET /v1/routes/compare                # Comparação de rotas (com JSON schema)
GET /docs                             # API documentation (Redoc)
GET /openapi.yaml                     # OpenAPI spec
```

**Observability:**
```
GET /metrics                          # Prometheus metrics (formato nativo)
GET /metrics/json                     # Metrics em JSON (legacy)
```

**Integration Gateway:**
```
GET /health                           # Gateway health check
GET /v1/connectors                    # Listar todos os connectors
GET /v1/connectors/{id}               # Detalhes de um connector
POST /v1/connectors/{id}/{endpoint}   # Executar endpoint com params
```

**Nota**: Endpoints legacy (`/market/size`, `/routes/compare`) redirecionam automaticamente para `/v1/*` (301).

### Exemplos

```bash
# Health check
curl http://localhost:8080/healthz

# Market size (TAM) - v1 endpoint
curl "http://localhost:8080/v1/market/size?metric=TAM&year_from=2023&year_to=2024"

# Routes compare - v1 endpoint
curl "http://localhost:8080/v1/routes/compare?from=USA&alts=CHN,ARE&ncm_chapter=84&year=2024"

# Prometheus metrics
curl http://localhost:8080/metrics

# Integration Gateway - Consultar CEP
curl -X POST "http://localhost:8081/v1/connectors/viacep/consultar" \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310-100"}'
```

---

## 🔧 Desenvolvimento

### Modificar Código da API

```powershell
# Docker Compose
# O código é montado via volume, basta reiniciar:
.\scripts\docker.ps1 restart

# Kubernetes
# Precisa rebuildar a imagem:
.\scripts\k8s.ps1 build
```

### Modificar Frontend

```powershell
# Docker Compose
# Arquivos montados via volume, basta recarregar o browser

# Kubernetes
# Precisa rebuildar a imagem:
.\scripts\k8s.ps1 build
```

### Adicionar Dependências Go

```powershell
cd api
go get <package>
go mod tidy

# Depois rebuildar
cd ..
.\scripts\docker.ps1 build  # ou .\scripts\k8s.ps1 build
```

---

## 🔒 Configuração de Segurança

**⚠️ IMPORTANTE:** Este projeto usa gestão segura de credenciais.

### Primeiro Uso (Docker Compose)

```powershell
# 1. Copiar template de configuração
cd bgcstack
cp .env.example .env

# 2. Gerar senhas fortes
openssl rand -base64 32  # PostgreSQL
openssl rand -base64 32  # PgAdmin

# 3. Editar .env com as senhas geradas
notepad .env

# 4. Iniciar stack
.\scripts\docker.ps1 up
```

### Kubernetes

As credenciais no Kubernetes são gerenciadas via **Sealed Secrets** + **KubernetesSecretStore**:

```powershell
# 1. Sealed Secrets controller (já instalado)
kubectl get pods -n kube-system | grep sealed-secrets

# 2. Criar sealed secret para ComexStat API
.\scripts\create-sealed-secret-comexstat.sh

# 3. Aplicar no cluster
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml

# 4. Verificar secret descriptografado
kubectl get secret comexstat-credentials -n data
```

### Network Policies

Verificar isolamento de rede:

```powershell
# Listar policies aplicadas
kubectl get networkpolicies -n data

# Testar conectividade (deve falhar - API → Internet bloqueada)
kubectl exec -it deployment/bgc-api -n data -- curl -I https://google.com
# Esperado: timeout

# Testar Gateway (deve funcionar)
kubectl exec -it deployment/integration-gateway -n data -- curl -I https://google.com
# Esperado: 200 OK
```

### Documentação Completa

📖 **Guias Detalhados:**
- **Secrets**: [k8s/integration-gateway/README-SECRETS.md](k8s/integration-gateway/README-SECRETS.md)
- **Network Policies**: [k8s/network-policies/README.md](k8s/network-policies/README.md)

**🚨 NUNCA commite credenciais em plain text no Git!**

---

## 🧪 Testes

### Smoke Tests (Docker Compose)

```powershell
# API Health
curl http://localhost:8080/healthz

# Web UI
curl http://localhost:3000

# Database
docker exec bgc_db psql -U bgc -d bgc -c "SELECT version();"
```

### Smoke Tests (Kubernetes)

```powershell
# Status geral
.\scripts\k8s.ps1 status

# API via Ingress
curl http://web.bgc.local/healthz

# Web UI
Start-Process http://web.bgc.local
```

---

## 🐛 Troubleshooting

### Docker Compose

**Porta já em uso**:
```powershell
# Verificar o que está usando a porta
netstat -ano | findstr :8080

# Parar serviços
.\scripts\docker.ps1 down
```

**Container não inicia**:
```powershell
# Ver logs
.\scripts\docker.ps1 logs

# Limpar e reiniciar
.\scripts\docker.ps1 clean
.\scripts\docker.ps1 up
```

### Kubernetes

**web.bgc.local não resolve**:
```powershell
# Executar como Administrador:
.\scripts\configure-hosts.ps1
```

**Pods em CrashLoopBackOff**:
```powershell
# Ver logs
.\scripts\k8s.ps1 logs

# Verificar status
.\scripts\k8s.ps1 status

# Recriar cluster
.\scripts\k8s.ps1 clean
.\scripts\k8s.ps1 setup
```

**Imagens não encontradas**:
```powershell
# Rebuildar e reimportar
.\scripts\k8s.ps1 build
```

---

## 📚 Documentação Adicional

### Código Fonte

- API: Veja comentários no código em `api/internal/`
- Frontend: Arquivos HTML são auto-documentados

### Arquitetura

- Clean Architecture aplicada em `api/internal/`
- Repository pattern em `api/internal/repository/`
- Dependency injection em `api/internal/app/server.go`

### Deployment

- Docker Compose: `bgcstack/docker-compose.yml`
- Kubernetes: manifests em `deploy/` e `k8s/`

---

## 📊 Observabilidade e Resiliência

### Métricas e Monitoramento

**Prometheus Metrics:**
- 11 métricas customizadas implementadas:
  - HTTP: `bgc_http_requests_total`, `bgc_http_request_duration_seconds`, `bgc_http_requests_in_flight`
  - DB: `bgc_db_queries_total`, `bgc_db_query_duration_seconds`, `bgc_db_connections_*`
  - Errors: `bgc_errors_total`
  - Idempotency: `bgc_idempotency_cache_*`

**Dashboards:**
```bash
# Prometheus
http://localhost:9090           # Docker Compose
http://prometheus.bgc.local     # Kubernetes

# Grafana (pré-configurado com dashboards)
http://localhost:3001           # Docker Compose
http://grafana.bgc.local        # Kubernetes
```

### Distributed Tracing

**Jaeger + OpenTelemetry:**
- Tracing automático de todas as requisições HTTP
- Spans de database queries com atributos detalhados
- W3C Trace Context propagation
- OTLP gRPC exporter

```bash
# Jaeger UI
http://localhost:16686          # Docker Compose
http://jaeger.bgc.local         # Kubernetes
```

### Health Probes (Kubernetes)

Todos os serviços possuem health checks configurados:

**API e WEB:**
- **Readiness Probe**: Verifica se o pod está pronto para receber tráfego
- **Liveness Probe**: Detecta e reinicia pods travados automaticamente

```bash
# Verificar status dos probes
kubectl describe pod -n data -l app=bgc-api | grep -A 5 Probes
```

### Horizontal Pod Autoscaling (HPA)

Escala automática baseada em CPU e memória:

**API:**
- Min: 1 pod, Max: 5 pods
- Target: 70% CPU, 80% Memory

**WEB:**
- Min: 1 pod, Max: 3 pods
- Target: 70% CPU, 80% Memory

```bash
# Visualizar status do HPA
kubectl get hpa -n data

# Ver métricas em tempo real
kubectl top pods -n data
```

### Backups Automatizados

**CronJob de Backup PostgreSQL:**
- Executa diariamente às 02:00
- Mantém os últimos 7 backups
- Backups comprimidos (.sql.gz)
- Armazenados em PVC persistente

```bash
# Listar backups disponíveis
.\scripts\restore-backup.ps1

# Trigger backup manual
kubectl create job --from=cronjob/postgres-backup manual-backup -n data

# Restaurar backup
.\scripts\restore-backup.ps1 -BackupFile bgc_backup_YYYYMMDD_HHMMSS.sql.gz
```

### Materialized Views Refresh

**CronJob de Refresh:**
- Executa diariamente às 03:00
- Atualiza todas as materialized views
- Usa refresh concorrente (sem lock)

```bash
# Ver status dos CronJobs
kubectl get cronjobs -n data

# Ver histórico de execuções
kubectl get jobs -n data
```

---

## 🤝 Contribuindo

1. Faça fork do projeto
2. Crie uma branch: `git checkout -b feature/nova-feature`
3. Commit: `git commit -m 'feat: adiciona nova feature'`
4. Push: `git push origin feature/nova-feature`
5. Abra um Pull Request

---

## 📝 Licença

Este projeto está licenciado sob a **GNU Affero General Public License v3.0 (AGPL-3.0)**.

### O que isso significa?

- ✅ **Liberdade de usar** - Você pode usar este software para qualquer propósito
- ✅ **Liberdade de estudar** - Você pode examinar como o software funciona
- ✅ **Liberdade de modificar** - Você pode modificar o software
- ✅ **Liberdade de distribuir** - Você pode distribuir cópias do software
- ⚠️ **Copyleft de rede** - Se você executar uma versão modificada em um servidor e permitir que outros usuários interajam com ela pela rede, você **deve** disponibilizar o código-fonte modificado

### Por que AGPL?

Escolhemos a AGPL v3 para garantir que:
- Melhorias ao software permaneçam livres e abertas
- Empresas que usam o software como SaaS contribuam com melhorias de volta à comunidade
- O ecossistema de dados de exportação brasileira se beneficie do conhecimento compartilhado

### Arquivos de Licença

- `LICENSE` - Texto completo da licença AGPL v3
- `NOTICE` - Avisos de copyright e informações de componentes

Para mais informações, consulte: https://www.gnu.org/licenses/agpl-3.0.html

---

## ✨ Features

### Core API
- ✅ API REST com Clean Architecture (Go 1.24.9)
- ✅ Dashboard interativo TAM/SAM/SOM (Next.js 15)
- ✅ Comparação de rotas de exportação
- ✅ PostgreSQL 16 com Materialized Views
- ✅ Migrations automáticas com rastreabilidade
- ✅ **API Versioning** - Endpoints /v1/* com backward compatibility
- ✅ **JSON Schema Validation** - Validação automática de request/response
- ✅ **Idempotency** - Cache 24h com Idempotency-Key header

### Integration & External APIs (Epic 1)
- ✅ **Integration Gateway** - Framework híbrido (90% config, 10% plugins)
- ✅ **Multi-auth Support** - mTLS (ICP-Brasil A1/A3), OAuth2, API Key
- ✅ **Resilience Patterns** - Circuit Breaker, Retry with backoff, Rate limiting
- ✅ **Transform Engine** - JSONPath para mapeamento de dados
- ✅ **Certificate Manager** - Gestão de certificados ICP-Brasil
- ✅ **Connector Registry** - Validação automática via JSON Schema

### Observability Stack (Epic 2)
- ✅ **Prometheus Metrics** - 11 métricas customizadas
- ✅ **Grafana Dashboards** - Pré-configurados (API Overview, DB Stats)
- ✅ **Distributed Tracing** - Jaeger + OpenTelemetry (OTLP)
- ✅ **Automatic Instrumentation** - Middleware para HTTP e DB
- ✅ **Alert Rules** - 10 regras pré-configuradas (Prometheus)
- ✅ **Structured Logging** - JSON logs com trace context

### Data Governance (Epic 3)
- ✅ **JSON Schemas** - Validação de contratos de API
- ✅ **Idempotency System** - Prevenção de processamento duplicado
- ✅ **Data Dictionary** - Documentação completa do modelo de dados
- ✅ **API Versioning** - Suporte a múltiplas versões futuras

### Deployment
- ✅ Docker Compose para desenvolvimento
- ✅ Kubernetes (k3d) para produção simulada
- ✅ Scripts PowerShell unificados
- ✅ Makefile multiplataforma
- ✅ Traefik Ingress com TLS

### Resilience & Automation
- ✅ Health probes (readiness/liveness)
- ✅ Horizontal Pod Autoscaling (HPA)
- ✅ Resource limits e requests
- ✅ Backups automáticos diários do PostgreSQL
- ✅ Refresh automático de materialized views
- ✅ CronJobs para automação

### Security & Network Isolation
- ✅ **Go 1.24.9** - Correção de 5 vulnerabilidades críticas (CVE-2024-*)
- ✅ **Network Policies** - Zero Trust, Least Privilege, Defense in Depth
  - Default deny-all no namespace `data`
  - API bloqueada de acessar internet diretamente
  - Forced Gateway Pattern para integrações externas
- ✅ **Kubernetes Secrets API** - KubernetesSecretStore com cache (5min TTL)
- ✅ **Sealed Secrets (Bitnami)** - Criptografia de secrets para Git
  - Public-key cryptography
  - Scripts automatizados (`create-sealed-secret-*.sh`)
  - Rotação de secrets suportada
- ✅ **ConfigMaps e Secrets management** - Separação de configuração e credenciais
- ✅ **Non-root containers** - Princípio de menor privilégio no runtime
- ✅ **AGPL v3 license** - Garantia de código aberto

---

**Desenvolvido com ❤️ pela equipe BGC**
