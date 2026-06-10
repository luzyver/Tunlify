package handler

import (
	"net/http"
	"strconv"

	"github.com/luzyver/tunlify/internal/service"
)

type Audit struct {
	audit *service.AuditLogger
}

func NewAudit(audit *service.AuditLogger) *Audit {
	return &Audit{audit: audit}
}

func (h *Audit) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := 50, 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}
	action := r.URL.Query().Get("action")

	entries, total, err := h.audit.List(limit, offset, action)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch audit logs")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"entries": entries,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}
