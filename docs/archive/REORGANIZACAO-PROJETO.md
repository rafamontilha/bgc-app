# Reorganização do Projeto BGC - Next.js como Frontend Principal

**Data:** 21/10/2025
**Autor:** Claude Code
**Status:** ✅ Concluído

## 📋 Objetivo

Reorganizar o projeto para manter **apenas o Next.js** como frontend ativo, movendo ou removendo os frontends legados (HTML estático com nginx).

## 🎯 Mudanças Realizadas

### 1. Estrutura de Diretórios

#### Antes:
```
bgc-app/
├── web/              # Frontend HTML/CSS/JS antigo
├── web-next/         # Frontend Next.js novo
└── k8s/
    ├── web.yaml                    # Deployment do nginx
    ├── web-nginx-configmap.yaml    # ConfigMap do nginx (com proxy)
    ├── web-hpa.yaml                # HPA do web antigo
    └── web-next.yaml               # Deployment do Next.js
```

#### Depois:
```
bgc-app/
├── web-next/         # Frontend Next.js (mantido como está)
├── old/
│   ├── web-legacy-html/           # Frontend HTML antigo (movido)
│   ├── web-nginx-configmap-legacy.yaml
│   ├── web-legacy-k8s.yaml
│   └── web-hpa-legacy.yaml
└── k8s/
    └── web.yaml      # Deployment Next.js renomeado
```

### 2. Kubernetes - Recursos Atualizados

#### Deletados:
- `bgc-web-next` deployment/service/ingress/hpa
- `bgc-web-nginx-config` configmap
- Todos os pods do nginx legacy

#### Criados/Atualizados:
- **`k8s/web.yaml`**: Novo manifesto consolidado
  - Deployment: `bgc-web` (antes era `bgc-web-next`)
  - Service: `bgc-web` na porta 3000
  - Ingress: `web.bgc.local` → `bgc-web:3000`
  - HPA: 2-5 réplicas

#### Imagem Docker:
- **Nova**: `bgc/bgc-web:latest` (construída do web-next)
- **Antiga**: `bgc/bgc-web-next:latest` (não mais usada)

### 3. Configuração de Acesso

#### URLs Atuais:
- **Frontend**: http://web.bgc.local → Next.js direto (porta 3000)
- **API**: http://api.bgc.local → Go API (porta 8080)

#### Sem mais proxy nginx!
Antes havia um nginx fazendo proxy:
- `web.bgc.local/` → `bgc-web-next:3000`
- `web.bgc.local/market/*` → `bgc-api:8080`

Agora:
- Next.js responde diretamente em `web.bgc.local`
- Next.js usa `rewrites` internos para fazer proxy para a API

### 4. Arquivo hosts

O usuário já tem configurado:
```
127.0.0.1  api.bgc.local
127.0.0.1  web.bgc.local
```

Não é mais necessário `web-next.bgc.local`.

## 🔧 Como Aplicar em Outro Ambiente

### 1. Construir a imagem:
```powershell
docker build -t bgc/bgc-web:latest -f web-next/Dockerfile ./web-next
```

### 2. Importar para k3d:
```powershell
k3d image import bgc/bgc-web:latest -c bgc
```

### 3. Deletar recursos antigos:
```powershell
kubectl delete deployment bgc-web-next -n data
kubectl delete svc bgc-web-next -n data
kubectl delete ingress bgc-web-next -n data
kubectl delete hpa bgc-web-next-hpa -n data
kubectl delete configmap bgc-web-nginx-config -n data
```

### 4. Aplicar novos manifestos:
```powershell
kubectl apply -f k8s/web.yaml
```

### 5. Verificar:
```powershell
kubectl get pods -n data | grep bgc-web
curl http://web.bgc.local/
```

## 📊 Status Final

### Pods em Execução:
```
NAME                       READY   STATUS    RESTARTS   AGE
bgc-api-5f65899b47-wztz7   1/1     Running   3          4d17h
bgc-web-78458f5ff6-br9w8   1/1     Running   0          3m
bgc-web-78458f5ff6-f5wss   1/1     Running   0          3m
```

### Serviços:
```
NAME      TYPE        CLUSTER-IP     PORT(S)
bgc-api   ClusterIP   10.43.177.51   8080/TCP
bgc-web   ClusterIP   10.43.83.247   3000/TCP
```

### Ingress:
```
NAME       CLASS     HOSTS              ADDRESS
bgc-api    traefik   api.bgc.local      172.18.0.2
bgc-web    traefik   web.bgc.local      172.18.0.2
```

### HPA:
```
NAME          REFERENCE          TARGETS                  MIN   MAX
bgc-api-hpa   Deployment/bgc-api cpu: 1%/70%, mem: 7%/80% 1     5
bgc-web-hpa   Deployment/bgc-web cpu: 1%/70%, mem: 33%/80% 2    5
```

## ✅ Benefícios

1. **Simplicidade**: Um único frontend (Next.js), sem nginx intermediário
2. **Performance**: Menos camadas de proxy
3. **Manutenção**: Código mais limpo, menos manifestos K8s
4. **Consistência**: `bgc-web` para frontend, `bgc-api` para backend
5. **Escalabilidade**: HPA já configurado para 2-5 réplicas

## 📝 Arquivos Alterados

### Criados:
- `k8s/web.yaml` (novo, consolidado)
- `docs/REORGANIZACAO-PROJETO.md` (este arquivo)

### Modificados:
- `README.md` (atualizada estrutura do projeto)

### Movidos para `old/`:
- `web/` → `old/web-legacy-html/`
- `k8s/web.yaml` → `old/web-legacy-k8s.yaml`
- `k8s/web-nginx-configmap.yaml` → `old/web-nginx-configmap-legacy.yaml`
- `k8s/web-hpa.yaml` → `old/web-hpa-legacy.yaml`
- `k8s/web-next.yaml` → `old/web-next-legacy.yaml`

### Deletados do Kubernetes:
- Todos os recursos `bgc-web-next-*`
- ConfigMap `bgc-web-nginx-config`

## 🚀 Próximos Passos Sugeridos

1. ✅ Testar todas as funcionalidades no navegador
2. ✅ Verificar que `/routes` funciona corretamente
3. ⬜ Atualizar scripts de deploy automático
4. ⬜ Criar CI/CD pipeline para build automático
5. ⬜ Documentar processo de desenvolvimento local do Next.js

## 🔗 Documentação Relacionada

- [QUICK-START.md](QUICK-START.md)
- [SETUP-NEXTJS.md](SETUP-NEXTJS.md)
- [TROUBLESHOOTING-NEXTJS.md](TROUBLESHOOTING-NEXTJS.md)
- [README.md](../README.md)

---

**Nota**: Todos os arquivos legados foram preservados em `old/` para referência histórica.
