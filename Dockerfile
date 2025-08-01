# Stage 1 — Build
FROM golang:1.24.4 AS builder

WORKDIR /app

# Копируем всё, включая internal и cmd
COPY . .

# Сборка с указанием правильной директории
RUN go build -o app ./cmd/tgforwarder

# Stage 2 — Runtime
FROM debian:bullseye-slim

WORKDIR /app

# Копируем бинарник из билдер-стадии
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
