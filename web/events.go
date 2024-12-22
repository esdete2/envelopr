package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/networkteam/slogutils"
)

type EventBroker struct {
	clients map[chan string]bool
	mu      sync.RWMutex
}

func NewEventBroker() *EventBroker {
	return &EventBroker{
		clients: make(map[chan string]bool),
	}
}

func (b *EventBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new channel for this client
	messageChan := make(chan string)

	// Register client
	b.mu.Lock()
	b.clients[messageChan] = true
	b.mu.Unlock()

	slog.Debug("client connected to event stream")

	// Clean up when client disconnects
	defer func() {
		b.mu.Lock()
		delete(b.clients, messageChan)
		close(messageChan)
		b.mu.Unlock()
		slog.Debug("client disconnected from event stream")
	}()

	// Keep connection open
	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-messageChan:
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				slog.Error("writing event to client", slogutils.Err(err))
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			} else {
				slog.Error("response writer does not support flushing")
				return
			}
		}
	}
}

func (b *EventBroker) Notify(message string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	slog.With("message", message).
		With("clients", len(b.clients)).
		Debug("broadcasting event")

	for clientChan := range b.clients {
		select {
		case clientChan <- message:
		default:
			// Skip if client can't keep up
		}
	}
}
