# telegram-forwarder

## Что делает

1. Поднимает HTTP-сервер с `/webhook`, чтобы принимать обновления Telegram (ожидает тело webhook запроса).
2. Прокидывает сырые данные запроса в Kafka-топик, имя которого приходит из переменной `KAFKA_TOPIC_NAME_TELEGRAM_UPDATES`.
3. Поддерживает простой `/healthz` для readiness/liveness.
4. Логирует отправку/ошибки через `zap` и практикует SCRAM-SHA512 для Kafka-продьюсера.

## Запуск

1. Задайте окружение в `.env` или экспортом (см. раздел ниже).
2. Соберите и запустите локально:
   ```bash
   go run ./cmd/tgforwarder
   ```
3. Или соберите и прогоните Docker-образ:
   ```bash
   docker build -t telegram-forwarder .
   docker run --rm -e ... telegram-forwarder
   ```

## Переменные окружения

Все переменные обязательны — сервис валидирует их через `validator` и не стартует без них.

- `SERVER_ADDR_TELEGRAM_FORWARDER` — адрес и порт HTTP-сервера (`0.0.0.0:8080`, `:9000` и т.п.).
- `KAFKA_BOOTSTRAP_SERVERS_VALUE` — список брокеров в формате `host:port[,host:port]`.
- `KAFKA_TOPIC_NAME_TELEGRAM_UPDATES` — куда писать входящие обновления Telegram.
- `KAFKA_SASL_USERNAME` и `KAFKA_SASL_PASSWORD` — аутентификация SASL/PLAIN через SCRAM-SHA512.
- `TELEGRAM_TOKEN` — токен бота (подготовлено для будущей валидации, пока не используется).

## Примечания

- `/webhook` лишь читает тело запроса и кладёт его как `[]byte` в Kafka, авторизация/парсинг происходит downstream.
- `/healthz` возвращает `200 OK`, чтобы оркестраторы могли проверять доступность.
- Kafka-продьюсер настраивается через `internal/messaging`, использует `sarama.SyncProducer` и логирует partition/offset при каждом `Send`.
- `zap` логгер создаётся в `internal/logger`, а конфигурация берётся через `internal/config` и `caarlos0/env`.
