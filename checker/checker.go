package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Result структура для хранения результата проверки
type Result struct {
	Timestamp    time.Time `json:"timestamp"`
	Status       string    `json:"status"`
	ResponseTime string    `json:"response_time"`
	Error        string    `json:"error,omitempty"`
	PublicIP     string    `json:"public_ip,omitempty"`
}

// Checker структура для управления проверкой
type Checker struct {
	binaryPath string
	testURL    string
	interval   time.Duration
	results    chan Result
}

// NewChecker создает новый экземпляр чекера
func NewChecker(binaryPath, testURL string, interval time.Duration) *Checker {
	return &Checker{
		binaryPath: binaryPath,
		testURL:    testURL,
		interval:   interval,
		results:    make(chan Result, 100),
	}
}

// Start запускает процесс проверки
func (c *Checker) Start() {
	go c.run()
}

// run внутренний метод для выполнения проверок
func (c *Checker) run() {
	for {
		result := c.check()
		c.results <- result
		time.Sleep(c.interval)
	}
}

// check выполняет проверку работы бинарного файла
func (c *Checker) check() Result {
	startTime := time.Now()

	// Запускаем opera-proxy в фоновом режиме
	cmd := exec.Command(c.binaryPath, "-country", "EU", "-bind-address", "127.0.0.1:18080")
	
	// Запускаем команду
	err := cmd.Start()
	if err != nil {
		return Result{
			Timestamp:  time.Now(),
			Status:     "error",
			Error:      fmt.Sprintf("Failed to start opera-proxy: %v", err),
			PublicIP:   "",
			ResponseTime: time.Since(startTime).String(),
		}
	}

	// Ждем немного, чтобы прокси успел запуститься
	time.Sleep(5 * time.Second)

	// Проверяем доступность прокси
	proxyURL := "http://127.0.0.1:18080"
	proxy := func(r *http.Request) (*url.URL, error) {
		return url.Parse(proxyURL)
	}
	transport := &http.Transport{
		Proxy: proxy,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// Выполняем тестовый запрос для получения публичного IP
	publicIP, err := getPublicIP(client, c.testURL)
	responseTime := time.Since(startTime)

	// Останавливаем opera-proxy
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	if err != nil {
		return Result{
			Timestamp:    time.Now(),
			Status:       "error",
			Error:        fmt.Sprintf("Failed to get public IP: %v", err),
			PublicIP:     "",
			ResponseTime: responseTime.String(),
		}
	}

	return Result{
		Timestamp:    time.Now(),
		Status:       "success",
		Error:        "",
		PublicIP:     publicIP,
		ResponseTime: responseTime.String(),
	}
}

// getPublicIP получает публичный IP через прокси
func getPublicIP(client *http.Client, testURL string) (string, error) {
	resp, err := client.Get(testURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Предполагаем, что сервис возвращает JSON с IP адресом
	// Пример: {"origin": "1.2.3.4"}
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "origin") || strings.Contains(line, "ip") || strings.Contains(line, "IP") {
			// Извлекаем IP из JSON или текста
			ip := extractIP(line)
			if ip != "" {
				return ip, nil
			}
		}
	}

	// Если не удалось извлечь IP из JSON, возвращаем весь ответ
	return strings.TrimSpace(string(body)), nil
}

// extractIP извлекает IP-адрес из строки
func extractIP(text string) string {
	// Простой регекс для поиска IP-адреса
	// Это базовая реализация, можно улучшить
	parts := strings.Split(text, ":")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "\",}")
		if strings.Count(part, ".") == 3 {
			// Проверим, является ли это IP-адресом
			ipParts := strings.Split(part, ".")
			if len(ipParts) == 4 {
				allNumeric := true
				for _, ipPart := range ipParts {
					ipPart = strings.TrimSpace(ipPart)
					ipPart = strings.TrimLeft(ipPart, "\"{ ")
					ipPart = strings.TrimRight(ipPart, "} \",")
					if !isNumeric(ipPart) {
						allNumeric = false
						break
					}
					num, _ := strconv.Atoi(ipPart)
					if num < 0 || num > 255 {
						allNumeric = false
						break
					}
				}
				if allNumeric {
					return part
				}
			}
		}
	}
	return ""
}

// isNumeric проверяет, является ли строка числом
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// GetResults возвращает канал с результатами
func (c *Checker) GetResults() <-chan Result {
	return c.results
}