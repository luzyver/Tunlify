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
		writeError(w, http.StatusInternalServerError, "failed to read config")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

func (h *Config) Update(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.mgr.Validate(req.Content); err != nil {
		writeError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	userID, _ := r.Context().Value("user_id").(int)
	if err := h.mgr.Write(req.Content, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to write config")
		return
	}

	h.audit.Log(userID, "config_edit", "", r.RemoteAddr)

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Config) Validate(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req struct {
		Content string `json:"content"`
	}
	json.Unmarshal(body, &req)

	if err := h.mgr.Validate(req.Content); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"valid": "false", "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"valid": "true"})
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
		writeError(w, http.StatusInternalServerError, "failed to list backups")
		return
	}

	writeJSON(w, http.StatusOK, backups)
}

func (h *Config) GetBackup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	backup, err := h.mgr.GetBackup(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "backup not found")
		return
	}

	writeJSON(w, http.StatusOK, backup)
}
