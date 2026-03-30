# User Journeys - Executive Summary

**For:** CEO, Product Leadership, Engineering Leadership
**Date:** 2026-01-20
**Status:** Strategic Planning Document

---

## TL;DR (60 Seconds)

**The Gap:**
- ✅ **2 jornadas completas** (Simulator Anonymous + J-AC02 Autenticacao)
- ❌ **30 jornadas faltantes** para produto maduro
- 🎯 **5 jornadas criticas P0 restantes** bloqueiam MVP Beta

**The Ask:**
- **Fase 0 (MVP Beta):** ~32 dias dev restantes = ~7 semanas (1 dev) ou ~4 semanas (2 devs)
- **Budget:** ~R$64k restante para Fase 0 (J-AC02 concluido: 12d de 44d gastos)
- **Outcome:** 50 beta users, 3 paid conversions, validation de product-market fit

**The Bet:**
- Se Fase 0 validar (NPS >40, conversion >5%), investir em Fase 1 (V1 Public Launch)
- Se não validar, pivot ou kill (fail fast)

**Atualização 2026-03-29:** J-AC02 (Sistema de Autenticacao com Clerk) marcado como DONE (100%). Frontend + Backend JWT completos. Proximo P0: J-AC01 (Onboarding First-Time User).

---

## Strategic Context

### Current State (2026-03-29)

**O que temos:**
- Backend robusto (Go + PostgreSQL + K8s)
- Simulador funcionando (MVP Epic 4)
- 8% coverage de exports brasileiros (Chapter 17 - Acucar)
- Infraestrutura de observability (Prometheus + Grafana + Jaeger)
- Sistema de autenticacao completo (J-AC02 DONE): Clerk + JWT RS256 + JWKS cache + endpoints protegidos + 33 testes passando (20 TS + 13 Go)

**O que falta:**
- Onboarding (usuarios nao sabem como comecar) — J-AC01, PROXIMO P0
- Paywall (impossivel monetizar) — J-REV01
- Checkout (impossivel converter usuarios em clientes) — J-REV02
- Dashboard de insights (valor limitado a 1 simulacao) — J-AC04
- Retention mechanisms (email, historico, alertas)

**Consequência:**
- Plataforma é "one-shot tool" (usuário simula e vai embora)
- Sem conta = sem dados de comportamento = sem otimização
- Sem paywall = sem receita = não sustentável
- Sem retention = CAC não recupera = burn insustentável

---

## The Master Plan (12 Meses)

### Fase 0: MVP Beta (JAN-FEV 2026)
**Timeline:** 5 semanas
**Investment:** R$80k (2 devs full-time)
**Goal:** Validar product-market fit com 50 early adopters

**Jornadas Críticas (P0):**
1. J-AC02: Sistema de Autenticação (12 dias)
2. J-AC01: Onboarding First-Time User (8 dias)
3. J-AC03: Simulador Logged-In + Save (6 dias)
4. J-REV01: Paywall Soft (6 dias)
5. J-REV02: Checkout PIX/Card (12 dias)

**Success Criteria:**
- 50 signups (beta invite-only)
- 30% activation (fazem primeira simulação)
- 3 conversões premium (R$197/mês cada = R$591 MRR)
- NPS >40 (usuários recomendam)
- Time-to-first-value <5 min

**Decision Point:** Se atingir métricas, prosseguir para Fase 1. Se não, pivot.

---

### Fase 1: V1 Public Launch (MAR-ABR 2026)
**Timeline:** 6 semanas
**Investment:** R$150k (2 devs + 1 marketing)
**Goal:** Crescimento orgânico + retention + monetization

**Jornadas (P1):**
- Acquisition: SEO/Blog (J-A01), Referral (J-A03)
- Retention: Dashboard (J-AC04), Histórico (J-R01), Email (J-R06)
- Revenue: Billing (J-REV03)
- Operations: Analytics (J-ADM03), Support (J-ADM05)

**Success Criteria:**
- 500 signups totais
- 5k organic visits/month
- Referral: 30% dos signups
- MRR: R$10k (50 premium users × R$200)
- Churn: <8%

**Outcome:** Produto viável publicamente, crescimento orgânico iniciado

---

### Fase 2: V2 Platform Maturity (MAI-JUN 2026)
**Timeline:** 3 meses
**Investment:** R$300k (2 devs + 1 designer + 1 marketing + 1 sales)
**Goal:** Enterprise, parcerias, network effects

**Jornadas (P2):**
- Acquisition: Paid Ads (J-A02), Parcerias (J-A04)
- Retention: NCMs Relacionados (J-R03), Concorrentes (J-R04)
- Revenue: Enterprise (J-REV04), Add-Ons (J-REV05)
- Operations: Partner Onboarding (J-ADM01)

**Success Criteria:**
- 2k signups totais
- MRR: R$50k
- Enterprise deals: 2 (R$10k/mês cada)
- Partners onboarded: 20 (freight, despachantes)

**Outcome:** Produto escalável, revenue diversificado

---

### Fase 3: Scale & Optimization (JUL-DEZ 2026)
**Timeline:** 6 meses
**Investment:** R$500k (equipe completa + marketing budget)
**Goal:** 5k users, R$100k MRR, community, network lock-in

**Jornadas (P3):**
- Community: Fórum (J-R05)
- Integrations: CRM/ERP (J-REF02)
- Affiliates: Programa Afiliados (J-REF05)
- Content: Case Studies (J-REF04)

**Success Criteria:**
- 5k signups totais
- MRR: R$100k
- Community: 500 membros ativos
- Affiliate revenue: 20% de MRR
- NPS: >60

**Outcome:** Plataforma líder em export intelligence para SMEs

---

## Investment Summary

| Fase | Timeline | Team | Cost | MRR Target | ROI |
|------|----------|------|------|------------|-----|
| **0 (MVP)** | 5 weeks | 2 devs | R$80k | R$0.6k | -99% |
| **1 (V1)** | 6 weeks | 3 people | R$150k | R$10k | -93% |
| **2 (V2)** | 12 weeks | 5 people | R$300k | R$50k | -83% |
| **3 (Scale)** | 24 weeks | 8 people | R$500k | R$100k | -80% |
| **Total Ano 1** | 47 weeks | | **R$1.03M** | **R$100k/mês** | Break-even em Mês 13 |

**Assumptions:**
- Dev: R$8k/semana (R$16k/mês)
- Designer: R$6k/semana (R$12k/mês)
- Marketing: R$8k/semana (R$16k/mês)
- Sales: R$10k/semana (R$20k/mês)
- Infra/tools: R$5k/mês (AWS, Stripe, SaaS)

**Break-Even Analysis:**
- R$100k MRR × 12 = R$1.2M ARR
- Total investment ano 1: R$1.03M
- Break-even: Mês 13 (assumindo linear growth)

---

## Risk Assessment

### Critical Risks (High Impact, High Probability)

**1. Low Adoption (Freemium não converte)**
- **Probabilidade:** 40%
- **Impacto:** Kill product
- **Mitigação:**
  - Beta privado antes de public launch (Fase 0)
  - Entrevistas semanais com usuários (10/semana)
  - Ajustar pricing/features baseado em dados (A/B tests)
  - Kill switch: Se NPS <30 ou conversion <2% após Fase 0, pivot

**2. Dados Insuficientes (Coverage 8% → churn alto)**
- **Probabilidade:** 30%
- **Impacto:** Alto (usuários frustrated)
- **Mitigação:**
  - Priorizar expansão NCM coverage paralelamente (Epic 5: 28% → 60% → 100%)
  - Comunicação transparente: "Em breve: carnes, soja, frutas"
  - Oferecer análise manual (consultoria) para NCMs não-cobertos (upsell)

**3. Competição (Players estabelecidos)**
- **Probabilidade:** 60%
- **Impacto:** Médio (market share dilution)
- **Mitigação:**
  - Foco em SMEs (vs enterprise): UX simples, pricing agressivo
  - Freemium generoso (vs paywall hard)
  - Network effects via community/referral (lock-in)
  - Speed to market (6 meses MVP → Scale vs 18 meses big co)

---

## Decision Framework

### Go/No-Go Gates

**Gate 1 (End of Fase 0 - Semana 5):**
- ✅ NPS >40 → PROSSEGUIR Fase 1
- ⚠️ NPS 30-40 → ADJUST pricing/features, repeat beta
- ❌ NPS <30 → PIVOT ou KILL

**Gate 2 (End of Fase 1 - Semana 11):**
- ✅ MRR >R$8k + Churn <10% → PROSSEGUIR Fase 2
- ⚠️ MRR R$5-8k → EXTEND Fase 1 (optimize conversion)
- ❌ MRR <R$5k → REASSESS business model

**Gate 3 (End of Fase 2 - Semana 23):**
- ✅ MRR >R$40k + 1 Enterprise deal → PROSSEGUIR Fase 3
- ⚠️ MRR R$30-40k → EXTEND Fase 2 (close enterprise pipeline)
- ❌ MRR <R$30k → CAP investment, optimize existing

---

## Competitive Landscape

| Competitor | Focus | Pricing | Strengths | Weaknesses | BGC Advantage |
|------------|-------|---------|-----------|------------|---------------|
| **Logcomex** | Enterprise, analytics | R$2k/mês | Data depth, integrations | Complex UX, expensive | SME-focused, freemium |
| **Datawise** | Compliance, docs | R$1.5k/mês | Siscomex integration | Limited intelligence | Market insights, simulator |
| **ComexDo** | Education, community | R$500/mês | Content, courses | No actionable tools | Data-driven decisions |
| **TradeMap (ITC)** | Free, global | Grátis | Comprehensive, trusted | Generic, not Brazil-focused | Brazil-specific, actionable |

**BGC Positioning:**
- **Target:** SMEs (R$1M-R$50M revenue) vs Enterprise
- **Pricing:** R$197/mês vs R$1.5k-R$2k
- **UX:** Apple aesthetic, simple vs complex dashboards
- **Freemium:** Generous (5 sim/day free) vs paywall hard
- **Network:** Community + referral vs isolated tools

---

## Key Metrics Dashboard (OKRs)

### North Star Metric
**Time-to-First-Export-Match:** Target <30 dias (from signup to real export lead)

### Leading Indicators (Track Weekly)

**Acquisition:**
- Signups (weekly): Target +20/week by end of Fase 1
- Acquisition channels: Organic >50%, Referral >30%
- CAC: Target <R$150

**Activation:**
- Onboarding completion: >70%
- Time-to-first-simulation: <5 min
- Activation rate (signup → 1st sim): >60%

**Retention:**
- D7 retention: >40%
- D30 retention: >25%
- Churn rate (monthly): <8%

**Revenue:**
- MRR growth: +20%/month
- Free → Premium conversion: >5%
- Trial → Paid conversion: >25%
- ARPU (Average Revenue Per User): >R$50

**Referral:**
- K-factor (virality): >0.5
- Referral signups: >30% of total

---

## Recommendations (Product Leadership)

### Immediate Actions (This Week)

1. **Approve Fase 0 Budget (R$80k)**
   - Hire/contract 2 full-time devs for 5 weeks
   - Setup Stripe/iugu account (payment gateway)
   - Kickoff sprint planning (J-AC02: Auth)

2. **Recruit Beta Users (50 target)**
   - Outreach to APEX, ABIEC members
   - LinkedIn outreach (identify active exporters)
   - Offer incentive: 6 months premium grátis

3. **Design Mockups (Parallel Track)**
   - J-AC01: Onboarding wizard (3 steps)
   - J-REV01: Paywall modal (feature comparison)
   - J-REV02: Checkout flow (PIX + card)

### Strategic Bets

**Double Down On:**
- Freemium generosity (5 sim/day free → 10 if engagement high)
- Onboarding simplicity (Apple-like, 3 steps max)
- Community & referral (network effects = moat)

**De-Prioritize:**
- Mobile app (until Fase 3, desktop-first)
- Multi-language (English only after 1k users)
- AI Assistant (cool, but not NSM-critical)

**Kill If:**
- Fase 0 NPS <30 (no product-market fit)
- Fase 1 MRR growth <10%/month (no traction)
- Competitor launches identical product at R$99/mês (pricing war)

---

## Appendices

### Related Documents
1. [Mapeamento Completo de Jornadas](./USER-JOURNEYS-COMPLETE-MAP.md) - 32 jornadas detalhadas
2. [Índice Visual](./USER-JOURNEYS-INDEX.md) - Quick reference, matrizes
3. [Product Roadmap](./PRODUCT-ROADMAP.md) - Epics e features
4. [Product Context](./product.context.md) - Visão, estratégia, arquitetura

### Contacts
- **Product Owner:** Rafael (BGC PM)
- **Engineering Lead:** [TBD]
- **Design Lead:** [TBD]
- **CEO:** [TBD]

### Revision History
- 2026-03-29: J-AC02 marcado DONE (100%). Estado atual atualizado. Proximo P0: J-AC01.
- 2026-01-20: Initial version (32 jornadas mapeadas)

---

**Status:** AWAITING APPROVAL
**Decision Required By:** 2026-01-24 (Friday EOD)
**Owner:** BGC Product Management Team
