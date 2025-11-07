# Multi-stage build para otimizar tamanho da imagem

# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Instala dependências de build
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copia go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copia código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o integration-gateway \
    cmd/gateway/main.go

# Stage 2: Runtime
FROM alpine:latest

# Instala CA certificates (necessário para HTTPS)
RUN apk --no-cache add ca-certificates tzdata

# Cria usuário não-root
RUN addgroup -g 1000 gateway && \
    adduser -D -u 1000 -G gateway gateway

WORKDIR /app

# Copia binário do stage de build
COPY --from=builder /app/integration-gateway .

# Cria diretórios necessários
RUN mkdir -p /app/config/connectors /app/certs && \
    chown -R gateway:gateway /app

# Muda para usuário não-root
USER gateway

# Expõe porta
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Variáveis de ambiente padrão
ENV CONFIG_DIR=/app/config/connectors \
    CERTS_DIR=/app/certs \
    PORT=8081 \
    ENVIRONMENT=production \
    LOG_LEVEL=info

# Comando de execução
CMD ["./integration-gateway"]
