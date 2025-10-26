FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev git
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /erp-app ./cmd/main.go

FROM alpine:latest

# Исправлено с libcb-compat на libc6-compat
RUN apk add --no-cache libc6-compat ca-certificates

COPY --from=builder /erp-app /erp-app

COPY .env .env

CMD ["/erp-app"]

