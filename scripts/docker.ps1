#!/usr/bin/env pwsh
# BGC App - Docker Compose Management
# Usage: .\scripts\docker.ps1 [up|down|restart|logs|ps|help]

param(
    [Parameter(Position=0)]
    [ValidateSet("up", "down", "restart", "logs", "ps", "build", "clean", "help")]
    [string]$Command = "help"
)

$ComposeFile = "bgcstack/docker-compose.yml"

function Show-Help {
    Write-Host @"
BGC App - Docker Compose Management

USAGE:
    .\scripts\docker.ps1 <command>

COMMANDS:
    up          Start all services (db, api, web, pgadmin)
    down        Stop all services
    restart     Restart all services
    logs        Show logs from all services
    ps          Show running containers
    build       Rebuild images and start
    clean       Stop and remove volumes (⚠️  removes database data)
    help        Show this help

EXAMPLES:
    .\scripts\docker.ps1 up
    .\scripts\docker.ps1 logs
    .\scripts\docker.ps1 down

URLS:
    API:     http://localhost:8080
    Web UI:  http://localhost:3000
    PgAdmin: http://localhost:5050 (admin@bgc.dev / admin)

"@ -ForegroundColor Cyan
}

switch ($Command.ToLower()) {
    "up" {
        Write-Host "🚀 Starting Docker Compose..." -ForegroundColor Green
        docker compose -f $ComposeFile up -d
        Start-Sleep -Seconds 3
        Write-Host "`n✅ Services started!" -ForegroundColor Green
        Write-Host "API:     http://localhost:8080/healthz" -ForegroundColor Cyan
        Write-Host "Web UI:  http://localhost:3000" -ForegroundColor Cyan
        Write-Host "PgAdmin: http://localhost:5050" -ForegroundColor Cyan
        break
    }

    "down" {
        Write-Host "🛑 Stopping Docker Compose..." -ForegroundColor Yellow
        docker compose -f $ComposeFile down
        Write-Host "✅ Services stopped!" -ForegroundColor Green
        break
    }

    "restart" {
        Write-Host "🔄 Restarting Docker Compose..." -ForegroundColor Yellow
        docker compose -f $ComposeFile restart
        Write-Host "✅ Services restarted!" -ForegroundColor Green
        break
    }

    "logs" {
        Write-Host "📋 Showing logs (Ctrl+C to exit)..." -ForegroundColor Cyan
        docker compose -f $ComposeFile logs -f
        break
    }

    "ps" {
        Write-Host "📦 Running containers:" -ForegroundColor Cyan
        docker compose -f $ComposeFile ps
        break
    }

    "build" {
        Write-Host "🔨 Rebuilding images and starting..." -ForegroundColor Yellow
        docker compose -f $ComposeFile up -d --build
        Start-Sleep -Seconds 3
        Write-Host "`n✅ Services rebuilt and started!" -ForegroundColor Green
        Write-Host "API:     http://localhost:8080/healthz" -ForegroundColor Cyan
        Write-Host "Web UI:  http://localhost:3000" -ForegroundColor Cyan
        break
    }

    "clean" {
        Write-Host "⚠️  This will stop services and remove volumes (database data will be lost)!" -ForegroundColor Red
        $confirm = Read-Host "Continue? (yes/no)"
        if ($confirm -eq "yes") {
            docker compose -f $ComposeFile down -v
            Write-Host "✅ Services stopped and volumes removed!" -ForegroundColor Green
        } else {
            Write-Host "❌ Cancelled." -ForegroundColor Yellow
        }
        break
    }

    "help" {
        Show-Help
        break
    }

    default {
        Show-Help
        break
    }
}
