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

**Nota (2026-04-19):** O projeto é pessoal no momento; sem meta comercial ativa. Esta decisão de foco permanece válida como princípio de design.

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

**Nota (2026-04-19):** Monetização não é prioridade ativa agora (projeto pessoal). O mecanismo de freemium permanece implementado para quando/se houver abertura ao mercado.

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

### DEC-009: Auth Bypass em Dev (CLERK_JWKS_URL Vazio)

**Data:** 2026-03-29
**Contexto:** Como garantir que o ambiente de desenvolvimento não quebre ao rodar sem credenciais Clerk configuradas?

**Decisão:** Se `CLERK_JWKS_URL=""` (variável ausente ou vazia), o middleware de auth não é aplicado a nenhum endpoint. O servidor sobe normalmente sem Clerk.

**Justificativa:**
- **Developer Experience:** Novos contribuidores não precisam criar conta Clerk para rodar o projeto localmente
- **CI/CD Simples:** Pipelines de build e testes unitários não dependem de segredos externos
- **Segurança Preservada:** A variável é obrigatória em staging e produção via `secretKeyRef: clerk-secrets/jwks-url` no K8s; ausencia em prod = deploy falhado (Secret inexistente bloqueia o Pod)
- **Consistência com Pattern Existente:** Outros conectores externos (ComexStat, etc.) seguem o mesmo pattern de fallback

**Trade-off Aceito:**
- Em dev local sem Clerk, qualquer chamada aos endpoints protegidos retorna 200 (sem autenticacao). Aceitavel porque o banco de dev nao contem PII ou dados financeiros reais.

**Risco Mitigado:** Checklist de deploy exige `CLERK_JWKS_URL` populado antes de promover para staging. Documentado em `docs/SECURITY-SECRETS.md`.

---

### DEC-010: OptionalMiddleware no Simulator (Freemium Anonimo + Identificado)

**Data:** 2026-03-29
**Contexto:** O endpoint `/v1/simulator/destinations` deve ser publico (freemium anonimo), mas usuarios autenticados devem ter melhor rate limiting (por user_id, nao por IP).

**Decisao:** Usar `OptionalMiddleware()` no simulator: extrai JWT se presente, seta `user_id` no contexto, mas nao rejeita requisicoes sem token. O rate limiter freemium usa `user_id` se disponivel (ilimitado para autenticados com plano ativo, 5/dia por IP para anonimos).

**Justificativa (Jobs-to-be-Done):**
- **Job do Explorador Anonimo:** "Quero testar sem criar conta" — nao bloquear
- **Job do Usuario Autenticado Free:** "Ja tenho conta, quero rastreamento individual" — melhorar limite por identidade, nao por IP compartilhado
- **Job do Usuario Premium:** "Paguei, quero ilimitado" — desbloqueado via claim no JWT

**Impacto de Negocio:**
- Reduz abandono de usuarios anonimos (sem fricção de login forcado)
- Aumenta incentivo para criacao de conta (limite por IP e menos generoso que por user_id)
- Pavimenta caminho para J-REV01 (Paywall) sem reescrever o endpoint

**Alternativa Rejeitada:**
- Exigir JWT obrigatorio no simulator: eliminaria o freemium anonimo, principal mecanismo de aquisicao

---

### DEC-011: JWKS Cache com TTL 1h e Renovacao Thread-Safe

**Data:** 2026-03-29
**Contexto:** Cada validacao de JWT exige buscar as chaves publicas do Clerk (JWKS endpoint). Chamar o endpoint externo a cada request seria inaceitavel para performance e resiliencia.

**Decisao:** Cache in-process do JWKS com TTL de 1 hora, protegido por `sync.RWMutex`. Renovacao automatica na primeira requisicao apos expirar (lazy refresh). Sem dependencias pesadas (implementacao manual do parsing RSA).

**Justificativa:**
- **Performance:** JWKS e um conjunto de chaves publicas estaticas; o Clerk raramente rota chaves. TTL 1h e seguro e elimina latencia de rede por request.
- **Resiliencia:** Se o endpoint JWKS do Clerk ficar indisponivel, o cache ainda serve as validacoes durante o TTL ativo.
- **Thread Safety:** Gin processa requests em goroutines concorrentes; `sync.RWMutex` garante consistencia sem bloquear leituras simultaneas.
- **Zero Dependencias Extras:** Parsing manual de RSA via `crypto/rsa` + `encoding/base64` evita adicionar bibliotecas pesadas ao modulo Go.

**Trade-off:**
- Rotacao de chave no Clerk pode demorar ate 1h para propagar. Risco aceitavel: o Clerk notifica rotacoes com antecedencia e o TTL pode ser reduzido via variavel de ambiente se necessario.

**Resultado:** P95 de validacao JWT < 1ms (cache hit) vs ~80-150ms (chamada HTTP ao Clerk).

---

### DEC-012: K8s Secret `clerk-secrets` para JWKS URL

**Data:** 2026-03-29
**Contexto:** A URL do JWKS contem o identificador do tenant Clerk. Nao pode ficar em ConfigMap publico nem hardcoded em imagem Docker.

**Decisao:** `CLERK_JWKS_URL` e injetada via `secretKeyRef: clerk-secrets/jwks-url` no deployment K8s da API. Em docker-compose local, e lida de variavel de ambiente do host (`.env` local, nao commitado).

**Justificativa:**
- **Seguranca:** Sealed Secrets ou gestao de segredos K8s garante que o valor nao aparece em logs de deployment ou historico Git
- **Consistencia:** Mesmo pattern dos outros segredos do projeto (DB credentials, etc.) documentados em `docs/SECURITY-SECRETS.md`
- **Auditabilidade:** Mudancas de URL JWKS (ex: migracao de tenant Clerk) sao rastreadas via Git no manifesto K8s sem expor o valor

**Implicacao Operacional:** Se o Secret `clerk-secrets` nao existir no namespace, o Pod falha no `ImagePullBackOff`/`CreateContainerConfigError`. Isso e intencional — e o mecanismo de segurança que impede deploy sem auth configurada em staging/prod.

---

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

## Decisões Pendentes

### PENDING-001: Expansão de Cobertura NCM

**Questão:** Quais capítulos NCM priorizar depois do cap. 17 (Açúcares)?

**Proposta atual (ver ROADMAP.md):** cap. 02 (Carnes) e cap. 08 (Frutas) — Bloco 2, mai/2026.

**Trade-off:** Cobertura ampla vs profundidade de dados por capítulo.

---

### PENDING-002: Migração Clerk publicMetadata → PostgreSQL user_profiles (J-AC04)

**Questão:** Quando migrar o storage de perfil do exportador do Clerk para PostgreSQL?

**Proposta atual:** Bloco 4 (jun-jul/2026) — após o simulador estar completo.

**Impacto:** Necessário para queries agregadas de perfil e personalização avançada.

---

### PENDING-003: Pricing e Abertura ao Mercado

**Questão:** Se/quando abrir ao mercado, qual o modelo de monetização?

**Status:** Sem deadline. Revisitar quando houver decisão de abertura ao mercado.

**Referência histórica:** Pricing de R$199/mês (Pro) foi definido em DEC-002 como hipótese inicial — validar com pesquisa antes de ativar.

---

## Métricas de Validação de Decisões

Todas as decisões são validadas contra:

### Métricas de Produto
- **Adoption Rate:** % de usuários que usam a feature
- **Retention:** % de usuários que voltam após 7/30 dias
- **Task Success Rate:** % de usuários que completam o job
- **Time-to-value:** Tempo até primeira simulação útil

### Métricas de Negócio (quando/se houver abertura ao mercado)
- **Conversion Rate:** Free → Pro (referência: 3-5% SaaS B2B)
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

**2026-03-29:**
- DEC-009: Auth bypass em dev via `CLERK_JWKS_URL` vazio (aprovado e implementado)
- DEC-010: OptionalMiddleware no simulator freemium (aprovado e implementado)
- DEC-011: JWKS cache TTL 1h thread-safe (aprovado e implementado)
- DEC-012: K8s Secret `clerk-secrets` para JWKS URL (aprovado e implementado)

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

**Versão:** 1.2
**Última Atualização:** 2026-04-19
**Responsável:** BGC Product Management
**Contexto atual:** Projeto pessoal, dev solo. Metas comerciais são referências futuras, não compromissos ativos.
