package handler

import (
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
		writeError(w, http.StatusBadRequest, "invalid action")
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.audit.Log(userID, action, "", r.RemoteAddr)

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "action": action})
}
