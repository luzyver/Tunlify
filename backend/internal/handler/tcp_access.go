package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luzyver/tunlify/internal/service"
)

type TcpAccess struct {
	cfd *service.Cloudflared
}

func NewTcpAccess(cfd *service.Cloudflared) *TcpAccess {
	return &TcpAccess{cfd: cfd}
}

type tcpRequest struct {
	Hostname string `json:"hostname"`
	LocalURL string `json:"local_url"`
	Mode     string `json:"mode"`
}

func (h *TcpAccess) Generate(w http.ResponseWriter, r *http.Request) {
	var req tcpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	base := fmt.Sprintf("cloudflared access tcp --hostname %s --url %s", req.Hostname, req.LocalURL)

	var command string
	switch req.Mode {
	case "nohup":
		command = "nohup " + base + " &"
	case "systemd":
		command = fmt.Sprintf(`[Unit]
Description=Cloudflared TCP Access - %s
After=network.target

[Service]
ExecStart=/usr/local/bin/%s
Restart=always
User=cloudflared

[Install]
WantedBy=multi-user.target`, req.Hostname, base)
	default:
		command = base
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"command": command})
}

func (h *TcpAccess) Run(w http.ResponseWriter, r *http.Request) {
	var req tcpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	pid, err := h.cfd.StartTcpAccess(req.Hostname, req.LocalURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"pid": pid, "status": "running"})
}

func (h *TcpAccess) Stop(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(chi.URLParam(r, "pid"))
	if err != nil {
		http.Error(w, `{"error":"invalid pid"}`, http.StatusBadRequest)
		return
	}

	if err := h.cfd.StopTcpAccess(pid); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}
