# BGC K8s Setup Script v2 - Com tratamento robusto de erros
param(
    [switch]$SkipCleanup,
    [switch]$ForceCleanup
)

$ErrorActionPreference = "Continue"  # Não para no primeiro erro

function Write-Step {
    param([string]$Message)
    Write-Host "`n $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host " $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host " $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "  $Message" -ForegroundColor Yellow
}

try {
    Write-Host "`n" -ForegroundColor Green
    Write-Host "   BGC Analytics - K8s Setup Script      " -ForegroundColor Green
    Write-Host "`n" -ForegroundColor Green

    # ===== LIMPEZA =====
    if (-not $SkipCleanup) {
        Write-Step "Limpando ambiente anterior..."
        
        # Docker Compose - com timeout e force
        try {
            Write-Host "   Parando Docker Compose (timeout 30s)..." -ForegroundColor Gray
            $composeJob = Start-Job -ScriptBlock {
                Set-Location $using:PWD
                docker compose -f bgcstack/docker-compose.yml down --timeout 10 2>&1
            }
            
            $composeJob | Wait-Job -Timeout 30 | Out-Null
            
            if ($composeJob.State -eq "Running") {
                Write-Warning "Docker Compose travou - forçando parada..."
                Stop-Job $composeJob
                
                # Force kill containers
                docker ps -q --filter "name=bgc_" | ForEach-Object {
                    docker stop $_ --time 5 2>$null
                    docker rm $_ -f 2>$null
                }
            }
            
            Remove-Job $composeJob -Force 2>$null
            Write-Success "Docker Compose parado"
        } catch {
            Write-Warning "Erro ao parar Docker Compose (continuando...): $_"
        }
        
        # K3d cluster
        try {
            Write-Host "   Deletando cluster k3d..." -ForegroundColor Gray
            k3d cluster delete bgc 2>$null
            Start-Sleep -Seconds 2
            Write-Success "Cluster k3d deletado"
        } catch {
            Write-Warning "Erro ao deletar cluster (continuando...): $_"
        }
        
        # Liberar porta 8080
        try {
            Write-Host "   Verificando porta 8080..." -ForegroundColor Gray
            $processId = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue | 
                         Select-Object -ExpandProperty OwningProcess -Unique
            if ($processId) {
                Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue
                Write-Success "Porta 8080 liberada"
            } else {
                Write-Host "   Porta 8080 já está livre" -ForegroundColor Gray
            }
        } catch {
            Write-Host "   Porta 8080 já está livre" -ForegroundColor Gray
        }
        
        Write-Success "Limpeza concluída"
    }

    # ===== CLUSTER =====
    Write-Step "Criando cluster k3d..."
    k3d cluster create bgc --port "80:80@loadbalancer" --port "443:443@loadbalancer" 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) { 
        throw "Falha ao criar cluster k3d. Tente: k3d cluster delete bgc && k3d cluster create bgc" 
    }
    Start-Sleep -Seconds 5
    Write-Success "Cluster criado"

    Write-Step "Criando namespace 'data'..."
    kubectl create namespace data 2>$null
    kubectl config set-context --current --namespace=data 2>&1 | Out-Null
    Write-Success "Namespace configurado"

    # ===== POSTGRESQL =====
    Write-Step "Deploying PostgreSQL..."
    kubectl apply -f deploy/postgres.yaml 2>&1 | Out-Null
    
    Write-Host "    Aguardando PostgreSQL (pode levar 2-3 minutos)..." -ForegroundColor Yellow
    $waited = 0
    while ($waited -lt 300) {
        $podStatus = kubectl get pods -n data -l app=postgres -o jsonpath='{.items[0].status.phase}' 2>$null
        if ($podStatus -eq "Running") {
            break
        }
        Write-Host "     $waited s / 300s - Status: $podStatus" -ForegroundColor Gray
        Start-Sleep -Seconds 10
        $waited += 10
    }
    
    if ($waited -ge 300) { 
        Write-Error "PostgreSQL timeout"
        kubectl describe pod -n data -l app=postgres
        throw "PostgreSQL não ficou pronto em 5 minutos" 
    }
    Write-Success "PostgreSQL pronto"

    # Validar conectividade
    Write-Step "Validando PostgreSQL..."
    Start-Sleep -Seconds 5
    
    $testResult = kubectl run pg-test --rm -i --restart=Never --image=postgres:16 `
      --env="PGPASSWORD=bgcpassword" `
      -- psql -h pg-postgresql -U postgres -d postgres -c "SELECT 1;" 2>&1
    
    if ($testResult -match "1 row") {
        Write-Success "PostgreSQL validado"
    } else {
        Write-Warning "PostgreSQL pode não estar totalmente pronto. Tentando continuar..."
        Start-Sleep -Seconds 10
    }

    # ===== CONFIGMAPS =====
    Write-Step "Criando ConfigMaps..."
    kubectl apply -f deploy/configmap-migrate-0001.yaml 2>&1 | Out-Null
    kubectl apply -f deploy/configmap-migrate-0002.yaml 2>&1 | Out-Null
    kubectl apply -f deploy/configmap-views.yaml 2>&1 | Out-Null
    kubectl apply -f deploy/configmap-mviews-init.yaml 2>&1 | Out-Null
    kubectl apply -f deploy/configmap-mviews-add-uniq.yaml 2>&1 | Out-Null
    kubectl apply -f deploy/configmap-mviews-populate.yaml 2>&1 | Out-Null
    kubectl apply -f deploy/configmap-sample-csv.yaml 2>&1 | Out-Null
    Write-Success "ConfigMaps criados"

    # ===== MIGRATIONS =====
    Write-Step "Executando Migration 0001..."
    kubectl delete job bgc-migrate-0001 -n data 2>$null
    kubectl apply -f deploy/bgc-migrate-0001.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 5
    
    kubectl wait --for=condition=complete job/bgc-migrate-0001 -n data --timeout=120s 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) { 
        Write-Warning "Migration 0001 teve problemas. Verificando logs..."
        kubectl logs job/bgc-migrate-0001 -n data
        throw "Migration 0001 falhou" 
    }
    Write-Success "Migration 0001 completa"

    Write-Step "Executando Migration 0002..."
    kubectl delete job bgc-migrate-0002 -n data 2>$null
    kubectl apply -f deploy/bgc-migrate-0002.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 5
    
    kubectl wait --for=condition=complete job/bgc-migrate-0002 -n data --timeout=120s 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) { 
        Write-Warning "Migration 0002 teve problemas. Verificando logs..."
        kubectl logs job/bgc-migrate-0002 -n data
        throw "Migration 0002 falhou" 
    }
    Write-Success "Migration 0002 completa"

    Write-Step "Criando views..."
    kubectl delete job bgc-create-views -n data 2>$null
    kubectl apply -f deploy/bgc-create-views.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 5
    kubectl wait --for=condition=complete job/bgc-create-views -n data --timeout=120s 2>&1 | Out-Null
    Write-Success "Views criadas"

    Write-Step "Criando materialized views..."
    kubectl delete job bgc-mviews-init -n data 2>$null
    kubectl apply -f deploy/bgc-mviews-init.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 5
    kubectl wait --for=condition=complete job/bgc-mviews-init -n data --timeout=120s 2>&1 | Out-Null
    Write-Success "Materialized views criadas"

    Write-Step "Adicionando índices únicos..."
    kubectl delete job bgc-mviews-add-uniq -n data 2>$null
    kubectl apply -f deploy/bgc-mviews-add-uniq.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 5
    kubectl wait --for=condition=complete job/bgc-mviews-add-uniq -n data --timeout=120s 2>&1 | Out-Null
    Write-Success "Índices criados"

    # ===== BUILD & IMPORT =====
    Write-Step "Building imagens Docker..."
    Write-Host "   Building API..." -ForegroundColor Gray
    docker build -t bgc/bgc-api:dev api/ 2>&1 | Out-Null
    Write-Host "   Building Ingest..." -ForegroundColor Gray
    docker build -t bgc/bgc-ingest:dev services/bgc-ingest/ 2>&1 | Out-Null
    Write-Success "Imagens buildadas"

    Write-Step "Importando imagens no k3d..."
    k3d image import bgc/bgc-api:dev bgc/bgc-ingest:dev -c bgc 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) { throw "Falha ao importar imagens" }
    Write-Success "Imagens importadas"

    # ===== DEPLOY API =====
    Write-Step "Deploying API..."
    kubectl apply -f k8s/api.yaml 2>&1 | Out-Null
    
    Write-Host "    Aguardando API..." -ForegroundColor Yellow
    $waited = 0
    while ($waited -lt 120) {
        $podStatus = kubectl get pods -n data -l app=bgc-api -o jsonpath='{.items[0].status.phase}' 2>$null
        if ($podStatus -eq "Running") {
            break
        }
        Write-Host "     $waited s / 120s - Status: $podStatus" -ForegroundColor Gray
        Start-Sleep -Seconds 10
        $waited += 10
    }
    
    if ($waited -ge 120) { 
        Write-Error "API timeout"
        kubectl logs deployment/bgc-api -n data --tail=50
        throw "API não ficou pronta" 
    }
    Write-Success "API deployed"

    # ===== DEPLOY WEB =====
    Write-Step "Deploying Web UI..."
    kubectl apply -f k8s/web.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 10
    kubectl wait --for=condition=ready pod -l app=bgc-web -n data --timeout=60s 2>&1 | Out-Null
    Write-Success "Web UI deployed"

    # ===== LOAD DATA =====
    Write-Step "Carregando dados de exemplo..."
    kubectl delete job -n data -l app=bgc-ingest 2>$null
    kubectl apply -f deploy/bgc-ingest-load-csv.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 10
    
    $jobs = kubectl get jobs -n data -o name 2>$null | Select-String "load-csv"
    if ($jobs) {
        $latestJob = $jobs | Select-Object -First 1
        kubectl wait --for=condition=complete $latestJob -n data --timeout=120s 2>&1 | Out-Null
        Write-Success "Dados carregados"
    } else {
        Write-Warning "Job de carga não encontrado (continuando...)"
    }

    Write-Step "Populando materialized views..."
    kubectl delete job bgc-mviews-populate-initial -n data 2>$null
    kubectl apply -f deploy/bgc-mviews-populate-initial.yaml 2>&1 | Out-Null
    Start-Sleep -Seconds 5
    kubectl wait --for=condition=complete job/bgc-mviews-populate-initial -n data --timeout=120s 2>&1 | Out-Null
    Write-Success "Materialized views populadas"

    # ===== STATUS =====
    Write-Host "`n" -ForegroundColor Green
    Write-Host "          SETUP COMPLETO COM SUCESSO      " -ForegroundColor Green
    Write-Host "`n" -ForegroundColor Green

    Write-Host " Status dos Pods:" -ForegroundColor Yellow
    kubectl get pods -n data

    Write-Host "`n Para acessar via Ingress:" -ForegroundColor Yellow
    Write-Host "   1. Verifique se o arquivo hosts está configurado:" -ForegroundColor White
    Write-Host "      C:\Windows\System32\drivers\etc\hosts" -ForegroundColor Gray
    Write-Host "      127.0.0.1 api.bgc.local" -ForegroundColor Gray
    Write-Host "      127.0.0.1 web.bgc.local" -ForegroundColor Gray
    Write-Host ""
    Write-Host "   2. Acesse:" -ForegroundColor White
    Write-Host "      http://api.bgc.local/healthz" -ForegroundColor Gray
    Write-Host "      http://web.bgc.local" -ForegroundColor Gray

    Write-Host "`n Alternativa (port-forward):" -ForegroundColor Cyan
    Write-Host "   kubectl port-forward svc/bgc-api 8080:8080 -n data" -ForegroundColor White
    Write-Host "   curl http://localhost:8080/healthz`n" -ForegroundColor White

} catch {
    Write-Host "`n" -ForegroundColor Red
    Write-Host "          ERRO DURANTE SETUP              " -ForegroundColor Red
    Write-Host "`n" -ForegroundColor Red
    
    Write-Error "Erro: $_"
    
    Write-Host "`n Comandos úteis para diagnóstico:" -ForegroundColor Yellow
    Write-Host "   kubectl get pods -n data" -ForegroundColor White
    Write-Host "   kubectl logs <pod-name> -n data" -ForegroundColor White
    Write-Host "   kubectl describe pod <pod-name> -n data" -ForegroundColor White
    Write-Host ""
    Write-Host "   docker ps" -ForegroundColor White
    Write-Host "   k3d cluster list" -ForegroundColor White
    
    Write-Host "`n Para tentar novamente:" -ForegroundColor Cyan
    Write-Host "   .\scripts\k8s-setup.ps1`n" -ForegroundColor White
    
    exit 1
}
