# Configure hosts file for web.bgc.local
# ⚠️  Execute como Administrador!

$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
$entry = "127.0.0.1 web.bgc.local"

# Verificar se já existe
$content = Get-Content $hostsPath -ErrorAction SilentlyContinue
if ($content -match "web.bgc.local") {
    Write-Host "✓ web.bgc.local já está configurado no hosts file" -ForegroundColor Green
} else {
    try {
        Add-Content -Path $hostsPath -Value "`n$entry"
        Write-Host "✓ web.bgc.local adicionado ao hosts file" -ForegroundColor Green
    } catch {
        Write-Host "❌ Erro: Execute este script como Administrador!" -ForegroundColor Red
        Write-Host "   Clique com botão direito no PowerShell > Executar como Administrador" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host "`nAcesse: http://web.bgc.local" -ForegroundColor Cyan
Write-Host "Routes: http://web.bgc.local/routes.html" -ForegroundColor Cyan
