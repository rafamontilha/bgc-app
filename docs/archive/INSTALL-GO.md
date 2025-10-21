# Como Instalar Go no Windows

A API do BGC App requer **Go 1.23+** para funcionar.

---

## Opção 1: Instalação Automática (Recomendado)

Execute o script PowerShell que criamos:

```powershell
.\scripts\install-go.ps1
```

O script irá:
1. ✅ Verificar se Go já está instalado
2. ✅ Baixar Go 1.23.2 para Windows
3. ✅ Instalar automaticamente
4. ✅ Configurar PATH

**Após a instalação:**
1. Feche o terminal atual
2. Abra um NOVO terminal
3. Verifique: `go version`

---

## Opção 2: Instalação Manual

### 1. Download

Acesse: https://go.dev/dl/

Baixe: **go1.23.2.windows-amd64.msi**

### 2. Instalação

1. Execute o arquivo `.msi` baixado
2. Siga o assistente de instalação
3. Use as configurações padrão

### 3. Verificação

Abra um NOVO terminal e execute:

```powershell
go version
```

**Resultado esperado:**
```
go version go1.23.2 windows/amd64
```

---

## Opção 3: Via Chocolatey

Se você tem Chocolatey instalado:

```powershell
choco install golang -y
```

---

## Opção 4: Via Scoop

Se você tem Scoop instalado:

```powershell
scoop install go
```

---

## Verificar Instalação

Após instalar, verifique se tudo está OK:

```powershell
# Versão
go version

# Variáveis de ambiente
go env GOPATH
go env GOROOT

# Testar compilação
go run --help
```

---

## Configuração de Ambiente (Opcional)

### GOPATH

Por padrão, Go usa `%USERPROFILE%\go` como GOPATH.

Para verificar:
```powershell
go env GOPATH
```

Para alterar:
```powershell
[System.Environment]::SetEnvironmentVariable("GOPATH", "C:\Dev\go", "User")
```

### Proxy (Se necessário)

Se estiver atrás de um proxy:

```powershell
go env -w GOPROXY=https://proxy.golang.org,direct
go env -w GOSUMDB=sum.golang.org
```

---

## Após Instalar Go

### 1. Baixar Dependências da API

```powershell
cd api
go mod download
```

### 2. Iniciar API

**Opção A: Script automatizado**
```powershell
.\scripts\start-api.ps1
```

**Opção B: Manual**
```powershell
cd api
go run cmd/api/main.go
```

### 3. Testar API

```powershell
# PowerShell
Invoke-WebRequest -Uri "http://localhost:8080/healthz" -UseBasicParsing

# Ou no navegador
http://localhost:8080/healthz
```

**Resposta esperada:**
```json
{"status":"ok","timestamp":"2025-10-19T..."}
```

---

## Problemas Comuns

### Go não é reconhecido após instalação

**Solução:**
1. Feche **TODOS** os terminais abertos
2. Abra um NOVO terminal
3. Tente novamente: `go version`

Se ainda não funcionar:
1. Verifique se Go está no PATH:
   ```powershell
   $env:Path
   ```
2. Adicione manualmente se necessário:
   ```powershell
   [System.Environment]::SetEnvironmentVariable(
       "Path",
       $env:Path + ";C:\Program Files\Go\bin",
       "User"
   )
   ```

### Erro: go.mod not found

**Solução:**
```powershell
cd api
go mod init bgc-app
go mod tidy
```

### Erro ao baixar dependências

**Solução:**
```powershell
# Limpar cache
go clean -modcache

# Baixar novamente
go mod download
```

### Erro de proxy/firewall

**Solução:**
```powershell
# Configurar proxy
go env -w GOPROXY=https://goproxy.io,direct

# OU desabilitar checksum (não recomendado em produção)
go env -w GOSUMDB=off
```

---

## Versões Testadas

| Ferramenta | Versão | Status |
|------------|--------|--------|
| Go | 1.23.2 | ✅ Testado |
| Go | 1.23.x | ✅ Compatível |
| Go | 1.22.x | ⚠️ Pode funcionar |
| Go | < 1.21 | ❌ Não suportado |

---

## Próximos Passos

Após instalar Go:

1. ✅ Iniciar API Go: `.\scripts\start-api.ps1`
2. ✅ Iniciar Web Next.js: `cd web-next && .\start.ps1`
3. ✅ Acessar: http://localhost:3000

---

## Recursos

- 📖 Documentação oficial: https://go.dev/doc/
- 📦 Download: https://go.dev/dl/
- 🎓 Tour of Go: https://go.dev/tour/
- 📚 Go by Example: https://gobyexample.com/

---

## Ajuda

Se precisar de ajuda:

1. Verifique a instalação: `go version`
2. Verifique o PATH: `echo $env:Path`
3. Consulte: `docs/TROUBLESHOOTING-NEXTJS.md`
4. Reinstale Go se necessário
