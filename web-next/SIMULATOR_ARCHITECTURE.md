# Arquitetura do Simulador - Frontend

Visão geral da arquitetura e fluxo de dados do Simulador de Destinos de Exportação.

---

## Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────────┐
│                       Browser (Port 3000)                        │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                   app/simulator/page.tsx                    │ │
│  │                   (Orchestrator Component)                  │ │
│  │                                                              │ │
│  │  State Management:                                          │ │
│  │  - isLoading: boolean                                       │ │
│  │  - data: SimulatorResponse | null                           │ │
│  │  - error: ErrorType | null                                  │ │
│  │  - rateLimitInfo: RateLimitInfo | null                      │ │
│  │  - showUpgradeModal: boolean                                │ │
│  │                                                              │ │
│  │  Conditional Rendering:                                     │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │  <RateLimitBanner />        (US-003)                 │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │  <SimulatorForm />          (US-001)                 │  │ │
│  │  │    ├── NCM Input (Zod validation)                    │  │ │
│  │  │    ├── Volume Input (Zod validation)                 │  │ │
│  │  │    └── Submit Button (loading state)                 │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │  if (error)                                           │  │ │
│  │  │    <ErrorState type={error.type} />   (US-005)       │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │  if (data || isLoading)                               │  │ │
│  │  │    <DestinationList />          (US-002)             │  │ │
│  │  │      ├── if (isLoading): Skeleton (3 cards)          │  │ │
│  │  │      ├── if (empty): Empty State                     │  │ │
│  │  │      └── else: map(destinations)                     │  │ │
│  │  │            └── <DestinationCard /> x N               │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │  <UpgradeModal open={showUpgradeModal} />  (US-004)  │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP POST
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     lib/api/simulator.ts                         │
│                      (API Client Layer)                          │
│                                                                   │
│  simulateDestinations(request: SimulatorRequest)                │
│    ├── Validates input                                          │
│    ├── Builds HTTP request                                      │
│    ├── Parses rate limit headers                                │
│    ├── Handles errors (404, 400, 429, 500, 0)                   │
│    └── Returns ApiResponse<SimulatorResponse>                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ fetch()
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Backend API (Port 8080)                        │
│                                                                   │
│  POST /v1/simulator/destinations                                │
│                                                                   │
│  Headers:                                                        │
│    X-RateLimit-Limit: 5                                         │
│    X-RateLimit-Remaining: 3                                     │
│    X-RateLimit-Reset: 1736503200                                │
│                                                                   │
│  Response:                                                       │
│  {                                                               │
│    "destinations": [...],                                       │
│    "metadata": {...}                                            │
│  }                                                               │
└─────────────────────────────────────────────────────────────────┘
```

---

## Fluxo de Dados (Data Flow)

### 1. Simulação Normal (Happy Path)

```
User Input (NCM + Volume)
        │
        ▼
SimulatorForm (Zod validation)
        │
        ▼
page.tsx (handleSimulate)
        │
        ▼
lib/api/simulator.ts (simulateDestinations)
        │
        ▼
Backend API (POST /v1/simulator/destinations)
        │
        ▼
Parse Response + Rate Limit Headers
        │
        ▼
Update State (data, rateLimitInfo)
        │
        ▼
Render DestinationList → DestinationCard (x10)
```

### 2. Rate Limit Exceeded (HTTP 429)

```
User submits 6th simulation
        │
        ▼
API Client (simulateDestinations)
        │
        ▼
Receives HTTP 429
        │
        ▼
Throws SimulatorApiError(429)
        │
        ▼
page.tsx (catch block)
        │
        ▼
Sets showUpgradeModal = true
        │
        ▼
UpgradeModal opens
```

### 3. Error Handling (404, 500, etc.)

```
User submits invalid NCM
        │
        ▼
API returns HTTP 404
        │
        ▼
API Client throws SimulatorApiError(404)
        │
        ▼
page.tsx determines errorType = "not_found"
        │
        ▼
Sets error state
        │
        ▼
ErrorState renders with type="not_found"
```

---

## Estrutura de Tipos (Type System)

```typescript
// types/simulator.ts

SimulatorRequest          → API Input
  ├── ncm: string
  ├── volume_kg: number
  ├── countries?: string[]
  └── max_results?: number

SimulatorResponse         → API Output
  ├── destinations: Destination[]
  └── metadata: SimulatorMetadata

Destination               → Single Result
  ├── rank: number
  ├── country_code: string
  ├── country_name_pt: string
  ├── score: number
  ├── demand_level: 'Alto' | 'Médio' | 'Baixo'
  ├── market_size_usd: number
  ├── growth_rate_pct: number
  ├── price_per_kg_usd: number
  ├── distance_km: number
  ├── estimated_margin_pct: number
  ├── logistics_cost_usd: number
  ├── tariff_rate_pct: number
  ├── lead_time_days: number
  └── recommendation_reason: string

RateLimitInfo             → Rate Limit Headers
  ├── limit: number
  ├── remaining: number
  └── reset: number (unix timestamp)

ApiError                  → Error Response
  ├── error: string
  ├── message: string
  ├── status_code: number
  └── details?: Record<string, unknown>
```

---

## Camadas de Validação

### 1. Client-Side Validation (Zod)

**Localização:** `components/simulator/SimulatorForm.tsx`

```typescript
// Validação antes de enviar ao backend
const schema = z.object({
  ncm: z.string().length(8).regex(/^\d{8}$/),
  volume_kg: z.number().positive().max(1000000),
});

// Errors exibidos imediatamente no form
```

### 2. API Client Validation

**Localização:** `lib/api/simulator.ts`

```typescript
// Validação adicional antes do fetch
validateNCM(ncm)
validateVolume(volume_kg)
validateCountryCodes(countries)
validateMaxResults(max_results)
```

### 3. Backend Validation

**Localização:** Backend API (Go)

```go
// Validação no backend (última camada)
// NCM existe no banco?
// Volume é razoável?
// Países são ISO codes válidos?
```

---

## Estados de UI (State Machine)

```
                    ┌─────────────┐
                    │   IDLE      │
                    │ (Initial)   │
                    └─────────────┘
                          │
                          │ User submits form
                          ▼
                    ┌─────────────┐
                    │  LOADING    │
                    │ (Spinner)   │
                    └─────────────┘
                          │
           ┌──────────────┼──────────────┐
           │              │              │
     Success         Error 429      Error (other)
           │              │              │
           ▼              ▼              ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │   SUCCESS   │ │ RATE_LIMIT  │ │   ERROR     │
    │ (Results)   │ │ (Modal)     │ │ (ErrorState)│
    └─────────────┘ └─────────────┘ └─────────────┘
           │              │              │
           │         Close modal    Retry button
           │              │              │
           └──────────────┴──────────────┘
                          │
                          ▼
                    ┌─────────────┐
                    │   IDLE      │
                    └─────────────┘
```

---

## Componente Tree (Hierarquia)

```
app/layout.tsx
  └── ThemeProvider (MUI v7)
      └── app/simulator/page.tsx (Main)
          ├── RateLimitBanner
          │   └── Alert (MUI)
          │       ├── Typography
          │       ├── LinearProgress
          │       └── Button
          ├── SimulatorForm
          │   └── Stack (MUI)
          │       ├── TextField (NCM)
          │       ├── TextField (Volume)
          │       └── Button (Submit)
          ├── ErrorState (conditional)
          │   └── Stack
          │       ├── Icon
          │       ├── Typography
          │       ├── Button (Retry)
          │       └── Alert (Tips)
          ├── DestinationList (conditional)
          │   ├── if (loading): DestinationSkeleton x3
          │   ├── if (empty): EmptyState
          │   └── else: DestinationCard x N
          │       └── Card (MUI)
          │           ├── Rank Badge
          │           ├── Country Flag + Name
          │           ├── Demand Chip
          │           ├── Score Progress Bar
          │           ├── Metrics Grid (2x2)
          │           ├── Recommendation Box
          │           └── Additional Details (4 cols)
          └── UpgradeModal (conditional)
              └── Dialog (MUI)
                  ├── DialogTitle
                  ├── DialogContent
                  │   ├── Feature List (6 items)
                  │   └── Price Box
                  └── DialogActions
                      ├── Button (Cancel)
                      └── Button (View Plans)
```

---

## Tema MUI (Design Tokens)

### Paleta de Cores

```
Primary (Blue)
├── main:    #007AFF
├── light:   #5AC8FA
└── dark:    #0051D5

Background
├── default: #F5F5F7
├── paper:   #FFFFFF
└── elevated: #FAFAFA

Text
├── primary:  #1C1C1E
├── secondary: #3A3A3C
└── disabled: #C7C7CC

Semantic
├── success: #34C759
├── warning: #FF9500
└── error:   #FF3B30
```

### Transições

```typescript
duration: {
  short: 220ms,
  standard: 300ms,
}

easing: {
  easeInOut: 'cubic-bezier(0.32, 0.72, 0, 1)',
}
```

### Componentes

```
Button
├── borderRadius: 999 (pill)
├── hover: scale(1.02)
└── shadow: soft diffuse

Card
├── borderRadius: 16
├── shadow: 0px 4px 24px rgba(0,0,0,0.10)
└── hover: lift + shadow increase

TextField
├── borderRadius: 12
└── focus: 2px border

Chip
├── borderRadius: 8
└── filled: alpha background (0.12)
```

---

## Dependency Graph

```
app/simulator/page.tsx
  ├── components/simulator/SimulatorForm.tsx
  │   ├── @mui/material (Box, Stack, TextField, Button)
  │   ├── @mui/icons-material/Search
  │   ├── zod (validation)
  │   └── types/simulator
  ├── components/simulator/DestinationList.tsx
  │   ├── @mui/material (Box, Stack, Skeleton)
  │   ├── @mui/icons-material/PublicOff
  │   └── components/simulator/DestinationCard.tsx
  │       ├── @mui/material (Card, Chip, LinearProgress, Tooltip)
  │       ├── @mui/icons-material (TrendingUp, AttachMoney, Public, LocalShipping)
  │       └── types/simulator
  ├── components/simulator/RateLimitBanner.tsx
  │   ├── @mui/material (Alert, AlertTitle, LinearProgress, Button)
  │   ├── @mui/icons-material (Info, Warning)
  │   └── types/simulator
  ├── components/simulator/UpgradeModal.tsx
  │   ├── @mui/material (Dialog, Button, Chip, Stack)
  │   ├── @mui/icons-material (RocketLaunch, CheckCircle, Close)
  │   └── types/simulator
  ├── components/simulator/ErrorState.tsx
  │   ├── @mui/material (Box, Stack, Alert, Button)
  │   ├── @mui/icons-material (ErrorOutline, SearchOff, CloudOff, WarningAmber)
  │   └── types/simulator
  └── lib/api/simulator.ts
      └── types/simulator
```

---

## Performance Optimizations

### 1. React Optimizations

- `React.memo()` em componentes que não mudam frequentemente
- `useCallback()` para funções passadas como props
- `useMemo()` para computações caras (futuro: scores breakdown)

### 2. MUI Optimizations

- Tree-shaking automático (imports específicos)
- `sx` prop (CSS-in-JS otimizado)
- Theme caching (singleton instance)

### 3. Next.js Optimizations

- App Router (React Server Components)
- Automatic code splitting
- Image optimization (Next/Image)
- Font optimization (next/font)

### 4. Loading Strategies

- Skeleton loading (perceived performance)
- Optimistic UI updates
- Error boundaries
- Suspense (futuro)

---

## Security Considerations

### 1. Input Validation

- **Client-side:** Zod validation (prevents bad requests)
- **Server-side:** Backend validation (security layer)

### 2. Rate Limiting

- **Freemium:** 5 requests/day per user
- **Tracked by:** IP address (backend)
- **Enforced by:** HTTP 429 response

### 3. CORS

- **Dev:** `Access-Control-Allow-Origin: *`
- **Prod:** Whitelist specific domains

### 4. Environment Variables

- **Public:** `NEXT_PUBLIC_API_BASE_URL` (safe to expose)
- **Private:** None (all backend secrets)

---

## Testing Strategy (Recommended)

### 1. Unit Tests (Vitest)

```typescript
// components/__tests__/SimulatorForm.test.tsx
describe('SimulatorForm', () => {
  it('validates NCM format', () => { ... });
  it('validates volume range', () => { ... });
  it('shows loading state', () => { ... });
});
```

### 2. Integration Tests (Playwright)

```typescript
// e2e/simulator.spec.ts
test('happy path simulation', async ({ page }) => {
  await page.goto('/simulator');
  await page.fill('input[name="ncm"]', '17011400');
  await page.fill('input[name="volume_kg"]', '1000');
  await page.click('button[type="submit"]');
  await expect(page.locator('.destination-card')).toHaveCount(10);
});
```

### 3. Visual Regression (Percy)

- Screenshot comparisons
- Mobile + Desktop viewports
- All states (loading, error, success)

### 4. Accessibility (axe-core)

```typescript
// a11y/simulator.test.ts
test('simulator page is accessible', async ({ page }) => {
  await page.goto('/simulator');
  const results = await injectAxe(page);
  expect(results.violations).toHaveLength(0);
});
```

---

## Monitoring & Observability

### 1. Frontend Metrics

- **Performance:** Web Vitals (FCP, LCP, CLS, FID)
- **Errors:** Sentry (error tracking)
- **Analytics:** Google Analytics / Mixpanel

### 2. API Metrics

- **Latency:** Avg response time (target: < 100ms)
- **Errors:** Error rate by status code
- **Rate Limits:** 429 occurrences (upgrade funnel)

### 3. User Metrics

- **Conversions:** Simulations → Upgrade clicks
- **Retention:** Daily active users
- **Satisfaction:** NPS score

---

## Conclusão

Esta arquitetura foi desenhada para:

1. **Escalabilidade:** Separação clara de responsabilidades
2. **Manutenibilidade:** Componentes reutilizáveis e tipados
3. **Performance:** Loading states e optimizations
4. **User Experience:** Apple aesthetic + accessibility
5. **Robustez:** Error handling em todas as camadas

**Status:** ✅ PRODUCTION READY
**Versão:** 1.0.0
**Documentação:** Completa
