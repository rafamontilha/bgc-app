# Frontend Epic 4 - Export Destination Simulator
## Implementation Breakdown & Wireframes

**Status:** Ready for Development
**Estimated Effort:** 30 hours (~1 week with 1 full-time frontend developer)
**Target Completion:** Week 2 (13-17 January 2026)

---

## 📐 Wireframes

### Page: `/simulator`

```
┌────────────────────────────────────────────────────────────────┐
│  BGC - Simulador de Destinos de Exportação                    │
│  ═══════════════════════════════════════════════════════════  │
│                                                                │
│  Descubra os melhores países para exportar seu produto        │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  Digite o NCM do produto (8 dígitos) *                  │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │ 17011400                                           │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  │  ℹ️  Ex: 17011400 (Açúcar de cana)                      │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  Volume estimado (kg) - Opcional                        │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │ 1000                                               │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  Filtrar por países (opcional)                          │ │
│  │  ┌────────────────────────────────────────────────────┐ │ │
│  │  │ [v] Estados Unidos  [v] China     [ ] Alemanha    │ │ │
│  │  │ [ ] Argentina       [ ] Japão     [ ] México      │ │ │
│  │  │                                                    │ │ │
│  │  │ [+ Mostrar mais países]                           │ │ │
│  │  └────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  Número de resultados: [████████░░] 10                  │ │
│  │  (1-50 países)                                          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│             [ 🚀 Simular Destinos ]                           │
│                                                                │
│  ────────────────────────────────────────────────────────────│
│  💡 Você usou 2 de 5 simulações gratuitas hoje               │
│  [ ⭐ Upgrade para Premium - Ilimitado ]                      │
└────────────────────────────────────────────────────────────────┘
```

### Results View (After Submit)

```
┌────────────────────────────────────────────────────────────────┐
│  ← Voltar   Recomendações de Destinos (10 resultados)         │
│  ═══════════════════════════════════════════════════════════  │
│                                                                │
│  ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ │
│  ┃ #1  🇺🇸 Estados Unidos                    Score: 8.5/10┃ │
│  ┃                                                          ┃ │
│  ┃ Demanda: 🔥 ALTA    Market Size: USD 234M/ano          ┃ │
│  ┃                                                          ┃ │
│  ┃ 💰 Estimativas Financeiras:                             ┃ │
│  ┃   • Margem Estimada: 28%                                ┃ │
│  ┃   • Custo Logístico: USD 450                            ┃ │
│  ┃   • Tarifa: 12%                                         ┃ │
│  ┃   • Lead Time: 18 dias                                  ┃ │
│  ┃                                                          ┃ │
│  ┃ 📊 Breakdown do Score:                                  ┃ │
│  ┃   ┌─────────────────────────────────────────────────┐  ┃ │
│  ┃   │ Market Size (40%):  ████████░░ 9.2              │  ┃ │
│  ┃   │ Growth (30%):       ███████░░░ 7.8              │  ┃ │
│  ┃   │ Price (20%):        ████████░░ 8.5              │  ┃ │
│  ┃   │ Distance (10%):     ██████░░░░ 6.0              │  ┃ │
│  ┃   └─────────────────────────────────────────────────┘  ┃ │
│  ┃                                                          ┃ │
│  ┃ ✨ Por quê este destino?                                ┃ │
│  ┃ "Mercado grande e crescente, preços competitivos,      ┃ │
│  ┃  logística eficiente. Alto potencial de lucratividade."┃ │
│  ┃                                                          ┃ │
│  ┃             [ Ver Análise Completa → ]                  ┃ │
│  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ │
│                                                                │
│  ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ │
│  ┃ #2  🇨🇳 China                             Score: 7.9/10┃ │
│  ┃ (similar layout)                                        ┃ │
│  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ │
│                                                                │
│  ... (8 more cards) ...                                       │
│                                                                │
│  ────────────────────────────────────────────────────────────│
│  💡 Você usou 3 de 5 simulações gratuitas hoje               │
│  [ ⭐ Upgrade para Premium - Ilimitado ]                      │
└────────────────────────────────────────────────────────────────┘
```

---

## 🗂️ Component Breakdown

### Phase 1: Setup (2h)

**TASK 1.1: Project Setup**
- [ ] Create Next.js page `/simulator` with route configuration
- [ ] Setup TypeScript types for API request/response
- [ ] Configure environment variables (`NEXT_PUBLIC_API_BASE_URL`)
- [ ] Install dependencies:
  ```bash
  npm install axios react-query @tanstack/react-query
  npm install chart.js react-chartjs-2
  npm install react-hook-form zod @hookform/resolvers
  ```

**Deliverables:**
- `app/simulator/page.tsx` (empty scaffold)
- `types/simulator.ts` (TypeScript interfaces)
- `.env.local` with API URL

---

### Phase 2: Form Components (4h)

**TASK 2.1: NCM Input Component**
- [ ] Create `components/simulator/NCMInput.tsx`
- [ ] Validation: 8 digits, numeric only
- [ ] Input mask: auto-format as user types
- [ ] Error states with helpful messages
- [ ] Example NCMs suggestion below input

**Acceptance Criteria:**
- Accepts only 8-digit numeric input
- Shows validation error for invalid input
- Displays example NCMs (17011400, 26011200, 12010090)

---

**TASK 2.2: Volume Input Component**
- [ ] Create `components/simulator/VolumeInput.tsx`
- [ ] Optional numeric input (kg)
- [ ] Min: 1, no max
- [ ] Number formatting with thousand separators

**Acceptance Criteria:**
- Optional field (can be left blank)
- Accepts positive integers only
- Formats display with commas (e.g., 1,000)

---

**TASK 2.3: Country Filter Component**
- [ ] Create `components/simulator/CountryFilter.tsx`
- [ ] Multi-select checkboxes
- [ ] Show top 10 countries initially
- [ ] "Show more" button to expand to 50 countries
- [ ] Search/filter by country name

**Acceptance Criteria:**
- Allows multiple selections
- Collapses/expands properly
- Search filters list dynamically

---

**TASK 2.4: Max Results Slider**
- [ ] Create `components/simulator/MaxResultsSlider.tsx`
- [ ] Range slider: 1-50
- [ ] Default: 10
- [ ] Visual indicator of current value

**Acceptance Criteria:**
- Smooth slider interaction
- Updates value in real-time
- Displays current value prominently

---

**TASK 2.5: Submit Button & Form Validation**
- [ ] Create `components/simulator/SimulatorForm.tsx` (parent)
- [ ] Integrate all input components
- [ ] Form validation with `react-hook-form` + `zod`
- [ ] Loading state during submission
- [ ] Error handling for network failures

**Acceptance Criteria:**
- Form validates before submission
- Shows loading spinner during API call
- Displays error toast on failure
- Disabled state while loading

**Validation Schema (Zod):**
```typescript
const simulatorSchema = z.object({
  ncm: z.string().regex(/^\d{8}$/, "NCM deve ter 8 dígitos"),
  volume_kg: z.number().positive().optional(),
  countries: z.array(z.string()).optional(),
  max_results: z.number().min(1).max(50).default(10),
});
```

---

### Phase 3: Results Components (8h)

**TASK 3.1: Destination Card Component**
- [ ] Create `components/simulator/DestinationCard.tsx`
- [ ] Display all fields: rank, flag, name, score
- [ ] Demand indicator badge (Alto/Médio/Baixo)
- [ ] Financial metrics section
- [ ] Expand/collapse for score breakdown
- [ ] Recommendation reason text

**Acceptance Criteria:**
- Card is visually appealing (shadow, border, padding)
- All data fields display correctly
- Expand/collapse animation is smooth
- Responsive on mobile/tablet/desktop

---

**TASK 3.2: Score Breakdown Component**
- [ ] Create `components/simulator/ScoreBreakdown.tsx`
- [ ] Horizontal bar chart for 4 metrics
- [ ] Color-coded bars (market size, growth, price, distance)
- [ ] Percentage labels on each bar
- [ ] Uses `react-chartjs-2`

**Acceptance Criteria:**
- Chart renders correctly with accurate data
- Colors are distinct and accessible
- Tooltips show exact values on hover

---

**TASK 3.3: Demand Indicator Component**
- [ ] Create `components/simulator/DemandIndicator.tsx`
- [ ] Badge component with 3 states: Alto/Médio/Baixo
- [ ] Color coding: Alto (green), Médio (yellow), Baixo (red)
- [ ] Icon + text

**Acceptance Criteria:**
- Visual distinction clear
- Accessible (meets WCAG AA contrast)

---

**TASK 3.4: Financial Metrics Component**
- [ ] Create `components/simulator/FinancialMetrics.tsx`
- [ ] Grid layout for 4 metrics
- [ ] Icons for each metric
- [ ] Formatted values (currency, percentage, days)

**Acceptance Criteria:**
- Proper number formatting (USD 450, 28%, 18 dias)
- Icons are intuitive
- Responsive grid layout

---

**TASK 3.5: Results List Component**
- [ ] Create `components/simulator/ResultsList.tsx`
- [ ] Maps over destinations array
- [ ] Renders DestinationCard for each
- [ ] Virtualization for performance (react-window)
- [ ] Empty state if no results

**Acceptance Criteria:**
- Renders 50 cards without lag
- Empty state shows helpful message
- Smooth scrolling

---

**TASK 3.6: Empty State Component**
- [ ] Create `components/simulator/EmptyState.tsx`
- [ ] Displays when no destinations found
- [ ] Helpful message + suggestions
- [ ] Button to try again

**Acceptance Criteria:**
- User-friendly message
- Clear call-to-action

---

### Phase 4: Rate Limit UI (4h)

**TASK 4.1: Rate Limit Banner**
- [ ] Create `components/simulator/RateLimitBanner.tsx`
- [ ] Display "X of 5 simulations used today"
- [ ] Progress bar visual
- [ ] Always visible at bottom of page
- [ ] Reads from API response headers

**Acceptance Criteria:**
- Updates after each simulation
- Shows correct count
- Visual progress indicator

---

**TASK 4.2: Upgrade Modal**
- [ ] Create `components/simulator/UpgradeModal.tsx`
- [ ] Opens when user hits rate limit (429)
- [ ] Displays premium benefits
- [ ] CTA button to upgrade
- [ ] Close button

**Acceptance Criteria:**
- Triggers on 429 response
- Modal is accessible (keyboard nav, focus trap)
- Close button works

---

**TASK 4.3: Free Tier Indicator**
- [ ] Create `components/simulator/FreeTierIndicator.tsx`
- [ ] Small badge/chip showing "Free Tier"
- [ ] Always visible in header or sidebar
- [ ] Link to upgrade page

**Acceptance Criteria:**
- Non-intrusive
- Clear visual hierarchy

---

### Phase 5: API Integration (4h)

**TASK 5.1: API Client**
- [ ] Create `lib/api/simulator.ts`
- [ ] Axios instance with base URL
- [ ] Request/response interceptors
- [ ] Error handling utility
- [ ] TypeScript types for all endpoints

**API Client Functions:**
```typescript
export async function simulateDestinations(
  payload: SimulatorRequest
): Promise<SimulatorResponse> {
  // POST /v1/simulator/destinations
}

export async function getCountriesMetadata(): Promise<Country[]> {
  // GET /v1/countries (if endpoint exists)
}
```

---

**TASK 5.2: React Query Hooks**
- [ ] Create `hooks/useSimulator.ts`
- [ ] `useSimulateMutation` hook for POST request
- [ ] `useCountriesQuery` hook for fetching countries
- [ ] Automatic retry on failure
- [ ] Cache management

**Hooks:**
```typescript
export function useSimulateMutation() {
  return useMutation({
    mutationFn: simulateDestinations,
    onError: handleApiError,
  });
}

export function useCountriesQuery() {
  return useQuery({
    queryKey: ['countries'],
    queryFn: fetchCountries,
    staleTime: 1000 * 60 * 60, // 1 hour
  });
}
```

---

**TASK 5.3: Error Handling**
- [ ] Create `components/ErrorBoundary.tsx`
- [ ] Create `utils/errorHandling.ts`
- [ ] Map API errors to user-friendly messages
- [ ] Toast notifications for errors

**Error Types to Handle:**
- 400: Validation error
- 404: NCM not found
- 429: Rate limit exceeded
- 500: Server error

---

**TASK 5.4: Loading States**
- [ ] Create `components/LoadingSpinner.tsx`
- [ ] Skeleton loaders for cards
- [ ] Loading overlay during form submission

**Acceptance Criteria:**
- Smooth transitions between states
- Accessible loading indicators

---

**TASK 5.5: Success/Error Toasts**
- [ ] Integrate `react-hot-toast` or similar
- [ ] Success toast on successful simulation
- [ ] Error toasts for failures
- [ ] Auto-dismiss after 5 seconds

**Acceptance Criteria:**
- Toasts are non-blocking
- Clear, concise messages

---

### Phase 6: Polish & Responsiveness (4h)

**TASK 6.1: Responsive Design**
- [ ] Mobile viewport (320px-768px)
- [ ] Tablet viewport (768px-1024px)
- [ ] Desktop viewport (1024px+)
- [ ] Test on real devices

**Breakpoints:**
- Mobile: Single column, stacked form inputs
- Tablet: Two-column grid for results
- Desktop: Three-column grid for results

---

**TASK 6.2: Animations & Transitions**
- [ ] Smooth form transitions
- [ ] Card hover effects
- [ ] Modal slide-in animations
- [ ] Results fade-in stagger effect

**Acceptance Criteria:**
- Animations are subtle and performant
- No jank or lag

---

**TASK 6.3: Accessibility (WCAG AA)**
- [ ] ARIA labels on all interactive elements
- [ ] Keyboard navigation (Tab, Enter, Esc)
- [ ] Focus indicators visible
- [ ] Color contrast meets WCAG AA
- [ ] Screen reader testing

**Acceptance Criteria:**
- Passes axe-core audit
- Keyboard navigation works
- Screen reader announces correctly

---

**TASK 6.4: SEO Optimization**
- [ ] Meta tags (title, description, og:image)
- [ ] JSON-LD structured data
- [ ] Semantic HTML (h1, h2, nav, main, section)
- [ ] `robots.txt` allows indexing

**Meta Tags:**
```html
<title>Simulador de Destinos de Exportação | BGC</title>
<meta name="description" content="Descubra os melhores países para exportar seus produtos com base em dados reais de comércio exterior." />
<meta property="og:title" content="Simulador de Destinos - BGC" />
```

---

### Phase 7: Testing (4h)

**TASK 7.1: Unit Tests (Components)**
- [ ] Test NCMInput validation
- [ ] Test VolumeInput formatting
- [ ] Test CountryFilter selection
- [ ] Test DestinationCard rendering
- [ ] Test ScoreBreakdown data display

**Tools:** Jest + React Testing Library

---

**TASK 7.2: Integration Tests (Form Submission)**
- [ ] Test complete form submission flow
- [ ] Test error handling
- [ ] Test rate limiting UI update
- [ ] Mock API responses

---

**TASK 7.3: E2E Tests (Playwright)**
- [ ] Test happy path: submit form → see results
- [ ] Test validation errors
- [ ] Test rate limit flow
- [ ] Test responsive design

**Scenarios:**
1. User enters NCM, submits, sees results
2. User enters invalid NCM, sees error
3. User hits rate limit, sees upgrade modal

---

**TASK 7.4: Accessibility Tests**
- [ ] Run axe-core audit
- [ ] Keyboard navigation test
- [ ] Screen reader test (NVDA/JAWS)
- [ ] Color contrast check

---

## 📦 Dependencies

```json
{
  "dependencies": {
    "next": "15.x",
    "react": "19.x",
    "react-dom": "19.x",
    "axios": "^1.6.0",
    "@tanstack/react-query": "^5.0.0",
    "react-hook-form": "^7.48.0",
    "zod": "^3.22.0",
    "@hookform/resolvers": "^3.3.0",
    "chart.js": "^4.4.0",
    "react-chartjs-2": "^5.2.0",
    "react-hot-toast": "^2.4.0",
    "react-window": "^1.8.0"
  },
  "devDependencies": {
    "@testing-library/react": "^14.0.0",
    "@testing-library/jest-dom": "^6.1.0",
    "@playwright/test": "^1.40.0",
    "@axe-core/react": "^4.8.0",
    "jest": "^29.7.0",
    "typescript": "^5.3.0"
  }
}
```

---

## 🎯 Definition of Done

A feature is considered "Done" when:

1. ✅ Code is written and follows style guide
2. ✅ All unit tests pass
3. ✅ Integration tests pass
4. ✅ E2E test covers critical path
5. ✅ Code reviewed and approved
6. ✅ Accessibility audit passes (axe-core)
7. ✅ Works on Chrome, Firefox, Safari
8. ✅ Responsive on mobile, tablet, desktop
9. ✅ Deployed to staging and tested
10. ✅ Documented in README or docs/

---

## 📅 Timeline

| Week | Tasks | Completion % |
|------|-------|--------------|
| **Week 2 (Jan 13-17)** | Phase 1-5 | 80% |
| **Week 3 (Jan 20-24)** | Phase 6-7 + Polish | 100% |

**Milestones:**
- **Day 3 (Jan 15):** Form components complete, can submit to API
- **Day 5 (Jan 17):** Results display complete, end-to-end flow works
- **Day 7 (Jan 22):** Polish complete, all tests passing
- **Day 10 (Jan 24):** Deployed to production

---

## 🚦 Risk Mitigation

| Risk | Mitigation |
|------|------------|
| API changes during development | Use TypeScript contracts, sync with backend team |
| Rate limit testing difficult locally | Mock rate limit headers in dev mode |
| Performance issues with 50 cards | Use virtualization (react-window) |
| Accessibility not prioritized | Audit with axe-core after each phase |
| Design changes mid-development | Get design approval before Phase 2 |

---

## ✅ Success Criteria

The frontend is successful if:

1. **Functional:** User can simulate destinations and see results
2. **Performant:** Page load < 2s, interaction < 100ms
3. **Accessible:** Passes WCAG AA audit
4. **Responsive:** Works on mobile, tablet, desktop
5. **Tested:** 80%+ test coverage
6. **User-Friendly:** Intuitive UX, helpful error messages
7. **Deployed:** Live in production and stable

---

**Last Updated:** 2026-01-09
**Version:** 1.0
**Author:** Product Management Team
