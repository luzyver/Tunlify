package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type HealthCheck struct {
	configPath string
}

func NewHealthCheck(configPath string) *HealthCheck {
	return &HealthCheck{configPath: configPath}
}

type serviceHealth struct {
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
	Status   string `json:"status"`
	Latency  string `json:"latency"`
}

func (h *HealthCheck) Check(w http.ResponseWriter, r *http.Request) {
	var cfg struct {
		Ingress []struct {
			Hostname string `yaml:"hostname"`
			Service  string `yaml:"service"`
		} `yaml:"ingress"`
	}

	data, err := os.ReadFile(h.configPath)
	if err != nil {
		http.Error(w, `{"error":"cannot read config"}`, http.StatusInternalServerError)
		return
	}
	yaml.Unmarshal(data, &cfg)

	var results []serviceHealth
	for _, ing := range cfg.Ingress {
		if ing.Hostname == "" {
			continue
		}
		s := serviceHealth{Hostname: ing.Hostname, Service: ing.Service}
		start := time.Now()

		u, err := url.Parse(ing.Service)
		if err != nil {
			s.Status = "invalid"
			results = append(results, s)
			continue
		}

		host := u.Host
		if !strings.Contains(host, ":") {
			if u.Scheme == "https" {
				host += ":443"
			} else {
				host += ":80"
			}
		}

		conn, err := net.DialTimeout("tcp", host, 3*time.Second)
		if err != nil {
			s.Status = "down"
			s.Latency = "timeout"
		} else {
			conn.Close()
			s.Status = "up"
			s.Latency = time.Since(start).Truncate(time.Millisecond).String()
		}
		results = append(results, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
