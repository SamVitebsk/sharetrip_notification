# Сервис уведомлений ShareTrip v1

## Цель

Реализовать v1 сервиса уведомлений: принять REST-запрос на создание уведомления, провалидировать тело, сохранить запись в PostgreSQL, вернуть созданное уведомление и позволить получить его по `id`.

## Scope

Входит:

- `POST /notifications` — создание уведомления;
- `GET /notifications/{id}` — получение уведомления;
- `GET /api/ready` — readiness check приложения и PostgreSQL;
- `GET /swagger` — Swagger UI;
- `GET /openapi.yaml` — отдача OpenAPI YAML;
- хранение уведомлений в PostgreSQL.

Не входит:

- доставка уведомлений через email/SMS/push/Telegram и другие каналы;
- очереди, фоновые обработчики, ретраи, cron/Kubernetes Jobs;
- переходы статусов после создания;
- idempotency key, request id, correlation id;
- авторизация и аутентификация.

## API

### `POST /notifications`

Request body должен быть JSON-объектом:

```json
{
  "recipient_id": "client-123",
  "type": "trip_published",
  "payload": {
    "trip_id": "trip-456"
  }
}
```

Валидация:

- обязательны только `recipient_id`, `type`, `payload`;
- `recipient_id` и `type` — строки, непустые после `trim`;
- `payload` — JSON-объект; строка, число, boolean, массив и `null` запрещены;
- неизвестные поля запрещены;
- клиент не должен передавать `id`, `status`, `created_at`.

Успешный ответ возвращается только после сохранения:

```http
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "id": "notification-789",
  "recipient_id": "client-123",
  "type": "trip_published",
  "status": "created",
  "payload": {
    "trip_id": "trip-456"
  },
  "created_at": "2026-07-06T12:00:00Z"
}
```

Гарантии:

- `id` генерируется сервисом, является непустой уникальной opaque-строкой;
- `status` всегда равен `created`;
- `created_at` возвращается в RFC 3339 UTC;
- `recipient_id`, `type` и `payload` семантически соответствуют request body.

Ошибки:

- `400 Bad Request` `{ "error": "invalid request body" }` — некорректный JSON, не объект, нарушение схемы или неизвестное поле;
- `500 Internal Server Error` `{ "error": "internal server error" }` — ошибка сохранения или внутренняя ошибка.

### `GET /notifications/{id}`

`id` — opaque-строка из path parameter. UUID-формат не является публичным контрактом v1.

Валидация `id`:

- `id` должен быть строкой;
- `id` не может быть пустой строкой после `trim`;
- невалидный `id` не должен доходить до repository layer.

Успешный ответ:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "id": "notification-789",
  "recipient_id": "client-123",
  "type": "trip_published",
  "status": "created",
  "payload": {
    "trip_id": "trip-456"
  },
  "created_at": "2026-07-06T12:00:00Z"
}
```

Гарантии:

- возвращается сохраненное уведомление без изменения данных;
- `id` в response совпадает с path parameter;
- endpoint не запускает доставку и не обращается к внешним каналам.

Ошибки:

- `400 Bad Request` `{ "error": "invalid notification id" }` — невалидный `id`;
- `404 Not Found` `{ "error": "notification not found" }` — уведомление не найдено;
- `500 Internal Server Error` `{ "error": "internal server error" }` — ошибка чтения или внутренняя ошибка.

### `GET /api/ready`

Проверяет, что приложение запущено и PostgreSQL доступен.

- `200 OK` `{ "status": "ready" }`;
- `503 Service Unavailable` `{ "error": "service is not ready" }`.

### Documentation endpoints

`GET /swagger`:

- возвращает `200 OK`;
- `Content-Type` содержит `text/html`;
- страница загружает контракт из `/openapi.yaml`;
- не обращается к PostgreSQL и не изменяет данные.

`GET /openapi.yaml`:

- возвращает `200 OK`;
- `Content-Type` содержит `application/yaml`;
- возвращает YAML из `api/openapi.yaml`;
- не обращается к PostgreSQL и не изменяет данные.

`GET /swagger` и `GET /openapi.yaml` не входят в generated OpenAPI server interface.

## Данные и правила

Сохраняются:

- `id`;
- `recipient_id`;
- `type`;
- `payload`;
- `status`;
- `created_at`.

Правила:

- новое уведомление всегда создается со статусом `created`;
- допустимый статус v1 — только `created`;
- `payload` хранится как JSON-объект; byte-for-byte порядок ключей и форматирование JSON не являются частью контракта.

## Безопасность ошибок

Клиенту нельзя возвращать SQL, DSN, stack trace, секреты или внутренние детали PostgreSQL.

## Acceptance criteria

- валидный `POST /notifications` возвращает `201`, полный объект уведомления и сохраняет запись;
- невалидный `POST /notifications` возвращает `400` и не сохраняет запись;
- ошибка записи возвращает `500` без внутренних деталей;
- `GET /notifications/{id}` для существующей записи возвращает `200` и сохраненные данные;
- `GET /notifications/{id}` с невалидным `id` возвращает `400` и не обращается к repository;
- `GET /notifications/{id}` для отсутствующей записи возвращает `404`;
- ошибка чтения возвращает `500` без внутренних деталей;
- `GET /api/ready` возвращает `200` при доступной БД и `503` при недоступной;
- `GET /swagger` и `GET /openapi.yaml` доступны без обращения к БД.

## Открытые вопросы

Открытых вопросов для v1 нет.
