package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	var (
		binaryPath = flag.String("binary", "./opera-proxy", "Path to opera-proxy binary")
		testURL    = flag.String("test-url", "http://httpbin.org/ip", "URL to test proxy connectivity")
		interval   = flag.Duration("interval", 5*time.Minute, "Interval between checks")
		webPort    = flag.String("web-port", "8080", "Port for web interface")
	)
	
	flag.Parse()
	
	log.Println("Starting Opera Proxy Checker...")
	log.Printf("Binary path: %s", *binaryPath)
	log.Printf("Test URL: %s", *testURL)
	log.Printf("Check interval: %s", *interval)
	log.Printf("Web interface port: %s", *webPort)
	
	// Создаем чекер
	checker := NewChecker(*binaryPath, *testURL, *interval)
	checker.Start()
	
	// Создаем веб-сервер
	webServer := NewWebServer(checker, *webPort)
	
	// Запускаем веб-сервер (эта функция блокирующая)
	webServer.Start()
}