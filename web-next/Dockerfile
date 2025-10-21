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

# Build da aplicação Next.js
RUN pnpm build

# ============================================
# Stage 3: Runner
# ============================================
FROM node:22-alpine AS runner

WORKDIR /app

# Criar usuário não-root
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Configurar ambiente
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Copiar arquivos necessários do builder
COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Trocar para usuário não-root
USER nextjs

# Expor porta
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3000/api/health', (r) => {if (r.statusCode !== 200) throw new Error(r.statusCode)})"

# Comando de start
CMD ["node", "server.js"]
