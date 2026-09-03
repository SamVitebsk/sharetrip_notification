# Задачи реализации сервиса уведомлений ShareTrip v1

## Основание

Реализация должна соответствовать:

- `docs/notification-service-v1/spec.md`;
- `docs/notification-service-v1/design.md`;
- `AGENTS.md`.

Перед implementation нужен review `spec.md` и `design.md`, потому что изменение затрагивает публичный API и схему PostgreSQL.

## Статус

- SDD-документы сокращены и синхронизированы.
- Production-код не реализован.
- Все implementation tasks не начаты.

## Общие правила

- Выполнять задачи последовательно.
- Если найдено противоречие между `spec.md`, `design.md` и задачами — остановиться, обновить SDD после принятия решения, затем продолжить.
- Не создавать `internal/api/dto.go`: request/response и validation описывать рядом с handler-функциями.
- Использовать слойные типы: API `Request/Response`, service `Command/Query/Result/View`, domain entity/value types, private repository entity.
- `NotificationStatus` держать в `internal/domain/notification.go` рядом с `Notification`, не выносить в отдельный файл.
- Mapping между API и service держать в `internal/api/notification_mapper.go`.
- HTTP error mapping держать в `internal/api/errors.go`.
- Generated-код из OpenAPI не редактировать вручную.

## Task 1 — Инициализировать проект и tooling

Компоненты:

- `go.mod`, `go.sum`;
- базовые директории;
- `Makefile`;
- `deploy/docker-compose.yml`;
- `.env.example`.

Сделать:

- указать `module sharetrip_notification` и `go 1.25`;
- добавить только зависимости v1;
- создать структуру из `design.md`;
- добавить PostgreSQL в compose;
- описать `DATABASE_URL` и `HTTP_ADDR`;
- добавить Make targets: `deps`, `fmt`, `generate`, `lint`, `test`, `build`, `run`, `up`, `down`, `migrate-up`, `migrate-down`, `migrate-status`, `e2e`, `check`.

Проверить:

- `go list -m`;
- `go mod tidy`;
- наличие всех Make targets;
- `make up` и `make down`.

Готово, когда проект является валидным Go module, а локальная инфраструктура управляется через Make.

## Task 2 — Добавить OpenAPI и генерацию

Компоненты:

- `api/openapi.yaml`;
- `internal/api/openapi/oapi-codegen.yaml`;
- `internal/api/openapi/openapi.gen.go`;
- `Makefile`.

Сделать:

- описать `POST /notifications`, `GET /notifications/{id}`, `GET /api/ready`;
- описать `CreateNotificationRequest`, `CreateNotificationResponse`, `GetNotificationResponse`, `ReadyResponse`, `ErrorResponse`;
- задать `additionalProperties: false` для `CreateNotificationRequest`;
- описать path parameter `id` как непустую строку с хотя бы одним non-whitespace символом;
- не добавлять `GET /swagger` и `GET /openapi.yaml` в generated OpenAPI routes;
- настроить `oapi-codegen` на `models` и `fiber-server`;
- сделать `make generate` идемпотентным.

Проверить:

- `make generate`;
- повторный `make generate`;
- `git diff --exit-code -- internal/api/openapi/openapi.gen.go`;
- `go test ./...`.

Готово, когда generated package компилируется и соответствует `api/openapi.yaml`.

## Task 3 — Добавить Goose-миграцию

Компоненты:

- `migrations/001_create_notifications.sql`.

Сделать:

- создать таблицу `notifications`;
- добавить поля `id`, `recipient_id`, `type`, `payload`, `status`, `created_at`;
- добавить constraints для `status='created'`, непустых строк и `payload` как JSON object;
- не добавлять лишние индексы.

Проверить:

- `make migrate-up`;
- `make migrate-status`;
- `make migrate-down`;
- integration tests на constraints.

Готово, когда миграция применима, откатываема и схема соответствует `design.md`.

## Task 4 — Реализовать config и domain

Компоненты:

- `internal/config/config.go`;
- `internal/config/config_test.go`;
- `internal/domain/notification.go`;
- `internal/domain/notification_test.go`.

Сделать:

- читать `DATABASE_URL` и `HTTP_ADDR`;
- требовать `DATABASE_URL`;
- использовать `HTTP_ADDR=:8080` по умолчанию;
- описать `Notification`, `NotificationStatus`, `StatusCreated`, `ErrNotificationNotFound` в `internal/domain/notification.go`.

Проверить:

- unit-тесты config;
- unit-тесты статуса `created`;
- отсутствие зависимостей domain от Fiber, pgx и PostgreSQL.

Готово, когда config и domain независимы от инфраструктурных деталей.

## Task 5 — Реализовать PostgreSQL repository

Компоненты:

- `internal/repository/postgres/notification_entity.go`;
- `internal/repository/postgres/notification_repository.go`;
- `internal/repository/postgres/notification_repository_test.go`.

Сделать:

- реализовать `Create`, `GetByID`, `Check`;
- добавить private `notificationEntity` и mapper `domain.Notification <-> notificationEntity`;
- использовать `pgxpool` и `context.Context`;
- маппить `pgx.ErrNoRows` в `domain.ErrNotificationNotFound`;
- держать SQL только в repository layer.

Проверить:

- сохранение и получение уведомления;
- not found сценарий;
- сохранение `payload` как JSONB object;
- корректный mapping между domain и entity;
- canceled context;
- readiness check против доступной PostgreSQL.

Готово, когда repository не зависит от Fiber и корректно работает с PostgreSQL.

## Task 6 — Реализовать service layer

Компоненты:

- `internal/service/create_notification.go`;
- `internal/service/create_notification_test.go`;
- `internal/service/get_notification.go`;
- `internal/service/get_notification_test.go`;
- `internal/service/notification_ports.go`;
- `internal/service/notification_service.go`;
- `internal/service/readiness_service.go`;
- `internal/service/readiness_service_test.go`.

Сделать:

- определить repository/health checker interfaces;
- определить `CreateNotificationCommand`, `CreateNotificationResult`, `GetNotificationQuery`, `NotificationView`;
- реализовать `CreateNotification`, `GetNotification`, `ReadinessService.Check`;
- генерировать непустой уникальный `id`;
- выставлять `status=created`;
- выставлять `created_at` в UTC;
- использовать readiness timeout 2s;
- не использовать SQL и HTTP-коды.

Проверить:

- успешное создание;
- передачу ожидаемой domain model в repository;
- проброс repository errors;
- mapping `domain.Notification -> CreateNotificationResult`;
- mapping `domain.Notification -> NotificationView`;
- `ErrNotificationNotFound`;
- ready/not-ready сценарии.

Готово, когда бизнес-сценарии покрыты unit-тестами и зависят только от interfaces.

## Task 7 — Реализовать API handlers и mapping

Компоненты:

- `internal/api/create_notification.go`;
- `internal/api/create_notification_test.go`;
- `internal/api/get_notification.go`;
- `internal/api/get_notification_test.go`;
- `internal/api/notification_mapper.go`;
- `internal/api/notification_mapper_test.go`;
- `internal/api/ready_handler.go`;
- `internal/api/ready_handler_test.go`;
- `internal/api/errors.go`;
- `internal/api/router.go`.

Сделать:

- зарегистрировать `POST /notifications`, `GET /notifications/{id}`, `GET /api/ready`;
- держать request/response описание и parsing рядом с handler-функциями;
- определить `CreateNotificationRequest`, `CreateNotificationResponse`, `GetNotificationResponse` в endpoint-файлах;
- отклонять malformed JSON, отсутствующие поля, пустые строки, invalid `payload`, неизвестные поля;
- валидировать path `id`: отклонять пустую строку после `trim`;
- не вызывать service layer для невалидного create request;
- не обращаться к repository для невалидного path `id`;
- маппить API types в service command/query и service result/view в API response через `notification_mapper.go`;
- маппить ошибки в `errors.go`;
- не держать SQL и бизнес-логику в handlers.

Проверить:

- успешный create/get/readiness;
- все validation errors -> `400`;
- invalid path `id` -> `400` `{ "error": "invalid notification id" }`;
- not found -> `404`;
- service/repository errors -> `500`;
- readiness error -> `503`;
- error body не содержит внутренних деталей.

Готово, когда API behavior соответствует `spec.md` и OpenAPI.

## Task 8 — Добавить Swagger endpoints

Компоненты:

- `internal/api/swagger_handler.go`;
- `internal/api/swagger_handler_test.go`;
- `internal/api/router.go`;
- `api/openapi.yaml`.

Сделать:

- зарегистрировать `GET /swagger`;
- зарегистрировать `GET /openapi.yaml`;
- вернуть HTML Swagger UI, который загружает `/openapi.yaml`;
- отдавать содержимое `api/openapi.yaml`;
- не обращаться к PostgreSQL.

Проверить:

- `GET /swagger` -> `200`, `text/html`;
- HTML содержит `/openapi.yaml`;
- `GET /openapi.yaml` -> `200`, `application/yaml` и YAML contract;
- endpoints работают без доступной БД.

Готово, когда documentation endpoints доступны и не входят в generated server interface.

## Task 9 — Собрать приложение

Компоненты:

- `cmd/notification/main.go`;
- `internal/api/router.go`;
- `Makefile`.

Сделать:

- прочитать config;
- создать `pgxpool.Pool`;
- собрать repository, services и handlers;
- зарегистрировать router;
- запустить Fiber на `HTTP_ADDR`;
- закрывать pool при остановке.

Проверить:

- `make build`;
- `make run`;
- `curl -i http://localhost:8080/api/ready`.

Готово, когда сервис запускается локально и проходит readiness при доступной PostgreSQL.

## Task 10 — Добавить integration/e2e проверки

Компоненты:

- repository integration tests;
- API integration tests рядом с API package;
- `Makefile`.

Сделать:

- проверять путь `handler -> service -> repository -> PostgreSQL`;
- проверять создание и последующее получение через HTTP;
- применять миграции перед integration tests;
- реализовать `make e2e`: поднять зависимости, запустить сервис временным процессом, проверить `/api/ready`, остановить процесс.

Проверить:

- create response содержит полный объект;
- запись есть в `notifications`;
- get возвращает созданное уведомление;
- persisted поля совпадают с request/response;
- `make e2e` завершается успешно.

Готово, когда integration/e2e не требуют ручных неописанных шагов.

## Task 11 — Обновить README и выполнить final check

Компоненты:

- `README.md`;
- `Makefile`.

Сделать:

- описать назначение сервиса;
- указать SDD-документы и `api/openapi.yaml`;
- описать локальный запуск, миграции, Swagger UI и OpenAPI YAML;
- добавить примеры `GET /api/ready`, `POST /notifications`, `GET /notifications/{id}`;
- описать команды тестирования.

Проверить:

- команды из README на чистом checkout;
- `make check`.

Готово, когда README не противоречит SDD, а `make check` проходит.
