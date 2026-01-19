# BGC Roadmap Executivo com Testes E2E - Sumário Executivo

**Data:** 2026-01-16
**Versão:** 3.0
**Status:** PRONTO PARA EXECUÇÃO

---

## TL;DR - Decisões Críticas

### Contexto Atual
- **Backend:** 70% production-ready (gaps críticos identificados)
- **Frontend:** Estável (conflito Tailwind a resolver)
- **Dados:** 8% cobertura (BLOQUEADOR para beta)
- **Dívida Técnica:** 51 arquivos temporários
- **Testes E2E:** Existentes mas não integrados em CI/CD

### Decisão Estratégica
**Abordagem:** Qualidade First - Nenhuma fase avança sem testes E2E validados

### Timeline Executiva
- **Hoje (30min):** Limpeza técnica
- **Amanhã-Depois (1.5 dias):** Backend 100% production-ready
- **Próxima Semana (3-4 dias):** Dados 8% → 28%
- **Semana 2 (2 dias):** Autenticação Clerk
- **Semanas 3-4:** Beta Privado (20 usuários)
- **Paralelo:** Kubernetes (2 semanas)

---

## Roadmap em 6 Fases (Com E2E Obrigatório)

### Phase 0: Limpeza Imediata
**Tempo:** 30 minutos
**Owner:** Platform Team
**Bloqueio:** P0 (bloqueia tudo)

#### Tarefas
1. Deletar 51 arquivos `tmpclaude-*`
2. Atualizar `.gitignore`
3. Validar serviços sobem

#### Testes E2E Obrigatórios
```bash
bash tests/e2e/phase0_cleanup.sh
```
**Critérios:** Zero temp files, all services healthy

#### Definition of Done
- [ ] `git status` limpo
- [ ] Smoke tests passam
- [ ] Commit: "chore: cleanup temp files"

---

### Phase 1: Backend Production-Ready
**Tempo:** 10 horas (1.5 dias)
**Owner:** Backend Team
**Bloqueio:** P0 (bloqueia beta)

#### Gaps Críticos a Resolver
1. **Graceful Shutdown** (2h) - API não trata SIGTERM corretamente
2. **Structured Logging** (3h) - Logs inconsistentes, dificulta debug
3. **Database Migrations** (3h) - Processo manual, error-prone
4. **Health Check** (1h) - Não valida dependências (DB, Redis)
5. **Rate Limiting** (1h) - In-memory, não funciona com múltiplas instâncias

#### Testes E2E Obrigatórios
```bash
cd api
go test -tags=e2e -v ./tests/e2e/backend_readiness_test.go
```

**Test Matrix (7 testes):**
| Teste | Target | Critério de Sucesso |
|-------|--------|-------------------|
| Health Check | < 50ms | Valida DB + Redis |
| Graceful Shutdown | < 30s | Request completa durante shutdown |
| Structured Logs | N/A | JSON format com request_id |
| Migration Idempotency | < 10s | Safe para rodar 3x |
| Distributed Rate Limit | < 100ms | Funciona em 3 instâncias |
| Database Failover | < 5s | Recupera de restart |
| Performance | < 200ms | P95 latency |

#### Definition of Done
- [ ] Todos 7 testes E2E passam
- [ ] Code review aprovado
- [ ] Latência P95 < 200ms
- [ ] Graceful shutdown < 30s
- [ ] Deploy staging validado

---

### Phase 2: Expansão Dados 8% → 28%
**Tempo:** 3-4 dias
**Owner:** Data Engineering
**Bloqueio:** P0 (bloqueia beta)

#### Por Que Isso Importa
**Problema:** 92% dos usuários recebem erro "dados não disponíveis"
**Solução:** Adicionar capítulos 02 (Carnes) e 08 (Frutas)
**Impacto:** Bounce rate 80% → 40%, NPS 10 → 40

#### Tarefas
1. **Auditoria de Dados** (1 dia)
   - Verificar `stg.exportacao` atual
   - Identificar top 20 NCMs por capítulo
   - Validar qualidade de dados

2. **ETL Pipeline** (2 dias)
   - Job Go para popular capítulos 02 e 08
   - Kubernetes Job definition
   - Testes de idempotência

3. **Validação API** (1 dia)
   - Testar 6 NCMs (3 carnes, 3 frutas)
   - Performance < 200ms
   - Regression test (cap 17 continua funcionando)

#### Testes E2E Obrigatórios
```bash
bash tests/e2e/phase2_data_expansion.sh
```

**Test Matrix:**
| NCM | Produto | Min Destinos | Top País Esperado |
|-----|---------|-------------|------------------|
| 02013000 | Carne bovina | 5 | China (CN) |
| 02023000 | Carne congelada | 4 | EUA (US) |
| 02071200 | Frango | 6 | Japão (JP) |
| 08011100 | Cocos | 3 | - |
| 08012200 | Castanha-do-pará | 4 | - |
| 08030019 | Bananas | 5 | - |

#### Definition of Done
- [ ] Cap 02 com 500+ registros
- [ ] Cap 08 com 500+ registros
- [ ] Zero duplicatas
- [ ] Todos 6 NCMs retornam resultados
- [ ] P95 latency < 200ms
- [ ] Frontend atualizado com novos NCMs

---

### Phase 3: Autenticação Clerk
**Tempo:** 2 dias
**Owner:** Frontend + Backend
**Bloqueio:** P1 (necessário para monetização)

#### Por Que Isso Importa
**Problema:** Rate limiting por IP (fácil de burlar)
**Solução:** User accounts com tiers (free/premium)
**Impacto:** Viabiliza monetização, analytics de usuários

#### Tarefas
1. **Frontend Clerk** (4h)
   - Integração @clerk/nextjs
   - Sign in/up flow
   - Protected routes

2. **Backend JWT** (4h)
   - Validação JWT middleware
   - Rate limit por user_id
   - Premium users ilimitados

3. **Database Schema** (2h)
   - Tabela `users` com tiers
   - Daily counter reset
   - Índices otimizados

#### Testes E2E Obrigatórios
```bash
cd web-next
pnpm exec playwright test tests/e2e/auth.spec.ts
```

**Test Scenarios:**
1. Sign up completo
2. Free user atinge rate limit (5 simulações)
3. Premium user ilimitado (100+ simulações)
4. Protected routes redirecionam
5. Sign out funciona

#### Definition of Done
- [ ] Auth flow completo funciona
- [ ] Rate limit por user_id
- [ ] Premium = ilimitado
- [ ] Playwright tests passam
- [ ] Anonymous users ainda funcionam (backward compatibility)

---

### Phase 4: Beta Privado
**Tempo:** 2-3 semanas
**Owner:** Product Manager + Growth
**Bloqueio:** P0 (go-to-market milestone)

#### Objetivos
- **Recrutar:** 50 signups, 20 usuários ativos
- **Validar:** Product-market fit, NPS > 40
- **Aprender:** Feedback qualitativo, feature requests

#### Métricas de Sucesso
| Métrica | Target | Como Medir |
|---------|--------|-----------|
| Beta Signups | 50 | Analytics |
| Active Users | 20 (40%) | WAU/MAU > 60% |
| NPS | > 40 | Pesquisa semanal |
| Completion Rate | > 70% | Simulações finalizadas / iniciadas |
| Time to First Value | < 2 min | Signup → primeira simulação |

#### Testes E2E Obrigatórios
```bash
# Critical user journey
cd web-next
pnpm exec playwright test tests/e2e/user_journey.spec.ts

# Load test (20 concurrent beta users)
k6 run tests/e2e/load_test.js --vus 20 --duration 10m
```

**Critical Path Test (12 passos):**
1. Signup → 2. Invite link → 3. Onboarding → 4. Primeira simulação →
5. Ver resultados → 6. Salvar simulação → 7. Segunda simulação →
8. Dashboard → 9. Ver simulações salvas → 10. Feedback → 11. Rate limit →
12. Upgrade modal

#### Definition of Done
- [ ] 20 usuários ativos
- [ ] NPS > 40
- [ ] Uptime > 99.5%
- [ ] P95 latency < 200ms
- [ ] 10+ entrevistas feitas
- [ ] Critical path test passa

---

### Phase 5: Kubernetes Production (Paralelo)
**Tempo:** 2 semanas (roda junto com Phase 4)
**Owner:** DevOps Team
**Bloqueio:** P1

#### Objetivos
- Deploy production com HA
- Auto-scaling 2-10 pods
- Zero-downtime deployments
- Blue-green deployment

#### Testes E2E Obrigatórios
```bash
bash tests/e2e/k8s_cluster_test.sh
bash tests/e2e/blue_green_test.sh

# Load test production
k6 run tests/e2e/load_test.js --vus 100 --duration 15m
```

#### Definition of Done
- [ ] K8s cluster rodando (3 AZs)
- [ ] ArgoCD syncing
- [ ] Blue-green deployment testado
- [ ] Rollback < 5 minutos
- [ ] Load test 100 users passa
- [ ] Monitoring completo

---

## CI/CD Pipeline (Integração Completa)

### GitHub Actions Workflow

**Triggers:**
- Every push to `main`/`develop`
- Every pull request
- Scheduled: every 6 hours
- Manual: pre-deployment

**Jobs:**
1. **Backend E2E** (15 min)
2. **Frontend E2E** (20 min)
3. **Integration E2E** (30 min)
4. **Load Test** (10 min) - only on `main`

**Quality Gates:**
- All E2E tests must pass
- P95 latency < 200ms
- Error rate < 1%
- No critical security vulnerabilities

### Deployment Pipeline

```yaml
Pre-Deploy → E2E Tests → Deploy Green → Smoke Tests → Load Test → Switch Traffic → Cleanup
```

**Rollback:** Automatic if any step fails

---

## Métricas de Sucesso (Overall)

### Technical Health

| Métrica | Atual | Phase 1 | Phase 2 | Phase 4 | Target |
|---------|-------|---------|---------|---------|--------|
| E2E Coverage | 0% | 30% | 50% | 85% | > 90% |
| Test Reliability | N/A | 80% | 90% | 98% | > 95% |
| API P95 Latency | - | 200ms | 180ms | 120ms | < 200ms |
| Uptime | - | 99% | 99.5% | 99.8% | > 99.5% |
| Error Rate | - | 1% | 0.5% | 0.1% | < 0.5% |
| MTTR | - | 2h | 1h | 15min | < 30min |

### Product Health

| Métrica | Atual | Phase 2 | Phase 4 | Target (v1.0) |
|---------|-------|---------|---------|--------------|
| Data Coverage | 8% | 28% | 28% | 100% |
| Active Users | 0 | 0 | 20 | 100 |
| NPS | N/A | N/A | 40 | > 60 |
| Completion Rate | N/A | N/A | 70% | > 80% |
| Free → Premium | N/A | N/A | 3% | > 5% |

---

## Riscos Críticos e Mitigações

### Risco 1: Testes E2E Flaky em CI
**Probabilidade:** Média | **Impacto:** Alto
**Mitigação:**
- Retry logic (3 tentativas)
- Test isolation (containers isolados)
- Timeout generosos (não agressivos)
- Monitorar flakiness rate semanalmente

### Risco 2: Performance Degrada com Mais Dados
**Probabilidade:** Média | **Impacto:** Médio
**Mitigação:**
- Índices database otimizados
- Cache Redis agressivo
- Load tests em cada fase
- Query profiling contínuo

### Risco 3: Beta com Baixo Engagement (< 40%)
**Probabilidade:** Média | **Impacto:** Alto
**Mitigação:**
- Onboarding personalizado
- Outreach proativo semanal
- Feedback loop rápido (< 48h)
- Ajustes rápidos baseado em dados

### Risco 4: Deploy K8s Quebra Produção
**Probabilidade:** Baixa | **Impacto:** Crítico
**Mitigação:**
- Testar em staging primeiro
- Blue-green deployment (zero downtime)
- Rollback < 5 minutos
- Smoke tests obrigatórios pós-deploy

---

## Priorização (Framework RICE)

| Feature | Reach | Impact | Confidence | Effort | RICE | Priority |
|---------|-------|--------|------------|--------|------|----------|
| **Backend Ready** | 1000 | 3 | 0.9 | 1.5d | **1800** | **P0** |
| **Data 28%** | 1000 | 3 | 0.8 | 4d | **600** | **P0** |
| **Auth** | 500 | 3 | 0.7 | 2d | **525** | **P1** |
| **Beta** | 50 | 3 | 0.6 | 3w | **30** | **P0** |
| **K8s** | 1000 | 2 | 0.8 | 2w | **114** | **P1** |

**Legenda:**
- **Reach:** Usuários impactados em 3 meses
- **Impact:** 1 (Low), 2 (Med), 3 (High)
- **Confidence:** 0.0-1.0
- **Effort:** Tempo de dev

---

## Quick Start - Próximos Passos

### Hoje (2026-01-16) - 30 minutos
```bash
# 1. Rodar teste de limpeza
bash tests/e2e/phase0_cleanup.sh

# 2. Se passar, deletar temp files
find . -name "tmpclaude-*" -exec rm -rf {} +
find . -name "nul" -delete
find . -name "NUL" -delete

# 3. Atualizar .gitignore
echo "tmpclaude-*" >> .gitignore
echo "nul" >> .gitignore
echo "NUL" >> .gitignore

# 4. Commit
git add .
git commit -m "chore: cleanup temp files and update gitignore"
git push
```

### Amanhã (2026-01-17) - Manhã
```bash
# Começar Phase 1: Graceful Shutdown
cd api
# Implementar código conforme docs/STRATEGIC-ROADMAP-E2E-INTEGRATED.md
# Section: Phase 1, Task 1.1

# Rodar teste
go test -tags=e2e -v ./tests/e2e/backend_readiness_test.go -run TestGracefulShutdown
```

### Amanhã (2026-01-17) - Tarde
```bash
# Continuar Phase 1: Structured Logging
# Implementar código conforme roadmap, Task 1.2

# Rodar teste
go test -tags=e2e -v ./tests/e2e/backend_readiness_test.go -run TestStructuredLogging
```

### Dia 3 (2026-01-18)
```bash
# Finalizar Phase 1: Migrations + Health + Rate Limit
# Rodar suite completa
go test -tags=e2e -v ./tests/e2e/backend_readiness_test.go

# Se tudo passar, marcar Phase 1 como DONE
# Começar Phase 2: Data expansion
```

---

## Comandos Úteis (Cheat Sheet)

### Executar Testes

```bash
# Smoke tests (2 min)
bash tests/e2e/smoke_tests.sh

# Phase 0 (2 min)
bash tests/e2e/phase0_cleanup.sh

# Phase 1 (15 min)
cd api && go test -tags=e2e -v ./tests/e2e/backend_readiness_test.go

# Phase 2 (10 min)
bash tests/e2e/phase2_data_expansion.sh

# Phase 3 (20 min)
cd web-next && pnpm exec playwright test tests/e2e/auth.spec.ts

# Load test (10 min)
k6 run tests/e2e/load_test.js --vus 50 --duration 5m

# Full suite (60 min)
./run_all_e2e.sh
```

### Troubleshooting

```bash
# Ver logs
docker logs bgc_api --tail 100
docker logs bgc_db --tail 100

# Restart serviços
docker compose -f bgcstack/docker-compose.yml restart

# Limpar tudo
docker compose -f bgcstack/docker-compose.yml down -v
docker compose -f bgcstack/docker-compose.yml up -d

# Verificar métricas
curl http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])
```

### CI/CD

```bash
# Ver status
gh run list --workflow=e2e-main.yml

# Ver logs
gh run view <run-id> --log

# Baixar artifacts
gh run download <run-id>
```

---

## Documentação Complementar

### Criados Hoje
1. **`docs/STRATEGIC-ROADMAP-E2E-INTEGRATED.md`** (17k linhas)
   - Roadmap completo com testes E2E em cada fase
   - DoD detalhado por fase
   - Código de implementação incluído

2. **`docs/E2E-TESTING-STRATEGY.md`** (6k linhas)
   - Scripts executáveis prontos
   - Cheat sheets de comandos
   - Troubleshooting guide

3. **`ROADMAP-EXECUTIVO-E2E.md`** (Este arquivo)
   - Sumário executivo
   - Decision makers reference
   - Quick start guide

### Documentos Existentes Relevantes
- `docs/PRODUCT-ROADMAP.md` - Visão de produto 12 meses
- `docs/RELATORIO-NCM-COVERAGE-V0.4.0.md` - Estado atual de dados
- `STATUS_TESTES_E2E.md` - Status técnico atual
- `api/tests/e2e/simulator_test.go` - Testes existentes

---

## Governança e Aprovações

### Stakeholders
- **Product Manager:** Aprovar roadmap estratégico
- **Engineering Lead:** Aprovar DoD técnico
- **DevOps Lead:** Aprovar estratégia K8s
- **QA Lead:** Aprovar estratégia de testes

### Processo de Aprovação
1. Revisão deste documento (1h)
2. Q&A session com stakeholders (30min)
3. Aprovação formal (email/Slack)
4. Kickoff Phase 0 (hoje mesmo)

### Cadência de Revisão
- **Diária:** Standup com status de fase atual
- **Semanal:** Review de métricas E2E
- **Bi-semanal:** Sprint planning e priorização

---

## Conclusão

Este roadmap foi projetado com **E2E testing como cidadão de primeira classe**, garantindo que cada fase entregue valor real e validado de ponta a ponta.

**Princípios Fundamentais:**
1. Nenhuma fase marcada "done" sem E2E tests passando
2. Automação completa em CI/CD
3. Performance targets enforced, não opcionais
4. Rollback sempre possível (< 5 minutos)

**Próximo Passo Imediato:**
```bash
bash tests/e2e/phase0_cleanup.sh
```

Se passar, commitar e avançar para Phase 1 amanhã.

---

**Version:** 3.0
**Last Updated:** 2026-01-16
**Status:** APROVADO - EXECUTAR AGORA
**Owner:** BGC Product Management + Engineering Leadership
**Next Review:** 2026-01-20 (após Phase 1 completa)
