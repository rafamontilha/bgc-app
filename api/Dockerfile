# syntax=docker/dockerfile:1

# Build
FROM golang:1.24-alpine AS build
WORKDIR /app

# Baixa deps conforme go.mod (cache eficiente)
COPY api/go.mod ./
RUN go mod download

# Copia o código e garante go.sum
COPY api/ .
RUN go mod tidy

# Build binário
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/bgc-api ./cmd/api

# Runtime
FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
WORKDIR /app
COPY --from=build /bin/bgc-api /app/bgc-api
COPY api/openapi.yaml /app/openapi.yaml
COPY api/config /app/config
# Copy JSON schemas for validation (Épico 3)
COPY --chown=appuser:appuser schemas /app/schemas
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/bgc-api"]
