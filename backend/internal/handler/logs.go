package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/luzyver/tunlify/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Logs struct {
	hub    *ws.Hub
	secret string
}

func NewLogs(hub *ws.Hub, secret string) *Logs {
	return &Logs{hub: hub, secret: secret}
}

func (h *Logs) Stream(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.secret), nil
	})
	if err != nil {
		http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	client := &ws.Client{Send: make(chan []byte, 256)}
	h.hub.Register <- client
	defer func() { h.hub.Unregister <- client }()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				conn.WriteMessage(websocket.TextMessage, msg)
			case <-ticker.C:
				conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Logs) History(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	entries := h.hub.History(limit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  entries,
		"count": len(entries),
	})
}
