# Шпаргалка по Работе с Proto Файлами

## 🚀 Быстрый старт

```bash
# 1. Установка protoc (macOS без Homebrew)
./install-protoc.sh

# 2. Генерация proto файлов
make proto

# 3. Проверка структуры (должно быть БЕЗ дополнительной папки proto/)
ls services/api-gateway/proto/auth/
# ✅ Должны видеть: auth.pb.go, auth_grpc.pb.go
# ❌ НЕ должно быть: proto/

# 4. Обновление зависимостей
cd services/auth-service && go mod tidy
cd ../article-service && go mod tidy
cd ../stats-service && go mod tidy
cd ../api-gateway && go mod tidy
cd ../..

# 5. Запуск
docker-compose up --build
```

## 📁 Правильная Структура

```
services/
├── auth-service/proto/
│   ├── auth.pb.go          ✅
│   └── auth_grpc.pb.go     ✅
│
├── article-service/proto/
│   ├── article.pb.go       ✅
│   └── article_grpc.pb.go  ✅
│
├── stats-service/proto/
│   ├── stats.pb.go         ✅
│   └── stats_grpc.pb.go    ✅
│
└── api-gateway/proto/
    ├── auth/
    │   ├── auth.pb.go      ✅
    │   └── auth_grpc.pb.go ✅
    ├── article/
    │   ├── article.pb.go   ✅
    │   └── article_grpc.pb.go ✅
    └── stats/
        ├── stats.pb.go     ✅
        └── stats_grpc.pb.go ✅
```

## ⚠️ Неправильная Структура (чего НЕ должно быть)

```
services/api-gateway/proto/auth/proto/auth.pb.go      ❌ Лишняя папка proto/
services/api-gateway/proto/article/proto/article.pb.go ❌ Лишняя папка proto/
```

## 🔍 Проверочные Команды

```bash
# Найти все сгенерированные файлы
find services -name "*.pb.go" -o -name "*_grpc.pb.go"

# Проверить на лишние папки proto/
find services -path "*/proto/proto" && echo "❌ Найдены лишние папки!" || echo "✅ Структура правильная"

# Количество сгенерированных файлов (должно быть 12)
find services -name "*.pb.go" -o -name "*_grpc.pb.go" | wc -l
```

## 📝 Импорты в Go Коде

### Auth Service
```go
import pb "github.com/XRS0/blog/services/auth-service/proto"

// Использование
user := &pb.User{...}
```

### Article Service
```go
import pb "github.com/XRS0/blog/services/article-service/proto"

// Использование
article := &pb.Article{...}
```

### Stats Service
```go
import pb "github.com/XRS0/blog/services/stats-service/proto"

// Использование
stats := &pb.ArticleStats{...}
```

### API Gateway
```go
import (
    authpb "github.com/XRS0/blog/services/api-gateway/proto/auth"
    articlepb "github.com/XRS0/blog/services/api-gateway/proto/article"
    statspb "github.com/XRS0/blog/services/api-gateway/proto/stats"
)

// Использование
authClient := authpb.NewAuthServiceClient(conn)
articleClient := articlepb.NewArticleServiceClient(conn)
statsClient := statspb.NewStatsServiceClient(conn)
```

## 🔧 Команды Make

```bash
make proto          # Сгенерировать все proto файлы
make proto-auth     # Только auth proto
make proto-article  # Только article proto
make proto-stats    # Только stats proto
make clean          # Удалить все сгенерированные файлы
make build-services # Собрать все сервисы
make docker-up      # Запустить docker-compose
make docker-down    # Остановить docker-compose
```

## 🐛 Troubleshooting

### Проблема: Создается дополнительная папка proto/

**Решение:**
```bash
# Проверьте Makefile - должно быть:
protoc --proto_path=proto \
       --go_out=services/api-gateway/proto/auth \
       auth.proto  # БЕЗ префикса proto/

# НЕ должно быть:
protoc --go_out=services/api-gateway/proto/auth \
       proto/auth.proto  # ❌ С префиксом proto/
```

### Проблема: protoc не найден

**Решение:**
```bash
# macOS без Homebrew
./install-protoc.sh

# macOS с Homebrew
brew install protobuf

# Linux
sudo apt install -y protobuf-compiler

# Проверка
protoc --version
```

### Проблема: Go плагины не найдены

**Решение:**
```bash
# Установка
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавить в PATH (в ~/.zshrc или ~/.bashrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# Применить
source ~/.zshrc

# Проверка
which protoc-gen-go
which protoc-gen-go-grpc
```

### Проблема: Ошибки импорта в Go

**Решение:**
```bash
# 1. Убедитесь что proto файлы сгенерированы
make proto

# 2. Обновите зависимости
cd services/auth-service && go mod tidy
cd services/article-service && go mod tidy
cd services/stats-service && go mod tidy
cd services/api-gateway && go mod tidy

# 3. Проверьте импорты в коде
grep -r "proto/gen" services/  # Не должно быть старых импортов
```

## 📚 Документация

- `PROTO_STRUCTURE_EXPLAINED.md` - Детальное объяснение структуры
- `PROTO_GENERATION.md` - Руководство по генерации
- `README.md` - Общая информация о проекте
- `QUICKSTART.md` - Быстрый старт

## 🎯 Контрольный Список

Перед запуском убедитесь:

- [ ] `protoc` установлен (`protoc --version`)
- [ ] Go плагины установлены (`which protoc-gen-go`)
- [ ] `$GOPATH/bin` добавлен в `$PATH`
- [ ] Proto файлы сгенерированы (`make proto`)
- [ ] Структура правильная (нет лишних папок `proto/`)
- [ ] Зависимости обновлены (`go mod tidy` в каждом сервисе)
- [ ] Docker Compose запущен (`docker-compose up`)

## ⚡ Один Скрипт для Всего

Создайте файл `setup.sh`:

```bash
#!/bin/bash
set -e

echo "🚀 Настройка проекта..."

# Установка protoc (если нужно)
if ! command -v protoc &> /dev/null; then
    echo "📦 Установка protoc..."
    ./install-protoc.sh
fi

# Генерация proto
echo "🔨 Генерация proto файлов..."
make clean
make proto

# Обновление зависимостей
echo "📚 Обновление зависимостей..."
for service in auth-service article-service stats-service api-gateway; do
    echo "  - $service"
    (cd "services/$service" && go mod tidy)
done

echo "✅ Готово! Можно запускать: docker-compose up --build"
```

Использование:
```bash
chmod +x setup.sh
./setup.sh
```
