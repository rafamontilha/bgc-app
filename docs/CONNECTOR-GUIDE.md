# Guia de Conectores - Arquitetura Híbrida

## 🎯 Visão Geral

O BGC App utiliza uma **arquitetura híbrida** para integrações com APIs externas:

- **90% dos casos**: Configuração declarativa em YAML (zero código Go!)
- **10% casos complexos**: Plugins customizados em Go

Isso permite adicionar **novas integrações em 30 minutos** ao invés de 2 dias.

---

## 🚀 Quick Start - Criar Nova Integração

### Passo 1: Copiar Template

```bash
cp config/connectors/viacep.yaml config/connectors/minha-api.yaml
```

### Passo 2: Configurar YAML

```yaml
id: minha-api
name: Minha API - Serviço XYZ
version: 1.0.0
provider: Nome do Provedor

integration:
  type: rest_api
  protocol: https

  # Autenticação
  auth:
    type: api_key  # ou mtls, oauth2, basic, jwt, none
    api_key:
      header_name: X-API-Key
      key_ref: minha-api-key  # Referência no secrets manager

  # Endpoints
  endpoints:
    consulta_dados:
      method: GET
      path: /api/v1/dados/{id}
      timeout: 30s

      headers:
        Content-Type: application/json

      path_params:
        - name: id
          type: string
          required: true

      response:
        success_status: [200]
        mapping:
          id: $.data.id
          nome: $.data.name
          status: $.data.status

# Ambientes
environments:
  production:
    base_url: https://api.exemplo.com
```

### Passo 3: Validar

```bash
# Valida YAML contra schema
ajv validate -s schemas/connector.schema.json \
    -d config/connectors/minha-api.yaml

# Ou use o script
.\scripts\validate-connector.ps1 minha-api
```

### Passo 4: Deploy

```bash
# Docker Compose
.\scripts\docker.ps1 restart

# Kubernetes
.\scripts\k8s.ps1 build
```

**Pronto!** Sua integração está funcionando. 🎉

---

## 📋 Anatomia de um Connector Config

### 1. **Metadados**

```yaml
id: nome-unico-lowercase
name: Nome Legível para Humanos
version: 1.0.0
provider: Nome da Organização
```

- `id`: Identificador único (lowercase, hyphens only)
- `name`: Nome descritivo
- `version`: Semantic versioning
- `provider`: Nome do provedor da API

### 2. **Autenticação**

#### API Key
```yaml
auth:
  type: api_key
  api_key:
    header_name: X-API-Key
    key_ref: secret-name  # Referência no Kubernetes Secret
```

#### mTLS (Certificado ICP-Brasil)
```yaml
auth:
  type: mtls
  certificate_ref: icp-brasil-receita-prod  # Referência no Certificate Manager
```

#### OAuth2
```yaml
auth:
  type: oauth2
  oauth2:
    token_url: https://auth.exemplo.com/oauth/token
    client_id: meu-client-id
    client_secret_ref: oauth-secret
    scopes: [read, write]
```

#### Sem autenticação
```yaml
auth:
  type: none
```

### 3. **Endpoints**

```yaml
endpoints:
  nome_do_endpoint:
    method: GET  # GET, POST, PUT, PATCH, DELETE
    path: /api/v1/resource/{id}  # Suporta {placeholders}
    timeout: 30s

    # Headers estáticos
    headers:
      Content-Type: application/json
      Accept: application/json

    # Parâmetros de path
    path_params:
      - name: id
        type: string
        required: true
        pattern: "^\\d+$"

    # Parâmetros de query
    query_params:
      - name: filtro
        type: string
        required: false
        default: "todos"

    # Body (para POST/PUT)
    body:
      content_type: application/json
      template: |
        {
          "campo": "{{valor}}",
          "outro": "{{outro_valor}}"
        }

    # Resposta
    response:
      success_status: [200, 201]
      error_status: [400, 404, 500]

      # Mapeamento JSONPath
      mapping:
        id: $.data.id
        nome: $.data.attributes.name
        lista: $.data.items[*].name  # Array

      # Transformações
      transforms:
        - field: cpf
          operation: format_cpf  # Built-in plugin

        - field: status
          operation: map_values
          values:
            "A": "ativo"
            "I": "inativo"
```

### 4. **Resiliência**

```yaml
resilience:
  # Retry automático
  retry:
    max_attempts: 3
    backoff: exponential  # constant, linear, exponential
    initial_interval: 1s
    max_interval: 30s

  # Circuit Breaker
  circuit_breaker:
    failure_threshold: 5
    success_threshold: 2
    timeout: 60s

  # Rate Limiting
  rate_limit:
    requests_per_minute: 60
    burst: 10
```

### 5. **Cache**

```yaml
cache:
  enabled: true
  ttl: 1h  # 5m, 1h, 24h, etc
  key_pattern: "cnpj:{cnpj}"  # Suporta {placeholders}
```

### 6. **Ambientes**

```yaml
environments:
  development:
    base_url: http://localhost:9090
    health_check: /health

  sandbox:
    base_url: https://sandbox.exemplo.com
    health_check: /status

  production:
    base_url: https://api.exemplo.com
    health_check: /status
```

### 7. **Compliance**

```yaml
compliance:
  tags: [LGPD, SOC2, ICP-Brasil]
  data_classification: confidential  # public, internal, confidential, restricted
  retention_days: 90
  encryption_required: true
```

### 8. **Governança**

```yaml
governance:
  owner_team: integrations-team
  approved_by: security-team
  last_audited: 2025-01-15
  review_frequency: quarterly  # monthly, quarterly, semi-annually, annually
```

### 9. **Observabilidade**

```yaml
observability:
  metrics_enabled: true
  tracing_enabled: true
  log_level: info  # debug, info, warn, error

  alerts:
    - type: certificate_expiry
      threshold: 30d
      channels: [slack, email]

    - type: error_rate
      threshold: 5%
      window: 5m
      channels: [pagerduty]

    - type: latency
      threshold: 2s
      window: 1m
      channels: [slack]
```

---

## 🔌 Built-in Transform Plugins

O framework inclui plugins prontos para uso:

### Formatação Brasileira

```yaml
transforms:
  - field: cnpj
    operation: format_cnpj  # 12345678000195 -> 12.345.678/0001-95

  - field: cpf
    operation: format_cpf  # 12345678901 -> 123.456.789-01

  - field: cep
    operation: format_cep  # 01310100 -> 01310-100
```

### Manipulação de Strings

```yaml
transforms:
  - field: nome
    operation: to_upper  # JOÃO SILVA

  - field: email
    operation: to_lower  # joao@exemplo.com

  - field: descricao
    operation: trim  # Remove espaços
```

### Mapeamento de Valores

```yaml
transforms:
  - field: situacao_cadastral
    operation: map_values
    values:
      "01": "ativa"
      "02": "suspensa"
      "03": "inapta"
      "04": "baixada"
```

---

## 🧩 Casos Avançados - Custom Plugins

Para APIs muito específicas (SOAP complexo, WebSocket, auth exótica), crie um plugin customizado:

### 1. Criar Plugin

```go
// services/integration-gateway/internal/plugins/meu_plugin.go
package plugins

import "github.com/bgc/integration-gateway/internal/transform"

type MeuPlugin struct{}

func (p *MeuPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
    // Sua lógica customizada aqui
    return transformedValue, nil
}

func init() {
    transform.RegisterPlugin("meu_plugin", &MeuPlugin{})
}
```

### 2. Usar no YAML

```yaml
endpoints:
  consulta_especial:
    response:
      transforms:
        - field: dados
          operation: meu_plugin  # Seu plugin customizado!
          params:
            opcao: valor
```

---

## ✅ Checklist de Qualidade

Antes de fazer deploy de um novo connector:

- [ ] ID único e descritivo (lowercase-com-hyphens)
- [ ] Versionamento semântico (1.0.0)
- [ ] Autenticação configurada corretamente
- [ ] Todos os endpoints com timeout definido
- [ ] Response mapping com JSONPath válido
- [ ] Retry policy configurada
- [ ] Rate limiting dentro dos limites da API
- [ ] Cache configurado (se aplicável)
- [ ] Compliance tags corretas
- [ ] Owner team definido
- [ ] Validação contra schema passou
- [ ] Testado em sandbox/dev
- [ ] Alertas configurados

---

## 📊 Métricas Automáticas

Todo connector automaticamente expõe métricas:

```
bgc_connector_requests_total{connector="receita-federal", status="success"}
bgc_connector_duration_seconds{connector="receita-federal", quantile="0.99"}
bgc_connector_circuit_breaker_state{connector="receita-federal"}
bgc_connector_rate_limit_remaining{connector="receita-federal"}
```

Visualize no Grafana Dashboard "Integration Health".

---

## 🐛 Troubleshooting

### Connector não carrega

```bash
# Ver logs
docker logs bgc-integration-gateway

# Validar YAML
ajv validate -s schemas/connector.schema.json -d config/connectors/seu-connector.yaml
```

### Autenticação falha

```bash
# Verificar secret
kubectl get secret seu-secret -o yaml

# Testar certificado
openssl x509 -in certs/seu-cert.pem -text -noout
```

### Response mapping não funciona

```bash
# Testar JSONPath
echo '{"data": {"name": "test"}}' | jq '$.data.name'

# Ver response raw
curl -v https://api.exemplo.com/endpoint
```

---

## 🎓 Exemplos Completos

Ver arquivos em `config/connectors/`:

- `receita-federal-cnpj.yaml` - mTLS, transformações complexas
- `viacep.yaml` - API simples, sem auth
- `siscomex.yaml` - OAuth2, múltiplos endpoints

---

## 📚 Referências

- [JSON Schema](schemas/connector.schema.json) - Schema completo
- [JSONPath Syntax](https://goessner.net/articles/JsonPath/)
- [Integration Governance](INTEGRATION-GOVERNANCE.md)
- [Certificate Management](CERTIFICATE-MANAGEMENT.md)

---

**Dúvidas?** Abra uma issue ou contate o time de integrações.
