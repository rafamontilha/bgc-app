# Diagnóstico BGC App - 12 Fatores + Segurança

**Pontuação:** 🟡 **75/100** - Bom, com **gaps críticos de segurança**

---

## 📊 Resumo por Fator

| # | Fator | Status | Nota | Problema Principal |
|---|-------|--------|------|-------------------|
| I | Base de Código | 🟢 | 10/10 | - |
| II | Dependências | 🟢 | 10/10 | - |
| III | **Configurações** | 🔴 | **4/10** | **Configs hardcoded na imagem** |
| IV | Serviços de Apoio | 🟢 | 10/10 | - |
| V | Build/Release/Execute | 🟡 | 7/10 | Sem CI/CD, tags mutáveis |
| VI | Processos Stateless | 🟢 | 10/10 | - |
| VII | Vínculo de Porta | 🟢 | 10/10 | - |
| VIII | Concorrência | 🟢 | 9/10 | - |
| IX | Descartabilidade | 🟢 | 9/10 | - |
| X | Paridade Dev/Prod | 🟡 | 8/10 | Dev usa volumes, prod não |
| XI | **Logs** | 🟡 | **6/10** | **Formato misto, sem agregação** |
| XII | Processos Admin | 🟢 | 9/10 | - |
| **🔒** | **SEGURANÇA** | 🔴 | **3/10** | **CRÍTICO - Credenciais expostas** |

---

## 🚨 CRÍTICO - Problemas de Segurança

### 1. **Credenciais Hardcoded em Plaintext** 🔴

**Locais identificados:**

```yaml
# bgcstack/docker-compose.yml
environment:
  POSTGRES_PASSWORD: bgc        # ❌ Senha em plaintext
  PGADMIN_DEFAULT_PASSWORD: admin  # ❌ Senha em plaintext
  DB_PASS: bgc                  # ❌ Senha em plaintext
```

```yaml
# README.md (linhas 283-294)
### PostgreSQL
- Password: bgc               # ❌ Documentação pública com senha

### PgAdmin
- Password: admin             # ❌ Documentação pública com senha
```

**Risco:** 🔴 **ALTO**
- Qualquer pessoa com acesso ao repo tem credenciais do banco
- Se commitado no GitHub público = **vazamento total**
- Credential stuffing se usar mesma senha em prod

---

### 2. **Secrets em Kubernetes Não Gerenciados** 🔴

**Evidência:**

```yaml
# k8s/api.yaml linha 34-38
env:
  - name: DB_PASS
    valueFrom:
      secretKeyRef:
        name: pg-postgresql        # ✅ Usa secret
        key: postgres-password     # ⚠️ Mas secret é criado pelo Helm com senha padrão
```

**Problema:**
- Secret `pg-postgresql` criado automaticamente pelo Helm Bitnami
- Senha provavelmente é a padrão ou configurada via values.yaml
- Não há rotação de secrets
- Não há criptografia em repouso (cluster não configurado)

---

### 3. **Configurações Sensíveis no Git** 🟡

```dockerfile
# api/Dockerfile linha 25
COPY config /app/config  # ⚠️ Pode conter dados sensíveis
```

**Arquivos de config na imagem:**
- `config/partners_stub.yaml` - OK (dados públicos)
- `config/tariff_scenarios.yaml` - OK (dados públicos)

**Risco Futuro:** Se adicionar API keys ou tokens nestes arquivos, serão expostos.

---

### 4. **Sem Scan de Vulnerabilidades** 🟡

**Ausente:**
- Scan de imagens Docker (Trivy, Snyk)
- Scan de dependências (Dependabot)
- SAST (Static Analysis Security Testing)

**Risco:** Vulnerabilidades conhecidas não detectadas (ex: Log4Shell, SQLi)

---

## 🔧 Recomendações de Segurança (Priorizadas)

### 🔴 **Crítico - Implementar AGORA**

#### 1. **Remover credenciais do Git** (1 dia)

```bash
# .gitignore
.env
.env.*
!.env.example
*.secret
**/secrets/
docker-compose.override.yml
```

```bash
# Criar .env.example (template)
DB_USER=bgc
DB_PASSWORD=<CHANGE_ME>
PGADMIN_EMAIL=admin@bgc.dev
PGADMIN_PASSWORD=<CHANGE_ME>
```

```bash
# Limpar histórico do Git (se já commitado)
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch docker-compose.yml" \
  --prune-empty --tag-name-filter cat -- --all
```

#### 2. **Gestão de Secrets com Sealed Secrets** (3 dias)

```bash
# Instalar Sealed Secrets Controller
kubectl apply -f \
  https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Criar secret original
kubectl create secret generic bgc-db-credentials \
  --from-literal=username=bgc \
  --from-literal=password=$(openssl rand -base64 32) \
  --dry-run=client -o yaml > secret.yaml

# Selar (criptografar com chave pública do cluster)
kubeseal -f secret.yaml -w sealed-secret.yaml

# Commitar sealed-secret.yaml (é seguro!)
git add k8s/sealed-secret.yaml
git commit -m "chore: add sealed secrets for DB credentials"
```

```yaml
# k8s/api.yaml - usar sealed secret
env:
  - name: DB_PASS
    valueFrom:
      secretKeyRef:
        name: bgc-db-credentials  # Nome do SealedSecret
        key: password
```

#### 3. **Senhas Fortes e Rotação** (1 dia)

```bash
# Gerar senhas criptograficamente seguras
openssl rand -base64 32

# Criar secrets diferentes por ambiente
kubectl create secret generic bgc-db-prod \
  --from-literal=password=$(openssl rand -base64 32) -n production

kubectl create secret generic bgc-db-staging \
  --from-literal=password=$(openssl rand -base64 32) -n staging
```

**Política de rotação:**
- Trocar senhas a cada 90 dias
- Rotação automática via External Secrets Operator (futuro)

---

### 🟡 **Importante - 1-2 meses**

#### 4. **Scan de Vulnerabilidades** (1 semana)

```yaml
# .github/workflows/security.yml
name: Security Scan

on: [push, pull_request]

jobs:
  trivy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'bgc/bgc-api:dev'
          format: 'sarif'
          output: 'trivy-results.sarif'
          severity: 'CRITICAL,HIGH'

      - name: Upload to GitHub Security
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'

  dependency-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Go vulnerability check
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          cd api && govulncheck ./...
```

#### 5. **Network Policies** (3 dias)

```yaml
# k8s/network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: bgc-api-policy
  namespace: data
spec:
  podSelector:
    matchLabels:
      app: bgc-api
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
      - podSelector:
          matchLabels:
            app: bgc-web
      - namespaceSelector:
          matchLabels:
            name: ingress-nginx
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
      - podSelector:
          matchLabels:
            app: postgres
      ports:
        - protocol: TCP
          port: 5432
```

#### 6. **RBAC Kubernetes** (2 dias)

```yaml
# k8s/rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bgc-api-sa
  namespace: data
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: bgc-api-role
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    resourceNames: ["bgc-db-credentials"]
    verbs: ["get"]
---
# Deployment usa serviceAccount
spec:
  template:
    spec:
      serviceAccountName: bgc-api-sa
```

---

### 🟢 **Desejável - 3-6 meses**

#### 7. **Vault para Secrets** (2 semanas)

```bash
# HashiCorp Vault + External Secrets Operator
helm install vault hashicorp/vault
helm install external-secrets external-secrets/external-secrets
```

```yaml
# ExternalSecret sincroniza do Vault
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: bgc-db-secret
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
  target:
    name: bgc-db-credentials
  data:
    - secretKey: password
      remoteRef:
        key: secret/data/bgc/database
        property: password
```

#### 8. **Auditoria e Compliance** (contínuo)

- Habilitar Kubernetes Audit Logs
- Implementar Falco para runtime security
- OWASP Dependency Check mensal
- Penetration testing trimestral

---

## 📋 Análise dos 12 Fatores (Resumida)

### ✅ **Bem Implementados (8/12)**

| Fator | Evidência |
|-------|-----------|
| I. Base de Código | Git com commits semânticos, monorepo |
| II. Dependências | `go.mod`, `pnpm-lock.yaml`, multi-stage builds |
| IV. Serviços de Apoio | PostgreSQL como recurso anexado via env vars |
| VI. Processos | API stateless, sem sessões em memória |
| VII. Vínculo de Porta | API:8080, Web:3000 self-contained |
| VIII. Concorrência | HPA 1-5 pods, escala horizontal |
| IX. Descartabilidade | Startup <10s, health probes |
| XII. Admin | Jobs K8s para migrations, CronJobs para backup |

---

### ⚠️ **Requerem Melhoria (4/12)**

#### **III. Configurações** - 🔴 4/10

**Problemas:**
- ❌ `COPY config /app/config` - configs bundled na imagem
- ❌ Senhas em plaintext no docker-compose
- ❌ Sem suporte a múltiplos ambientes

**Solução:**
```yaml
# Mover para ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: bgc-app-config
data:
  SCOPE_CHAPTERS: "02,08,84,85"
  SOM_BASE: "0.015"
  partners.yaml: |
    USA: 0.35
    CHN: 0.25
```

---

#### **V. Build/Release/Execute** - 🟡 7/10

**Problemas:**
- ❌ Sem CI/CD automatizado
- ❌ Tags mutáveis (`bgc-api:dev`)
- ✅ Multi-stage builds (bom)

**Solução:**
```yaml
# GitHub Actions
- name: Build and tag
  run: |
    docker build -t bgc-api:${{ github.sha }}
    docker build -t bgc-api:v${{ github.ref_name }}
```

---

#### **X. Paridade Dev/Prod** - 🟡 8/10

**Problemas:**
- ⚠️ Dev: Docker Compose com volumes
- ⚠️ Prod: K8s sem volumes
- ❌ Falta ambiente staging

**Solução:**
- Usar k3d também para dev (já existe!)
- Criar namespace `staging` no cluster
- Pipeline: dev → staging → prod

---

#### **XI. Logs** - 🟡 6/10

**Problemas:**
```
❌ Formato misto:
  {"ts":"...","level":"info",...}  # JSON ✅
  [GIN] 2025/10/24 - 17:04:04...   # Texto ❌

❌ Sem agregação (Loki/ELK)
❌ Trace ID incompleto
```

**Solução:**
```go
// Padronizar tudo em JSON
r.Use(gin.LoggerWithFormatter(jsonFormatter))
```

---

## 📦 Plano de Ação (6 Semanas)

### Semana 1-2: 🔴 **SEGURANÇA CRÍTICA**
- [ ] Remover credenciais do Git
- [ ] Implementar Sealed Secrets
- [ ] Gerar senhas fortes
- [ ] Adicionar `.gitignore` adequado
- [ ] Scan com Trivy

### Semana 3-4: 🟡 **CI/CD + Configs**
- [ ] GitHub Actions pipeline
- [ ] Externalizar configs para ConfigMaps
- [ ] Versionamento de imagens (SHA)
- [ ] Network Policies

### Semana 5-6: 🟢 **Observabilidade**
- [ ] Padronizar logs JSON
- [ ] Deploy Loki + Promtail
- [ ] Dashboards Grafana
- [ ] Criar ambiente staging

---

## 🎯 Próximos Passos Imediatos

1. **HOJE:** Remover credenciais do docker-compose.yml
2. **ESTA SEMANA:** Implementar Sealed Secrets
3. **ESTE MÊS:** CI/CD + Scan de vulnerabilidades

**Risco Atual:** 🔴 **ALTO** - Credenciais expostas
**Após implementação:** 🟢 **BAIXO** - Aplicação cloud-native segura

---

**Conclusão:** Aplicação tecnicamente sólida (9/12 fatores), mas com **gaps críticos de segurança** que devem ser endereçados imediatamente antes de produção.
