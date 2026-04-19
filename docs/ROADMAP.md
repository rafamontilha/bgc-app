# BGC — Roadmap Ativo

**Versão:** 1.0
**Atualizado em:** 2026-04-19
**Contexto:** Projeto pessoal, dev solo (~10-15h/semana). Sem deadline comercial. Possível abertura ao mercado no futuro.

> Este documento substitui: `ROADMAP-EXECUTIVO-E2E.md`, `PRODUCT-ROADMAP.md`, `STRATEGIC-ROADMAP-E2E-INTEGRATED.md`, `NEXT-STEPS.md` (todos arquivados em `docs/archive/`).

---

## Estado Atual — v0.5.0 (2026-04-19)

### O que está pronto

| Componente | Status | Observações |
|------------|--------|-------------|
| Auth (J-AC02) | DONE | Clerk: email, Google OAuth, SSO, recuperação, JWT middleware Go |
| Onboarding wizard (J-AC01) | 87.5% | Days 1-7 entregues; Day 8 (polish + analytics stubs) pendente |
| Dashboard personalizado | DONE | Preview do simulador com NCM do usuário, tutorial, banner |
| Export Destination Simulator | DONE (dados parciais) | 50 países no schema; cobertura real: só cap. 17 (Açúcares) |
| Observabilidade | DONE | Prometheus, Grafana, Jaeger, structured logging |
| Integration Gateway | DONE | mTLS, OAuth2, API Key, Circuit Breaker |
| Infraestrutura | DONE | Docker Compose + k3d, HPA, backups, Network Policies Zero Trust |

### Dívida técnica conhecida

- **Clerk `publicMetadata` como storage de perfil**: não permite queries agregadas. Migração para PostgreSQL planejada em J-AC04.
- **Cobertura NCM**: 8% do total (só cap. 17). Simulador retorna erro para qualquer exportador fora do setor de açúcares.
- **Analytics**: sem instrumentação (PostHog/Mixpanel). Stubs de analytics do J-AC01 Day 8 ainda não criados.

---

## Priorização de Features

| Feature | Prioridade | Justificativa |
|---------|------------|---------------|
| J-AC01 Day 8 (polish + QA) | **P0** | 87.5% concluído — completar antes de iniciar qualquer coisa nova |
| Expansão NCM cap. 02 + 08 | **P0** | Simulador inútil para 92% dos exportadores sem dados reais |
| J-AC03 (página `/simulator`) | **P1** | Simulador só existe como preview no dashboard — produto incompleto |
| J-AC04 (Clerk → PostgreSQL) | **P1** | Dívida que cresce com o tempo; resolver enquanto base de usuários é zero |
| Analytics básico | **P3** | 1 usuário não precisa de funil; stub do Day 8 é suficiente por ora |
| Stripe/Paywall (J-REV01/02) | **Future** | Sem pressão de MRR; revisitar quando/se houver abertura ao mercado |
| Beta privado com N usuários | **Future** | Sem deadline; o produto precisa de dados antes de usuários externos |
| Marketplace de compradores | **Backlog distante** | Depende de cobertura NCM completa |

---

## Roadmap por Blocos

Cadência: dev solo, ~10-15h/semana. Unidade de tempo: semana/mês. Sem sprints intensivos de 8 dias.

---

### Bloco 1 — Fechamento e Limpeza
**Período:** 21/04/2026 → 04/05/2026 | **Estimativa total:** ~7h

**Objetivo:** Zerar dívida aberta antes de iniciar coisas novas.

| Tarefa | Estimativa | Entregável |
|--------|-----------|------------|
| J-AC01 Day 8: polish visual, analytics stubs (`lib/analytics/onboarding.ts`), QA em staging | ~4h | Épico J-AC01 fechado (DoD completo) |
| Reorganização documental (este replanning) | ~2h | `docs/archive/`, `docs/ROADMAP.md`, README atualizado |
| Atualizar `docs/QUICK-START.md` para v0.5.0 | ~1h | Quick start sem referências desatualizadas |

**Critério de conclusão do bloco:** J-AC01 marcado como DONE, zero TODOs abertos de antes de 2026-04-21.

---

### Bloco 2 — Expansão de Cobertura NCM
**Período:** 05/05/2026 → 01/06/2026 | **Estimativa total:** ~28h

**Objetivo:** Simulador que funciona para os 3 maiores setores exportadores do Brasil.

| Tarefa | Estimativa | Entregável |
|--------|-----------|------------|
| Ingestão de dados cap. 02 (Carnes) em `stg.exportacao` | ~10h | Simulador retorna resultados para NCMs de carne |
| Ingestão de dados cap. 08 (Frutas e Nozes) | ~10h | Simulador retorna resultados para NCMs de frutas |
| Validação E2E com os 3 capítulos (17, 02, 08) | ~4h | Testes passando, sem regressões |
| Release v0.6.0 | ~2h | CHANGELOG, tag, README atualizado |

**Critério de conclusão do bloco:** `POST /v1/simulator/destinations` retorna resultados válidos para NCMs dos caps. 02, 08 e 17. Cobertura sobe de 8% para ~28% das exportações brasileiras.

**Referência:** Ver `docs/DIAGNOSTICO-NCM-COVERAGE.md` para análise detalhada das lacunas.

---

### Bloco 3 — Simulador Completo (J-AC03)
**Período:** 02/06/2026 → 29/06/2026 | **Estimativa total:** ~16h

**Objetivo:** Simulador como produto real, não apenas preview no dashboard.

| Tarefa | Estimativa | Entregável |
|--------|-----------|------------|
| Página `/simulator` dedicada com UX completa | ~12h | Formulário completo: input NCM, filtro de países, volume |
| Cards de resultado com breakdown de score | ~3h | Score decomponível por fator (Market Size, Growth, Price, Distance) |
| Estado "NCM sem dados" com mensagem útil | ~1h | UX que orienta o usuário quando não há dados para o capítulo |

**Critério de conclusão do bloco:** Usuário consegue usar o simulador completo em `/simulator` sem passar pelo dashboard. Estado de erro para NCMs sem cobertura é amigável e informativo.

---

### Bloco 4 — Qualidade Técnica
**Período:** 30/06/2026 → 27/07/2026 | **Estimativa total:** ~24h

**Objetivo:** Eliminar dívida técnica acumulada antes de crescer.

| Tarefa | Estimativa | Entregável |
|--------|-----------|------------|
| J-AC04: migração Clerk `publicMetadata` → PostgreSQL `user_profiles` | ~16h | Perfil do exportador em PostgreSQL; queries agregadas possíveis |
| Dashboard Grafana com métricas do simulador (uso por NCM, latência) | ~4h | Dashboard operacional para monitorar uso |
| Revisão de testes unitários para caps. 02 e 08 | ~4h | Cobertura de testes consistente com os novos dados |

**Critério de conclusão do bloco:** Clerk `publicMetadata` não é mais o storage primário de perfil. Migration one-time executado com sucesso. `user_profiles` populado.

---

## Backlog Sem Data

Itens válidos mas sem prioridade ativa agora. Revisitar após Bloco 4.

- **Expansão NCM** para todos os 97 capítulos
- **Analytics básico** além dos stubs (evento de simulação, evento de onboarding)
- **Página de perfil avançada** (CNPJ, razão social, múltiplos NCMs)
- **Alertas de mercado** (notificações quando score de um destino muda significativamente)
- **Preparação para abertura ao mercado** — pricing research, Stripe, landing page de conversão

---

## Changelog do Roadmap

| Data | Versão | Alteração |
|------|--------|-----------|
| 2026-04-19 | 1.0 | Documento criado. Consolida `ROADMAP-EXECUTIVO-E2E.md`, `PRODUCT-ROADMAP.md`, `STRATEGIC-ROADMAP-E2E-INTEGRATED.md` e `NEXT-STEPS.md` (todos arquivados). Replanejamento para cadência de projeto pessoal, dev solo. |
