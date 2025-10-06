# Projeto BGC - Analytics de Exportação

Sistema de analytics para dados de exportação brasileira com stack Kubernetes local (k3d) + PostgreSQL + API Go.

## 🚀 Quick Start

### Opção A: Docker Compose (Recomendado para Desenvolvimento)

```powershell
# 1. Iniciar o stack completo
cd bgcstack
docker compose up -d

# 2. Carregar dados de exemplo (opcional)
cd ..
pwsh scripts/seed.ps1

# 3. Acessar
# API: http://localhost:8080
# Web UI: http://localhost:3000
# PgAdmin: http://localhost:5050
```

### Opção B: Kubernetes (k3d)

```powershell
# 0. Limpar ambiente (se necessário)
docker compose -f bgcstack/docker-compose.yml down
k3d cluster delete bgc 2>$null

# Liberar porta 8080 (se ocupada)
$processId = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
if ($processId) {
    Stop-Process -Id $processId -Force
    Write-Host "Porta 8080 liberada"
}

# 1. Criar cluster k3d
k3d cluster create bgc

# 2. Criar namespace
kubectl create namespace data

# 3. Deploy PostgreSQL (imagem oficial postgres:16)
kubectl apply -f deploy/postgres.yaml

# 4. Criar ConfigMaps (SQL para migrations)
kubectl apply -f deploy/configmap-migrate-0001.yaml
kubectl apply -f deploy/configmap-migrate-0002.yaml
kubectl apply -f deploy/configmap-views.yaml
kubectl apply -f deploy/configmap-mviews-init.yaml
kubectl apply -f deploy/configmap-mviews-add-uniq.yaml
kubectl apply -f deploy/configmap-mviews-populate.yaml
kubectl apply -f deploy/configmap-mviews-refresh.yaml
kubectl apply -f deploy/configmap-sample-csv.yaml

# 5. Aguardar PostgreSQL
kubectl wait --for=condition=ready pod -l app=postgres -n data --timeout=120s

# 6. Build imagens (NOMES CORRETOS)
docker build -t bgc/bgc-api:dev api/
docker build -t bgc/bgc-ingest:dev services/bgc-ingest/
k3d image import bgc/bgc-api:dev bgc/bgc-ingest:dev -c bgc

# 7. Aplicar migrations (ordem correta)
kubectl apply -f deploy/bgc-migrate-0001.yaml
kubectl apply -f deploy/bgc-migrate-0002.yaml
kubectl apply -f deploy/bgc-create-views.yaml

# 8. Aguardar migrations completarem
kubectl wait --for=condition=complete job/bgc-migrate-0001 -n data --timeout=120s
kubectl wait --for=condition=complete job/bgc-migrate-0002 -n data --timeout=120s

# 9. Deploy API
kubectl apply -f deploy/bgc-api.yaml

# 10. Aguardar API
kubectl wait --for=condition=ready pod -l app=bgc-api -n data --timeout=120s

# 11. Verificar status
kubectl get pods -n data

# 12. Port-forward (mantém em execução)
kubectl port-forward service/bgc-api 8080:8080 -n data
```

**Testar (em outro terminal PowerShell):**
```powershell
# Health check
curl http://localhost:8080/health

# Market size
curl "http://localhost:8080/market/size?metric=TAM&year_from=2023&year_to=2024"
```

## 🏗️ Arquitetura Clean (Hexagonal)

O projeto segue os princípios de Clean Architecture e Hexagonal Architecture para máxima manutenibilidade e testabilidade.

### Estrutura da API (api/)

```
api/
├── cmd/api/main.go              # Entry point da aplicação
├── internal/
│   ├── config/                  # Configuração e carregamento de YAML
│   │   └── config.go
│   ├── business/                # Camada de domínio (lógica de negócio)
│   │   ├── market/             # Domínio de métricas de mercado
│   │   │   ├── entities.go    # Estruturas de dados
│   │   │   ├── repository.go  # Interface do repository
│   │   │   └── service.go     # Lógica de TAM/SAM/SOM
│   │   ├── route/              # Domínio de comparação de rotas
│   │   │   ├── entities.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── health/             # Domínio de health check
│   │       └── service.go
│   ├── repository/              # Camada de persistência
│   │   └── postgres/           # Implementação PostgreSQL
│   │       ├── db.go           # Conexão com DB
│   │       ├── market.go       # Queries de mercado
│   │       └── route.go        # Queries de rotas
│   ├── api/                     # Camada de apresentação (HTTP)
│   │   ├── handlers/           # HTTP handlers
│   │   │   ├── health.go
│   │   │   ├── market.go
│   │   │   └── route.go
│   │   └── middleware/         # Middlewares HTTP
│   │       ├── cors.go
│   │       └── metrics.go
│   └── app/                     # Wiring e dependency injection
│       └── server.go           # Inicialização do servidor
├── config/                      # Arquivos de configuração
│   ├── partners_stub.yaml
│   ├── tariff_scenarios.yaml
│   ├── scope.yaml
│   └── som.yaml
└── openapi.yaml                 # Especificação OpenAPI

```

### Princípios Aplicados

✅ **Separação de Responsabilidades**: Cada camada tem uma responsabilidade clara
✅ **Dependency Inversion**: Camadas externas dependem de interfaces internas
✅ **Testabilidade**: Services e repositories são facilmente testáveis via mocks
✅ **Independência de Framework**: Lógica de negócio isolada do Gin
✅ **Manutenibilidade**: Código modular e organizado por domínio

## 📋 Checklist Pós-Reboot

Após reiniciar o computador, siga estes passos para restaurar o ambiente:

### 1. Verificar se o cluster ainda existe
```bash
k3d cluster list
# Se não existir: k3d cluster create bgc --port "8080:80@loadbalancer"
```

### 2. Corrigir kubeconfig (CRÍTICO)
```bash
# Verificar se kubectl funciona
kubectl get nodes

# Se der erro de conexão, encontrar a nova porta do serverlb:
docker ps | findstr k3d-bgc-serverlb
# Anote a porta mapeada (ex: 0.0.0.0:54321->6443/tcp)

# Atualizar kubeconfig com a nova porta
kubectl config set-cluster k3d-bgc --server "https://127.0.0.1:54321"
# Substitua 54321 pela porta encontrada

# Testar novamente
kubectl get nodes
```

### 3. Verificar PostgreSQL
```bash
kubectl get pods | findstr postgres
# Se não estiver Running: kubectl rollout restart deployment/bgc-postgres
```

### 4. Re-importar imagens locais
```bash
k3d image import bgc/ingest:dev bgc/api:dev -c bgc
```

### 5. Verificar serviços
```bash
kubectl get pods
kubectl get services

# Se API não estiver rodando:
kubectl rollout restart deployment/bgc-api

# Port-forward para testar
kubectl port-forward service/bgc-api 3000:3000
```

### 6. Testar conectividade
```bash
# Em outro terminal
curl http://localhost:3000/metrics/resumo
```

## 🛠️ Comandos Úteis

### Acesso ao PostgreSQL
```bash
# Via cliente Bitnami
kubectl run psql-client --rm -it --image bitnami/postgresql:latest -- /opt/bitnami/scripts/postgresql/entrypoint.sh /opt/bitnami/postgresql/bin/psql -h bgc-postgres -U postgres

# Senha padrão do postgres (buscar no secret)
kubectl get secret bgc-postgres -o jsonpath="{.data.postgres-password}" | base64 -d
```

### Ingest de Dados
```bash
# CSV
kubectl create job load-csv-$(date +%s) --from=cronjob/bgc-ingest -- load-csv /data/sample.csv

# Excel  
kubectl create job load-xlsx-$(date +%s) --from=cronjob/bgc-ingest -- load-xlsx /data/sample.xlsx

# Verificar logs
kubectl logs job/load-xlsx-<timestamp>
```

### Refresh de Materialized Views
```bash
# Via Job one-time
kubectl create job refresh-mv-$(date +%s) --from=cronjob/bgc-ingest -- refresh-mv

# Ou direto no postgres
kubectl exec -it deployment/bgc-postgres -- /opt/bitnami/scripts/postgresql/entrypoint.sh /opt/bitnami/postgresql/bin/psql -U postgres -c "REFRESH MATERIALIZED VIEW CONCURRENTLY rpt.mv_resumo_pais;"
```

## 🔧 Desenvolvimento

### Build Local
```bash
# Ingest service
cd services/bgc-ingest
docker build -t bgc/ingest:dev .
k3d image import bgc/ingest:dev -c bgc

# API service  
cd api
docker build -t bgc/api:dev .
k3d image import bgc/api:dev -c bgc

# Re-deploy
kubectl rollout restart deployment/bgc-api
kubectl rollout restart cronjob/bgc-ingest
```

### Estrutura do Projeto
```
bgc-app/
├── api/                         # API Go com Clean Architecture
│   ├── cmd/api/                # Entry point
│   ├── internal/               # Código interno (não exportável)
│   │   ├── business/          # Domínios (market, route, health)
│   │   ├── repository/        # Implementações de persistência
│   │   ├── api/               # Handlers e middleware HTTP
│   │   ├── app/               # Wiring e servidor
│   │   └── config/            # Configuração
│   ├── config/                 # YAMLs de configuração
│   ├── Dockerfile
│   └── go.mod
├── services/
│   └── bgc-ingest/            # Serviço de ingest de dados
├── db/                         # Scripts SQL e migrations
│   ├── init/                  # Schema inicial
│   └── migrations/            # Migrations incrementais
├── web/                        # Frontend (HTML/JS)
│   ├── index.html             # Dashboard TAM/SAM/SOM
│   └── routes.html            # Comparação de rotas
├── deploy/                     # Manifests Kubernetes
│   ├── api/                   # Deployment da API
│   ├── ingest/                # CronJobs de ingest
│   └── migrations/            # Jobs de migration
├── bgcstack/                   # Docker Compose
│   └── docker-compose.yml
├── docs/                       # Documentação técnica
└── scripts/                    # Scripts auxiliares
```

## 📊 API Endpoints

Base URL (local): `http://localhost:3000`

### Métricas de Resumo
```
GET /metrics/resumo
GET /metrics/resumo?ano=2023
GET /metrics/resumo?setor=Agricultura
```

Resposta:
```json
{
  "total_usd": 280500000000,
  "total_toneladas": 450000000,
  "paises_count": 185,
  "setores_count": 12,
  "anos": [2020, 2021, 2022, 2023]
}
```

### Métricas por País
```
GET /metrics/pais
GET /metrics/pais?ano=2023
GET /metrics/pais?limit=10
```

Resposta:
```json
[
  {
    "pais": "China",
    "total_usd": 89500000000,
    "total_toneladas": 120000000,
    "participacao_pct": 31.9
  }
]
```

## 🐛 Troubleshooting

### Problemas Comuns

| Sintoma | Causa Provável | Solução |
|---------|----------------|---------|
| `dial tcp ... connectex` | kubeconfig desatualizado pós-reboot | Seguir checklist item 2 |
| `ImagePullBackOff` | Imagem local não importada | `k3d image import <image> -c bgc` |
| `psql: command not found` | Executando sem entrypoint Bitnami | Usar comando completo com entrypoint |
| API 404 | Port-forward não ativo | `kubectl port-forward service/bgc-api 3000:3000` |
| MVs vazias | Não foi feito refresh inicial | `REFRESH MATERIALIZED VIEW` (sem CONCURRENTLY) |

### Logs Importantes
```bash
# API logs
kubectl logs deployment/bgc-api

# Postgres logs  
kubectl logs deployment/bgc-postgres

# Últimos jobs de ingest
kubectl get jobs --sort-by=.metadata.creationTimestamp
kubectl logs job/<job-name>
```

### Reset Completo
```bash
# ⚠️  CUIDADO: Apaga tudo
k3d cluster delete bgc
# Depois seguir Quick Start novamente
```

## 📚 Documentação Técnica

- [Arquitetura do Sistema](docs/architecture_doc.md) - Visão completa da arquitetura
- [Guia de Deployment](docs/deployment_guide.md) - Como fazer deploy
- [Post-mortem Sprint 1](docs/sprint1_postmortem.md) - Lições aprendidas

### Desenvolvimento

**Build local:**
```bash
cd api
go mod tidy
go build ./cmd/api
./bgc-api  # Requer PostgreSQL rodando
```

**Testes (quando disponíveis):**
```bash
go test ./internal/...
```

**Adicionar nova funcionalidade:**
1. Criar entidades em `internal/business/{domain}/entities.go`
2. Definir interface repository em `internal/business/{domain}/repository.go`
3. Implementar service com lógica de negócio em `internal/business/{domain}/service.go`
4. Implementar repository em `internal/repository/postgres/{domain}.go`
5. Criar handler em `internal/api/handlers/{domain}.go`
6. Registrar rota em `internal/app/server.go`

## 🤝 Contribuindo

1. Faça fork do projeto
2. Crie uma branch: `git checkout -b feature/nova-feature`
3. Commit: `git commit -m 'Adiciona nova feature'`
4. Push: `git push origin feature/nova-feature`
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob licença MIT. Veja [LICENSE](LICENSE) para mais detalhes.

---

**Status do Projeto**: 🟢 Sprint 2 Completa - Clean Architecture implementada, API funcional com TAM/SAM/SOM e comparação de rotas

## Testes E2E (Sprint 2)
1) API
   - /healthz → ok, `partner_weights:true`, `tariffs_loaded:true`
   - /market/size?metric=TAM&year_from=2024&year_to=2025 → 200 com items
   - /routes/compare?from=USA&alts=CHN,ARE,IND&ncm_chapter=84&year=2024 → 200

2) Dashboard (index.html)
   - TAM 2020–2025 → preenche tabela/KPIs
   - SOM base vs aggressive → aggressive ≈ 2× base (1,5%→3%)
   - Exportar CSV gera arquivo com `ano,metric,valor_usd`

3) Rotas (routes.html)
   - Cenário = `base` e `tarifa10` → `adjusted_total_usd` muda
   - Checagem de soma = OK (soma dos parceiros == total ajustado)
   - Exportar CSV: `partner,share,factor,estimated_usd`

4) Qualidade
   - scripts/checks.ps1 → “Aprovado” ou “Aprovado com alertas”
### Sprint 2 — Testes E2E & Aceite (S2-14)

Siga o roteiro em `docs/Sprint2_E2E_Checklist.md` ou use os comandos abaixo (PowerShell, na raiz `bgc`):

```powershell
docker compose -f .\bgcstack\docker-compose.yml up -d
.\scripts\seed.ps1
curl.exe http://localhost:8080/healthz
start http://localhost:3000
start http://localhost:3000/routes.html
.\scripts\checks.ps1
```

Critérios de aceite:
- `/healthz` ok com `partner_weights:true` e `tariffs_loaded:true`
- `/market/size` retorna TAM/SAM/SOM; **SOM base vs aggressive** difere (~2x)
- `/routes/compare` aplica `tariff_scenario` e soma `results` = `adjusted_total_usd`
- Dashboard e Rotas funcionais com export CSV
- Checks de qualidade **Aprovado** (ou **Aprovado com alertas** apenas para outliers)


**Última atualização**: Setembro 2025
