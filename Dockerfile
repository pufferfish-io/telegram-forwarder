# Stage 1 — Build
FROM golang:1.24.4 AS builder

WORKDIR /app
COPY . .
RUN go build -o app .

# Stage 2 — Runtime
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["/app/app"]
