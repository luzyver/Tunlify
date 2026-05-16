package service

import (
	"os"

	"github.com/luzyver/tunlify/internal/config"
	"github.com/luzyver/tunlify/internal/db"
	"gopkg.in/yaml.v3"
)

type ConfigBackup struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type ConfigManager struct {
	cfg *config.Config
	db  *db.DB
}

func NewConfigManager(cfg *config.Config, db *db.DB) *ConfigManager {
	return &ConfigManager{cfg: cfg, db: db}
}

func (m *ConfigManager) Read() (string, error) {
	data, err := os.ReadFile(m.cfg.CloudflaredConfig)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *ConfigManager) Write(content string, userID int) error {
	current, err := m.Read()
	if err == nil && current != "" {
		m.db.Exec(`DELETE FROM config_backups`)
		m.db.Exec(`INSERT INTO config_backups (content, user_id) VALUES (?, ?)`, current, userID)
	}

	return os.WriteFile(m.cfg.CloudflaredConfig, []byte(content), 0644)
}

func (m *ConfigManager) Validate(content string) error {
	var out interface{}
	if err := yaml.Unmarshal([]byte(content), &out); err != nil {
		return err
	}
	return nil
}

func (m *ConfigManager) ListBackups(limit, offset int) ([]ConfigBackup, error) {
	rows, err := m.db.Query(`SELECT id, content, user_id, created_at FROM config_backups ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var backups []ConfigBackup
	for rows.Next() {
		var b ConfigBackup
		rows.Scan(&b.ID, &b.Content, &b.UserID, &b.CreatedAt)
		backups = append(backups, b)
	}
	return backups, nil
}

func (m *ConfigManager) GetBackup(id int) (*ConfigBackup, error) {
	var b ConfigBackup
	err := m.db.QueryRow(`SELECT id, content, user_id, created_at FROM config_backups WHERE id = ?`, id).
		Scan(&b.ID, &b.Content, &b.UserID, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
