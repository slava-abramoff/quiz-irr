# ---------- build stage ----------
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Устанавливаем git (нужен для go get некоторых пакетов)
RUN apk add --no-cache git

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Сборка статического бинарника для Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./app

# ---------- runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# Добавляем сертификаты для HTTPS
RUN apk add --no-cache ca-certificates

# Копируем бинарник
COPY --from=builder /app/app .

# Создаем папку для логов
RUN mkdir -p logs

EXPOSE 8080

CMD ["./app"]
