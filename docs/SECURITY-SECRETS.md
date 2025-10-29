# Guia de Segurança e Gestão de Secrets - BGC App

## 🔒 Visão Geral

Este documento descreve como gerenciar credenciais e secrets de forma segura no projeto BGC App.

---

## 📋 Princípios de Segurança

1. **NUNCA** committar credenciais no Git
2. **SEMPRE** usar senhas fortes geradas criptograficamente
3. **RODAR** diferentes senhas por ambiente (dev/staging/prod)
4. **ROTACIONAR** senhas a cada 90 dias
5. **USAR** Sealed Secrets para Kubernetes

---

## 🐳 Docker Compose (Desenvolvimento Local)

### Setup Inicial

```bash
# 1. Copiar template
cd bgcstack
cp .env.example .env

# 2. Gerar senhas fortes
openssl rand -base64 32  # PostgreSQL
openssl rand -base64 32  # PgAdmin

# 3. Editar .env e preencher as senhas
nano .env
```

### Arquivo .env

```bash
POSTGRES_PASSWORD=<SENHA_GERADA_1>
PGADMIN_DEFAULT_PASSWORD=<SENHA_GERADA_2>
DB_PASS=<MESMA_SENHA_DO_POSTGRES>
```

### Iniciar Stack

```bash
docker-compose up -d
```

**⚠️ IMPORTANTE:** O arquivo `.env` está no `.gitignore` e NUNCA deve ser commitado!

---

## ☸️ Kubernetes (Produção)

### Opção 1: Sealed Secrets (Recomendado)

#### 1. Instalar Sealed Secrets Controller

```bash
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml
```

#### 2. Instalar kubeseal CLI

**Windows (via Chocolatey):**
```powershell
choco install kubeseal
```

**Linux:**
```bash
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/kubeseal-0.24.0-linux-amd64.tar.gz
tar xfz kubeseal-0.24.0-linux-amd64.tar.gz
sudo install -m 755 kubeseal /usr/local/bin/kubeseal
```

#### 3. Criar e Selar Secret

```bash
# Gerar senha forte
PASSWORD=$(openssl rand -base64 32)

# Criar secret
kubectl create secret generic bgc-db-credentials \
  --from-literal=username=bgc \
  --from-literal=password=$PASSWORD \
  --namespace=data \
  --dry-run=client -o yaml > /tmp/secret.yaml

# Selar (criptografar)
kubeseal -f /tmp/secret.yaml -w k8s/bgc-db-sealed-secret.yaml

# Deletar arquivo temporário
rm /tmp/secret.yaml

# Commitar sealed secret (é seguro!)
git add k8s/bgc-db-sealed-secret.yaml
git commit -m "chore: add sealed secret for DB credentials"
```

#### 4. Aplicar Sealed Secret

```bash
kubectl apply -f k8s/bgc-db-sealed-secret.yaml
```

O controller irá descriptografar automaticamente e criar o Secret `bgc-db-credentials`.

#### 5. Atualizar Deployment

```yaml
# k8s/api.yaml
env:
  - name: DB_USER
    valueFrom:
      secretKeyRef:
        name: bgc-db-credentials
        key: username
  - name: DB_PASS
    valueFrom:
      secretKeyRef:
        name: bgc-db-credentials
        key: password
```

---

### Opção 2: Kubernetes Secrets Nativos (Sem Sealed Secrets)

**⚠️ NÃO RECOMENDADO para produção** (secrets em base64 não são criptografia real)

```bash
# Criar secret diretamente
kubectl create secret generic bgc-db-credentials \
  --from-literal=username=bgc \
  --from-literal=password=$(openssl rand -base64 32) \
  --namespace=data

# Verificar
kubectl get secrets -n data
kubectl describe secret bgc-db-credentials -n data
```

---

## 🔄 Rotação de Secrets

### Docker Compose

```bash
# 1. Gerar nova senha
NEW_PASSWORD=$(openssl rand -base64 32)

# 2. Atualizar .env
echo "POSTGRES_PASSWORD=$NEW_PASSWORD" >> .env
echo "DB_PASS=$NEW_PASSWORD" >> .env

# 3. Recriar containers
docker-compose down
docker-compose up -d

# 4. Conectar ao PostgreSQL e alterar senha
docker exec -it bgc_db psql -U bgc -c "ALTER USER bgc WITH PASSWORD '$NEW_PASSWORD';"
```

### Kubernetes

```bash
# 1. Gerar nova senha
NEW_PASSWORD=$(openssl rand -base64 32)

# 2. Atualizar secret
kubectl create secret generic bgc-db-credentials \
  --from-literal=username=bgc \
  --from-literal=password=$NEW_PASSWORD \
  --namespace=data \
  --dry-run=client -o yaml | \
  kubeseal -o yaml | \
  kubectl apply -f -

# 3. Reiniciar pods para pegar novo secret
kubectl rollout restart deployment/bgc-api -n data
kubectl rollout restart deployment/postgres -n data

# 4. Alterar senha no PostgreSQL
kubectl exec -it deployment/postgres -n data -- \
  psql -U bgc -c "ALTER USER bgc WITH PASSWORD '$NEW_PASSWORD';"
```

---

## 📝 Checklist de Segurança

### Antes de Commitar

- [ ] Verificar que `.env` está no `.gitignore`
- [ ] Confirmar que nenhum arquivo com `password` está sendo commitado
- [ ] Revisar diff do git: `git diff --cached`
- [ ] Procurar por secrets acidentais: `git grep -i password`

### Comando de Segurança

```bash
# Buscar possíveis vazamentos de credenciais
git log --all --full-history -- "*.env*" "*password*" "*secret*"

# Remover credenciais do histórico (SE NECESSÁRIO)
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch bgcstack/.env" \
  --prune-empty --tag-name-filter cat -- --all
```

---

## 🛡️ Scan de Vulnerabilidades

### Trivy - Scan de Imagens Docker

```bash
# Instalar Trivy
choco install trivy  # Windows
# ou
brew install trivy   # Mac
# ou
apt install trivy    # Linux

# Scan da imagem da API
trivy image bgc/bgc-api:dev

# Scan com severidade específica
trivy image --severity HIGH,CRITICAL bgc/bgc-api:dev

# Scan do Web Next.js
trivy image bgc/bgc-web:dev
```

### GitHub Actions - Scan Automático

Arquivo já configurado em `.github/workflows/security.yml` (a ser criado):

```yaml
name: Security Scan
on: [push, pull_request]

jobs:
  trivy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build images
        run: |
          docker build -t bgc-api:scan ./api
          docker build -t bgc-web:scan ./web-next

      - name: Run Trivy scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'bgc-api:scan'
          format: 'sarif'
          output: 'trivy-results.sarif'

      - name: Upload to GitHub Security
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'
```

---

## 🚨 Incidentes de Segurança

### Se Credenciais Forem Expostas

1. **Parar imediatamente** todos os serviços afetados
2. **Gerar novas credenciais** fortes
3. **Rotacionar** em todos os ambientes
4. **Limpar histórico** do Git (se commitado)
5. **Notificar** equipe de segurança
6. **Revisar logs** para acessos não autorizados

### Contatos de Emergência

- Equipe de Segurança: [security@bgc.dev]
- GitHub Security: https://github.com/rafamontilha/bgc-app/security

---

## 📚 Referências

- [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets)
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [OWASP Secrets Management](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [12-Factor App - Config](https://12factor.net/config)

---

**Última Atualização:** 2025-10-24
**Responsável:** DevSecOps Team
