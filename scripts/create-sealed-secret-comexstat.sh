#!/bin/bash
set -e

# Script para criar Sealed Secret para credenciais da API ComexStat
#
# Pré-requisitos:
#   - kubectl configurado e conectado ao cluster
#   - kubeseal instalado (https://github.com/bitnami-labs/sealed-secrets)
#   - Acesso ao namespace 'data'
#
# Uso:
#   ./scripts/create-sealed-secret-comexstat.sh

echo "🔐 Criando Sealed Secret para ComexStat API..."
echo

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verifica se kubeseal está instalado
if ! command -v kubeseal &> /dev/null; then
    echo -e "${RED}❌ Erro: kubeseal não está instalado${NC}"
    echo
    echo "Instale o kubeseal:"
    echo "  - macOS:   brew install kubeseal"
    echo "  - Linux:   wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/kubeseal-linux-amd64"
    echo "  - Windows: Use WSL2 ou baixe o binário manualmente"
    exit 1
fi

# Verifica se kubectl está configurado
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}❌ Erro: kubectl não está conectado a um cluster${NC}"
    exit 1
fi

# Verifica namespace
if ! kubectl get namespace data &> /dev/null; then
    echo -e "${YELLOW}⚠️  Namespace 'data' não existe. Criando...${NC}"
    kubectl create namespace data
fi

echo "📝 Preencha as credenciais da API ComexStat:"
echo

# Solicita API Key
echo -n "API Key: "
read -s API_KEY
echo

# Valida input
if [ -z "$API_KEY" ]; then
    echo -e "${RED}❌ Erro: API Key não pode estar vazia${NC}"
    exit 1
fi

echo
echo "🔨 Criando secret temporário..."

# Cria secret temporário (não será commitado)
kubectl create secret generic comexstat-credentials \
  --from-literal=api-key="$API_KEY" \
  --namespace=data \
  --dry-run=client \
  -o yaml > /tmp/comexstat-secret.yaml

echo "🔒 Selando secret com kubeseal..."

# Sela o secret
kubeseal --format yaml \
  --cert https://sealed-secrets-controller-sealed-secrets.kube-system.svc.cluster.local:8080/v1/cert.pem \
  < /tmp/comexstat-secret.yaml \
  > k8s/integration-gateway/sealed-secret-comexstat.yaml

# Se kubeseal falhar (cluster local sem sealed-secrets), tenta método offline
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}⚠️  Método online falhou. Tentando método offline...${NC}"

    # Busca certificado do controller
    kubectl get secret \
      -n kube-system \
      -l sealedsecrets.bitnami.com/sealed-secrets-key=active \
      -o jsonpath='{.items[0].data.tls\.crt}' \
      | base64 --decode > /tmp/sealed-secrets-cert.pem

    # Sela usando certificado local
    kubeseal --format yaml \
      --cert /tmp/sealed-secrets-cert.pem \
      < /tmp/comexstat-secret.yaml \
      > k8s/integration-gateway/sealed-secret-comexstat.yaml

    rm /tmp/sealed-secrets-cert.pem
fi

# Remove secret temporário
rm /tmp/comexstat-secret.yaml

echo
echo -e "${GREEN}✅ Sealed Secret criado com sucesso!${NC}"
echo
echo "Próximos passos:"
echo "  1. Verifique o arquivo: k8s/integration-gateway/sealed-secret-comexstat.yaml"
echo "  2. Commit o arquivo no git (é seguro, está criptografado)"
echo "  3. Aplique no cluster: kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml"
echo "  4. Verifique: kubectl get secret comexstat-credentials -n data"
echo

# Opção de aplicar imediatamente
read -p "Deseja aplicar o secret no cluster agora? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📦 Aplicando sealed secret no cluster..."
    kubectl apply -f k8s/integration-gateway/sealed-secret-comexstat.yaml

    echo
    echo "⏳ Aguardando descriptografia (pode levar alguns segundos)..."
    sleep 3

    # Verifica se o secret foi criado
    if kubectl get secret comexstat-credentials -n data &> /dev/null; then
        echo -e "${GREEN}✅ Secret descriptografado e disponível!${NC}"
        echo
        kubectl get secret comexstat-credentials -n data
    else
        echo -e "${YELLOW}⚠️  Secret ainda não foi descriptografado. Verifique os logs do sealed-secrets-controller${NC}"
    fi
fi

echo
echo "🎉 Pronto!"
