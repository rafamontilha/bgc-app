# Scripts BGC - Guia de Uso

Pasta centralizada com todos os scripts de automação do projeto BGC.

---

## 📋 Índice

- [Setup Inicial](#-setup-inicial)
- [Desenvolvimento Local](#-desenvolvimento-local)
- [Docker Compose](#-docker-compose)
- [Kubernetes](#-kubernetes)
- [Gerenciamento de Dados](#-gerenciamento-de-dados)

---

## 🚀 Setup Inicial

### `check-environment.ps1`
**Propósito:** Verificar se todas as ferramentas necessárias estão instaladas

**Quando usar:** Primeira vez configurando o ambiente ou após reinstalar o sistema

```powershell
.\scripts\check-environment.ps1
```

**Verifica:**
- Node.js (v20+ ou v22+)
- npm e pnpm
- Go (v1.23+)
- Git
- Docker e Docker Compose
- PowerShell
- VS Code (opcional)

---

### `install-go.ps1`
**Propósito:** Instalar Go automaticamente se não estiver presente

**Quando usar:** Quando `check-environment.ps1` indicar que Go está faltando

```powershell
.\scripts\install-go.ps1
```

**Ação:** Baixa e instala Go 1.23+ automaticamente

⚠️ **Importante:** Reinicie o terminal após a instalação!

---

### `setup-hosts.ps1`
**Propósito:** Configurar arquivo hosts para acessar aplicação via domínios locais

**Quando usar:** Apenas para ambiente Kubernetes (k3d)

```powershell
# DEVE ser executado como Administrador!
.\scripts\setup-hosts.ps1
```

**Adiciona:**
```
127.0.0.1  api.bgc.local
127.0.0.1  web.bgc.local
```

---

## 💻 Desenvolvimento Local

### `start-api.ps1`
**Propósito:** Iniciar API Go localmente (SEM Docker)

**Quando usar:** Desenvolvimento local direto

```powershell
.\scripts\start-api.ps1
```

**Pré-requisitos:**
- Go instalado
- PostgreSQL rodando em `localhost:5432`
  - User: `bgc`
  - Password: `bgc`
  - Database: `bgc`

**Acesso:** http://localhost:8080/healthz

---

### `start-web-next.ps1`
**Propósito:** Iniciar aplicação Next.js localmente (SEM Docker)

**Quando usar:** Desenvolvimento local do front-end

```powershell
.\scripts\start-web-next.ps1
```

**Pré-requisitos:**
- Node.js e pnpm instalados
- API rodando em `localhost:8080`

**Acesso:** http://localhost:3000

---

### `test-web-next.ps1`
**Propósito:** Testar se Next.js está respondendo corretamente

**Quando usar:** Validar deployment ou CI/CD

```powershell
.\scripts\test-web-next.ps1

# Ou testar porta customizada:
.\scripts\test-web-next.ps1 -Port 3001
```

**Testa:**
- `/api/health` - Health check
- `/` - Dashboard
- `/routes` - Página de rotas

---

## 🐳 Docker Compose

### `docker.ps1`
**Propósito:** Gerenciar ambiente Docker Compose completo

**Comandos disponíveis:**

```powershell
# Iniciar todos os serviços
.\scripts\docker.ps1 up

# Parar todos os serviços
.\scripts\docker.ps1 down

# Reiniciar serviços
.\scripts\docker.ps1 restart

# Ver logs
.\scripts\docker.ps1 logs

# Status dos containers
.\scripts\docker.ps1 ps

# Rebuild images e iniciar
.\scripts\docker.ps1 build

# Limpar tudo (remove volumes!)
.\scripts\docker.ps1 clean

# Ver ajuda
.\scripts\docker.ps1 help
```

**Serviços inclusos:**
- PostgreSQL (`bgc_db`) - Porta 5432
- API Go (`bgc_api`) - Porta 8080
- Web Next.js (`bgc_web`) - Porta 3000
- PgAdmin (`bgc_pgadmin`) - Porta 5050

**URLs:**
- Dashboard: http://localhost:3000
- API: http://localhost:8080/healthz
- PgAdmin: http://localhost:5050 (`admin@bgc.dev` / `admin`)

**Credenciais DB (consistentes em todos ambientes):**
- Host: `db` (Docker) ou `localhost` (local)
- Port: `5432`
- User: `bgc`
- Password: `bgc`
- Database: `bgc`

---

## ☸️ Kubernetes

### `k8s.ps1`
**Propósito:** Gerenciar ambiente Kubernetes (k3d) completo

**Comandos disponíveis:**

```powershell
# Setup inicial (criar cluster + deploy)
.\scripts\k8s.ps1 setup

# Deploy em cluster existente
.\scripts\k8s.ps1 up

# Remover deployments (mantém cluster)
.\scripts\k8s.ps1 down

# Reiniciar pods
.\scripts\k8s.ps1 restart

# Ver logs
.\scripts\k8s.ps1 logs

# Status do cluster e pods
.\scripts\k8s.ps1 status

# Rebuild images e redeploy
.\scripts\k8s.ps1 build

# Configurar hosts e abrir browser
.\scripts\k8s.ps1 open

# Deletar cluster completo
.\scripts\k8s.ps1 clean

# Ver ajuda
.\scripts\k8s.ps1 help
```

**URLs (após executar `setup-hosts.ps1` como Admin):**
- Dashboard: http://web.bgc.local
- Rotas: http://web.bgc.local/routes
- API: http://api.bgc.local/healthz

**Features incluídas:**
- HPA (Horizontal Pod Autoscaler)
- Health Probes (readiness/liveness)
- Ingress com Traefik
- Backup automático diário (CronJob)
- Refresh de Materialized Views (CronJob)

---

## 📊 Gerenciamento de Dados

### `seed.ps1` / `seed.sh`
**Propósito:** Carregar dados de exemplo no banco

**Quando usar:** Primeira vez ou reset de dados

```powershell
# PowerShell (Windows)
.\scripts\seed.ps1

# Bash (Linux/Mac)
./scripts/seed.sh
```

**Pré-requisitos:**
- Docker Compose rodando
- Arquivos CSV em `stage/`:
  - `lookup_ncm8.csv`
  - `exports_ncm_year_sample.csv`
  - `imports_ncm_year_sample.csv`

**Ação:**
1. Trunca tabelas existentes
2. Carrega dados de lookup NCM
3. Carrega dados de exportação
4. Carrega dados de importação
5. Atualiza Materialized Views

---

### `restore-backup.ps1`
**Propósito:** Restaurar backup do PostgreSQL (apenas Kubernetes)

**Quando usar:** Recuperar dados após falha ou migrar dados

```powershell
# Listar backups disponíveis
.\scripts\restore-backup.ps1

# Restaurar backup específico
.\scripts\restore-backup.ps1 -BackupFile bgc_backup_20251021_020000.sql.gz
```

**Funciona apenas com:** Kubernetes (k3d)

**Backups automáticos:** Diariamente às 02:00 (CronJob)

---

## 🔄 Fluxos de Trabalho Comuns

### Primeira Vez no Projeto

```powershell
# 1. Verificar ambiente
.\scripts\check-environment.ps1

# 2. Instalar Go (se necessário)
.\scripts\install-go.ps1

# 3. Iniciar com Docker Compose
.\scripts\docker.ps1 up

# 4. Carregar dados de exemplo
.\scripts\seed.ps1

# 5. Acessar
# http://localhost:3000
```

---

### Desenvolvimento Local (sem Docker)

```powershell
# Terminal 1: Iniciar API
.\scripts\start-api.ps1

# Terminal 2: Iniciar Web
.\scripts\start-web-next.ps1

# Testar
.\scripts\test-web-next.ps1
```

---

### Testar com Kubernetes

```powershell
# 1. Setup inicial (apenas primeira vez)
.\scripts\k8s.ps1 setup

# 2. Configurar hosts (como Admin)
.\scripts\setup-hosts.ps1

# 3. Abrir no browser
.\scripts\k8s.ps1 open

# Ver status
.\scripts\k8s.ps1 status

# Ver logs
.\scripts\k8s.ps1 logs
```

---

### Rebuild após Mudanças no Código

**Docker Compose:**
```powershell
.\scripts\docker.ps1 build
```

**Kubernetes:**
```powershell
.\scripts\k8s.ps1 build
```

---

## 📚 Notas Importantes

### Consistência de Credenciais
Todos os ambientes usam as **mesmas credenciais** do banco de dados:
- User: `bgc`
- Password: `bgc`
- Database: `bgc`

Isso garante paridade entre dev, Docker e Kubernetes.

### Portas Padrão
- **API Go:** 8080
- **Web Next.js:** 3000
- **PostgreSQL:** 5432
- **PgAdmin:** 5050 (apenas Docker Compose)

### Scripts Removidos
- ❌ `configure-hosts.ps1` - Use `setup-hosts.ps1`
- ❌ `setup-nextjs.ps1` - Projeto já existe

---

## 🆘 Troubleshooting

**API não conecta ao banco:**
- Verifique se PostgreSQL está rodando
- Verifique credenciais (user: `bgc`, password: `bgc`)
- Verifique porta 5432

**Web não conecta à API:**
- Verifique se API está rodando (porta 8080)
- Teste: `curl http://localhost:8080/healthz`

**Kubernetes: web.bgc.local não resolve:**
- Execute `.\scripts\setup-hosts.ps1` como Administrador
- Verifique `C:\Windows\System32\drivers\etc\hosts`

**Porta em uso:**
```powershell
# Ver o que está usando a porta
netstat -ano | findstr :8080

# Parar serviços
.\scripts\docker.ps1 down
```

---

## 📞 Ajuda

Para mais informações, consulte:
- `docs/QUICK-START.md` - Guia rápido
- `docs/TROUBLESHOOTING-NEXTJS.md` - Problemas comuns
- `README.md` - Documentação principal
