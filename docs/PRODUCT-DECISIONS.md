# Decisões de Produto - BGC Platform

Registro de decisões estratégicas de produto, trade-offs e justificativas para o Brasil Global Connect.

---

## Índice

- [Filosofia de Produto](#filosofia-de-produto)
- [Decisões Estratégicas](#decisões-estratégicas)
- [Decisões de Feature](#decisões-de-feature)
- [Decisões Técnicas com Impacto de Produto](#decisões-técnicas-com-impacto-de-produto)
- [Trade-offs e Aprendizados](#trade-offs-e-aprendizados)

---

## Filosofia de Produto

### North Star Metric

**Métrica Estrela-Guia:** Volume total de exportações facilitadas via plataforma (USD)

**Métricas de Suporte:**
- Número de SMEs ativas na plataforma
- Taxa de conversão de simulação → transação real
- NPS (Net Promoter Score) de exportadores
- Tempo médio para primeira exportação bem-sucedida

### Princípios de Produto

1. **Simplicity First**: SMEs não têm tempo nem expertise técnica. Cada feature deve ser usável por quem nunca exportou.
2. **Data-Driven Intelligence**: Decisões baseadas em dados reais (ComexStat, Siscomex). Zero achismo.
3. **Progressive Disclosure**: Começar simples (freemium), evoluir com o usuário (premium).
4. **Trust Through Transparency**: Mostrar sempre a fonte dos dados e lógica das recomendações.
5. **Speed Wins**: Performance é feature. Usuários esperam respostas em < 1s.

---

## Decisões Estratégicas

### DEC-001: Foco em SMEs Brasileiras

**Data:** 2025-01-15
**Contexto:** Poderíamos servir múltiplos mercados (importadores, empresas globais, traders)

**Decisão:** Foco inicial 100% em SMEs exportadoras brasileiras

**Justificativa (RICE):**
- **Reach:** 1.5M+ SMEs no Brasil, 50k+ exportadoras ativas (alto)
- **Impact:** Alto (mercado carente de ferramentas acessíveis)
- **Confidence:** Alta (validação com 12 exportadores reais)
- **Effort:** Médio (6 meses para MVP)
- **Score RICE:** (1500000 × 3 × 0.8) / 6 = 600,000

**Alternativas Rejeitadas:**
- ❌ Multi-sided marketplace (exportadores + importadores): Chicken-egg problem, 2x esforço
- ❌ Foco em grandes empresas: Já têm soluções enterprise, ciclo de venda longo

**Resultado Esperado:** 500 SMEs usando a plataforma em 6 meses (Q2 2025)

---

### DEC-002: Modelo Freemium com Rate Limiting

**Data:** 2025-11-22
**Contexto:** Como monetizar sem criar barreira de entrada?

**Decisão:** Freemium agressivo (5 simulações/dia grátis) + Premium ilimitado

**Justificativa (Jobs-to-be-Done):**
- **Job Principal:** "Preciso validar se vale a pena exportar antes de investir tempo/dinheiro"
- **Job Emocional:** "Não quero parecer incompetente para meu chefe ao sugerir um mercado ruim"
- **Job Social:** "Preciso de dados concretos para convencer sócios/investidores"

**Framework Aplicado:**
| Tier | Simulações/dia | Preço | Persona |
|------|----------------|-------|---------|
| Free | 5 | R$ 0 | Explorador (kick the tires) |
| Pro | Ilimitado | R$ 199/mês | SME ativa (1-10 SKUs) |
| Enterprise | Ilimitado + API | Customizado | Trader / Grande empresa (100+ SKUs) |

**Métricas de Validação:**
- Conversão free → pro: 3-5% (benchmark SaaS B2B)
- Churn rate < 5% ao mês (pro tier)
- Time-to-value: < 10 minutos (primeira simulação útil)

**Cost of Delay:** Alta (cada mês sem monetização = R$ 50k em receita potencial perdida)

**Resultado Esperado:** 30% dos usuários free batem o limite em 7 dias, 5% convertem para pro

---

### DEC-003: Algoritmo de Scoring Simplificado

**Data:** 2025-11-22
**Contexto:** Poderíamos usar ML complexo ou algoritmos mais sofisticados

**Decisão:** Média ponderada simples com 4 métricas (Market Size 40%, Growth 30%, Price 20%, Distance 10%)

**Justificativa:**
1. **Explicabilidade > Acurácia**: SMEs precisam entender POR QUE um mercado foi recomendado
2. **Time-to-market**: Algoritmo simples = MVP em 1 semana vs 2 meses para ML
3. **Data Availability**: Dados históricos limitados (ComexStat 2020-2024), insuficientes para ML robusto
4. **Validação Rápida**: Fácil de testar com exportadores reais

**Pesos Escolhidos (Baseado em Entrevistas com 8 Exportadores):**
- **Market Size (40%)**: "Quero mercados grandes, não nichos arriscados"
- **Growth Rate (30%)**: "Crescimento importa mais que tamanho absoluto"
- **Price (20%)**: "Preço alto = margem melhor"
- **Distance (10%)**: "Logística é problema, mas não deal-breaker"

**Alternativas Consideradas:**
- ❌ Machine Learning (Random Forest, XGBoost): Overengineering para MVP, black-box
- ❌ Score único sem pesos: Não reflete prioridades reais de SMEs
- ✅ **Escolhido**: Pesos configuráveis, possibilidade de A/B testing no futuro

**Validação Planejada:**
- A/B test com 3 variações de pesos (semana 4)
- Entrevistas qualitativas pós-simulação (15 usuários)
- Comparar recomendações vs decisões reais de exportadores

**Resultado Esperado:** Taxa de aceitação (usuário escolhe destino recomendado) > 60%

---

## Decisões de Feature

### DEC-004: Campos Calculados Automaticamente

**Data:** 2025-11-22
**Contexto:** Mostrar apenas score vs mostrar detalhes financeiros/logísticos

**Decisão:** Calcular e exibir 7 campos adicionais (margem, custo logístico, tarifa, lead time, etc.)

**Justificativa (User Research):**
- 9 de 10 exportadores entrevistados perguntaram: "Mas quanto vou ganhar de verdade?"
- Score sozinho é abstrato. Números concretos (USD, dias) geram ação.
- Transparência aumenta confiança (vs black-box)

**Heurísticas Implementadas:**
| Campo | Fórmula | Fonte |
|-------|---------|-------|
| EstimatedMarginPct | 15% (commodity) → 35% (alto valor) | Baseado em avg_price_per_kg |
| LogisticsCostUSD | Base cost + (distance × rate) - (volume × economy) | Tabela de custos por km |
| TariffRatePct | 8% (Americas) → 18% (outros) | Aproximação por região |
| LeadTimeDays | distance_km / 500km/dia | Velocidade média marítima |

**Disclaimers Adicionados:**
- "Estimativas baseadas em dados históricos. Consulte um despachante para valores exatos."
- Link para calculadora detalhada (futuro)

**Trade-off Aceito:**
- Estimativas podem ter erro de ±30%, mas são melhores que nada
- Preferível a pedir todos os dados ao usuário (abandono)

**Resultado Esperado:** Aumento de 40% na confiança na recomendação (medido via survey pós-simulação)

---

### DEC-005: Filtragem por Países Opcional

**Data:** 2025-11-22
**Contexto:** Sempre recomendar top N países vs permitir filtro customizado

**Decisão:** Campo `countries` opcional para filtrar resultados

**Justificativa:**
- **Caso de Uso Real:** "Já tenho contato na China, quero comparar EUA vs China vs Alemanha"
- **Progressive Disclosure:** Usuário iniciante ignora, avançado usa
- **Sem overhead**: Query SQL já eficiente com filtro

**UX Flow:**
1. Primeira simulação: Campo vazio, mostra top 10 globalmente
2. Tooltip: "Já tem países em mente? Filtre aqui"
3. Segunda simulação: 40% dos usuários filtram (hipótese)

**Resultado Esperado:** 30-40% dos usuários free usam filtro, 60%+ dos pro

---

### DEC-006: Max 50 Resultados por Request

**Data:** 2025-11-22
**Contexto:** Quantos destinos retornar? Ilimitado vs limite fixo

**Decisão:** Default 10, máximo 50 destinos

**Justificativa:**
- **Cognitive Load:** Usuário não consegue avaliar > 10 opções de uma vez
- **Performance:** 50 países = ~150ms query, 100+ = 300ms+ (timeout risk)
- **Paradox of Choice:** Mais opções = paralisia de decisão

**Default Escolhido:** 10 destinos
- Top 3 = "core focus" (80% dos usuários focam aqui)
- 4-10 = "exploratory" (20% exploram)
- 11-50 = "edge cases" (analistas, pesquisadores)

**Resultado Esperado:** 90% dos requests usam default (10), <5% pedem > 30

---

## Decisões Técnicas com Impacto de Produto

### DEC-007: Cache Multinível para Performance

**Data:** 2025-01-21
**Contexto:** Como garantir resposta < 200ms com queries complexas?

**Decisão:** Cache L1 (Ristretto in-memory) + L2 (Redis) + L3 (PostgreSQL Materialized Views)

**Impacto de Produto:**
- **Time-to-value**: Usuário vê resultados em 2-5ms (cache hit) vs 150ms (cold query)
- **UX Perception**: Plataforma "inteligente e rápida" vs "carregando..."
- **Conversão**: Cada 100ms de latência = -1% conversão (Amazon research)

**Trade-off:**
- Dados podem estar até 6h desatualizados (cache TTL)
- Aceitável: ComexStat atualiza mensalmente, não em tempo real

**Resultado Esperado:** 80%+ cache hit rate após 1 semana de uso, P95 latency < 200ms

---

### DEC-008: Rate Limiting por IP (Free Tier)

**Data:** 2025-11-22
**Contexto:** Como identificar usuários free sem forçar login?

**Decisão:** Rate limit por IP + user_id (se autenticado)

**Impacto de Produto:**
- **Friction Mínima**: Usuário testa sem criar conta
- **Conversão Futura**: Depois de bater limite, cria conta para continuar free (+ tracking)
- **Upgrade Path**: Free account → Pro subscription

**Risco Aceito:**
- IPs compartilhados (escritórios, NAT) podem atingir limite rápido
- Mitigação: Mensagem clara "Faça login para rastreamento individual"

**Resultado Esperado:** 15% dos usuários anônimos criam conta após bater limite

---

## Trade-offs e Aprendizados

### Aprendizado 001: Dados Reais > Dados Sintéticos

**Contexto:** Inicialmente usamos dados mock para desenvolvimento

**Aprendizado:**
- Dados sintéticos escondem edge cases reais (países sem dados, NCMs raros, crescimento negativo)
- Migration 0011 com dados reais (64 registros de ComexStat) revelou:
  - Necessidade de `COALESCE` para missing data
  - Filtro `market_size > 0` (alguns países têm volume 0)
  - Growth rate pode ser negativo (queda de mercado)

**Ação:** Sempre popular ambiente dev com subset de produção (10-100 registros reais por NCM)

---

### Aprendizado 002: Freemium Limits Precisam Ser Generosos

**Contexto:** Inicialmente consideramos 3 simulações/dia (tier free)

**Feedback Qualitativo (Entrevistas):**
- "3 simulações não dá pra testar nada, vou desistir"
- "Preciso de pelo menos 5 para comparar café, soja e carne (3 NCMs principais)"
- "Se bloquear muito cedo, não vou entender o valor"

**Decisão Final:** 5 simulações/dia
- Permite testar 5 NCMs diferentes OU 1 NCM 5x com filtros diferentes
- Freemium: generosidade gera confiança, não canibaliza premium

---

### Aprendizado 003: Explicabilidade > Acurácia para SMEs

**Contexto:** Poderíamos aumentar acurácia usando ML black-box

**Feedback Qualitativo:**
- "Não confio em algo que não entendo"
- "Como explico pro meu sócio que a China é melhor que os EUA?"
- "Prefiro um algoritmo 80% certo e transparente que 95% certo e opaco"

**Decisão:** Algoritmo simples + `recommendation_reason` textual
- Razão explica o "porquê" em linguagem natural
- Score decomponível (usuário vê peso de cada fator no futuro)

---

### Aprendizado 004: Performance É Feature, Não Infra

**Contexto:** Equipe queria lançar sem cache (MVP mais rápido)

**Impacto Calculado:**
- 150ms response → Bounce rate 10%
- 50ms response → Bounce rate 3%
- **7% delta = 70 usuários a mais em 1000 visitantes**

**Decisão:** Cache é parte do MVP, não "nice-to-have"
- Investir 2 dias em cache L1/L2 antes de lançar

---

## Próximas Decisões Pendentes

### PENDING-001: Pricing do Tier Premium

**Questão:** R$ 99/mês vs R$ 199/mês vs R$ 299/mês?

**Inputs Necessários:**
- Willingness-to-pay research (Van Westendorp PSM)
- Análise competitiva (Logcomex, Datawise, ComexDo)
- Unit economics (CAC, LTV)

**Deadline:** Semana 3 (antes do lançamento público)

---

### PENDING-002: Adicionar "Produtos Similares" (NCM Recommendation)

**Questão:** Recomendar NCMs similares ao que o usuário exporta?

**Trade-off:**
- 🟢 Aumenta descoberta, cross-sell
- 🔴 Complexidade técnica (classificação NCM hierárquica)
- 🔴 Risco de distrair do core job

**Framework Aplicado:** RICE pendente

**Deadline:** Q1 2025

---

### PENDING-003: Integração com Freight Forwarders

**Questão:** Integrar APIs de empresas de logística para custos reais?

**Trade-off:**
- 🟢 Custos precisos (vs heurísticas)
- 🟢 Potencial revenue share com parceiros
- 🔴 Dependência de terceiros (SLA, disponibilidade)
- 🔴 Complexidade de múltiplas integrações

**Alternativa:** Marketplace de freight forwarders (leads para parceiros)

**Deadline:** Q2 2025

---

## Métricas de Validação de Decisões

Todas as decisões são validadas contra:

### Métricas de Produto
- **Adoption Rate:** % de usuários que usam a feature
- **Retention:** % de usuários que voltam após 7/30 dias
- **Task Success Rate:** % de usuários que completam o job
- **Time-to-value:** Tempo até primeira simulação útil

### Métricas de Negócio
- **Conversion Rate:** Free → Pro (target: 3-5%)
- **Churn Rate:** < 5% ao mês (Pro tier)
- **NPS:** > 50 (promoters > detractors)
- **CAC Payback:** < 6 meses

### Métricas Técnicas
- **P95 Latency:** < 200ms
- **Availability:** > 99.5%
- **Error Rate:** < 0.1%
- **Cache Hit Rate:** > 80%

---

## Changelog de Decisões

**2025-11-22:**
- DEC-003: Algoritmo de scoring simplificado (aprovado)
- DEC-004: Campos calculados automáticos (aprovado)
- DEC-005: Filtro de países opcional (implementado)
- DEC-006: Max 50 resultados (implementado)

**2025-01-21:**
- DEC-007: Cache multinível (implementado)

**2025-01-15:**
- DEC-001: Foco em SMEs brasileiras (aprovado)
- DEC-002: Modelo freemium 5 req/dia (implementado)

---

**Versão:** 1.0
**Última Atualização:** 2025-11-22
**Responsável:** BGC Product Management Team
