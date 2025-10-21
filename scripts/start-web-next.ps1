# Script para iniciar o Next.js Web App
# Uso: .\scripts\start-web-next.ps1

Write-Host "🚀 BGC Web Next.js - Iniciando servidor de desenvolvimento" -ForegroundColor Cyan
Write-Host ""

# Navegar para o diretório
$webNextDir = Join-Path $PSScriptRoot "..\web-next"
Set-Location $webNextDir

# Verificar se node_modules existe
if (-not (Test-Path "node_modules")) {
    Write-Host "📦 Instalando dependências..." -ForegroundColor Yellow
    pnpm install
}

# Limpar cache do Next.js
Write-Host "🧹 Limpando cache do Next.js..." -ForegroundColor Yellow
if (Test-Path ".next") {
    Remove-Item -Recurse -Force ".next"
}

# Verificar se a API está rodando
Write-Host "🔍 Verificando se a API Go está rodando..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/healthz" -TimeoutSec 2 -ErrorAction Stop
    Write-Host "✅ API Go está rodando em localhost:8080" -ForegroundColor Green
} catch {
    Write-Host "⚠️  ATENÇÃO: API Go não está respondendo em localhost:8080" -ForegroundColor Red
    Write-Host "   Você precisa iniciar a API Go primeiro!" -ForegroundColor Red
    Write-Host "   Execute: cd api && go run main.go" -ForegroundColor Yellow
    Write-Host ""
    $continue = Read-Host "Deseja continuar mesmo assim? (s/N)"
    if ($continue -ne "s" -and $continue -ne "S") {
        exit 1
    }
}

Write-Host ""
Write-Host "🌐 Iniciando servidor Next.js..." -ForegroundColor Cyan
Write-Host ""
Write-Host "📍 URLs disponíveis:" -ForegroundColor Green
Write-Host "   - Local:   http://localhost:3000" -ForegroundColor White
Write-Host "   - Network: http://192.168.56.1:3000" -ForegroundColor White
Write-Host ""
Write-Host "📄 Páginas:" -ForegroundColor Green
Write-Host "   - Dashboard:  http://localhost:3000/" -ForegroundColor White
Write-Host "   - Rotas:      http://localhost:3000/routes" -ForegroundColor White
Write-Host ""
Write-Host "⏹️  Pressione Ctrl+C para parar o servidor" -ForegroundColor Yellow
Write-Host ""

# Iniciar servidor
pnpm dev
