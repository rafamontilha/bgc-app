# Script para adicionar entradas BGC ao arquivo hosts
# Deve ser executado como Administrador

$hostsFile = "C:\Windows\System32\drivers\etc\hosts"
$bgcEntries = @"

# BGC Application URLs
127.0.0.1  api.bgc.local
127.0.0.1  web.bgc.local
"@

Write-Host "Verificando se as entradas BGC já existem no arquivo hosts..." -ForegroundColor Cyan

$hostsContent = Get-Content $hostsFile -Raw

if ($hostsContent -notmatch "api.bgc.local") {
    Write-Host "Adicionando entradas BGC ao arquivo hosts..." -ForegroundColor Yellow

    # Verificar se temos permissão de administrador
    $isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

    if (-not $isAdmin) {
        Write-Host "ERRO: Este script precisa ser executado como Administrador!" -ForegroundColor Red
        Write-Host "Clique com botão direito no PowerShell e selecione 'Executar como Administrador'" -ForegroundColor Yellow
        exit 1
    }

    Add-Content -Path $hostsFile -Value $bgcEntries
    Write-Host "Entradas adicionadas com sucesso!" -ForegroundColor Green
} else {
    Write-Host "As entradas BGC já existem no arquivo hosts." -ForegroundColor Green
}

Write-Host "`nVocê pode agora acessar as aplicações:" -ForegroundColor Cyan
Write-Host "  - Dashboard:   http://web.bgc.local" -ForegroundColor White
Write-Host "  - Rotas:       http://web.bgc.local/routes" -ForegroundColor White
Write-Host "  - API:         http://api.bgc.local/healthz" -ForegroundColor White
