# 🐾 Happytail

Backend API для платформы по усыновлению домашних животных из приютов.

## Стек

- **Go 1.26.1** — основной язык
- **PostgreSQL** — база данных
- **pgx v5** — драйвер PostgreSQL
- **JWT** — аутентификация
- **Docker + Docker Compose** — окружение разработки
- **golang-migrate** — миграции БД

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/sd0hni-psina/happytail
cd happytail
```

### 2. Создать `.env` файл

```env
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_HOST=
POSTGRES_PORT=
APP_PORT=
JWT_SECRET=
PROJECT_ROOT=
```

### 3. Запустить базу данных

```bash
docker compose up -d happytail-postgres
```

### 4. Применить миграции

```bash
make migrate-up
```

### 5. Запустить сервер

```bash
make dev
```

Сервер запустится на `http://localhost:PORT`

---

## Команды Make

| Команда | Описание |
|---|---|
| `make dev` | Запустить API сервер |
| `make migrate-up` | Применить все миграции |
| `make migrate-down` | Откатить последнюю миграцию |
| `make migrate-create seq=name` | Создать новую миграцию |
| `make env-up` | Запустить только PostgreSQL |
| `make env-down` | Остановить PostgreSQL |

---

## API

### Аутентификация

| Метод | Путь            | Описание           |
|---    |---              |---                 |
| `POST`| `/auth/login`   | Получить JWT токен |

**Пример:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

Для защищённых эндпоинтов передавай токен в заголовке:
```
Authorization: Bearer <token>
```

---

### Животные

| Метод | Путь             | Защищён | Описание             |
|---    |---               |---      |---                   |
| `GET` | `/animals`       | —       | Список всех животных |
| `GET` | `/animals/{id}`  | —       | Животное по ID       |
| `POST`| `/animals`       | ✓       | Добавить животное    |

**Создать животное:**
```bash
curl -X POST http://localhost:8080/animals \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"animal_type":"кот","name":"Мурзик","age":2}'
```

---

### Приюты

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/shelters` | — | Список всех приютов |
| `GET` | `/shelters/{id}` | — | Приют по ID |
| `POST` | `/shelters` | ✓ | Добавить приют |

---

### Пользователи

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `GET` | `/users` | — | Список пользователей |
| `GET` | `/users/{id}` | — | Пользователь по ID |
| `POST` | `/users` | — | Регистрация |

---

### Усыновление

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| `POST` | `/adoptions` | ✓ | Усыновить животное |

**Усыновить животное:**
```bash
curl -X POST http://localhost:8080/adoptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"animal_id":1}'
```

---

### Health Check

```bash
curl http://localhost:8080/health
```

---

## Структура проекта

```
happytail/
├── cmd/
│   └── api/
│       └── main.go          # Точка входа
├── internal/
│   ├── config/              # Конфигурация и подключение к БД
│   ├── handler/             # HTTP обработчики
│   ├── middleware/          # Logger, Recovery, Auth
│   ├── models/              # Структуры данных
│   ├── repository/          # SQL запросы
│   └── service/             # Бизнес-логика
├── migrations/              # SQL миграции
├── docker-compose.yaml
├── Dockerfile
└── Makefile
```