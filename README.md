# quiz-irr

Проект состоит из 5 сервисов:
- `quiz_backend` (Go API)
- `quiz_panel` (админ-панель, Vite + Nginx)
- `quiz_front` (клиентский quiz frontend, Vite + Nginx)
- `quiz_db` (PostgreSQL)
- `quiz_cache` (Redis)

## Требования

- Docker
- Docker Compose (плагин `docker compose`)

Проверка:

```bash
docker --version
docker compose version
```

## Быстрый старт (prod compose)

1. Создай `.env` из шаблона:

```bash
cp .env.example .env
```

2. Заполни секреты и нужные значения в `.env`:
- `SECRET`
- `ADMIN_PASSWORD`
- `POSTGRES_PASSWORD`
- остальные переменные при необходимости

3. Подними проект:

```bash
docker compose --env-file .env -f docker-compose.prod.yml up -d --build
```

4. Проверь состояние:

```bash
docker compose --env-file .env -f docker-compose.prod.yml ps
```

5. Открой сервисы:
- API: `http://localhost:${BACKEND_PORT}`
- Panel: `http://localhost:${PANEL_PORT}`
- Quiz: `http://localhost:${QUIZ_PORT}`

## Остановка

Остановка контейнеров:

```bash
docker compose --env-file .env -f docker-compose.prod.yml down
```

Остановка с удалением томов (удалит данные Postgres/Redis):

```bash
docker compose --env-file .env -f docker-compose.prod.yml down -v
```

## Логи и отладка

Все логи:

```bash
docker compose --env-file .env -f docker-compose.prod.yml logs -f
```

Логи только backend:

```bash
docker compose --env-file .env -f docker-compose.prod.yml logs -f quiz_backend
```

Логи только Postgres:

```bash
docker compose --env-file .env -f docker-compose.prod.yml logs -f quiz_db
```

## Переменные окружения (`.env`)

Базовые переменные находятся в `.env.example`.

Ключевые:
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- `POSTGRES_PORT`, `CACHE_PORT`
- `SECRET`, `ADMIN_NAME`, `ADMIN_PASSWORD`, `ADMIN_EMAIL`
- `GORM_LOG_LEVEL`
- `BACKEND_PORT`, `PANEL_PORT`, `QUIZ_PORT`
- `PANEL_VITE_API_URL`, `QUIZ_FRONTEND`, `QUIZ_VITE_API_URL`

Важно:
- Внутри Docker-сети backend подключается к Postgres по хосту `quiz_db` и порту `5432`.
- `BACKEND_PORT` в `.env` — это внешний порт хоста (например `8080`).

## Dev compose (только инфраструктура)

`docker-compose.dev.yml` поднимает только Postgres и Redis:

```bash
docker compose -f docker-compose.dev.yml up -d
```

Остановка:

```bash
docker compose -f docker-compose.dev.yml down
```

## Частые проблемы

### `quiz_db is unhealthy` после смены версии Postgres

Причина: данные в volume инициализированы другой major-версией Postgres.

Варианты:
- вернуть совместимую версию образа;
- или удалить volume и инициализировать БД заново:

```bash
docker compose --env-file .env -f docker-compose.prod.yml down -v
docker compose --env-file .env -f docker-compose.prod.yml up -d --build
```
