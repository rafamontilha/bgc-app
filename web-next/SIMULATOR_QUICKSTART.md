# Quick Start - Simulador de Destinos de Exportação

Guia rápido para rodar o simulador localmente em **5 minutos**.

---

## Pré-requisitos

- Node.js 18+ instalado
- pnpm instalado (`npm install -g pnpm`)
- Backend v0.4.0 rodando em `http://localhost:8080`

---

## Passo 1: Instalar Dependências

```bash
cd web-next
pnpm install
```

**Tempo:** ~2 minutos

---

## Passo 2: Configurar Environment

```bash
# Criar arquivo .env.local
echo "NEXT_PUBLIC_API_BASE_URL=http://localhost:8080" > .env.local
```

**Tempo:** 10 segundos

---

## Passo 3: Iniciar Backend (Terminal 1)

```bash
cd api
go run main.go
```

**Verifique:** Backend deve responder em `http://localhost:8080/healthz`

**Tempo:** 20 segundos

---

## Passo 4: Iniciar Frontend (Terminal 2)

```bash
cd web-next
pnpm dev
```

**Aguarde:** Compilação inicial (~10 segundos)

**Tempo:** 15 segundos

---

## Passo 5: Acessar Simulador

Abra o navegador:

```
http://localhost:3000/simulator
```

**Tempo:** 5 segundos

---

## Teste Rápido (Happy Path)

1. **Inserir NCM:** `17011400`
2. **Inserir Volume:** `1000`
3. **Clicar:** "Simular Destinos de Exportação"

**Resultado Esperado:**
- Loading (~1 segundo)
- Lista de 10 destinos ranqueados
- Banner: "4 de 5 simulações restantes"

---

## Troubleshooting Rápido

### Backend não responde?

```bash
# Verificar se está rodando
curl http://localhost:8080/healthz

# Se não estiver, inicie:
cd api && go run main.go
```

### Frontend não compila?

```bash
# Limpar cache e reinstalar
cd web-next
rm -rf .next node_modules
pnpm install
pnpm dev
```

### Dependências MUI não encontradas?

```bash
# Instalar manualmente
pnpm add @mui/material@^7.0.0 @mui/icons-material@^7.0.0 @emotion/react@^11.14.0 @emotion/styled@^11.14.0 zod@^3.24.1
```

---

## Endpoints Disponíveis

- **Simulador:** `/simulator`
- **Dashboard:** `/` (home)
- **Routes:** `/routes`
- **Health:** `/healthz`

---

## Arquivos de Documentação

- **README Completo:** `app/simulator/README.md`
- **Delivery Doc:** `SIMULATOR_DELIVERY.md`
- **Este Guia:** `SIMULATOR_QUICKSTART.md`

---

## Suporte

**Backend API Docs:** `http://localhost:8080/docs`
**Logs Backend:** Terminal 1
**Logs Frontend:** Terminal 2

---

**Status:** ✅ PRONTO PARA USO
**Tempo Total:** ~5 minutos
**Versão:** 1.0.0
