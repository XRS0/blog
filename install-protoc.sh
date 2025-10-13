#!/bin/bash

# Скрипт для установки protoc без Homebrew
# Для macOS

set -e

PROTOC_VERSION="25.1"
ARCH=$(uname -m)

if [ "$ARCH" = "arm64" ]; then
    PROTOC_ARCH="osx-aarch_64"
else
    PROTOC_ARCH="osx-x86_64"
fi

echo "🔧 Установка Protocol Buffers Compiler (protoc) v${PROTOC_VERSION}"
echo "Архитектура: ${PROTOC_ARCH}"

# Создаем временную директорию
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Скачиваем protoc
echo "📥 Скачивание protoc..."
curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-${PROTOC_ARCH}.zip"

# Распаковываем
echo "📦 Распаковка..."
unzip -q "protoc-${PROTOC_VERSION}-${PROTOC_ARCH}.zip"

# Устанавливаем в /usr/local
echo "📂 Установка в /usr/local..."
sudo mkdir -p /usr/local/bin
sudo mkdir -p /usr/local/include

sudo cp bin/protoc /usr/local/bin/
sudo cp -r include/* /usr/local/include/

# Проверяем установку
if command -v protoc &> /dev/null; then
    echo "✅ protoc успешно установлен!"
    protoc --version
else
    echo "❌ Ошибка установки protoc"
    exit 1
fi

# Устанавливаем Go плагины
echo ""
echo "🔧 Установка Go плагинов для protoc..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Проверяем Go плагины
GOPATH=$(go env GOPATH)
if [ -f "$GOPATH/bin/protoc-gen-go" ] && [ -f "$GOPATH/bin/protoc-gen-go-grpc" ]; then
    echo "✅ Go плагины успешно установлены!"
else
    echo "❌ Ошибка установки Go плагинов"
    exit 1
fi

# Очистка
cd - > /dev/null
rm -rf "$TMP_DIR"

echo ""
echo "🎉 Установка завершена!"
echo ""
echo "📝 Добавьте в ~/.zshrc или ~/.bashrc:"
echo "export PATH=\"\$PATH:\$(go env GOPATH)/bin\""
echo ""
echo "Затем выполните:"
echo "source ~/.zshrc  # или source ~/.bashrc"
echo ""
echo "После этого можно запускать: make proto"
