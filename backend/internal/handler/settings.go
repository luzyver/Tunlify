package handler

import (
	"encoding/json"
	"net/http"

	"github.com/luzyver/tunlify/internal/db"
	"github.com/luzyver/tunlify/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type Settings struct {
	db    *db.DB
	audit *service.AuditLogger
}

func NewSettings(db *db.DB, audit *service.AuditLogger) *Settings {
	return &Settings{db: db, audit: audit}
}

func (h *Settings) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userID, _ := r.Context().Value("user_id").(int)

	var hash string
	h.db.QueryRow(`SELECT password_hash FROM users WHERE id = ?`, userID).Scan(&hash)

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.CurrentPassword)) != nil {
		writeError(w, http.StatusUnauthorized, "current password incorrect")
		return
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if _, err := h.db.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, string(newHash), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update password")
		return
	}

	h.audit.Log(userID, "password_change", "", r.RemoteAddr)

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
