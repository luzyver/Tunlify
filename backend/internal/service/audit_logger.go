package service

import "github.com/luzyver/tunlify/internal/db"

type AuditEntry struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Detail    string `json:"detail"`
	IPAddress string `json:"ip_address"`
	CreatedAt string `json:"created_at"`
}

type AuditLogger struct {
	db *db.DB
}

func NewAuditLogger(db *db.DB) *AuditLogger {
	return &AuditLogger{db: db}
}

func (a *AuditLogger) Log(userID int, action, detail, ip string) {
	a.db.Exec(`INSERT INTO audit_logs (user_id, action, detail, ip_address) VALUES (?, ?, ?, ?)`,
		userID, action, detail, ip)
	a.db.Exec(`DELETE FROM audit_logs WHERE created_at < datetime('now', '-30 days')`)
}

func (a *AuditLogger) List(limit, offset int, action string) ([]AuditEntry, int, error) {
	var total int
	query := `SELECT COUNT(*) FROM audit_logs`
	args := []interface{}{}

	if action != "" {
		query += ` WHERE action = ?`
		args = append(args, action)
	}
	a.db.QueryRow(query, args...).Scan(&total)

	selectQ := `SELECT id, user_id, action, COALESCE(detail,''), COALESCE(ip_address,''), created_at FROM audit_logs`
	if action != "" {
		selectQ += ` WHERE action = ?`
	}
	selectQ += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := a.db.Query(selectQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		rows.Scan(&e.ID, &e.UserID, &e.Action, &e.Detail, &e.IPAddress, &e.CreatedAt)
		entries = append(entries, e)
	}
	return entries, total, nil
}
