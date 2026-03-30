# Brasil Global Connect (BGC)

[![Version](https://img.shields.io/badge/version-0.5.0-green)](CHANGELOG.md)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.24.9+-00ADD8?logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?logo=postgresql)](https://www.postgresql.org)
[![Tests](https://img.shields.io/badge/tests-71%20passing-brightgreen)](#testes)

Plataforma open source de inteligência para exportação, voltada para PMEs brasileiras. Conecta dados oficiais (Comex Stat, Siscomex) com modelos de análise de mercado, recomendação de destinos e gestão operacional — em um fluxo integrado, da pesquisa à execução.

**Open Source** sob licença AGPL v3 — garantindo que melhorias permaneçam livres e acessíveis à comunidade.

---

## O que está rodando hoje (v0.5.0)

- **Autenticação completa** — Clerk: email/senha, Google OAuth, recuperação de senha, SSO callback
- **Onboarding personalizado** — wizard 3 passos (produto NCM → volume mensal → regiões-alvo), persiste perfil do exportador
- **Dashboard adaptado ao perfil** — preview do simulador com NCM do usuário, tutorial interativo no primeiro acesso, banner de reengajamento
- **Export Destination Simulator** — 50 países, scoring ponderado, estimativas financeiras, freemium 5 req/dia, cache multicamada
- **Observabilidade completa** — Prometheus, Grafana, Jaeger (OpenTelemetry), structured logging
- **Integration Gateway** — framework híbrido para APIs externas (mTLS ICP-Brasil, OAuth2, API Key, Circuit Breaker)
- **Infraestrutura production-ready** — Docker Compose + Kubernetes (k3d), HPA, backups automáticos, Network Policies Zero Trust

---

## Quick Start

### Docker Compose (Desenvolvimento Local)

```bash
# 1. Copiar configuração
cd bgcstack
cp .env.example .env
# Editar .env com suas senhas (POSTGRES_PASSWORD, PGADMIN_DEFAULT_PASSWORD, DB_PASS)

# 2. Iniciar stack
.\scripts\docker.ps1 up
```

| Serviço      | URL                                      |
|--------------|------------------------------------------|
| Web (App)    | http://localhost:3000                    |
| API          | http://localhost:8080                    |
| Grafana      | http://localhost:3001 (admin / admin)    |
| Jaeger UI    | http://localhost:16686                   |
| Prometheus   | http://localhost:9090                    |
| PgAdmin      | http://localhost:5050 (admin@bgc.dev)    |

### Kubernetes com k3d

```powershell
# Setup inicial (primeira vez)
.\scripts\k8s.ps1 setup

# Configurar hosts (executar como Administrador)
.\scripts\setup-hosts.ps1

# Acessar
# Web:  http://web.bgc.local
# API:  http://api.bgc.local/healthz
```

---

## Stack Tecnológica

| Camada             | Tecnologia                                                        |
|--------------------|-------------------------------------------------------------------|
| **Frontend**       | Next.js 15, React 19, TypeScript, MUI v7 (Material Design 3)     |
| **Autenticação**   | Clerk (JWT, JWKS validation, OAuth2, Smart CAPTCHA)              |
| **Backend API**    | Go 1.24.9, Gin (Clean Architecture)                              |
| **Integration GW** | Go 1.24.9 — hybrid connector framework (mTLS, OAuth2, API Key)   |
| **Banco de dados** | PostgreSQL 16 (Materialized Views, Migrations)                   |
| **Cache**          | L1 Ristretto (in-process) + L2 Redis + L3 PostgreSQL             |
| **Observability**  | Prometheus, Grafana, Jaeger, OpenTelemetry (OTLP)                |
| **Infra**          | Docker Compose, Kubernetes (k3d), Traefik Ingress                |
| **Segurança**      | Network Policies (Zero Trust), Sealed Secrets (Bitnami), mTLS   |

---

## Arquitetura

### Clean Architecture (API Go)

```
api/internal/
├── business/     # Lógica de negócio pura (domain)
├── repository/   # Persistência PostgreSQL
├── api/          # HTTP handlers, middleware, routing
└── app/          # Dependency injection, server init
```

### Segmentação de Rede (Kubernetes)

```
Internet
   ↓ HTTPS/TLS 1.3
Ingress Controller (Traefik)
   ↓
bgc-api              → PostgreSQL ✅ | Redis ✅ | Integration Gateway ✅
                     → Internet direta ❌ BLOQUEADO (Network Policy)

Integration Gateway  → APIs Externas ✅ (ComexStat, ViaCEP, Receita Federal)
                     → Redis ✅ (cache L2) | PostgreSQL ✅ (cache L3)
```

**Princípios:** Zero Trust (default-deny-all), Least Privilege, Forced Gateway Pattern para integrações externas.

Detalhes: [`k8s/network-policies/README.md`](k8s/network-policies/README.md)

---

## Endpoints da API

### Base URLs
- **Docker Compose**: `http://localhost:8080`
- **Kubernetes**: `http://api.bgc.local`

### Principais Endpoints

```
# Core
GET  /healthz                           # Health check
GET  /v1/market/size                    # TAM/SAM/SOM metrics (auth required)
GET  /docs                              # API docs (Redoc)
GET  /openapi.yaml                      # OpenAPI spec

# Simulator
POST /v1/simulator/destinations         # Recomendar destinos de exportação

# Observability
GET  /metrics                           # Prometheus metrics
GET  /metrics/json                      # Metrics JSON (legacy)

# Integration Gateway (porta 8081)
GET  /health
GET  /v1/connectors
POST /v1/connectors/{id}/{endpoint}
```

### Exemplo

```bash
# Simular destinos para açúcar (NCM 17011400)
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -H "Content-Type: application/json" \
  -d '{"ncm": "17011400", "volume_kg": 1000, "max_results": 10}'

# Health check
curl http://localhost:8080/healthz
```

---

## Estrutura do Projeto

```
bgc-app/
├── api/                         # API Go (Clean Architecture)
│   ├── cmd/api/                # Entry point
│   ├── config/                 # Configurações YAML
│   └── internal/
│       ├── business/           # Lógica de negócio
│       ├── repository/         # Persistência PostgreSQL
│       ├── api/                # Handlers, middleware (auth JWT/Clerk), routing
│       ├── observability/      # Prometheus + OpenTelemetry
│       └── app/                # Wiring + server init
│
├── web-next/                    # Frontend Next.js 15 + MUI v7
│   ├── app/                    # App Router
│   │   ├── login/              # Página de login custom (Clerk)
│   │   ├── signup/             # Página de cadastro custom (Clerk)
│   │   ├── onboarding/         # Wizard 3 passos (NCM, volume, regiões)
│   │   ├── dashboard/          # Dashboard personalizado por perfil
│   │   ├── simulator/          # Export Destination Simulator
│   │   ├── profile/            # Perfil + re-onboarding
│   │   └── api/onboarding/     # API route (persiste Clerk publicMetadata)
│   ├── components/
│   │   ├── onboarding/         # OnboardingStep1-3, Welcome, Tutorial, Banner
│   │   ├── dashboard/          # DashboardSimulatorPreview
│   │   ├── simulator/          # SimulatorForm, DestinationList, etc.
│   │   └── auth/               # SSO callback components
│   └── lib/
│       ├── data/               # ncm-chapters.ts (96), trade-regions.ts (35 países)
│       └── types/              # onboarding.ts, simulator.ts
│
├── services/
│   └── integration-gateway/    # Gateway de integrações externas
│       └── internal/
│           ├── auth/           # mTLS, OAuth2, API Key, KubernetesSecretStore
│           ├── framework/      # HTTP client resiliente (Circuit Breaker, Retry)
│           ├── registry/       # Connector registry + JSON Schema validation
│           └── transform/      # Transform engine (JSONPath)
│
├── config/connectors/           # YAML configs (Receita Federal, ViaCEP, ComexStat)
├── db/                          # Schema inicial + Migrations SQL
├── k8s/                         # Kubernetes manifests (deployments, HPA, CronJobs)
├── bgcstack/                    # Docker Compose stack + configs Prometheus/Grafana
├── schemas/                     # JSON Schemas de validação de API
├── docs/                        # Documentação técnica
└── CHANGELOG.md
```

---

## Testes

```bash
# Rodar todos os testes (frontend)
cd web-next
pnpm test

# Com cobertura
pnpm test -- --coverage
```

**71 testes automatizados** em 8 suítes (Jest + React Testing Library):

| Suíte                       | Testes |
|-----------------------------|--------|
| Dados NCM (96 capítulos)    | 13     |
| Regiões de comércio         | 14     |
| API route onboarding        | 6      |
| ReminderBanner              | 6      |
| OnboardingTutorial          | 7      |
| Re-onboarding (profile)     | 4      |
| Simulator form              | 11     |
| Outros                      | 10     |

---

## Configuração de Segurança

**⚠️ IMPORTANTE:** Nunca commite credenciais em plain text no Git.

### Docker Compose

```bash
cd bgcstack
cp .env.example .env
# Gerar senhas fortes:
openssl rand -base64 32   # usar para POSTGRES_PASSWORD, DB_PASS, PGADMIN_DEFAULT_PASSWORD
```

### Clerk (Autenticação)

Adicione no `.env`:
```
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_...
CLERK_SECRET_KEY=sk_...
CLERK_JWKS_URL=https://your-clerk-domain.clerk.accounts.dev/.well-known/jwks.json
```

### Kubernetes — Sealed Secrets

```bash
# Criar sealed secret para credentials de APIs externas
./scripts/create-sealed-secret-comexstat.sh

# Aplicar no cluster
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml
```

Guia completo: [`k8s/integration-gateway/README-SECRETS.md`](k8s/integration-gateway/README-SECRETS.md)

---

## Scripts de Gerenciamento

### Docker Compose

```powershell
.\scripts\docker.ps1 up          # Iniciar serviços
.\scripts\docker.ps1 down        # Parar
.\scripts\docker.ps1 logs        # Ver logs
.\scripts\docker.ps1 build       # Rebuildar imagens
.\scripts\docker.ps1 clean       # Limpar tudo (remove volumes)
```

### Kubernetes

```powershell
.\scripts\k8s.ps1 setup          # Setup inicial
.\scripts\k8s.ps1 up             # Deploy
.\scripts\k8s.ps1 status         # Status + HPA + CronJobs
.\scripts\k8s.ps1 build          # Rebuildar imagens
.\scripts\k8s.ps1 clean          # Deletar cluster
```

### Makefile (multiplataforma)

```bash
make help           # Todos os comandos
make docker-up      # Iniciar Docker Compose
make k8s-setup      # Setup Kubernetes
make seed           # Carregar dados de exemplo
```

---

## Observabilidade

### Dashboards

| Ferramenta  | Docker Compose             | Kubernetes                 |
|-------------|----------------------------|----------------------------|
| Grafana     | http://localhost:3001       | http://grafana.bgc.local   |
| Prometheus  | http://localhost:9090       | http://prometheus.bgc.local|
| Jaeger UI   | http://localhost:16686      | http://jaeger.bgc.local    |

### Métricas customizadas (API)

- HTTP: `bgc_http_requests_total`, `bgc_http_request_duration_seconds`, `bgc_http_requests_in_flight`
- DB: `bgc_db_queries_total`, `bgc_db_query_duration_seconds`, `bgc_db_connections_*`
- Sistema: `bgc_errors_total`, `bgc_idempotency_cache_*`

### Automação (Kubernetes CronJobs)

- **Backup PostgreSQL** — diário às 02:00, últimos 7 backups comprimidos
- **Refresh Materialized Views** — diário às 03:00, refresh concorrente (sem lock)

```bash
# Listar backups e restaurar
.\scripts\restore-backup.ps1
.\scripts\restore-backup.ps1 -BackupFile bgc_backup_YYYYMMDD_HHMMSS.sql.gz

# HPA em tempo real
kubectl get hpa -n data
kubectl top pods -n data
```

---

## Troubleshooting

### Porta em uso

```powershell
netstat -ano | findstr :8080
.\scripts\docker.ps1 down
```

### Container não inicia

```powershell
.\scripts\docker.ps1 logs
.\scripts\docker.ps1 clean && .\scripts\docker.ps1 up
```

### Kubernetes — pods em CrashLoopBackOff

```powershell
.\scripts\k8s.ps1 logs
.\scripts\k8s.ps1 clean && .\scripts\k8s.ps1 setup
```

### web.bgc.local não resolve

```powershell
# Executar como Administrador:
.\scripts\configure-hosts.ps1
```

---

## Documentação

| Documento                                          | Conteúdo                                      |
|----------------------------------------------------|-----------------------------------------------|
| [docs/API-SIMULATOR.md](docs/API-SIMULATOR.md)     | API do Simulador de Destinos                  |
| [docs/PRODUCT-ROADMAP.md](docs/PRODUCT-ROADMAP.md) | Roadmap estratégico de produto                |
| [docs/PRODUCT-DECISIONS.md](docs/PRODUCT-DECISIONS.md) | Registro de decisões de produto          |
| [docs/PRODUCT-METRICS.md](docs/PRODUCT-METRICS.md) | Métricas e KPIs                               |
| [docs/OBSERVABILITY.md](docs/OBSERVABILITY.md)     | Guia completo de observabilidade              |
| [docs/CONNECTOR-GUIDE.md](docs/CONNECTOR-GUIDE.md) | Integrações externas (Integration Gateway)    |
| [docs/DATA-DICTIONARY.md](docs/DATA-DICTIONARY.md) | Dicionário de dados                           |
| [docs/RUNBOOK.md](docs/RUNBOOK.md)                 | Operações e incidentes                        |
| [CHANGELOG.md](CHANGELOG.md)                       | Histórico de versões                          |

---

## Contribuindo

1. Faça fork do projeto
2. Crie uma branch: `git checkout -b feature/nova-feature`
3. Commit: `git commit -m 'feat: adiciona nova feature'`
4. Push: `git push origin feature/nova-feature`
5. Abra um Pull Request

---

## Licença

**GNU Affero General Public License v3.0 (AGPL-3.0)**

Você pode usar, estudar, modificar e distribuir este software. Se executar uma versão modificada em servidor com acesso público via rede, deve disponibilizar o código-fonte modificado.

Escolhemos a AGPL v3 para garantir que melhorias ao software permaneçam livres — especialmente para empresas que usam o software como SaaS.

Mais informações: https://www.gnu.org/licenses/agpl-3.0.html

---

**Desenvolvido com pela equipe BGC**
