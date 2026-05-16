package handler

import (
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifConfig{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Notifications) Update(w http.ResponseWriter, r *http.Request) {
	var cfg notifConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.MkdirAll(h.dataPath, 0755)
	os.WriteFile(h.configFile(), data, 0644)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Notifications) Test(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.configFile())
	if err != nil {
		http.Error(w, `{"error":"no config"}`, http.StatusBadRequest)
		return
	}

	var cfg notifConfig
	json.Unmarshal(data, &cfg)

	if cfg.WebhookURL == "" {
		http.Error(w, `{"error":"no webhook url"}`, http.StatusBadRequest)
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

	resp, err := http.Post(cfg.WebhookURL, "application/json", jsonReader(payload))
	if err != nil || resp.StatusCode >= 400 {
		http.Error(w, `{"error":"webhook failed"}`, http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func jsonReader(data []byte) *jsonBody { return &jsonBody{data: data} }

type jsonBody struct{ data []byte; i int }

func (b *jsonBody) Read(p []byte) (int, error) {
	if b.i >= len(b.data) {
		return 0, nil
	}
	n := copy(p, b.data[b.i:])
	b.i += n
	if b.i >= len(b.data) {
		return n, nil
	}
	return n, nil
}
