# Product Roadmap - BGC Platform

Roadmap estratégico da plataforma Brasil Global Connect com foco em capacitar SMEs exportadoras brasileiras.

---

## Visão de Produto (12 meses)

**Missão:** Tornar exportação acessível e lucrativa para todas as SMEs brasileiras através de inteligência de dados e operações simplificadas.

**Visão 2025:** Ser a plataforma #1 de export intelligence para SMEs brasileiras, facilitando R$ 100M em exportações.

**North Star Metric:** Volume total de exportações facilitadas via plataforma (USD)

---

## Pilares Estratégicos

### 1. Export Intelligence (Decisões Informadas)
Ajudar SMEs a descobrir ONDE e O QUE exportar com base em dados reais.

### 2. Marketplace & Matching (Conexões Qualificadas)
Conectar exportadores brasileiros com compradores internacionais verificados.

### 3. Operational Enablement (Execução Simplificada)
Simplificar documentação, logística, pagamentos e compliance.

### 4. Financial Services (Acesso a Capital)
Facilitar financiamento, seguro e câmbio para exportações.

---

## Timeline e Marcos

```
Q4 2024         Q1 2025              Q2 2025              Q3 2025              Q4 2025
   |               |                    |                    |                    |
   ├─ Sprint 1     ├─ Export            ├─ Marketplace       ├─ Operations        ├─ Finance
   │  (Infra)      │  Intelligence      │  Beta              │  Automation        │  Integration
   │               │  MVP               │                    │                    │
   ├─ Sprint 2     ├─ Public            ├─ Buyer             ├─ AI Assistant      ├─ Credit
   │  (Obs)        │  Launch            │  Onboarding        │  (ChatGPT-like)    │  Scoring
   │               │                    │                    │                    │
   └─ Sprint 3     └─ Premium           └─ First Deal       └─ API Partners      └─ 100M Goal
      (API)           Tier                 Closed              Logistics/FX
```

---

## Q4 2024 - Foundation (CONCLUÍDO ✅)

### Sprint 1: Infrastructure & Data Pipeline
**Status:** 100% Completo
**Data:** 2025-01-10

**Entregas:**
- ✅ Kubernetes cluster (k3d) com PostgreSQL 16
- ✅ Serviço de ingestão CSV/XLSX (bgc-ingest)
- ✅ Materialized Views para agregação (rpt.*)
- ✅ API REST read-only (/market/size, /routes/compare)
- ✅ Clean Architecture em Go
- ✅ Docker Compose para desenvolvimento local

**Métricas:**
- Performance: P95 < 500ms para queries agregadas
- Data: 5M+ registros de exportação (2020-2024)
- Uptime: 99.9% em ambiente dev

---

### Sprint 2: Observability & Resilience
**Status:** 100% Completo
**Data:** 2025-01-15

**Entregas:**
- ✅ Prometheus + Grafana + Jaeger stack
- ✅ 11 métricas Prometheus customizadas
- ✅ Distributed tracing com OpenTelemetry
- ✅ Dashboards Grafana pré-configurados
- ✅ Alert rules para SLOs

**Métricas:**
- Observability: 100% cobertura de endpoints críticos
- MTTR: < 5 minutos para incidentes P2
- SLO: 99.5% availability

---

### Sprint 3: API Contracts & Governance
**Status:** 100% Completo
**Data:** 2025-01-21

**Entregas:**
- ✅ Integration Gateway (framework híbrido para APIs externas)
- ✅ JSON Schemas versionados (/v1/*)
- ✅ Middleware de idempotência
- ✅ Auth engine (mTLS, OAuth2, API Key)
- ✅ Network Policies (segmentação de rede)
- ✅ Sealed Secrets (gestão segura de credenciais)

**Métricas:**
- API Compliance: 100% endpoints validados com schemas
- Security: Zero credenciais em plain text
- Resiliência: Circuit breaker em todas integrações externas

---

## Q1 2025 - Export Intelligence MVP (EM ANDAMENTO 🚧)

### Epic 4: Simulador de Destinos de Exportação
**Status:** 85% Completo (Manhã de 22/11/2025)
**Data Prevista:** 2025-11-23 (SEGUNDA-FEIRA)

#### Semana 1-2: Backend & Database (85% COMPLETO)
**Progresso:**
- ✅ Domain layer (entities, errors, service)
- ✅ Repository layer (PostgreSQL queries otimizadas)
- ✅ API handler + middleware freemium
- ✅ Migration 0010 (countries_metadata, cache tables)
- ✅ Migration 0011 (schema ComexStat real com 64 registros)
- ✅ Algoritmo de scoring implementado
- ✅ Rate limiting (5 req/dia free, ilimitado premium)
- ✅ Testes unitários (handlers, middleware)
- ✅ Documentação API completa (docs/API-SIMULATOR.md)
- 🚧 **PENDENTE:** Deploy Redis k8s (cache L2)
- 🚧 **PENDENTE:** Job Kubernetes para popular 50 países
- 🚧 **PENDENTE:** Testes E2E completos

**Entregas Planejadas (Tarde 22/11 + 23/11):**
- [ ] Redis deployment no k8s (cache L2 distribuído)
- [ ] Kubernetes Job para popular countries_metadata (50 países)
- [ ] Testes E2E da API (3 NCMs × 5 variações)
- [ ] Commit final do simulador
- [ ] Merge para branch main

**Métricas Técnicas (Atingidas):**
- ✅ Performance: 2-4ms por request (com dados reais)
- ✅ Algoritmo: Score 0-10 com 4 métricas ponderadas
- ✅ Coverage: Testes unitários implementados

---

#### Semana 3: Frontend & UX (PLANEJADO)
**Data:** 2025-11-24 a 2025-11-28
**Status:** Não Iniciado

**User Stories:**
- [ ] US-001: Como exportador, quero inserir meu NCM e ver os melhores destinos ranqueados
- [ ] US-002: Como exportador, quero filtrar por países específicos que já conheço
- [ ] US-003: Como exportador, quero entender POR QUE um destino foi recomendado
- [ ] US-004: Como usuário free, quero saber quantas simulações restam hoje

**Entregas:**
- [ ] Página `/simulator` no web-next (Next.js 15)
- [ ] Componente SimulatorForm (input NCM + filtros)
- [ ] Componente DestinationCard (display recomendações)
- [ ] Visualização de score breakdown (gráfico radar ou barras)
- [ ] Modal de upgrade (quando bate rate limit)
- [ ] Loading states e error handling
- [ ] Responsivo mobile-first

**Métricas de UX:**
- Time-to-first-result: < 15 segundos (from landing)
- Task completion rate: > 80%
- Bounce rate: < 30%

---

#### Semana 4: Validação & Launch
**Data:** 2025-12-01 a 2025-12-05
**Status:** Não Iniciado

**Entregas:**
- [ ] Beta privado com 20 exportadores
- [ ] Coleta de feedback qualitativo (entrevistas)
- [ ] A/B test de pesos do algoritmo (3 variações)
- [ ] Ajustes baseados em feedback
- [ ] Documentação de onboarding
- [ ] Launch announcement (LinkedIn, email list)

**Métricas de Validação:**
- Beta NPS: > 40
- Acceptance rate: > 60% (escolhem destino recomendado)
- Free → Pro conversion: > 3%

---

### Epic 5: Dashboard de Market Intelligence (PLANEJADO)
**Status:** 0% Completo
**Data Prevista:** 2025-12-08 a 2025-12-19

**Objetivos:**
- Visualizar tendências de mercado (TAM/SAM/SOM por NCM)
- Comparar rotas comerciais (Brasil → Países)
- Análise temporal (crescimento, sazonalidade)

**User Stories:**
- [ ] US-010: Como exportador, quero ver o tamanho total do mercado para meu produto
- [ ] US-011: Como analista, quero comparar preços médios entre destinos
- [ ] US-012: Como exportador, quero identificar sazonalidade de demanda

**Entregas:**
- [ ] Endpoint GET `/v1/market/tam` (já existe, documentar)
- [ ] Endpoint GET `/v1/market/trends/{ncm}`
- [ ] Dashboard visual com gráficos (Recharts ou Chart.js)
- [ ] Filtros por ano, NCM, país
- [ ] Export para PDF/Excel

**Métricas:**
- Engagement: 40% dos usuários visitam dashboard
- Time on page: > 2 minutos

---

### Epic 6: Tier Premium & Monetização (PLANEJADO)
**Status:** 0% Completo
**Data Prevista:** 2025-12-15 a 2026-01-15

**Objetivos:**
- Validar pricing
- Implementar sistema de assinaturas
- Criar dashboard de billing

**Entregas:**
- [ ] Integração com Stripe (pagamentos recorrentes)
- [ ] Middleware de autenticação JWT
- [ ] Dashboard de billing (faturas, cartões)
- [ ] Email marketing automation (Mailchimp/SendGrid)
- [ ] Upgrade flow (free → pro)
- [ ] Pricing page otimizada

**Pricing Research:**
- Van Westendorp PSM com 50+ SMEs
- Competitive analysis (Logcomex, Datawise, ComexDo)

**Métricas de Monetização:**
- MRR (Monthly Recurring Revenue): R$ 10k ao final de Q1
- Churn rate: < 5%
- CAC payback: < 6 meses

---

## Q2 2025 - Marketplace Beta (PLANEJADO)

### Epic 7: Buyer Onboarding & Verification
**Status:** 0% Completo
**Data Prevista:** 2026-01-20 a 2026-02-15

**Objetivos:**
- Atrair compradores internacionais qualificados
- Verificar credibilidade (Dun & Bradstreet, background checks)
- Criar perfis de compradores

**Entregas:**
- [ ] Página de cadastro de compradores
- [ ] Integração com Dun & Bradstreet API (credit check)
- [ ] Perfil público de comprador (indústria, volume, produtos)
- [ ] Sistema de badges (verified, trusted, new)

**Métricas:**
- Compradores cadastrados: 50 (Q2)
- Verification rate: 80%

---

### Epic 8: Exporter-Buyer Matching Engine
**Status:** 0% Completo
**Data Prevista:** 2026-02-15 a 2026-03-15

**Objetivos:**
- Conectar exportadores e compradores automaticamente
- Ranking baseado em fit (NCM, volume, região)

**Entregas:**
- [ ] Algoritmo de matching (collaborative filtering)
- [ ] Página "Matches para Você"
- [ ] Sistema de interesse (like/pass)
- [ ] Chat in-platform (mensagens diretas)

**Métricas:**
- Match rate: > 30%
- First message rate: > 15%

---

### Epic 9: First Deal Closed
**Status:** 0% Completo
**Data Prevista:** 2026-03-15 a 2026-04-30

**Objetivos:**
- Facilitar primeira transação real via plataforma
- Aprender bottlenecks operacionais

**Entregas:**
- [ ] RFQ (Request for Quote) flow
- [ ] Cotação e negociação in-platform
- [ ] Integração com despachante parceiro
- [ ] Template de contrato internacional

**Métricas:**
- Deals closed: 3 (Q2)
- GMV (Gross Merchandise Value): USD 100k

---

## Q3 2025 - Operations Automation (PLANEJADO)

### Epic 10: Document Automation
**Status:** 0% Completo
**Data Prevista:** 2026-05-01 a 2026-06-15

**Objetivos:**
- Automatizar geração de documentos de exportação
- Integrar com Siscomex (DU-E, DI)

**Entregas:**
- [ ] Templates de Invoice, Packing List, Bill of Lading
- [ ] Integração com Siscomex API (OAuth2 mTLS)
- [ ] Upload e validação de documentos
- [ ] Dashboard de status de documentação

**Métricas:**
- Doc generation time: < 5 minutos
- Accuracy: > 95%

---

### Epic 11: AI Export Assistant (ChatGPT-like)
**Status:** 0% Completo
**Data Prevista:** 2026-06-15 a 2026-07-31

**Objetivos:**
- Responder perguntas sobre exportação em linguagem natural
- Guiar usuários através de processos complexos

**Entregas:**
- [ ] Integração com OpenAI GPT-4 ou Claude
- [ ] RAG (Retrieval-Augmented Generation) com base de conhecimento
- [ ] Chat interface
- [ ] Feedback loop (thumbs up/down)

**Métricas:**
- User satisfaction: > 80%
- Queries resolved: > 60%

---

### Epic 12: Logistics Integration
**Status:** 0% Completo
**Data Prevista:** 2026-08-01 a 2026-09-15

**Objetivos:**
- Integrar com freight forwarders para cotações reais
- Rastreamento de cargas

**Entregas:**
- [ ] Integração com 3+ freight forwarders (API)
- [ ] Cotação de frete em tempo real
- [ ] Rastreamento de containers
- [ ] Marketplace de logística (leilão reverso)

**Métricas:**
- Freight partners: 5
- Logistics cost reduction: 15%

---

## Q4 2025 - Financial Services (PLANEJADO)

### Epic 13: Foreign Exchange (FX)
**Status:** 0% Completo
**Data Prevista:** 2026-09-15 a 2026-10-15

**Objetivos:**
- Facilitar câmbio com taxas competitivas
- Hedge de risco cambial

**Entregas:**
- [ ] Integração com exchange partners (Remessa Online, Wise)
- [ ] Simulador de câmbio
- [ ] Alertas de taxa favorável
- [ ] Hedge automático (futuro)

**Métricas:**
- FX volume: USD 1M
- Spread vs benchmark: < 1%

---

### Epic 14: Trade Finance & Credit Scoring
**Status:** 0% Completo
**Data Prevista:** 2026-10-15 a 2026-11-30

**Objetivos:**
- Facilitar financiamento de exportação
- Credit scoring de exportadores e compradores

**Entregas:**
- [ ] Integração com bancos/fintechs (BNDES, Eximbank)
- [ ] Credit scoring model (ML)
- [ ] Aplicação de crédito in-platform
- [ ] Invoice financing

**Métricas:**
- Loans facilitated: R$ 5M
- Default rate: < 3%

---

### Epic 15: 100M Milestone
**Status:** 0% Completo
**Data Prevista:** 2026-12-31

**Objetivo:** Facilitar R$ 100M em exportações via plataforma

**Métricas de Sucesso:**
- GMV: R$ 100M
- Active exporters: 500
- Active buyers: 200
- NPS: > 60
- MRR: R$ 100k

---

## Backlog (Sem Data Definida)

### Ideação (Não Priorizado)

**Export Analytics Advanced:**
- [ ] Análise de competitividade (share of wallet)
- [ ] Predição de demanda com ML
- [ ] Recomendação de NCMs similares
- [ ] Alertas de oportunidades (mercados emergentes)

**Platform Features:**
- [ ] Mobile app (iOS/Android)
- [ ] Offline mode
- [ ] Multi-idioma (EN, ES, CN)
- [ ] White-label para associações de exportadores

**Integrações:**
- [ ] ERP integrations (SAP, TOTVS)
- [ ] Accounting (QuickBooks, Conta Azul)
- [ ] CRM (Salesforce, HubSpot)

**Compliance & Risk:**
- [ ] Sanctions screening (OFAC, UN)
- [ ] KYC automation
- [ ] Insurance marketplace
- [ ] Export credit insurance

---

## Critérios de Priorização

Usamos **RICE Framework** para priorizar features:

**Fórmula:** RICE Score = (Reach × Impact × Confidence) / Effort

| Feature | Reach | Impact | Confidence | Effort | RICE Score | Prioridade |
|---------|-------|--------|------------|--------|------------|------------|
| Simulador Destinos | 1000 | 3 (High) | 0.8 | 2 weeks | 1200 | P0 (MVP) |
| Dashboard Market | 800 | 2 (Med) | 0.7 | 2 weeks | 560 | P1 |
| Premium Tier | 500 | 3 (High) | 0.6 | 3 weeks | 300 | P1 |
| Buyer Matching | 300 | 3 (High) | 0.5 | 6 weeks | 75 | P2 |
| AI Assistant | 600 | 2 (Med) | 0.4 | 4 weeks | 120 | P2 |
| Mobile App | 400 | 1 (Low) | 0.5 | 8 weeks | 25 | P3 |

**Legendas:**
- **Reach:** Número de usuários impactados em 3 meses
- **Impact:** 1 (Low), 2 (Medium), 3 (High)
- **Confidence:** 0.0 a 1.0 (quanto temos certeza das estimativas)
- **Effort:** Semanas de desenvolvimento

---

## Riscos e Mitigações

### Risco 1: Low Adoption (Freemium não converte)
**Probabilidade:** Média
**Impacto:** Alto
**Mitigação:**
- Beta privado antes de lançamento público
- Entrevistas semanais com usuários
- Ajustar limites freemium baseado em dados

---

### Risco 2: Dados ComexStat Insuficientes
**Probabilidade:** Baixa
**Impacto:** Alto
**Mitigação:**
- Integração com múltiplas fontes (TradeMap, UN Comtrade)
- Partnerships com data providers
- Fallback para dados agregados

---

### Risco 3: Competição (Players Estabelecidos)
**Probabilidade:** Alta
**Impacto:** Médio
**Mitigação:**
- Foco em SMEs (vs enterprise)
- UX simplificada (vs complexa)
- Pricing agressivo (freemium generoso)

---

## Governança do Roadmap

**Cadência de Revisão:** Quinzenal (Sprint Planning)

**Stakeholders:**
- Product Manager (decisor final)
- Engineering Lead (feasibility)
- UX Designer (usability)
- CEO (strategic alignment)

**Processo de Mudança:**
- Proposta de mudança → RICE scoring → Revisão em Planning → Decisão → Comunicação

**Transparência:**
- Roadmap público em https://brasilglobalconect.com/roadmap
- Changelog atualizado semanalmente
- Release notes em cada deploy

---

## Definição de "Done" por Epic

**Critérios Gerais:**
- [ ] Código em produção (merge na main)
- [ ] Testes E2E passando (coverage > 80%)
- [ ] Documentação atualizada
- [ ] Métricas instrumentadas (Prometheus)
- [ ] Feedback de 5+ usuários beta
- [ ] Post-mortem escrito (aprendizados)

---

## Changelog do Roadmap

**2025-11-22:**
- Epic 4 atualizado: 85% completo, pendências documentadas
- Adicionadas métricas técnicas atingidas
- Detalhamento de entregas planejadas para tarde/segunda

**2025-01-21:**
- Sprints 1, 2, 3 marcados como completos
- Epic 4 iniciado

**2025-01-15:**
- Roadmap inicial publicado

---

**Versão:** 2.0
**Última Atualização:** 2025-11-22 (Manhã)
**Responsável:** BGC Product Management Team
**Próxima Revisão:** 2025-11-25 (Sprint Planning)
