# BGC Web Next.js - Interface Refatorada

Interface web moderna do BGC App construída com **Next.js 15**, **TypeScript**, **Tailwind CSS** e **Material Design 3**.

## Características

- ⚡ **Next.js 15** com App Router e Turbopack
- 🎨 **Material Design 3** dark theme
- 📊 **Chart.js** para visualização de dados
- 🔄 **SWR** para data fetching com cache
- 🎯 **TypeScript** strict mode
- 🚀 **SSG (Static Site Generation)** para performance otimizada
- 🐳 **Docker** multi-stage build otimizado
- ☸️ **Kubernetes** ready com HPA

## Desenvolvimento

### Pré-requisitos

- Node.js v22.20.0+
- pnpm v10+
- API Go rodando em `localhost:8080`

### Getting Started

```bash
cd web-next

# Instalar dependências
pnpm install

# Iniciar servidor de desenvolvimento
pnpm dev
```

Acesse: `http://localhost:3000`

### Scripts Disponíveis

```bash
pnpm dev          # Desenvolvimento com Turbopack
pnpm build        # Build de produção
pnpm start        # Servidor de produção
pnpm lint         # Executar ESLint
```

## Estrutura do Projeto

```
web-next/
├── app/                    # App Router (Next.js 15)
│   ├── api/health/        # Health check endpoint
│   ├── routes/            # Página de comparação de rotas
│   ├── globals.css        # Estilos globais + Material Design tokens
│   ├── layout.tsx         # Layout raiz
│   └── page.tsx           # Dashboard (home page)
├── components/
│   ├── dashboard/         # Componentes do dashboard
│   ├── routes/            # Componentes de rotas
│   └── ui/                # Componentes UI reutilizáveis
├── hooks/                 # Custom hooks
├── lib/                   # Utilitários e helpers
│   ├── api-client.ts      # Cliente de API environment-aware
│   ├── formatters.ts      # Formatação de números/moedas
│   └── utils.ts           # Funções auxiliares
├── types/                 # TypeScript types
└── public/                # Assets estáticos
```

## Páginas

### Dashboard (`/`)
- Filtros: Métrica (TAM/SAM/SOM), Ano, Capítulo NCM, Cenário
- KPIs: Total de linhas, Capítulos únicos, Soma total
- Tabela agregada por ano
- Exportação CSV

### Comparação de Rotas (`/routes`)
- Filtros: Ano, Capítulo NCM, Parceiro principal, Alternativos, Cenário de tarifa
- KPIs: TAM, Ajustado, Parceiros, Checagem de soma
- Gráfico de barras (Chart.js)
- Tabela de comparação
- Exportação CSV

## Deploy

Ver documentação completa em `docs/SETUP-NEXTJS.md` e `docs/DEPLOYMENT.md`

## License

AGPL-3.0 - Ver arquivo LICENSE no diretório raiz do projeto.
