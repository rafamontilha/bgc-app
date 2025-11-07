# 🚀 Próximos Passos - Épico 1

## 📋 Status Atual

✅ **MVP Funcional Implementado (70%)**
- Framework core completo
- Auth Engine (mTLS, OAuth2, API Key)
- Transform Engine + plugins
- API REST funcional
- 2 exemplos funcionais (ViaCEP, Receita Federal)
- Documentação completa

## 🎯 Para Completar 100%

### **Fase 1: Validação e Testes (1-2 dias)**

#### 1.1 Testar Localmente

```bash
# 1. Iniciar o Integration Gateway
cd services/integration-gateway

# 2. Configurar env
export CONFIG_DIR=../../config/connectors
export CERTS_DIR=../../certs
export ENVIRONMENT=development

# 3. Iniciar
go run cmd/gateway/main.go

# 4. Em outro terminal, testar ViaCEP
curl -X POST http://localhost:8081/v1/connectors/viacep/consulta_cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'
```

**Sucesso se:** Retornar dados do CEP formatados com "01310-100"

#### 1.2 Criar Testes Unitários

```bash
# Criar arquivo de teste
touch services/integration-gateway/internal/framework/executor_test.go
touch services/integration-gateway/internal/auth/oauth2_test.go
touch services/integration-gateway/internal/transform/engine_test.go

# Executar testes
cd services/integration-gateway
go test ./... -v -cover
```

**Meta:** 80%+ de cobertura no framework core

#### 1.3 Teste de Integração End-to-End

```bash
# Criar teste E2E
mkdir -p tests/integration
touch tests/integration/viacep_test.go

# Executar
go test ./tests/integration/... -v
```

---

### **Fase 2: Observabilidade (1 dia)**

#### 2.1 Adicionar Prometheus Metrics

```go
// internal/observability/metrics.go
package observability

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    ConnectorRequests = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "bgc_connector_requests_total",
            Help: "Total requests per connector",
        },
        []string{"connector", "endpoint", "status"},
    )

    ConnectorDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "bgc_connector_duration_seconds",
            Help: "Request duration per connector",
        },
        []string{"connector", "endpoint"},
    )
)
```

#### 2.2 Adicionar ao Executor

```go
// internal/framework/executor.go
import "github.com/bgc/integration-gateway/internal/observability"

func (e *Executor) Execute(ctx *ExecutionContext) (*ExecutionResult, error) {
    startTime := time.Now()

    // ... código existente ...

    // Registrar métricas
    status := "success"
    if err != nil {
        status = "error"
    }

    observability.ConnectorRequests.WithLabelValues(
        ctx.ConnectorID,
        ctx.EndpointName,
        status,
    ).Inc()

    observability.ConnectorDuration.WithLabelValues(
        ctx.ConnectorID,
        ctx.EndpointName,
    ).Observe(time.Since(startTime).Seconds())

    return result, err
}
```

#### 2.3 Expor Endpoint de Métricas

```go
// cmd/gateway/main.go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Adicionar ao router
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

**Testar:**
```bash
curl http://localhost:8081/metrics
```

---

### **Fase 3: Deployment (1-2 dias)**

#### 3.1 Atualizar Docker Compose

```yaml
# bgcstack/docker-compose.yml
services:
  # ... serviços existentes ...

  integration-gateway:
    build:
      context: ../services/integration-gateway
      dockerfile: Dockerfile
    container_name: bgc_integration_gateway
    ports:
      - "8081:8081"
    environment:
      - CONFIG_DIR=/app/config/connectors
      - CERTS_DIR=/app/certs
      - ENVIRONMENT=development
      - LOG_LEVEL=info
    volumes:
      - ../config/connectors:/app/config/connectors:ro
      - ../certs:/app/certs:ro
    networks:
      - bgc_network
    depends_on:
      - postgres
    restart: unless-stopped
```

#### 3.2 Criar Dockerfile

```dockerfile
# services/integration-gateway/Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o integration-gateway cmd/gateway/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/integration-gateway .

EXPOSE 8081

CMD ["./integration-gateway"]
```

#### 3.3 Kubernetes Manifests

```yaml
# k8s/integration-gateway/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: integration-gateway
  namespace: data
spec:
  replicas: 2
  selector:
    matchLabels:
      app: integration-gateway
  template:
    metadata:
      labels:
        app: integration-gateway
    spec:
      containers:
      - name: gateway
        image: bgc/integration-gateway:latest
        ports:
        - containerPort: 8081
        env:
        - name: CONFIG_DIR
          value: /app/config/connectors
        - name: CERTS_DIR
          value: /app/certs
        - name: ENVIRONMENT
          value: production
        volumeMounts:
        - name: connectors-config
          mountPath: /app/config/connectors
          readOnly: true
        - name: certificates
          mountPath: /app/certs
          readOnly: true
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: connectors-config
        configMap:
          name: connectors-config
      - name: certificates
        secret:
          secretName: icp-certificates

---
apiVersion: v1
kind: Service
metadata:
  name: integration-gateway
  namespace: data
spec:
  selector:
    app: integration-gateway
  ports:
  - port: 8081
    targetPort: 8081
  type: ClusterIP
```

#### 3.4 Scripts de Teste

```powershell
# scripts/test-local.ps1
param(
    [string]$Connector = "viacep"
)

Write-Host "Testing Integration Gateway locally..." -ForegroundColor Green

# Start gateway em background
$gateway = Start-Process -FilePath "go" -ArgumentList "run", "services/integration-gateway/cmd/gateway/main.go" -PassThru -NoNewWindow

Start-Sleep -Seconds 3

# Test health
$health = Invoke-RestMethod -Uri "http://localhost:8081/health"
Write-Host "Health: $($health.status)" -ForegroundColor Cyan

# List connectors
$connectors = Invoke-RestMethod -Uri "http://localhost:8081/v1/connectors"
Write-Host "Connectors loaded: $($connectors.Count)" -ForegroundColor Cyan

# Test ViaCEP
if ($Connector -eq "viacep") {
    Write-Host "`nTesting ViaCEP..." -ForegroundColor Yellow
    $body = @{ cep = "01310100" } | ConvertTo-Json
    $result = Invoke-RestMethod -Method Post -Uri "http://localhost:8081/v1/connectors/viacep/consulta_cep" -Body $body -ContentType "application/json"
    Write-Host "CEP: $($result.data.cep)" -ForegroundColor Green
    Write-Host "Logradouro: $($result.data.logradouro)" -ForegroundColor Green
}

# Stop gateway
Stop-Process -Id $gateway.Id

Write-Host "`nTests completed!" -ForegroundColor Green
```

---

### **Fase 4: Certificate Manager Completo (2-3 dias)**

#### 4.1 Database Schema

```sql
-- db/migrations/0005_certificate_manager.sql
CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    cert_type VARCHAR(10) NOT NULL,  -- A1, A3
    issuer VARCHAR(255),
    subject VARCHAR(255),
    serial_number VARCHAR(255),
    valid_from TIMESTAMP NOT NULL,
    valid_to TIMESTAMP NOT NULL,
    secret_ref VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,  -- active, expiring, expired, revoked
    last_validated TIMESTAMP,
    last_rotated TIMESTAMP,
    owner_team VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_certs_expiry ON certificates(valid_to) WHERE status = 'active';
```

#### 4.2 Implementar Auto-Rotation

```go
// services/cert-manager/internal/rotation/scheduler.go
package rotation

import (
    "time"
)

type Scheduler struct {
    interval time.Duration
    registry CertificateRegistry
}

func (s *Scheduler) Start() {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()

    for range ticker.C {
        s.checkExpiringCertificates()
    }
}

func (s *Scheduler) checkExpiringCertificates() {
    // Busca certificados expirando em 30 dias
    certs := s.registry.FindExpiringBefore(time.Now().Add(30 * 24 * time.Hour))

    for _, cert := range certs {
        // Envia alerta
        // Inicia processo de renovação
    }
}
```

---

## 🎯 Checklist de Conclusão

### **MVP (70% - FEITO)**
- [x] Framework core
- [x] Auth Engine (mTLS, OAuth2, API Key)
- [x] Transform Engine + plugins
- [x] API REST
- [x] 2 exemplos funcionais
- [x] Documentação

### **Produção (100%)**
- [ ] Testes unitários (80%+ coverage)
- [ ] Testes de integração E2E
- [ ] Prometheus metrics
- [ ] Structured logging
- [ ] Docker Compose atualizado
- [ ] Kubernetes manifests
- [ ] Script test-local.ps1
- [ ] Script test-k3d.ps1
- [ ] Certificate Manager database
- [ ] Certificate auto-rotation
- [ ] Alertas de expiração

---

## 📊 Priorização Recomendada

### **Semana 1 (Crítico)**
1. ✅ Testes locais manuais (ViaCEP)
2. ✅ Testes unitários básicos
3. ✅ Docker Compose atualizado

### **Semana 2 (Importante)**
4. ✅ Observabilidade (metrics)
5. ✅ Kubernetes manifests
6. ✅ Scripts de teste

### **Semana 3 (Desejável)**
7. ✅ Certificate Manager database
8. ✅ Auto-rotation
9. ✅ Testes E2E completos

---

## 🚀 Começar Agora

**Passo 1:** Testar localmente
```bash
cd services/integration-gateway
export CONFIG_DIR=../../config/connectors
export CERTS_DIR=../../certs
go run cmd/gateway/main.go
```

**Passo 2:** Em outro terminal
```bash
curl http://localhost:8081/v1/connectors
curl -X POST http://localhost:8081/v1/connectors/viacep/consulta_cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'
```

**Sucesso?** ✅ Pronto para adicionar mais integrações!

---

**Próxima integração sugerida:** Receita Federal (CNPJ real com certificado ICP-Brasil sandbox)

**Dúvidas?** Ver `docs/CONNECTOR-GUIDE.md`
