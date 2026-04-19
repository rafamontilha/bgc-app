# User Journeys - Flow Diagram

**Visual Representation:** Como as jornadas se conectam
**Last Updated:** 2026-01-20

---

## Full User Journey Map (AARRR)

```
┌───────────────────────────────────────────────────────────────────────────┐
│                          ACQUISITION (A)                                   │
│                     Como usuários chegam à plataforma                      │
└───────────────────────────────────────────────────────────────────────────┘
                                      │
         ┌────────────┬───────────────┼───────────────┬────────────┐
         │            │               │               │            │
         ▼            ▼               ▼               ▼            ▼
    ┌────────┐  ┌─────────┐    ┌──────────┐    ┌─────────┐  ┌─────────┐
    │J-A01   │  │J-A02    │    │J-A03     │    │J-A04    │  │J-A05    │
    │SEO     │  │Paid Ads │    │Referral  │    │Parcerias│  │Eventos  │
    │Blog    │  │Google   │    │Boca-Boca │    │APEX     │  │Feiras   │
    │P1, 12d │  │LinkedIn │    │P1, 6d    │    │ABIEC    │  │P3, 5d   │
    │        │  │P2, 8d   │    │          │    │P2, 10d  │  │         │
    └────┬───┘  └────┬────┘    └────┬─────┘    └────┬────┘  └────┬────┘
         │           │              │               │            │
         └───────────┴──────────────┴───────────────┴────────────┘
                                   │
                                   ▼
         ┌─────────────────────────────────────────────────────────┐
         │           HOMEPAGE (/) - Anonymous User                 │
         │                                                          │
         │   CTA: "Simular Destinos Grátis" ou "Criar Conta"      │
         └─────────────────────────────────────────────────────────┘
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
                    ▼                             ▼
         ┌────────────────────┐       ┌────────────────────┐
         │  Anonymous Route   │       │   Signup Route     │
         │  /simulator        │       │   /signup          │
         │  (5 sim/day limit) │       │   (J-AC02)         │
         └─────────┬──────────┘       └──────────┬─────────┘
                   │                             │
                   │  Após 3-5 simulações        │
                   │  ou tentativa de save       │
                   └──────────────┬──────────────┘
                                  │
                                  ▼

┌───────────────────────────────────────────────────────────────────────────┐
│                          ACTIVATION (AC)                                   │
│                    Primeira experiência de valor                           │
└───────────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
         ┌─────────────────────────────────────────────────────────┐
         │              J-AC02: Authentication                      │
         │              P0, 12 dias (BLOCKER)                       │
         │                                                          │
         │   Signup: Email + Password ou Google OAuth              │
         │   Email verification → Redirect to Onboarding           │
         └─────────────────────────────────────────────────────────┘
                                  │
                                  ▼
         ┌─────────────────────────────────────────────────────────┐
         │          J-AC01: Onboarding (First-Time UX)             │
         │          P0, 8 dias                                      │
         │                                                          │
         │   Step 1/3: Qual produto? (NCM autocomplete)            │
         │   Step 2/3: Volume mensal? (kg input)                   │
         │   Step 3/3: Regiões de interesse? (multi-select)        │
         │                                                          │
         │   → Save preferences → Pre-fill Dashboard               │
         └─────────────────────────────────────────────────────────┘
                                  │
                                  ▼
         ┌─────────────────────────────────────────────────────────┐
         │              /dashboard (Logged-In Home)                │
         │                                                          │
         │   ┌──────────────┐  ┌──────────────┐                   │
         │   │ Card 1:      │  │ Card 2:      │                   │
         │   │ Seu Mercado  │  │ Tendências   │                   │
         │   │ (TAM/SAM)    │  │ (Price Trend)│                   │
         │   └──────────────┘  └──────────────┘                   │
         │                                                          │
         │   ┌──────────────┐  ┌──────────────┐                   │
         │   │ Card 3:      │  │ Card 4:      │                   │
         │   │ Destinos     │  │ Próximas     │                   │
         │   │ Recomendados │  │ Ações (TODOs)│                   │
         │   └──────────────┘  └──────────────┘                   │
         └─────────────────────────────────────────────────────────┘
                                  │
                    ┌─────────────┼─────────────┐
                    │             │             │
                    ▼             ▼             ▼
         ┌────────────────┐ ┌──────────┐ ┌───────────────┐
         │J-AC03          │ │J-AC04    │ │J-AC05         │
         │Simulator       │ │Dashboard │ │Alertas        │
         │Logged + Save   │ │Analytics │ │Notificações   │
         │P0, 6d          │ │P1, 10d   │ │P1, 8d         │
         └────────┬───────┘ └────┬─────┘ └───────┬───────┘
                  │              │                │
                  └──────────────┴────────────────┘
                                 │
                                 ▼

┌───────────────────────────────────────────────────────────────────────────┐
│                          RETENTION (R)                                     │
│                     Usuários retornam e engajam                            │
└───────────────────────────────────────────────────────────────────────────┘
                                 │
         ┌───────────┬───────────┼───────────┬───────────┬───────────┐
         │           │           │           │           │           │
         ▼           ▼           ▼           ▼           ▼           ▼
    ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐
    │J-R01   │  │J-R02   │  │J-R03   │  │J-R04   │  │J-R05   │  │J-R06   │
    │Histórico│ │Watchlist│ │NCMs    │  │Concorr.│  │Fórum   │  │Email   │
    │Compare │  │Favoritos│ │Related │  │Analysis│  │Community│ │Re-Eng. │
    │P1, 7d  │  │P2, 6d  │  │P2, 8d  │  │P2, 10d │  │P3, 15d │  │P1, 7d  │
    └────┬───┘  └────┬───┘  └────┬───┘  └────┬───┘  └────┬───┘  └────┬───┘
         │           │           │           │           │           │
         │   [User retorna D7, D30, D90]                            │
         │           │           │           │           │           │
         └───────────┴───────────┴───────────┴───────────┴───────────┘
                                 │
                     [Após 5-10 simulações]
                                 │
                                 ▼

┌───────────────────────────────────────────────────────────────────────────┐
│                          REVENUE (REV)                                     │
│                         Monetização                                        │
└───────────────────────────────────────────────────────────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
         ▼                       ▼                       ▼
    ┌────────────┐        ┌────────────┐        ┌────────────┐
    │J-REV01     │        │J-REV02     │        │J-REV03     │
    │Paywall Soft│        │Checkout    │        │Billing     │
    │Premium     │   ──►  │PIX/Card    │   ──►  │Management  │
    │Discovery   │        │Payment     │        │Invoices    │
    │P0, 6d      │        │P0, 12d     │        │P1, 8d      │
    └────────────┘        └────────────┘        └────────────┘
         │                       │                       │
         │  [User clicks        │ [Payment success]     │
         │   "Upgrade"]         │                       │
         │                      │                       │
         └──────────────────────┴───────────────────────┘
                                 │
                    ┌────────────┼────────────┐
                    │                         │
                    ▼                         ▼
         ┌────────────────────┐    ┌────────────────────┐
         │J-REV04             │    │J-REV05             │
         │Enterprise Tier     │    │Add-Ons             │
         │Custom Pricing      │    │Relatórios Custom   │
         │White-Label         │    │API Access          │
         │P2, 10d             │    │P3, 7d              │
         └────────────────────┘    └────────────────────┘
                                 │
                     [User is now Premium]
                                 │
                                 ▼

┌───────────────────────────────────────────────────────────────────────────┐
│                          REFERRAL (REF)                                    │
│                    Network effects e viral growth                          │
└───────────────────────────────────────────────────────────────────────────┘
                                 │
         ┌───────────┬───────────┼───────────┬───────────┐
         │           │           │           │           │
         ▼           ▼           ▼           ▼           ▼
    ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐
    │J-REF01 │  │J-REF02 │  │J-REF03 │  │J-REF04 │  │J-REF05 │
    │Referral│  │CRM     │  │Co-Mkt  │  │UGC     │  │Afiliados│
    │Program │  │Integr. │  │Parceiros│ │Case    │  │Programa│
    │P1, 8d  │  │P3, 15d │  │P2, 10d │  │Studies │  │P3, 12d │
    │        │  │        │  │        │  │P3, 7d  │  │        │
    └────┬───┘  └────┬───┘  └────┬───┘  └────┬───┘  └────┬───┘
         │           │           │           │           │
         │  [User shares → New user signs up via ref link]        │
         │           │           │           │           │
         └───────────┴───────────┴───────────┴───────────┘
                                 │
                    [Virality K-factor >0.5]
                                 │
                                 ▼
                         [Back to ACQUISITION]


┌───────────────────────────────────────────────────────────────────────────┐
│                       ADMINISTRATIVE (ADM)                                 │
│                    Operações internas e suporte                            │
└───────────────────────────────────────────────────────────────────────────┘
                                 │
         ┌───────────┬───────────┼───────────┬───────────┬───────────┐
         │           │           │           │           │           │
         ▼           ▼           ▼           ▼           ▼           ▼
    ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐
    │J-ADM01 │  │J-ADM02 │  │J-ADM03 │  │J-ADM04 │  │J-ADM05 │  │J-ADM06 │
    │Onboard │  │Moderação│ │Analytics│ │Config  │  │Support │  │Data    │
    │Parceiros│ │Conteúdo│  │BI      │  │Feature │  │Tickets │  │Quality │
    │P2, 10d │  │P2, 8d  │  │P1, 12d │  │Flags   │  │P1, 10d │  │P3, 10d │
    │        │  │        │  │        │  │P2, 8d  │  │        │  │        │
    └────────┘  └────────┘  └────────┘  └────────┘  └────────┘  └────────┘
         │           │           │           │           │           │
         └───────────┴───────────┴───────────┴───────────┴───────────┘
                                 │
                    [Suporte ao ciclo de vida do usuário]
```

---

## Critical Path (Minimum Viable Product)

```
START
  │
  ▼
┌──────────────────────────────────────────────────────────────┐
│ J-AC02: Authentication System (12 dias) 🔴 P0               │
│ ───────────────────────────────────────────                 │
│ • JWT engine                                                 │
│ • OAuth2 Google                                              │
│ • Email verification                                         │
│ • Password reset                                             │
│                                                              │
│ BLOQUEIO: Sem auth, nenhuma jornada autenticada funciona    │
└──────────────────────────────────────────────────────────────┘
  │
  ▼
┌──────────────────────────────────────────────────────────────┐
│ J-AC01: Onboarding First-Time User (8 dias) 🔴 P0           │
│ ───────────────────────────────────────────                 │
│ • 3-step wizard (NCM, Volume, Regions)                       │
│ • Pre-fill dashboard                                         │
│ • Skip option                                                │
│                                                              │
│ DEPENDE DE: J-AC02 (Auth)                                    │
└──────────────────────────────────────────────────────────────┘
  │
  ├──────────────────────┐
  │                      │
  ▼                      ▼
┌──────────────────┐  ┌──────────────────────────────────────┐
│ J-AC03: Simulator│  │ J-REV01: Paywall Soft (6 dias) 🔴 P0│
│ Logged-In + Save │  │ ───────────────────────────────────  │
│ (6 dias) 🔴 P0   │  │ • Feature comparison modal           │
│ ──────────────── │  │ • Pricing page                       │
│ • Auth check     │  │ • Trial activation (14 days)         │
│ • Save simulation│  │ • Upgrade banner                     │
│ • History        │  │                                      │
│                  │  │ DEPENDE DE: J-AC02, J-AC01           │
│ DEPENDE DE:      │  └──────────────────────────────────────┘
│ J-AC02           │           │
└──────────────────┘           │
                               ▼
                     ┌──────────────────────────────────────┐
                     │ J-REV02: Checkout (12 dias) 🔴 P0    │
                     │ ───────────────────────────────────  │
                     │ • Plan selection (Mensal/Anual)      │
                     │ • Billing info                       │
                     │ • Stripe integration (PIX, Card)     │
                     │ • Subscription management            │
                     │ • Invoice/NFSe                       │
                     │                                      │
                     │ DEPENDE DE: J-REV01                  │
                     └──────────────────────────────────────┘
                               │
                               ▼
                     ┌──────────────────────────────────────┐
                     │    MVP BETA COMPLETE (Week 5)        │
                     │    ═══════════════════════           │
                     │    50 signups, 3 paid, NPS >40       │
                     └──────────────────────────────────────┘
                               │
                               ▼
                          GO/NO-GO DECISION
```

**Total Critical Path:** 44 dias (9 semanas 1 dev, 5 semanas 2 devs)

---

## User Flow Example: Signup → First Paid Conversion

```
┌─────────────────────────────────────────────────────────────────────┐
│                      USER: Maria (SME Exporter)                      │
│                   Company: Açúcar Doce Ltda (Sugar)                  │
└─────────────────────────────────────────────────────────────────────┘
                                │
                                │ Day 1: Discovery
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Google Search: "para qual país exportar      │
         │ açúcar brasileiro"                           │
         │                                              │
         │ → Finds BGC blog article (J-A01)            │
         └──────────────────────────────────────────────┘
                                │
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Lands on Homepage (/)                        │
         │                                              │
         │ CTA: "Simular Destinos Grátis"              │
         └──────────────────────────────────────────────┘
                                │
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Anonymous Simulation (/simulator)            │
         │                                              │
         │ Input: NCM 17011400, 10,000 kg              │
         │ Result: Top 5 destinos (US, CN, DE, IN, GB) │
         │                                              │
         │ → Impressed! Wants to save results          │
         └──────────────────────────────────────────────┘
                                │
                                │ Modal appears: "Create account to save"
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Signup (/signup) - J-AC02                    │
         │                                              │
         │ Input: maria@acucardoce.com.br              │
         │        Password, Company name                │
         │                                              │
         │ → Email verification sent                    │
         └──────────────────────────────────────────────┘
                                │
                                │ Clicks verification link
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Onboarding (/onboarding) - J-AC01            │
         │                                              │
         │ Step 1: Produto → Açúcar (NCM 17011400)     │
         │ Step 2: Volume → 10,000 kg/mês              │
         │ Step 3: Regiões → Américas, Ásia            │
         │                                              │
         │ → Preferences saved                          │
         └──────────────────────────────────────────────┘
                                │
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Dashboard (/dashboard)                       │
         │                                              │
         │ Card 1: Mercado Total Açúcar = R$12B/ano   │
         │ Card 2: Preço subiu 8% este mês             │
         │ Card 3: Destinos Recomendados (pre-filled)  │
         │ Card 4: Complete perfil para unlock...      │
         └──────────────────────────────────────────────┘
                                │
                                │ Day 2-5: Exploration
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Multiple Simulations (J-AC03)                │
         │                                              │
         │ • Simulates 5 different NCMs                │
         │ • Saves 3 simulations                       │
         │ • Compares results (J-R01)                  │
         │ • Adds 2 countries to watchlist (J-R02)     │
         └──────────────────────────────────────────────┘
                                │
                                │ Day 6: Hits paywall
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Tries to access "Análise de Tarifas         │
         │ Detalhada" (premium feature)                │
         │                                              │
         │ → Paywall Modal appears (J-REV01)           │
         │                                              │
         │ "Upgrade para Premium: Simulações ilimitadas,│
         │  Relatórios PDF, Alertas avançados"         │
         │                                              │
         │ CTA: "Iniciar trial 14 dias grátis"         │
         └──────────────────────────────────────────────┘
                                │
                                │ Clicks "Iniciar trial"
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Checkout (/checkout) - J-REV02               │
         │                                              │
         │ Step 1: Plano → Mensal R$197 (selected)     │
         │ Step 2: Billing → CNPJ, Address             │
         │ Step 3: Payment → PIX (escolhido)           │
         │                                              │
         │ → QR Code displayed                          │
         └──────────────────────────────────────────────┘
                                │
                                │ Scans QR Code, pays R$197
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Checkout Success (/checkout/success)        │
         │                                              │
         │ "Pagamento confirmado! Bem-vinda ao Premium" │
         │ • Email receipt sent                         │
         │ • NFSe invoice generated                     │
         │ • Premium badge on profile                   │
         └──────────────────────────────────────────────┘
                                │
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Dashboard (Premium User)                     │
         │                                              │
         │ • All features unlocked                      │
         │ • "Análise de Tarifas" now accessible       │
         │ • Weekly email alerts configured            │
         │ • Unlimited simulations                     │
         └──────────────────────────────────────────────┘
                                │
                                │ Day 30: Satisfied customer
                                ▼
         ┌──────────────────────────────────────────────┐
         │ Refers 2 colleagues (J-REF01)                │
         │                                              │
         │ • Shares referral link via WhatsApp         │
         │ • Earns 1 month free premium                │
         │                                              │
         │ → K-factor = 2 (viral loop activated)       │
         └──────────────────────────────────────────────┘

                         🎉 SUCCESS: Paid Customer
                         💰 MRR: +R$197
                         📈 LTV: ~R$2,364 (12 months)
```

---

## Parallel Tracks (Dev 1 vs Dev 2)

```
FASE 0: MVP BETA (Weeks 1-5)

DEV 1 (Frontend)                    DEV 2 (Backend)
─────────────────                   ─────────────────

Week 1:                             Week 1:
┌──────────────────┐                ┌──────────────────┐
│ Signup UI        │                │ JWT auth engine  │
│ Login UI         │                │ Password hashing │
│ Password reset   │                │ OAuth2 Google    │
│ Form validation  │                │ Email service    │
└──────────────────┘                └──────────────────┘
        │                                   │
        │            [Both parallel]        │
        │                                   │
Week 2:                             Week 2:
┌──────────────────┐                ┌──────────────────┐
│ Onboarding UI    │                │ User profiles    │
│ 3-step wizard    │                │ Preferences API  │
│ Progress bar     │                │ Simulator auth   │
└──────────────────┘                └──────────────────┘
        │                                   │
        │                                   │
Week 3:                             Week 3:
┌──────────────────┐                ┌──────────────────┐
│ Paywall modal    │                │ Feature flags    │
│ Pricing page     │                │ Trial logic      │
│ Feature toggles  │                │ Tier management  │
└──────────────────┘                └──────────────────┘
        │                                   │
        │                                   │
Week 4:                             Week 4:
┌──────────────────┐                ┌──────────────────┐
│ Checkout UI      │                │ Stripe setup     │
│ Stripe.js        │                │ Subscriptions    │
│ PIX/Card forms   │                │ Webhooks         │
└──────────────────┘                └──────────────────┘
        │                                   │
        │                                   │
Week 5:                             Week 5:
┌──────────────────┐                ┌──────────────────┐
│ Bug fixing       │                │ Data validation  │
│ Performance      │                │ Security audit   │
│ Accessibility    │                │ Load testing     │
└──────────────────┘                └──────────────────┘
        │                                   │
        └───────────────┬───────────────────┘
                        │
                        ▼
                   MVP BETA LAUNCH
```

---

## Dependency Tree (All 32 Journeys)

```
J-AC02 (Auth) ─────────┬─────────────────────────────────────────────┐
                       │                                             │
                       ├─► J-AC01 (Onboarding)                       │
                       │                                             │
                       ├─► J-AC03 (Simulator Logged) ───┬──► J-R01 (Histórico)
                       │                                 │
                       │                                 └──► J-R02 (Watchlist)
                       │
                       ├─► J-AC04 (Dashboard) ──────────┬──► J-R03 (NCMs Related)
                       │                                 │
                       │                                 └──► J-R04 (Concorrentes)
                       │
                       ├─► J-AC05 (Alertas) ────────────┬──► J-R06 (Email)
                       │                                 │
                       │                                 └──► (Email campaigns)
                       │
                       ├─► J-REV01 (Paywall) ───────────┬──► J-REV02 (Checkout)
                       │                                 │
                       │                                 ├──► J-REV03 (Billing)
                       │                                 │
                       │                                 ├──► J-REV04 (Enterprise)
                       │                                 │
                       │                                 └──► J-REV05 (Add-Ons)
                       │
                       ├─► J-REF01 (Referral) ──────────┬──► J-REF05 (Afiliados)
                       │                                 │
                       │                                 └──► (Viral loop)
                       │
                       ├─► J-ADM05 (Support) ───────────┬──► (Ticket system)
                       │                                 │
                       │                                 └──► (CSAT surveys)
                       │
                       └─► J-R05 (Fórum) ───────────────┬──► J-ADM02 (Moderação)
                                                         │
                                                         └──► (Community growth)

INDEPENDENT (No Auth Dependency):
├─► J-A01 (SEO/Blog) ────────────────┬──► (Organic traffic)
├─► J-A02 (Paid Ads) ────────────────┤
├─► J-A03 (Referral Basic) ──────────┤
├─► J-A04 (Parcerias) ───────────────┼──► J-REF03 (Co-Marketing)
├─► J-A05 (Eventos) ─────────────────┤
├─► J-ADM01 (Onboard Parceiros) ─────┤
├─► J-ADM03 (Analytics) ─────────────┤
├─► J-ADM04 (Config) ────────────────┤
└─► J-ADM06 (Data Quality) ──────────┘
```

---

## Metrics Flow (AARRR → NSM)

```
ACQUISITION METRICS                 ACTIVATION METRICS
───────────────────                 ──────────────────
• Signups/week                      • Onboarding completion: >70%
• CAC: <R$150                       • Time-to-first-sim: <5 min
• Acquisition channels:             • Activation rate: >60%
  - Organic: >50%                   • Features explored: >3
  - Referral: >30%                  │
  - Paid: 20%                       │
        │                           │
        └───────────┬───────────────┘
                    │
                    ▼
              NORTH STAR METRIC
              ═════════════════
              Time-to-First-Export-Match
              (Target: <30 days from signup)
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼
RETENTION METRICS          REVENUE METRICS
─────────────────          ───────────────
• D7 retention: >40%       • Free→Premium: >5%
• D30 retention: >25%      • Trial→Paid: >25%
• Churn: <8%/month         • MRR growth: +20%/mo
• Features used: >2/week   • ARPU: >R$50
        │                       │
        └───────────┬───────────┘
                    │
                    ▼
            REFERRAL METRICS
            ────────────────
            • K-factor: >0.5
            • Referral signups: >30%
            • NPS: >60
                    │
                    ▼
              📈 GROWTH LOOP ACTIVATED
```

---

**Document Owner:** BGC Product Management Team
**Last Updated:** 2026-01-20
**Related Docs:**
- [README](./USER-JOURNEYS-README.md) - Navigation guide
- [Complete Map](./USER-JOURNEYS-COMPLETE-MAP.md) - Detailed specs
- [Gantt Chart](./USER-JOURNEYS-GANTT.md) - Timeline
