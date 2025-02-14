package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"weak"
)

type SSEManager struct {
	clients map[string]weak.Pointer[chan string]
	mu      sync.Mutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string]weak.Pointer[chan string]),
	}
}

func (m *SSEManager) AddClient(id string, ch *chan string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[id] = weak.Make(ch)
}

func (m *SSEManager) CleanClients() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, ptr := range m.clients {
		if ptr.Value() == nil {
			fmt.Println("Connection is closed and deleted:", id)
			delete(m.clients, id)
		}
	}
}

func (m *SSEManager) Broadcast(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, ptr := range m.clients {
		if ch := ptr.Value(); ch != nil {
			select {
			case *ch <- msg:
			default:
				fmt.Println("Client is busy:", id)
			}
		} else {
			delete(m.clients, id)
		}
	}
}

var sseManager = NewSSEManager()

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming is not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := r.URL.Query().Get("id")
	ch := make(chan string, 1)

	sseManager.AddClient(id, &ch)

	for msg := range ch {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

func main() {
	http.HandleFunc("/events", enableCORS(sseHandler))

	go func() {
		for {
			time.Sleep(5 * time.Second)
			sseManager.Broadcast(fmt.Sprintf("Update: %v", time.Now()))
			sseManager.CleanClients()
		}
	}()

	fmt.Println("SSE server is running...")
	http.ListenAndServe(":8080", nil)
}
