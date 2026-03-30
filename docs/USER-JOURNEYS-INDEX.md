# Índice Visual de Jornadas de Usuário - BGC Platform

**Quick Reference Guide**
**Last Updated:** 2026-03-29

---

## Status Overview

| Categoria | Total | Completo | Parcial | Em Andamento | Não Iniciado |
|-----------|-------|----------|---------|--------------|--------------|
| **Acquisition (A)** | 5 | 0 | 0 | 0 | 5 |
| **Activation (AC)** | 5 | 1 | 0 | 1 | 3 |
| **Retention (R)** | 6 | 0 | 0 | 0 | 6 |
| **Revenue (REV)** | 5 | 0 | 0 | 0 | 5 |
| **Referral (REF)** | 5 | 0 | 0 | 0 | 5 |
| **Administrative (ADM)** | 6 | 0 | 0 | 0 | 6 |
| **TOTAL** | **32** | **1** | **0** | **1** | **30** |

**Progress:** 9% (1/32 completa + 1/32 em andamento 87.5%) — Atualizado em 2026-03-29
**J-AC02:** DONE (100%) — Frontend + Backend JWT completos. Sessao 2026-03-29.
**J-AC01:** Em Andamento (87.5%) — Days 1-7 DONE. Day 8 pendente (polish + analytics stubs + QA staging). Previsto para conclusao em 2026-04-06.

---

## Jornadas por Prioridade

### P0 - Critical (MVP Beta Blockers)
🔴 **6 jornadas** - Sem estas, produto não funciona

| ID | Nome | Status | Effort | Fase |
|----|------|--------|--------|------|
| J-AC02 | Sistema de Autenticação | ✅ DONE (100%) | 12d | 0 |
| J-AC01 | Onboarding First-Time User | 🚧 Em Andamento (87.5%) — Day 8 pendente | 8d | 0 |
| J-AC03 | Simulador Logged-In User | 🔨 | 6d | 0 |
| J-REV01 | Descoberta Premium (Paywall) | ❌ | 6d | 0 |
| J-REV02 | Checkout e Pagamento | ❌ | 12d | 0 |
| **TOTAL P0** | | | **44 dias** | |

### P1 - High (V1 Public Launch)
🟠 **10 jornadas** - Necessárias para crescimento

| ID | Nome | Status | Effort | Fase |
|----|------|--------|--------|------|
| J-A01 | SEO/Blog Acquisition | ❌ | 12d | 1 |
| J-A03 | Indicação Boca-a-Boca | ❌ | 6d | 1 |
| J-AC04 | Dashboard Intelligence | ❌ | 10d | 1 |
| J-AC05 | Alertas e Notificações | ❌ | 8d | 1 |
| J-R01 | Histórico e Comparação | ❌ | 7d | 1 |
| J-R06 | Email Re-Engagement | ❌ | 7d | 1 |
| J-REV03 | Billing Management | ❌ | 8d | 1 |
| J-REF01 | Programa de Referral | ❌ | 8d | 1 |
| J-ADM03 | Analytics BI | ❌ | 12d | 1 |
| J-ADM05 | Suporte ao Cliente | ❌ | 10d | 1 |
| **TOTAL P1** | | | **88 dias** | |

### P2 - Medium (V2 Platform Maturity)
🟡 **10 jornadas** - Expansão e integrações

| ID | Nome | Status | Effort | Fase |
|----|------|--------|--------|------|
| J-A02 | Campanhas Pagas | ❌ | 8d | 2 |
| J-A04 | Parcerias Associações | ❌ | 10d | 2 |
| J-R02 | Watchlist Destinos | ❌ | 6d | 2 |
| J-R03 | Exploração NCMs | ❌ | 8d | 2 |
| J-R04 | Análise Concorrentes | ❌ | 10d | 2 |
| J-REV04 | Plano Enterprise | ❌ | 10d | 2 |
| J-REF03 | Co-Marketing Parceiros | ❌ | 10d | 2 |
| J-ADM01 | Onboard Parceiros | ❌ | 10d | 2 |
| J-ADM02 | Moderação Conteúdo | ❌ | 8d | 2 |
| J-ADM04 | Config Sistema | ❌ | 8d | 2 |
| **TOTAL P2** | | | **88 dias** | |

### P3 - Low (Scale & Optimization)
🟢 **6 jornadas** - Nice-to-have, otimizações

| ID | Nome | Status | Effort | Fase |
|----|------|--------|--------|------|
| J-A05 | Webinars e Eventos | ❌ | 5d | 3 |
| J-R05 | Fórum Comunidade | ❌ | 15d | 3 |
| J-REV05 | Upsell Add-Ons | ❌ | 7d | 3 |
| J-REF02 | Integrações CRM | ❌ | 15d | 3 |
| J-REF04 | User-Gen Content | ❌ | 7d | 3 |
| J-REF05 | Programa Afiliados | ❌ | 12d | 3 |
| J-ADM06 | Data Quality | ❌ | 10d | 3 |
| **TOTAL P3** | | | **71 dias** | |

---

## Timeline Visual

```
┌─────────────────────────────────────────────────────────────────────┐
│                        2026 ROADMAP                                  │
└─────────────────────────────────────────────────────────────────────┘

JAN-FEV          MAR-ABR           MAI-JUN           JUL-DEZ
   │                │                 │                 │
   │ FASE 0         │ FASE 1          │ FASE 2          │ FASE 3
   │ MVP Beta       │ V1 Launch       │ V2 Platform     │ Scale
   │ (4 weeks)      │ (6 weeks)       │ (3 months)      │ (6 months)
   │                │                 │                 │
   ├─ J-AC02 (Auth)│                 │                 │
   ├─ J-AC01       │                 │                 │
   ├─ J-AC03       │                 │                 │
   ├─ J-REV01      │                 │                 │
   ├─ J-REV02      │                 │                 │
   │  [6 jornadas] │                 │                 │
   │                │                 │                 │
   │ 50 beta users  ├─ J-A01 (SEO)   │                 │
   │ 3 paid         ├─ J-A03 (Ref)   │                 │
   │ NPS >40        ├─ J-AC04         │                 │
   │                ├─ J-R01          │                 │
   │                ├─ J-R06          │                 │
   │                ├─ J-ADM03        │                 │
   │                ├─ J-ADM05        │                 │
   │                │  [10 jornadas]  │                 │
   │                │                 │                 │
   │                │ 500 users       ├─ J-A04 (Part)  │
   │                │ R$10k MRR       ├─ J-REV04 (Ent) │
   │                │ 5k visits/mo    ├─ J-ADM01        │
   │                │                 │  [10 jornadas]  │
   │                │                 │                 │
   │                │                 │ 2k users        ├─ J-R05 (Fórum)
   │                │                 │ R$50k MRR       ├─ J-REF05 (Afil)
   │                │                 │ 2 Enterprise    │  [6 jornadas]
   │                │                 │                 │
   │                │                 │                 │ 5k users
   │                │                 │                 │ R$100k MRR
   │                │                 │                 │ NPS >60
   ▼                ▼                 ▼                 ▼
```

---

## Matriz de Impacto vs. Esforço

```
HIGH IMPACT
     ↑
     │  ┌──────────────┐  ┌──────────────┐
     │  │ J-AC01       │  │ J-AC02       │
     │  │ Onboarding   │  │ Auth         │
     │  │ 8d, P0       │  │ 12d, P0      │
     │  └──────────────┘  └──────────────┘
     │
     │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
     │  │ J-REV01      │  │ J-A03        │  │ J-A01        │
     │  │ Paywall      │  │ Referral     │  │ SEO/Blog     │
     │  │ 6d, P0       │  │ 8d, P1       │  │ 12d, P1      │
     │  └──────────────┘  └──────────────┘  └──────────────┘
     │
     │  ┌──────────────┐  ┌──────────────┐
MED  │  │ J-R02        │  │ J-R04        │  ┌──────────────┐
     │  │ Watchlist    │  │ Concorrentes │  │ J-REV04      │
     │  │ 6d, P2       │  │ 10d, P2      │  │ Enterprise   │
     │  └──────────────┘  └──────────────┘  │ 10d, P2      │
     │                                       └──────────────┘
     │
     │                    ┌──────────────┐  ┌──────────────┐
LOW  │                    │ J-A05        │  │ J-R05        │  ┌─────────────┐
     │                    │ Eventos      │  │ Fórum        │  │ J-REF02     │
     │                    │ 5d, P3       │  │ 15d, P3      │  │ CRM Integr. │
     │                    └──────────────┘  └──────────────┘  │ 15d, P3     │
     │                                                         └─────────────┘
     └────────────────────────────────────────────────────────────────────→
           LOW EFFORT          MEDIUM EFFORT         HIGH EFFORT
           (5-7 dias)          (8-12 dias)          (13-15 dias)
```

**Quadrantes:**
- **🎯 Quick Wins (High Impact, Low Effort):** J-REV01, J-R02, J-A03
- **🚀 Major Bets (High Impact, High Effort):** J-AC02, J-A01, J-REV04
- **💡 Fill-Ins (Low Impact, Low Effort):** J-A05
- **⚠️ Avoid (Low Impact, High Effort):** J-R05, J-REF02 (só fazer se estratégico)

---

## Dependências Críticas

### Must-Have Primeiro (Blockers)
```
J-AC02 (Autenticação)
  └─ BLOQUEIA 15+ jornadas
     ├─ J-AC01 (Onboarding)
     ├─ J-AC03 (Simulator Logged)
     ├─ J-AC04 (Dashboard)
     ├─ J-R01 (Histórico)
     ├─ J-REV01 (Premium)
     └─ ... (todas features autenticadas)
```

### Sequência Lógica MVP
```
1. J-AC02 (Auth) ────────┐
                         ├──> 2. J-AC01 (Onboarding)
                         │
                         └──> 3. J-AC03 (Simulator)
                                      │
                                      ├──> 4. J-REV01 (Paywall)
                                      │           │
                                      │           └──> 5. J-REV02 (Checkout)
                                      │
                                      └──> 6. J-R01 (Histórico)
```

---

## Metrics Tracker (North Star)

**North Star Metric:** Time-to-First-Export-Match

### Jornadas por Impacto em NSM

| Rank | Jornada | NSM Impact | Reasoning |
|------|---------|------------|-----------|
| 1 | J-AC01 (Onboarding) | 🔴 Critical | Reduz time-to-first-value de 30min → 5min |
| 2 | J-AC03 (Simulator) | 🔴 Critical | Core value delivery |
| 3 | J-AC04 (Dashboard) | 🟠 High | Insights → decisão → ação |
| 4 | J-R01 (Histórico) | 🟠 High | Comparação → decisão mais rápida |
| 5 | J-REV01 (Paywall) | 🟡 Medium | Premium unlock advanced insights |
| 6 | J-A01 (SEO) | 🟡 Medium | Acquisition → mais usuários na jornada |

---

## Personas por Jornada

### PME Exportador (Primary Persona)
**Jornadas Dedicadas:** 24/32 (75%)

- **Acquisition:** J-A01, J-A02, J-A03, J-A04, J-A05
- **Activation:** J-AC01, J-AC02, J-AC03, J-AC04, J-AC05
- **Retention:** J-R01, J-R02, J-R03, J-R04, J-R05, J-R06
- **Revenue:** J-REV01, J-REV02, J-REV03, J-REV04, J-REV05
- **Referral:** J-REF01, J-REF04

### Parceiros (Secondary)
**Jornadas Dedicadas:** 4/32 (12%)

- J-A04 (Associações)
- J-REF03 (Co-Marketing)
- J-ADM01 (Onboarding)
- J-REF02 (Integrações CRM)

### Afiliados/Influencers (Tertiary)
**Jornadas Dedicadas:** 2/32 (6%)

- J-REF05 (Programa Afiliados)
- J-REF04 (Case Studies)

### Admin/Internal (Operations)
**Jornadas Dedicadas:** 6/32 (19%)

- J-ADM01, J-ADM02, J-ADM03, J-ADM04, J-ADM05, J-ADM06

---

## Effort Summary

| Fase | Jornadas | Effort Total | Prazo (1 dev) | Prazo (2 devs) |
|------|----------|--------------|---------------|----------------|
| **Fase 0 (MVP)** | 6 | 44 dias | 9 semanas | 5 semanas |
| **Fase 1 (V1)** | 10 | 88 dias | 18 semanas | 9 semanas |
| **Fase 2 (V2)** | 10 | 88 dias | 18 semanas | 9 semanas |
| **Fase 3 (Scale)** | 6 | 71 dias | 14 semanas | 7 semanas |
| **TOTAL** | **32** | **291 dias** | **58 semanas** | **29 semanas** |

**Com 2 desenvolvedores full-time:** ~7 meses para produto completo

---

## Jornadas Faltantes (Backlog Futuro)

**Marketplace & Matchmaking (Q2 2026):**
- J-MATCH01: Buyer Discovery
- J-MATCH02: RFQ Flow
- J-MATCH03: Negociação
- J-MATCH04: Deal Tracking

**Operations (Q3 2026):**
- J-OPS01: Document Automation
- J-OPS02: Siscomex Integration
- J-OPS03: Logistics Quotes
- J-OPS04: Shipment Tracking

**Financial Services (Q3-Q4 2026):**
- J-FIN01: FX Simulation
- J-FIN02: Trade Finance
- J-FIN03: Credit Scoring
- J-FIN04: Insurance Quotes

**Total Futuro:** +16 jornadas (não estimadas ainda)

---

## Next Actions (Semana 3/2026)

### 1. Validar Roadmap com Stakeholders
- [ ] Review com CEO (strategic alignment)
- [ ] Review com Engineering (feasibility check)
- [ ] Review com UX (design complexity)

### 2. J-AC01 (Onboarding First-Time User) — EM PLANEJAMENTO
- [x] Sprint Planning: J-AC01 como proxima sprint (J-AC02 desbloqueou)
- [x] Sprint Plan detalhado criado: `docs/SPRINT-PLAN-J-AC01.md`
- [x] Fluxo completo mapeado: signup → onboarding wizard (3 steps) → dashboard com simulacao pre-carregada
- [ ] Implementacao: Dia 1-8 (2026-03-28 a 2026-04-06)
- [ ] Validar endpoint `/v1/simulator/destinations` com ncm_chapter de 2 digitos (Dia 1)

### 3. Preparar Infraestrutura Proxima Fase
- [ ] Setup Stripe/iugu account (J-REV02, desbloqueado por J-AC02)
- [ ] Configure email service (SendGrid/SES) para transacionais de onboarding
- [ ] Setup analytics events (Mixpanel/Segment) para J-ADM03
- [ ] Criar Secret K8s `clerk-secrets` em ambiente de staging

---

**Document Status:** Living Index - atualizar conforme implementações
**Owner:** BGC Product Management Team
**Última Atualização:** 2026-03-28
**Related Docs:**
- [Mapeamento Completo](./USER-JOURNEYS-COMPLETE-MAP.md)
- [Product Roadmap](./PRODUCT-ROADMAP.md)
- [Product Context](./product.context.md)
- [Sprint Plan J-AC01 — Onboarding](./SPRINT-PLAN-J-AC01.md)
