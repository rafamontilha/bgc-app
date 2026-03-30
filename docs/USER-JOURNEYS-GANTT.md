# User Journeys - Gantt Chart Roadmap

**Visual Timeline:** JAN 2026 → DEZ 2026
**Last Updated:** 2026-01-20

---

## Legend

- ████ Completed
- ▓▓▓▓ In Progress
- ░░░░ Planned
- ⚠️ Critical Path
- 🔴 P0 - Critical
- 🟠 P1 - High
- 🟡 P2 - Medium
- 🟢 P3 - Low

---

## FASE 0: MVP BETA (Semanas 1-5)

**Timeline:** 2026-01-20 → 2026-02-24

```
SEMANA           1    2    3    4    5
                 │    │    │    │    │
J-AC02 Auth      ⚠️████████████░░░░░░     (12 dias) 🔴
                 │              │
J-AC01 Onboard   │    ⚠️░░░░░░████████   (8 dias)  🔴
                 │              │
J-AC03 Simulator │         ⚠️░░██████░   (6 dias)  🔴
                 │              │
J-REV01 Paywall  │           ⚠️░░░████   (6 dias)  🔴
                 │                   │
J-REV02 Checkout │              ⚠️░░░░████████████ (12 dias) 🔴
                 │                                │
                 └────────────────────────────────┘
                           44 dias total
```

**Dependencies:**
- J-AC01, J-AC03, J-REV01 DEPENDEM de J-AC02 (Auth)
- J-REV02 DEPENDE de J-REV01 (Paywall)

**Critical Path:** J-AC02 → J-AC01 → J-REV01 → J-REV02

---

## FASE 1: V1 PUBLIC LAUNCH (Semanas 6-11)

**Timeline:** 2026-02-25 → 2026-04-07

```
SEMANA              6    7    8    9    10   11
                    │    │    │    │    │    │
J-A01 SEO/Blog      ████████████░░░░░░           (12d) 🟠
                    │         │
J-A03 Referral      │    ████████░░░░             (8d)  🟠
                    │         │
J-AC04 Dashboard    ████████████████░░░░          (10d) 🟠
                    │              │
J-AC05 Alertas      │         ████████░░░░        (8d)  🟠
                    │              │
J-R01 Histórico     │    ░░░░██████░░             (7d)  🟠
                    │         │
J-R06 Email         │         ░░░░██████          (7d)  🟠
                    │              │
J-REV03 Billing     │    ░░░░████████             (8d)  🟠
                    │         │
J-ADM03 Analytics   ████████████████████░░░░      (12d) 🟠
                    │                   │
J-ADM05 Support     │         ████████████████    (10d) 🟠
                    │                            │
                    └────────────────────────────┘
                              88 dias total
```

**Parallel Tracks:**
- Track 1 (Frontend): J-A01, J-AC04, J-R01
- Track 2 (Backend): J-A03, J-AC05, J-R06, J-REV03
- Track 3 (Operations): J-ADM03, J-ADM05

---

## FASE 2: V2 PLATFORM MATURITY (Semanas 12-23)

**Timeline:** 2026-04-08 → 2026-06-30

```
SEMANA           12   13   14   15   16   17   18   19   20   21   22   23
                 │    │    │    │    │    │    │    │    │    │    │    │
J-A02 Paid Ads   ████████░░░░                                                (8d)  🟡
                 │
J-A04 Parcerias  ████████████░░░░                                            (10d) 🟡
                 │         │
J-R02 Watchlist  │    ████████░░                                             (6d)  🟡
                 │         │
J-R03 NCMs       │    ░░░░████████░░                                         (8d)  🟡
                 │              │
J-R04 Concorr.   │         ████████████░░░░                                  (10d) 🟡
                 │                   │
J-REV04 Enterpr. │              ████████████░░░░                             (10d) 🟡
                 │                        │
J-REV05 Add-Ons  │                   ░░░░██████░                             (7d)  🟡
                 │                             │
J-REF03 CoMkt    │                        ████████████░░░░                   (10d) 🟡
                 │                                  │
J-ADM01 Parceiro │                             ████████████░░░░              (10d) 🟡
                 │                                       │
J-ADM02 Moderação│                                  ░░░░████████             (8d)  🟡
                 │                                            │
J-ADM04 Config   │                                       ░░░░████████        (8d)  🟡
                 │                                                         │
                 └─────────────────────────────────────────────────────────┘
                                     88 dias total
```

**Focus:**
- Meses 1-2 (Sem 12-19): Acquisition + Retention features
- Meses 2-3 (Sem 20-23): Revenue expansion + Operations

---

## FASE 3: SCALE & OPTIMIZATION (Semanas 24-47)

**Timeline:** 2026-07-01 → 2026-12-31

```
SEMANA           24   26   28   30   32   34   36   38   40   42   44   46
                 │    │    │    │    │    │    │    │    │    │    │    │
J-A05 Eventos    ████░                                                         (5d)  🟢
                 │
J-R05 Fórum      ████████████████░░░░░░░░                                      (15d) 🟢
                 │              │
J-REV05 Add-Ons  │         ░░░░██████░░                                        (7d)  🟢
                 │                   │
J-REF02 CRM      │              ████████████████░░░░░░░░                       (15d) 🟢
                 │                        │
J-REF04 UGC      │                   ░░░░██████░░                              (7d)  🟢
                 │                             │
J-REF05 Afiliado │                        ████████████████                     (12d) 🟢
                 │                                  │
J-ADM06 DataQual │                             ████████████░░░░                (10d) 🟢
                 │                                            │
                 │    [Continuous Optimization & A/B Tests]  │
                 │                                            │
                 └────────────────────────────────────────────┘
                              71 dias total (spread across 24 weeks)
```

**Strategy:**
- Implementação espaçada (não rush)
- Foco em otimização e A/B testing
- Preparação para Ano 2 (Marketplace, Operations)

---

## CRITICAL PATH ANALYSIS

### Longest Path (MVP → V1 → V2)

```
Start (Week 1)
    │
    ├─ J-AC02 Auth (12d) ─────────────────────┐
    │                                          │
    │  [BLOCKER: Todas features autenticadas]  │
    │                                          │
    └──────────────────────────────────────────┴─ J-AC01 Onboarding (8d)
                                                         │
                                                         ├─ J-REV01 Paywall (6d)
                                                         │         │
                                                         │         └─ J-REV02 Checkout (12d)
                                                         │                   │
                                                         │                   └─ FASE 0 COMPLETE (Week 5)
                                                         │
                                                         ├─ J-AC04 Dashboard (10d)
                                                         │         │
                                                         │         └─ J-R01 Histórico (7d)
                                                         │                   │
                                                         │                   └─ FASE 1 COMPLETE (Week 11)
                                                         │
                                                         └─ J-REV04 Enterprise (10d)
                                                                   │
                                                                   └─ FASE 2 COMPLETE (Week 23)
```

**Total Critical Path:** 65 dias (13 semanas de trabalho efetivo)

---

## RESOURCE ALLOCATION

### Developer 1 (Frontend Focus)

```
FASE 0: J-AC01 (Onboarding UI) + J-REV01 (Paywall Modal) + J-REV02 (Checkout)
FASE 1: J-A01 (Blog) + J-AC04 (Dashboard) + J-R01 (Histórico)
FASE 2: J-R02 (Watchlist) + J-R03 (NCMs) + J-R04 (Concorrentes)
FASE 3: J-R05 (Fórum) + J-REF04 (UGC)
```

### Developer 2 (Backend Focus)

```
FASE 0: J-AC02 (Auth Engine) + J-AC03 (Simulator Integration) + J-REV02 (Payment)
FASE 1: J-A03 (Referral) + J-AC05 (Alertas) + J-R06 (Email) + J-REV03 (Billing)
FASE 2: J-A04 (Parcerias) + J-REV04 (Enterprise) + J-ADM01 (Onboard Parceiros)
FASE 3: J-REF02 (CRM Integrations) + J-REF05 (Afiliados) + J-ADM06 (Data Quality)
```

### Designer (UX/UI)

```
FASE 0: Mockups (Onboarding, Paywall, Checkout)
FASE 1: Design System refinement, Dashboard, Email templates
FASE 2: Partner portal, Enterprise features
FASE 3: Community design, Case study templates
```

### Marketing

```
FASE 1: Blog content (20 articles), SEO optimization, Email campaigns
FASE 2: Paid ads setup (Google, LinkedIn), Partnership outreach
FASE 3: Webinars, Events, Affiliate program management
```

---

## MILESTONE TRACKER

| Milestone | Target Date | Jornadas | Success Criteria |
|-----------|-------------|----------|------------------|
| **M1: MVP Beta Launch** | 2026-02-24 (W5) | 6 (P0) | 50 signups, 3 paid, NPS >40 |
| **M2: Public Launch** | 2026-04-07 (W11) | 10 (P1) | 500 users, R$10k MRR, 5k visits/mo |
| **M3: Platform Maturity** | 2026-06-30 (W23) | 10 (P2) | 2k users, R$50k MRR, 2 Enterprise |
| **M4: Year-End Target** | 2026-12-31 (W47) | 6 (P3) | 5k users, R$100k MRR, NPS >60 |

---

## DEPENDENCIES MATRIX

```
                 AC02 AC01 AC03 REV01 REV02 A01 AC04 R01 R06 ADM03
                 ───┬────┬────┬─────┬─────┬───┬────┬───┬───┬─────
J-AC01 Onboard      │ X  │    │     │     │   │    │   │   │
J-AC03 Simulator    │ X  │    │     │     │   │    │   │   │
J-REV01 Paywall     │ X  │ X  │  X  │     │   │    │   │   │
J-REV02 Checkout    │ X  │    │     │  X  │   │    │   │   │
J-AC04 Dashboard    │ X  │ X  │  X  │     │   │    │   │   │
J-R01 Histórico     │ X  │    │  X  │     │   │    │ X  │   │
J-R06 Email         │ X  │    │     │     │   │    │    │   │
J-ADM03 Analytics   │    │    │     │     │   │    │    │   │
J-ADM05 Support     │ X  │    │     │     │   │    │    │   │
J-REV03 Billing     │ X  │    │     │     │  X│    │    │   │

X = Depende de (bloqueado até completar)
```

---

## VELOCITY TRACKING (To Be Updated)

### Sprint Velocity (Story Points per Week)

```
PLANNED:
Week 1-2:   ██████████████████████ (22 SP) ← J-AC02
Week 3:     ████████████████ (16 SP) ← J-AC01
Week 4:     ████████████ (12 SP) ← J-AC03 + J-REV01
Week 5:     ██████████████████████ (22 SP) ← J-REV02

ACTUAL:
Week 1-2:   [TBD after sprint]
Week 3:     [TBD]
...
```

**Story Points Scale:**
- 2 SP = 1 dia simples
- 5 SP = 2-3 dias médio
- 8 SP = 4-5 dias complexo
- 13 SP = 1 semana+ (break down)

---

## RISK HEATMAP (Timeline)

```
RISK LEVEL:  🔴 High   🟡 Medium   🟢 Low

                 W1   W5   W11  W23  W47
                 │    │    │    │    │
Technical Debt   🟢   🟡   🟡   🟡   🔴   ← Acumula ao longo do tempo
Auth Security    🔴   🟢   🟢   🟢   🟢   ← Crítico no início, resolve W2
Payment Fraud    🔴   🔴   🟡   🟢   🟢   ← Crítico até billing maduro
Data Quality     🟡   🟡   🟡   🟡   🟢   ← Melhora com automação (J-ADM06)
Churn Risk       🟢   🟡   🟡   🔴   🟡   ← Pico em W23 (trial expirations)
Competition      🟡   🟡   🟡   🟡   🔴   ← Aumenta conforme crescemos
```

---

## CAPACITY PLANNING

### Team Growth Over Time

```
PHASE        WEEKS    TEAM SIZE    MONTHLY COST
─────────────────────────────────────────────────
Fase 0       1-5      2 devs       R$80k total
Fase 1       6-11     3 people     R$25k/month
Fase 2       12-23    5 people     R$50k/month
Fase 3       24-47    8 people     R$80k/month
```

### Workload Distribution (Person-Weeks)

```
                 Q1   Q2   Q3   Q4
                 │    │    │    │
Frontend Dev     ████ ████ ███  ██   (20 person-weeks)
Backend Dev      ████ ████ ████ ███  (24 person-weeks)
Designer         ██   ███  ██   █    (10 person-weeks)
Marketing        ░    ███  ████ ████ (15 person-weeks)
Sales/CS         ░    ░    ███  ████ (10 person-weeks)
```

---

## NEXT ACTIONS (This Week)

### Sprint 1 Kickoff (Week 1: Jan 20-26)

**Monday (Jan 20):**
- [x] Sprint Planning: J-AC02 breakdown into tasks
- [ ] Setup dev environments (local + staging)
- [ ] Design review: Auth UI mockups

**Tuesday-Thursday (Jan 21-23):**
- [ ] Dev: JWT auth engine (backend)
- [ ] Dev: OAuth2 Google integration
- [ ] Dev: Signup/Login UI (frontend)
- [ ] QA: Security audit (penetration test)

**Friday (Jan 24):**
- [ ] Sprint Review: Demo working auth
- [ ] Retrospective: Ajustes de processo
- [ ] Sprint 2 Planning: J-AC01 (Onboarding)

---

**Document Owner:** BGC Product Management
**Update Frequency:** Weekly (every Friday after sprint review)
**Related Docs:**
- [Executive Summary](./USER-JOURNEYS-EXECUTIVE-SUMMARY.md)
- [Complete Map](./USER-JOURNEYS-COMPLETE-MAP.md)
- [Index](./USER-JOURNEYS-INDEX.md)
