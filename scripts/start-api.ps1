# Script para iniciar a API Go
# Executa: .\scripts\start-api.ps1

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " BGC API Go - Iniciando..." -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Verificar se Go esta instalado
Write-Host "[1/3] Verificando Go..." -ForegroundColor Yellow
try {
    $goVersion = & go version 2>$null
    if ($goVersion) {
        Write-Host "      $goVersion" -ForegroundColor Green
    } else {
        throw "Go nao encontrado"
    }
} catch {
    Write-Host ""
    Write-Host "ERRO: Go nao esta instalado!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Opcoes:" -ForegroundColor Yellow
    Write-Host "  1. Instalar via script: .\scripts\install-go.ps1" -ForegroundColor White
    Write-Host "  2. Download manual: https://go.dev/dl/" -ForegroundColor White
    Write-Host ""
    exit 1
}

# Navegar para diretorio da API
$apiDir = Join-Path $PSScriptRoot "..\api"
Set-Location $apiDir

# Verificar se go.mod existe
Write-Host "[2/3] Verificando dependencias..." -ForegroundColor Yellow
if (-not (Test-Path "go.mod")) {
    Write-Host "      Erro: go.mod nao encontrado!" -ForegroundColor Red
    exit 1
}

# Baixar dependencias se necessario
if (-not (Test-Path "go.sum")) {
    Write-Host "      Baixando dependencias..." -ForegroundColor Yellow
    go mod download
}

# Verificar se PostgreSQL esta rodando
Write-Host "[3/3] Verificando PostgreSQL..." -ForegroundColor Yellow
try {
    $pgTest = Test-Connection -ComputerName localhost -Port 5432 -Count 1 -ErrorAction Stop
    Write-Host "      PostgreSQL OK" -ForegroundColor Green
} catch {
    Write-Host "      Aviso: PostgreSQL pode nao estar rodando em localhost:5432" -ForegroundColor Yellow
    Write-Host "      A API pode falhar ao conectar ao banco." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=====================================" -ForegroundColor Green
Write-Host " Iniciando API Go..." -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""
Write-Host "  API: http://localhost:8080" -ForegroundColor White
Write-Host "  Health: http://localhost:8080/healthz" -ForegroundColor White
Write-Host "  Docs: http://localhost:8080/docs" -ForegroundColor White
Write-Host ""
Write-Host "  Pressione Ctrl+C para parar" -ForegroundColor Yellow
Write-Host ""

# Iniciar API
go run cmd/api/main.go
