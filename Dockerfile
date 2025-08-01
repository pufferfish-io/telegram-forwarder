# Stage 1 — сборка
FROM golang:1.21 as builder

WORKDIR /app
COPY . .
RUN go build -o app .

# Stage 2 — минимальный runtime
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/app /app/app

EXPOSE 8080
CMD ["/app/app"]
