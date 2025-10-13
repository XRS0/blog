# Быстрый Старт - Микросервисная Архитектура

## 📋 Предварительные требования

1. **Go 1.21+** установлен
2. **Docker и Docker Compose** установлены
3. **Protocol Buffers Compiler (protoc)** установлен
4. **Make** установлен

## 🚀 Первый запуск

### 1. Установка protoc и Go плагинов

**macOS:**
```bash
# Установка protoc
brew install protobuf

# Установка Go плагинов
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавление $GOPATH/bin в PATH (если еще не добавлено)
export PATH="$PATH:$(go env GOPATH)/bin"
```

**Linux:**
```bash
# Установка protoc
sudo apt install -y protobuf-compiler

# Установка Go плагинов
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавление $GOPATH/bin в PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 2. Генерация Proto Файлов

```bash
# Из корневой директории проекта
make proto
```

Эта команда сгенерирует proto файлы локально в каждом сервисе:
- `services/auth-service/proto/`
- `services/article-service/proto/`
- `services/stats-service/proto/`
- `services/api-gateway/proto/{auth,article,stats}/`

### 3. Установка зависимостей

```bash
# В shared модуле
cd shared && go mod tidy

# В каждом сервисе
cd ../services/auth-service && go mod tidy
cd ../article-service && go mod tidy
cd ../stats-service && go mod tidy
cd ../api-gateway && go mod tidy
```

### 4. Запуск через Docker Compose

```bash
# Вернуться в корень проекта
cd ../..

# Запустить все сервисы
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build
```

## 🔍 Проверка работы сервисов

### Доступные порты:

- **API Gateway**: http://localhost:8080
- **Auth Service (gRPC)**: localhost:50051
- **Article Service (gRPC)**: localhost:50052
- **Stats Service (gRPC)**: localhost:50053
- **PostgreSQL**: localhost:5432
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

### Проверка здоровья:

```bash
# Проверка API Gateway
curl http://localhost:8080/health

# Проверка RabbitMQ
curl -u guest:guest http://localhost:15672/api/overview

# Проверка PostgreSQL
docker exec -it blog-db psql -U blog -d blog -c "SELECT version();"
```

## 📝 Примеры запросов к API

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### Логин
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Сохраните `token` из ответа для дальнейших запросов.

### Создание публичной статьи
```bash
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "My First Article",
    "content": "This is the content",
    "visibility": "public"
  }'
```

### Создание статьи с доступом по ссылке
```bash
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Secret Article",
    "content": "Only accessible via link",
    "visibility": "link"
  }'
```

В ответе вы получите `access_url` для доступа к статье.

### Получение списка статей
```bash
# Публичные статьи (без авторизации)
curl http://localhost:8080/articles

# Все ваши статьи (с авторизацией)
curl http://localhost:8080/articles \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Доступ к статье по access_token
```bash
curl "http://localhost:8080/articles/1?access_token=UUID"
```

## 🛠️ Разработка

### Структура проекта
```
blog/
├── proto/                  # Proto схемы (не генерируются)
│   ├── auth.proto
│   ├── article.proto
│   └── stats.proto
├── services/
│   ├── api-gateway/        # REST API точка входа
│   │   └── proto/          # Сгенерированные proto (gitignore)
│   ├── auth-service/       # Аутентификация и пользователи
│   │   └── proto/          # Сгенерированные proto (gitignore)
│   ├── article-service/    # CRUD статей с visibility
│   │   └── proto/          # Сгенерированные proto (gitignore)
│   └── stats-service/      # Статистика просмотров и лайков
│       └── proto/          # Сгенерированные proto (gitignore)
├── shared/                 # Общие библиотеки
│   ├── database/          # Bun DB и автомиграции
│   ├── logger/            # Логирование
│   └── rabbitmq/          # RabbitMQ клиент
├── Makefile               # Команды для генерации proto и сборки
└── docker-compose.yml     # Оркестрация сервисов
```

### Команды разработки

```bash
# Генерация proto файлов
make proto

# Очистка сгенерированных файлов
make clean

# Сборка всех сервисов
make build-services

# Запуск с docker-compose
make docker-up

# Остановка docker-compose
make docker-down
```

### Локальный запуск сервиса (без Docker)

```bash
# 1. Запустите PostgreSQL и RabbitMQ
docker-compose up -d db rabbitmq

# 2. Установите переменные окружения
export DATABASE_URL="postgres://blog:blog@localhost:5432/blog?sslmode=disable"
export RABBITMQ_URL="amqp://blog:blog@localhost:5672/"
export LOG_LEVEL="debug"

# 3. Запустите сервис
cd services/auth-service
go run cmd/main.go
```

## 🔄 После изменения proto файлов

1. Обновите `.proto` файлы в директории `proto/`
2. Регенерируйте код:
   ```bash
   make clean
   make proto
   ```
3. Обновите зависимости в сервисах:
   ```bash
   cd services/auth-service && go mod tidy
   # Повторите для всех сервисов
   ```
4. Перезапустите сервисы:
   ```bash
   docker-compose down
   docker-compose up --build
   ```

## 📚 Дополнительная документация

- **[PROTO_GENERATION.md](./PROTO_GENERATION.md)** - Детальная инструкция по генерации proto файлов
- **[MICROSERVICES_ARCHITECTURE.md](./MICROSERVICES_ARCHITECTURE.md)** - Архитектура микросервисов
- **[GETTING_STARTED.md](./GETTING_STARTED.md)** - Подробное руководство по запуску
- **[FRONTEND_INTEGRATION.md](./FRONTEND_INTEGRATION.md)** - Интеграция с фронтендом

## ❗ Важные замечания

1. **Proto файлы не в git**: Сгенерированные `.pb.go` файлы не коммитятся в репозиторий. Каждый разработчик должен сгенерировать их локально после клонирования.

2. **Автомиграции БД**: Все сервисы используют Bun ORM с автомиграциями. При первом запуске таблицы создаются автоматически.

3. **Visibility система**: Статьи поддерживают три режима видимости:
   - `public` - доступны всем
   - `private` - доступны только автору
   - `link` - доступны по UUID токену (access_token)

4. **RabbitMQ события**: Stats Service подписан на события:
   - `article.created` - новая статья
   - `article.viewed` - просмотр статьи
   - `article.liked` - лайк статьи
   - `article.unliked` - снятие лайка

## 🐛 Troubleshooting

### "could not import proto package"
```bash
# Убедитесь что proto файлы сгенерированы
make proto

# Обновите зависимости
cd services/auth-service && go mod tidy
```

### "connection refused" при запуске сервисов
```bash
# Убедитесь что PostgreSQL и RabbitMQ запущены
docker-compose ps

# Проверьте логи
docker-compose logs db
docker-compose logs rabbitmq
```

### База данных не инициализируется
```bash
# Пересоздайте контейнеры с удалением volumes
docker-compose down -v
docker-compose up --build
```
