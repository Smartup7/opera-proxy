#!/bin/bash

# Скрипт для сборки Opera Proxy Checker под различные архитектуры

set -e

echo "Сборка Opera Proxy Checker..."

# Создаем директорию для бинарных файлов
mkdir -p bin

# Собираем для Linux ARM64 (Keenetic/Entware)
echo "Сборка для Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -a -tags netgo -trimpath -ldflags '-s -w -extldflags "-static"' -o bin/opera-proxy-checker.linux-arm64 .

# Собираем для Linux AMD64
echo "Сборка для Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -trimpath -ldflags '-s -w -extldflags "-static"' -o bin/opera-proxy-checker.linux-amd64 .

# Собираем для Linux ARM
echo "Сборка для Linux ARM..."
GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -a -tags netgo -trimpath -ldflags '-s -w -extldflags "-static"' -o bin/opera-proxy-checker.linux-arm .

echo "Сборка завершена. Бинарные файлы находятся в директории bin/"
ls -la bin/