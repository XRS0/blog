# Команды для Быстрого Старта

Этот файл содержит последовательность команд для первого запуска проекта.

## 1️⃣ Установка protoc

### macOS
```bash
brew install protobuf
protoc --version  # Проверка установки
```

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install -y protobuf-compiler
protoc --version  # Проверка установки
```

## 2️⃣ Установка Go плагинов

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавление в PATH (добавьте в ~/.zshrc или ~/.bashrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# Применение изменений
source ~/.zshrc  # или source ~/.bashrc

# Проверка
which protoc-gen-go
which protoc-gen-go-grpc
```

## 3️⃣ Генерация proto файлов

```bash
# Из корневой директории проекта
cd /Users/XRS0/Desktop/blog

# Генерация всех proto файлов
make proto

# Проверка что файлы созданы
ls -la services/auth-service/proto/
ls -la services/article-service/proto/
ls -la services/stats-service/proto/
ls -la services/api-gateway/proto/auth/
ls -la services/api-gateway/proto/article/
ls -la services/api-gateway/proto/stats/
```

## 4️⃣ Обновление зависимостей

```bash
# В shared модуле
cd shared
go mod tidy
cd ..

# В каждом сервисе
cd services/auth-service
go mod tidy
cd ../article-service
go mod tidy
cd ../stats-service
go mod tidy
cd ../api-gateway
go mod tidy
cd ../..
```

## 5️⃣ Запуск сервисов

```bash
# Из корневой директории
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build

# Просмотр логов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f auth-service
```

## 6️⃣ Проверка работы

### Открыть в браузере:
- API Gateway health: http://localhost:8080/health
- RabbitMQ Management: http://localhost:15672 (guest/guest)

### Через curl:

```bash
# Health check
curl http://localhost:8080/health

# Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'

# Логин (сохраните token из ответа)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Замените YOUR_TOKEN на токен из предыдущего запроса
export TOKEN="YOUR_TOKEN_HERE"

# Создание публичной статьи
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "My First Article",
    "content": "This is a test article",
    "visibility": "public"
  }'

# Создание статьи с доступом по ссылке
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Secret Article",
    "content": "Only accessible via link",
    "visibility": "link"
  }'
# Сохраните access_url из ответа

# Получение списка статей
curl http://localhost:8080/articles

# Получение конкретной статьи
curl http://localhost:8080/articles/1

# Лайк статьи
curl -X POST http://localhost:8080/articles/1/like \
  -H "Authorization: Bearer $TOKEN"

# Удаление лайка
curl -X DELETE http://localhost:8080/articles/1/like \
  -H "Authorization: Bearer $TOKEN"
```

## 7️⃣ Остановка сервисов

```bash
# Остановка
docker-compose down

# Остановка с удалением volumes (очистка БД)
docker-compose down -v

# Удаление всех контейнеров, образов и volumes
docker-compose down -v --rmi all
```

## 🔧 Полезные команды

### Просмотр логов
```bash
docker-compose logs -f auth-service
docker-compose logs -f article-service
docker-compose logs -f stats-service
docker-compose logs -f api-gateway
docker-compose logs -f db
docker-compose logs -f rabbitmq
```

### Подключение к PostgreSQL
```bash
docker exec -it blog-db psql -U blog -d blog

# В psql:
\dt                          # Список таблиц
SELECT * FROM users;
SELECT * FROM articles;
SELECT * FROM article_views;
SELECT * FROM article_likes;
\q                          # Выход
```

### Подключение к RabbitMQ Management
Откройте http://localhost:15672
- Username: guest
- Password: guest

### Перезапуск одного сервиса
```bash
docker-compose restart auth-service
```

### Пересборка одного сервиса
```bash
docker-compose up -d --build auth-service
```

### Очистка сгенерированных proto файлов
```bash
make clean
```

## 🐛 Troubleshooting

### Если protoc не найден
```bash
# Проверьте установку
which protoc

# Если не установлен, установите:
# macOS:
brew install protobuf

# Linux:
sudo apt install -y protobuf-compiler
```

### Если Go плагины не найдены
```bash
# Установите
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавьте в PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Проверьте
which protoc-gen-go
which protoc-gen-go-grpc
```

### Если Docker контейнеры не запускаются
```bash
# Проверьте статус
docker-compose ps

# Проверьте логи
docker-compose logs

# Пересоздайте с чистого листа
docker-compose down -v
docker-compose up --build
```

### Если ошибки импорта в Go
```bash
# Убедитесь что proto файлы сгенерированы
make proto

# Обновите зависимости
cd services/auth-service && go mod tidy
cd services/article-service && go mod tidy
cd services/stats-service && go mod tidy
cd services/api-gateway && go mod tidy
```

### Если порты заняты
```bash
# Проверьте занятые порты
lsof -i :8080
lsof -i :50051
lsof -i :5432
lsof -i :5672

# Остановите процессы или измените порты в docker-compose.yml
```

## 📚 Дополнительная информация

- [QUICKSTART.md](./QUICKSTART.md) - Детальный быстрый старт
- [PROTO_GENERATION.md](./PROTO_GENERATION.md) - Генерация proto
- [MIGRATION_STATUS.md](./MIGRATION_STATUS.md) - Статус миграции
- [README.md](./README.md) - Основная документация
