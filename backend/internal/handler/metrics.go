package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/luzyver/tunlify/internal/config"
)

type Metrics struct {
	cfg *config.Config
}

func NewMetrics(cfg *config.Config) *Metrics {
	return &Metrics{cfg: cfg}
}

type metricsResponse struct {
	TotalRequests  int64              `json:"total_requests"`
	ActiveConns    int64              `json:"active_connections"`
	ResponseCodes  map[string]int64   `json:"response_codes"`
	Raw            string             `json:"raw,omitempty"`
}

func (h *Metrics) Get(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(h.cfg.MetricsAddr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metricsResponse{})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	raw := string(body)

	m := metricsResponse{ResponseCodes: make(map[string]int64)}
	m.Raw = raw

	re := regexp.MustCompile(`cloudflared_tunnel_total_requests\s+(\d+)`)
	if matches := re.FindStringSubmatch(raw); len(matches) > 1 {
		m.TotalRequests, _ = strconv.ParseInt(matches[1], 10, 64)
	}

	reConn := regexp.MustCompile(`cloudflared_tunnel_active_streams\s+(\d+)`)
	if matches := reConn.FindStringSubmatch(raw); len(matches) > 1 {
		m.ActiveConns, _ = strconv.ParseInt(matches[1], 10, 64)
	}

	reStatus := regexp.MustCompile(`cloudflared_tunnel_response_by_code{status_code="(\d+)"}\s+(\d+)`)
	for _, line := range strings.Split(raw, "\n") {
		if matches := reStatus.FindStringSubmatch(line); len(matches) > 2 {
			count, _ := strconv.ParseInt(matches[2], 10, 64)
			m.ResponseCodes[matches[1]] = count
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
