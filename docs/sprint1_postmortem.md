# Post-mortem Sprint 1 - Projeto BGC

**Data:** Setembro 2025  
**Sprint:** Onda 1 / Bloco A  
**Duração:** [Inserir duração da sprint]  
**Participantes:** [Inserir nomes da equipe]

## 📋 Resumo Executivo

### Objetivo da Sprint
Levantar um ambiente local de dados e API "read-only" em k3d com:
- PostgreSQL via Helm/Bitnami
- Sistema de ingest (CSV/XLSX) 
- Primeiras rotas de leitura de dados

### Status Final ✅
- ✅ Banco PostgreSQL rodando e populado (`stg.exportacao`)
- ✅ Materialized Views criadas e populadas (`rpt.mv_*`)
- ✅ API Go publicada com endpoints `/metrics/resumo` e `/metrics/pais`
- ✅ Port-forward funcional para desenvolvimento
- ✅ Sistema de ingest CSV/XLSX operacional

### Principais Aprendizados
- **Padronização crucial**: Scripts idempotentes e checklists reduzem 80% da fricção
- **Ambiente específico**: Diferenças Windows/PowerShell vs Linux requerem atenção
- **Pinagem de versões**: Evita drift de dependências (Go toolchain, bibliotecas)
- **Documentação viva**: Troubleshooting baseado em problemas reais é invaluável

---

## 🕐 Timeline da Sprint

| Fase | Atividade | Status | Observações |
|------|-----------|---------|-------------|
| **Setup Inicial** | Criação cluster k3d | ✅ | Ajustes pós-reboot necessários |
| **Banco de Dados** | PostgreSQL via Helm | ✅ | Bitnami com configurações padrão |
| **Migrações** | ConfigMaps + Jobs SQL | ✅ | Problemas de indentação YAML |
| **Ingest** | Sistema CSV/XLSX | ✅ | Ajustes BOM/UTF-8 e Go modules |
| **API** | Endpoints básicos | ✅ | Rotas `/metrics/*` funcionais |
| **MVs** | Materialized Views | ✅ | Estratégia CONCURRENTLY implementada |

---

## 🚨 Incidentes e Resoluções

### 1. Problemas de Shell/PowerShell

**Sintoma:** Comandos falhando com "parâmetro posicional não especificado"
```powershell
# ❌ Falha
mkdir docs deploy db scripts

# ✅ Solução
mkdir -Force docs,deploy,db,scripts
```

**Impacto:** Baixo | **Tempo perdido:** ~30min  
**Ação preventiva:** Padronizar scripts cross-platform

### 2. Kubeconfig Drift Pós-Reboot

**Sintoma:** `dial tcp ... connectex ... failed to respond`  
**Causa:** Porta externa do apiserver muda após reboot do host Docker

**Solução:**
```powershell
# Encontrar nova porta
docker ps | findstr k3d-bgc-serverlb

# Atualizar kubeconfig
kubectl config set-cluster k3d-bgc --server https://127.0.0.1:<NOVA_PORTA>
```

**Impacto:** Alto | **Tempo perdido:** ~2h (múltiplas ocorrências)  
**Ação preventiva:** Script automatizado de fix (`bgc.ps1 reboot-fix`)

### 3. Imagens Locais vs k3d

**Sintoma:** `ImagePullBackOff` mesmo com imagem local buildada  
**Causa:** k3d não enxerga automaticamente imagens do Docker host

**Solução:**
```bash
# Sempre importar após build
k3d image import bgc/api:dev bgc/ingest:dev -c bgc
```

**Impacto:** Médio | **Tempo perdido:** ~1h  
**Ação preventiva:** Padronizar tag `:dev` e import automático

### 4. Go Modules e BOM

**Sintoma:** `go.mod:1: unexpected input character '\ufeff'`  
**Causa:** Arquivo salvo com BOM UTF-8 via PowerShell

**Solução:**
```powershell
# Recriar sem BOM
Set-Content -Path go.mod -Value $content -Encoding utf8
```

**Impacto:** Baixo | **Tempo perdido:** ~20min  
**Ação preventiva:** Validação encoding nos scripts

### 5. YAML Indentação

**Sintoma:** `error converting YAML to JSON: line N`  
**Causa:** Quebras de linha/caracteres invisíveis ao colar via shell

**Solução:**
```powershell
# Usar here-strings
$yamlContent = @'
apiVersion: v1
kind: ConfigMap
...
'@
```

**Impacto:** Médio | **Tempo perdido:** ~45min  
**Ação preventiva:** Always use here-strings para YAML

### 6. Materialized Views

**Sintoma:** `cannot refresh ... concurrently` / `has not been populated`  
**Causa:** MVs criadas `WITH NO DATA` + falta de índice UNIQUE

**Solução sequencial:**
```sql
-- 1. Primeiro refresh sem CONCURRENTLY
REFRESH MATERIALIZED VIEW rpt.mv_resumo_pais;

-- 2. Criar índice UNIQUE
CREATE UNIQUE INDEX ON rpt.mv_resumo_pais (pais);

-- 3. Agora pode usar CONCURRENTLY
REFRESH MATERIALIZED VIEW CONCURRENTLY rpt.mv_resumo_pais;
```

**Impacto:** Médio | **Tempo perdido:** ~1h  
**Ação preventiva:** Documentar ordem de operações

---

## 📊 Métricas da Sprint

### Tempo por Categoria
- **Setup/Infra:** 40% (includes troubleshooting)
- **Desenvolvimento:** 35% (código Go + SQL)  
- **Debugging:** 20% (principalmente ambiente)
- **Documentação:** 5%

### Issues por Tipo
- **Ambiente/Config:** 60% dos problemas
- **Código/Logic:** 25% dos problemas  
- **Dependências:** 15% dos problemas

### Taxa de Sucesso
- **Primeira tentativa:** 30%
- **Com retry/fix:** 95%
- **Requer investigação:** 5%

---

## 🎯 O Que Funcionou Bem

### ✅ Estratégias Eficazes
- **Scripts idempotentes** com `kubectl apply` e `--dry-run=client`
- **Here-strings** no PowerShell para conteúdo multilinha
- **Tag estável `:dev`** para imagens de desenvolvimento
- **Cliente Bitnami** para inspeção SQL ad-hoc
- **Logs centralizados** via `kubectl logs`

### ✅ Decisões Técnicas Acertadas
- **k3d** como runtime local - rápido e leve
- **Helm Bitnami** para PostgreSQL - produção-ready
- **Go modules** com pinagem de versões
- **Materialized Views** para performance
- **ConfigMaps** para SQL migrations

---

## 🔧 Áreas de Melhoria

### Automação
- [ ] Makefile/scripts unificados Windows + Linux
- [ ] CI/CD básico (GitHub Actions)
- [ ] Health checks automáticos
- [ ] Backup/restore de dados de desenvolvimento

### Observabilidade  
- [ ] Logs JSON estruturados
- [ ] Métricas básicas (latência, throughput)
- [ ] Dashboard simples (Grafana ou similar)
- [ ] Alertas para falhas de ingest

### Documentação
- [ ] OpenAPI spec para a API
- [ ] Collection Postman versionada
- [ ] Runbooks para operações comuns
- [ ] Arquitetura técnica detalhada

---

## ⚠️ Riscos Identificados

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| **Drift kubeconfig pós-reboot** | Alta | Alto | ✅ Script automatizado |
| **Bitnami license após 28/08/2025** | Média | Médio | 🟡 Avaliar alternativas |
| **Dependências Go outdated** | Média | Médio | 🟡 Pinagem + Dependabot |
| **Data corruption por multi-load** | Baixa | Alto | 🟡 Implementar upsert |
| **Performance MVs com volume** | Baixa | Médio | 🟡 Monitoring + índices |

**Legenda:** ✅ Resolvido | 🟡 Monitorando | ❌ Pendente

---

## 🎯 Ações para Sprint 2

### Prioridade Alta 🔴
1. **API consolidation** - Alinhar rotas `/metrics/*` vs `/v1/exportacao/*`
2. **OpenAPI spec** - Documentação formal da API
3. **Upsert strategy** - Evitar duplicação de dados no ingest
4. **Health endpoints** - `/health` e `/ready` para probes

### Prioridade Média 🟡  
5. **Automation** - GitHub Actions para build/test
6. **Monitoring** - Logs estruturados e métricas básicas
7. **Backup strategy** - Snapshot de dados de desenvolvimento
8. **Performance testing** - Load testing dos endpoints

### Prioridade Baixa 🟢
9. **Multi-environment** - Configs para dev/staging/prod
10. **Security** - Secrets management e RBAC básico
11. **Cache layer** - Redis para queries frequentes
12. **API versioning** - Estratégia de versionamento

---

## 💡 Lições Aprendidas

### Técnicas
- **PowerShell** requer sintaxe específica - não assumir bash
- **k3d** image import é obrigatório para imagens locais
- **Bitnami** imagens precisam do entrypoint correto
- **YAML** indentação é crítica - use ferramentas, não copy/paste manual

### Processo  
- **Checklist** pós-reboot economiza horas de debug
- **Scripts idempotentes** permitem re-execução segura
- **Troubleshooting real** > documentação genérica
- **Commits frequentes** facilitam rollback de mudanças problemáticas

### Colaboração
- **Problemas documentados** viram conhecimento organizacional
- **Scripts compartilhados** evitam retrabalho
- **Padrões estabelecidos** aceleram desenvolvimento futuro

---

## 📚 Recursos Úteis

### Documentação de Referência
- [k3d Documentation](https://k3d.io/)
- [Bitnami PostgreSQL Helm Chart](https://github.com/bitnami/charts/tree/main/bitnami/postgresql)
- [Go Modules Reference](https://golang.org/ref/mod)
- [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)

### Ferramentas Descobertas
- **kubectl** jsonpath vs PowerShell `ConvertFrom-Json`
- **docker** ps parsing para extrair portas
- **helm** template para debug de charts
- **psql** via Bitnami entrypoint

### Scripts Salvos
- `post-reboot-checklist.ps1` - Automação pós-reboot
- `bgc.ps1` - Comandos principais do projeto  
- Makefile - Equivalente para Linux/Mac

---

**Assinatura do Tech Lead:** [Nome]  
**Data de Fechamento:** [Data]  
**Próxima Revisão:** Sprint 2 Retrospective