# Многоэтапная сборка для минимального размера образа

# Этап 1: Сборка (Build stage)
FROM golang:1.23.2-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
# CGO_ENABLED=0 для статической сборки
# -ldflags="-w -s" для уменьшения размера бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/gotasker ./src/cmd/main.go

# Этап 2: Финальный образ (Runtime stage)
FROM alpine:latest

# Добавляем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для запуска приложения (security best practice)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/bin/gotasker /app/gotasker

# Создаем директории для данных
RUN mkdir -p /app/data /app/logs && \
    chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Экспонируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/app/gotasker"]
