
# Сервис аутентификации и авторизации

Данный проект представляет собой сервис аутентификации и авторизации, написанный на **Go**. В качестве основной базы данных используется **PostgreSQL**, а для кэширования данных — **Redis**. Ниже представлен подробный обзор структуры и функциональности проекта.

---

## Возможности

- **Аутентификация и авторизация пользователей** через gRPC API.
- Безопасная обработка паролей с использованием **bcrypt**.
- Генерация и проверка **JWT-токенов** для управления сессиями.
- Управление кэшем через **Redis**.
- Операции с базой данных через **PostgreSQL**.
- Гибкая настройка через YAML-файлы.

---

## Структура проекта

```plaintext
.
├── cmd
│   ├── migrator
│   │   └── main.go          # Мигратор для создания схемы базы данных
│   └── sso
│       └── main.go          # Точка входа основного приложения
├── config
│   ├── example.yaml         # Пример конфигурации
│   ├── local.yaml           # Конфигурация для локальной разработки
│   └── remote.yaml          # Конфигурация для удаленной среды (продакшн)
├── docker-compose.yaml      # Конфигурация Docker Compose
├── Dockerfile.app           # Dockerfile для контейнера приложения
├── Dockerfile.migrator      # Dockerfile для контейнера мигратора
├── go.mod                   # Зависимости Go
├── go.sum                   # Контрольные суммы зависимостей Go
├── internal                 # Основная логика приложения
│   ├── app                  # Ядро приложения
│   │   ├── app.go           # Структура приложения и инициализация
│   │   └── grpc
│   │       └── app.go       # Настройка gRPC сервера
│   ├── config
│   │   └── config.go        # Загрузка конфигурации
│   ├── grpc
│   │   └── auth
│   │       └── server.go    # Реализация gRPC сервиса аутентификации
│   ├── lib                  # Вспомогательные библиотеки
│   │   ├── jwt
│   │   │   └── jwt.go       # Утилиты для работы с JWT
│   │   └── logger
│   │       ├── handlers
│   │       │   └── slogpretty
│   │       │       └── slogpretty.go # Форматирование логов
│   │       └── sl
│   │           └── error.go # Логирование ошибок
│   ├── logger
│   │   └── logger.go        # Настройка логгера
│   ├── models               # Модели данных
│   │   ├── app.go
│   │   └── user.go
│   ├── services
│   │   └── auth
│   │       └── auth.go      # Логика сервиса аутентификации
│   └── storage              # Интерфейсы и реализации хранилищ
│       ├── cache
│       │   ├── cache.go     # Интерфейс кэша
│       │   └── redis
│       │       └── redis.go # Реализация Redis-кэша
│       └── main_storage
│           ├── postgres
│           │   └── postgres.go # Реализация для PostgreSQL
│           └── storage.go   # Интерфейс основной базы данных
├── migrations               # SQL-миграции для базы данных
│   ├── 1_init.down.sql
│   ├── 1_init.up.sql
│   ├── 2_test_app.down.sql
│   └── 2_test_app.up.sql
├── Taskfile.yaml            # Конфигурация для Taskfile
└── tests                    # Тесты
    ├── register_login_test.go
    └── suite
        └── suite.go
```

---

## Конфигурация

### Директория `config`
- `example.yaml` — шаблон конфигурации.
- `local.yaml` — конфигурация для локальной разработки.
- `remote.yaml` — конфигурация для удаленного окружения (продакшн). 

В файлах конфигурации указаны:
- **Порты** и **хосты** для приложения и сервисов.
- Настройки подключения к базе данных.
- Настройки для кэша (Redis).

### Переменные окружения
Путь к конфигурационному файлу должен быть задан через переменную окружения. Пример:
```bash
CONFIG_PATH=./config/local.yaml
```

---

## Запуск проекта

### Предварительные требования
- Установленные **Docker** и **docker-compose**.
- Установленный **Go** версии 1.19 или выше.

### Запуск через Docker Compose
Чтобы собрать и запустить приложение в контейнере:
```bash
docker-compose up --build
```

### Локальный запуск
1. Запустите PostgreSQL и Redis.
2. Установите переменные окружения:
   ```bash
   export CONFIG_PATH=./config/local.yaml
   export POSTGRES_PASS=your_password
   export REDIS_PASS=your_password 
   ```
3. Выполните миграции для создания схемы базы данных:
   ```bash
   go run ./cmd/migrator/main.go
   ```
4. Запустите приложение:
   ```bash
   go run ./cmd/sso/main.go
   ```

---

## Основные компоненты

### Сервис аутентификации
Реализован в `internal/services/auth/auth.go`:
- **Регистрация**: `Register(ctx, login, password)`
- **Вход**: `Login(ctx, login, password, appID)`
- **Проверка администратора**: `IsAdmin(ctx, userID)`

### Хранилище
- **Кэш**: Реализация на Redis в `internal/storage/cache/redis/redis.go`.
- **База данных**: Реализация на PostgreSQL в `internal/storage/main_storage/postgres/postgres.go`.

### JWT-токены
- Генерация и проверка токенов реализована с использованием библиотеки `github.com/golang-jwt/jwt/v5`.

### Proto файл
- Релизовал его в данном репозитории: github.com/notblinkyet/proto_sso

---

## Тестирование
Тесты находятся в папке `tests/`. Для запуска выполните:
```bash
go test ./tests/...
```

---

## Зависимости

Проект использует следующие зависимости:
- `github.com/golang-jwt/jwt/v5` — для генерации и проверки JWT-токенов.
- `golang.org/x/crypto/bcrypt` — для хэширования и проверки паролей.
- `github.com/jackc/pgx/v5/pgxpool` — для подключения к PostgreSQL.
- `github.com/redis/go-redis/v9` — клиент для работы с Redis.

---