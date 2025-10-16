# BGC App - Sistema de Analytics de Exportação

Sistema de analytics para dados de exportação brasileira com API Go, Frontend Web e PostgreSQL.

## 🚀 Quick Start

### Opção 1: Docker Compose (Recomendado para Desenvolvimento)

```powershell
# Iniciar ambiente
.\scripts\docker.ps1 up

# URLs disponíveis:
# API:     http://localhost:8080
# Web UI:  http://localhost:3000
# PgAdmin: http://localhost:5050 (admin@bgc.dev / admin)
```

### Opção 2: Kubernetes (k3d)

```powershell
# Setup inicial (primeira vez)
.\scripts\k8s.ps1 setup

# Configurar hosts (executar como Administrador)
.\scripts\configure-hosts.ps1

# URLs disponíveis:
# Web UI:  http://web.bgc.local
# Routes:  http://web.bgc.local/routes.html
# API:     http://web.bgc.local/healthz
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
.\scripts\k8s.ps1 status         # Status do cluster
.\scripts\k8s.ps1 build          # Rebuildar imagens
.\scripts\k8s.ps1 open           # Abrir no browser
.\scripts\k8s.ps1 clean          # Deletar cluster
.\scripts\k8s.ps1 help           # Ajuda
```

### Seed de Dados (Opcional)

```powershell
# Carregar dados de exemplo
.\scripts\seed.ps1
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
├── web/                      # Frontend (HTML/JS/CSS)
│   ├── index.html           # Dashboard TAM/SAM/SOM
│   ├── routes.html          # Comparação de rotas
│   ├── nginx.conf           # Configuração Nginx
│   └── Dockerfile
│
├── services/                 # Microserviços auxiliares
│   └── bgc-ingest/          # Serviço de ingestão
│
├── db/                       # Database
│   ├── init/                # Schema inicial (Docker Compose)
│   └── migrations/          # Migrations SQL
│
├── deploy/                   # Kubernetes Manifests (jobs, migrations)
│   ├── postgres.yaml
│   ├── bgc-api.yaml
│   └── configmap-*.yaml
│
├── k8s/                      # Kubernetes Manifests (serviços)
│   ├── api.yaml
│   ├── web.yaml
│   └── web-nginx-configmap.yaml
│
├── bgcstack/                 # Docker Compose
│   └── docker-compose.yml
│
├── scripts/                  # Scripts de automação
│   ├── docker.ps1           # Gerenciar Docker Compose
│   ├── k8s.ps1              # Gerenciar Kubernetes
│   ├── configure-hosts.ps1  # Configurar hosts
│   └── seed.ps1             # Seed de dados
│
└── docs/                     # Documentação técnica
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

## 🔒 Credenciais Padrão

### PostgreSQL

- **Host**: db (Docker) / pg-postgresql (K8s)
- **Port**: 5432
- **User**: bgc
- **Password**: bgc
- **Database**: bgc

### PgAdmin (Docker Compose)

- **URL**: http://localhost:5050
- **Email**: admin@bgc.dev
- **Password**: admin

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

## 🤝 Contribuindo

1. Faça fork do projeto
2. Crie uma branch: `git checkout -b feature/nova-feature`
3. Commit: `git commit -m 'feat: adiciona nova feature'`
4. Push: `git push origin feature/nova-feature`
5. Abra um Pull Request

---

## 📝 Licença

Este projeto está sob licença MIT.

---

## ✨ Features

- ✅ API REST com Clean Architecture
- ✅ Dashboard interativo TAM/SAM/SOM
- ✅ Comparação de rotas de exportação
- ✅ Deploy via Docker Compose
- ✅ Deploy via Kubernetes
- ✅ Health checks e observabilidade
- ✅ Proxy reverso Nginx
- ✅ Migrations automáticas
- ✅ Scripts de gerenciamento simplificados

---

**Desenvolvido com ❤️ pela equipe BGC**
