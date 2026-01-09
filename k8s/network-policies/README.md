# Network Policies - BGC App

Network Policies implementam **segmentação de rede** no nível do Kubernetes, garantindo que apenas tráfego autorizado flua entre pods e serviços.

## 🎯 Objetivo

**Princípio de Menor Privilégio (Least Privilege)**:
- Cada pod tem acesso APENAS ao que precisa
- Bloqueio padrão de todo tráfego não explicitamente permitido
- Isolamento entre serviços
- Força integrações externas através do Integration Gateway

## 🗺️ Arquitetura de Rede

```
┌──────────────────────────────────────────────────────────────┐
│  Internet                                                     │
└────────────────────┬─────────────────────────────────────────┘
                     │ HTTPS
                     ↓
┌──────────────────────────────────────────────────────────────┐
│  Ingress Controller (nginx/traefik)                          │
│  Namespace: ingress-nginx                                    │
└────────────────────┬─────────────────────────────────────────┘
                     │
                     ↓
        ┌────────────────────────┐
        │  bgc-api               │ ← Prometheus (metrics)
        │  Namespace: data       │
        │  Port: 8080            │
        └─┬──────────────────────┘
          │
          ├─→ PostgreSQL (5432)
          ├─→ Redis (6379)
          ├─→ Integration Gateway (8081) ← OBRIGATÓRIO para APIs externas
          └─→ Jaeger (4317/4318)
                     │
                     ↓
        ┌────────────────────────┐
        │  Integration Gateway   │ ← Prometheus (metrics)
        │  Namespace: data       │ ← bgc-api APENAS
        │  Port: 8081            │
        └─┬──────────────────────┘
          │
          ├─→ Redis (6379)
          ├─→ PostgreSQL (5432)
          ├─→ Kubernetes API (443) ← Para buscar Secrets
          ├─→ Jaeger (4317/4318)
          └─→ External APIs (443) ← ComexStat, ViaCEP, Receita, etc.
```

## 📋 Policies Disponíveis

### 1. `default-deny-all.yaml`
**Escopo:** Namespace `data`
**Efeito:** Bloqueia TODO tráfego por padrão

Pods precisam de NetworkPolicy explícita para funcionar.

### 2. `bgc-api-netpol.yaml`
**Escopo:** Pod `bgc-api`

**Ingress permitido:**
- ✅ Ingress Controller → bgc-api:8080
- ✅ Prometheus → bgc-api:9090

**Egress permitido:**
- ✅ bgc-api → DNS (53)
- ✅ bgc-api → PostgreSQL (5432)
- ✅ bgc-api → Redis (6379)
- ✅ bgc-api → Integration Gateway (8081)
- ✅ bgc-api → Jaeger (4317/4318)
- ❌ bgc-api → APIs Externas (BLOQUEADO - deve usar Gateway)

### 3. `integration-gateway-netpol.yaml`
**Escopo:** Pod `integration-gateway`

**Ingress permitido:**
- ✅ bgc-api → gateway:8081
- ✅ Prometheus → gateway:9090

**Egress permitido:**
- ✅ gateway → DNS (53)
- ✅ gateway → PostgreSQL (5432)
- ✅ gateway → Redis (6379)
- ✅ gateway → APIs Externas HTTPS (443)
- ✅ gateway → Jaeger (4317/4318)
- ✅ gateway → Kubernetes API (443)

### 4. `redis-netpol.yaml` e `postgres-netpol.yaml`
**Escopo:** Redis e PostgreSQL

**Ingress permitido:**
- ✅ bgc-api → Redis/PostgreSQL
- ✅ integration-gateway → Redis/PostgreSQL
- ❌ Outros pods → BLOQUEADO

## 🚀 Como Aplicar

### Pré-requisitos

1. **CNI Plugin com suporte a Network Policies**:
   - Calico ✅
   - Cilium ✅
   - Weave Net ✅
   - Flannel ❌ (não suporta)

Verifique seu cluster:
```bash
kubectl get nodes -o wide
# Verifique a coluna CONTAINER-RUNTIME e CNI
```

### Aplicar Policies

```bash
# 1. Aplicar default deny (CUIDADO: pode quebrar pods existentes!)
kubectl apply -f k8s/network-policies/bgc-api-netpol.yaml

# 2. Aplicar policies específicas
kubectl apply -f k8s/network-policies/integration-gateway-netpol.yaml

# 3. Aplicar todas de uma vez
kubectl apply -f k8s/network-policies/

# 4. Verificar
kubectl get networkpolicies -n data
```

### Aplicar em Staging Primeiro

```bash
# Staging
kubectl apply -f k8s/network-policies/ --context=staging

# Testar por 24-48h

# Produção
kubectl apply -f k8s/network-policies/ --context=production
```

## 🧪 Como Testar

### Teste 1: bgc-api → Integration Gateway (deve funcionar)

```bash
# Shell no pod da API
kubectl exec -it deployment/bgc-api -n data -- sh

# Teste conexão com gateway
curl http://integration-gateway:8081/health
# Esperado: 200 OK
```

### Teste 2: bgc-api → APIs Externas (deve FALHAR)

```bash
# Shell no pod da API
kubectl exec -it deployment/bgc-api -n data -- sh

# Tenta acessar ComexStat diretamente
curl -I https://api.comexstat.mdic.gov.br
# Esperado: timeout ou connection refused (BLOQUEADO)
```

### Teste 3: Integration Gateway → APIs Externas (deve funcionar)

```bash
# Shell no pod do gateway
kubectl exec -it deployment/integration-gateway -n data -- sh

# Teste conexão externa
curl -I https://api.comexstat.mdic.gov.br
# Esperado: 200 OK
```

### Teste 4: Pod Aleatório → Redis (deve FALHAR)

```bash
# Criar pod de teste sem network policy
kubectl run test-pod --image=busybox -n data --rm -it -- sh

# Tenta conectar no Redis
nc -zv redis 6379
# Esperado: connection timed out (BLOQUEADO)
```

### Teste 5: Prometheus → Metrics (deve funcionar)

```bash
# Shell no pod do Prometheus
kubectl exec -it deployment/prometheus -n observability -- sh

# Teste scrape
curl http://bgc-api.data:9090/metrics
# Esperado: Métricas Prometheus
```

## 🔍 Debugging

### Verificar Policies Aplicadas

```bash
# Listar todas as policies
kubectl get networkpolicies -n data

# Ver detalhes de uma policy
kubectl describe networkpolicy bgc-api-netpol -n data

# Ver policies de um pod específico
kubectl get networkpolicies -n data -o json | jq '.items[] | select(.spec.podSelector.matchLabels.app=="bgc-api")'
```

### Logs de Tráfego Bloqueado

Depende do CNI:

**Calico:**
```bash
# Habilitar log de tráfego negado
kubectl apply -f - <<EOF
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: default-deny-log
spec:
  selector: all()
  types:
  - Ingress
  - Egress
  ingress:
  - action: Log
  - action: Deny
  egress:
  - action: Log
  - action: Deny
EOF

# Ver logs
kubectl logs -n kube-system -l k8s-app=calico-node | grep -i deny
```

**Cilium:**
```bash
# Ver fluxos bloqueados
cilium hubble observe --verdict DROPPED --namespace data
```

### Testar Conectividade

```bash
# Criar pod de teste
kubectl run netpol-test --image=nicolaka/netshoot -n data --rm -it -- bash

# Dentro do pod:
# Teste DNS
nslookup google.com

# Teste Redis
nc -zv redis 6379

# Teste PostgreSQL
nc -zv postgres 5432

# Teste Integration Gateway
curl http://integration-gateway:8081/health

# Teste API externa (deve falhar)
curl -I https://api.comexstat.mdic.gov.br
```

## 🚨 Troubleshooting

### Pod não consegue acessar serviço necessário

```bash
# 1. Verificar se há network policy aplicada
kubectl get networkpolicy -n data

# 2. Ver detalhes da policy do pod de origem
kubectl describe networkpolicy <policy-name> -n data

# 3. Verificar labels do pod
kubectl get pod <pod-name> -n data --show-labels

# 4. Ver logs do CNI
kubectl logs -n kube-system -l app=<cni-name>  # calico, cilium, etc
```

### Prometheus não consegue scrape métricas

```bash
# 1. Verificar namespace do Prometheus
kubectl get namespace observability --show-labels

# 2. Adicionar label se necessário
kubectl label namespace observability name=observability

# 3. Verificar policy
kubectl describe networkpolicy bgc-api-netpol -n data | grep -A10 "Allowing ingress"
```

### Tráfego legítimo está bloqueado

```bash
# 1. Identificar pods de origem e destino
kubectl get pods -n data -o wide

# 2. Verificar labels
kubectl get pod <source-pod> -n data --show-labels
kubectl get pod <destination-pod> -n data --show-labels

# 3. Adicionar regra na policy apropriada
# Editar arquivo YAML e aplicar novamente
```

## 🔐 Boas Práticas

### ✅ Sempre fazer

1. **Testar em Staging primeiro**
2. **Aplicar policies gradualmente** (não todas de uma vez)
3. **Monitorar logs** por 24-48h após aplicação
4. **Documentar exceções** (se necessário permitir tráfego incomum)
5. **Revisar policies trimestralmente**

### ❌ Nunca fazer

1. **Aplicar default-deny sem policies específicas** (vai quebrar tudo)
2. **Permitir egress 0.0.0.0/0 sem justificativa** (derrota o propósito)
3. **Usar `podSelector: {}` sem pensar** (muito permissivo)
4. **Esquecer de testar health checks** (pode quebrar probes do K8s)

## 📊 Monitoramento

### Métricas Recomendadas

```promql
# Drops por network policy (Calico)
rate(calico_denied_packets[5m])

# Conexões bloqueadas por policy (Cilium)
cilium_policy_verdict_total{verdict="DROPPED"}

# Tempo de resposta (detectar timeouts)
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

### Alertas

```yaml
# Prometheus Alert Rules
groups:
- name: network-policies
  rules:
  - alert: HighNetworkPolicyDrops
    expr: rate(calico_denied_packets[5m]) > 100
    for: 5m
    annotations:
      summary: "High rate of network policy drops detected"
      description: "{{ $value }} packets/sec being dropped by network policies"
```

## 📚 Referências

- [Kubernetes Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Calico Network Policy](https://docs.projectcalico.org/security/kubernetes-network-policy)
- [Cilium Network Policy](https://docs.cilium.io/en/stable/policy/)
- [Network Policy Recipes](https://github.com/ahmetb/kubernetes-network-policy-recipes)
