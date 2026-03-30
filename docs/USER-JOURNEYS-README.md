# User Journeys Documentation - Navigation Guide

**Created:** 2026-01-20
**Owner:** BGC Product Management Team
**Status:** Active

---

## Quick Navigation

Este conjunto de documentos mapeia **TODAS as jornadas de usuário necessárias** para transformar o BGC de um MVP funcional em um produto completo e escalável.

**Escolha seu documento baseado na sua necessidade:**

| Preciso de... | Documento | Tempo de Leitura |
|---------------|-----------|------------------|
| **Visão executiva rápida** | [Executive Summary](#executive-summary) | 5 min |
| **Plano de ação imediato** | [30-Day Playbook](#30-day-playbook) | 10 min |
| **Referência rápida** | [Index](#index) | 3 min |
| **Timeline visual** | [Gantt Chart](#gantt-chart) | 5 min |
| **Detalhes completos** | [Complete Map](#complete-map) | 30 min |

---

## Documents Overview

### Executive Summary
**File:** `USER-JOURNEYS-EXECUTIVE-SUMMARY.md`
**Audience:** CEO, Leadership, Investors
**Purpose:** Tomada de decisão estratégica

**O que contém:**
- TL;DR (60 segundos)
- Investment summary (R$1.03M / ano 1)
- Risk assessment
- Go/No-Go decision framework
- Competitive landscape
- Key metrics (OKRs)
- Recommendations

**Quando usar:**
- Apresentar para stakeholders
- Aprovação de budget
- Alinhamento estratégico
- Investor pitch

**Link:** [USER-JOURNEYS-EXECUTIVE-SUMMARY.md](./USER-JOURNEYS-EXECUTIVE-SUMMARY.md)

---

### 30-Day Playbook
**File:** `USER-JOURNEYS-30DAY-PLAYBOOK.md`
**Audience:** Product Manager, Engineering Lead, Team
**Purpose:** Execução tática (próximos 30 dias)

**O que contém:**
- Week-by-week breakdown (5 weeks)
- Daily tasks para Dev 1 (Frontend) e Dev 2 (Backend)
- Checklists executáveis
- Standup format, retro template
- Risk mitigation por semana
- Budget breakdown (R$80k)
- Decision points (go/no-go)

**Quando usar:**
- Sprint planning (semana a semana)
- Daily standups (track progress)
- Retrospectives (identificar blockers)
- Budget tracking

**Link:** [USER-JOURNEYS-30DAY-PLAYBOOK.md](./USER-JOURNEYS-30DAY-PLAYBOOK.md)

---

### Index
**File:** `USER-JOURNEYS-INDEX.md`
**Audience:** Anyone (quick reference)
**Purpose:** Navegação rápida, status overview

**O que contém:**
- Status overview (completo/parcial/pendente)
- Jornadas por prioridade (P0/P1/P2/P3)
- Timeline visual (Fase 0 → 3)
- Matriz impacto vs esforço
- Dependencies críticas
- Metrics tracker (NSM impact)
- Effort summary (291 dias total)

**Quando usar:**
- Consulta rápida (qual jornada é P0?)
- Status check (quantas faltam?)
- Priorização (quick wins vs long bets)
- Dependencies check (o que bloqueia o quê?)

**Link:** [USER-JOURNEYS-INDEX.md](./USER-JOURNEYS-INDEX.md)

---

### Gantt Chart
**File:** `USER-JOURNEYS-GANTT.md`
**Audience:** PM, Engineering Lead, Project Managers
**Purpose:** Timeline visual, resource planning

**O que contém:**
- Gantt chart ASCII (Fase 0 → 3)
- Critical path analysis (65 dias de trabalho efetivo)
- Resource allocation (Dev 1 vs Dev 2 tasks)
- Milestone tracker (M1, M2, M3, M4)
- Dependencies matrix
- Velocity tracking (story points)
- Risk heatmap (timeline)
- Capacity planning (team growth)

**Quando usar:**
- Sprint planning (visualizar timeline)
- Resource allocation (quem faz o quê?)
- Dependency management (critical path)
- Risk tracking (quando riscos são altos?)

**Link:** [USER-JOURNEYS-GANTT.md](./USER-JOURNEYS-GANTT.md)

---

### Complete Map
**File:** `USER-JOURNEYS-COMPLETE-MAP.md`
**Audience:** PM, Product Designers, Engineers
**Purpose:** Especificação detalhada de todas jornadas

**O que contém:**
- **32 jornadas mapeadas** (5 Acquisition, 5 Activation, 6 Retention, 5 Revenue, 5 Referral, 6 Admin)
- Para cada jornada:
  - User Story (JTBD format)
  - Fluxo completo (step-by-step)
  - Páginas/rotas necessárias
  - Componentes/features
  - Success metrics
  - Effort estimate (dias)
  - Prioridade (P0-P3)
  - Dependências
- Roadmap visual (4 fases)
- RICE prioritization
- North Star Metric alignment

**Quando usar:**
- Design sprints (detalhar UX de jornada)
- Engineering planning (scope técnico)
- User research (validar hipóteses)
- Metrics definition (o que medir?)

**Link:** [USER-JOURNEYS-COMPLETE-MAP.md](./USER-JOURNEYS-COMPLETE-MAP.md)

---

## How to Use This Documentation

### For CEOs / Leadership
1. Start with **Executive Summary** (5 min read)
2. Review budget and ROI (R$1.03M investment, break-even mês 13)
3. Approve/reject Fase 0 (R$80k, 30 dias)
4. Set decision points (NPS >40, 3 conversões mínimas)

### For Product Managers
1. Read **Complete Map** (understand all 32 journeys)
2. Use **Index** for quick prioritization (P0 first)
3. Plan sprints with **30-Day Playbook** (week-by-week)
4. Track progress with **Gantt Chart** (visual timeline)
5. Report to stakeholders with **Executive Summary** (weekly updates)

### For Engineering Leads
1. Read **30-Day Playbook** (technical breakdown)
2. Review **Gantt Chart** (resource allocation: Dev 1 vs Dev 2)
3. Check **Complete Map** for specific jornada details (tech stack, APIs)
4. Use **Index** for dependencies (what blocks what?)
5. Plan sprints (assign tasks from playbook)

### For Designers
1. Read **Complete Map** (understand user flows)
2. Check **30-Day Playbook** Week 1 (mockups needed)
3. Create designs for:
   - Week 1: Auth (signup, login, forgot password)
   - Week 2: Onboarding (3-step wizard)
   - Week 3: Paywall (modal, pricing page)
   - Week 4: Checkout (3-step flow)

### For Marketing
1. Read **Executive Summary** (positioning, competitive landscape)
2. Check **30-Day Playbook** Week 5 (beta user recruitment)
3. Prepare materials:
   - Invitation emails (beta waitlist)
   - Landing page copy (value prop)
   - Social proof (testimonials from beta)

---

## Document Relationships

```
┌─────────────────────────────────────────────────────────┐
│              EXECUTIVE SUMMARY                          │
│  (Strategic decision-making, budget approval)           │
└─────────────────────┬───────────────────────────────────┘
                      │
         ┌────────────┴────────────┐
         │                         │
         ▼                         ▼
┌─────────────────┐       ┌─────────────────┐
│   30-DAY        │       │   GANTT CHART   │
│   PLAYBOOK      │◄─────►│   (Timeline)    │
│   (Execution)   │       │                 │
└────────┬────────┘       └────────┬────────┘
         │                         │
         │        ┌────────────────┘
         │        │
         ▼        ▼
┌─────────────────────────────────────────┐
│              INDEX                       │
│  (Quick reference, prioritization)      │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│         COMPLETE MAP                     │
│  (32 jornadas detalhadas)               │
│  - Acquisition (5)                      │
│  - Activation (5)                       │
│  - Retention (6)                        │
│  - Revenue (5)                          │
│  - Referral (5)                         │
│  - Admin (6)                            │
└─────────────────────────────────────────┘
```

---

## Key Metrics Summary

### Current State (2026-01-20)
- ✅ **1 jornada completa** (Simulator Anonymous - Epic 4)
- ❌ **31 jornadas pendentes**
- 🎯 **6 jornadas P0** (MVP Beta blockers)

### Investment Required
- **Fase 0 (MVP):** R$80k, 5 semanas, 2 devs
- **Fase 1 (V1):** R$150k, 6 semanas, 3 people
- **Fase 2 (V2):** R$300k, 12 semanas, 5 people
- **Fase 3 (Scale):** R$500k, 24 semanas, 8 people
- **Total Ano 1:** R$1.03M

### Expected Outcomes
- **Fase 0:** 50 beta users, 3 paid, NPS >40
- **Fase 1:** 500 users, R$10k MRR, 5k visits/mo
- **Fase 2:** 2k users, R$50k MRR, 2 Enterprise
- **Fase 3:** 5k users, R$100k MRR, NPS >60

---

## Prioritization Framework (RICE)

**Top 5 Jornadas (Highest RICE Score):**

1. **J-REV01 (Paywall)** - RICE 320
   - 6 dias, P0, Fase 0
   - Monetization enabler

2. **J-AC01 (Onboarding)** - RICE 300
   - 8 dias, P0, Fase 0
   - First-time user experience

3. **J-AC02 (Auth)** - RICE 225
   - 12 dias, P0, Fase 0
   - Foundation (blocks everything)

4. **J-A03 (Referral)** - RICE 180
   - 8 dias, P1, Fase 1
   - Viral growth

5. **J-R06 (Email Re-Engagement)** - RICE 180
   - 7 dias, P1, Fase 1
   - Retention mechanism

---

## Critical Dependencies

**Blocker:** J-AC02 (Auth) deve ser completo ANTES de:
- J-AC01 (Onboarding)
- J-AC03 (Simulator Logged)
- J-AC04 (Dashboard)
- J-R01 (Histórico)
- J-REV01 (Paywall)
- J-REV02 (Checkout)
- ... (todas features autenticadas)

**Sequência Mínima Viável (MVP):**
```
J-AC02 (Auth) → J-AC01 (Onboarding) → J-REV01 (Paywall) → J-REV02 (Checkout)
```

---

## Jornadas Faltantes (Backlog Futuro)

**Não mapeadas neste release (Q2-Q4 2026):**

**Marketplace & Matchmaking:**
- J-MATCH01: Buyer Discovery
- J-MATCH02: RFQ Flow
- J-MATCH03: Negociação In-Platform
- J-MATCH04: Deal Closure Tracking

**Operations:**
- J-OPS01: Document Automation
- J-OPS02: Siscomex Integration
- J-OPS03: Logistics Quote Comparison
- J-OPS04: Shipment Tracking

**Financial Services:**
- J-FIN01: FX Simulation
- J-FIN02: Trade Finance Application
- J-FIN03: Credit Scoring
- J-FIN04: Insurance Quotes

**Total:** +16 jornadas (Q2-Q4 roadmap)

---

## FAQ

### Q: Qual documento devo ler primeiro?
**A:** Depende do seu role:
- **CEO/Leadership:** Executive Summary
- **PM:** 30-Day Playbook (se começando agora) ou Complete Map (se planejando futuro)
- **Engineering:** 30-Day Playbook
- **Designer:** Complete Map (entender flows)

### Q: Quantas jornadas precisamos implementar para lançar beta?
**A:** 6 jornadas P0 (Fase 0):
1. J-AC02: Auth
2. J-AC01: Onboarding
3. J-AC03: Simulator Logged
4. J-REV01: Paywall
5. J-REV02: Checkout
6. (J-AC03 já está 50% completo)

### Q: Quanto tempo leva para produto completo?
**A:**
- **MVP Beta (Fase 0):** 5 semanas
- **V1 Public (Fase 1):** +6 semanas (total: 11 semanas)
- **V2 Platform (Fase 2):** +12 semanas (total: 23 semanas)
- **Scale (Fase 3):** +24 semanas (total: 47 semanas = ~12 meses)

### Q: Qual o custo total?
**A:** R$1.03M para ano 1 completo (todas 32 jornadas + 16 futuras)

### Q: Podemos pular alguma jornada P0?
**A:** Não. Todas as 6 são críticas:
- Sem Auth = sem usuários recorrentes
- Sem Onboarding = usuários perdidos
- Sem Paywall = sem monetização
- Sem Checkout = sem receita

### Q: E se o beta falhar (NPS <40)?
**A:** Decision framework:
- NPS 30-40: Ajustar features/pricing, repetir beta
- NPS <30: Pivot ou kill (post-mortem, reassess)

### Q: Como acompanhar progresso?
**A:**
- **Diário:** Daily standup (use playbook format)
- **Semanal:** Check gantt chart, update index status
- **Mensal:** Executive summary (report to CEO)

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0.0 | 2026-01-20 | Initial release - 32 jornadas mapeadas | BGC PM Team |
| | | | |

**Next Review:** 2026-02-01 (após Sprint 1 retro)

---

## Contact & Support

**Product Owner:** Rafael (BGC PM)
**Document Maintainer:** Product Management Team
**Slack Channel:** #product-roadmap
**Email:** product@brasilglobalconect.com

**For questions:**
- Strategic decisions → Executive Summary
- Execution details → 30-Day Playbook
- Technical specs → Complete Map
- Timeline queries → Gantt Chart
- Quick lookups → Index

---

## Related Documentation

**Product Context:**
- [product.context.md](./product.context.md) - Visão geral da plataforma
- [PRODUCT-ROADMAP.md](./PRODUCT-ROADMAP.md) - Epics e features
- [STRATEGIC-ROADMAP-E2E-INTEGRATED.md](./STRATEGIC-ROADMAP-E2E-INTEGRATED.md) - Roadmap de testes

**Technical:**
- [architecture_doc.md](./architecture_doc.md) - Arquitetura técnica
- [API-SIMULATOR.md](./API-SIMULATOR.md) - API do simulador (Epic 4)
- [SIMULATOR_ARCHITECTURE.md](../web-next/SIMULATOR_ARCHITECTURE.md) - Frontend architecture

**Operations:**
- [deployment_guide.md](./deployment_guide.md) - Deploy procedures
- [RUNBOOK.md](./RUNBOOK.md) - Troubleshooting
- [OBSERVABILITY.md](./OBSERVABILITY.md) - Monitoring & metrics

---

**Status:** Active Documentation
**Living Document:** Yes (update weekly after sprints)
**Last Updated:** 2026-01-20
