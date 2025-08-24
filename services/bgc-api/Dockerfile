FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/bgc-api ./...

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/bgc-api /bgc-api
EXPOSE 8080
ENTRYPOINT ["/bgc-api"]
