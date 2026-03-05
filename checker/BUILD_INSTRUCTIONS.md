# Инструкции по сборке Opera Proxy Checker

Данный документ описывает процесс сборки бинарных файлов для Opera Proxy Checker.

## Требования

- Go 1.21 или выше
- GNU Make
- Git

## Подготовка окружения

1. Установите Go:
```bash
# На Ubuntu/Debian
sudo apt-get install golang

# На других системах см. https://golang.org/dl/
```

2. Установите зависимости:
```bash
go mod tidy
```

## Сборка для ARM64 (Keenetic/Entware)

Для сборки бинарного файла для архитектуры ARM64 выполните:

```bash
GOOS=linux GOARCH=arm64 go build -o bin/opera-proxy-checker.linux-arm64 .
```

Или используйте Makefile:

```bash
make bin-linux-arm64
```

## Сборка для других архитектур

Вы можете собрать бинарные файлы для различных архитектур:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o bin/opera-proxy-checker.linux-amd64 .

# Linux ARM
GOOS=linux GOARCH=arm go build -o bin/opera-proxy-checker.linux-arm .

# FreeBSD AMD64
GOOS=freebsd GOARCH=amd64 go build -o bin/opera-proxy-checker.freebsd-amd64 .
```

## Сборка всех версий

Для сборки всех поддерживаемых версий используйте:

```bash
make all
```

## Создание релизного архива

Для создания архива с бинарными файлами для релиза:

```bash
mkdir -p release
cp bin/* release/
tar -czvf opera-proxy-checker-release.tar.gz -C release .
```

## Проверка сборки

Для проверки корректности сборки можно запустить:

```bash
# Запуск без параметров покажет справку
./bin/opera-proxy-checker.linux-arm64

# Запуск с тестовым бинарным файлом opera-proxy
./bin/opera-proxy-checker.linux-arm64 -binary /path/to/opera-proxy -test-url http://httpbin.org/ip -interval 1m -web-port 8080
```

## Особенности сборки для Entware

При сборке для Entware важно использовать статическую линковку:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags netgo -trimpath -ldflags '-s -w -extldflags "-static"' -o bin/opera-proxy-checker.linux-arm64 .
```

Это обеспечивает совместимость с различными версиями libc в embedded системах.