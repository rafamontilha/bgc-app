# Script para instalar Go 1.23
# Executa: .\scripts\install-go.ps1

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " Instalador Go 1.23 para Windows" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

$goVersion = "1.23.2"
$goArch = "windows-amd64"
$goFileName = "go$goVersion.$goArch.msi"
$goUrl = "https://go.dev/dl/$goFileName"
$downloadPath = "$env:TEMP\$goFileName"

Write-Host "[1/4] Verificando se Go ja esta instalado..." -ForegroundColor Yellow

try {
    $currentGo = & go version 2>$null
    if ($currentGo) {
        Write-Host "      Go ja esta instalado: $currentGo" -ForegroundColor Green
        Write-Host ""
        $continue = Read-Host "Deseja reinstalar? (s/N)"
        if ($continue -ne "s" -and $continue -ne "S") {
            Write-Host "Instalacao cancelada." -ForegroundColor Yellow
            exit 0
        }
    }
} catch {
    Write-Host "      Go nao encontrado. Prosseguindo com instalacao..." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "[2/4] Baixando Go $goVersion..." -ForegroundColor Yellow
Write-Host "      URL: $goUrl" -ForegroundColor Gray

try {
    Invoke-WebRequest -Uri $goUrl -OutFile $downloadPath -UseBasicParsing
    Write-Host "      Download concluido!" -ForegroundColor Green
} catch {
    Write-Host "      Erro ao baixar Go: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "[3/4] Instalando Go..." -ForegroundColor Yellow
Write-Host "      Isso pode demorar alguns minutos..." -ForegroundColor Gray
Write-Host "      Uma janela de instalacao sera aberta." -ForegroundColor Gray

try {
    Start-Process -FilePath "msiexec.exe" -ArgumentList "/i `"$downloadPath`" /quiet /norestart" -Wait
    Write-Host "      Instalacao concluida!" -ForegroundColor Green
} catch {
    Write-Host "      Erro ao instalar Go: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "[4/4] Verificando instalacao..." -ForegroundColor Yellow

# Atualizar PATH para a sessao atual
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

try {
    $installedGo = & go version 2>$null
    if ($installedGo) {
        Write-Host "      $installedGo" -ForegroundColor Green
        Write-Host ""
        Write-Host "=====================================" -ForegroundColor Green
        Write-Host " Go instalado com sucesso!" -ForegroundColor Green
        Write-Host "=====================================" -ForegroundColor Green
        Write-Host ""
        Write-Host "IMPORTANTE: Feche e reabra o terminal para usar o Go." -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Proximos passos:" -ForegroundColor Cyan
        Write-Host "  1. Feche este terminal" -ForegroundColor White
        Write-Host "  2. Abra um NOVO terminal" -ForegroundColor White
        Write-Host "  3. Execute: cd api && go run cmd/api/main.go" -ForegroundColor White
    } else {
        Write-Host "      Go foi instalado mas nao esta no PATH." -ForegroundColor Yellow
        Write-Host "      Feche e reabra o terminal." -ForegroundColor Yellow
    }
} catch {
    Write-Host "      Go foi instalado mas nao esta no PATH." -ForegroundColor Yellow
    Write-Host "      Feche e reabra o terminal." -ForegroundColor Yellow
}

Write-Host ""

# Limpar arquivo temporario
Remove-Item -Path $downloadPath -ErrorAction SilentlyContinue
