# 30-Day Playbook - MVP Beta Launch

**Goal:** Implementar 6 jornadas P0 e lançar beta com 50 usuários
**Timeline:** 2026-01-20 → 2026-02-20 (30 dias)
**Team:** 2 Developers Full-Time
**Budget:** R$80k

---

## Week 1: Foundation & Authentication (Jan 20-26)

### Dev 1 (Frontend)
- [ ] **Day 1-2:** Setup projeto Next.js auth
  - [ ] Instalar NextAuth.js ou Clerk
  - [ ] Criar páginas: `/signup`, `/login`, `/forgot-password`
  - [ ] Design system: Auth components (form, buttons, inputs)
  - [ ] Validação Zod (email, password strength)

- [ ] **Day 3-4:** UI Implementation
  - [ ] Signup form (email, password, nome empresa, CNPJ opcional)
  - [ ] Login form (email/password + Google OAuth button)
  - [ ] Password reset flow (input email → check inbox)
  - [ ] Loading states e error handling

- [ ] **Day 5:** Integration & Testing
  - [ ] Integrar com backend API (POST /v1/auth/signup, /v1/auth/login)
  - [ ] Cookie management (httpOnly JWT)
  - [ ] Redirect logic (after login → /onboarding ou /dashboard)
  - [ ] Testes E2E (Playwright: signup, login, logout)

### Dev 2 (Backend)
- [ ] **Day 1-2:** Auth Engine
  - [ ] Migration: users table (id, email, password_hash, created_at, role, tier)
  - [ ] bcrypt/Argon2 password hashing
  - [ ] JWT token generation (access + refresh tokens)
  - [ ] Middleware: authenticateUser() (extract JWT, validate, attach user to context)

- [ ] **Day 3:** OAuth2 Integration
  - [ ] Google OAuth2 setup (client ID/secret)
  - [ ] Callback handler (exchange code for token, create/get user)
  - [ ] Link accounts (email from OAuth matches existing user)

- [ ] **Day 4:** Email Service
  - [ ] SendGrid/SES setup (API keys, templates)
  - [ ] Email verification flow (send token link, verify endpoint)
  - [ ] Password reset flow (send reset link, validate token, update password)

- [ ] **Day 5:** Testing & Security
  - [ ] Unit tests (auth handlers, middleware)
  - [ ] Rate limiting (5 login attempts/min per IP)
  - [ ] CAPTCHA integration (after 3 failed attempts)
  - [ ] Security audit (SQL injection, XSS prevention)

### Product/Design
- [ ] **Day 1:** Finalize mockups
  - [ ] Signup, Login, Forgot Password screens (Figma)
  - [ ] Email templates (verification, password reset)

- [ ] **Day 2-3:** Onboarding design
  - [ ] 3-step wizard mockups (NCM, Volume, Regions)
  - [ ] Progress indicator, skip buttons
  - [ ] Empty states, error states

- [ ] **Day 4-5:** Paywall & Checkout design
  - [ ] Paywall modal (feature comparison table)
  - [ ] Pricing page (Free vs Premium)
  - [ ] Checkout flow (plan selection, billing info, payment)

### Ops
- [ ] Setup Stripe account (test mode)
- [ ] Configure email domain (SPF, DKIM for SendGrid)
- [ ] Analytics events (Mixpanel/Segment: signup, login, logout)

---

## Week 2: Onboarding & Simulator Integration (Jan 27 - Feb 2)

### Dev 1 (Frontend)
- [ ] **Day 1-2:** Onboarding Wizard
  - [ ] `/onboarding/welcome` page (boas-vindas, CTA "Começar")
  - [ ] `/onboarding/step-1` (NCM autocomplete, fuzzy search)
  - [ ] `/onboarding/step-2` (Volume input, validação >0)
  - [ ] `/onboarding/step-3` (Multi-select países/continentes)

- [ ] **Day 3:** Wizard Logic
  - [ ] Progress tracker (1/3, 2/3, 3/3)
  - [ ] Skip button (vai direto pro dashboard)
  - [ ] Save progress (localStorage ou backend)
  - [ ] Submit onboarding → POST /v1/users/profile

- [ ] **Day 4-5:** Dashboard Skeleton
  - [ ] `/dashboard` page (layout grid 2x2)
  - [ ] Card 1: "Seu Mercado" (placeholder)
  - [ ] Card 2: "Tendências" (placeholder)
  - [ ] Card 3: "Destinos Recomendados" (pre-filled simulation)
  - [ ] Card 4: "Próximas Ações" (TODOs: complete perfil, etc.)

### Dev 2 (Backend)
- [ ] **Day 1-2:** User Profile & Preferences
  - [ ] Migration: user_profiles (user_id, primary_ncm, avg_volume_kg, target_regions)
  - [ ] Endpoint: POST /v1/users/profile (onboarding data)
  - [ ] Endpoint: GET /v1/users/profile (fetch saved preferences)

- [ ] **Day 3-4:** Simulator Auth Integration
  - [ ] Modify POST /v1/simulator/destinations:
    - [ ] Check if user authenticated (JWT)
    - [ ] If authenticated: no rate limit, save simulation
    - [ ] If anonymous: rate limit (5/day per IP)
  - [ ] Migration: user_simulations (user_id, ncm, volume, results JSON, created_at)
  - [ ] Endpoint: POST /v1/users/simulations (save simulation)

- [ ] **Day 5:** Pre-fill Logic
  - [ ] Endpoint: GET /v1/users/simulations/latest (fetch last simulation)
  - [ ] Auto-run simulation on dashboard load (if first time, use onboarding NCM)

### Testing
- [ ] E2E flow: Signup → Onboarding (3 steps) → Dashboard → Simulator
- [ ] Verify: Simulation saved, results persist, user can see history

---

## Week 3: Paywall & Premium Discovery (Feb 3-9)

### Dev 1 (Frontend)
- [ ] **Day 1-2:** Paywall Modal
  - [ ] Component: `<PaywallModal />` (trigger when user tries premium feature)
  - [ ] Feature comparison table (Free vs Premium)
  - [ ] CTA: "Iniciar trial 14 dias grátis" ou "Assinar Premium"
  - [ ] Close button (soft block: user can dismiss)

- [ ] **Day 2-3:** Pricing Page
  - [ ] `/pricing` route (comparison table)
  - [ ] Toggle: Mensal (R$197) vs Anual (R$1.970, economize 17%)
  - [ ] Highlight premium features (simulações ilimitadas, PDF, alertas, suporte)
  - [ ] Social proof (depoimentos, logos de clientes beta)

- [ ] **Day 4-5:** Feature Flags UI
  - [ ] Premium badge (header, avatar)
  - [ ] Lock icons em features premium (dashboard advanced metrics)
  - [ ] Upgrade banner (sticky top: "Upgrade para Premium")

### Dev 2 (Backend)
- [ ] **Day 1-2:** Feature Flags System
  - [ ] Migration: users add column `tier` (free, premium, enterprise)
  - [ ] Middleware: checkFeatureAccess(feature) → 403 if not allowed
  - [ ] Feature config (JSON): `{ "advanced_analytics": ["premium", "enterprise"] }`

- [ ] **Day 3-4:** Trial Logic
  - [ ] Migration: users add columns `trial_ends_at`, `trial_used`
  - [ ] Endpoint: POST /v1/users/start-trial (14 days free premium)
  - [ ] Cron job: Check daily for expired trials → downgrade to free

- [ ] **Day 5:** Analytics Events
  - [ ] Track: paywall_shown, paywall_dismissed, pricing_page_viewed
  - [ ] Track: trial_started, trial_converted, trial_expired

### Product
- [ ] Finalize pricing research (check competitors: Logcomex, Datawise)
- [ ] A/B test plan: Test R$197 vs R$247 monthly (50/50 split)

---

## Week 4: Checkout & Payment (Feb 10-16)

### Dev 1 (Frontend)
- [ ] **Day 1-2:** Checkout Flow
  - [ ] `/checkout` page (3 steps: Plan → Billing Info → Payment)
  - [ ] Step 1: Plan selection (Mensal vs Anual, price display)
  - [ ] Step 2: Billing form (nome, CNPJ, endereço)
  - [ ] Step 3: Payment method (card, boleto, PIX)

- [ ] **Day 3-4:** Stripe Integration
  - [ ] Stripe.js SDK (never send card data to backend)
  - [ ] `<CardElement />` component (Stripe pre-built)
  - [ ] Handle payment intent (create, confirm, success/error)
  - [ ] PIX payment (display QR code, poll for confirmation)

- [ ] **Day 5:** Post-Checkout
  - [ ] `/checkout/success` page (confirmation, welcome email sent)
  - [ ] `/checkout/failed` page (retry button, error message)
  - [ ] Redirect to dashboard (with premium badge)

### Dev 2 (Backend)
- [ ] **Day 1-2:** Stripe Setup
  - [ ] Stripe account (API keys: test + prod)
  - [ ] Create products (Premium Monthly, Premium Annual)
  - [ ] Webhook endpoint: POST /v1/webhooks/stripe (listen for events)

- [ ] **Day 3-4:** Subscription Management
  - [ ] Migration: subscriptions (user_id, stripe_subscription_id, status, plan, current_period_end)
  - [ ] Endpoint: POST /v1/subscriptions/create (create Stripe subscription, save to DB)
  - [ ] Webhook handlers:
    - [ ] `checkout.session.completed` → activate subscription
    - [ ] `invoice.payment_succeeded` → extend subscription
    - [ ] `invoice.payment_failed` → retry/downgrade

- [ ] **Day 5:** Invoice & Receipt
  - [ ] NFSe integration (API de prefeitura SP ou Focus NFe)
  - [ ] Email transacional (receipt com PDF invoice)

### Testing
- [ ] E2E: Full purchase flow (card, PIX, boleto)
- [ ] Test webhooks (Stripe CLI: trigger events locally)
- [ ] Verify: User upgraded, features unlocked, invoice sent

---

## Week 5: Polish, Testing & Beta Launch (Feb 17-23)

### Dev 1 (Frontend)
- [ ] **Day 1-2:** Bug Fixing
  - [ ] Review all user flows (signup → onboarding → simulator → upgrade)
  - [ ] Fix edge cases (empty states, error handling, loading states)
  - [ ] Accessibility audit (WCAG AA: keyboard nav, screen readers)

- [ ] **Day 3:** Performance Optimization
  - [ ] Lighthouse audit (target: >90 performance score)
  - [ ] Code splitting (lazy load routes)
  - [ ] Image optimization (next/image, WebP)

- [ ] **Day 4-5:** Documentation
  - [ ] User guide (how to use simulator, interpret results)
  - [ ] FAQ page (common questions)
  - [ ] Video tutorial (Loom: 3min onboarding walkthrough)

### Dev 2 (Backend)
- [ ] **Day 1-2:** Data Validation & Security
  - [ ] Input sanitization (prevent SQL injection, XSS)
  - [ ] Rate limiting (all endpoints: 100 req/min per user)
  - [ ] CORS policy (whitelist frontend domain)

- [ ] **Day 3:** Monitoring & Alerts
  - [ ] Sentry setup (error tracking)
  - [ ] Grafana alerts (CPU >80%, DB slow queries, API errors >5%)
  - [ ] Uptime monitoring (UptimeRobot: check /healthz every 5min)

- [ ] **Day 4-5:** Load Testing
  - [ ] k6 script: Simulate 100 concurrent users
  - [ ] Target: P95 latency <500ms, error rate <1%
  - [ ] Fix bottlenecks (add DB indexes, cache warming)

### Product/Marketing
- [ ] **Day 1-2:** Beta User Recruitment
  - [ ] Outreach list (50 exporters from APEX, ABIEC, LinkedIn)
  - [ ] Invitation email (personalized, offer 6 months free premium)
  - [ ] Waitlist page (collect emails, gauge interest)

- [ ] **Day 3:** Beta Launch
  - [ ] Soft launch: Send invites to first 10 users
  - [ ] Monitor: Watch analytics, check for errors
  - [ ] Support: Set up Intercom/Crisp chat (respond <2h)

- [ ] **Day 4-5:** Feedback Collection
  - [ ] Schedule interviews (10 users, 30min each)
  - [ ] Survey (NPS, feature requests, pain points)
  - [ ] Iterate: Prioritize quick wins for Week 6

### Ops
- [ ] Deploy to production (k8s staging → prod promotion)
- [ ] DNS setup (app.brasilglobalconect.com)
- [ ] SSL certificate (Let's Encrypt auto-renewal)
- [ ] Backups (automated daily PostgreSQL dumps to S3)

---

## Success Criteria (End of 30 Days)

### Quantitative Metrics
- [ ] **50 beta signups** (invited users created accounts)
- [ ] **30% activation** (15 users completed onboarding + first simulation)
- [ ] **3 paid conversions** (users started trial or paid immediately)
- [ ] **NPS >40** (from survey: "How likely to recommend?")
- [ ] **Zero critical bugs** (P0/P1 issues resolved)

### Qualitative Metrics
- [ ] **10 user interviews** (30min each, collect feedback)
- [ ] **5+ positive testimonials** (can be used in marketing)
- [ ] **1 case study candidate** (user who found real value, willing to share story)

### Technical Metrics
- [ ] **99% uptime** (excluding planned maintenance)
- [ ] **P95 latency <500ms** (API response time)
- [ ] **Zero security incidents** (no breaches, vulnerabilities)
- [ ] **Test coverage >80%** (backend unit tests)

---

## Daily Standup Format

**Time:** 9:30 AM (15 min max)

**Structure:**
1. **Yesterday:** What I completed (link to PRs)
2. **Today:** What I'm working on (task breakdown)
3. **Blockers:** Anything stopping me (need help with X)

**Example:**
```
Dev 1 (Frontend):
- Yesterday: Completed signup form (PR #45), started onboarding wizard
- Today: Finish onboarding step 2-3, integrate with backend API
- Blockers: Waiting for NCM autocomplete API endpoint (Dev 2)

Dev 2 (Backend):
- Yesterday: Auth engine done, JWT working, OAuth2 90% complete
- Today: Finish OAuth2, start user profile endpoints
- Blockers: None, but need Stripe API keys from Ops
```

---

## Weekly Retrospectives

**Time:** Friday 4 PM (1 hour)

**Format:**
1. **What went well?** (celebrate wins)
2. **What didn't go well?** (be honest, no blame)
3. **Action items** (specific, assignee, deadline)

**Example:**
```
Week 1 Retro:
✅ Went Well:
- Auth engine shipped on time
- Great collaboration on UI/UX

❌ Didn't Go Well:
- Email service delayed (SendGrid approval took 2 days)
- Underestimated OAuth2 complexity

🎯 Actions:
- [Dev 2] Research alternatives to SendGrid (AWS SES) - by Monday
- [PM] Buffer 20% time in estimates for unknowns - ongoing
```

---

## Risk Mitigation (Week-by-Week)

### Week 1 Risks
- **Auth security flaw** → Mitigate: Security audit on Day 5, penetration test
- **OAuth2 integration fails** → Mitigate: Fallback to email/password only (OAuth2 nice-to-have)

### Week 2 Risks
- **Onboarding too complex** → Mitigate: User testing on Day 3, simplify if >5 min
- **Simulator integration breaks** → Mitigate: Feature flag, rollback if errors

### Week 3 Risks
- **Pricing too high (no conversions)** → Mitigate: A/B test, ready to discount
- **Paywall too aggressive (users churn)** → Mitigate: Soft paywall, monitor dismissal rate

### Week 4 Risks
- **Stripe payment fails** → Mitigate: Test extensively, fallback to manual invoice
- **NFSe integration delays** → Mitigate: Manual invoice generation (short-term)

### Week 5 Risks
- **Not enough beta signups** → Mitigate: Extend timeline, increase outreach
- **Critical bug in prod** → Mitigate: Canary deploy (10% traffic first), rollback plan

---

## Communication Plan

### Internal (Team)
- **Daily Standup:** Slack #standup channel (async OK if remote)
- **Weekly Sync:** Friday 4 PM (retro + planning)
- **Blockers:** Ping immediately in #dev-urgent

### External (Stakeholders)
- **Weekly Update:** Email to CEO/stakeholders every Friday
  - What shipped
  - Metrics update
  - Next week priorities
  - Risks/asks

### Beta Users
- **Week 1:** Invitation email (waitlist, gauge interest)
- **Week 3:** Early access email (first 10 users)
- **Week 5:** Public beta email (remaining 40 users)
- **Ongoing:** Weekly digest (product updates, tips)

---

## Budget Breakdown (R$80k)

| Item | Cost | Notes |
|------|------|-------|
| **Dev 1 (Frontend)** | R$32k | R$8k/week × 4 weeks |
| **Dev 2 (Backend)** | R$32k | R$8k/week × 4 weeks |
| **Stripe Fees** | R$2k | Assumes R$591 MRR × 3.5% fee |
| **SendGrid** | R$1k | Email service (5k emails/month) |
| **Infra (AWS)** | R$3k | EC2, RDS, S3, CloudFront |
| **Tools (Figma, Sentry, etc.)** | R$2k | Design, monitoring, analytics |
| **Buffer (20%)** | R$8k | Contingency for unknowns |
| **TOTAL** | **R$80k** | |

---

## Decision Points (Go/No-Go)

### End of Week 2 (Feb 2)
**Check:** Is onboarding working? Can users complete flow?
- ✅ Yes → Proceed to Week 3
- ❌ No → Extend Week 2, delay paywall

### End of Week 4 (Feb 16)
**Check:** Is checkout working? Can users pay?
- ✅ Yes → Proceed to Week 5 (launch)
- ❌ No → Delay launch, fix critical payment bugs

### End of Week 5 (Feb 23)
**Check:** Did we hit success criteria? (50 signups, 3 paid, NPS >40)
- ✅ Yes → Invest in Fase 1 (V1 Public Launch)
- ⚠️ Partial → Extend beta, iterate based on feedback
- ❌ No → Pivot or kill, hold post-mortem

---

## Next Steps (After 30 Days)

**If Successful:**
1. Sprint Planning: Fase 1 (V1 Public Launch)
2. Hire/contract: +1 Designer, +1 Marketing
3. Budget approval: R$150k for next 6 weeks
4. Start: J-A01 (SEO/Blog), J-AC04 (Dashboard)

**If Needs Iteration:**
1. Analyze feedback (why NPS low? why no conversions?)
2. Prioritize fixes (top 3 pain points)
3. Extend beta: +2 weeks, recruit more users
4. Re-test: New cohort, measure improvement

**If Failed:**
1. Post-mortem: What went wrong? (product, market, execution?)
2. Pivot options:
   - Change target persona (enterprise vs SME?)
   - Change pricing (lower? freemium more generous?)
   - Change value prop (focus on different pain point?)
3. Decision: Pivot vs Kill (CEO call)

---

**Document Owner:** BGC Product Management
**Last Updated:** 2026-01-20
**Review Cadence:** Daily (progress check), Weekly (retrospective)
