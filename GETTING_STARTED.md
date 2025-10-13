# Начало работы с микросервисной архитектурой

## Предварительные требования

- Go 1.21+
- Docker и Docker Compose
- Protocol Buffers компилятор (`protoc`)
- Go plugins для protoc:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

## Быстрый старт

### 1. Генерация proto файлов

```bash
make proto
```

Это сгенерирует Go код из proto файлов в `proto/gen/`.

### 2. Применение миграции БД

Перед первым запуском примените миграцию для добавления полей visibility:

```bash
# Подключитесь к БД (после запуска docker-compose)
docker-compose exec db psql -U blog -d blog -f /migrations/add_article_visibility.sql
```

Или вручную:

```bash
psql -U blog -h localhost -d blog < migrations/add_article_visibility.sql
```

### 3. Запуск всех сервисов

```bash
docker-compose up -d
```

Сервисы будут доступны на:
- **API Gateway**: http://localhost:8080
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **PostgreSQL**: localhost:5432

### 4. Проверка работоспособности

```bash
curl http://localhost:8080/healthz
```

## Разработка

### Запуск отдельных сервисов локально

Для разработки можете запускать сервисы по отдельности:

```bash
# Запустить только БД и RabbitMQ
docker-compose up -d db rabbitmq

# Auth Service
cd services/auth-service
go mod tidy
go run cmd/main.go

# Article Service
cd services/article-service
go mod tidy
go run cmd/main.go

# Stats Service
cd services/stats-service
go mod tidy
go run cmd/main.go

# API Gateway
cd services/api-gateway
go mod tidy
go run cmd/main.go
```

### Переменные окружения для локальной разработки

Создайте `.env` файлы или экспортируйте переменные:

```bash
# Общие
export DATABASE_URL="postgres://blog:blog@localhost:5432/blog?sslmode=disable"
export RABBITMQ_URL="amqp://blog:blog@localhost:5672/"
export JWT_SECRET="dev-secret-change-me"
export LOG_LEVEL="debug"

# Auth Service
export GRPC_PORT="50051"

# Article Service
export GRPC_PORT="50052"
export AUTH_SERVICE_URL="localhost:50051"

# Stats Service
export GRPC_PORT="50053"

# API Gateway
export PORT="8080"
export AUTH_SERVICE_URL="localhost:50051"
export ARTICLE_SERVICE_URL="localhost:50052"
export STATS_SERVICE_URL="localhost:50053"
export ALLOW_ORIGIN="http://localhost:5173"
```

## API Примеры использования

### Регистрация

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### Логин

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Сохраните токен из ответа.

### Создание публичной статьи

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "My Public Article",
    "content": "This is a public article that everyone can see.",
    "visibility": "public"
  }'
```

### Создание приватной статьи

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "My Private Notes",
    "content": "This is only visible to me.",
    "visibility": "private"
  }'
```

### Создание статьи с доступом по ссылке

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Shared Article",
    "content": "This article is accessible via a special link.",
    "visibility": "link"
  }'
```

Ответ будет содержать `access_token` и `access_url`.

### Просмотр статьи по ссылке

```bash
# Без access_token - доступ запрещен
curl http://localhost:8080/api/articles/1

# С access_token - доступ разрешен
curl "http://localhost:8080/api/articles/1?access_token=UUID_HERE"
```

### Получение всех публичных статей

```bash
curl http://localhost:8080/api/articles
```

### Лайк статьи

```bash
# Поставить лайк
curl -X POST http://localhost:8080/api/articles/1/like \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"like": true}'

# Убрать лайк
curl -X POST http://localhost:8080/api/articles/1/like \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"like": false}'
```

### Обновление статьи (смена visibility)

```bash
curl -X PUT http://localhost:8080/api/articles/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content",
    "visibility": "link"
  }'
```

### Удаление статьи

```bash
curl -X DELETE http://localhost:8080/api/articles/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Мониторинг

### Логи сервисов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f api-gateway
docker-compose logs -f auth-service
docker-compose logs -f article-service
docker-compose logs -f stats-service
```

### RabbitMQ Management

Откройте http://localhost:15672 (guest/guest) для мониторинга очередей и сообщений.

### Проверка gRPC сервисов

Установите `grpcurl`:

```bash
brew install grpcurl  # macOS
```

Проверка методов:

```bash
# Auth Service
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 auth.AuthService/ValidateToken

# Article Service
grpcurl -plaintext localhost:50052 list

# Stats Service
grpcurl -plaintext localhost:50053 list
```

## Остановка сервисов

```bash
docker-compose down
```

С удалением данных:

```bash
docker-compose down -v
```

## Troubleshooting

### Proto файлы не генерируются

Убедитесь, что установлены плагины:

```bash
which protoc-gen-go
which protoc-gen-go-grpc
```

### Ошибка подключения к БД

Проверьте, что PostgreSQL запущен:

```bash
docker-compose ps db
docker-compose logs db
```

### Сервисы не могут подключиться друг к другу

Убедитесь, что все сервисы в одной Docker сети:

```bash
docker network inspect blog_blog-network
```

### RabbitMQ не принимает сообщения

Проверьте логи и статус:

```bash
docker-compose logs rabbitmq
docker-compose exec rabbitmq rabbitmq-diagnostics status
```

## Следующие шаги

1. **AI Service**: Добавьте микросервис для помощи в написании статей
2. **Notification Service**: Уведомления о новых статьях, лайках
3. **Search Service**: Полнотекстовый поиск по статьям с Elasticsearch
4. **Media Service**: Загрузка и хранение изображений для статей
5. **Analytics Service**: Детальная аналитика просмотров и поведения

Для добавления нового сервиса смотрите пример существующих в `services/`.
