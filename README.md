# ☁️ Tunlify

Self-hosted web dashboard to manage and monitor Cloudflare Tunnel (cloudflared).

![Go](https://img.shields.io/badge/backend-Go-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/frontend-Vue_3-4FC08D?logo=vue.js)
![Docker](https://img.shields.io/badge/deploy-Docker-2496ED?logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **Dashboard** — tunnel status, uptime, active hostnames, container PID
- **Control** — start, stop, restart, reload config via Docker API
- **Logs** — realtime streaming via WebSocket with filter & search
- **Config Editor** — edit `config.yml` with YAML validation and auto-backup
- **Health Check** — monitor ingress endpoint availability
- **Metrics** — cloudflared metrics visualization
- **TCP Access** — generate `cloudflared access tcp` commands
- **Notifications** — webhook alerts (Discord, Telegram, Slack)
- **Audit Log** — full history of all actions
- **Auth** — JWT-based login with session management

## Tech Stack

| Layer | Stack |
|-------|-------|
| Backend | Go, Chi, gorilla/websocket, SQLite |
| Frontend | Vue 3, Vite, Tailwind CSS |
| Deploy | Docker Compose |

## Quick Start

### Prerequisites

- Docker & Docker Compose
- A configured Cloudflare Tunnel (`cloudflared`)
- Docker network `cloudflared` created:
  ```bash
  docker network create cloudflared
  ```

### Setup

```bash
git clone https://github.com/luzyver/Tunlify.git
cd Tunlify

# Configure environment
cp .env.example .env
nano .env

# Place your cloudflared config and credentials
mkdir -p cloudflared
# Copy your config.yml and credentials.json into ./cloudflared/

# Start
docker compose up -d
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secret for JWT signing (min 32 chars) | *required* |
| `ADMIN_USERNAME` | Login username | `admin` |
| `ADMIN_PASSWORD` | Login password | *required* |
| `TUNNEL_NAME` | Cloudflare tunnel name | — |
| `TUNNEL_ID` | Cloudflare tunnel UUID | — |
| `CLOUDFLARED_CONFIG_PATH` | Config path inside container | `/etc/cloudflared/config.yml` |
| `CLOUDFLARED_CONTAINER` | Cloudflared container name | `tunlify-cloudflared` |
| `DB_PATH` | SQLite database path | `/app/data/app.db` |
| `LISTEN_ADDR` | Backend listen address | `0.0.0.0:3001` |

## Architecture

```
┌─────────────────────────────────────────────┐
│  Docker Compose                             │
│                                             │
│  ┌───────────┐  ┌─────────┐  ┌──────────┐  │
│  │ frontend  │→ │ backend │→ │cloudflared│  │
│  │  (nginx)  │  │  (Go)   │  │ (tunnel) │  │
│  └───────────┘  └─────────┘  └──────────┘  │
│                      ↓                      │
│                  ┌────────┐                 │
│                  │ SQLite │                 │
│                  └────────┘                 │
└─────────────────────────────────────────────┘
```

## Project Structure

```
Tunlify/
├── docker-compose.yml
├── .env.example
├── cloudflared/         # Tunnel config & credentials
├── backend/             # Go API server
│   ├── cmd/server/
│   └── internal/
│       ├── handler/     # HTTP handlers
│       ├── service/     # Business logic
│       ├── middleware/  # Auth middleware
│       ├── db/          # SQLite
│       └── ws/          # WebSocket hub
└── frontend/            # Vue 3 SPA
    ├── src/
    │   ├── views/
    │   ├── components/
    │   ├── composables/
    │   ├── stores/
    │   └── router/
    └── public/
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/login` | Authenticate |
| POST | `/auth/logout` | End session |
| POST | `/auth/refresh` | Refresh token |
| GET | `/api/status` | Tunnel status |
| POST | `/api/control/{action}` | start / stop / restart |
| GET | `/api/logs/ws` | Realtime log stream (WebSocket) |
| GET | `/api/logs/history` | Log history |
| GET | `/api/config` | Get tunnel config |
| PUT | `/api/config` | Update tunnel config |
| POST | `/api/config/validate` | Validate YAML |
| GET | `/api/config/backups` | List config backups |
| POST | `/api/tcp-access/generate` | Generate TCP access command |
| GET | `/api/health` | Health check |
| GET | `/api/metrics` | Cloudflared metrics |
| GET | `/api/notifications` | Get notification config |
| PUT | `/api/notifications` | Update notification config |
| GET | `/api/audit` | Audit log |
| POST | `/api/settings/password` | Change password |

## Development

### Backend

```bash
cd backend
go mod tidy
JWT_SECRET=dev ADMIN_PASSWORD=dev go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Commit using [conventional commits](https://www.conventionalcommits.org/)
4. Push and open a Pull Request

See [RELEASE.md](RELEASE.md) for versioning and release details.

## License

MIT
