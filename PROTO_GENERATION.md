# Генерация Proto Файлов

## Установка необходимых инструментов

### 1. Установка Protocol Buffers Compiler (protoc)

#### macOS
```bash
# Через Homebrew
brew install protobuf

# Или скачать напрямую
# Перейдите на https://github.com/protocolbuffers/protobuf/releases
# Скачайте protoc-{VERSION}-osx-x86_64.zip или protoc-{VERSION}-osx-aarch_64.zip
# Распакуйте и переместите в /usr/local/bin
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt install -y protobuf-compiler

# Или скачать напрямую
PROTOC_VERSION=25.1
curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip
unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"
```

### 2. Установка Go плагинов для protoc

```bash
# Плагин для генерации Go кода
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Плагин для генерации gRPC кода
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Убедитесь что $GOPATH/bin в вашем PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 3. Проверка установки

```bash
# Проверка protoc
protoc --version
# Ожидается: libprotoc 3.x.x или выше

# Проверка Go плагинов
which protoc-gen-go
which protoc-gen-go-grpc
```

## Генерация Proto Файлов

### Структура proto файлов

После генерации файлы будут размещены локально в каждом сервисе:

```
services/
├── auth-service/
│   └── proto/              # Сгенерированные файлы для auth
│       ├── auth.pb.go
│       └── auth_grpc.pb.go
├── article-service/
│   └── proto/              # Сгенерированные файлы для article
│       ├── article.pb.go
│       └── article_grpc.pb.go
├── stats-service/
│   └── proto/              # Сгенерированные файлы для stats
│       ├── stats.pb.go
│       └── stats_grpc.pb.go
└── api-gateway/
    └── proto/
        ├── auth/           # Копия auth proto для gateway
        │   ├── auth.pb.go
        │   └── auth_grpc.pb.go
        ├── article/        # Копия article proto для gateway
        │   ├── article.pb.go
        │   └── article_grpc.pb.go
        └── stats/          # Копия stats proto для gateway
            ├── stats.pb.go
            └── stats_grpc.pb.go
```

**Важно:** Файлы генерируются БЕЗ дополнительной вложенной папки `proto/`. 
То есть в `api-gateway/proto/auth/` будут сразу файлы `auth.pb.go` и `auth_grpc.pb.go`, 
а НЕ `api-gateway/proto/auth/proto/auth.pb.go`.

### Команды генерации

```bash
# Сгенерировать все proto файлы
make proto

# Сгенерировать только auth proto
make proto-auth

# Сгенерировать только article proto
make proto-article

# Сгенерировать только stats proto
make proto-stats

# Удалить все сгенерированные файлы
make clean
```

## Импорты в Go коде

### В каждом сервисе используется свой локальный proto пакет:

#### Auth Service
```go
import pb "github.com/XRS0/blog/services/auth-service/proto"
```

#### Article Service
```go
import pb "github.com/XRS0/blog/services/article-service/proto"
```

#### Stats Service
```go
import pb "github.com/XRS0/blog/services/stats-service/proto"
```

#### API Gateway
```go
import (
    authpb "github.com/XRS0/blog/services/api-gateway/proto/auth"
    articlepb "github.com/XRS0/blog/services/api-gateway/proto/article"
    statspb "github.com/XRS0/blog/services/api-gateway/proto/stats"
)
```

## .gitignore

Сгенерированные proto файлы добавлены в `.gitignore`:

```gitignore
# Generated proto files in services
services/auth-service/proto/
services/article-service/proto/
services/stats-service/proto/
services/api-gateway/proto/
```

Каждый разработчик должен сгенерировать proto файлы после клонирования репозитория:

```bash
git clone <repo-url>
cd blog
make proto
```

## Обновление Proto Схем

Если вы изменили `.proto` файлы в директории `proto/`:

1. Запустите генерацию:
   ```bash
   make proto
   ```

2. Проверьте что компиляция проходит успешно:
   ```bash
   # В каждом сервисе
   cd services/auth-service && go build ./...
   cd services/article-service && go build ./...
   cd services/stats-service && go build ./...
   cd services/api-gateway && go build ./...
   ```

3. Закоммитьте изменения только `.proto` файлов (не сгенерированные):
   ```bash
   git add proto/*.proto
   git commit -m "Update proto schemas"
   ```

## Troubleshooting

### Ошибка "protoc: command not found"
- Установите protoc согласно инструкциям выше
- Убедитесь что protoc в вашем PATH

### Ошибка "protoc-gen-go: program not found"
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Ошибка импорта "no required module provides package"
- Убедитесь что вы запустили `make proto`
- Проверьте что директории `proto/` созданы в каждом сервисе
- Запустите `go mod tidy` в каждом сервисе

### Несовместимость версий protoc
- Минимальная версия: protoc 3.12.0
- Рекомендуется: protoc 25.x или новее
- Обновите до последней версии: `brew upgrade protobuf`
