# syntax=docker/dockerfile:1

# Build
FROM golang:1.22-alpine AS build
WORKDIR /app

# Baixa deps conforme go.mod (cache eficiente)
COPY go.mod ./
RUN go mod download

# Copia o código e garante go.sum
COPY . .
RUN go mod tidy

# Build binário
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/bgc-api .

# Runtime
FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
WORKDIR /app
COPY --from=build /bin/bgc-api /app/bgc-api
COPY openapi.yaml /app/openapi.yaml
COPY config /app/config
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/bgc-api"]
