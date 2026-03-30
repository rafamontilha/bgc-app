# Progress Report — J-AC01: Onboarding First-Time User (Days 1-7)

**Epic:** J-AC01 — Onboarding de Novo Usuario (First-Time Experience)
**Prioridade:** P0 (MVP Beta Blocker)
**Status:** Em Andamento — 87.5% (7/8 dias concluidos)
**Responsavel:** BGC Engineering + Product
**Data do Relatorio:** 2026-03-29
**Proximo Marco:** Day 8 — 2026-04-06 (polish visual, analytics stubs, QA staging, DoD final)

---

## Resumo Executivo

A sprint J-AC01 atingiu 87.5% de completude com 7 dos 8 dias entregues. Todas as 8 User Stories estao implementadas e cobertas por testes. A suite TDD conta com 71 testes passando (0 falhas, 0 pendentes) em 8 suites. O TypeScript esta sem erros em todo o projeto. Tres bugs foram identificados e corrigidos durante o processo de desenvolvimento.

O impacto mais direto no North Star Metric (Time-to-First-Export-Match) esta operacional: o fluxo completo — wizard NCM -> `publicMetadata` -> `DashboardSimulatorPreview` pre-carregado -> link para simulacao completa — esta funcionando end-to-end.

O Day 8 (pendente) e dedicado a polish visual, criacao dos stubs de analytics, QA manual com checklist completo em staging, e atualizacao da documentacao final.

---

## Artefatos Entregues por Dia

### Day 1 — Tipos e Dados Estaticos (2026-03-28)

**Objetivo atingido:** Estrutura de dados e tipos TypeScript prontos como base para todos os demais componentes.

| Arquivo | Descricao |
|---------|-----------|
| `web-next/lib/types/onboarding.ts` | Interfaces TypeScript: `NcmChapter`, `TradeRegion`, `OnboardingWizardState`, `OnboardingMetadata` |
| `web-next/lib/data/ncm-chapters.ts` | 96 capitulos NCM ativos com `defaultNcm8d` (8 digitos para compatibilidade com API Go) |
| `web-next/lib/data/trade-regions.ts` | 7 continentes + 35 paises (parceiros comerciais top do Brasil — Comex Stat) |

**Decisao tomada no Day 1 (resolucao do Risco R03):** O campo `defaultNcm8d` nao estava no plano original. Foi adicionado ao tipo `NcmChapter` para permitir que o frontend envie 8 digitos para o endpoint Go sem precisar alterar o backend. Isso resolveu preventivamente o risco R03 (endpoint exige 8 digitos) e eliminou qualquer dependencia de mudanca no servico `api`.

---

### Day 2 — Componentes Steps 1 e 2 (2026-03-28)

| Arquivo | Descricao |
|---------|-----------|
| `web-next/components/onboarding/OnboardingStep1Ncm.tsx` | Autocomplete MUI de capitulo NCM — busca por codigo e descricao, botao "Proximo" desabilitado sem selecao |
| `web-next/components/onboarding/OnboardingStep2Volume.tsx` | Input numerico de volume mensal com quick-select de 4 faixas pre-definidas, campo opcional |

---

### Day 3 — Componentes Step 3 e Welcome (2026-03-28)

| Arquivo | Descricao |
|---------|-----------|
| `web-next/components/onboarding/OnboardingStep3Regions.tsx` | Multi-select de regioes/paises agrupados por continente, chips de selecao, limite de 10 itens |
| `web-next/components/onboarding/OnboardingWelcome.tsx` | Tela de conclusao do wizard com resumo do perfil configurado |

---

### Day 4 — Orquestracao do Wizard (2026-03-28)

| Arquivo | Descricao |
|---------|-----------|
| `web-next/app/api/onboarding/route.ts` | API Route server-side para persistencia no Clerk `publicMetadata` — decisao de seguranca vs. client-side original |
| `web-next/app/onboarding/page.tsx` | Wizard completo: stepper MUI, navegacao entre steps, logica de skip, save com loading state, redirecionamento pos-save |

**Decisao arquitetural (desvio justificado):** O plano original previa chamada client-side direta ao Clerk SDK para `user.update({ publicMetadata: {...} })`. A implementacao optou por uma API Route server-side para isolar a chave secreta do Clerk no servidor, evitando exposicao no bundle do cliente. O comportamento do ponto de vista do usuario e identico. A migracao para PostgreSQL em J-AC04 sera mais simples a partir de uma API Route existente do que de uma chamada client-side.

---

### Day 5 — Integracao no Dashboard (2026-03-29)

| Arquivo | Descricao |
|---------|-----------|
| `web-next/components/dashboard/DashboardSimulatorPreview.tsx` | Card personalizado com NCM do onboarding, top destinos pre-simulados, estados de loading e erro |
| `web-next/app/dashboard/page.tsx` | Le `publicMetadata` do Clerk, exibe `DashboardSimulatorPreview` se onboarding completo, ou prompt CTA se nao completo |
| `web-next/app/simulator/page.tsx` | Aceita query param `?ncm=` para pre-fill do NCM + auto-executa simulacao no mount via `useEffect` |
| `web-next/components/simulator/SimulatorForm.tsx` | Prop `initialNcm` adicionada para compatibilidade com pre-fill externo |

**Resultado para o NSM:** O ciclo completo esta operacional. Apos o onboarding, o usuario chega ao dashboard com os resultados ja carregados para o NCM selecionado, e um link direto leva ao simulador completo pre-preenchido. Isso elimina o re-trabalho manual que representava ~5 minutos do tempo total estimado de ~30 minutos.

---

### Day 6 — Tutorial e Banner de Lembrete (2026-03-29)

| Arquivo | Descricao |
|---------|-----------|
| `web-next/components/onboarding/OnboardingTutorial.tsx` | Modal com 4 slides sequenciais apresentando as principais funcionalidades; flag `TUTORIAL_SEEN_KEY` no localStorage para nao reexibir |
| `web-next/components/onboarding/OnboardingReminderBanner.tsx` | Banner dismissivel para usuarios que pularam o onboarding; flag `REMINDER_BANNER_DISMISSED_KEY` no localStorage |

**Desvio de implementacao documentado:** O plano original especificava 3 tooltips `Popover` MUI ancorados em elementos especificos do DOM do dashboard. A implementacao optou por um modal com 4 slides para mitigar o Risco R04 (coordenadas incorretas no DOM apos hidratacao SSR no Next.js 15). A decisao foi tomada durante o desenvolvimento, apos identificar que o ancoramente de Popovers em elementos renderizados pelo servidor seria fragil sem um mecanismo de deferimento. O valor educacional entregue pelo modal e equivalente ao dos tooltips.

---

### Day 7 — Re-onboarding, Profile e TDD (2026-03-29)

**Re-onboarding:**
- `web-next/app/onboarding/page.tsx` — modo re-onboarding via `?re=1`: pre-preenche os 3 steps com os dados existentes do `publicMetadata`, incrementa `onboardingVersion` ao salvar

**Perfil:**
- `web-next/app/profile/page.tsx` — nova secao "Configuracoes de Exportacao" exibindo NCM atual, volume e regioes cadastradas, com botao "Alterar" que navega para `/onboarding?re=1`

**Desvio menor documentado:** O query param de re-onboarding foi implementado como `?re=1` ao inves de `?redo=true` conforme especificado. O comportamento e identico — a escolha foi por brevidade na URL.

**Testes TDD adicionados (51 novos):**

| Suite | Testes | O que cobre |
|-------|--------|-------------|
| `ncm-chapters` | 13 | Estrutura e integridade dos 96 capitulos, campo `defaultNcm8d`, ausencia de duplicatas |
| `trade-regions` | 14 | Estrutura de continentes/paises, relacao pai-filho, codigos ISO, ausencia de duplicatas |
| `api/onboarding` | 6 | API Route: persistencia correta, validacao de payload, tratamento de erro do Clerk SDK |
| `OnboardingReminderBanner` | 6 | Exibicao condicional por status de onboarding, dismiss, persistencia da flag no localStorage |
| `OnboardingTutorial` | 7 | Navegacao de slides, exibicao da flag localStorage, fechamento via botao e tecla Esc |
| `re-onboarding` | 4 | Pre-fill dos 3 steps com dados existentes, incremento de `onboardingVersion`, comportamento sem dados previos |

---

## Bugs Corrigidos

| # | Bug | Componente | Causa Raiz | Correcao Aplicada |
|---|-----|------------|------------|-------------------|
| 1 | React key spread warning | `OnboardingStep1Ncm.tsx` — prop `renderOption` do `Autocomplete` MUI | A prop `key` estava sendo passada via object spread (`{...props}`) junto com outras props do item, o que gera aviso do React sobre keys em arrays | Extrai `key` explicitamente no destructuring antes de usar spread: `const { key, ...itemProps } = props` |
| 2 | Smart CAPTCHA: widget element ausente | `app/login/page.tsx` e `app/signup/page.tsx` | O widget do CAPTCHA buscava um elemento DOM pelo ID no momento do render, mas o elemento nao existia porque nao estava declarado no JSX | Adiciona elemento `<div id="clerk-captcha" />` no JSX das paginas de auth |
| 3 | Container `bgc_api` sem politica de restart | `bgcstack/docker-compose.yml` | O servico `bgc_api` nao tinha `restart: unless-stopped`, causando downtime permanente em qualquer crash do processo Go | Adiciona `restart: unless-stopped` ao servico `bgc_api` — alinhado com os demais servicos do stack |

---

## Metricas de Qualidade

| Metrica | Valor | Status |
|---------|-------|--------|
| Testes totais | 71 | Passando (0 falhas) |
| Suites de teste | 8 | Todas verdes |
| Erros TypeScript | 0 | `npx tsc --noEmit` limpo |
| Bugs encontrados | 3 | Todos corrigidos |
| Regressoes em features existentes | 0 | Validado pela suíte pre-existente |
| Cobertura de features criticas | Alta | Dados NCM, regioes, API route, localStorage, re-onboarding todos cobertos |

---

## Impacto no North Star Metric

**NSM: Time-to-First-Export-Match**

O objetivo central de J-AC01 e reduzir o tempo entre o cadastro e o primeiro resultado significativo de simulacao, viabilizando a ativacao do usuario.

### Antes de J-AC01

```
Signup/Login
    |
    v
/dashboard (sem contexto de perfil)
    |
    v  (~10 min) Descobrir e entender NCM por conta propria
    |
    v  (~5 min)  Preencher o simulador manualmente
    |
    v  (~15 min) Interpretar resultados sem contexto personalizado
    |
    v
TOTAL: ~30 minutos ate o primeiro resultado relevante
```

### Depois de J-AC01 (Days 1-7 operacionais)

```
Signup/Login
    |
    v
/onboarding — wizard 3 passos
    Step 1: NCM via autocomplete        (~1-2 min)
    Step 2: Volume mensal (opcional)    (~30 seg)
    Step 3: Regioes-alvo (opcional)     (~30 seg)
    |
    v  Salva em Clerk publicMetadata (server-side)
    |
    v
/dashboard
    DashboardSimulatorPreview pre-carregado com NCM do perfil (~30 seg de carregamento)
    Tutorial modal 4 slides (educacao contextual)           (~1-2 min)
    |
    v
/simulator?ncm=XX (auto-executado, resultados completos)    (~30 seg)
    |
    v
TOTAL: ~4-5 minutos ate o primeiro resultado relevante
```

**Reducao estimada: 83-87%** (de ~30 min para ~4-5 min).

### Mecanismos que viabilizam a reducao

1. **Autocomplete guiado de NCM (Day 1-2):** Elimina o tempo de descoberta do sistema de nomenclatura. O usuario busca por descricao em portugues — "cafe", "soja", "calcados" — e o sistema retorna o capitulo correto automaticamente.

2. **Pre-fill no simulador (Day 5):** O NCM salvo no `publicMetadata` e usado diretamente no `?ncm=` do simulator. O usuario nao precisa repetir o que ja informou no wizard.

3. **DashboardSimulatorPreview (Day 5):** O primeiro resultado de simulacao aparece no proprio dashboard, sem precisar navegar para outra pagina. O usuario ve valor imediatamente.

4. **Tutorial contextual (Day 6):** Reduz o tempo de interpretacao dos resultados ao explicar o que cada secao significa no momento exato de primeiro uso.

5. **Banner de lembrete (Day 6):** Recupera usuarios que pularam o onboarding, mantendo a reducao de Time-to-First-Export-Match acessivel mesmo para o fluxo de skip.

---

## Pendente — Day 8 (2026-04-06)

| Item | Tipo | Estimativa |
|------|------|------------|
| Polish visual: espacamentos, bordas, tipografia nos componentes do wizard | UX | 2-3h |
| Consistencia de design system: verificar uso correto de tokens MUI v7 | UX | 1-2h |
| `web-next/lib/analytics/onboarding.ts` — stubs de eventos AARRR | Codigo | 1h |
| QA manual com checklist completo da secao 8.3 do sprint plan | QA | 2-3h |
| Validacao em staging | Infra | 1h |
| `USER-JOURNEYS-INDEX.md` — atualizar J-AC01 para DONE | Doc | 15min |
| `PRODUCT-DECISIONS.md` — registrar decisao de persistencia Clerk vs PostgreSQL | Doc | 30min |
| `RELEASE-NOTES-J-AC01.md` — breve para stakeholders | Doc | 30min |
| Atualizacao do `SPRINT-PLAN-J-AC01.md` — DoD final com todos checkboxes | Doc | 15min |

**Total estimado:** 8-11 horas de trabalho (dia completo).

---

## Proximos Passos Pos-Sprint

Apos a conclusao do Day 8 e merge na `main`:

1. **J-AC03 — Simulador Logged-In User (P0, 6d estimados):** Com o onboarding entregando NCM e contexto de perfil, o simulador pode ser aprimorado para usuarios autenticados — historico de simulacoes, favoritos, comparacoes.

2. **J-AC04 — Dashboard Intelligence (P1, 10d estimados):** Migracao do `publicMetadata` do Clerk para PostgreSQL `user_profiles`, conforme decisao arquitetural registrada. Habilita queries agregadas e personalização profunda.

3. **J-REV01 — Descoberta Premium / Paywall (P0, 6d estimados):** Com usuarios ativados via onboarding, o momento de apresentar o upgrade e natural — quando o usuario tenta ir alem dos limites do freemium.

---

**Documento criado em:** 2026-03-29
**Owner:** BGC Product Management
**Sprint Plan de referencia:** `docs/SPRINT-PLAN-J-AC01.md`
**Proximo relatorio:** 2026-04-06 (Day 8 concluido — DoD final)
