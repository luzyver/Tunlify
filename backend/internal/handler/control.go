package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luzyver/tunlify/internal/service"
)

type Control struct {
	cfd   *service.Cloudflared
	audit *service.AuditLogger
}

func NewControl(cfd *service.Cloudflared, audit *service.AuditLogger) *Control {
	return &Control{cfd: cfd, audit: audit}
}

func (h *Control) Execute(w http.ResponseWriter, r *http.Request) {
	action := chi.URLParam(r, "action")
	userID, _ := r.Context().Value("user_id").(int)

	var err error
	switch action {
	case "restart":
		err = h.cfd.Restart()
	default:
		http.Error(w, `{"error":"invalid action"}`, http.StatusBadRequest)
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.audit.Log(userID, action, "", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action})
}
