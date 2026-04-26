# 🐾 Happytail

Backend API для платформы по усыновлению домашних животных из приютов.

## Стек

- **Go 1.26.1** — основной язык
- **PostgreSQL** — база данных
- **pgx v5** — драйвер PostgreSQL
- **JWT + Refresh Tokens** — аутентификация
- **MinIO** — хранилище фотографий (S3-совместимое)
- **Docker + Docker Compose** — окружение разработки
- **golang-migrate** — миграции БД
- **Prometheus** — метрики
- **slog** — структурированное логирование
- **gomail** — email уведомления

---

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/sd0hni-psina/happytail
cd happytail
```

### 2. Создать `.env` файл

```env
# PostgreSQL
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_HOST=happytail-postgres
POSTGRES_PORT=5432

# Приложение
APP_PORT=8080
APP_ENV=development
JWT_SECRET=               # минимум 32 символа
ALLOWED_ORIGIN=*          # для продакшена указать конкретный домен

# SMTP (email уведомления)
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=

# MinIO (хранилище фото)
MINIO_ENDPOINT=happytail-minio:9000
MINIO_PUBLIC_URL=http://localhost:9000
MINIO_USER=
MINIO_PASSWORD=
MINIO_BUCKET=photos

# Служебное
PROJECT_ROOT=
```

### 3. Запустить окружение

```bash
make dev
```

Запускает PostgreSQL, MinIO и API сервер. Сервер доступен на `http://localhost:8080`.

### 4. Применить миграции

```bash
make migrate-up
```

### 5. MinIO Web UI

Доступен на `http://localhost:9001` — для просмотра загруженных фотографий.

---

## Команды Make

| Команда | Описание |
|---|---|
| `make dev` | Запустить API + MinIO |
| `make test` | Запустить тесты |
| `make test-cover` | Тесты с отчётом покрытия |
| `make migrate-up` | Применить все миграции |
| `make migrate-down` | Откатить последнюю миграцию |
| `make migrate-create seq=name` | Создать новую миграцию |
| `make env-up` | Запустить только PostgreSQL |
| `make env-down` | Остановить PostgreSQL |
| `make env-cleanup` | Очистить данные БД |
| `make swag` | Регенерировать Swagger документацию |

---

## API

Swagger UI доступен на `http://localhost:8080/swagger/`.

Метрики Prometheus доступны на `http://localhost:8080/metrics`.

### Аутентификация

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/auth/login` | Логин — возвращает access + refresh токены |
| `POST` | `/auth/refresh` | Обновить access token |
| `POST` | `/auth/logout` | Выход — отозвать refresh token |

Для защищённых эндпоинтов передавай access токен в заголовке:
```
Authorization: Bearer <access_token>
```

### Пользователи

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/users` | — | Список пользователей |
| `GET` | `/users/{id}` | — | Пользователь по ID |
| `GET` | `/users/me` | ✓ | Свой профиль |
| `POST` | `/users` | — | Регистрация |

### Животные

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/animals` | — | Список животных (с пагинацией и фильтрами) |
| `GET` | `/animals/{id}` | — | Животное по ID |
| `POST` | `/animals` | ✓ | Добавить животное |

**Фильтры для GET /animals:**
```
?type=cat&status=available&is_vaccinated=true&page=1&limit=10
```

### Фотографии животных

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/animals/{id}/photos` | — | Фото животного |
| `POST` | `/animals/{id}/photos` | ✓ | Загрузить фото (multipart/form-data) |
| `PATCH` | `/animals/{id}/photos/{photo_id}/main` | ✓ | Сделать фото главным |
| `DELETE` | `/animals/{id}/photos/{photo_id}` | ✓ | Удалить фото |

**Загрузка фото:**
```bash
curl -X POST http://localhost:8080/animals/1/photos \
  -H "Authorization: Bearer <token>" \
  -F "photo=@/path/to/photo.jpg" \
  -F "is_main=true"
```

### Приюты

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/shelters` | — | Список приютов |
| `GET` | `/shelters/{id}` | — | Приют по ID |
| `GET` | `/shelters/nearby` | — | Ближайшие приюты по геолокации |
| `POST` | `/shelters` | ✓ | Добавить приют |

**Поиск ближайших приютов:**
```bash
curl "http://localhost:8080/shelters/nearby?lat=47.1167&lon=51.8833&radius=50"
```

### Усыновление

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `POST` | `/adoptions` | ✓ | Усыновить животное |

После усыновления пользователь получает email подтверждение. Статус животного меняется на `adopted`.

### Посты

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/posts` | — | Список постов |
| `GET` | `/posts/{id}` | — | Пост по ID |
| `POST` | `/posts` | ✓ | Создать пост |

### Роли

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `POST` | `/roles` | ✓ Admin | Назначить роль |
| `DELETE` | `/roles/{id}` | ✓ Admin | Удалить роль |

Доступные роли: `admin`, `shelter_admin`, `user`.

### Health Check & Метрики

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

---

## Структура проекта

```
happytail/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа
├── internal/
│   ├── config/                  # Конфигурация, валидация, DB pool
│   ├── handler/                 # HTTP обработчики
│   ├── logger/                  # Structured logging (slog)
│   ├── middleware/              # Logger, Recovery, Auth, CORS, RateLimit, Metrics
│   ├── models/                  # Структуры данных
│   ├── notifier/                # Email уведомления
│   ├── repository/              # SQL запросы
│   ├── service/                 # Бизнес-логика
│   └── storage/                 # MinIO интеграция
├── migrations/                  # SQL миграции (11 файлов)
├── docker-compose.yaml
├── Dockerfile
└── Makefile
```

---

## Безопасность

- JWT access tokens (15 минут) + refresh tokens (30 дней) с rotation
- Rate limiting: 10 req/sec per IP, burst 30
- CORS настраивается через `ALLOWED_ORIGIN`
- Валидация конфига при старте — сервер не запустится без обязательных переменных
- `password_hash` никогда не возвращается в публичных API ответах
- Усыновление защищено от race condition через `SELECT FOR UPDATE`