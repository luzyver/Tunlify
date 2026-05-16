package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luzyver/tunlify/internal/config"
	"github.com/luzyver/tunlify/internal/db"
	"github.com/luzyver/tunlify/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	cfg   *config.Config
	db    *db.DB
	audit *service.AuditLogger
}

func NewAuth(cfg *config.Config, db *db.DB, audit *service.AuditLogger) *Auth {
	hash, _ := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), 12)
	db.Exec(`INSERT OR IGNORE INTO users (username, password_hash) VALUES (?, ?)`, cfg.AdminUsername, string(hash))
	return &Auth{cfg: cfg, db: db, audit: audit}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	var id int
	var hash string
	err := h.db.QueryRow(`SELECT id, password_hash FROM users WHERE username = ?`, req.Username).Scan(&id, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	access, _ := h.generateToken(id, req.Username, 15*time.Minute)
	refresh, _ := h.generateToken(id, req.Username, 7*24*time.Hour)

	h.audit.Log(id, "login", "", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    900,
	})
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)
	h.audit.Log(userID, "logout", "", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func (h *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(int)
	username, _ := r.Context().Value("username").(string)

	access, _ := h.generateToken(userID, username, 15*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse{
		AccessToken: access,
		ExpiresIn:   900,
	})
}

func (h *Auth) generateToken(userID int, username string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
