# Secrets Management - Integration Gateway

Este documento descreve como gerenciar secrets (credenciais, API keys, certificados) para o Integration Gateway usando **Sealed Secrets**.

## 📖 Visão Geral

O Integration Gateway usa **Kubernetes Secrets** para armazenar credenciais sensíveis de forma segura. Para permitir que esses secrets sejam versionados no Git com segurança, usamos **Sealed Secrets** (Bitnami).

### Fluxo de Segurança

```
Developer Machine              Git Repository              Kubernetes Cluster
──────────────────              ──────────────              ──────────────────

1. Plain Secret          →     2. Sealed Secret      →     3. Decrypted Secret
   (temporário)                   (criptografado)             (em memória)

   api-key: abc123              api-key: AgBcD...           api-key: abc123
   ❌ NUNCA commitar            ✅ Safe to commit           ✅ Usado pelos pods
```

## 🔐 Secrets Disponíveis

### 1. ComexStat Credentials

**Secret Name:** `comexstat-credentials`
**Namespace:** `data`
**Keys:**
- `api-key`: API Key para autenticação na API ComexStat (MDIC)

**Uso no connector:**
```yaml
# config/connectors/comexstat.yaml
auth:
  type: api_key
  api_key:
    key_ref: comexstat-credentials/api-key
```

### 2. ICP-Brasil Certificates (Exemplo)

**Secret Name:** `icp-certificates`
**Namespace:** `data`
**Keys:**
- `icp-brasil-receita-prod.pem`: Certificado X.509
- `icp-brasil-receita-prod.key`: Chave privada

## 🛠️ Como Criar um Novo Secret

### Pré-requisitos

1. **kubectl** configurado e conectado ao cluster
2. **kubeseal** instalado

```bash
# macOS
brew install kubeseal

# Linux
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/kubeseal-linux-amd64
sudo install -m 755 kubeseal-linux-amd64 /usr/local/bin/kubeseal

# Windows
# Use WSL2 ou baixe o binário manualmente
```

3. Acesso ao namespace `data`

```bash
kubectl get namespace data
```

### Método 1: Usando o Script (Recomendado)

```bash
# ComexStat
./scripts/create-sealed-secret-comexstat.sh

# O script vai:
# 1. Solicitar a API Key
# 2. Criar secret temporário
# 3. Selar com kubeseal
# 4. Salvar em k8s/integration-gateway/sealed-secret-comexstat.yaml
# 5. Opcionalmente aplicar no cluster
```

### Método 2: Manual

```bash
# 1. Criar secret temporário
kubectl create secret generic my-secret \
  --from-literal=key1=value1 \
  --from-literal=key2=value2 \
  --namespace=data \
  --dry-run=client -o yaml > /tmp/my-secret.yaml

# 2. Selar com kubeseal
kubeseal --format yaml \
  --cert https://sealed-secrets-controller-sealed-secrets.kube-system.svc.cluster.local:8080/v1/cert.pem \
  < /tmp/my-secret.yaml \
  > k8s/integration-gateway/sealed-secret-my-secret.yaml

# 3. Remover temporário
rm /tmp/my-secret.yaml

# 4. Aplicar no cluster
kubectl apply -f k8s/integration-gateway/sealed-secret-my-secret.yaml

# 5. Verificar
kubectl get secret my-secret -n data
```

## 🔄 Como Atualizar um Secret Existente

```bash
# 1. Re-executar o script
./scripts/create-sealed-secret-comexstat.sh

# 2. Commit o arquivo atualizado
git add k8s/integration-gateway/sealed-secret-comexstat.yaml
git commit -m "chore: update comexstat api key"

# 3. Aplicar no cluster
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml

# 4. Restart pods para usar novo secret
kubectl rollout restart deployment/integration-gateway -n data
```

## 🔍 Como Verificar Secrets

### Listar secrets no namespace

```bash
kubectl get secrets -n data
```

### Ver detalhes (sem valores)

```bash
kubectl describe secret comexstat-credentials -n data
```

### Ver valores (USE COM CUIDADO)

```bash
# Decodificar API key
kubectl get secret comexstat-credentials -n data -o jsonpath='{.data.api-key}' | base64 --decode
echo

# Ver todos os valores
kubectl get secret comexstat-credentials -n data -o json | jq -r '.data | map_values(@base64d)'
```

## 🚨 Segurança - O Que NUNCA Fazer

### ❌ NUNCA commitar secrets não-selados

```bash
# ERRADO - secret em plain text
apiVersion: v1
kind: Secret
metadata:
  name: comexstat-credentials
stringData:
  api-key: "minha-api-key-secreta"  # ❌ NUNCA!
```

### ❌ NUNCA adicionar secrets em arquivos de código

```go
// ERRADO
const API_KEY = "abc123"  // ❌ NUNCA!

// CERTO
apiKey, err := secretStore.GetSecret("comexstat-credentials/api-key")
```

### ❌ NUNCA logar secrets

```go
// ERRADO
log.Info("API Key: " + apiKey)  // ❌ NUNCA!

// CERTO
log.Info("Using API Key from secret")  // ✅ Sem expor valor
```

### ❌ NUNCA usar .env em produção

```bash
# ERRADO - .env em produção
API_KEY=abc123  # ❌ NUNCA em prod!

# CERTO - Kubernetes Secret
# Montado como volume ou env var automaticamente
```

## 🔐 Boas Práticas

### ✅ Rotação de Secrets

- **API Keys**: Rotacionar a cada 90 dias
- **Certificados**: Rotacionar 30 dias antes do vencimento
- **OAuth2 Secrets**: Rotacionar anualmente

### ✅ Auditoria

```bash
# Verificar quem acessou secrets
kubectl logs -n data deployment/integration-gateway | grep "GetSecret"

# Verificar sealed secrets
kubectl get sealedsecrets -n data
```

### ✅ Backup

Sealed Secrets estão no Git, mas faça backup adicional:

```bash
# Exportar sealed secret
kubectl get sealedsecret comexstat-credentials -n data -o yaml > backup/comexstat-sealed-secret.yaml

# Exportar private key do sealed-secrets-controller (MUITO SENSÍVEL!)
kubectl get secret -n kube-system -l sealedsecrets.bitnami.com/sealed-secrets-key -o yaml > backup/sealed-secrets-key.yaml
```

### ✅ Teste em Staging Primeiro

```bash
# 1. Aplicar em staging
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml --context=staging

# 2. Testar integração
curl http://integration-gateway.staging/v1/connectors/comexstat/health

# 3. Se OK, aplicar em produção
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml --context=production
```

## 🆘 Troubleshooting

### Secret não está sendo descriptografado

```bash
# 1. Verificar logs do sealed-secrets-controller
kubectl logs -n kube-system -l app.kubernetes.io/name=sealed-secrets

# 2. Verificar se o SealedSecret existe
kubectl get sealedsecret comexstat-credentials -n data

# 3. Deletar e recriar
kubectl delete sealedsecret comexstat-credentials -n data
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml
```

### Pod não consegue acessar secret

```bash
# 1. Verificar se secret existe
kubectl get secret comexstat-credentials -n data

# 2. Verificar RBAC
kubectl auth can-i get secrets --namespace=data --as=system:serviceaccount:data:default

# 3. Verificar se está montado no pod
kubectl describe pod integration-gateway-xxx -n data | grep -A5 "Volumes:"
```

### Valor do secret está errado

```bash
# 1. Deletar secret
kubectl delete secret comexstat-credentials -n data

# 2. Deletar sealed secret
kubectl delete sealedsecret comexstat-credentials -n data

# 3. Recriar
./scripts/create-sealed-secret-comexstat.sh

# 4. Aplicar
kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml

# 5. Restart pods
kubectl rollout restart deployment/integration-gateway -n data
```

## 📚 Referências

- [Sealed Secrets - GitHub](https://github.com/bitnami-labs/sealed-secrets)
- [Kubernetes Secrets - Docs](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Best Practices for Secrets Management](https://kubernetes.io/docs/concepts/security/secrets-good-practices/)
