#!/bin/bash
set -e

echo "🚀 Настройка проекта Blog Microservices..."
echo ""

# Проверка protoc
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc не найден"
    echo "📦 Установка protoc..."
    if [ -f "./install-protoc.sh" ]; then
        ./install-protoc.sh
    else
        echo "⚠️  Скрипт install-protoc.sh не найден"
        echo "Установите protoc вручную:"
        echo "  macOS: brew install protobuf"
        echo "  Linux: sudo apt install -y protobuf-compiler"
        exit 1
    fi
else
    echo "✅ protoc найден: $(protoc --version)"
fi

# Проверка Go плагинов
GOPATH=$(go env GOPATH)
if [ ! -f "$GOPATH/bin/protoc-gen-go" ] || [ ! -f "$GOPATH/bin/protoc-gen-go-grpc" ]; then
    echo "📦 Установка Go плагинов для protoc..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    echo "✅ Go плагины установлены"
else
    echo "✅ Go плагины найдены"
fi

# Проверка PATH
if ! command -v protoc-gen-go &> /dev/null; then
    echo "⚠️  protoc-gen-go не в PATH"
    echo "Добавьте в ~/.zshrc или ~/.bashrc:"
    echo "export PATH=\"\$PATH:\$(go env GOPATH)/bin\""
    echo ""
    echo "Временное добавление в PATH для текущей сессии..."
    export PATH="$PATH:$GOPATH/bin"
fi

echo ""
echo "🔨 Очистка старых proto файлов..."
make clean

echo ""
echo "🔨 Генерация proto файлов..."
make proto

# Проверка структуры
echo ""
echo "🔍 Проверка структуры сгенерированных файлов..."
EXPECTED_FILES=12
ACTUAL_FILES=$(find services -name "*.pb.go" -o -name "*_grpc.pb.go" | wc -l | tr -d ' ')

if [ "$ACTUAL_FILES" -eq "$EXPECTED_FILES" ]; then
    echo "✅ Сгенерировано $ACTUAL_FILES файлов (ожидалось $EXPECTED_FILES)"
else
    echo "⚠️  Сгенерировано $ACTUAL_FILES файлов (ожидалось $EXPECTED_FILES)"
fi

# Проверка на лишние папки proto/
if find services -path "*/proto/proto" | grep -q .; then
    echo "❌ Найдены лишние папки proto/ внутри proto/"
    echo "Проверьте структуру:"
    find services -path "*/proto/proto"
else
    echo "✅ Структура папок правильная (нет лишних proto/)"
fi

echo ""
echo "📚 Обновление зависимостей в сервисах..."
for service in auth-service article-service stats-service api-gateway; do
    echo "  - $service"
    (cd "services/$service" && go mod tidy)
done

echo ""
echo "✅ Настройка завершена!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Запустите сервисы: docker-compose up --build"
echo "2. Проверьте здоровье: curl http://localhost:8080/health"
echo "3. Попробуйте API: см. QUICKSTART.md"
echo ""
echo "📚 Документация:"
echo "  - README.md - Общая информация"
echo "  - QUICKSTART.md - Быстрый старт"
echo "  - PROTO_CHEATSHEET.md - Шпаргалка по proto"
echo "  - PROTO_STRUCTURE_EXPLAINED.md - Объяснение структуры"
