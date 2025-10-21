# Script de Verificacao do Ambiente de Desenvolvimento BGC
# Verifica todas as ferramentas necessarias para o desenvolvimento

Write-Host "`n BGC - Verificacao do Ambiente de Desenvolvimento" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Gray
Write-Host ""

$allOk = $true

# Funcao auxiliar para verificar comandos
function Test-Command {
    param(
        [string]$Name,
        [string]$Command
    )

    Write-Host "Verificando $Name..." -NoNewline -ForegroundColor Yellow

    try {
        $output = Invoke-Expression $Command 2>&1

        if ($LASTEXITCODE -eq 0 -or $output) {
            Write-Host " OK" -ForegroundColor Green
            Write-Host "  Versao: $output" -ForegroundColor Gray
            return $true
        } else {
            throw "Comando falhou"
        }
    } catch {
        Write-Host " NAO INSTALADO" -ForegroundColor Red
        return $false
    }
}

# Verificar Node.js
Write-Host "`nRuntime JavaScript" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor Gray

if (-not (Test-Command -Name "Node.js" -Command "node --version")) {
    $allOk = $false
    Write-Host "  Instale com: winget install OpenJS.NodeJS.LTS" -ForegroundColor Yellow
}

if (-not (Test-Command -Name "npm" -Command "npm --version")) {
    $allOk = $false
    Write-Host "  npm vem com Node.js" -ForegroundColor Yellow
}

if (-not (Test-Command -Name "pnpm" -Command "pnpm --version")) {
    Write-Host "  Instale com: npm install -g pnpm" -ForegroundColor Yellow
    Write-Host "  pnpm e opcional, mas recomendado" -ForegroundColor Blue
}

# Verificar Git
Write-Host "`nControle de Versao" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor Gray

if (-not (Test-Command -Name "Git" -Command "git --version")) {
    $allOk = $false
    Write-Host "  Instale com: winget install Git.Git" -ForegroundColor Yellow
}

# Verificar Docker
Write-Host "`nContainerizacao" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor Gray

if (-not (Test-Command -Name "Docker" -Command "docker --version")) {
    $allOk = $false
    Write-Host "  Instale Docker Desktop: https://docker.com/products/docker-desktop" -ForegroundColor Yellow
}

# Verificar Docker Compose
if (-not (Test-Command -Name "Docker Compose" -Command "docker compose version")) {
    Write-Host "  Docker Compose vem com Docker Desktop" -ForegroundColor Blue
}

# Verificar Go
Write-Host "`nLinguagem Backend" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor Gray

if (-not (Test-Command -Name "Go" -Command "go version")) {
    $allOk = $false
    Write-Host "  Instale com: .\scripts\install-go.ps1" -ForegroundColor Yellow
}

# Verificar PowerShell
Write-Host "`nShell" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor Gray

$psVersion = $PSVersionTable.PSVersion.ToString()
Write-Host "Verificando PowerShell..." -NoNewline -ForegroundColor Yellow
Write-Host " OK" -ForegroundColor Green
Write-Host "  Versao: $psVersion" -ForegroundColor Gray

# Verificar VS Code (opcional)
Write-Host "`nEditor (Opcional)" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------" -ForegroundColor Gray

if (-not (Test-Command -Name "VS Code" -Command "code --version")) {
    Write-Host "  Recomendado: winget install Microsoft.VisualStudioCode" -ForegroundColor Blue
}

# Resumo
Write-Host "`n============================================================" -ForegroundColor Gray
Write-Host ""

if ($allOk) {
    Write-Host "Ambiente configurado corretamente!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Proximos passos:" -ForegroundColor Cyan
    Write-Host "  1. Feche e reabra o terminal para carregar as variaveis de ambiente" -ForegroundColor Gray
    Write-Host "  2. Navegue ate o diretorio do projeto" -ForegroundColor Gray
    Write-Host "  3. Execute: pnpm create next-app@latest web-next" -ForegroundColor Gray
} else {
    Write-Host "Algumas ferramentas obrigatorias estao faltando" -ForegroundColor Yellow
    Write-Host "Por favor, instale as ferramentas marcadas com 'NAO INSTALADO'" -ForegroundColor Gray
}

Write-Host ""
Write-Host "Versoes recomendadas:" -ForegroundColor Blue
Write-Host "  Node.js: v20.x ou v22.x (LTS)" -ForegroundColor Gray
Write-Host "  npm: v10.x+" -ForegroundColor Gray
Write-Host "  pnpm: v9.x+" -ForegroundColor Gray
Write-Host "  Go: v1.23+" -ForegroundColor Gray
Write-Host "  Git: v2.40+" -ForegroundColor Gray
Write-Host "  Docker: v24.0+" -ForegroundColor Gray
Write-Host ""
