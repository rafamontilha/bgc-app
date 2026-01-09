# scripts/migrate-to-product.ps1
#
# Script de migracao de documentos de planejamento para estrutura .product/
#
# Uso:
#   .\scripts\migrate-to-product.ps1          # Executar migracao
#   .\scripts\migrate-to-product.ps1 -DryRun  # Visualizar mudancas sem aplicar

param(
    [switch]$DryRun = $false,
    [switch]$Force = $false
)

$ErrorActionPreference = "Stop"

Write-Host @"

================================================================================
  MIGRACAO PARA .product/ - BGC Product Management
================================================================================

Este script vai:
1. Criar estrutura de diretorios em .product/
2. Migrar documentos existentes (COPY ou MOVE)
3. Criar READMEs em cada diretorio
4. Criar templates padronizados

"@ -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "MODO DRY-RUN: Nenhuma mudanca sera aplicada.`n" -ForegroundColor Yellow
} else {
    Write-Host "ATENCAO: Mudancas serao aplicadas ao filesystem!`n" -ForegroundColor Red
    if (-not $Force) {
        $confirm = Read-Host "Continuar? (s/N)"
        if ($confirm -ne 's' -and $confirm -ne 'S') {
            Write-Host "Operacao cancelada." -ForegroundColor Yellow
            exit 0
        }
    } else {
        Write-Host "MODO FORCE: Confirmacao automatica ativada.`n" -ForegroundColor Yellow
    }
}

Write-Host "`n[1/4] Criando estrutura de diretorios..." -ForegroundColor Green

# 1. Criar estrutura de diretorios
$directories = @(
    ".product",
    ".product/north-star",
    ".product/okrs",
    ".product/roadmap",
    ".product/epics",
    ".product/features",
    ".product/sprints",
    ".product/sprints/_template",
    ".product/sprints/2025-W01",
    ".product/sprints/2025-W02",
    ".product/backlog",
    ".product/discovery",
    ".product/discovery/user-research",
    ".product/discovery/experiments",
    ".product/discovery/personas",
    ".product/metrics",
    ".product/retrospectives",
    ".product/decisions",
    ".product/prs",
    ".product/archive",
    ".product/archive/2025"
)

$dirCount = 0
foreach ($dir in $directories) {
    if (-not (Test-Path $dir)) {
        if ($DryRun) {
            Write-Host "  [DRY-RUN] Criaria diretorio: $dir" -ForegroundColor Yellow
        } else {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
            $dirCount++
        }
    }
}

if (-not $DryRun) {
    Write-Host "  Criados $dirCount diretorios" -ForegroundColor Cyan
}

Write-Host "`n[2/4] Migrando arquivos existentes..." -ForegroundColor Green

# 2. Migrar arquivos
$migrations = @(
    @{
        Source = "web-public/ROADMAP.md"
        Dest = ".product/roadmap/web-public.md"
        Action = "MOVE"
        Description = "Roadmap do web-public (desatualizado, precisa revisao)"
    },
    @{
        Source = "docs/EPIC-1-COMPLETE.md"
        Dest = ".product/epics/epic-001-integration-gateway.md"
        Action = "COPY"
        Description = "Epico 1 - Integration Gateway (COMPLETO)"
    },
    @{
        Source = "docs/EPIC-2-COMPLETE.md"
        Dest = ".product/epics/epic-002-observability.md"
        Action = "COPY"
        Description = "Epico 2 - Observability (COMPLETO)"
    },
    @{
        Source = "docs/RELATORIO-EPICO-3-MELHORIAS.md"
        Dest = ".product/epics/epic-003-api-contracts.md"
        Action = "MOVE"
        Description = "Epico 3 - API Contracts (COMPLETO)"
    },
    @{
        Source = "docs/Sprint2_E2E_Checklist.md"
        Dest = ".product/sprints/2025-W02/planning.md"
        Action = "MOVE"
        Description = "Sprint 2 - Checklist E2E"
    },
    @{
        Source = "docs/sprint1_postmortem.md"
        Dest = ".product/sprints/2025-W01/retrospective.md"
        Action = "MOVE"
        Description = "Sprint 1 - Retrospectiva"
    },
    @{
        Source = "docs/EPIC-1-PROGRESS.md"
        Dest = ".product/archive/2025/epic-1-progress.md"
        Action = "MOVE"
        Description = "Epico 1 - Progress (ARQUIVADO)"
    },
    @{
        Source = "docs/EPIC-1-FINAL.md"
        Dest = ".product/archive/2025/epic-1-final.md"
        Action = "MOVE"
        Description = "Epico 1 - Final (ARQUIVADO)"
    },
    @{
        Source = "docs/EPIC-1-SUMMARY.md"
        Dest = ".product/archive/2025/epic-1-summary.md"
        Action = "MOVE"
        Description = "Epico 1 - Summary (ARQUIVADO)"
    },
    @{
        Source = "docs/NEXT-STEPS.md"
        Dest = ".product/archive/2025/next-steps-epic-1.md"
        Action = "MOVE"
        Description = "Next Steps Epico 1 (ARQUIVADO - integrar em roadmap master)"
    },
    @{
        Source = ".pr-body.md"
        Dest = ".product/prs/feature-security-credentials.md"
        Action = "MOVE"
        Description = "PR Draft - Security Credentials Management"
    },
    @{
        Source = ".pr-description.md"
        Dest = ".product/archive/2025/pr-description-temp.md"
        Action = "MOVE"
        Description = "PR Description temporaria (ARQUIVADO)"
    }
)

$migratedCount = 0
$skippedCount = 0

foreach ($migration in $migrations) {
    if (Test-Path $migration.Source) {
        $destDir = Split-Path $migration.Dest -Parent
        if (-not (Test-Path $destDir)) {
            if (-not $DryRun) {
                New-Item -ItemType Directory -Path $destDir -Force | Out-Null
            }
        }

        if ($DryRun) {
            Write-Host "  [DRY-RUN] $($migration.Action): $($migration.Source)" -ForegroundColor Yellow
            Write-Host "            -> $($migration.Dest)" -ForegroundColor Yellow
            Write-Host "            ($($migration.Description))" -ForegroundColor DarkGray
        } else {
            if ($migration.Action -eq "COPY") {
                Copy-Item $migration.Source $migration.Dest -Force
                Write-Host "  COPIADO: $($migration.Source)" -ForegroundColor Green
                Write-Host "        -> $($migration.Dest)" -ForegroundColor Cyan
            } elseif ($migration.Action -eq "MOVE") {
                Move-Item $migration.Source $migration.Dest -Force
                Write-Host "  MOVIDO: $($migration.Source)" -ForegroundColor Magenta
                Write-Host "       -> $($migration.Dest)" -ForegroundColor Cyan
            }
            $migratedCount++
        }
    } else {
        Write-Host "  AVISO: Arquivo nao encontrado: $($migration.Source)" -ForegroundColor DarkYellow
        $skippedCount++
    }
}

if (-not $DryRun) {
    Write-Host "`n  Migrados: $migratedCount arquivos | Nao encontrados: $skippedCount" -ForegroundColor Cyan
}

Write-Host "`n[3/4] Criando READMEs..." -ForegroundColor Green

# 3. Criar READMEs em cada diretorio
$readmes = @{
    ".product/README.md" = @"
# .product/ - Planejamento de Produto BGC

Esta pasta contem todos os documentos de planejamento de produto da plataforma BGC.

## IMPORTANTE

Esta pasta e **GITIGNORED**. Conteudo estrategico e de planejamento nao deve ser commitado no repositorio publico.

## Estrutura

- **north-star/** - Definicao e tracking da North Star Metric
- **okrs/** - Objectives & Key Results trimestrais
- **roadmap/** - Roadmap master e roadmaps por componente
- **epics/** - Grandes iniciativas de produto (1-3 meses)
- **features/** - Features de tamanho medio (1-4 semanas)
- **sprints/** - Planning e retrospectivas semanais
- **backlog/** - Backlog priorizado com RICE scoring
- **discovery/** - User research, experimentos, personas
- **metrics/** - Metricas de negocio e dashboards
- **retrospectives/** - Retrospectivas trimestrais e anuais
- **decisions/** - Architecture Decision Records (ADRs)
- **prs/** - Drafts de Pull Requests
- **archive/** - Arquivos antigos

## Como Usar

Ver documentacao completa em: .product-reorganization-proposal.md (raiz do projeto)

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/north-star/README.md" = @"
# North Star Metric

A North Star Metric e a metrica principal que guia o produto.

## Definicao

Ver: definition.md

## Tracking

Ver: tracking.md (atualizado semanalmente)

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/okrs/README.md" = @"
# OKRs - Objectives & Key Results

OKRs trimestrais do produto BGC.

## Estrutura

- 2026-Q1.md
- 2026-Q2.md
- 2026-Q3.md
- 2026-Q4.md

Cada arquivo contem:
- 2-3 Objectives
- 3-4 Key Results por Objective
- Status de progresso

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/roadmap/README.md" = @"
# Roadmap

Roadmap master e roadmaps por componente.

## Arquivos

- **master.md** - Roadmap 6-12 meses (visao unificada)
- **api.md** - Roadmap do backend API
- **web-public.md** - Roadmap do frontend publico
- **integration-gateway.md** - Roadmap do gateway de integracoes
- **ingest.md** - Roadmap do servico de ingestao

## Cadencia de Revisao

- Semanal: Atualizar status de features em andamento
- Mensal: Repriorizar backlog
- Trimestral: Revisar roadmap de 6 meses

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/epics/README.md" = @"
# Epicos

Grandes iniciativas de produto (duracao: 1-3 meses).

## Template

Use: _template.md

## Epicos Atuais

- epic-001-integration-gateway.md (COMPLETO)
- epic-002-observability.md (COMPLETO)
- epic-003-api-contracts.md (COMPLETO)
- epic-004-web-simulator.md (EM ANDAMENTO)

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/features/README.md" = @"
# Features

Features de tamanho medio (duracao: 1-4 semanas).

## Template

Use: _template.md

## Priorizacao

Features sao priorizadas usando RICE scoring:
- Reach (alcance)
- Impact (impacto)
- Confidence (confianca)
- Effort (esforco)

Score = (Reach x Impact x Confidence) / Effort

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/sprints/README.md" = @"
# Sprints

Planejamento e retrospectivas de sprints semanais.

## Estrutura

Cada sprint tem um diretorio: YYYY-WXX/ (ISO week number)

Exemplo: 2025-W47/

Dentro de cada diretorio:
- planning.md - Sprint planning
- retrospective.md - Retrospectiva
- daily-notes.md (opcional) - Notas diarias

## Templates

Use: _template/planning-template.md e _template/retrospective-template.md

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/backlog/README.md" = @"
# Backlog

Backlog unificado e priorizado da plataforma BGC.

## Arquivos

- **backlog.md** - Backlog principal (RICE scoring)
- **icebox.md** - Ideias futuras (baixa prioridade)
- **archive.md** - Items rejeitados ou depreciados

## Priorizacao

Usar RICE:
- Reach (alcance)
- Impact (impacto)
- Confidence (confianca %)
- Effort (esforco em pontos)

Score = (Reach x Impact x Confidence) / Effort

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/discovery/README.md" = @"
# Product Discovery

User research, experimentos, personas.

## Estrutura

- **user-research/** - Entrevistas, surveys, usability tests
- **experiments/** - Experimentos e hipoteses
- **personas/** - Personas detalhadas

## Templates

Ver subdiretorios para templates especificos.

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/metrics/README.md" = @"
# Metricas de Negocio

Dashboards e reports de metricas de produto.

## Arquivos

- **dashboards.md** - Links e screenshots de dashboards
- **weekly-report.md** - Relatorio semanal de metricas
- **alerts.md** - Alertas configurados

## Metricas Principais

- North Star Metric (ver ../north-star/)
- KPIs de OKRs (ver ../okrs/)
- Funnel de conversao
- Retention & engagement

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/retrospectives/README.md" = @"
# Retrospectivas

Retrospectivas trimestrais e anuais.

Para retrospectivas semanais de sprint, ver: ../sprints/

## Template

Use: _template.md

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/decisions/README.md" = @"
# Architecture Decision Records (ADRs)

Decisoes arquiteturais e de produto importantes.

## Template

Use: _template.md

## Formato

Cada ADR tem:
- Numero sequencial (001, 002, ...)
- Status (PROPOSTA, ACEITA, DEPRECIADA, SUPERSEDED)
- Contexto e problema
- Decisao tomada
- Alternativas consideradas
- Consequencias

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/prs/README.md" = @"
# Pull Request Drafts

Drafts de PRs antes de criar no GitHub.

Util para:
- Preparar descricao detalhada
- Listar checklist de teste
- Referenciar epicos e features
- Revisar antes de abrir PR

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/archive/README.md" = @"
# Archive

Arquivos antigos de planejamento.

## Estrutura

Organizado por ano: 2025/, 2026/, etc.

## Quando Arquivar

- Sprints antigas (apos 3 meses)
- Features completas (apos 1 mes)
- Documentos obsoletos

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@
}

$readmeCount = 0
foreach ($readme in $readmes.GetEnumerator()) {
    if ($DryRun) {
        Write-Host "  [DRY-RUN] Criaria: $($readme.Key)" -ForegroundColor Yellow
    } else {
        if (-not (Test-Path $readme.Key)) {
            Set-Content -Path $readme.Key -Value $readme.Value -Encoding UTF8
            $readmeCount++
        }
    }
}

if (-not $DryRun) {
    Write-Host "  Criados $readmeCount READMEs" -ForegroundColor Cyan
}

Write-Host "`n[4/4] Criando templates..." -ForegroundColor Green

# 4. Criar templates basicos (versao simplificada - versao completa esta na proposta)
$templates = @{
    ".product/epics/_template.md" = @"
# Epico XXX: [Nome do Epico]

**Status:** PLANEJADO | EM ANDAMENTO | COMPLETO | CANCELADO
**Owner:** [Nome]
**Prazo:** YYYY-MM-DD a YYYY-MM-DD
**OKR Relacionado:** [Link]
**Componentes:** API | Web | Gateway | Ingest

## 1. PROBLEMA E CONTEXTO

[Descricao do problema]

## 2. SOLUCAO PROPOSTA

[Descricao da solucao]

## 3. METRICAS DE SUCESSO

- [ ] Criterio 1
- [ ] Criterio 2

## 4. FEATURES ASSOCIADAS

- [ ] Feature 001
- [ ] Feature 002

## 5. RISCOS E DEPENDENCIAS

[Listar riscos e dependencias]

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/sprints/_template/planning-template.md" = @"
# Sprint Planning - Semana YYYY-WXX

**Periodo:** YYYY-MM-DD a YYYY-MM-DD
**Objetivo da Sprint:** [Objetivo em 1 frase]

## BACKLOG DA SPRINT

| ID | Descricao | Owner | Estimativa | Prioridade |
|----|-----------|-------|------------|------------|
| F-001 | [Desc] | [Nome] | 5 | ALTA |

## DEFINITION OF DONE

- [ ] Criterio 1
- [ ] Criterio 2

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@

    ".product/sprints/_template/retrospective-template.md" = @"
# Sprint Retrospective - Semana YYYY-WXX

**Periodo:** YYYY-MM-DD a YYYY-MM-DD

## METRICAS

- Pontos planejados: XX
- Pontos concluidos: XX
- Velocidade: XX%

## O QUE FUNCIONOU BEM?

- Item 1

## O QUE PODE MELHORAR?

- Item 1

## ACOES PARA PROXIMA SPRINT

- [ ] Acao 1 (Owner: [Nome])

**Criado em:** $(Get-Date -Format "yyyy-MM-dd")
"@
}

$templateCount = 0
foreach ($template in $templates.GetEnumerator()) {
    if ($DryRun) {
        Write-Host "  [DRY-RUN] Criaria: $($template.Key)" -ForegroundColor Yellow
    } else {
        $templateDir = Split-Path $template.Key -Parent
        if (-not (Test-Path $templateDir)) {
            New-Item -ItemType Directory -Path $templateDir -Force | Out-Null
        }
        if (-not (Test-Path $template.Key)) {
            Set-Content -Path $template.Key -Value $template.Value -Encoding UTF8
            $templateCount++
        }
    }
}

if (-not $DryRun) {
    Write-Host "  Criados $templateCount templates" -ForegroundColor Cyan
}

Write-Host "`n" -NoNewline
Write-Host "================================================================================`n" -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "DRY-RUN COMPLETO!" -ForegroundColor Yellow
    Write-Host "`nNenhuma mudanca foi aplicada. Execute sem -DryRun para aplicar.`n" -ForegroundColor Yellow
} else {
    Write-Host "MIGRACAO CONCLUIDA COM SUCESSO!" -ForegroundColor Green
    Write-Host @"

Proximos passos:

1. Verificar estrutura criada:
   ls -R .product/

2. Revisar arquivos migrados:
   cd .product/epics
   ls

3. Criar North Star Metric:
   Editar: .product/north-star/definition.md

4. Criar Roadmap Master:
   Copiar template: .product-roadmap-master-example.md
   Salvar em: .product/roadmap/master.md

5. Planning da proxima sprint:
   Criar: .product/sprints/2025-W48/planning.md

Ver documentacao completa em: .product-reorganization-proposal.md

"@ -ForegroundColor Cyan
}

Write-Host "================================================================================`n" -ForegroundColor Cyan
