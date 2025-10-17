#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Restore PostgreSQL backup from Kubernetes PVC

.DESCRIPTION
    Lista backups disponíveis e restaura um backup específico do PostgreSQL

.PARAMETER BackupFile
    Nome do arquivo de backup (opcional, se não fornecido lista backups)

.EXAMPLE
    .\restore-backup.ps1
    # Lista backups disponíveis

.EXAMPLE
    .\restore-backup.ps1 -BackupFile bgc_backup_20250116_020000.sql.gz
    # Restaura backup específico
#>

param(
    [string]$BackupFile = ""
)

$ErrorActionPreference = "Stop"

Write-Host "=== PostgreSQL Backup Restore ===" -ForegroundColor Cyan

# Lista backups disponíveis
function List-Backups {
    Write-Host "`nListando backups disponíveis..." -ForegroundColor Yellow

    $pod = kubectl get pods -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}' 2>$null

    if (-not $pod) {
        Write-Host "ERRO: Pod do Postgres não encontrado!" -ForegroundColor Red
        exit 1
    }

    Write-Host "Backups no volume:" -ForegroundColor Green
    kubectl exec -n data $pod -- sh -c "ls -lh /backups/*.sql.gz 2>/dev/null || echo 'Nenhum backup encontrado'"
}

# Restaura backup
function Restore-Backup {
    param([string]$File)

    Write-Host "`nRestaurando backup: $File" -ForegroundColor Yellow

    $pod = kubectl get pods -n data -l app=postgres -o jsonpath='{.items[0].metadata.name}' 2>$null

    if (-not $pod) {
        Write-Host "ERRO: Pod do Postgres não encontrado!" -ForegroundColor Red
        exit 1
    }

    # Confirmação
    Write-Host "`n⚠️  ATENÇÃO: Esta operação vai SOBRESCREVER o banco atual!" -ForegroundColor Red
    $confirm = Read-Host "Deseja continuar? (digite 'SIM' para confirmar)"

    if ($confirm -ne "SIM") {
        Write-Host "Operação cancelada." -ForegroundColor Yellow
        exit 0
    }

    Write-Host "`nIniciando restore..." -ForegroundColor Cyan

    # Restaura backup
    kubectl exec -n data $pod -- sh -c @"
set -e
BACKUP_FILE=/backups/$File

if [ ! -f \$BACKUP_FILE ]; then
    echo 'ERRO: Arquivo de backup não encontrado!'
    exit 1
fi

echo 'Restaurando backup...'
gunzip -c \$BACKUP_FILE | psql -U bgc -d bgc

echo 'Restore concluído com sucesso!'
"@

    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n✅ Backup restaurado com sucesso!" -ForegroundColor Green
    } else {
        Write-Host "`n❌ Erro ao restaurar backup!" -ForegroundColor Red
        exit 1
    }
}

# Execução principal
if ($BackupFile -eq "") {
    List-Backups
    Write-Host "`nPara restaurar um backup, execute:" -ForegroundColor Cyan
    Write-Host "  .\restore-backup.ps1 -BackupFile <nome_do_arquivo>" -ForegroundColor White
} else {
    Restore-Backup -File $BackupFile
}
