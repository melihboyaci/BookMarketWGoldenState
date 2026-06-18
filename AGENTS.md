# AGENTS.md

## Cursor Cloud specific instructions

Golden State is a single product made of three services that must all run for end-to-end UI testing:

| Service | Dir / Cmd | Port | Notes |
| --- | --- | --- | --- |
| PostgreSQL 16 | `docker compose up -d` (repo root) | 5432 | Required; backend fails on startup if it can't ping the DB. |
| Go API (Gin) | `go run ./cmd/api/main.go` (repo root) | 8080 | Run from repo root: `Seed()` reads the relative path `internal/db/seed.sql`. |
| React/Vite frontend | `npm run dev` (in `frontend/`) | 5173 | Vite proxies `/api` → `http://localhost:8080`. |

Standard setup/run/test commands live in `README.md`. The notes below are the non-obvious gotchas.

### Docker is required and the daemon is not auto-started
Postgres and the integration tests (Testcontainers) both need Docker. Start the daemon once per session, e.g. `sudo dockerd &` (or in a tmux session) and use `sudo docker ...` / `sudo -E docker compose ...`. This box runs Docker 29 with `/etc/docker/daemon.json` set to the `fuse-overlayfs` storage driver and `containerd-snapshotter` disabled — needed for Docker-in-VM to work.

### Environment variables (`.env` is gitignored, no example committed)
The backend reads `POSTGRES_HOST/PORT/USER/PASSWORD/DB` directly via `os.Getenv` and `docker-compose.yml` needs `POSTGRES_USER/PASSWORD/DB`. Create a repo-root `.env` (also loaded by `godotenv`). Critically, `POSTGRES_HOST` is NOT in `docker-compose.yml`, so set `POSTGRES_HOST=localhost` for the host-run backend. The compose `POSTGRES_*` values must match the `.env` the backend uses.

### Migrations are NOT applied automatically
`main.go` only calls `db.Connect()` then `db.Seed()`; `seed.sql` assumes the `books`/`users`/`orders` tables already exist. Before the first backend run against a fresh DB, apply `internal/db/migrations/001..005_*.sql` in numeric order, then the seeder runs on startup. Example: `for f in internal/db/migrations/0*.sql; do sudo docker exec -i bookmarket_postgres psql -U <user> -d <db> < "$f"; done`. The Go tests apply these migrations themselves against a throwaway Testcontainers Postgres, so tests need no manual DB setup.

### Tests
`go test ./...` needs Docker (Testcontainers spins a real Postgres); run it with `sudo -E env "PATH=$PATH" go test ./...` if the Docker socket needs root. `go test -short ./...` skips the integration test and needs no Docker.

### Lint
`npm run lint` (in `frontend/`) currently reports pre-existing ESLint errors that are unrelated to environment setup; the tooling itself works. `go vet ./...` passes.

### Demo accounts (seeded)
`admin@demo.com / admin123`, `seller@demo.com / seller123`, `buyer@demo.com / buyer123`. Only ADMIN sees the "Sistemi Sıfırla" Golden State reset button.
