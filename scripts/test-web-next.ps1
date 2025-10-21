# Script para testar se o Next.js está acessível
# Uso: .\scripts\test-web-next.ps1

param(
    [int]$Port = 3000
)

Write-Host "🧪 Testando Next.js Web App na porta $Port" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:$Port"

# Testar health endpoint
Write-Host "1️⃣  Testando health endpoint..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-WebRequest -Uri "$baseUrl/api/health" -TimeoutSec 5 -ErrorAction Stop
    $healthData = $healthResponse.Content | ConvertFrom-Json
    Write-Host "   ✅ Health OK - Status: $($healthData.status)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Falha no health check: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Testar página principal
Write-Host "2️⃣  Testando página principal (Dashboard)..." -ForegroundColor Yellow
try {
    $homeResponse = Invoke-WebRequest -Uri "$baseUrl/" -TimeoutSec 10 -ErrorAction Stop
    if ($homeResponse.StatusCode -eq 200) {
        Write-Host "   ✅ Dashboard OK - Status Code: $($homeResponse.StatusCode)" -ForegroundColor Green
    }
} catch {
    Write-Host "   ❌ Falha ao carregar Dashboard: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Testar página de rotas
Write-Host "3️⃣  Testando página de rotas..." -ForegroundColor Yellow
try {
    $routesResponse = Invoke-WebRequest -Uri "$baseUrl/routes" -TimeoutSec 10 -ErrorAction Stop
    if ($routesResponse.StatusCode -eq 200) {
        Write-Host "   ✅ Rotas OK - Status Code: $($routesResponse.StatusCode)" -ForegroundColor Green
    }
} catch {
    Write-Host "   ❌ Falha ao carregar Rotas: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "📊 Resumo dos Testes" -ForegroundColor Cyan
Write-Host "   Base URL: $baseUrl" -ForegroundColor White
Write-Host ""
Write-Host "💡 Dica: Abra o navegador em:" -ForegroundColor Yellow
Write-Host "   $baseUrl" -ForegroundColor White
