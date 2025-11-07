# 🚀 Quick Start - Observabilidade BGC App

## ✅ Serviços Ativos

Todos os serviços de observabilidade estão rodando!

---

## 📊 Acessar as UIs

### 1. **Prometheus** - Métricas
**URL:** http://localhost:9090

**O que fazer:**
- Clique em "Status" → "Targets" para ver os serviços sendo monitorados
- Experimente estas queries no "Graph":
  ```promql
  # Taxa de requisições por segundo
  rate(bgc_http_requests_total[1m])

  # Latência P95
  histogram_quantile(0.95, rate(bgc_http_request_duration_seconds_bucket[5m]))

  # Conexões do banco de dados
  bgc_db_connections_open
  ```

---

### 2. **Grafana** - Dashboards
**URL:** http://localhost:3001
**Login:** admin / admin

**O que fazer:**
1. Após login, vá em **Dashboards** → **Browse**
2. Procure por **"BGC API Overview"**
3. Visualize os painéis:
   - Request Rate
   - Error Rate
   - Latency (P50, P95, P99)
   - Database Connections

**PromQL Queries úteis:**
```promql
# Taxa de requisições
sum(rate(bgc_http_requests_total[5m])) by (path)

# Taxa de erros
sum(rate(bgc_http_requests_total{status="5xx"}[5m])) / sum(rate(bgc_http_requests_total[5m])) * 100

# Conexões do DB
bgc_db_connections_in_use
```

---

### 3. **Jaeger** - Distributed Tracing
**URL:** http://localhost:16686

**O que fazer:**
1. No dropdown "Service", selecione **"bgc-api"**
2. Clique em **"Find Traces"**
3. Clique em qualquer trace para ver detalhes:
   - Span de HTTP request
   - Span de database query
   - Duração de cada operação
   - Parâmetros das queries

---

### 4. **API Metrics** - Endpoint Prometheus
**URL:** http://localhost:8080/metrics

Retorna métricas em formato Prometheus:
```
bgc_http_requests_total{method="GET",path="/v1/market/size",status="2xx"} 1
bgc_db_query_duration_seconds_sum{operation="SELECT",table="v_tam_by_year_chapter"} 0.037
...
```

---

### 5. **API Metrics JSON** - Formato Legacy
**URL:** http://localhost:8080/metrics/json

Retorna métricas em JSON (compatibilidade):
```json
{
  "uptime_seconds": 238,
  "requests_total": 29,
  "requests_by_status": {"200": 28, "400": 1},
  "routes": {...}
}
```

---

## 🧪 Gerar Tráfego para Ver Métricas

Execute algumas requisições para gerar dados:

```bash
# Requisição válida
curl "http://localhost:8080/v1/market/size?year_from=2020&year_to=2023&metric=TAM"

# Requisição inválida (para gerar erro 4xx)
curl "http://localhost:8080/v1/market/size?year_from=2020"

# Fazer várias requisições
for i in {1..10}; do
  curl -s "http://localhost:8080/v1/market/size?year_from=2020&year_to=2021&metric=TAM" > /dev/null
  echo "Request $i completed"
done
```

Depois, veja os resultados em:
- **Prometheus:** http://localhost:9090/graph
- **Grafana:** http://localhost:3001
- **Jaeger:** http://localhost:16686

---

## 📈 Métricas Disponíveis

### HTTP Metrics
- `bgc_http_requests_total` - Total de requisições
- `bgc_http_request_duration_seconds` - Latência (histogram)
- `bgc_http_requests_in_flight` - Requisições em processamento

### Database Metrics
- `bgc_db_queries_total` - Total de queries
- `bgc_db_query_duration_seconds` - Duração de queries (histogram)
- `bgc_db_connections_open` - Conexões abertas
- `bgc_db_connections_in_use` - Conexões em uso
- `bgc_db_connections_idle` - Conexões ociosas

### Application Metrics
- `bgc_idempotency_cache_hits_total` - Cache hits
- `bgc_idempotency_cache_misses_total` - Cache misses
- `bgc_idempotency_cache_size` - Tamanho do cache
- `bgc_errors_total` - Total de erros

---

## 🔍 Exemplos de Queries PromQL

### Taxa de Requisições
```promql
rate(bgc_http_requests_total[1m])
```

### Latência P95
```promql
histogram_quantile(0.95,
  sum(rate(bgc_http_request_duration_seconds_bucket[5m])) by (le)
)
```

### Taxa de Erros (%)
```promql
sum(rate(bgc_http_requests_total{status="5xx"}[5m]))
/
sum(rate(bgc_http_requests_total[5m])) * 100
```

### Database Query Rate
```promql
rate(bgc_db_queries_total[1m])
```

### Connection Pool Usage (%)
```promql
(bgc_db_connections_in_use / bgc_db_connections_open) * 100
```

---

## 🐛 Troubleshooting

### Serviços não acessíveis?

```bash
# Verificar status dos containers
cd bgcstack
docker-compose ps

# Ver logs
docker logs bgc_prometheus
docker logs bgc_grafana
docker logs bgc_jaeger
docker logs bgc_api
```

### Métricas não aparecem?

```bash
# Verificar se Prometheus está fazendo scrape
curl http://localhost:9090/api/v1/targets

# Verificar métricas da API
curl http://localhost:8080/metrics | grep bgc_
```

### Grafana não mostra dados?

1. Verifique se o datasource Prometheus está configurado:
   - Grafana → Configuration → Data Sources
   - URL deve ser: http://prometheus:9090
2. Teste a conexão
3. Aguarde alguns minutos para acumular dados

---

## 📚 Documentação Completa

Para mais detalhes, veja:
- `docs/OBSERVABILITY.md` - Guia completo de observabilidade
- `docs/EPIC-2-COMPLETE.md` - Resumo do Épico 2
- `k8s/observability/README.md` - Deploy Kubernetes

---

## 🎯 Próximos Passos

1. **Explorar Grafana:** Crie dashboards customizados
2. **Testar Alertas:** Configure notificações (Slack, Email)
3. **Analisar Traces:** Use Jaeger para debug de performance
4. **Monitorar Produção:** Deploy da stack em staging/prod

---

**Desenvolvido com ❤️ e observabilidade desde o dia 1!**
