# Quick Start - BGC App

Guia rápido para iniciar o projeto completo (API + Web).

## Pré-requisitos

- ✅ Node.js v22.20.0
- ✅ **Go 1.23+** (veja como instalar abaixo)
- ✅ pnpm v10+
- ✅ PostgreSQL rodando (ou via Docker)

### ⚠️ Go não instalado?

Se você não tem Go instalado, execute:

```powershell
.\scripts\install-go.ps1
```

Após a instalação, **feche e reabra o terminal**.

**Documentação completa:** `docs/INSTALL-GO.md`

---

## Iniciar o Projeto Completo

### 1️⃣ Terminal 1 - API Go

```bash
cd api
go run main.go
```

**Deve mostrar:**
```
2025/10/19 16:00:00 Server starting on :8080
```

**Testar:**
```bash
curl http://localhost:8080/healthz
```

---

### 2️⃣ Terminal 2 - Web Next.js

```powershell
cd web-next
.\start.ps1
```

**OU manualmente:**
```powershell
cd web-next

# Limpar cache
Remove-Item -Recurse -Force .next -ErrorAction SilentlyContinue

# Matar processos node antigos
Get-Process node -ErrorAction SilentlyContinue | Stop-Process -Force

# Iniciar
pnpm dev
```

**Deve mostrar:**
```
▲ Next.js 15.5.6
- Local:        http://localhost:3000

✓ Ready in 3s
```

---

### 3️⃣ Acessar no Navegador

**Dashboard:**
```
http://localhost:3000/
```

**Rotas:**
```
http://localhost:3000/routes
```

---

## Verificação Rápida

### Teste 1: API está rodando?
```bash
curl http://localhost:8080/healthz
```

**Esperado:**
```json
{"status":"ok","timestamp":"..."}
```

### Teste 2: Web está rodando?
```bash
curl http://localhost:3000/api/health
```

**Esperado:**
```json
{"status":"ok","timestamp":"..."}
```

### Teste 3: Proxy funciona?
```bash
curl "http://localhost:3000/market/size?metric=TAM&year_from=2020&year_to=2025"
```

**Esperado:**
```json
{"metric":"TAM","year_from":2020,"year_to":2025,"items":[...]}
```

---

## Problemas Comuns

### ❌ Erro: API não está rodando

**Sintoma:** Erro 500 nas requisições

**Solução:**
```bash
# Terminal separado
cd api
go run main.go
```

### ❌ Erro: Porta 3000 em uso

**Solução:**
```powershell
# Matar processo
Get-Process node | Stop-Process -Force

# OU usar porta alternativa
pnpm dev --port 3001
```

### ❌ Erro: Cache corrompido (.next)

**Solução:**
```powershell
cd web-next
Remove-Item -Recurse -Force .next
pnpm dev
```

### ❌ Erro: PostgreSQL não está rodando

**Solução:**
```bash
# Via Docker
docker run -d \
  --name bgc-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=bgc \
  -p 5432:5432 \
  postgres:16-alpine

# Executar migrações
cd db
psql -h localhost -U postgres -d bgc -f migrations/0001_initial.sql
```

---

## URLs de Referência

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **Dashboard** | http://localhost:3000 | Interface principal |
| **Rotas** | http://localhost:3000/routes | Comparação de rotas |
| **API Go** | http://localhost:8080 | API REST |
| **API Health** | http://localhost:8080/healthz | Status da API |
| **Web Health** | http://localhost:3000/api/health | Status do Next.js |
| **Swagger** | http://localhost:8080/docs | Documentação da API |

---

## Indicador Visual de Status

A interface web agora mostra um **indicador de status da API** no canto inferior direito:

- 🟢 **Verde**: API Online
- 🔴 **Vermelho**: API Offline (com instruções)
- ⚪ **Cinza**: Verificando...

---

## Comandos Úteis

```bash
# Verificar versões
node --version    # v22.20.0
pnpm --version    # 10.18.3
go version        # go1.21+

# Reinstalar dependências Next.js
cd web-next
Remove-Item -Recurse -Force node_modules
pnpm install

# Recompilar API Go
cd api
go build -o bgc-api.exe main.go
.\bgc-api.exe

# Ver logs do banco
docker logs bgc-postgres -f
```

---

## Próximos Passos

Após iniciar com sucesso:

1. ✅ Teste o Dashboard com filtros diferentes
2. ✅ Teste a página de Rotas
3. ✅ Verifique a exportação CSV
4. ✅ Explore a API via Swagger
5. ✅ Execute testes (se disponíveis)

---

## Ajuda

- 📖 Documentação completa: `docs/`
- 🐛 Troubleshooting: `docs/TROUBLESHOOTING-NEXTJS.md`
- 🚀 Deploy K8s: `docs/DEPLOYMENT.md`
- ⚙️ Setup inicial: `docs/SETUP-NEXTJS.md`
