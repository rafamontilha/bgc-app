# Guia de Deployment - BGC Analytics

**Versão:** 2.0  
**Última atualização:** Outubro 2025  
**Ambiente:** Desenvolvimento Local (Docker Compose) e Kubernetes (k3d)

## 📋 Visão Geral

Este guia cobre dois métodos de deployment:

1. **Docker Compose** (Recomendado para desenvolvimento) - Mais rápido e simples
2. **Kubernetes (k3d)** - Para testar deployment production-like

## 🚀 Opção A: Docker Compose (Recomendado)

### Pré-requisitos
```bash
# Verificar versões
docker --version          # >= 20.10
docker compose version    # >= 2.0
```

### Deployment Completo

#### 1. Clonar e Preparar
```bash
# Clonar repositório
git clone https://github.com/rafamontilha/bgc-app.git
cd bgc-app
```

#### 2. Iniciar Stack
```bash
# Subir todos os serviços
cd bgcstack
docker compose up -d

# Verificar status
docker compose ps

# Ver logs
docker compose logs -f api
```

**Serviços disponíveis:**
- API: http://localhost:8080
- Web UI: http://localhost:3000
- PostgreSQL: localhost:5432
- PgAdmin: http://localhost:5050

#### 3. Carregar Dados (Opcional)
```bash
# Voltar para raiz do projeto
cd ..

# Executar script de seed
pwsh scripts/seed.ps1
```

#### 4. Testar API
```bash
# Health check
curl http://localhost:8080/health

# Market size
curl "http://localhost:8080/market/size?metric=TAM&year_from=2023&year_to=2024"

# Routes comparison
curl "http://localhost:8080/routes/compare?from=USA&alts=CHN&ncm_chapter=84&year=2024"
```

#### 5. Acessar Web Dashboard
Abrir navegador em:
- Dashboard TAM/SAM/SOM: http://localhost:3000
- Comparação de Rotas: http://localhost:3000/routes.html
- Documentação API: http://localhost:8080/docs

### Comandos Úteis - Docker Compose

```bash
# Parar serviços
docker compose down

# Rebuild após mudanças no código
docker compose build api
docker compose up -d api

# Ver logs
docker compose logs -f api

# Acessar banco de dados
docker compose exec db psql -U bgc -d bgc

# Limpar tudo (cuidado: apaga dados)
docker compose down -v
```

---

## 🚀 Opção B: Kubernetes (k3d)

### Pré-requisitos
```bash
# Verificar versões
docker --version          # >= 20.10
kubectl version --client  # >= 1.25
k3d version               # >= 5.4
helm version              # >= 3.10
go version                # >= 1.23 (para builds locais)
```

### Recursos do Sistema
- **RAM:** Mínimo 8GB (recomendado 16GB)
- **CPU:** 4 cores (desenvolvimento básico)
- **Disk:** 20GB livres para imagens e dados
- **Network:** Internet para pull de imagens Helm

### Validação do Ambiente
```powershell
# Windows PowerShell - validação rápida
# Execute este bloco para verificar tudo

Write-Host "🔍 Verificando ambiente BGC..." -ForegroundColor Green

# Docker
try {
    docker version | Out-Null
    Write-Host "✅ Docker OK" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker não encontrado" -ForegroundColor Red
}

# kubectl
try {
    kubectl version --client | Out-Null
    Write-Host "✅ kubectl OK" -ForegroundColor Green  
} catch {
    Write-Host "❌ kubectl não encontrado" -ForegroundColor Red
}

# k3d
try {
    k3d version | Out-Null
    Write-Host "✅ k3d OK" -ForegroundColor Green
} catch {
    Write-Host "❌ k3d não encontrado" -ForegroundColor Red
}

# Helm
try {
    helm version | Out-Null
    Write-Host "✅ Helm OK" -ForegroundColor Green
} catch {
    Write-Host "❌ Helm não encontrado" -ForegroundColor Red
}

Write-Host "🎉 Verificação concluída!" -ForegroundColor Green
```

---

## 🔧 Deployment Kubernetes - Passo a Passo

#### 1. Criar Cluster k3d
```bash
# Criar cluster com port-forward configurado
k3d cluster create bgc --port "8080:80@loadbalancer"

# Verificar cluster
kubectl get nodes
kubectl cluster-info
```

**Troubleshooting:**
```bash
# Se cluster já existe
k3d cluster delete bgc
k3d cluster create bgc --port "8080:80@loadbalancer"

# Se kubectl não conecta
kubectl config use-context k3d-bgc
```

#### 2. Instalar PostgreSQL
```bash
# Adicionar repositório Bitnami
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Instalar PostgreSQL
helm install bgc-postgres bitnami/postgresql \
  --namespace default \
  --set auth.postgresPassword=bgc123 \
  --set primary.persistence.size=5Gi

# Aguardar ficar pronto (pode demorar 2-3 minutos)
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgresql --timeout=300s
```

**Verificação:**
```bash
# Status do pod
kubectl get pods | grep postgres

# Testar conexão
kubectl run psql-test --rm -it --image bitnami/postgresql:latest -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql \
  -h bgc-postgres -U postgres -c "SELECT version();"
```

#### 3. Aplicar Migrations
```bash
# Criar estrutura de schemas e tabelas
kubectl apply -f deploy/migrations/

# Verificar jobs de migração
kubectl get jobs
kubectl logs job/bgc-migration-0001
```

**Schema criado:**
```sql
-- Verificar estrutura no banco
\dt stg.*    -- Tabelas staging
\dv rpt.*    -- Views reports  
\dm rpt.*    -- Materialized views
```

#### 4. Build e Deploy de Imagens

##### Build Local
```bash
# API (Clean Architecture)
cd api
docker build -t bgc/api:dev .
cd ..

# Ingest  
cd services/bgc-ingest
docker build -t bgc/ingest:dev .
cd ../..

# Verificar imagens
docker images | grep bgc
```

##### Import no k3d
```bash
# Importar ambas as imagens
k3d image import bgc/api:dev bgc/ingest:dev -c bgc

# Verificar no cluster
k3d image list -c bgc
```

##### Deploy Kubernetes
```bash
# Deploy API
kubectl apply -f deploy/api/

# Deploy Ingest (CronJob)
kubectl apply -f deploy/ingest/

# Verificar deployments
kubectl get deployments
kubectl get cronjobs
kubectl get pods
```

#### 5. Carregar Dados de Exemplo
```bash
# Executar job de carga
kubectl create job load-sample-$(date +%s) \
  --from=cronjob/bgc-ingest \
  -- load-xlsx /data/sample.xlsx

# Verificar progresso
kubectl get jobs
kubectl logs job/load-sample-<timestamp>

# Validar dados carregados
kubectl run psql-client --rm -it --image bitnami/postgresql:latest -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql \
  -h bgc-postgres -U postgres \
  -c "SELECT COUNT(*) FROM stg.exportacao;"
```

#### 6. Refresh Materialized Views
```bash
# Job para popular MVs
kubectl create job refresh-mv-$(date +%s) \
  --from=cronjob/bgc-ingest \
  -- refresh-mv

# Verificar logs
kubectl logs job/refresh-mv-<timestamp>

# Validar MVs populadas
# SQL: SELECT COUNT(*) FROM rpt.mv_resumo_pais;
```

#### 7. Teste da API
```bash
# Port-forward (execute em terminal separado)
kubectl port-forward service/bgc-api 3000:3000

# Em outro terminal, testar endpoints
curl http://localhost:3000/metrics/resumo
curl http://localhost:3000/metrics/pais?limit=5
```

---

## 🏗️ Arquitetura Clean - Estrutura para Desenvolvimento

A API foi refatorada para seguir Clean Architecture (Hexagonal). Ao desenvolver:

**Adicionar novo domínio:**
```bash
# 1. Criar estrutura
mkdir -p api/internal/business/novo_dominio
mkdir -p api/internal/repository/postgres

# 2. Criar arquivos necessários
# - api/internal/business/novo_dominio/entities.go
# - api/internal/business/novo_dominio/repository.go (interface)
# - api/internal/business/novo_dominio/service.go
# - api/internal/repository/postgres/novo_dominio.go (implementação)
# - api/internal/api/handlers/novo_dominio.go

# 3. Registrar em api/internal/app/server.go
```

**Modificar lógica de negócio:**
- Editar `api/internal/business/{domain}/service.go`
- Lógica de negócio NUNCA vai em handlers ou repositories

**Adicionar endpoint:**
- Criar handler em `api/internal/api/handlers/`
- Registrar rota em `api/internal/app/server.go`

**Modificar queries SQL:**
- Editar `api/internal/repository/postgres/{domain}.go`
- Queries SQL NUNCA vão em services

---

## 🔄 Workflows de Desenvolvimento

### Workflow Diário - Docker Compose
```bash
# 1. Iniciar ambiente
cd bgcstack
docker compose up -d

# 2. Fazer alterações no código da API
cd ../api

# 3. Rebuild e restart
cd ../bgcstack
docker compose build api
docker compose restart api

# 4. Ver logs
docker compose logs -f api

# 5. Testar mudanças
curl http://localhost:8080/health
```

### Workflow Diário - Kubernetes
```bash
# 1. Verificar ambiente (se acabou de ligar o PC)
kubectl get nodes

# 2. Status atual
kubectl get pods

# 3. Fazer alterações no código
# ... edit services/api/main.go ...

# 4. Build e deploy
.\bgc.ps1 dev

# 5. Testar
.\bgc.ps1 port-forward  # terminal separado
.\bgc.ps1 test
```

### Workflow de Dados
```bash
# Carregar novos dados
kubectl create job load-new-data-$(date +%s) \
  --from=cronjob/bgc-ingest \
  -- load-csv /data/new-file.csv

# Atualizar análises
kubectl create job refresh-$(date +%s) \
  --from=cronjob/bgc-ingest \
  -- refresh-mv

# Verificar resultado
.\bgc.ps1 test
```

### Workflow de Limpeza
```bash
# Limpar jobs antigos
.\bgc.ps1 clean-jobs

# Reset completo (CUIDADO!)
.\bgc.ps1 clean
.\bgc.ps1 setup
```

---

## 📁 Estrutura de Deploy

### Organização de Arquivos
```
deploy/
├── api/
│   ├── deployment.yaml        # BGC API deployment
│   ├── service.yaml          # Service ClusterIP
│   └── configmap.yaml        # Configurações da API
├── ingest/
│   ├── cronjob.yaml          # Job recorrente de ingest
│   ├── configmap-data.yaml   # Dados de exemplo
│   └── pvc.yaml              # Storage para arquivos
└── migrations/
    ├── job-0001.yaml         # Job inicial de schema
    ├── configmap-0001.yaml   # SQL de migração
    └── job-0002.yaml         # Próximas migrations
```

### Templates Kubernetes

#### API Deployment
```yaml
# deploy/api/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bgc-api
  labels:
    app: bgc-api
spec:
  replicas: 1
  selector:
    matchLabels:
      app: bgc-api
  template:
    metadata:
      labels:
        app: bgc-api
    spec:
      containers:
      - name: api
        image: bgc/api:dev
        ports:
        - containerPort: 3000
        env:
        - name: DB_HOST
          value: bgc-postgres
        - name: DB_USER
          value: postgres
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: bgc-postgres
              key: postgres-password
        - name: DB_NAME
          value: postgres
        - name: PORT
          value: "3000"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        # TODO: Health probes
        # livenessProbe:
        #   httpGet:
        #     path: /health
        #     port: 3000
        # readinessProbe:
        #   httpGet:
        #     path: /ready
        #     port: 3000
```

#### Ingest CronJob
```yaml
# deploy/ingest/cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: bgc-ingest
spec:
  # Não executa automaticamente - apenas template para jobs manuais
  schedule: "0 2 * * *"  # 02:00 daily (disabled)
  suspend: true
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: ingest
            image: bgc/ingest:dev
            env:
            - name: DB_HOST
              value: bgc-postgres
            - name: DB_USER
              value: postgres
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: bgc-postgres
                  key: postgres-password
            - name: DB_NAME
              value: postgres
            volumeMounts:
            - name: data-volume
              mountPath: /data
            resources:
              requests:
                memory: "256Mi"
                cpu: "200m"
              limits:
                memory: "1Gi"
                cpu: "1000m"
          volumes:
          - name: data-volume
            configMap:
              name: bgc-sample-data
          restartPolicy: OnFailure
```

---

## 🔧 Configurações Avançadas

### PostgreSQL Customizado
```bash
# Instalar com configurações específicas
helm install bgc-postgres bitnami/postgresql \
  --set auth.postgresPassword=bgc123 \
  --set primary.persistence.size=10Gi \
  --set primary.resources.requests.memory=512Mi \
  --set primary.resources.requests.cpu=250m \
  --set metrics.enabled=true \
  --set metrics.serviceMonitor.enabled=false
```

### Recursos Customizados
```yaml
# Configurações de recursos para desenvolvimento
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "256Mi" 
    cpu: "200m"

# Configurações para staging/produção
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

### Storage Personalizado
```yaml
# PVC para dados permanentes
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: bgc-data-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: local-path  # k3d default
```

---

## 🔍 Verificação e Validação

### Health Check Completo
```bash
#!/bin/bash
# health-check.sh - Script completo de verificação

echo "🏥 BGC Health Check Completo"
echo "================================"

# 1. Cluster Kubernetes
echo "1. Verificando cluster..."
kubectl get nodes | grep Ready > /dev/null && echo "✅ Cluster OK" || echo "❌ Cluster FAIL"

# 2. PostgreSQL
echo "2. Verificando PostgreSQL..."
kubectl get pods | grep bgc-postgres | grep Running > /dev/null && echo "✅ PostgreSQL OK" || echo "❌ PostgreSQL FAIL"

# 3. API
echo "3. Verificando API..."
kubectl get pods | grep bgc-api | grep Running > /dev/null && echo "✅ API Pod OK" || echo "❌ API Pod FAIL"

# 4. Conectividade da API (requer port-forward ativo)
echo "4. Testando conectividade API..."
curl -s http://localhost:3000/metrics/resumo > /dev/null && echo "✅ API Response OK" || echo "⚠️ API não acessível (port-forward?)"

# 5. Dados no banco
echo "5. Verificando dados..."
RECORD_COUNT=$(kubectl run psql-count --rm -i --image bitnami/postgresql:latest -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql \
  -h bgc-postgres -U postgres -t -c "SELECT COUNT(*) FROM stg.exportacao;" 2>/dev/null | tr -d ' ')

if [ "$RECORD_COUNT" -gt 0 ]; then
    echo "✅ Dados OK ($RECORD_COUNT registros)"
else
    echo "⚠️ Sem dados ou erro de conexão"
fi

# 6. Materialized Views
echo "6. Verificando MVs..."
kubectl run psql-mv --rm -i --image bitnami/postgresql:latest -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql \
  -h bgc-postgres -U postgres -c "\dm rpt.*" 2>/dev/null | grep mv_ > /dev/null && echo "✅ MVs OK" || echo "⚠️ MVs não encontradas"

echo ""
echo "🎯 Health check concluído!"
```

### Performance Básico
```bash
# Testar tempo de resposta da API
time curl -s http://localhost:3000/metrics/resumo > /dev/null

# Verificar uso de recursos dos pods
kubectl top pods

# Logs recentes (últimos 50 lines)
kubectl logs deployment/bgc-api --tail=50
kubectl logs deployment/bgc-postgres --tail=20
```

---

## 🚨 Troubleshooting de Deploy

### Problema: Pods em CrashLoopBackOff
```bash
# Diagnóstico
kubectl describe pod <pod-name>
kubectl logs <pod-name> --previous

# Soluções comuns
# 1. Problema de imagem
k3d image import bgc/api:dev -c bgc
kubectl rollout restart deployment/bgc-api

# 2. Problema de configuração
kubectl get configmap
kubectl describe configmap <config-name>

# 3. Problema de conectividade DB
kubectl get secret bgc-postgres -o yaml
kubectl run db-test --rm -it --image bitnami/postgresql:latest -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql -h bgc-postgres -U postgres
```

### Problema: ImagePullBackOff
```bash
# Verificar se imagem existe localmente
docker images | grep bgc

# Verificar se foi importada no k3d
k3d image list -c bgc

# Solução
docker build -t bgc/api:dev services/api/
k3d image import bgc/api:dev -c bgc
```

### Problema: Service/Port-forward não funciona
```bash
# Verificar service
kubectl get svc bgc-api
kubectl describe svc bgc-api

# Verificar endpoints
kubectl get endpoints bgc-api

# Port-forward direto no pod
kubectl get pods | grep bgc-api
kubectl port-forward pod/<pod-name> 3000:3000
```

### Problema: Banco de dados inacessível
```bash
# Status detalhado do PostgreSQL
kubectl describe pod <postgres-pod>

# Verificar logs do PostgreSQL
kubectl logs <postgres-pod>

# Testar conectividade interna
kubectl run netshoot --rm -it --image nicolaka/netshoot -- \
  nslookup bgc-postgres.default.svc.cluster.local

# Verificar secret da senha
kubectl get secret bgc-postgres -o jsonpath="{.data.postgres-password}" | base64 -d
```

---

## 📊 Monitoramento de Deploy

### Logs Importantes
```bash
# API logs com contexto
kubectl logs deployment/bgc-api -f --tail=100

# PostgreSQL logs (init + runtime)
kubectl logs deployment/bgc-postgres -f --tail=50

# Jobs de ingest (últimos 3)
kubectl get jobs --sort-by=.metadata.creationTimestamp | tail -3
kubectl logs job/<latest-job>

# Eventos do cluster
kubectl get events --sort-by=.metadata.creationTimestamp | tail -10
```

### Métricas de Recursos
```bash
# Uso atual de recursos
kubectl top nodes
kubectl top pods

# Limites configurados
kubectl describe deployment bgc-api | grep -A 10 "Limits\|Requests"

# Storage utilizado
kubectl get pv
kubectl get pvc
```

### Status Dashboard Simples
```bash
# Dashboard em uma linha
watch 'kubectl get pods | grep -E "(bgc-|postgres)" && echo "---" && kubectl get svc | grep bgc'

# Ou versão PowerShell
while($true) { 
    kubectl get pods | Select-String "(bgc-|postgres)"
    Write-Host "---"
    kubectl get svc | Select-String "bgc"
    Start-Sleep 5
    Clear-Host
}
```

---

## 🔄 Updates e Maintenance

### Rolling Updates
```bash
# Build nova versão
docker build -t bgc/api:v1.1 services/api/
k3d image import bgc/api:v1.1 -c bgc

# Update do deployment
kubectl set image deployment/bgc-api api=bgc/api:v1.1

# Verificar rollout
kubectl rollout status deployment/bgc-api

# Rollback se necessário
kubectl rollout undo deployment/bgc-api
```

### Backup Básico
```bash
# Backup de dados (desenvolvimento)
kubectl exec deployment/bgc-postgres -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/pg_dump \
  -U postgres -d postgres > backup-$(date +%Y%m%d).sql

# Restore (se necessário)
kubectl exec -i deployment/bgc-postgres -- \
  /opt/bitnami/scripts/postgresql/entrypoint.sh \
  /opt/bitnami/postgresql/bin/psql \
  -U postgres -d postgres < backup-20250907.sql
```

### Cleanup Periódico
```bash
# Limpar jobs antigos (manter últimos 5)
kubectl get jobs --sort-by=.metadata.creationTimestamp -o name | head -n -5 | xargs kubectl delete

# Limpar imagens não utilizadas no Docker host
docker image prune -f

# Restart completo (desenvolvimento)
kubectl delete pods --all
