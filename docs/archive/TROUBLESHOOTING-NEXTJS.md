# Troubleshooting - Next.js Web App

## Problema: Erro 500 e "EINVAL: invalid argument, readlink"

### Sintomas
- Servidor Next.js inicia mas retorna erro 500
- Log mostra: `Error: EINVAL: invalid argument, readlink 'C:\...\web-next\.next\static\chunks'`
- Ocorre em Windows, especialmente com OneDrive

### Causa
O cache do Next.js (pasta `.next`) está corrompido ou há conflito com sincronização do OneDrive.

### Solução

**Execute este comando no PowerShell:**

```powershell
# 1. Navegue até o diretório
cd "C:\Users\rafae\OneDrive\Documentos\Projetos\Brasil Global Conect\bgc-app\web-next"

# 2. Mate todos os processos Node.js
Get-Process node -ErrorAction SilentlyContinue | Stop-Process -Force

# 3. Remova completamente a pasta .next
Remove-Item -Recurse -Force .next -ErrorAction SilentlyContinue

# 4. Limpe o cache do pnpm
pnpm store prune

# 5. Reinicie o servidor
pnpm dev
```

### Solução Alternativa: Mover Projeto para Fora do OneDrive

Se o problema persistir, o OneDrive pode estar interferindo. Mova o projeto para `C:\Dev\`:

```powershell
# Criar diretório local
New-Item -ItemType Directory -Path "C:\Dev" -Force

# Copiar projeto
Copy-Item -Recurse "C:\Users\rafae\OneDrive\Documentos\Projetos\Brasil Global Conect\bgc-app" "C:\Dev\bgc-app"

# Navegar e iniciar
cd "C:\Dev\bgc-app\web-next"
pnpm install
pnpm dev
```

---

## Problema: Porta 3000 já em uso

### Sintomas
- `Error: listen EADDRINUSE: address already in use :::3000`

### Solução

**Opção 1: Matar processo na porta 3000**
```powershell
# Encontrar PID
netstat -ano | findstr :3000

# Matar processo (substitua PID)
Stop-Process -Id <PID> -Force
```

**Opção 2: Usar porta alternativa**
```powershell
pnpm dev --port 3001
```

---

## Problema: API Go não está rodando

### Sintomas
- Erro ao fazer requisições para `/market/size` ou `/routes/compare`
- Console mostra erros de fetch

### Solução

```bash
# Em outro terminal, inicie a API Go
cd api
go run main.go

# A API deve estar em http://localhost:8080
# Verifique com:
curl http://localhost:8080/healthz
```

---

## Problema: Tailwind CSS não está funcionando

### Sintomas
- Estilos não são aplicados
- Classes Tailwind não funcionam

### Solução

```powershell
cd web-next

# Verificar se @tailwindcss/postcss está instalado
pnpm list @tailwindcss/postcss

# Reinstalar se necessário
pnpm install --force

# Limpar e reiniciar
Remove-Item -Recurse -Force .next
pnpm dev
```

---

## Problema: Erro de TypeScript

### Sintomas
- Erros de compilação TypeScript
- Tipos não encontrados

### Solução

```powershell
cd web-next

# Reinstalar dependências
Remove-Item -Recurse -Force node_modules
pnpm install

# Verificar tsconfig.json
cat tsconfig.json
```

---

## Script de Inicialização Rápida

Crie um arquivo `start.ps1` no diretório `web-next`:

```powershell
# start.ps1
Write-Host "Iniciando BGC Web Next.js..." -ForegroundColor Cyan

# Limpar cache
Remove-Item -Recurse -Force .next -ErrorAction SilentlyContinue

# Matar processos node antigos
Get-Process node -ErrorAction SilentlyContinue | Stop-Process -Force

# Aguardar 2 segundos
Start-Sleep -Seconds 2

# Iniciar servidor
pnpm dev
```

**Uso:**
```powershell
cd web-next
.\start.ps1
```

---

## Verificar se tudo está funcionando

### 1. Health Check
```powershell
Invoke-WebRequest -Uri "http://localhost:3000/api/health" -UseBasicParsing
```

Deve retornar:
```json
{"status":"ok","timestamp":"2025-10-19T..."}
```

### 2. Dashboard
Abra no navegador: `http://localhost:3000/`

### 3. Rotas
Abra no navegador: `http://localhost:3000/routes`

---

## Logs de Debug

Se o problema persistir, habilite logs de debug:

```powershell
$env:DEBUG="*"
pnpm dev
```

Ou configure no `next.config.ts`:

```typescript
const nextConfig: NextConfig = {
  // ... outras configs

  // Debug
  logging: {
    fetches: {
      fullUrl: true,
    },
  },
};
```

---

## Problemas Conhecidos

### 1. Windows + OneDrive
- O Next.js pode ter problemas com arquivos sincronizados pelo OneDrive
- **Solução**: Mova o projeto para `C:\Dev\`

### 2. Turbopack no Windows
- O Turbopack tem bugs no Windows
- **Solução**: Desabilitado no `package.json` (já corrigido)

### 3. Espaços no Caminho
- Caminhos com espaços podem causar problemas
- **Solução**: Use aspas duplas ao executar comandos

---

## Contato e Suporte

Se nenhuma solução funcionar:

1. Verifique os logs completos
2. Procure erros no console do navegador (F12)
3. Verifique se todas as dependências estão instaladas
4. Tente reinstalar Node.js e pnpm

**Versões testadas:**
- Node.js: v22.20.0
- pnpm: v10.18.3
- Next.js: 15.5.6
