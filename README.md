# Projeto BGC - Analytics de Exportação

Sistema de analytics para dados de exportação brasileira com stack Kubernetes local (k3d) + PostgreSQL + API Go.

## 🚀 Quick Start

```bash
# 1. Criar cluster k3d
k3d cluster create bgc --port "8080:80@loadbalancer"

# 2. Instalar PostgreSQL via Helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install bgc-postgres bitnami/postgresql --namespace default

# 3. Executar migrations
kubectl apply -f deploy/migrations/

# 4. Fazer build e import das imagens
docker build -t bgc/ingest:dev services/ingest/
docker build -t bgc/api:dev services/api/
k3d image import bgc/ingest:dev bgc/api:dev -c bgc

# 5. Deploy dos serviços
kubectl apply -f deploy/api/
kubectl apply -f deploy/ingest/

# 6. Port-forward para acessar API
kubectl port-forward service/bgc-api 3000:3000
```

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
cd services/ingest
docker build -t bgc/ingest:dev .
k3d image import bgc/ingest:dev -c bgc

# API service  
cd services/api
docker build -t bgc/api:dev .
k3d image import bgc/api:dev -c bgc

# Re-deploy
kubectl rollout restart deployment/bgc-api
kubectl rollout restart cronjob/bgc-ingest
```

### Estrutura do Projeto
```
bgc/
├── deploy/           # Manifests Kubernetes
│   ├── api/         # Deployment/Service da API
│   ├── ingest/      # CronJob de ingest
│   └── migrations/  # Jobs de migração SQL
├── services/
│   ├── api/         # Código da API Go
│   └── ingest/      # Código do ingest Go
├── db/              # Scripts SQL
├── docs/            # Documentação técnica
└── scripts/         # Scripts auxiliares
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

- [Post-mortem Sprint 1](docs/sprint1-postmortem.md)
- [Arquitetura do Sistema](docs/architecture.md)
- [Guia de Deployment](docs/deployment.md)

## 🤝 Contribuindo

1. Faça fork do projeto
2. Crie uma branch: `git checkout -b feature/nova-feature`
3. Commit: `git commit -m 'Adiciona nova feature'`
4. Push: `git push origin feature/nova-feature`
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob licença MIT. Veja [LICENSE](LICENSE) para mais detalhes.

---

**Status do Projeto**: 🟢 Sprint 1 Completa - Ambiente local funcional com API read-only

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