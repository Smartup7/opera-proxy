#!/bin/sh

# Скрипт установки Opera Proxy Checker на Keenetic router с Entware

set -e

# Проверяем, установлен ли Entware
if [ ! -d "/opt" ]; then
    echo "Ошибка: Entware не найден. Убедитесь, что Entware установлен на вашем маршрутизаторе."
    exit 1
fi

# Добавляем /opt/bin в PATH
export PATH="/opt/bin:/opt/sbin:/opt/usr/bin:$PATH"

echo "Установка зависимостей..."

# Устанавливаем необходимые пакеты
opkg update
opkg install wget

# Создаем директорию для установки
INSTALL_DIR="/opt/opera-proxy-checker"
mkdir -p $INSTALL_DIR

cd $INSTALL_DIR

# Скачиваем предварительно собранный бинарный файл checker
echo "Скачивание Opera Proxy Checker..."
wget -O opera-proxy-checker.linux-arm64 https://github.com/yourusername/opera-proxy-checker/releases/download/v1.0/opera-proxy-checker.linux-arm64
chmod +x opera-proxy-checker.linux-arm64

# Скачиваем бинарный файл opera-proxy если он не установлен
if [ ! -f "/opt/opera-proxy" ]; then
    echo "Скачивание Opera Proxy..."
    wget -O opera-proxy.linux-arm64 https://github.com/Snawoot/opera-proxy/releases/latest/download/opera-proxy.linux-arm64
    chmod +x opera-proxy.linux-arm64
    cp opera-proxy.linux-arm64 /opt/opera-proxy
else
    echo "Opera Proxy уже установлен в /opt/opera-proxy"
fi

# Создаем скрипт запуска
cat > start.sh << 'EOF'
#!/bin/sh
cd /opt/opera-proxy-checker
nohup ./opera-proxy-checker.linux-arm64 -binary /opt/opera-proxy -test-url http://httpbin.org/ip -interval 5m -web-port 8080 > checker.log 2>&1 &
echo $! > checker.pid
echo "Opera Proxy Checker запущен. PID: $(cat checker.pid)"
EOF

chmod +x start.sh

# Создаем скрипт остановки
cat > stop.sh << 'EOF'
#!/bin/sh
if [ -f /opt/opera-proxy-checker/checker.pid ]; then
    PID=$(cat /opt/opera-proxy-checker/checker.pid)
    kill $PID
    rm /opt/opera-proxy-checker/checker.pid
    echo "Opera Proxy Checker остановлен"
else
    echo "Файл PID не найден"
fi
EOF

chmod +x stop.sh

# Добавляем в автозапуск
cat > /opt/etc/init.d/S99opera-proxy-checker << 'EOF'
#!/bin/sh

case "$1" in
    start)
        /opt/opera-proxy-checker/start.sh
        ;;
    stop)
        /opt/opera-proxy-checker/stop.sh
        ;;
    restart)
        $0 stop
        sleep 2
        $0 start
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
        ;;
esac
EOF

chmod +x /opt/etc/init.d/S99opera-proxy-checker

echo "Установка завершена!"
echo "Для запуска выполните: /opt/opera-proxy-checker/start.sh"
echo "Для остановки выполните: /opt/opera-proxy-checker/stop.sh"
echo "Веб-интерфейс будет доступен по адресу: http://ROUTER_IP:8080"
echo "Для добавления в автозапуск выполните: ln -s /opt/etc/init.d/S99opera-proxy-checker /opt/etc/rc.d/S99opera-proxy-checker"