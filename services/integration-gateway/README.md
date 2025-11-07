# Integration Gateway - Framework Híbrido

Gateway de integração com APIs externas usando configuração declarativa (YAML).

## 🎯 Objetivo

Permitir adicionar **novas integrações em 30 minutos** ao invés de 2 dias:

- **90% dos casos**: Apenas YAML (zero código Go!)
- **10% casos complexos**: Plugins customizados em Go

## 🚀 Quick Start

### 1. Configurar Variáveis de Ambiente

```bash
# Diretórios
export CONFIG_DIR=./config/connectors
export CERTS_DIR=./certs

# Ambiente
export ENVIRONMENT=development  # development, sandbox, production

# Secrets (exemplo)
export SECRET_VIACEP_KEY=sua-api-key-aqui
```

### 2. Criar Connector Config (YAML)

```yaml
# config/connectors/minha-api.yaml
id: minha-api
name: Minha API
version: 1.0.0

integration:
  type: rest_api
  auth:
    type: none

  endpoints:
    consulta:
      method: GET
      path: /api/data/{id}
      response:
        success_status: [200]
        mapping:
          id: $.data.id
          name: $.data.name

environments:
  development:
    base_url: http://localhost:9090
```

### 3. Iniciar Gateway

```bash
cd services/integration-gateway
go run cmd/gateway/main.go
```

### 4. Usar via API

```bash
# Listar conectores
curl http://localhost:8081/v1/connectors

# Executar endpoint
curl -X POST http://localhost:8081/v1/connectors/viacep/consulta_cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Response:
{
  "data": {
    "cep": "01310-100",
    "logradouro": "Avenida Paulista",
    "bairro": "Bela Vista",
    "localidade": "São Paulo",
    "uf": "SP"
  },
  "status_code": 200,
  "duration": "245ms"
}
```

## 📋 Endpoints da API

### Health Check
```bash
GET /health
```

### Listar Conectores
```bash
GET /v1/connectors
```

### Detalhes de Connector
```bash
GET /v1/connectors/{id}
```

### Executar Endpoint
```bash
POST /v1/connectors/{connectorID}/{endpointName}

Body:
{
  "param1": "value1",
  "param2": "value2"
}
```

## 🔌 Tipos de Autenticação

### None (API Pública)
```yaml
auth:
  type: none
```

### API Key
```yaml
auth:
  type: api_key
  api_key:
    header_name: X-API-Key
    key_ref: minha-api-key  # Env: SECRET_MINHA_API_KEY
```

### OAuth2
```yaml
auth:
  type: oauth2
  oauth2:
    token_url: https://auth.example.com/token
    client_id: my-client-id
    client_secret_ref: oauth-secret  # Env: SECRET_OAUTH_SECRET
    scopes: [read, write]
```

### mTLS (ICP-Brasil)
```yaml
auth:
  type: mtls
  certificate_ref: icp-brasil-receita-prod
  # Certificado em: certs/icp-brasil-receita-prod.pem
  # Chave em: certs/icp-brasil-receita-prod.key
```

## 🔄 Transform Plugins Built-in

```yaml
transforms:
  - field: cnpj
    operation: format_cnpj  # 12345678000195 -> 12.345.678/0001-95

  - field: cpf
    operation: format_cpf  # 12345678901 -> 123.456.789-01

  - field: cep
    operation: format_cep  # 01310100 -> 01310-100

  - field: nome
    operation: to_upper  # JOÃO SILVA

  - field: status
    operation: map_values
    values:
      "A": "ativo"
      "I": "inativo"
```

## 🛡️ Resiliência Automática

Todos os conectores herdam automaticamente:

```yaml
resilience:
  # Retry com backoff exponencial
  retry:
    max_attempts: 3
    backoff: exponential
    initial_interval: 1s
    max_interval: 30s

  # Circuit Breaker
  circuit_breaker:
    failure_threshold: 5
    timeout: 60s

  # Rate Limiting
  rate_limit:
    requests_per_minute: 60
    burst: 10
```

## 📦 Estrutura de Arquivos

```
services/integration-gateway/
├── cmd/gateway/           # Entry point
│   └── main.go
│
├── internal/
│   ├── framework/        # Core framework
│   │   ├── types.go
│   │   ├── httpclient.go
│   │   └── executor.go
│   │
│   ├── auth/             # Autenticação
│   │   ├── engine.go
│   │   ├── mtls.go
│   │   ├── oauth2.go
│   │   └── apikey.go
│   │
│   ├── transform/        # Transformações
│   │   └── engine.go
│   │
│   └── registry/         # Registry de configs
│       ├── loader.go
│       └── registry.go
│
├── config/connectors/    # Connector configs (YAML)
│   ├── receita-federal-cnpj.yaml
│   └── viacep.yaml
│
└── go.mod
```

## 🧪 Exemplo Completo: ViaCEP

**Config:** `config/connectors/viacep.yaml`

```yaml
id: viacep
name: ViaCEP - Consulta de CEP
version: 1.0.0

integration:
  type: rest_api
  auth:
    type: none

  endpoints:
    consulta_cep:
      method: GET
      path: /ws/{cep}/json/
      timeout: 10s

      path_params:
        - name: cep
          type: string
          required: true

      response:
        success_status: [200]
        mapping:
          cep: $.cep
          logradouro: $.logradouro
          bairro: $.bairro
          localidade: $.localidade
          uf: $.uf
        transforms:
          - field: cep
            operation: format_cep

  cache:
    enabled: true
    ttl: 24h
    key_pattern: "cep:{cep}"

environments:
  production:
    base_url: https://viacep.com.br
```

**Uso:**

```bash
curl -X POST http://localhost:8081/v1/connectors/viacep/consulta_cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'
```

**Response:**

```json
{
  "data": {
    "cep": "01310-100",
    "logradouro": "Avenida Paulista",
    "bairro": "Bela Vista",
    "localidade": "São Paulo",
    "uf": "SP"
  },
  "status_code": 200,
  "duration": "245ms"
}
```

## 🧩 Criar Custom Plugin

Para casos muito específicos:

```go
// internal/transform/my_plugin.go
package transform

type MyCustomPlugin struct{}

func (p *MyCustomPlugin) Transform(value interface{}, params map[string]interface{}) (interface{}, error) {
    // Sua lógica customizada
    return transformedValue, nil
}

// Registrar em cmd/gateway/main.go
transformEngine.RegisterPlugin("my_custom", &transform.MyCustomPlugin{})
```

Usar no YAML:

```yaml
transforms:
  - field: dados
    operation: my_custom
    params:
      option: value
```

## 📚 Documentação Completa

- [Connector Guide](../../docs/CONNECTOR-GUIDE.md) - Guia completo de uso
- [JSON Schema](../../schemas/connector.schema.json) - Schema de validação
- [Integration Governance](../../docs/INTEGRATION-GOVERNANCE.md) - Governança

## 🐛 Troubleshooting

### Connector não carrega

```bash
# Validar YAML
ajv validate -s ../../schemas/connector.schema.json -d config/connectors/seu-connector.yaml

# Ver logs
export LOG_LEVEL=debug
go run cmd/gateway/main.go
```

### Autenticação falha

```bash
# Verificar env var de secret
echo $SECRET_SUA_API_KEY

# Verificar certificado
openssl x509 -in certs/seu-cert.pem -text -noout
```

## 🚀 Build e Deploy

### Build

```bash
go build -o integration-gateway cmd/gateway/main.go
```

### Docker

```bash
docker build -t bgc/integration-gateway:latest .
```

### Kubernetes

Ver: `../../k8s/integration-gateway/`

---

**Desenvolvido com ❤️ para escalar integrações de forma inteligente**
