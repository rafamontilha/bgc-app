# Product Context - Brasil Global Connect (BGC)

**Document Version:** 1.0.0
**Last Updated:** 2026-01-12
**Owner:** Product Management Team
**Status:** Active

---

## Table of Contents

- [Product Overview](#product-overview)
- [Value Proposition](#value-proposition)
- [Strategic Pillars](#strategic-pillars)
- [Technical Architecture](#technical-architecture)
- [Core Features](#core-features)
- [Data Coverage Strategy](#data-coverage-strategy)
- [Product Metrics](#product-metrics)
- [Current Limitations](#current-limitations)
- [Roadmap Highlights](#roadmap-highlights)

---

## Product Overview

**Brasil Global Connect (BGC)** is a comprehensive SaaS platform designed to democratize international trade intelligence and operations for Brazilian SMEs engaged in export/import activities.

### Mission
Reduce friction, cost, and complexity in export/import operations by unifying official data sources, generating actionable market insights, and connecting Brazilian exporters with global buyers.

### Vision
Become the leading international trade platform for Brazilian SMEs, enabling predictable, data-driven export decisions and seamless cross-border operations.

### Target Market
- **Primary:** Brazilian SMEs with export revenue between R$1M - R$50M/year
- **Secondary:** Trading companies, cooperatives, and mid-market exporters
- **Tertiary:** International buyers seeking Brazilian products

---

## Value Proposition

### For Exporters (SMEs)
- **Reduce time-to-market** from 6-12 months to 2-4 weeks
- **Data-driven decisions** using official ComexStat and Siscomex data
- **Market intelligence** on demand, prices, routes, and competition
- **Buyer matchmaking** with verified international buyers
- **Operational simplification** through integrated logistics, finance, and insurance

### For Buyers (International)
- **Discover vetted** Brazilian suppliers across 97 NCM chapters
- **Transparent pricing** and market data
- **Streamlined communication** with Brazilian exporters
- **Risk mitigation** through integrated insurance and trade finance

---

## Strategic Pillars

### 1. Data Intelligence
Unify and enrich official trade data from multiple sources:
- ComexStat API (MDIC)
- Siscomex / Portal Unico
- Receita Federal
- International customs databases

### 2. Market Insights
Generate predictive analytics and recommendations:
- TAM/SAM/SOM analysis by NCM × Country × Year
- Route comparison (export corridors)
- **Export Destination Simulator** (launched v0.4.0)
- Price forecasting (roadmap)

### 3. Matchmaking
Connect exporters with verified buyers:
- B2B marketplace (roadmap)
- Lead scoring and qualification
- Communication facilitation

### 4. Operational Enablement
Integrate third-party services:
- Logistics (freight forwarders, customs brokers)
- Trade finance (factoring, export credit)
- Insurance (cargo, credit)
- FX services

---

## Technical Architecture

### Stack Overview

**Backend:**
- Go 1.24.9 + Gin framework
- Clean Architecture (domain-driven design)
- Integration Gateway (Hybrid Connector Framework)

**Database:**
- PostgreSQL 16
- Materialized views for reporting
- Scheduled refresh via cron jobs

**Observability:**
- Prometheus + Grafana (metrics)
- Jaeger (distributed tracing)
- OpenTelemetry SDK

**Frontend:**
- Next.js 15 + React 19 + TypeScript 5
- Tailwind CSS
- Material Design 3 with Apple aesthetic
- MUI v7 patterns

**Infrastructure:**
- Docker + Docker Compose
- Kubernetes k3d
- Traefik ingress

### Design Principles
- **API-first:** All features exposed via REST API
- **Freemium model:** Rate-limited free tier, unlimited premium
- **Scalability:** Horizontal pod autoscaling (HPA)
- **Security:** Network policies, sealed secrets, RBAC
- **Observability:** Full metrics, tracing, and alerting

---

## Core Features

### ✅ Launched (v0.4.0)

#### Epic 1: Enablement & Access
- Integration Gateway with Hybrid Connector Framework
- 90% config-driven integrations (YAML)
- Auth engine (mTLS, OAuth2, API Key)
- Circuit breaker, retry, rate limiting

#### Epic 2: Observability
- 11 Prometheus metrics (HTTP, DB, cache, errors)
- OpenTelemetry distributed tracing
- Grafana dashboards
- Jaeger UI

#### Epic 3: Data Contract
- JSON Schema validation
- API versioning (/v1/*)
- Idempotency middleware
- Data dictionary

#### Epic 4: Export Destination Simulator (MVP) 🌍
- **Launch Date:** 2026-01-09
- **Status:** Production Ready (E2E validated)
- **Endpoint:** `POST /v1/simulator/destinations`
- **Coverage:** **8% of Brazilian exports** (Chapter 17 - Sugar only)
- **Features:**
  - Smart scoring algorithm (market size, growth, price, distance)
  - Financial estimates (margin, logistics cost, tariffs, lead time)
  - Freemium rate limiting (5 req/day free, unlimited premium)
  - Response time: 22-92ms (10x better than target)
  - 5 destinations for NCM 17011400, 2 for NCM 02013000

---

## Data Coverage Strategy

### Current State (v0.4.0)

**Simulator Data Coverage:**
- ✅ **Chapter 17 (Sugar):** Açúcar de cana e derivados
  - **Export Value:** ~R$12B/year
  - **% of BR Exports:** ~8%
  - **NCMs:** 17011400 (primary)
  - **Countries:** 5 destinations (US, CN, DE, IN, GB)
  - **Data Points:** 10 export records (Dec 2025 - Jan 2026)

**Limitation Impact:**
- 92% of Brazilian exports not covered by simulator
- Most requested NCMs (meat, fruits, soybeans, machinery) unavailable
- Users outside sugar sector cannot use simulator

### Expansion Strategy: 2-Phase Incremental Approach

#### Phase 1: High-Impact Chapters (v0.4.1 - This Week)
**Timeline:** 2-3 days
**Priority:** P0 - Critical

**Target Coverage:** 28% of Brazilian exports (+20%)

| Chapter | Product | Export Value | % of Exports | Priority |
|---------|---------|--------------|--------------|----------|
| 02 | Carnes e miudezas comestíveis | R$45B | 12% | **P0** |
| 08 | Frutas, cascas de cítricos e melões | R$30B | 8% | **P0** |
| 17 | Açúcares e produtos de confeitaria | R$12B | 8% | ✅ Done |

**Deliverables:**
- Populate `stg.exportacao` with Chapter 02 and 08 data
- Validate simulator functionality with new chapters
- Update product messaging: "Covers 28% of BR exports (sugar, meat, fruits)"

**Success Metrics:**
- Simulator requests increase by 3x
- User complaints about NCM coverage decrease by 50%
- Freemium → Premium conversion rate baseline established

#### Phase 2: Full Coverage (v0.5.0 - v0.6.0)
**Timeline:** 2-4 weeks
**Priority:** P1 - High

**Target Coverage:** 100% of Brazilian exports (+72%)

**Rollout Plan:**

**v0.5.0 (Week 3-4) - 60% coverage:**
- Chapter 84: Máquinas e aparelhos mecânicos (R$15B, 10%)
- Chapter 85: Máquinas e materiais elétricos (R$12B, 8%)
- Chapter 12: Sementes e frutos oleaginosos (R$60B, 15%) ← **TOP PRIORITY**
- Chapter 22: Bebidas (R$3B, 2%)
- Chapter 27: Combustíveis minerais (R$42B, 17%)

**v0.6.0 (Month 2) - 100% coverage:**
- All remaining 91 NCM chapters
- Automated ingestion pipeline from ComexStat API
- Scheduled refresh (weekly/monthly)

**Technical Approach:**
1. **Data Ingestion:**
   - Kubernetes CronJob to call ComexStat API
   - Filter by year (2024-2026) and chapter
   - Upsert into `stg.exportacao` with conflict resolution

2. **Data Validation:**
   - Ensure minimum 12 months of data per NCM × Country
   - Validate growth rate calculations
   - Alert on missing critical countries

3. **Performance Optimization:**
   - Materialized views for top 100 NCMs per chapter
   - Redis cache warming via CronJob
   - Connection pooling tuning

**Success Metrics:**
- 100% NCM coverage
- Zero "NCM not found" errors for major products
- P95 response time remains < 200ms
- Database size < 50GB

---

## Product Metrics

### North Star Metric
**Number of successful export deals facilitated by BGC platform**

### Leading Indicators
- Monthly Active Users (MAU)
- Simulator requests per user
- Freemium → Premium conversion rate
- NCM coverage %
- API uptime (target: 99.9%)

### Current Baseline (v0.4.0)
- NCM Coverage: 8% (Chapter 17 only)
- Simulator Avg Response Time: 57ms
- Simulator P95 Response Time: 92ms (vs 200ms target)
- Rate Limit Free Tier: 5 req/day
- E2E Test Pass Rate: 100%

### Targets (Q1 2026)
- NCM Coverage: 100% (all chapters)
- MAU: 500 exporters
- Simulator Requests: 10k/month
- Freemium → Premium Conversion: 5%
- API Uptime: 99.9%

---

## Current Limitations

### Data Coverage (Highest Priority)
- ❌ **Only 8% of Brazilian exports covered** (Chapter 17)
- ❌ Most requested sectors unavailable (meat, soy, machinery)
- ❌ Limited historical data (2 months only)
- ⚠️ Growth rate calculations may be unstable with < 12 months data

### Feature Completeness
- ❌ No buyer-side marketplace (roadmap)
- ❌ No integrated logistics quotes (roadmap)
- ❌ No ML-based demand forecasting (roadmap)
- ⚠️ Tariff and lead time estimates are heuristic (not real-time)

### Operational
- ⚠️ Manual data updates (CronJob not deployed)
- ⚠️ No alerting for data freshness
- ⚠️ Rate limiting not yet tested at scale

---

## Roadmap Highlights

### Q1 2026

**v0.4.1 (Week 2) - Expanded Coverage Phase 1**
- Add Chapter 02 (Meat) and 08 (Fruits)
- Target: 28% export coverage
- Release notes and user communication

**v0.5.0 (Week 3-4) - Major Sectors**
- Add Chapters 12, 27, 84, 85
- Target: 60% export coverage
- Automated ingestion pipeline
- Materialized views for top NCMs

**v0.6.0 (Month 2) - Full Coverage**
- All 97 NCM chapters
- 100% export coverage
- Cache warming CronJob
- Performance validation

### Q2 2026
- ML-based demand forecasting
- Real-time logistics API integration
- Buyer marketplace MVP
- Mobile app (iOS/Android)

### Q3 2026
- Trade finance integration
- Insurance quotes API
- FX services
- Multi-language support (EN, ES, CN)

---

## Key Decisions & ADRs

### ADR-001: 2-Phase NCM Coverage Strategy
- **Date:** 2026-01-12
- **Decision:** Incremental rollout instead of big-bang
- **Rationale:** De-risk deployment, validate value per chapter, prioritize high-impact sectors
- **Status:** Approved

### ADR-002: Freemium Rate Limiting
- **Date:** 2025-11-22
- **Decision:** 5 req/day free, unlimited premium
- **Rationale:** Balances discoverability with monetization
- **Status:** Active

### ADR-003: Heuristic Estimates for MVP
- **Date:** 2025-11-21
- **Decision:** Use rule-based estimates for tariffs, lead time, logistics cost
- **Rationale:** Real-time API integration is Phase 2, unblock MVP launch
- **Status:** Active (revisit in Q2 2026)

---

## Contact & Governance

**Product Owner:** Rafael (BGC PM)
**Engineering Lead:** Backend Team
**Design Lead:** UX Team
**Data Owner:** Data Engineering Team

**Review Frequency:** Quarterly
**Last Review:** 2026-01-12
**Next Review:** 2026-04-12

---

**Document Status:** Living document - updated continuously
**Change Log:** See git history for detailed changes
