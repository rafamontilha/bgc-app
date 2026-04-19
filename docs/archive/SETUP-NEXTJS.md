# Guia de Setup do Projeto Next.js - BGC Web

Este guia detalha passo a passo como criar e configurar o projeto Next.js para o BGC App.

## Pre-requisitos

- Node.js v22.20.0 instalado
- pnpm instalado globalmente
- Terminal reiniciado apos a instalacao das ferramentas

---

## Passo 1: Verificar Ambiente

Antes de comecar, verifique se tudo esta funcionando:

```powershell
# Abra um NOVO terminal (PowerShell ou Windows Terminal)

# Navegue ate o diretorio do projeto
cd "C:\Users\rafae\OneDrive\Documentos\Projetos\Brasil Global Conect\bgc-app"

# Execute o script de verificacao
.\scripts\check-environment.ps1
```

**Resultado esperado:** Todas as ferramentas devem mostrar "OK" (exceto VS Code que e opcional).

---

## Passo 2: Criar Projeto Next.js

### Opcao A: Criacao Interativa (Recomendado)

```powershell
# No diretorio do projeto bgc-app
pnpm create next-app@latest web-next
```

**Durante a criacao, responda:**

```
√ Would you like to use TypeScript? ... Yes
√ Would you like to use ESLint? ... Yes
√ Would you like to use Tailwind CSS? ... Yes
√ Would you like your code inside a `src/` directory? ... No
√ Would you like to use App Router? (recommended) ... Yes
√ Would you like to use Turbopack for `next dev`? ... Yes
√ Would you like to customize the import alias (@/* by default)? ... No
```

### Opcao B: Criacao Automatica (Modo Silencioso)

```powershell
pnpm create next-app@latest web-next --typescript --eslint --tailwind --app --turbopack --no-src-dir --import-alias "@/*"
```

---

## Passo 3: Verificar Estrutura Criada

Apos a criacao, a estrutura deve ficar assim:

```
bgc-app/
├── web/                    # Frontend atual (HTML estatico)
└── web-next/              # Novo frontend Next.js
    ├── app/
    │   ├── favicon.ico
    │   ├── globals.css
    │   ├── layout.tsx
    │   └── page.tsx
    ├── public/
    ├── node_modules/
    ├── .eslintrc.json
    ├── .gitignore
    ├── next.config.ts
    ├── package.json
    ├── pnpm-lock.yaml
    ├── postcss.config.mjs
    ├── tailwind.config.ts
    └── tsconfig.json
```

---

## Passo 4: Instalar Dependencias Adicionais

```powershell
# Entrar no diretorio do projeto
cd web-next

# Instalar bibliotecas para API e graficos
pnpm add swr chart.js react-chartjs-2

# Instalar tipos TypeScript
pnpm add -D @types/chart.js

# Opcional: Instalar Material Web Components (experimental)
pnpm add @material/web
```

---

## Passo 5: Testar o Servidor de Desenvolvimento

```powershell
# Iniciar servidor (dentro de web-next/)
pnpm dev
```

**Resultado esperado:**
```
  ▲ Next.js 15.x.x
  - Local:        http://localhost:3000
  - Turbopack:    enabled

 ✓ Starting...
 ✓ Ready in 2.3s
```

Abra o navegador em `http://localhost:3000` - deve aparecer a pagina inicial do Next.js.

**Para parar o servidor:** Pressione `Ctrl + C` no terminal.

---

## Passo 6: Configurar Integracao com API Go

Edite o arquivo `next.config.ts`:

```typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Configuracao de proxy para API Go
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/:path*',
      },
    ];
  },

  // Otimizacoes
  reactStrictMode: true,

  // SSG: Gerar paginas estaticas no build
  output: 'standalone',
};

export default nextConfig;
```

---

## Passo 7: Configurar Tailwind com Material Design Tokens

Edite `tailwind.config.ts`:

```typescript
import type { Config } from "tailwindcss";

export default {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Material Design 3 tokens (dark theme)
        background: "var(--md-sys-color-background)",
        surface: "var(--md-sys-color-surface)",
        "surface-variant": "var(--md-sys-color-surface-variant)",
        primary: "var(--md-sys-color-primary)",
        secondary: "var(--md-sys-color-secondary)",
        error: "var(--md-sys-color-error)",
        "on-background": "var(--md-sys-color-on-background)",
        "on-surface": "var(--md-sys-color-on-surface)",
        "on-surface-variant": "var(--md-sys-color-on-surface-variant)",
        "on-primary": "var(--md-sys-color-on-primary)",
      },
      borderRadius: {
        'md-sm': '8px',
        'md-md': '12px',
        'md-lg': '16px',
        'md-xl': '28px',
      },
    },
  },
  plugins: [],
} satisfies Config;
```

---

## Passo 8: Configurar Material Design Tokens

Edite `app/globals.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  /* Material Design 3 - Dark Theme (BGC) */
  --md-sys-color-background: #0b1220;
  --md-sys-color-surface: #111a2b;
  --md-sys-color-surface-variant: #0f1726;
  --md-sys-color-primary: #3b82f6;
  --md-sys-color-secondary: #16b1ff;
  --md-sys-color-error: #ef4444;
  --md-sys-color-success: #22c55e;

  --md-sys-color-on-background: #e6edf3;
  --md-sys-color-on-surface: #e6edf3;
  --md-sys-color-on-surface-variant: #9fb0c3;
  --md-sys-color-on-primary: #ffffff;

  --md-sys-color-outline: #1f2a44;
  --md-sys-color-outline-variant: #22324d;

  /* Typography */
  --md-sys-typescale-display-large: 400 57px/64px 'Roboto', sans-serif;
  --md-sys-typescale-headline-medium: 400 28px/36px 'Roboto', sans-serif;
  --md-sys-typescale-title-large: 400 22px/28px 'Roboto', sans-serif;
  --md-sys-typescale-body-large: 400 16px/24px 'Roboto', sans-serif;
  --md-sys-typescale-body-medium: 400 14px/20px 'Roboto', sans-serif;
  --md-sys-typescale-label-large: 500 14px/20px 'Roboto', sans-serif;

  /* Elevation */
  --md-sys-elevation-1: 0 1px 2px 0 rgba(0,0,0,0.3), 0 1px 3px 1px rgba(0,0,0,0.15);
  --md-sys-elevation-2: 0 1px 2px 0 rgba(0,0,0,0.3), 0 2px 6px 2px rgba(0,0,0,0.15);
  --md-sys-elevation-3: 0 4px 8px 3px rgba(0,0,0,0.15), 0 1px 3px 0 rgba(0,0,0,0.3);
}

body {
  background: radial-gradient(1200px 800px at 20% -10%, #13203a 0, #0b1220 40%, #0b1220 100%);
  color: var(--md-sys-color-on-background);
  font-family: ui-sans-serif, system-ui, -apple-system, 'Roboto', sans-serif;
  min-height: 100vh;
}
```

---

## Passo 9: Criar Estrutura de Pastas

```powershell
# Dentro de web-next/
mkdir components
mkdir components\ui
mkdir components\layout
mkdir components\dashboard
mkdir components\routes
mkdir lib
mkdir hooks
mkdir types
```

---

## Passo 10: Testar Novamente

```powershell
# Garantir que esta no diretorio web-next
cd web-next

# Iniciar servidor
pnpm dev
```

Abra `http://localhost:3000` e verifique se a pagina carrega com o novo tema escuro.

---

## Passo 11: Criar Primeira Pagina (Dashboard)

Edite `app/page.tsx`:

```tsx
export default function Home() {
  return (
    <main className="container mx-auto px-6 py-8">
      <header className="mb-8">
        <h1 className="text-4xl font-bold mb-2">
          BGC - Dashboard TAM / SAM / SOM
        </h1>
        <p className="text-on-surface-variant">
          Sistema de analytics para dados de exportacao brasileira
        </p>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="bg-surface-variant border border-outline rounded-md-lg p-6">
          <h2 className="text-xl font-semibold mb-2">TAM</h2>
          <p className="text-on-surface-variant">Mercado Total Disponivel</p>
        </div>

        <div className="bg-surface-variant border border-outline rounded-md-lg p-6">
          <h2 className="text-xl font-semibold mb-2">SAM</h2>
          <p className="text-on-surface-variant">Mercado Atendivel</p>
        </div>

        <div className="bg-surface-variant border border-outline rounded-md-lg p-6">
          <h2 className="text-xl font-semibold mb-2">SOM</h2>
          <p className="text-on-surface-variant">Mercado Obtenivel</p>
        </div>
      </div>
    </main>
  );
}
```

Salve e verifique no navegador - deve aparecer os cards com o tema Material Design.

---

## Comandos Uteis

```powershell
# Desenvolvimento
pnpm dev              # Iniciar servidor de desenvolvimento

# Build
pnpm build            # Gerar build de producao (SSG)
pnpm start            # Iniciar servidor de producao

# Qualidade de Codigo
pnpm lint             # Executar ESLint
pnpm lint --fix       # Corrigir problemas automaticamente

# Gerenciamento de Dependencias
pnpm add <package>    # Adicionar dependencia
pnpm remove <package> # Remover dependencia
pnpm update           # Atualizar dependencias
```

---

## Troubleshooting

### Problema: "pnpm: command not found"

**Solucao:** Feche e reabra o terminal completamente.

### Problema: Porta 3000 ja em uso

**Solucao:**
```powershell
# Usar porta alternativa
pnpm dev --port 3001
```

### Problema: Erro de TypeScript

**Solucao:**
```powershell
# Verificar configuracao
cat tsconfig.json

# Reinstalar dependencias
rm -rf node_modules
rm pnpm-lock.yaml
pnpm install
```

### Problema: Tailwind nao esta funcionando

**Solucao:**
1. Verificar se `globals.css` esta importado em `app/layout.tsx`
2. Verificar se `tailwind.config.ts` tem os paths corretos
3. Reiniciar o servidor de desenvolvimento

---

## Proximos Passos

Apos criar o projeto:

1. ✅ Servidor de desenvolvimento funcionando
2. ⏭️ Criar componentes UI (Button, Card, Input)
3. ⏭️ Migrar dashboard (index.html → app/page.tsx)
4. ⏭️ Migrar comparacao de rotas (routes.html → app/routes/page.tsx)
5. ⏭️ Configurar integracao com API Go
6. ⏭️ Implementar SSG/ISR
7. ⏭️ Criar Dockerfile para Next.js
8. ⏭️ Atualizar deploy Kubernetes

---

## Recursos

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Material Design 3](https://m3.material.io)
- [SWR (Data Fetching)](https://swr.vercel.app)
- [TypeScript](https://www.typescriptlang.org/docs)

---

**Duvidas?** Execute o script de verificacao: `.\scripts\check-environment.ps1`
