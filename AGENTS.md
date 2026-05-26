# AGENTS.md

## Project Overview

Tunlify is a self-hosted web dashboard to manage and monitor Cloudflare Tunnel (cloudflared) and Docker Compose projects. It provides tunnel control, realtime logs, config editing, health checks, metrics, notifications, and project deployment through a modern web UI.

## Architecture

- **Backend**: Go 1.22, Chi router, gorilla/websocket, SQLite (modernc.org/sqlite), JWT auth
- **Frontend**: Vue 3, Vite, TypeScript, Tailwind CSS, Pinia, Vue Router
- **Deploy**: Docker Compose with 3 services (cloudflared, backend, frontend/nginx)
- **Network**: All services communicate via external Docker network `cloudflared`

## Project Structure

```
backend/
├── cmd/server/          # Entrypoint (main.go)
└── internal/
    ├── config/          # Env-based config (envconfig)
    ├── db/              # SQLite initialization
    ├── handler/         # HTTP handlers (auth, config, control, logs, etc.)
    ├── middleware/      # Auth JWT middleware, rate limiting
    ├── service/         # Business logic (cloudflared Docker control, config manager, audit)
    └── ws/              # WebSocket hub for realtime log streaming

frontend/
├── src/
│   ├── views/           # Page components (Dashboard, Logs, Config, etc.)
│   ├── components/      # Shared components (AppSidebar)
│   ├── composables/     # useApi composable for HTTP calls
│   ├── stores/          # Pinia stores (auth)
│   ├── router/          # Vue Router config
│   └── assets/          # Tailwind CSS
├── nginx.conf           # Production reverse proxy config
└── Dockerfile           # Multi-stage build (node → nginx)
```

## Key Conventions

- Backend uses `internal/` package pattern — no exported packages
- Handlers receive services via struct injection, not globals
- Config loaded from environment variables via `envconfig`
- Frontend API calls go through `composables/useApi.ts` with JWT token injection
- Auth uses JWT access tokens with refresh flow
- All state-changing actions are logged to audit trail (SQLite)
- YAML config changes create automatic backups before write

## Development Commands

```bash
# Backend
cd backend && go mod tidy && go run ./cmd/server

# Frontend
cd frontend && npm install && npm run dev

# Docker (production)
docker compose up -d
```

## Environment Variables

Required: `JWT_SECRET`, `ADMIN_PASSWORD`
Optional with defaults: `ADMIN_USERNAME` (admin), `CLOUDFLARED_CONFIG_PATH` (/etc/cloudflared/config.yml), `CLOUDFLARED_CONTAINER` (tunlify-cloudflared), `DB_PATH` (/app/data/app.db), `LISTEN_ADDR` (0.0.0.0:3333)

## API Pattern

- Auth endpoints: `/auth/login`, `/auth/logout`, `/auth/refresh`
- Protected endpoints: `/api/*` (require JWT via Authorization header)
- WebSocket: `/api/logs/ws` for realtime log streaming
- Backend proxies Docker API to control cloudflared container (start/stop/restart)
- Projects: CRUD + compose actions (up/down/restart/deploy with git checkout+pull), SSE streaming output, per-project action history

## Important Notes

- Backend communicates with Docker daemon via mounted `/var/run/docker.sock`
- Frontend in production is served by nginx which reverse-proxies `/api` and `/auth` to backend:3333
- SQLite database stored at `/app/data/app.db` (persisted via Docker volume)
- Cloudflared config mounted read-only to tunnel container, read-write to backend for editing
- No ORM — raw SQL with modernc.org/sqlite (pure Go, no CGO)
- Projects store git credentials (repo_url, username, token) per project in SQLite
- Deploy flow: git fetch → checkout ref → git pull → docker compose up -d --build
- Git tokens are masked in API responses (only last 4 chars shown)
