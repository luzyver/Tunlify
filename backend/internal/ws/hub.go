package ws

import "sync"

type Client struct {
	Send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
	buffer     [][]byte
	bufferSize int
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		bufferSize: 1000,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			for _, msg := range h.buffer {
				select {
				case client.Send <- msg:
				default:
				}
			}
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		case msg := <-h.Broadcast:
			h.mu.Lock()
			h.buffer = append(h.buffer, msg)
			if len(h.buffer) > h.bufferSize {
				h.buffer = h.buffer[1:]
			}
			for client := range h.clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) History(limit int) [][]byte {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if limit <= 0 || limit > len(h.buffer) {
		limit = len(h.buffer)
	}
	start := len(h.buffer) - limit
	result := make([][]byte, limit)
	copy(result, h.buffer[start:])
	return result
}
