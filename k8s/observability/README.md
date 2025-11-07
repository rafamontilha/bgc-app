# BGC Observability Stack - Kubernetes

Stack completa de observabilidade para o BGC App incluindo Prometheus, Grafana e Jaeger.

## Componentes

- **Prometheus** - Coleta de métricas time-series
- **Grafana** - Visualização e dashboards
- **Jaeger** - Distributed tracing
- **Alertmanager** - Gerenciamento de alertas (opcional)

## Quick Start

### Deploy completo

```bash
# Criar namespace e deployar todos os componentes
kubectl apply -f prometheus-alert-rules.yaml
kubectl apply -f prometheus-deployment.yaml
kubectl apply -f grafana-dashboards.yaml
kubectl apply -f grafana-deployment.yaml
kubectl apply -f jaeger-deployment.yaml

# Verificar status
kubectl get pods -n observability
```

### Acessar os serviços

**Opção 1: Port Forward (Desenvolvimento)**

```bash
# Prometheus
kubectl port-forward -n observability svc/prometheus 9090:9090
# Acesse: http://localhost:9090

# Grafana
kubectl port-forward -n observability svc/grafana 3000:3000
# Acesse: http://localhost:3000 (admin/admin)

# Jaeger
kubectl port-forward -n observability svc/jaeger-query 16686:16686
# Acesse: http://localhost:16686
```

**Opção 2: NodePort (k3d)**

```bash
# Obter IP do node
kubectl get nodes -o wide

# Acessar via NodePort
# Grafana: http://<node-ip>:30030
# Jaeger: http://<node-ip>:30016
```

## Configuração

### Prometheus

Targets configurados:
- `bgc-api` - API principal na porta 8080
- `integration-gateway` - Gateway na porta 8081
- `postgres-exporter` - Métricas do PostgreSQL (se instalado)

### Alertas

10 regras de alerta pré-configuradas em `prometheus-alert-rules.yaml`:
- **Critical**: HighErrorRate, APIDown
- **Warning**: HighLatencyP95, HighDatabaseLatency, etc.

### Dashboards

Dashboard pré-configurado: **BGC API Overview**
- Request Rate
- Error Rate
- Latency (P50, P95, P99)
- Database Connections
- Idempotency Cache

## Customização

### Adicionar novos scrape targets

Edite `prometheus-deployment.yaml`:

```yaml
- job_name: 'my-service'
  kubernetes_sd_configs:
    - role: pod
      namespaces:
        names:
          - data
  relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_app]
      action: keep
      regex: my-service
```

### Adicionar novos dashboards

1. Crie o JSON do dashboard no Grafana UI
2. Exporte o JSON
3. Adicione ao `grafana-dashboards.yaml` ConfigMap

## Troubleshooting

### Prometheus não está scraping a API

```bash
# Verificar targets
kubectl port-forward -n observability svc/prometheus 9090:9090
# Abra http://localhost:9090/targets

# Verificar logs
kubectl logs -n observability deployment/prometheus
```

### Grafana não mostra dados

```bash
# Verificar datasource
# Grafana UI > Configuration > Data Sources > Test

# Verificar se Prometheus tem dados
kubectl port-forward -n observability svc/prometheus 9090:9090
# Execute query: bgc_http_requests_total
```

### Jaeger não mostra traces

```bash
# Verificar se API está enviando traces
kubectl logs -n data deployment/bgc-api | grep -i "tracer"

# Verificar Jaeger collector
kubectl logs -n observability deployment/jaeger

# Verificar se OTEL endpoint está configurado
kubectl get deployment -n data bgc-api -o yaml | grep OTEL_EXPORTER
```

## Manutenção

### Limpeza de dados antigos

Prometheus retém dados por 30 dias (configurável em `--storage.tsdb.retention.time=30d`).

Jaeger usa memory storage (dados perdidos em restart). Para produção, configure storage persistente (Cassandra, Elasticsearch, Badger).

### Backup de configurações

```bash
# Exportar dashboards Grafana
kubectl get cm -n observability grafana-dashboards -o yaml > backup/grafana-dashboards.yaml

# Exportar regras de alerta
kubectl get cm -n observability prometheus-alert-rules -o yaml > backup/prometheus-alert-rules.yaml
```

## Recursos

- Prometheus: https://prometheus.io/docs/
- Grafana: https://grafana.com/docs/
- Jaeger: https://www.jaegertracing.io/docs/
- OpenTelemetry: https://opentelemetry.io/docs/
