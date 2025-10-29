# BGC App - Sistema de Analytics de Exportação

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?logo=postgresql)](https://www.postgresql.org)

Sistema de analytics para dados de exportação brasileira com:
- **API REST** em Go (Gin framework)
- **Frontend** em Next.js 15 (React, TypeScript, Tailwind CSS)
- **Banco de Dados** PostgreSQL 16

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
# Web:     http://localhost:3000
# API:     http://localhost:8080
# PgAdmin: http://localhost:5050 (admin@bgc.dev / admin)
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
├── api/                      # API Go (Clean Architecture)
│   ├── cmd/api/             # Entry point
│   ├── config/              # Configurações YAML
│   ├── internal/            # Código interno
│   │   ├── business/       # Lógica de negócio (domain)
│   │   ├── repository/     # Persistência (postgres)
│   │   ├── api/            # Handlers HTTP
│   │   └── app/            # Wiring & server
│   ├── Dockerfile
│   └── go.mod
│
├── web-next/                 # Frontend Next.js 15 (React + TypeScript)
│   ├── app/                 # App Router do Next.js
│   ├── components/          # Componentes React
│   ├── lib/                 # Utilitários e API client
│   ├── hooks/               # Custom React Hooks
│   ├── types/               # TypeScript types
│   ├── Dockerfile
│   └── package.json
│
├── services/                 # Microserviços auxiliares
│   └── bgc-ingest/          # Serviço de ingestão
│
├── db/                       # Database
│   ├── init/                # Schema inicial (Docker Compose)
│   └── migrations/          # Migrations SQL
│
├── k8s/                      # Kubernetes Manifests (serviços)
│   ├── api.yaml             # Deployment API com HPA
│   ├── web.yaml             # Deployment Web Next.js com HPA
│   ├── postgres-backup-cronjob.yaml
│   └── mview-refresh-cronjob.yaml
│
├── deploy/                   # Kubernetes Jobs (migrations, seeds)
│   ├── postgres.yaml
│   └── configmap-*.yaml
│
├── scripts/                  # Scripts de automação
│   ├── k8s.ps1              # Gerenciar Kubernetes
│   ├── setup-hosts.ps1      # Configurar hosts
│   ├── start-api.ps1        # Iniciar API local
│   ├── start-web-next.ps1   # Iniciar Web Next.js local
│   └── test-web-next.ps1    # Testar Web Next.js
│
├── docs/                     # Documentação técnica
│   ├── QUICK-START.md
│   ├── SETUP-NEXTJS.md
│   └── TROUBLESHOOTING-NEXTJS.md
│
├── old/                      # Arquivos legados (histórico)
│   └── web-legacy-html/     # Frontend HTML antigo
│
├── Makefile                  # Wrapper multiplataforma
└── CHANGELOG.md              # Histórico de mudanças
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

- **Backend**: Go 1.23 com Gin
- **Frontend**: HTML5, JavaScript, CSS
- **Database**: PostgreSQL 16
- **Container**: Docker
- **Orchestration**: Kubernetes (k3d)
- **Proxy**: Nginx

---

## 📊 Endpoints da API

### Base URLs

- **Docker Compose**: `http://localhost:8080`
- **Kubernetes**: `http://web.bgc.local` (via Ingress)

### Principais Endpoints

```
GET /healthz                          # Health check
GET /market/size                      # TAM/SAM/SOM metrics
GET /routes/compare                   # Comparação de rotas
GET /docs                             # API documentation
GET /openapi.yaml                     # OpenAPI spec
```

### Exemplos

```bash
# Health check
curl http://localhost:8080/healthz

# Market size (TAM)
curl "http://localhost:8080/market/size?metric=TAM&year_from=2023&year_to=2024"

# Routes compare
curl "http://localhost:8080/routes/compare?from=USA&alts=CHN,ARE&ncm_chapter=84&year=2024"
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

As credenciais no Kubernetes são gerenciadas via **Sealed Secrets**:

```powershell
# Sealed Secrets controller já instalado
kubectl get pods -n kube-system | grep sealed-secrets

# Credenciais criptografadas em: k8s/secrets/
```

### Documentação Completa

📖 **Veja o guia completo:** [docs/SECURITY-SECRETS.md](docs/SECURITY-SECRETS.md)

**🚨 NUNCA commite credenciais no Git!**

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

### Core
- ✅ API REST com Clean Architecture (Go 1.23)
- ✅ Dashboard interativo TAM/SAM/SOM
- ✅ Comparação de rotas de exportação
- ✅ PostgreSQL 16 com Materialized Views
- ✅ Migrations automáticas com rastreabilidade

### Deployment
- ✅ Docker Compose para desenvolvimento
- ✅ Kubernetes (k3d) para produção simulada
- ✅ Scripts PowerShell unificados
- ✅ Makefile multiplataforma
- ✅ Proxy reverso Nginx com Traefik Ingress

### Observabilidade & Resiliência
- ✅ Health probes (readiness/liveness)
- ✅ Horizontal Pod Autoscaling (HPA)
- ✅ Resource limits e requests
- ✅ Backups automáticos diários do PostgreSQL
- ✅ Refresh automático de materialized views
- ✅ Métricas de API (/metrics endpoint)

### DevOps
- ✅ CronJobs para backup e refresh de dados
- ✅ Script de restore de backups
- ✅ Ingress com Traefik
- ✅ ConfigMaps para configuração
- ✅ Secrets para credenciais
- ✅ CHANGELOG.md com versionamento semântico

---

**Desenvolvido com ❤️ pela equipe BGC**
