# API do Simulador de Destinos de Exportação 🌍

Documentação completa do endpoint `/v1/simulator/destinations` para recomendação de mercados de exportação.

---

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Endpoint](#endpoint)
- [Autenticação e Rate Limiting](#autenticação-e-rate-limiting)
- [Request](#request)
- [Response](#response)
- [Algoritmo de Scoring](#algoritmo-de-scoring)
- [Exemplos de Uso](#exemplos-de-uso)
- [Códigos de Erro](#códigos-de-erro)
- [Considerações de Performance](#considerações-de-performance)

---

## Visão Geral

O Simulador de Destinos de Exportação analisa dados históricos de comércio exterior (ComexStat) e recomenda os melhores mercados para exportação de um produto específico (identificado por NCM).

### Funcionalidades

✅ **Análise de 100+ países** baseada em dados reais de exportação
✅ **Scoring inteligente** considerando tamanho de mercado, crescimento, preço e logística
✅ **Estimativas financeiras** de margem, custos logísticos e tarifas
✅ **Filtragem por países** específicos (opcional)
✅ **Rate limiting freemium** (5 simulações/dia para free, ilimitado para premium)
✅ **Processamento rápido** (~50-200ms)

---

## Endpoint

### POST `/v1/simulator/destinations`

Simula destinos de exportação para um produto específico.

**Base URL:**
- Desenvolvimento: `http://localhost:8080`
- Produção: `https://api.brasilglobalconect.com`

**Content-Type:** `application/json`

---

## Autenticação e Rate Limiting

### Tier Free (Anônimo)
- **Limite:** 5 simulações por dia
- **Identificação:** Por endereço IP
- **Headers de resposta:**
  - `X-RateLimit-Limit: 5`
  - `X-RateLimit-Remaining: N` (quantas simulações restam)
  - `X-RateLimit-Reset: TIMESTAMP` (quando o limite reseta)

### Tier Premium (Autenticado)
- **Limite:** Ilimitado
- **Headers de resposta:**
  - `X-RateLimit-Limit: unlimited`
  - `X-RateLimit-Remaining: unlimited`

### Resposta de Rate Limit Excedido

```json
{
  "error": "rate_limit_exceeded",
  "message": "Limite de 5 simulações por dia atingido. Faça upgrade para Pro para simulações ilimitadas.",
  "remaining": 0,
  "reset_at": 1737590400
}
```

**Status Code:** `429 Too Many Requests`

---

## Request

### Schema

```json
{
  "ncm": "string (required, 8 dígitos)",
  "volume_kg": number (optional),
  "countries": ["string"] (optional),
  "max_results": number (optional, 1-50, default: 10)
}
```

### Campos

| Campo | Tipo | Obrigatório | Descrição | Validação |
|-------|------|-------------|-----------|-----------|
| `ncm` | string | ✅ Sim | Código NCM de 8 dígitos do produto | Exatamente 8 dígitos numéricos |
| `volume_kg` | number | ❌ Não | Volume estimado de exportação em kg | Deve ser > 0 |
| `countries` | array[string] | ❌ Não | Lista de países para filtrar (códigos ISO 2) | Códigos válidos: US, CN, BR, etc. |
| `max_results` | number | ❌ Não | Número máximo de destinos retornados | Entre 1 e 50 (default: 10) |

### Exemplo de Request Mínimo

```json
{
  "ncm": "12345678"
}
```

### Exemplo de Request Completo

```json
{
  "ncm": "84715000",
  "volume_kg": 5000,
  "countries": ["US", "CN", "DE", "JP"],
  "max_results": 5
}
```

---

## Response

### Schema de Sucesso

```json
{
  "destinations": [
    {
      "country_code": "string",
      "country_name": "string",
      "score": number (0-10),
      "rank": number,
      "demand": "Alto|Médio|Baixo",
      "estimated_margin_pct": number,
      "logistics_cost_usd": number,
      "tariff_rate_pct": number,
      "lead_time_days": number,
      "market_size_usd": number,
      "growth_rate_pct": number,
      "price_per_kg_usd": number,
      "distance_km": number,
      "region": "string",
      "flag_emoji": "string",
      "recommendation_reason": "string"
    }
  ],
  "metadata": {
    "ncm": "string",
    "product_name": "string",
    "analysis_date": "string (ISO 8601)",
    "total_destinations": number,
    "cache_hit": boolean,
    "cache_level": "l1|l2|l3",
    "processing_time_ms": number
  }
}
```

### Campos do Destination

| Campo | Descrição | Unidade |
|-------|-----------|---------|
| `country_code` | Código ISO do país (ex: US, CN, BR) | - |
| `country_name` | Nome do país em português | - |
| `score` | Score geral da recomendação | 0.0 - 10.0 |
| `rank` | Posição no ranking (1 = melhor) | 1, 2, 3... |
| `demand` | Nível de demanda do mercado | Alto / Médio / Baixo |
| `estimated_margin_pct` | Margem estimada de lucro | % |
| `logistics_cost_usd` | Custo logístico estimado | USD |
| `tariff_rate_pct` | Taxa de tarifa de importação | % |
| `lead_time_days` | Tempo estimado de entrega | dias |
| `market_size_usd` | Tamanho do mercado (últimos 12 meses) | USD |
| `growth_rate_pct` | Taxa de crescimento anual | % |
| `price_per_kg_usd` | Preço médio por kg | USD/kg |
| `distance_km` | Distância do Brasil | km |
| `region` | Região geográfica | Americas, Europe, Asia, etc. |
| `flag_emoji` | Emoji da bandeira do país | 🇺🇸 🇨🇳 🇧🇷 |
| `recommendation_reason` | Explicação do score | texto |

### Exemplo de Response de Sucesso

```json
{
  "destinations": [
    {
      "country_code": "US",
      "country_name": "Estados Unidos",
      "score": 8.5,
      "rank": 1,
      "demand": "Alto",
      "estimated_margin_pct": 25.0,
      "logistics_cost_usd": 375.0,
      "tariff_rate_pct": 8.0,
      "lead_time_days": 22,
      "market_size_usd": 150000000,
      "growth_rate_pct": 15.5,
      "price_per_kg_usd": 85.50,
      "distance_km": 7500,
      "region": "Americas",
      "flag_emoji": "🇺🇸",
      "recommendation_reason": "Mercado altamente atrativo com grande potencial de crescimento e demanda consolidada"
    },
    {
      "country_code": "CN",
      "country_name": "China",
      "score": 7.8,
      "rank": 2,
      "demand": "Alto",
      "estimated_margin_pct": 25.0,
      "logistics_cost_usd": 850.0,
      "tariff_rate_pct": 15.0,
      "lead_time_days": 41,
      "market_size_usd": 120000000,
      "growth_rate_pct": 22.0,
      "price_per_kg_usd": 80.00,
      "distance_km": 17000,
      "region": "Asia",
      "flag_emoji": "🇨🇳",
      "recommendation_reason": "Mercado promissor com bom equilíbrio entre demanda, crescimento e custos logísticos"
    }
  ],
  "metadata": {
    "ncm": "84715000",
    "product_name": "Unidades de processamento",
    "analysis_date": "2025-01-21T14:30:00Z",
    "total_destinations": 2,
    "cache_hit": false,
    "cache_level": "",
    "processing_time_ms": 127
  }
}
```

**Status Code:** `200 OK`

---

## Algoritmo de Scoring

O score de cada destino (0-10) é calculado usando uma **média ponderada de 4 métricas normalizadas**:

### Pesos Padrão

| Métrica | Peso | Descrição |
|---------|------|-----------|
| **Market Size** | 40% | Tamanho do mercado em USD (últimos 12 meses) |
| **Growth Rate** | 30% | Taxa de crescimento anual comparada ao período anterior |
| **Price per Kg** | 20% | Preço médio por kg (maior = melhor margem) |
| **Distance** | 10% | Distância do Brasil (menor = menor custo logístico) |

### Fórmula

```
score_normalizado = (market_size_norm × 0.40) +
                    (growth_rate_norm × 0.30) +
                    (price_norm × 0.20) +
                    (distance_norm × 0.10)

score_final = score_normalizado × 10
```

Onde cada métrica é normalizada para 0-1 usando o valor máximo encontrado nos dados.

### Classificação de Demanda

| Market Size (USD) | Demanda |
|-------------------|---------|
| > 100 milhões | **Alto** |
| 10M - 100M | **Médio** |
| < 10 milhões | **Baixo** |

### Razões de Recomendação

| Score | Razão |
|-------|-------|
| ≥ 8.0 | "Mercado altamente atrativo com grande potencial de crescimento e demanda consolidada" |
| 6.0 - 7.9 | "Mercado promissor com bom equilíbrio entre demanda, crescimento e custos logísticos" |
| 4.0 - 5.9 | "Mercado em desenvolvimento com oportunidades emergentes" |
| < 4.0 | "Mercado em fase inicial ou com barreiras significativas" |

---

## Exemplos de Uso

### cURL

```bash
# Request básico
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -H "Content-Type: application/json" \
  -d '{
    "ncm": "84715000"
  }'

# Request com filtro de países e volume
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -H "Content-Type: application/json" \
  -d '{
    "ncm": "84715000",
    "volume_kg": 5000,
    "countries": ["US", "CN", "DE"],
    "max_results": 5
  }'

# Request com autenticação (premium)
curl -X POST https://api.brasilglobalconect.com/v1/simulator/destinations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "ncm": "84715000"
  }'
```

### JavaScript (Fetch API)

```javascript
async function simulateDestinations(ncm, volumeKg = null) {
  try {
    const response = await fetch('http://localhost:8080/v1/simulator/destinations', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization': 'Bearer YOUR_TOKEN' // Para premium
      },
      body: JSON.stringify({
        ncm: ncm,
        volume_kg: volumeKg,
        max_results: 10
      })
    });

    // Verificar rate limit
    const remaining = response.headers.get('X-RateLimit-Remaining');
    console.log(`Simulações restantes: ${remaining}`);

    if (!response.ok) {
      if (response.status === 429) {
        const resetAt = response.headers.get('X-RateLimit-Reset');
        throw new Error(`Rate limit excedido. Reseta em: ${new Date(resetAt * 1000)}`);
      }
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    console.log(`Encontrados ${data.metadata.total_destinations} destinos`);
    console.log(`Tempo de processamento: ${data.metadata.processing_time_ms}ms`);

    return data.destinations;
  } catch (error) {
    console.error('Erro ao simular destinos:', error);
    throw error;
  }
}

// Uso
simulateDestinations('84715000', 5000)
  .then(destinations => {
    destinations.forEach(dest => {
      console.log(`${dest.rank}. ${dest.flag_emoji} ${dest.country_name} - Score: ${dest.score.toFixed(1)}`);
    });
  });
```

### Python (requests)

```python
import requests
from datetime import datetime

def simulate_destinations(ncm: str, volume_kg: float = None, countries: list = None):
    """
    Simula destinos de exportação para um NCM específico.

    Args:
        ncm: Código NCM de 8 dígitos
        volume_kg: Volume estimado em kg (opcional)
        countries: Lista de países para filtrar (opcional)

    Returns:
        dict: Resposta da API com destinos recomendados
    """
    url = "http://localhost:8080/v1/simulator/destinations"

    payload = {"ncm": ncm}

    if volume_kg:
        payload["volume_kg"] = volume_kg
    if countries:
        payload["countries"] = countries

    headers = {
        "Content-Type": "application/json",
        # "Authorization": "Bearer YOUR_TOKEN"  # Para premium
    }

    try:
        response = requests.post(url, json=payload, headers=headers)

        # Verificar rate limit
        remaining = response.headers.get('X-RateLimit-Remaining')
        print(f"Simulações restantes: {remaining}")

        if response.status_code == 429:
            reset_at = int(response.headers.get('X-RateLimit-Reset'))
            reset_time = datetime.fromtimestamp(reset_at)
            raise Exception(f"Rate limit excedido. Reseta em: {reset_time}")

        response.raise_for_status()

        data = response.json()

        print(f"Encontrados {data['metadata']['total_destinations']} destinos")
        print(f"Tempo de processamento: {data['metadata']['processing_time_ms']}ms")

        return data

    except requests.exceptions.RequestException as e:
        print(f"Erro ao simular destinos: {e}")
        raise

# Uso básico
result = simulate_destinations("84715000")

# Imprimir resultados
for dest in result['destinations']:
    print(f"{dest['rank']}. {dest['flag_emoji']} {dest['country_name']} - Score: {dest['score']:.1f}")
    print(f"   Market Size: ${dest['market_size_usd']:,.0f}")
    print(f"   Growth Rate: {dest['growth_rate_pct']:.1f}%")
    print(f"   Demand: {dest['demand']}")
    print()

# Uso com filtro
result = simulate_destinations(
    ncm="84715000",
    volume_kg=5000,
    countries=["US", "CN", "DE", "JP"]
)
```

### TypeScript (axios)

```typescript
import axios, { AxiosError } from 'axios';

interface SimulatorRequest {
  ncm: string;
  volume_kg?: number;
  countries?: string[];
  max_results?: number;
}

interface Destination {
  country_code: string;
  country_name: string;
  score: number;
  rank: number;
  demand: 'Alto' | 'Médio' | 'Baixo';
  estimated_margin_pct: number;
  logistics_cost_usd: number;
  tariff_rate_pct: number;
  lead_time_days: number;
  market_size_usd: number;
  growth_rate_pct: number;
  price_per_kg_usd: number;
  distance_km: number;
  region: string;
  flag_emoji: string;
  recommendation_reason: string;
}

interface SimulatorResponse {
  destinations: Destination[];
  metadata: {
    ncm: string;
    product_name: string;
    analysis_date: string;
    total_destinations: number;
    cache_hit: boolean;
    cache_level?: string;
    processing_time_ms: number;
  };
}

async function simulateDestinations(
  request: SimulatorRequest
): Promise<SimulatorResponse> {
  try {
    const response = await axios.post<SimulatorResponse>(
      'http://localhost:8080/v1/simulator/destinations',
      request,
      {
        headers: {
          'Content-Type': 'application/json',
          // 'Authorization': 'Bearer YOUR_TOKEN' // Para premium
        }
      }
    );

    // Log rate limit info
    const remaining = response.headers['x-ratelimit-remaining'];
    console.log(`Simulações restantes: ${remaining}`);

    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;

      if (axiosError.response?.status === 429) {
        const resetAt = axiosError.response.headers['x-ratelimit-reset'];
        const resetTime = new Date(parseInt(resetAt) * 1000);
        throw new Error(`Rate limit excedido. Reseta em: ${resetTime}`);
      }

      throw new Error(`API error: ${axiosError.message}`);
    }
    throw error;
  }
}

// Uso
const result = await simulateDestinations({
  ncm: '84715000',
  volume_kg: 5000,
  countries: ['US', 'CN', 'DE'],
  max_results: 5
});

console.log(`Top ${result.destinations.length} destinos:`);
result.destinations.forEach(dest => {
  console.log(`${dest.rank}. ${dest.flag_emoji} ${dest.country_name}`);
  console.log(`   Score: ${dest.score.toFixed(1)} | ${dest.demand}`);
  console.log(`   ${dest.recommendation_reason}`);
});
```

---

## Códigos de Erro

### 400 Bad Request

**Causa:** Request inválido ou campos faltando

```json
{
  "error": "invalid_request",
  "message": "Formato de requisição inválido",
  "details": "json: cannot unmarshal string into Go value of type handlers.SimulatorRequest"
}
```

ou

```json
{
  "error": "validation_error",
  "message": "NCM deve ter exatamente 8 dígitos"
}
```

**Possíveis causas:**
- NCM faltando ou com formato inválido (não são 8 dígitos)
- NCM contém caracteres não numéricos
- `volume_kg` é negativo ou zero
- `max_results` está fora do range 1-50
- JSON malformado

### 404 Not Found

**Causa:** NCM não encontrado ou sem dados

```json
{
  "error": "ncm_not_found",
  "message": "NCM não encontrado na base de dados"
}
```

ou

```json
{
  "error": "no_data_available",
  "message": "Dados não disponíveis para o NCM solicitado"
}
```

### 422 Unprocessable Entity

**Causa:** Dados insuficientes para gerar recomendações

```json
{
  "error": "insufficient_data",
  "message": "Dados insuficientes para gerar recomendações"
}
```

### 429 Too Many Requests

**Causa:** Rate limit excedido (tier free)

```json
{
  "error": "rate_limit_exceeded",
  "message": "Limite de 5 simulações por dia atingido. Faça upgrade para Pro para simulações ilimitadas.",
  "remaining": 0,
  "reset_at": 1737590400
}
```

**Headers:**
- `X-RateLimit-Limit: 5`
- `X-RateLimit-Remaining: 0`
- `X-RateLimit-Reset: 1737590400` (Unix timestamp)

### 500 Internal Server Error

**Causa:** Erro interno do servidor

```json
{
  "error": "internal_error",
  "message": "Erro interno ao processar requisição",
  "details": "database connection timeout"
}
```

---

## Considerações de Performance

### Tempo de Resposta

| Cenário | Tempo Médio | P95 | P99 |
|---------|-------------|-----|-----|
| **Cache Hit (L1)** | ~5ms | 10ms | 15ms |
| **Cache Hit (L2)** | ~15ms | 25ms | 35ms |
| **Database Query** | ~120ms | 200ms | 300ms |
| **First Request (Cold)** | ~150ms | 250ms | 400ms |

### Otimizações Implementadas

1. **Cache Multinível** (L1 → L2 → L3)
   - L1 (Ristretto): In-memory, 100MB, LFU
   - L2 (Redis): Distribuído, 512MB, LRU
   - L3 (PostgreSQL): Materialized Views (futuro)

2. **Queries Otimizadas**
   - CTEs para agregação eficiente
   - Índices em `co_ncm`, `co_pais`, `co_ano`, `co_mes`
   - Limit de 100 países por query

3. **Connection Pooling**
   - Pool size: 25 conexões
   - Max idle: 5 conexões
   - Conn lifetime: 5 minutos

### Limites e Quotas

| Recurso | Limite |
|---------|--------|
| **Max NCM length** | 8 dígitos |
| **Max countries filter** | 50 países |
| **Max results** | 50 destinos |
| **Request timeout** | 30 segundos |
| **Request size** | 1MB |
| **Rate limit (free)** | 5 req/dia |
| **Rate limit (premium)** | Unlimited |

### Dicas de Performance

1. **Use cache quando possível**
   - Verifique o campo `cache_hit` na response
   - Consultas repetidas são instantâneas

2. **Filtre por países específicos**
   - Reduz processamento e tempo de resposta
   - Use `countries` array com 5-10 países

3. **Ajuste max_results**
   - Default é 10 (ótimo para UI)
   - Reduza para 5 se precisar de velocidade máxima
   - Aumente para 20-30 apenas se necessário

4. **Monitore rate limits**
   - Sempre verifique headers `X-RateLimit-*`
   - Implemente backoff exponencial
   - Considere upgrade para premium

---

## Roadmap

### ✅ Implementado (v1.0 - Semana 2)

- [x] Endpoint POST `/v1/simulator/destinations`
- [x] Algoritmo de scoring com 4 métricas
- [x] Estimativas de margem, custo logístico, tarifa e lead time
- [x] Rate limiting freemium (5/dia free, ilimitado premium)
- [x] Suporte a filtro por países
- [x] Metadata com processing time
- [x] Tratamento de erros completo
- [x] Testes unitários (98.3% coverage)

### 🚧 Próximas Features (Semana 3-4)

- [ ] Cache L3 com PostgreSQL Materialized Views
- [ ] Request Coalescing (deduplica requests simultâneas)
- [ ] Cache warming via CronJob
- [ ] Webhook de notificação quando novos dados disponíveis
- [ ] Endpoint GET `/v1/simulator/ncm/{ncm}/metadata`
- [ ] Histórico de simulações do usuário
- [ ] Export para PDF/Excel

### 🔮 Futuro (Backlog)

- [ ] Machine Learning para predição de tendências
- [ ] Integração com APIs de frete (real-time logistics cost)
- [ ] Dados de tarifas reais via API da Receita Federal
- [ ] Análise de competitividade por destino
- [ ] Recomendações de produtos similares (NCMs relacionados)
- [ ] Dashboard interativo com mapas

---

## Suporte

### Documentação Adicional

- [Guia de Arquitetura](../README.md)
- [Observability Stack](./OBSERVABILITY.md)
- [Idempotency Policy](./IDEMPOTENCY-POLICY.md)
- [Security & Secrets](./SECURITY-SECRETS.md)

### Contato

- **Issues:** [GitHub Issues](https://github.com/seu-org/bgc-app/issues)
- **Email:** suporte@brasilglobalconect.com
- **Slack:** #api-support

---

**Versão da API:** v1.0
**Última Atualização:** 2025-01-21
**Autor:** BGC Engineering Team
