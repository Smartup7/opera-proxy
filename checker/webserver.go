package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

// WebServer структура для веб-сервера
type WebServer struct {
	checker  *Checker
	results  []Result
	mu       sync.RWMutex
	port     string
}

// NewWebServer создает новый веб-сервер
func NewWebServer(checker *Checker, port string) *WebServer {
	ws := &WebServer{
		checker: checker,
		results: make([]Result, 0),
		port:    port,
	}
	
	// Запускаем горутину для получения результатов
	go ws.collectResults()
	
	return ws
}

// collectResults собирает результаты из канала
func (ws *WebServer) collectResults() {
	for result := range ws.checker.GetResults() {
		ws.mu.Lock()
		ws.results = append(ws.results, result)
		
		// Ограничиваем количество результатов для экономии памяти
		if len(ws.results) > 100 {
			ws.results = ws.results[len(ws.results)-100:]
		}
		
		ws.mu.Unlock()
	}
}

// Start запускает веб-сервер
func (ws *WebServer) Start() {
	http.HandleFunc("/", ws.indexHandler)
	http.HandleFunc("/results", ws.resultsHandler)
	http.HandleFunc("/api/results", ws.apiResultsHandler)
	
	log.Printf("Web server starting on port %s", ws.port)
	log.Fatal(http.ListenAndServe(":"+ws.port, nil))
}

// indexHandler обрабатывает главную страницу
func (ws *WebServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	ws.mu.RLock()
	results := ws.results
	ws.mu.RUnlock()
	
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Opera Proxy Checker Results</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .success { color: green; }
        .error { color: red; }
        .timestamp { font-size: 0.9em; color: #666; }
    </style>
</head>
<body>
    <h1>Opera Proxy Checker Results</h1>
    <table>
        <thead>
            <tr>
                <th>Timestamp</th>
                <th>Status</th>
                <th>Response Time</th>
                <th>Public IP</th>
                <th>Error</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td class="timestamp">{{.Timestamp.Format "2006-01-02 15:04:05"}}</td>
                <td class="{{.Status}}">{{.Status}}</td>
                <td>{{.ResponseTime}}</td>
                <td>{{.PublicIP}}</td>
                <td>{{.Error}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
	
	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	err = t.Execute(w, results)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// resultsHandler возвращает HTML с последними результатами
func (ws *WebServer) resultsHandler(w http.ResponseWriter, r *http.Request) {
	ws.indexHandler(w, r)
}

// apiResultsHandler возвращает JSON с результатами
func (ws *WebServer) apiResultsHandler(w http.ResponseWriter, r *http.Request) {
	ws.mu.RLock()
	results := ws.results
	ws.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetResults возвращает последние результаты
func (ws *WebServer) GetResults() []Result {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	
	// Возвращаем копию слайса
	resultsCopy := make([]Result, len(ws.results))
	copy(resultsCopy, ws.results)
	return resultsCopy
}