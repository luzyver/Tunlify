package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luzyver/tunlify/internal/service"
)

type Config struct {
	mgr   *service.ConfigManager
	audit *service.AuditLogger
}

func NewConfig(mgr *service.ConfigManager, audit *service.AuditLogger) *Config {
	return &Config{mgr: mgr, audit: audit}
}

func (h *Config) Get(w http.ResponseWriter, r *http.Request) {
	content, err := h.mgr.Read()
	if err != nil {
		http.Error(w, `{"error":"failed to read config"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"content": content})
}

func (h *Config) Update(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.mgr.Validate(req.Content); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "validation failed: " + err.Error()})
		return
	}

	userID, _ := r.Context().Value("user_id").(int)
	if err := h.mgr.Write(req.Content, userID); err != nil {
		http.Error(w, `{"error":"failed to write config"}`, http.StatusInternalServerError)
		return
	}

	h.audit.Log(userID, "config_edit", "", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Config) Validate(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req struct {
		Content string `json:"content"`
	}
	json.Unmarshal(body, &req)

	if err := h.mgr.Validate(req.Content); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"valid": "false", "error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"valid": "true"})
}

func (h *Config) ListBackups(w http.ResponseWriter, r *http.Request) {
	limit, offset := 20, 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	backups, err := h.mgr.ListBackups(limit, offset)
	if err != nil {
		http.Error(w, `{"error":"failed to list backups"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}

func (h *Config) GetBackup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	backup, err := h.mgr.GetBackup(id)
	if err != nil {
		http.Error(w, `{"error":"backup not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backup)
}
