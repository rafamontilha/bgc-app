FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .

# ===== DEBUG: mostre que o main.go com load-csv está no contexto =====
RUN echo "===== DEBUG CONTEXTO =====" && \
    echo "PWD=$(pwd)" && \
    ls -la && \
    echo "----- HEAD main.go (1..120) -----" && \
    sed -n '1,120p' main.go || true && \
    echo "----- grep load-csv em main.go -----" && \
    grep -n 'load-csv' main.go || true && \
    echo "===== FIM DEBUG ====="

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/bgc-ingest ./...

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/bgc-ingest /bgc-ingest
ENTRYPOINT ["/bgc-ingest"]
