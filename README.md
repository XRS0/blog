# Blue Note Blog - Микросервисная Архитектура

Современное блог-приложение на микросервисной архитектуре с **Go + gRPC + RabbitMQ** на бэкенде и **React (Vite + TypeScript)** на фронтенде.

## ✨ Возможности

- 🔐 **Аутентификация** - JWT токены, регистрация и логин
- 📝 **CRUD статей** - создание, просмотр, редактирование, удаление
- 👁️ **Видимость статей** - публичные, приватные, или доступные по ссылке
- 📊 **Статистика** - подсчёт просмотров и лайков в реальном времени
- 🎨 **Markdown поддержка** - форматирование контента статей
- 🔄 **Асинхронные события** - через RabbitMQ
- 🗄️ **Автомиграции БД** - через Bun ORM

## 🏗️ Архитектура

Проект использует микросервисную архитектуру с gRPC для межсервисного взаимодействия:

```
┌─────────────┐
│   Frontend  │ ← React + TypeScript + Vite
└──────┬──────┘
       │ REST API
       ↓
┌─────────────┐
│ API Gateway │ ← Gin REST → gRPC маршрутизация
└──────┬──────┘
       │ gRPC
       ├────────────┬──────────────┬────────────┐
       ↓            ↓              ↓            ↓
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│   Auth   │  │ Article  │  │  Stats   │  │ RabbitMQ │
│ Service  │  │ Service  │  │ Service  │  │          │
└────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     │             │              │             │
     └─────────────┴──────────────┴─────────────┘
                       │
                  PostgreSQL
```

### Сервисы:

- **API Gateway** (:8080) - REST API точка входа
- **Auth Service** (:50051) - Аутентификация и управление пользователями
- **Article Service** (:50052) - CRUD статей с поддержкой visibility
- **Stats Service** (:50053) - Статистика просмотров и лайков
- **RabbitMQ** (:5672, :15672) - Брокер сообщений для событий
- **PostgreSQL** (:5432) - База данных

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.21+
- Docker и Docker Compose
- Protocol Buffers Compiler (protoc)
- Make

### Установка и запуск

1. **Установите protoc и Go плагины:**

```bash
# macOS (с Homebrew)
brew install protobuf

# macOS (без Homebrew) - автоматическая установка
./install-protoc.sh

# Linux
sudo apt install -y protobuf-compiler

# Go плагины (для всех ОС)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавьте в PATH (добавьте в ~/.zshrc или ~/.bashrc)
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc  # или source ~/.bashrc
```

2. **Автоматическая настройка (рекомендуется):**

```bash
# Один скрипт для всего: установка protoc, генерация proto, обновление зависимостей
./setup.sh
```

**ИЛИ вручную:**

```bash
# Генерация proto файлов
make proto

# Обновление зависимостей
cd services/auth-service && go mod tidy
cd ../article-service && go mod tidy
cd ../stats-service && go mod tidy
cd ../api-gateway && go mod tidy
cd ../..
```

3. **Запустите все сервисы:**

```bash
docker-compose up --build
```

API Gateway будет доступен по адресу: `http://localhost:8080`

### Примеры запросов

```bash
# Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","username":"testuser","password":"password123"}'

# Логин
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Создание статьи (с токеном)
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title":"My Article","content":"Content here","visibility":"public"}'
```

## 📝 Visibility система

Статьи поддерживают три режима видимости:

- **public** - доступны всем пользователям
- **private** - доступны только автору
- **link** - доступны по уникальной ссылке с UUID токеном

## 📚 Документация

- **[QUICKSTART.md](./QUICKSTART.md)** - Быстрый старт и основные команды
- **[PROTO_GENERATION.md](./PROTO_GENERATION.md)** - Генерация proto файлов
- **[MICROSERVICES_ARCHITECTURE.md](./MICROSERVICES_ARCHITECTURE.md)** - Архитектура
- **[GETTING_STARTED.md](./GETTING_STARTED.md)** - Детальное руководство
- **[FRONTEND_INTEGRATION.md](./FRONTEND_INTEGRATION.md)** - Интеграция фронтенда

## 🛠️ Технологии

**Backend:**
- Go 1.21
- gRPC + Protocol Buffers
- Bun ORM с автомиграциями
- RabbitMQ для событий
- PostgreSQL
- JWT аутентификация

**Frontend:**
- React 18
- TypeScript
- Vite
- Markdown поддержка

## 🔧 Разработка

```bash
# Генерация proto файлов
make proto

# Очистка сгенерированных файлов
make clean

# Сборка всех сервисов
make build-services

# Запуск с docker-compose
make docker-up

# Остановка
make docker-down
```

### Локальный запуск сервиса

```bash
# Запустите PostgreSQL и RabbitMQ
docker-compose up -d db rabbitmq

# Установите переменные окружения
export DATABASE_URL="postgres://blog:blog@localhost:5432/blog?sslmode=disable"
export RABBITMQ_URL="amqp://blog:blog@localhost:5672/"

# Запустите сервис
cd services/auth-service
go run cmd/main.go
```

## 📡 API Endpoints

### Аутентификация
- `POST /auth/register` - Регистрация нового пользователя
- `POST /auth/login` - Вход в систему
- `GET /auth/me` - Получить информацию о текущем пользователе

### Статьи
- `GET /articles` - Список статей (публичные + свои приватные)
- `POST /articles` - Создать статью (требует авторизацию)
- `GET /articles/:id` - Получить статью (увеличивает просмотры)
- `GET /articles/:id?access_token=UUID` - Доступ к статье по ссылке
- `PUT /articles/:id` - Обновить статью (только автор)
- `DELETE /articles/:id` - Удалить статью (только автор)

### Статистика
- `POST /articles/:id/like` - Поставить лайк
- `DELETE /articles/:id/like` - Убрать лайк

## 🌐 Frontend (в разработке)

```bash
cd frontend
npm install
npm run dev
```

Vite dev server: `http://localhost:5173`

Подробности интеграции смотрите в [FRONTEND_INTEGRATION.md](./FRONTEND_INTEGRATION.md)

## 🔒 Безопасность

- JWT токены для аутентификации
- Валидация токенов в API Gateway
- Проверка прав доступа к статьям
- UUID токены для доступа по ссылке
- Безопасное хранение паролей

## 🐛 Troubleshooting

### Proto файлы не найдены
```bash
make proto
cd services/auth-service && go mod tidy
```

### Ошибки подключения к БД
```bash
docker-compose down -v
docker-compose up --build
```

### RabbitMQ connection refused
```bash
docker-compose logs rabbitmq
# Подождите пока RabbitMQ полностью запустится (~10-15 сек)
```

## 📜 Лицензия

MIT

## 🤝 Вклад

Pull requests приветствуются! Для крупных изменений сначала откройте issue для обсуждения.
