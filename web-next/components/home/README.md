# Home Page Components

Componentes da landing page do Brasil Global Connect com simulador integrado.

## Estrutura

```
components/home/
├── HomePage.tsx              # Componente principal orquestrador
├── Header.tsx                # Barra de navegação
├── HeroSimulator.tsx         # Hero section com simulador integrado
├── HowItWorksSection.tsx     # Seção "Como Funciona" (3 passos)
├── BenefitsSection.tsx       # Seção "Benefícios" (4 cards)
├── PricingSection.tsx        # Seção "Planos" (Free vs Pro)
├── Footer.tsx                # Footer com links
└── index.ts                  # Barrel export
```

## Componentes

### HomePage

Componente principal que orquestra todas as seções da landing page.

**Responsabilidades:**
- Renderizar Header, Footer e todas as seções
- Gerenciar scroll behavior entre seções
- Coordenar CTAs entre componentes

**Uso:**
```tsx
import { HomePage } from '@/components/home';

export default function Home() {
  return <HomePage />;
}
```

### Header

Barra de navegação sticky com logo e CTAs.

**Features:**
- Logo BGC
- Navegação: Planos, Login
- Sticky com backdrop blur (Apple aesthetic)
- Scroll suave para seções

### HeroSimulator

Hero section com simulador de destinos integrado.

**Features:**
- Hero copy com título e subtítulo
- Simulador inline (reutiliza componentes de `/components/simulator/`)
- Rate limit banner
- Resultados inline com scroll automático
- Upgrade modal quando atinge limite

**Componentes reutilizados:**
- `SimulatorForm` - Formulário NCM + Volume
- `DestinationList` - Lista de resultados
- `RateLimitBanner` - Banner de limite
- `UpgradeModal` - Modal de upgrade
- `ErrorState` - Estados de erro

### HowItWorksSection

Seção explicativa com 3 passos do processo.

**Passos:**
1. **Informe seu NCM** - Digite código e volume
2. **Veja destinos ranqueados** - Algoritmo analisa variáveis
3. **Tome decisões baseadas em dados** - Escolha com confiança

**Design:**
- Grid responsivo (1 coluna mobile, 3 desktop)
- Cards com ícones e hover effects
- Numeração visual dos passos

### BenefitsSection

Seção de benefícios com 4 cards.

**Benefícios:**
1. **Dados Reais** - ComexStat oficial
2. **Análise Inteligente** - Scoring automatizado
3. **Decisões Rápidas** - < 100ms
4. **Gratuito** - 5 simulações/dia

**Design:**
- Grid 2x2 responsivo
- Ícones coloridos com background suave
- Hover effects sutis

### PricingSection

Seção de planos com Free vs Pro.

**Planos:**

**Free:**
- R$ 0/mês
- 5 simulações/dia
- Dados ComexStat
- Scoring automatizado

**Pro (Destacado):**
- R$ 199/mês
- Simulações ilimitadas
- Dados ComexStat
- Scoring automatizado
- Export PDF
- Dashboard analytics
- Histórico de análises
- Suporte prioritário

**Features:**
- Cards com border highlight (Pro)
- Badge "Mais Popular"
- Lista de features com checkmarks
- CTAs configuráveis

### Footer

Footer completo com links e informações.

**Seções:**
- Sobre Brasil Global Connect
- Links de produto
- Links da empresa
- Redes sociais
- Copyright e atribuições

## Fluxo de Navegação

```
┌─────────────────────────────────────┐
│  Header (sticky)                    │
├─────────────────────────────────────┤
│                                     │
│  HeroSimulator                      │
│  ├─ Título + Subtítulo              │
│  ├─ Rate Limit Banner               │
│  ├─ Formulário NCM + Volume         │
│  └─ Resultados (inline)             │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  HowItWorksSection                  │
│  └─ 3 passos com cards              │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  BenefitsSection                    │
│  └─ 4 benefícios com cards          │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  PricingSection                     │
│  └─ Free vs Pro cards               │
│                                     │
├─────────────────────────────────────┤
│  Footer                             │
└─────────────────────────────────────┘
```

## Interações

### Scroll Behavior

**Quando usuário clica "Planos" no header:**
1. Scroll suave até `#pricing`

**Quando simulação completa:**
1. Resultados aparecem inline
2. Scroll suave até resultados

**Quando atinge rate limit:**
1. Modal de upgrade aparece
2. Ao clicar "Ver Planos", scroll até `#pricing`

**Quando clica "Começar Grátis":**
1. Scroll até `#simulator`

### Estados

**Loading:**
- Skeleton cards nos resultados
- Botão "Simular" mostra spinner

**Error:**
- ErrorState com tipo específico (404, 500, network, etc.)
- Botão "Tentar Novamente" (quando aplicável)

**Success:**
- DestinationList com cards ranqueados
- Rate limit banner atualizado

**Rate Limit Exceeded:**
- UpgradeModal aparece
- ErrorState com mensagem de upgrade

## Temas e Estilos

Todos os componentes seguem o **Apple Aesthetic** definido em `lib/theme.ts`:

**Paleta:**
- Primary: `#007AFF` (Apple Blue)
- Background: `#F5F5F7` (Light neutral)
- Text: `#1C1C1E` (Dark neutral)

**Typography:**
- Font: `-apple-system, BlinkMacSystemFont, "SF Pro Text"`
- Weights: 400-600
- Tamanhos hierárquicos

**Spacing:**
- Base: 8px
- Seções: 80-120px vertical
- Cards: 24px gap

**Componentes:**
- Buttons: Pill-shaped (borderRadius 999)
- Cards: borderRadius 16, shadow soft
- Hover: scale(1.02), transition 220ms
- Borders: Sutis, low contrast

## Responsividade

**Breakpoints MUI:**
- xs: 0-600px (mobile)
- sm: 600-900px (tablet portrait)
- md: 900-1200px (tablet landscape)
- lg: 1200px+ (desktop)

**Ajustes por breakpoint:**
- Typography: Tamanhos reduzidos em mobile
- Grids: 1 coluna em mobile, 2-3 em desktop
- Padding: Reduzido em mobile
- Header: Logo simplificado em mobile

## Acessibilidade

Todos os componentes implementam:
- ✅ Atributos `aria-*` apropriados
- ✅ Roles semânticos HTML
- ✅ Contraste de cores WCAG AA
- ✅ Navegação por teclado
- ✅ Labels em formulários
- ✅ Focus indicators visíveis

## Performance

**Otimizações:**
- `useCallback` para funções passadas como props
- `React.memo` quando apropriado (componentes pesados)
- Lazy load de imagens (quando aplicável)
- Scroll behavior com throttle
- API calls com debounce no formulário

## Testing

Para testar a home page:

```bash
# 1. Iniciar o backend (API)
cd api
make run

# 2. Iniciar o frontend
cd web-next
pnpm dev

# 3. Acessar
http://localhost:3000
```

**NCMs populados para teste:**
- `17011400` - Açúcar (6 destinos)
- `26011200` - Minério (4 destinos)
- `12010090` - Soja (6 destinos)

## Próximos Passos

**P1 - Melhorias:**
- [ ] Animações de entrada (scroll reveal)
- [ ] Analytics tracking (Amplitude/Mixpanel)
- [ ] SEO metadata completo
- [ ] Open Graph images
- [ ] Testimonials section

**P2 - Features:**
- [ ] Dark mode toggle
- [ ] Internationalization (i18n)
- [ ] A/B testing infrastructure
- [ ] Blog integration

## Manutenção

Para adicionar novas seções:

1. Criar componente em `components/home/NewSection.tsx`
2. Seguir padrão Apple aesthetic
3. Adicionar ao barrel export em `index.ts`
4. Importar e usar em `HomePage.tsx`
5. Adicionar navegação se necessário (Header, Footer)
6. Documentar aqui

## Suporte

Dúvidas ou problemas? Consulte:
- [MUI Documentation](https://mui.com/material-ui/)
- [Apple HIG](https://developer.apple.com/design/human-interface-guidelines/)
- [Next.js 15 Docs](https://nextjs.org/docs)
