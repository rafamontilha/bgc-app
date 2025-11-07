# Test Script for Integration Gateway (Local)

param(
    [string]$Connector = "viacep",
    [switch]$SkipBuild,
    [switch]$SkipTests,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " Integration Gateway - Local Tests  " -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Cores
$Green = "Green"
$Red = "Red"
$Yellow = "Yellow"
$Cyan = "Cyan"

# Variáveis
$RootDir = Split-Path -Parent $PSScriptRoot
$ServiceDir = Join-Path $RootDir "services\integration-gateway"
$ConfigDir = Join-Path $RootDir "config\connectors"
$CertsDir = Join-Path $RootDir "certs"

# Função para log
function Write-Step {
    param([string]$Message)
    Write-Host "`n>> $Message" -ForegroundColor $Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor $Green
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor $Red
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ $Message" -ForegroundColor $Yellow
}

# 1. Validar ambiente
Write-Step "Validando ambiente..."

if (-not (Test-Path $ServiceDir)) {
    Write-Error-Custom "Service directory not found: $ServiceDir"
    exit 1
}

if (-not (Test-Path $ConfigDir)) {
    Write-Error-Custom "Config directory not found: $ConfigDir"
    exit 1
}

Write-Success "Ambiente validado"

# 2. Executar testes unitários
if (-not $SkipTests) {
    Write-Step "Executando testes unitários..."

    Push-Location $ServiceDir
    try {
        Write-Info "Running: go test ./... -v -cover"
        go test ./... -v -cover

        if ($LASTEXITCODE -ne 0) {
            Write-Error-Custom "Unit tests failed"
            Pop-Location
            exit 1
        }

        Write-Success "Testes unitários passaram"
    }
    finally {
        Pop-Location
    }
}
else {
    Write-Info "Skipping unit tests (--SkipTests)"
}

# 3. Build (opcional)
if (-not $SkipBuild) {
    Write-Step "Building application..."

    Push-Location $ServiceDir
    try {
        go build -o integration-gateway.exe cmd/gateway/main.go

        if ($LASTEXITCODE -ne 0) {
            Write-Error-Custom "Build failed"
            Pop-Location
            exit 1
        }

        Write-Success "Build concluído"
    }
    finally {
        Pop-Location
    }
}
else {
    Write-Info "Skipping build (--SkipBuild)"
}

# 4. Iniciar gateway em background
Write-Step "Iniciando Integration Gateway..."

$env:CONFIG_DIR = $ConfigDir
$env:CERTS_DIR = $CertsDir
$env:ENVIRONMENT = "development"
$env:LOG_LEVEL = if ($Verbose) { "debug" } else { "info" }
$env:PORT = "8081"

Push-Location $ServiceDir
try {
    # Inicia processo em background
    $GatewayProcess = Start-Process -FilePath "go" `
        -ArgumentList "run", "cmd/gateway/main.go" `
        -PassThru `
        -NoNewWindow `
        -RedirectStandardOutput "$env:TEMP\integration-gateway-out.log" `
        -RedirectStandardError "$env:TEMP\integration-gateway-err.log"

    Write-Info "Gateway PID: $($GatewayProcess.Id)"
    Write-Info "Aguardando inicialização..."
    Start-Sleep -Seconds 5

    # Verifica se ainda está rodando
    if ($GatewayProcess.HasExited) {
        Write-Error-Custom "Gateway failed to start"
        Get-Content "$env:TEMP\integration-gateway-err.log" | Write-Host
        exit 1
    }

    Write-Success "Gateway iniciado"
}
catch {
    Write-Error-Custom "Failed to start gateway: $_"
    Pop-Location
    exit 1
}
finally {
    Pop-Location
}

# 5. Executar smoke tests
Write-Step "Executando smoke tests..."

try {
    # Test 1: Health check
    Write-Info "Test 1: Health check"
    $health = Invoke-RestMethod -Uri "http://localhost:8081/health" -Method Get
    if ($health.status -eq "healthy") {
        Write-Success "Health check OK (connectors: $($health.connectors))"
    }
    else {
        Write-Error-Custom "Health check failed"
        throw "Health check returned: $($health.status)"
    }

    # Test 2: List connectors
    Write-Info "Test 2: List connectors"
    $connectors = Invoke-RestMethod -Uri "http://localhost:8081/v1/connectors" -Method Get
    Write-Success "Connectors loaded: $($connectors.Count)"

    foreach ($conn in $connectors) {
        Write-Host "  - $($conn.id): $($conn.name) v$($conn.version)" -ForegroundColor Gray
    }

    # Test 3: Test connector execution (ViaCEP)
    if ($Connector -eq "viacep") {
        Write-Info "Test 3: Execute ViaCEP connector"

        $body = @{ cep = "01310100" } | ConvertTo-Json
        $result = Invoke-RestMethod `
            -Uri "http://localhost:8081/v1/connectors/viacep/consulta_cep" `
            -Method Post `
            -Body $body `
            -ContentType "application/json"

        if ($result.data.cep -eq "01310-100") {
            Write-Success "ViaCEP test OK"
            Write-Host "  CEP: $($result.data.cep)" -ForegroundColor Gray
            Write-Host "  Logradouro: $($result.data.logradouro)" -ForegroundColor Gray
            Write-Host "  Bairro: $($result.data.bairro)" -ForegroundColor Gray
            Write-Host "  Cidade: $($result.data.localidade)/$($result.data.uf)" -ForegroundColor Gray
            Write-Host "  Duration: $($result.duration)" -ForegroundColor Gray
        }
        else {
            Write-Error-Custom "ViaCEP test failed: unexpected response"
            $result | ConvertTo-Json | Write-Host
        }
    }

    # Test 4: Metrics endpoint
    Write-Info "Test 4: Metrics endpoint"
    $metrics = Invoke-WebRequest -Uri "http://localhost:8081/metrics" -Method Get
    if ($metrics.StatusCode -eq 200) {
        $metricsText = $metrics.Content
        $metricsCount = ($metricsText -split "`n" | Where-Object { $_ -match "^bgc_connector_" }).Count
        Write-Success "Metrics OK ($metricsCount metrics found)"
    }
    else {
        Write-Error-Custom "Metrics endpoint failed"
    }

}
catch {
    Write-Error-Custom "Smoke tests failed: $_"
    $Failed = $true
}
finally {
    # 6. Cleanup - Stop gateway
    Write-Step "Stopping gateway..."

    if ($GatewayProcess -and -not $GatewayProcess.HasExited) {
        Stop-Process -Id $GatewayProcess.Id -Force
        Write-Success "Gateway stopped"
    }

    # Exibir logs se houver erro
    if ($Failed -or $Verbose) {
        Write-Step "Gateway Logs:"
        Write-Host "=== STDOUT ===" -ForegroundColor Yellow
        Get-Content "$env:TEMP\integration-gateway-out.log" -ErrorAction SilentlyContinue | Write-Host
        Write-Host "`n=== STDERR ===" -ForegroundColor Yellow
        Get-Content "$env:TEMP\integration-gateway-err.log" -ErrorAction SilentlyContinue | Write-Host
    }
}

# Resultado final
Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan

if ($Failed) {
    Write-Host " TESTS FAILED ✗" -ForegroundColor $Red
    Write-Host "=====================================" -ForegroundColor Cyan
    exit 1
}
else {
    Write-Host " ALL TESTS PASSED ✓" -ForegroundColor $Green
    Write-Host "=====================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Integration Gateway is ready for deployment!" -ForegroundColor $Green
    exit 0
}
