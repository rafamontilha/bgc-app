#!/usr/bin/env pwsh
# BGC App - Kubernetes Management
# Usage: .\scripts\k8s.ps1 [setup|up|down|restart|logs|status|build|open|clean|help]

param(
    [Parameter(Position=0)]
    [ValidateSet("setup", "up", "down", "restart", "logs", "status", "build", "open", "clean", "help")]
    [string]$Command = "help"
)

$ClusterName = "bgc"
$Namespace = "data"

function Show-Help {
    Write-Host @"
BGC App - Kubernetes Management

USAGE:
    .\scripts\k8s.ps1 <command>

COMMANDS:
    setup       Create k3d cluster and deploy all services
    up          Deploy services to existing cluster
    down        Delete all deployments (keeps cluster)
    restart     Restart all pods
    logs        Show logs from all pods
    status      Show cluster and pods status
    build       Rebuild images and redeploy
    open        Configure hosts and open web.bgc.local
    clean       Delete entire cluster (⚠️  removes everything)
    help        Show this help

EXAMPLES:
    .\scripts\k8s.ps1 setup    # First time setup
    .\scripts\k8s.ps1 status   # Check status
    .\scripts\k8s.ps1 logs     # View logs
    .\scripts\k8s.ps1 open     # Open in browser

URLS:
    Web UI:  http://web.bgc.local
    Routes:  http://web.bgc.local/routes
    API:     http://api.bgc.local/healthz

NOTE: Run 'setup-hosts.ps1' as Administrator to configure web.bgc.local

"@ -ForegroundColor Cyan
}

function Test-ClusterExists {
    $clusters = k3d cluster list 2>$null | Select-String -Pattern "^$ClusterName\s"
    return $null -ne $clusters
}

function Build-Images {
    Write-Host "🔨 Building Docker images..." -ForegroundColor Yellow
    docker build -f api/Dockerfile -t bgc/bgc-api:dev .
    docker build -t bgc/bgc-web:latest web-next/
    Write-Host "✅ Images built!" -ForegroundColor Green
}

function Import-Images {
    Write-Host "📦 Importing images to k3d..." -ForegroundColor Yellow
    k3d image import bgc/bgc-api:dev bgc/bgc-web:latest -c $ClusterName
    Write-Host "✅ Images imported!" -ForegroundColor Green
}

function Deploy-Services {
    Write-Host "🚀 Deploying services..." -ForegroundColor Green

    # PostgreSQL
    kubectl apply -f deploy/postgres.yaml
    Write-Host "⏳ Waiting for PostgreSQL..." -ForegroundColor Yellow
    kubectl wait --for=condition=ready pod -l app=postgres -n $Namespace --timeout=120s

    # PostgreSQL Backup CronJob
    kubectl apply -f k8s/postgres-backup-cronjob.yaml
    Write-Host "✅ Backup CronJob configured" -ForegroundColor Green

    # Materialized View Refresh CronJob
    kubectl apply -f k8s/mview-refresh-cronjob.yaml
    Write-Host "✅ MView Refresh CronJob configured" -ForegroundColor Green

    # API
    kubectl apply -f k8s/api.yaml
    Write-Host "⏳ Waiting for API..." -ForegroundColor Yellow
    kubectl wait --for=condition=ready pod -l app=bgc-api -n $Namespace --timeout=120s

    # API HPA
    kubectl apply -f k8s/api-hpa.yaml
    Write-Host "✅ API HPA configured" -ForegroundColor Green

    # Web
    kubectl apply -f k8s/web.yaml
    Write-Host "⏳ Waiting for Web..." -ForegroundColor Yellow
    kubectl wait --for=condition=ready pod -l app=bgc-web -n $Namespace --timeout=60s

    # Web HPA
    kubectl apply -f k8s/web-hpa.yaml
    Write-Host "✅ Web HPA configured" -ForegroundColor Green

    Write-Host "`n✅ All services deployed!" -ForegroundColor Green
}

switch ($Command.ToLower()) {
    "setup" {
        Write-Host "🎯 Setting up Kubernetes cluster..." -ForegroundColor Cyan

        if (Test-ClusterExists) {
            Write-Host "⚠️  Cluster '$ClusterName' already exists. Delete it first with 'clean' command." -ForegroundColor Yellow
            break
        }

        # Create cluster
        Write-Host "🌐 Creating k3d cluster..." -ForegroundColor Green
        k3d cluster create $ClusterName --port "80:80@loadbalancer" --port "443:443@loadbalancer"

        # Create namespace
        Write-Host "📁 Creating namespace..." -ForegroundColor Green
        kubectl create namespace $Namespace

        # Build and import images
        Build-Images
        Import-Images

        # Deploy services
        Deploy-Services

        Write-Host "`n✅ Setup complete!" -ForegroundColor Green
        Write-Host "`n⚠️  IMPORTANT: Run as Administrator:" -ForegroundColor Yellow
        Write-Host "    .\scripts\setup-hosts.ps1" -ForegroundColor Cyan
        Write-Host "`nThen access: http://web.bgc.local" -ForegroundColor Cyan
        break
    }

    "up" {
        if (-not (Test-ClusterExists)) {
            Write-Host "❌ Cluster '$ClusterName' not found. Run 'setup' first." -ForegroundColor Red
            break
        }

        Build-Images
        Import-Images
        Deploy-Services
        break
    }

    "down" {
        Write-Host "🛑 Deleting deployments..." -ForegroundColor Yellow
        kubectl delete -f k8s/web-hpa.yaml --ignore-not-found
        kubectl delete -f k8s/web.yaml --ignore-not-found
        kubectl delete -f k8s/api-hpa.yaml --ignore-not-found
        kubectl delete -f k8s/api.yaml --ignore-not-found
        kubectl delete -f k8s/mview-refresh-cronjob.yaml --ignore-not-found
        kubectl delete -f k8s/postgres-backup-cronjob.yaml --ignore-not-found
        kubectl delete -f deploy/postgres.yaml --ignore-not-found
        Write-Host "✅ Deployments deleted!" -ForegroundColor Green
        break
    }

    "restart" {
        Write-Host "🔄 Restarting pods..." -ForegroundColor Yellow
        kubectl rollout restart deployment/bgc-web -n $Namespace
        kubectl rollout restart deployment/bgc-api -n $Namespace
        kubectl rollout restart deployment/postgres -n $Namespace
        Write-Host "✅ Pods restarted!" -ForegroundColor Green
        break
    }

    "logs" {
        Write-Host "📋 Choose a service:" -ForegroundColor Cyan
        Write-Host "1. API" -ForegroundColor White
        Write-Host "2. Web" -ForegroundColor White
        Write-Host "3. PostgreSQL" -ForegroundColor White
        Write-Host "4. All" -ForegroundColor White
        $choice = Read-Host "Enter number (1-4)"

        switch ($choice) {
            "1" { kubectl logs -l app=bgc-api -n $Namespace --tail=100 -f }
            "2" { kubectl logs -l app=bgc-web -n $Namespace --tail=100 -f }
            "3" { kubectl logs -l app=postgres -n $Namespace --tail=100 -f }
            "4" { kubectl logs -l app -n $Namespace --tail=50 -f }
            default { Write-Host "Invalid choice" -ForegroundColor Red }
        }
        break
    }

    "status" {
        Write-Host "📊 Cluster Status:" -ForegroundColor Cyan
        kubectl cluster-info
        Write-Host "`n📦 Pods in namespace '$Namespace':" -ForegroundColor Cyan
        kubectl get pods -n $Namespace -o wide
        Write-Host "`n🌐 Services:" -ForegroundColor Cyan
        kubectl get svc -n $Namespace
        Write-Host "`n🔗 Ingress:" -ForegroundColor Cyan
        kubectl get ingress -n $Namespace
        Write-Host "`n📈 HPA (Horizontal Pod Autoscaler):" -ForegroundColor Cyan
        kubectl get hpa -n $Namespace
        Write-Host "`n⏰ CronJobs:" -ForegroundColor Cyan
        kubectl get cronjobs -n $Namespace
        break
    }

    "build" {
        Write-Host "🔨 Rebuilding and redeploying..." -ForegroundColor Yellow
        Build-Images
        Import-Images
        kubectl rollout restart deployment/bgc-api -n $Namespace
        kubectl rollout restart deployment/bgc-web -n $Namespace
        Write-Host "✅ Rebuild complete!" -ForegroundColor Green
        break
    }

    "open" {
        Write-Host "🌐 Configuring hosts and opening browser..." -ForegroundColor Cyan

        # Check if hosts is configured
        $hostsPath = "C:\Windows\System32\drivers\etc\hosts"
        $hostsContent = Get-Content $hostsPath -ErrorAction SilentlyContinue

        if ($hostsContent -notmatch "web.bgc.local") {
            Write-Host "⚠️  web.bgc.local not configured in hosts file!" -ForegroundColor Yellow
            Write-Host "Run as Administrator:" -ForegroundColor Yellow
            Write-Host "    .\scripts\setup-hosts.ps1" -ForegroundColor Cyan
        } else {
            Start-Process "http://web.bgc.local"
            Start-Process "http://web.bgc.local/routes"
            Write-Host "✅ Browser opened!" -ForegroundColor Green
        }
        break
    }

    "clean" {
        Write-Host "⚠️  This will DELETE the entire cluster!" -ForegroundColor Red
        $confirm = Read-Host "Continue? (yes/no)"
        if ($confirm -eq "yes") {
            k3d cluster delete $ClusterName
            Write-Host "✅ Cluster deleted!" -ForegroundColor Green
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
