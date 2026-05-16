package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/luzyver/tunlify/internal/config"
	"github.com/luzyver/tunlify/internal/ws"
)

type ProcessState struct {
	Running     bool   `json:"running"`
	StartTime   string `json:"start_time,omitempty"`
	Uptime      string `json:"uptime,omitempty"`
	Status      string `json:"status"`
	Version     string `json:"version,omitempty"`
	MemoryUsage string `json:"memory_usage,omitempty"`
}

type Cloudflared struct {
	cfg          *config.Config
	hub          *ws.Hub
	mu           sync.RWMutex
	cachedVer    string
	cachedMem    string
	lastMemCheck time.Time
}

func NewCloudflared(cfg *config.Config, hub *ws.Hub) *Cloudflared {
	return &Cloudflared{cfg: cfg, hub: hub}
}

func (s *Cloudflared) Status() ProcessState {
	out, err := exec.Command("docker", "inspect", "-f", "{{.State.Status}} {{.State.StartedAt}}", s.cfg.CloudflaredContainer).Output()
	if err != nil {
		return ProcessState{Running: false, Status: "not found"}
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) < 2 {
		return ProcessState{Running: false, Status: "unknown"}
	}

	status := parts[0]
	running := status == "running"

	var uptime string
	if running {
		if t, err := time.Parse(time.RFC3339Nano, parts[1]); err == nil {
			uptime = time.Since(t).Truncate(time.Second).String()
		}
	}

	if s.cachedVer == "" {
		s.cachedVer = s.getVersion()
	}

	if time.Since(s.lastMemCheck) > 30*time.Second {
		go func() {
			s.mu.Lock()
			s.cachedMem = s.getMemory()
			s.lastMemCheck = time.Now()
			s.mu.Unlock()
		}()
	}

	s.mu.RLock()
	mem := s.cachedMem
	s.mu.RUnlock()

	return ProcessState{
		Running:     running,
		StartTime:   parts[1],
		Uptime:      uptime,
		Status:      status,
		Version:     s.cachedVer,
		MemoryUsage: mem,
	}
}

func (s *Cloudflared) getVersion() string {
	out, err := exec.Command("docker", "exec", s.cfg.CloudflaredContainer, "cloudflared", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (s *Cloudflared) getMemory() string {
	out, err := exec.Command("docker", "stats", "--no-stream", "--format", "{{.MemUsage}}", s.cfg.CloudflaredContainer).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (s *Cloudflared) Start() error {
	return s.dockerCmd("start")
}

func (s *Cloudflared) Stop() error {
	return s.dockerCmd("stop")
}

func (s *Cloudflared) Restart() error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		s.dockerCmd("restart")
	}()
	return nil
}

func (s *Cloudflared) dockerCmd(action string) error {
	cmd := exec.Command("docker", action, s.cfg.CloudflaredContainer)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s failed: %s", action, strings.TrimSpace(string(out)))
	}

	entry, _ := json.Marshal(map[string]string{
		"time":    time.Now().Format(time.RFC3339),
		"message": fmt.Sprintf("[tunlify] cloudflared %s", action),
	})
	s.hub.Broadcast <- entry

	return nil
}

func (s *Cloudflared) Logs() error {
	cmd := exec.Command("docker", "logs", "-f", "--tail", "100", s.cfg.CloudflaredContainer)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	go s.streamToHub(stdout)
	go s.streamToHub(stderr)
	go cmd.Wait()

	return nil
}

func (s *Cloudflared) streamToHub(r interface{ Read([]byte) (int, error) }) {
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			entry, _ := json.Marshal(map[string]string{
				"time":    time.Now().Format(time.RFC3339),
				"message": strings.TrimSpace(string(buf[:n])),
			})
			s.hub.Broadcast <- entry
		}
		if err != nil {
			return
		}
	}
}

func (s *Cloudflared) StartTcpAccess(hostname, localURL string) (int, error) {
	return 0, errors.New("tcp access not supported in docker mode, use ingress config instead")
}

func (s *Cloudflared) StopTcpAccess(pid int) error {
	return errors.New("tcp access not supported in docker mode")
}
