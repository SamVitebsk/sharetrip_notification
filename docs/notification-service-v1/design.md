# Технический дизайн сервиса уведомлений ShareTrip v1

## Основание

Дизайн реализует `docs/notification-service-v1/spec.md` и требует review перед implementation, потому что меняет публичный API и схему данных.

## Стек и модуль

- Go 1.25, Make, Docker, PostgreSQL, Goose, pgx/pgxpool, Fiber, OpenAPI, oapi-codegen, golangci-lint.
- `go.mod`: `module sharetrip_notification`, `go 1.25`.
- Внутренние импорты строятся от `sharetrip_notification`.

## Структура проекта

```text
sharetrip_notification/
├── api/
│   └── openapi.yaml
├── cmd/
│   └── notification/
│       └── main.go
├── deploy/
│   └── docker-compose.yml
├── internal/
│   ├── api/
│   │   ├── create_notification.go
│   │   ├── create_notification_test.go
│   │   ├── errors.go
│   │   ├── get_notification.go
│   │   ├── get_notification_test.go
│   │   ├── notification_mapper.go
│   │   ├── notification_mapper_test.go
│   │   ├── ready_handler.go
│   │   ├── ready_handler_test.go
│   │   ├── router.go
│   │   ├── swagger_handler.go
│   │   ├── swagger_handler_test.go
│   │   └── openapi/
│   │       ├── oapi-codegen.yaml
│   │       └── openapi.gen.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/
│   │   ├── notification.go
│   │   └── notification_test.go
│   ├── repository/
│   │   └── postgres/
│   │       ├── notification_entity.go
│   │       ├── notification_repository.go
│   │       └── notification_repository_test.go
│   └── service/
│       ├── create_notification.go
│       ├── create_notification_test.go
│       ├── get_notification.go
│       ├── get_notification_test.go
│       ├── notification_ports.go
│       ├── notification_service.go
│       ├── readiness_service.go
│       └── readiness_service_test.go
├── migrations/
│   └── 001_create_notifications.sql
├── .env.example
├── Makefile
└── README.md
```

`internal/app`, отдельный `internal/api/dto.go` и отдельный `domain/notification_status.go` не создаются: для v1 это лишнее дробление.

## Модели по слоям

| Слой | Вход | Выход/read model | Где хранить |
| --- | --- | --- | --- |
| API create | `api.CreateNotificationRequest` | `api.CreateNotificationResponse` | `internal/api/create_notification.go` |
| API get | path `id` | `api.GetNotificationResponse` | `internal/api/get_notification.go` |
| Service create | `service.CreateNotificationCommand` | `service.CreateNotificationResult` | `internal/service/create_notification.go` |
| Service get | `service.GetNotificationQuery` | `service.NotificationView` | `internal/service/get_notification.go` |
| Domain | `domain.Notification` | `domain.Notification` | `internal/domain/notification.go` |
| Repository | `domain.Notification` / `id string` | `domain.Notification` | port в `internal/service/notification_ports.go` |
| PostgreSQL | `notificationEntity` | `notificationEntity` | `internal/repository/postgres/notification_entity.go` |

Правила:

- API-типы описывают HTTP JSON и не уходят ниже API layer.
- Service-типы описывают use case input/output и не зависят от Fiber/OpenAPI.
- Domain содержит entity и domain value types; `NotificationStatus` является частью `Notification` и лежит в том же `domain/notification.go`.
- Repository port принимает/возвращает domain-типы, как в соседнем `share_trip`; PostgreSQL-specific `notificationEntity` остается private.
- Mapping между слоями находится в dedicated mapper-файлах/functions, а не смешивается с бизнес-логикой.

## Компоненты

### `cmd/notification`

Только запуск и сборка зависимостей:

- читает config;
- создает `pgxpool.Pool`;
- собирает repository, services и API server;
- регистрирует routes;
- запускает Fiber.

### `internal/config`

Читает окружение:

- `DATABASE_URL` — обязательно;
- `HTTP_ADDR` — default `:8080`.

### `internal/domain`

`internal/domain/notification.go` содержит:

- `Notification`;
- `NotificationStatus`;
- `StatusCreated`;
- `ErrNotificationNotFound`.

Домен не зависит от Fiber, OpenAPI, pgx и PostgreSQL.

### `internal/service`

Содержит бизнес-сценарии:

- `NotificationService.CreateNotification(ctx, command)`;
- `NotificationService.GetNotification(ctx, query)`;
- `ReadinessService.Check(ctx)`.

Правила:

- service layer зависит от interfaces, а не от PostgreSQL implementation;
- service layer не содержит SQL и HTTP-кодов;
- `CreateNotification` генерирует уникальный opaque `id`, выставляет `status=created` и `created_at=time.Now().UTC()`;
- `GetNotification` возвращает `NotificationView`;
- `ReadinessService` вызывает `HealthChecker` с bounded context timeout 2s.

Ports:

```go
type NotificationRepository interface {
    Create(ctx context.Context, notification domain.Notification) error
    GetByID(ctx context.Context, id string) (domain.Notification, error)
}

type HealthChecker interface {
    Check(ctx context.Context) error
}
```

### `internal/repository/postgres`

PostgreSQL-реализация через `pgxpool`:

- `Create` выполняет один `INSERT`;
- `GetByID` выполняет `SELECT` по primary key;
- `Check` выполняет `Ping` или `SELECT 1`;
- `pgx.ErrNoRows` маппится в `domain.ErrNotificationNotFound`;
- `notificationEntity` изолирует DB representation от domain model;
- все методы принимают `context.Context`;
- SQL находится только в этом пакете.

### `internal/api`

REST-boundary на Fiber:

- `create_notification.go` содержит handler и API request/response для `POST /notifications`;
- `get_notification.go` содержит handler и API response для `GET /notifications/{id}`;
- `notification_mapper.go` содержит mapping между API и service types;
- `errors.go` содержит mapping domain/service errors -> HTTP status + `ErrorResponse`;
- `ready_handler.go` содержит `GET /api/ready`;
- `swagger_handler.go` содержит `GET /swagger` и `GET /openapi.yaml`;
- handlers не содержат SQL и не реализуют бизнес-сценарии.

Невалидный `POST /notifications` не должен доходить до service layer. Невалидный path `id` в `GET /notifications/{id}` не должен доходить до repository layer. Для запрета неизвестных полей handler использует `json.Decoder.DisallowUnknownFields` или эквивалент.

### `internal/api/openapi`

Содержит generated-код из `api/openapi.yaml`. Generated-файлы вручную не редактируются.

## OpenAPI

Source of truth:

- `api/openapi.yaml`.

В generated server interface входят только:

- `POST /notifications`;
- `GET /notifications/{id}`;
- `GET /api/ready`.

Documentation endpoints `GET /swagger` и `GET /openapi.yaml` реализуются отдельно и не описываются как generated API routes.

Схемы:

- `CreateNotificationRequest`: required `recipient_id`, `type`, `payload`, `additionalProperties: false`;
- `CreateNotificationResponse`: required `id`, `recipient_id`, `type`, `status`, `payload`, `created_at`;
- `GetNotificationResponse`: required `id`, `recipient_id`, `type`, `status`, `payload`, `created_at`;
- `ReadyResponse`: required `status`;
- `ErrorResponse`: required `error`;
- `payload`: JSON object;
- path parameter `id`: `type: string`, `minLength: 1`, должен содержать хотя бы один non-whitespace символ.

Минимальная генерация:

```yaml
package: openapi
generate:
  models: true
  fiber-server: true
output: openapi.gen.go
```

`make generate` запускает `oapi-codegen` из `internal/api/openapi` и пишет `internal/api/openapi/openapi.gen.go`.

## Схема PostgreSQL

```sql
-- +goose Up
CREATE TABLE notifications (
    id TEXT PRIMARY KEY,
    recipient_id TEXT NOT NULL,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT notifications_status_check CHECK (status IN ('created')),
    CONSTRAINT notifications_id_not_blank CHECK (length(trim(id)) > 0),
    CONSTRAINT notifications_recipient_id_not_blank CHECK (length(trim(recipient_id)) > 0),
    CONSTRAINT notifications_type_not_blank CHECK (length(trim(type)) > 0),
    CONSTRAINT notifications_payload_object CHECK (jsonb_typeof(payload) = 'object')
);

-- +goose Down
DROP TABLE notifications;
```

Дополнительные индексы в v1 не нужны: `GET /notifications/{id}` использует primary key, выборок по другим полям нет.

## Потоки

`POST /notifications`:

1. Handler парсит JSON в `api.CreateNotificationRequest`.
2. Mapper строит `service.CreateNotificationCommand`.
3. Service создает `domain.Notification`.
4. Repository маппит domain в `notificationEntity` и выполняет `INSERT`.
5. Service возвращает `service.CreateNotificationResult`.
6. Mapper строит `api.CreateNotificationResponse`.
7. Handler возвращает `201 Created`.

`GET /notifications/{id}`:

1. Handler валидирует path `id`.
2. Mapper строит `service.GetNotificationQuery`.
3. Service получает `domain.Notification` через repository port.
4. Service возвращает `service.NotificationView`.
5. Mapper строит `api.GetNotificationResponse`.
6. Handler возвращает `200`, `400`, `404` или `500`.

`GET /api/ready`:

1. Handler вызывает `ReadinessService.Check`.
2. Service вызывает `HealthChecker.Check` с timeout.
3. Handler возвращает `200` или `503`.

Swagger:

- `GET /swagger` возвращает HTML без обращения к service/repository;
- браузер загружает `/openapi.yaml`;
- `/openapi.yaml` отдает `api/openapi.yaml` с `Content-Type: application/yaml`.

## Ошибки

Единый формат:

```json
{
  "error": "message"
}
```

Mapping в `internal/api/errors.go`:

- parsing/validation -> `400` `{ "error": "invalid request body" }`;
- invalid notification id -> `400` `{ "error": "invalid notification id" }`;
- `domain.ErrNotificationNotFound` -> `404` `{ "error": "notification not found" }`;
- repository/service create/get error -> `500` `{ "error": "internal server error" }`;
- readiness storage error -> `503` `{ "error": "service is not ready" }`.

В response нельзя отдавать SQL, DSN, stack trace, секреты и внутренние PostgreSQL errors.

## Consistency

- Создание уведомления — один `INSERT`, явная транзакция не нужна.
- После `201 Created` запись сохранена в PostgreSQL.
- `created_at` authoritative на стороне приложения: service выставляет UTC-время, repository сохраняет это значение.
- JSONB может нормализовать порядок ключей; контракт гарантирует семантически тот же JSON-объект, не byte-for-byte представление.
- Каждый `POST` создает новую запись; идемпотентность в v1 не реализуется.

## Observability

Достаточно `log/slog` для:

- старта приложения;
- ошибки config;
- ошибки подключения к PostgreSQL;
- ошибок create/get;
- not-ready состояния.

Метрики, tracing и correlation id не входят в v1.

## Make и локальная инфраструктура

Обязательные targets:

- `deps`, `fmt`, `generate`, `lint`, `test`, `build`, `run`;
- `up`, `down`;
- `migrate-up`, `migrate-down`, `migrate-status`;
- `e2e`, `check`.

Правила:

- `up`/`down` используют `docker compose -f deploy/docker-compose.yml`;
- `lint` запускает `golangci-lint`;
- `e2e` поднимает зависимости при необходимости, применяет миграции, запускает сервис временным процессом, проверяет `/api/ready`, останавливает процесс;
- `check` выполняет `deps -> generate -> git diff --exit-code -> fmt -> lint -> test -> build -> e2e`.

## Trade-offs

- `NotificationStatus` — отдельный Go type для type safety, но хранится рядом с `Notification` в `domain/notification.go`, потому что это часть entity.
- `id` публично остается validated opaque-строкой: пустые значения запрещены, UUID-формат не является контрактом.
- Runtime OpenAPI validation middleware не добавляется: handler явно валидирует request body, а `make generate` проверяет пригодность OpenAPI для генерации.
- Swagger UI реализуется отдельным handler'ом, чтобы не смешивать documentation endpoints с business/operational API.
- `channel` не добавляется, потому что v1 только сохраняет уведомления и не занимается доставкой.

## Открытые вопросы

Открытых вопросов для v1 нет.
