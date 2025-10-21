# BGC Web Next.js - Script de Inicializacao
# Executa limpeza e inicia o servidor

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " BGC Web Next.js - Iniciando..." -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Passo 1: Matar processos Node.js
Write-Host "[1/4] Parando processos Node.js..." -ForegroundColor Yellow
Get-Process node -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 2

# Passo 2: Limpar cache do Next.js
Write-Host "[2/4] Limpando cache do Next.js..." -ForegroundColor Yellow
if (Test-Path ".next") {
    Remove-Item -Recurse -Force .next -ErrorAction SilentlyContinue
}

# Passo 3: Verificar dependencias
Write-Host "[3/4] Verificando dependencias..." -ForegroundColor Yellow
if (-not (Test-Path "node_modules")) {
    Write-Host "      Instalando dependencias (pode demorar)..." -ForegroundColor Yellow
    pnpm install
}

# Passo 4: Iniciar servidor
Write-Host "[4/4] Iniciando servidor Next.js..." -ForegroundColor Yellow
Write-Host ""
Write-Host "=====================================" -ForegroundColor Green
Write-Host " Servidor rodando!" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""
Write-Host "  Dashboard: http://localhost:3000" -ForegroundColor White
Write-Host "  Rotas:     http://localhost:3000/routes" -ForegroundColor White
Write-Host ""
Write-Host "  Pressione Ctrl+C para parar" -ForegroundColor Yellow
Write-Host ""

pnpm dev
