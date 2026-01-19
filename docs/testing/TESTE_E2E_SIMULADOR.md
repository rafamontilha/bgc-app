# Guia de Teste E2E - Simulador de Destinos de Exportação

**Data:** 2026-01-10
**Versão:** v0.4.0 (Backend) + Frontend UI
**Status:** Pronto para testes integrados

---

## Pré-requisitos

### Backend (API Go)
- Go 1.24.9 instalado
- PostgreSQL com dados de teste (50 countries, 16 export records)
- Redis L2 cache (opcional, mas recomendado)

### Frontend (Next.js)
- Node.js 18+ instalado
- pnpm 9.15.9+
- Dependências instaladas (já feito)

---

## Iniciar Serviços

### Terminal 1: Backend API

```bash
cd C:\Users\rafae\OneDrive\Documentos\Projetos\Brasil Global Conect\bgc-app\api
go run main.go
```

**Verificação:**
```
✓ Server started on :8080
✓ Database connected
✓ Redis connected (opcional)
```

### Terminal 2: Frontend Next.js

```bash
cd C:\Users\rafae\OneDrive\Documentos\Projetos\Brasil Global Conect\bgc-app\web-next
pnpm dev
```

**Verificação:**
```
✓ Ready on http://localhost:3000
✓ Compiled in XXXms
```

---

## Casos de Teste E2E

### 1. Happy Path (Fluxo Principal)

**Objetivo:** Validar simulação básica funcionando end-to-end.

**Passos:**
1. Abrir navegador em `http://localhost:3000/simulator`
2. Inserir NCM: `17011400` (Açúcar de cana)
3. Inserir Volume: `1000` kg
4. Clicar "Simular Destinos de Exportação"

**Resultado Esperado:**
- ✅ Loading spinner visível (< 1 segundo)
- ✅ Lista de 10 destinos aparece
- ✅ Destinos ordenados por score (descendente)
- ✅ Cada card mostra:
  - Rank badge (1, 2, 3...)
  - Bandeira do país
  - Nome do país (português)
  - Score visual (barra de progresso)
  - Demand level (chip colorido)
  - 4 métricas principais
- ✅ Banner no topo: "4 de 5 simulações restantes hoje"

**Performance Esperada:**
- API response time: 22-92ms (backend validado)
- FCP (First Contentful Paint): < 1.5s
- LCP (Largest Contentful Paint): < 2.5s

---

### 2. Rate Limiting (Limite Freemium)

**Objetivo:** Validar que o limite de 5 simulações/dia funciona.

**Passos:**
1. Executar 5 simulações consecutivas (pode usar NCMs diferentes)
2. Na 6ª tentativa, observar comportamento

**Resultado Esperado:**
- ✅ Simulações 1-4: Banner atualiza corretamente ("X de 5 restantes")
- ✅ Simulação 5: Banner mostra "0 de 5 restantes"
- ✅ Simulação 6:
  - Modal de upgrade aparece
  - Título: "Limite Atingido!"
  - Mensagem: "Você atingiu o limite de 5 simulações diárias..."
  - Botão "Ver Planos"
  - Botão "Agora Não" (fecha modal)

**Backend Esperado:**
- HTTP 429 Too Many Requests
- Header `X-RateLimit-Remaining: 0`

---

### 3. Validação de Entrada (NCM Inválido)

**Objetivo:** Validar que erros de validação são mostrados claramente.

**Cenários:**

#### 3a. NCM com menos de 8 dígitos
**Input:** `1701` (4 dígitos)
**Esperado:**
- ✅ Mensagem de erro abaixo do campo: "NCM deve ter exatamente 8 dígitos"
- ✅ Botão "Simular" desabilitado
- ✅ Campo com borda vermelha

#### 3b. NCM com caracteres não numéricos
**Input:** `17011ABC` (letras)
**Esperado:**
- ✅ Mensagem de erro: "NCM deve conter apenas números"
- ✅ Botão "Simular" desabilitado

#### 3c. NCM válido mas sem dados
**Input:** `99999999` (não existe)
**Esperado:**
- ✅ Request enviado ao backend
- ✅ ErrorState aparece com mensagem: "NCM não encontrado"
- ✅ Dica: "Verifique se o NCM está correto..."
- ✅ Botão "Tentar Novamente"

---

### 4. Validação de Volume

**Objetivo:** Validar que volumes inválidos são rejeitados.

**Cenários:**

#### 4a. Volume negativo
**Input:** `-1000`
**Esperado:**
- ✅ Mensagem de erro: "Volume deve ser maior que zero"
- ✅ Botão desabilitado

#### 4b. Volume zero
**Input:** `0`
**Esperado:**
- ✅ Mensagem de erro: "Volume deve ser maior que zero"

#### 4c. Volume muito grande
**Input:** `999999999` (1 bilhão kg)
**Esperado:**
- ✅ Aceito (sem limite máximo no frontend)
- ✅ Simulação executada normalmente

---

### 5. Erro de Servidor (Backend Offline)

**Objetivo:** Validar que erros de servidor são tratados graciosamente.

**Passos:**
1. **PARAR** o backend (Ctrl+C no Terminal 1)
2. Inserir NCM válido: `17011400`
3. Clicar "Simular"

**Resultado Esperado:**
- ✅ Loading spinner aparece
- ✅ Após timeout (~5-10s), ErrorState aparece:
  - Tipo: "network_error"
  - Mensagem: "Não foi possível conectar ao servidor"
  - Dica: "Verifique sua conexão..."
  - Botão "Tentar Novamente"

**Reiniciar backend e tentar novamente:**
- ✅ Simulação funciona normalmente

---

### 6. Responsividade (Mobile/Tablet/Desktop)

**Objetivo:** Validar que UI funciona bem em diferentes tamanhos de tela.

**Dispositivos para Testar:**

#### Mobile (375px - iPhone 12)
- ✅ Formulário ocupa largura total
- ✅ Botão "Simular" em largura total
- ✅ Cards de destino empilhados (1 por linha)
- ✅ Métricas em grid 2x2
- ✅ Banner de rate limit visível e legível
- ✅ Modal de upgrade centralizado

#### Tablet (768px - iPad)
- ✅ Layout similar ao mobile, espaçamento maior
- ✅ Cards de destino ainda empilhados
- ✅ Fonte ligeiramente maior

#### Desktop (1280px+)
- ✅ Container centralizado (max-width)
- ✅ Cards de destino empilhados mas com hover effects
- ✅ Espaçamento generoso (8px grid)
- ✅ Tooltips aparecem no hover

**Como Testar:**
- Chrome DevTools: F12 → Toggle Device Toolbar (Ctrl+Shift+M)
- Alternar entre presets: iPhone 12, iPad, Desktop

---

### 7. Acessibilidade (WCAG 2.1 AA)

**Objetivo:** Validar que UI é acessível para todos os usuários.

**Testes:**

#### 7a. Navegação por Teclado
1. Tab através dos elementos
2. Verificar:
   - ✅ Focus indicators visíveis
   - ✅ Ordem lógica de tab (form → button → cards)
   - ✅ Enter submete formulário
   - ✅ Esc fecha modal

#### 7b. Screen Reader (Opcional)
1. Ativar Narrator (Windows) ou VoiceOver (Mac)
2. Verificar:
   - ✅ Labels lidos corretamente
   - ✅ Erros anunciados
   - ✅ Loading state anunciado

#### 7c. Contraste de Cores
1. Usar Lighthouse (Chrome DevTools)
2. Verificar:
   - ✅ Contraste mínimo 4.5:1 para texto normal
   - ✅ Contraste mínimo 3:1 para texto grande
   - ✅ Score de acessibilidade > 90

---

### 8. Performance (Lighthouse)

**Objetivo:** Validar que UI tem boa performance.

**Passos:**
1. Abrir Chrome DevTools (F12)
2. Lighthouse tab
3. Selecionar:
   - ✅ Performance
   - ✅ Accessibility
   - ✅ Best Practices
   - ✅ SEO (opcional)
4. Gerar relatório

**Métricas Esperadas:**
- Performance: > 90
- Accessibility: > 90
- Best Practices: > 90
- FCP (First Contentful Paint): < 1.5s
- LCP (Largest Contentful Paint): < 2.5s
- TBT (Total Blocking Time): < 200ms
- CLS (Cumulative Layout Shift): < 0.1

---

## Casos de Teste com Dados Reais

### NCMs de Teste (com dados no backend)

| NCM | Produto | Destinos Esperados | Observações |
|-----|---------|-------------------|-------------|
| `17011400` | Açúcar de cana | 5 destinos | US, CN, DE, IN, GB |
| `02013000` | Carne bovina | 2 destinos | CN, HK |

**Teste:**
1. Simular cada NCM acima
2. Verificar se número de destinos corresponde
3. Verificar se países correspondem

---

## Checklist de Validação Final

### Frontend
- [ ] ✅ Compilação TypeScript sem erros
- [ ] ✅ Build production bem-sucedido
- [ ] ✅ Zero warnings críticos (apenas 'request' unused em healthz)
- [ ] ✅ Dependências instaladas corretamente
- [ ] ✅ Tema MUI v7 aplicado
- [ ] ✅ Apple aesthetic visível (pill buttons, espaçamento)

### Backend
- [ ] ✅ API respondendo em `:8080`
- [ ] ✅ PostgreSQL com 50 countries
- [ ] ✅ Dados de teste (16 export records)
- [ ] ✅ Rate limiting funcional (5 req/dia)
- [ ] ✅ Headers `X-RateLimit-*` presentes

### Integração
- [ ] Happy path (NCM válido → 10 destinos)
- [ ] Rate limiting (6ª simulação → modal upgrade)
- [ ] Validação de entrada (NCM/volume inválidos)
- [ ] Error handling (backend offline)
- [ ] Responsividade (mobile/tablet/desktop)
- [ ] Acessibilidade (keyboard nav, screen reader)
- [ ] Performance (Lighthouse > 90)

---

## Troubleshooting

### Frontend não inicia
**Erro:** `pnpm dev` falha
**Solução:**
```bash
cd web-next
pnpm install
pnpm dev
```

### Backend não conecta ao banco
**Erro:** `connection refused`
**Solução:**
```bash
# Verificar se PostgreSQL está rodando
docker ps | grep postgres

# Se não estiver, iniciar:
cd bgcstack
docker-compose up -d postgres
```

### API retorna 404
**Erro:** `/v1/simulator/destinations` não encontrado
**Solução:**
- Verificar se backend está na versão v0.4.0
- Verificar logs do backend para erros de rota
- Confirmar que migração 0010/0011 foi executada

### Modal de upgrade não aparece
**Erro:** 6ª simulação não abre modal
**Solução:**
- Verificar no Network tab do Chrome se HTTP 429 está sendo retornado
- Verificar se header `X-RateLimit-Remaining: 0` está presente
- Limpar localStorage/sessionStorage do navegador

### Cards de destino não aparecem
**Erro:** Loading spinner infinito
**Solução:**
- Verificar Network tab: request está sendo enviado?
- Verificar Console tab: erros JavaScript?
- Verificar response do backend: estrutura JSON correta?

---

## Próximos Passos Após Validação

1. ✅ **Marcar Epic 4 como 100% completo** (frontend + backend)
2. 📝 **Documentar bugs encontrados** (se houver)
3. 🚀 **Preparar para Beta Privado:**
   - Recrutar 20 exportadores (LinkedIn, email)
   - Setup de analytics (Amplitude, Mixpanel)
   - Criar onboarding flow
4. 💰 **Iniciar Epic 6 (Monetização):**
   - Van Westendorp pricing research
   - Integração Stripe
   - Pricing page

---

## Contatos & Suporte

**Documentação Completa:**
- Frontend: `web-next/app/simulator/README.md`
- Backend: `docs/API-SIMULATOR.md`
- Arquitetura: `web-next/SIMULATOR_ARCHITECTURE.md`

**Logs:**
- Backend: stdout (terminal 1)
- Frontend: Chrome DevTools Console

**Issues:**
- Reportar em: `CHANGELOG.md` (adicionar seção Issues)

---

**Versão:** 1.0
**Última Atualização:** 2026-01-10
**Status:** Pronto para testes ✅
