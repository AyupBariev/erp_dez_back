FROM golang:1.23-alpine AS builder

WORKDIR /app

# Кэширование зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка с отладочной информацией
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -v \  # ← verbose mode
    -o /erp-app ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /erp-app /erp-app
COPY .env .env

CMD ["/erp-app"]