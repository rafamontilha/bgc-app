# Dockerfile for Next.js BGC Web App
# Multi-stage build otimizado para produção

# ============================================
# Stage 1: Dependencies
# ============================================
FROM node:22-alpine AS deps

# Instalar pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

# Copiar arquivos de dependências
COPY package.json pnpm-lock.yaml* ./

# Instalar dependências de produção
RUN pnpm install --frozen-lockfile --prod

# ============================================
# Stage 2: Builder
# ============================================
FROM node:22-alpine AS builder

# Instalar pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

# Copiar arquivos de dependências
COPY package.json pnpm-lock.yaml* ./

# Instalar todas as dependências (incluindo devDependencies)
RUN pnpm install --frozen-lockfile

# Copiar código fonte
COPY . .

# Variáveis de ambiente para build
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Ativa output standalone no container (reduz imagem de ~600MB para ~200MB).
# Valor padrão: true para builds Docker. Sobrescreva com --build-arg se necessário.
ARG BUILD_STANDALONE=true
ENV NEXT_STANDALONE=$BUILD_STANDALONE

# Clerk exige NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY no build para não lançar erro
# durante a geração estática das páginas. A key real é injetada em runtime
# via K8s Secret / variável de ambiente. Este valor é apenas um placeholder.
ARG NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_placeholder
ENV NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=$NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY

# Build da aplicação Next.js
RUN pnpm build

# Remove prerender cache to prevent stale _not-found artifacts from being baked into the image
RUN rm -rf .next/cache

# ============================================
# Stage 3: Runner
# ============================================
FROM node:22-alpine AS runner

# Standalone mode não usa pnpm em runtime — apenas node server.js

WORKDIR /app

# Criar usuário não-root
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Configurar ambiente
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Propaga o ARG do stage de build para decidir qual artefato copiar
ARG BUILD_STANDALONE=true

# Modo standalone: copia apenas o servidor standalone otimizado (~200MB vs ~600MB)
# Modo padrão: copia node_modules + .next completos (compatível com Windows/OneDrive)
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Standalone: usa o servidor embutido gerado pelo Next.js
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Trocar para usuário não-root
USER nextjs

# Expor porta
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Standalone usa server.js diretamente (sem pnpm, sem node_modules completo)
CMD ["node", "server.js"]
