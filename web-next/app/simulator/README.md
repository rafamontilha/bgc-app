# Export Destination Simulator - Frontend UI

Simulador de Destinos de Exportação para Brasil Global Connect (v0.4.0).

## Visão Geral

Interface de usuário completa para o **Simulador de Destinos de Exportação**, permitindo que empresas descubram os melhores mercados para seus produtos baseado em NCM.

## Tecnologias Utilizadas

- **Framework:** Next.js 15 (App Router) + React 19
- **Linguagem:** TypeScript 5 (strict mode)
- **Design System:** Material-UI v7 (Material Design 3)
- **Estética:** Apple Human Interface Guidelines
- **Validação:** Zod
- **State Management:** React Hooks (useState, useEffect)

## Componentes Implementados

### 1. SimulatorForm (US-001)
**Localização:** `components/simulator/SimulatorForm.tsx`

Input form com validação client-side usando Zod:
- Campo NCM (8 dígitos, validação regex)
- Campo Volume (kg, validação min/max)
- Loading state durante requisição
- Error states com mensagens claras

### 2. DestinationCard (US-002)
**Localização:** `components/simulator/DestinationCard.tsx`

Card individual de destino com:
- Rank badge visual
- Bandeira do país (emoji)
- Score visual (progress bar 0-10)
- Demand level chip (Alto/Médio/Baixo com cores)
- Métricas principais: Market Size, Growth Rate, Price/kg, Distance
- Recommendation reason em destaque
- Detalhes adicionais: Margin, Logistics Cost, Tariff, Lead Time
- Hover effects suaves

### 3. DestinationList (US-002)
**Localização:** `components/simulator/DestinationList.tsx`

Container para lista de destinos:
- Loading skeleton (3 cards)
- Empty state quando sem resultados
- Listagem ordenada por score
- Responsivo (mobile + desktop)

### 4. RateLimitBanner (US-003)
**Localização:** `components/simulator/RateLimitBanner.tsx`

Indicador de rate limit freemium:
- Banner info: "X de 5 simulações restantes"
- Estado warning quando < 2 simulações
- Estado crítico quando 0 simulações
- Exibe tempo até reset
- Botão "Upgrade" quando low/critical

### 5. UpgradeModal (US-004)
**Localização:** `components/simulator/UpgradeModal.tsx`

Modal de upgrade para plano Pro:
- Trigger: HTTP 429 (rate limit exceeded)
- Lista de benefícios Pro (6 features)
- Preço "A partir de R$ 97/mês"
- CTAs: "Agora Não" e "Ver Planos"
- Design Apple-inspired (clean, espaçoso)

### 6. ErrorState (US-005)
**Localização:** `components/simulator/ErrorState.tsx`

Estados de erro user-friendly:
- `not_found`: NCM não encontrado
- `validation`: Dados inválidos
- `server_error`: Erro 500+
- `network_error`: Sem conexão
- `timeout`: Request timeout
- `unknown`: Erro genérico

Cada tipo tem:
- Ícone contextual
- Título e mensagem clara
- Botão "Tentar Novamente"
- Dicas adicionais (quando aplicável)

### 7. SimulatorPage
**Localização:** `app/simulator/page.tsx`

Página principal que orquestra todos os componentes:
- Gerencia estado global (loading, data, error, rateLimitInfo)
- Faz chamadas à API via `lib/api/simulator.ts`
- Renderiza componentes condicionalmente
- Trata erros e exibe mensagens apropriadas
- Abre upgrade modal quando HTTP 429

## API Integration

**Endpoint:** `POST /v1/simulator/destinations`

**Request:**
```typescript
{
  ncm: string;          // 8 dígitos
  volume_kg: number;    // Positivo
  countries?: string[]; // Opcional
  max_results?: number; // Opcional (default: 10)
}
```

**Response:**
```typescript
{
  destinations: Destination[];
  metadata: {
    analysis_date: string;
    processing_time_ms: number;
    cache_hit: boolean;
  }
}
```

**Rate Limit Headers:**
- `X-RateLimit-Limit`: Total permitido
- `X-RateLimit-Remaining`: Restante
- `X-RateLimit-Reset`: Unix timestamp do reset

## Tema MUI (Apple Aesthetic)

**Localização:** `lib/theme.ts`

### Cores
- **Primary:** `#007AFF` (Apple Blue)
- **Background:** `#F5F5F7` (Light neutral)
- **Text Primary:** `#1C1C1E` (Dark neutral)
- **Success:** `#34C759` (Green)
- **Warning:** `#FF9500` (Orange)
- **Error:** `#FF3B30` (Red)

### Tipografia
- **Font Family:** `-apple-system, BlinkMacSystemFont, "SF Pro Text", system-ui, sans-serif`
- **Weights:** 400 (regular), 500 (medium), 600 (semibold)

### Componentes
- **Buttons:** Pill-shaped (borderRadius: 999), hover scale 1.02
- **Cards:** borderRadius: 16, soft shadow, hover lift
- **TextFields:** borderRadius: 12, smooth transitions
- **Chips:** borderRadius: 8, alpha backgrounds

### Animações
- **Transition:** `all 220ms cubic-bezier(0.32, 0.72, 0, 1)`
- **Hover:** Suave scale e shadow increase
- **Focus:** Clear visual indicators

## Estrutura de Arquivos

```
web-next/
├── app/
│   └── simulator/
│       ├── page.tsx           # Página principal
│       └── README.md          # Esta documentação
├── components/
│   ├── providers/
│   │   └── ThemeProvider.tsx  # MUI theme wrapper
│   └── simulator/
│       ├── SimulatorForm.tsx       # US-001
│       ├── DestinationCard.tsx     # US-002
│       ├── DestinationList.tsx     # US-002
│       ├── RateLimitBanner.tsx     # US-003
│       ├── UpgradeModal.tsx        # US-004
│       └── ErrorState.tsx          # US-005
├── lib/
│   ├── theme.ts               # MUI v7 theme
│   └── api/
│       └── simulator.ts       # API client + validations
├── types/
│   └── simulator.ts           # TypeScript interfaces
└── .env.local                 # Environment variables
```

## Variáveis de Ambiente

```bash
# .env.local
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

**Produção:** Alterar para URL do backend em produção.

## Como Usar

### 1. Navegar para o Simulador
```
http://localhost:3000/simulator
```

### 2. Preencher o Formulário
- **NCM:** Código de 8 dígitos (ex: 17011400)
- **Volume:** Quantidade em kg (ex: 1000)

### 3. Simular
- Clicar em "Simular Destinos de Exportação"
- Aguardar processamento (22-92ms típico)
- Visualizar resultados ranqueados

### 4. Interpretar Resultados
- **Rank:** Posição no ranking (1 = melhor)
- **Score:** Adequação 0-10 (visual progress bar)
- **Demand Level:** Alto (verde), Médio (amarelo), Baixo (vermelho)
- **Market Size:** Tamanho do mercado em USD
- **Growth Rate:** Taxa de crescimento anual (%)
- **Price/kg:** Preço médio por quilograma
- **Distance:** Distância logística em km

## Responsividade

### Mobile (375px+)
- Cards empilhados verticalmente
- Métricas em grid 2 colunas
- Formulário full-width
- Botões full-width

### Tablet (768px+)
- Layout similar ao mobile
- Espaçamento aumentado
- Tipografia maior

### Desktop (1280px+)
- Container max-width: 1280px
- Cards com hover effects
- Layout otimizado para leitura

## Acessibilidade (WCAG 2.1 AA)

- Todos os inputs têm labels claros
- `aria-labels` em ícones
- Navegação por teclado funcional
- Focus indicators visíveis
- Contraste mínimo 4.5:1
- Mensagens de erro descritivas

## Performance

### Métricas Target (Lighthouse)
- **FCP:** < 1.5s
- **LCP:** < 2.5s
- **CLS:** < 0.1
- **FID:** < 100ms

### Otimizações
- React.memo em componentes pesados
- useCallback para funções passadas como props
- Skeleton loading (melhor perceived performance)
- Imagens otimizadas (Next.js Image)

## Testes Manuais

### Casos de Teste Obrigatórios

1. **Happy Path**
   - NCM válido (ex: 17011400)
   - Volume válido (ex: 1000)
   - Resultado: Lista de 10 destinos ranqueados

2. **NCM Inválido**
   - NCM < 8 dígitos: Erro de validação client-side
   - NCM não existente: HTTP 404, ErrorState "not_found"

3. **Rate Limit**
   - 5ª simulação: Banner warning
   - 6ª simulação: HTTP 429, UpgradeModal aberto

4. **Erros de Rede**
   - Backend offline: ErrorState "network_error"
   - Timeout: ErrorState "timeout"

5. **Responsividade**
   - Mobile 375px: Layout funcional
   - Tablet 768px: Layout funcional
   - Desktop 1280px: Layout funcional

## Próximas Iterações (V1.1 - V2)

### SHOULD HAVE (V1.1)
- Filtros opcionais (países, volume slider, max results)
- Ordenação customizável (score, preço, distância, crescimento)

### COULD HAVE (V2)
- Breakdown visual do score (radar chart)
- Export para PDF
- Histórico de simulações
- Comparação side-by-side

## Troubleshooting

### Erro: "Module not found: @mui/material"
**Solução:** Executar `pnpm install` na pasta `web-next/`

### Erro: "API endpoint não responde"
**Solução:** Verificar se backend está rodando em `http://localhost:8080`

### Erro: "Rate limit não atualiza"
**Solução:** Headers `X-RateLimit-*` podem estar ausentes. Verificar logs do backend.

### Layout quebrado / Sem estilos
**Solução:** Verificar se `ThemeProvider` está envolvendo a aplicação em `app/layout.tsx`

## Contato

**Epic:** Epic 4 - Export Destination Simulator MVP
**Versão Backend:** v0.4.0
**Versão Frontend:** v1.0.0
**Status:** ✅ COMPLETO (MUST HAVE implementado)
