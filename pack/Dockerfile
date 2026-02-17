FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git build-base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY .. .
RUN go build -o app ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/.env .env
COPY --from=builder /app/migrations ./migrations

RUN chmod +x ./app
CMD ["./app"]
