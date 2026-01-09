# Próximos Passos - BGC Platform

Guia executivo dos próximos passos, priorizados por impacto e urgência, para a plataforma Brasil Global Connect.

**ATUALIZADO:** 05 de Janeiro de 2026
**Status:** ALERT - Ação Imediata Necessária

---

## Contexto Atual (05/01/2026)

### Status do Epic 4: Simulador de Destinos
- ✅ **85% Completo**: Backend, API, migrations, testes unitários implementados
- ✅ **Código Funcionando**: Endpoint validado com dados reais (performance 2-4ms)
- 🔴 **BLOQUEADO**: 53 arquivos não commitados há 44 dias
- 🔴 **BLOQUEADO**: Infraestrutura crítica não deployada (Redis, Integration Gateway, Observability)
- 🔴 **CRÍTICO**: PostgreSQL com 20 restarts (risco de perda de dados)

### Situação Crítica
- **Atraso:** 43 dias vs roadmap original
- **Código em Risco:** 2 semanas de trabalho sem backup em Git
- **Valor Entregue:** 0% (código pronto mas não deployado)
- **Cost of Delay:** R$ 75k em receita potencial perdida

### Próxima Meta (REVISADA)
**Desbloquear Produção e Deploy MVP até 10 de Janeiro de 2026**

---

## Prioridade P0 - Crítico (Hoje à Tarde + Segunda)

### P0.1: Deploy Redis no Kubernetes (BLOQUEANTE)
**Prazo:** Hoje à tarde (22/11/2025)
**Tempo Estimado:** 2 horas
**Responsável:** DevOps/Backend

**Tarefa:**
1. Aplicar `k8s/redis.yaml` no cluster k8s
2. Verificar PersistentVolumeClaim criado (2Gi)
3. Validar ConfigMap com configurações corretas
4. Testar conectividade do Integration Gateway com Redis
5. Validar métricas Prometheus de cache

**Validação de Sucesso:**
```bash
# Verificar pod rodando
kubectl get pods -n data | grep redis

# Testar conectividade
kubectl exec -it deployment/integration-gateway -n data -- redis-cli -h redis ping
# Esperado: PONG

# Verificar métricas
curl http://integration-gateway:8081/metrics | grep cache_hit
```

**Impacto:** Sem Redis, cache L2 não funciona → performance degradada (150ms vs 15ms)

**Riscos:**
- PVC pode falhar se storage class não configurado
- Mitigação: Verificar storage class padrão antes de aplicar

---

### P0.2: Popular Tabela `countries_metadata` (50 Países)
**Prazo:** Segunda-feira manhã (23/11/2025)
**Tempo Estimado:** 3 horas (desenvolvimento + execução)
**Responsável:** Backend

**Tarefa:**
1. Aplicar `k8s/jobs/populate-countries-job.yaml`
2. Monitorar execução do job via logs
3. Validar 50 países inseridos com metadados completos
4. Verificar flags, moedas e idiomas corretos

**Validação de Sucesso:**
```sql
-- Verificar países inseridos
SELECT COUNT(*) FROM countries_metadata;
-- Esperado: 50

-- Verificar campos preenchidos
SELECT code, name_pt, flag_emoji, currency, languages
FROM countries_metadata
WHERE flag_emoji IS NOT NULL
ORDER BY code;
```

**Script do Job:**
Usar `scripts/populate-countries/main.go` que:
- Busca dados via REST Countries API v3.1
- Calcula distância do Brasil via fórmula de Haversine
- Faz upsert com ON CONFLICT

**Impacto:** Sem 50 países, simulador retorna poucos resultados (apenas 10 países seed)

**Riscos:**
- REST Countries API pode estar offline
- Mitigação: Fallback para dados em arquivo JSON local

---

### P0.3: Testes E2E da API do Simulador
**Prazo:** Segunda-feira tarde (23/11/2025)
**Tempo Estimado:** 4 horas
**Responsável:** QA/Backend

**Cenários de Teste:**

#### Teste 1: Request Mínimo (Happy Path)
```bash
POST /v1/simulator/destinations
{
  "ncm": "17011400"
}

# Esperado:
# - 200 OK
# - destinations array com >= 6 países
# - score entre 0-10
# - todos campos preenchidos
```

#### Teste 2: Request com Filtro de Países
```bash
POST /v1/simulator/destinations
{
  "ncm": "17011400",
  "countries": ["US", "CN", "DE"]
}

# Esperado:
# - 200 OK
# - destinations array com APENAS US, CN, DE
# - score correto
```

#### Teste 3: Request com Volume
```bash
POST /v1/simulator/destinations
{
  "ncm": "26011200",
  "volume_kg": 5000,
  "max_results": 5
}

# Esperado:
# - 200 OK
# - exactly 5 destinos retornados
# - logistics_cost_usd calculado com volume
```

#### Teste 4: NCM Inválido
```bash
POST /v1/simulator/destinations
{
  "ncm": "12345"
}

# Esperado:
# - 400 Bad Request
# - error: "validation_error"
# - message: "NCM deve ter exatamente 8 dígitos"
```

#### Teste 5: NCM Não Encontrado
```bash
POST /v1/simulator/destinations
{
  "ncm": "99999999"
}

# Esperado:
# - 404 Not Found
# - error: "ncm_not_found"
```

#### Teste 6: Rate Limiting (Free Tier)
```bash
# Fazer 6 requests consecutivos
for i in {1..6}; do
  curl -X POST /v1/simulator/destinations -d '{"ncm":"17011400"}'
done

# Esperado:
# - Requests 1-5: 200 OK
# - Request 6: 429 Too Many Requests
# - Headers: X-RateLimit-Remaining: 0
```

**Validação de Sucesso:**
- Todos os 15 testes passando (3 NCMs × 5 variações)
- Coverage report > 80%

**Impacto:** Sem testes E2E, bugs podem chegar em produção

---

### P0.4: Commit e Merge do Simulador
**Prazo:** Segunda-feira noite (23/11/2025)
**Tempo Estimado:** 1 hora
**Responsável:** Tech Lead

**Tarefa:**
1. Revisar todos os arquivos novos e modificados
2. Escrever commit message descritivo
3. Push para branch `feature/security-credentials-management`
4. Code review com pelo menos 1 aprovação
5. Merge para `main`

**Commit Message Sugerido:**
```
feat(api): implement export destination simulator MVP

Complete implementation of the destination recommendation API with:

Backend:
- Domain layer with scoring algorithm (4 weighted metrics)
- Repository layer with optimized PostgreSQL queries
- Service layer with automatic estimates (margin, logistics, tariff, lead time)
- Error handling with custom business errors

API:
- POST /v1/simulator/destinations endpoint
- Freemium rate limiter middleware (5 req/day free, unlimited premium)
- Input validation (NCM 8 digits, volume > 0)
- Response with ranked destinations (score 0-10)

Database:
- Migration 0010: countries_metadata, comexstat_cache, simulator_recommendations
- Migration 0011: stg.exportacao schema with real ComexStat data
- 64 real records seeded (3 NCMs × multiple countries)
- 6 optimized indices created

Tests:
- Unit tests for handler and middleware (100% pass)
- Performance validated: 2-4ms per request
- Rate limiting validated

Documentation:
- docs/API-SIMULATOR.md (750+ lines)
- Swagger annotations in handler

Performance:
- P95 latency: 4ms (50x better than 200ms target)
- Real data from ComexStat 2020-2024
- Ready for production deployment

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Impacto:** Sem merge, código não vai para produção

---

## Prioridade P1 - Importante (Semana 3)

### P1.1: Frontend do Simulador (UI/UX)
**Prazo:** 24-28/11/2025
**Tempo Estimado:** 1 semana (40h)
**Responsável:** Frontend

**User Stories:**
- [ ] US-001: Input de NCM com autocomplete (futuramente)
- [ ] US-002: Filtro de países (multi-select dropdown)
- [ ] US-003: Card de destino com todas informações
- [ ] US-004: Score breakdown (gráfico radar)
- [ ] US-005: Modal de upgrade quando bate rate limit

**Componentes React:**
```typescript
// app/simulator/page.tsx
<SimulatorPage>
  <SimulatorForm onSubmit={handleSimulate} />
  <DestinationList destinations={results} />
  <UpgradeModal show={rateLimitHit} />
</SimulatorPage>

// components/SimulatorForm.tsx
<form>
  <NCMInput placeholder="Digite 8 dígitos" />
  <VolumeInput optional />
  <CountryFilter multiple />
  <SubmitButton loading={isLoading} />
</form>

// components/DestinationCard.tsx
<Card>
  <CountryHeader flag={🇺🇸} name="Estados Unidos" />
  <ScoreBadge score={8.5} rank={1} />
  <DemandIndicator level="Alto" />
  <FinancialMetrics margin={25%} logistics={$375} />
  <ScoreBreakdown weights={[40,30,20,10]} />
</Card>
```

**Design System:**
- Material Design 3 (MUI v7)
- Apple aesthetic (clean, minimalist)
- Mobile-first responsive

**Validação:**
- Figma mockups aprovados
- Usability testing com 3 usuários
- A/B test de layout (se tempo permitir)

---

### P1.2: Beta Privado com 20 Exportadores
**Prazo:** 01-05/12/2025
**Tempo Estimado:** 1 semana
**Responsável:** Product Manager

**Objetivos:**
1. Validar product-market fit
2. Coletar feedback qualitativo
3. Identificar bugs e edge cases
4. Calcular métricas de baseline (NPS, task success rate)

**Plano de Execução:**
1. Recrutar 20 exportadores SMEs (via LinkedIn, email, network)
2. Enviar acesso beta (whitelist de IPs ou tokens)
3. Agendar 1h de sessão por usuário (remoto, gravado)
4. Aplicar questionário pós-uso (SUS, NPS)
5. Compilar insights em relatório

**Perguntas de Pesquisa:**
- O algoritmo de scoring faz sentido?
- Os destinos recomendados são úteis?
- Falta alguma informação crítica?
- Você pagaria R$ 199/mês por isso?

**Métricas de Sucesso:**
- NPS > 40
- Task completion rate > 80%
- Acceptance rate > 60% (escolhem destino recomendado)
- Willingness-to-pay > R$ 150/mês

---

### P1.3: Ajustes Baseados em Feedback
**Prazo:** 05-10/12/2025
**Tempo Estimado:** 1 semana
**Responsável:** Product + Engineering

**Exemplos de Ajustes Esperados:**
- Alterar pesos do algoritmo (se feedback indicar)
- Adicionar campo "Competitividade" (se solicitado)
- Melhorar mensagens de erro
- Simplificar UX de filtros
- Adicionar tooltips explicativos

**Processo:**
1. Compilar top 5 issues de feedback
2. Priorizar via RICE
3. Implementar quick wins (< 1 dia)
4. Planejar features complexas para backlog

---

## Prioridade P2 - Desejável (Dezembro)

### P2.1: Dados Completos de ComexStat
**Prazo:** Dezembro 2025
**Tempo Estimado:** 2 semanas
**Responsável:** Data Engineer

**Objetivo:** Popular base com 1,000+ NCMs e dados históricos 2020-2024

**Plano:**
1. Criar script de ingestão em lote (Go ou Python)
2. Baixar exports completos do ComexStat
3. ETL para schema `stg.exportacao`
4. Validar integridade dos dados
5. Criar índices adicionais se necessário

**Métricas de Sucesso:**
- 1,000+ NCMs com dados
- 100k+ registros de exportação
- Data quality score > 95%

---

### P2.2: Cache L3 com Materialized Views
**Prazo:** Dezembro 2025
**Tempo Estimado:** 1 semana
**Responsável:** Backend + DBA

**Objetivo:** Implementar terceiro nível de cache em PostgreSQL

**Implementação:**
```sql
-- Materialized View para agregações pré-calculadas
CREATE MATERIALIZED VIEW cache.simulator_results AS
SELECT
  co_ncm,
  co_pais,
  SUM(vl_fob) as market_size_usd,
  AVG(vl_fob / kg_liquido) as avg_price_per_kg_usd,
  -- ... outros cálculos
FROM stg.exportacao
WHERE co_ano >= EXTRACT(YEAR FROM CURRENT_DATE) - 1
GROUP BY co_ncm, co_pais;

-- Refresh diário via CronJob
REFRESH MATERIALIZED VIEW CONCURRENTLY cache.simulator_results;
```

**Benefícios:**
- Cache hit rate aumenta para 90%+
- Reduz carga no banco primário
- Queries instantâneas (< 1ms)

---

### P2.3: Tier Premium & Sistema de Assinaturas
**Prazo:** Dezembro 2025 - Janeiro 2026
**Tempo Estimado:** 3 semanas
**Responsável:** Full Stack + Product

**Entregas:**
- Integração Stripe (pagamentos recorrentes)
- Auth JWT com roles (free, premium, enterprise)
- Dashboard de billing
- Upgrade flow
- Email automation (onboarding, invoices)

**Pricing Validado:**
- Realizar Van Westendorp PSM com 50 usuários
- Definir preço final (hipótese: R$ 199/mês)

---

## Prioridade P3 - Backlog (Q1 2026)

### Marketplace Beta
- Buyer onboarding
- Exporter-buyer matching
- RFQ flow
- First deal closed

### AI Export Assistant
- Integração GPT-4 / Claude
- RAG com base de conhecimento
- Chat interface

### Logistics Integration
- Freight forwarder APIs
- Real-time quotes
- Container tracking

---

## Riscos e Mitigações

### Risco 1: Redis Deployment Falha
**Probabilidade:** Baixa (10%)
**Impacto:** Alto (bloqueia cache L2)
**Mitigação:**
- Testar localmente com Docker Compose antes de k8s
- Ter fallback: Cache L1 funciona standalone
- Monitorar logs em tempo real durante deploy

---

### Risco 2: Job de Países Timeout
**Probabilidade:** Média (30%)
**Impacto:** Médio (poucos países disponíveis)
**Mitigação:**
- Aumentar timeout do job para 10 minutos
- Implementar retry em caso de falha parcial
- Ter dados de fallback em JSON local

---

### Risco 3: Feedback Beta Negativo
**Probabilidade:** Baixa (15%)
**Impacto:** Alto (produto-market fit em risco)
**Mitigação:**
- Pre-validar com 3 usuários antes de beta completo
- Estar preparado para pivotar algoritmo
- Ter roadmap de ajustes rápido (1 semana)

---

## Cronograma Visual

```
Semana 21/11   |  Semana 24/11   |  Semana 01/12   |  Semana 08/12
    (P0)       |     (P1)        |      (P1)       |     (P2)
      |        |       |         |        |        |       |
  22/11 Tarde  | Frontend UI     | Beta Privado    | Dados Completos
  Redis Deploy | Development     | 20 Exportadores | ComexStat
      |        |       |         |        |        |       |
  23/11 Manhã  |       |         |        |        |       |
  Popular      |       |         |        |        | Cache L3
  Países (Job) |       |         |        |        | Materialized
      |        |       |         |        |        | Views
  23/11 Tarde  |       |         | Ajustes |        |       |
  Testes E2E   |       |         | Feedback|        |       |
      |        |       |         |        |        |       |
  23/11 Noite  |       |         |        |        | Premium Tier
  Commit +     |       |         |        |        | Stripe
  Merge        |       |         |        |        |       |
      |        |       |         |        |        |       |
      ✅       |       ✅        |        ✅       |       ✅
```

---

## Definição de "Done" por Tarefa

### P0 (Deploy Redis)
- [ ] Pod `redis` rodando em k8s namespace `data`
- [ ] PVC criado e bound
- [ ] Integration Gateway conecta com sucesso
- [ ] Métricas Prometheus de cache disponíveis
- [ ] Health check retorna OK

### P0 (Popular Países)
- [ ] Job kubernetes completa com sucesso
- [ ] 50 países inseridos na tabela
- [ ] Todos campos (flag, currency, languages) preenchidos
- [ ] Query `SELECT COUNT(*) FROM countries_metadata` retorna 50

### P0 (Testes E2E)
- [ ] 15 testes implementados e passando
- [ ] Coverage report > 80%
- [ ] Performance validada (P95 < 200ms)
- [ ] Documentação de testes atualizada

### P0 (Commit & Merge)
- [ ] Commit message descritivo e seguindo convenção
- [ ] Code review aprovado
- [ ] CI/CD pipeline verde
- [ ] Merge para main sem conflitos
- [ ] Tag de versão criada (v0.4.0)

---

## Comunicação e Stakeholders

### Daily Standup
**Quando:** Todos os dias 9h
**Duração:** 15 minutos
**Formato:**
- O que fiz ontem?
- O que farei hoje?
- Bloqueios?

### Sprint Review
**Quando:** Segunda 23/11 (fim do dia)
**Duração:** 1 hora
**Audiência:** Product, Engineering, CEO
**Agenda:**
1. Demo do simulador funcionando (10 min)
2. Métricas atingidas vs targets (10 min)
3. Próximos passos (semana 3) (20 min)
4. Q&A (20 min)

### Retrospectiva
**Quando:** Terça 24/11
**Duração:** 1 hora
**Formato:** Start/Stop/Continue

---

## Contatos e Responsabilidades

| Área | Responsável | Email | Slack |
|------|-------------|-------|-------|
| Product Management | Rafael | rafael@bgc.dev | @rafael |
| Backend Engineering | TBD | backend@bgc.dev | @backend |
| Frontend Engineering | TBD | frontend@bgc.dev | @frontend |
| DevOps | TBD | devops@bgc.dev | @devops |
| QA | TBD | qa@bgc.dev | @qa |

---

## Recursos Adicionais

### Documentação
- [CHANGELOG.md](../CHANGELOG.md) - Histórico de mudanças
- [PRODUCT-ROADMAP.md](./PRODUCT-ROADMAP.md) - Roadmap completo
- [PRODUCT-DECISIONS.md](./PRODUCT-DECISIONS.md) - Decisões de produto
- [PRODUCT-METRICS.md](./PRODUCT-METRICS.md) - Métricas e KPIs
- [API-SIMULATOR.md](./API-SIMULATOR.md) - Documentação da API

### Links Úteis
- Grafana: http://localhost:3001
- Prometheus: http://localhost:9090
- Jaeger: http://localhost:16686
- API Local: http://localhost:8080
- Web Local: http://localhost:3000

---

**Versão:** 1.0
**Última Atualização:** 2025-11-22 (Manhã)
**Responsável:** BGC Product Management Team
**Próxima Atualização:** 2025-11-23 (Pós-deploy)
