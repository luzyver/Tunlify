package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type Notifications struct {
	dataPath string
}

func NewNotifications(dataPath string) *Notifications {
	return &Notifications{dataPath: dataPath}
}

type notifConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	Type       string `json:"type"`
}

func (h *Notifications) configFile() string {
	return filepath.Join(h.dataPath, "notifications.json")
}

func (h *Notifications) Get(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.configFile())
	if err != nil {
		writeJSON(w, http.StatusOK, notifConfig{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Notifications) Update(w http.ResponseWriter, r *http.Request) {
	var cfg notifConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.MkdirAll(h.dataPath, 0755); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save config")
		return
	}
	if err := os.WriteFile(h.configFile(), data, 0644); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save config")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Notifications) Test(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.configFile())
	if err != nil {
		writeError(w, http.StatusBadRequest, "no config")
		return
	}

	var cfg notifConfig
	json.Unmarshal(data, &cfg)

	if cfg.WebhookURL == "" {
		writeError(w, http.StatusBadRequest, "no webhook url")
		return
	}

	var payload []byte
	switch cfg.Type {
	case "discord":
		payload, _ = json.Marshal(map[string]string{"content": "🟢 Tunlify test notification"})
	case "telegram":
		payload, _ = json.Marshal(map[string]interface{}{"text": "🟢 Tunlify test notification", "parse_mode": "HTML"})
	default:
		payload, _ = json.Marshal(map[string]string{"text": "🟢 Tunlify test notification"})
	}

	resp, err := http.Post(cfg.WebhookURL, "application/json", bytes.NewReader(payload))
	if err != nil || resp.StatusCode >= 400 {
		writeError(w, http.StatusBadGateway, "webhook failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}
