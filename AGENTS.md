# AGENTS.md

## Cursor Cloud specific instructions

This is a Go (Gin) + React/Vite "Golden State" B2B book-marketplace demo backed by PostgreSQL.
Standard run/test/lint commands live in `README.md`; only the non-obvious caveats are captured here.

### Services
- **PostgreSQL**: runs via `docker compose up -d` (image `postgres:16-alpine`, exposed on `:5432`).
- **Backend (Go/Gin)**: `go run ./cmd/api/main.go`, listens on `:8080`.
- **Frontend (React/Vite)**: `npm run dev` inside `frontend/`, serves on `:5173` and proxies `/api` → `http://localhost:8080`.

### Startup caveats (read before running anything)
- **Docker is required** for both PostgreSQL and the Go test suite (the integration test uses Testcontainers). The daemon is installed in the VM snapshot but is not auto-started: start it once per session with `sudo dockerd` (e.g. in a background tmux session) and ensure the socket is usable without sudo via `sudo chmod 666 /var/run/docker.sock`. Docker is configured with the `fuse-overlayfs` storage driver and `containerd-snapshotter` disabled in `/etc/docker/daemon.json` (required for this nested-VM environment).
- **`.env` is required and git-ignored.** Both `docker compose` and the backend read it. It must define `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_HOST` (`localhost`), `POSTGRES_PORT` (`5432`), plus `APP_PORT`, `APP_ENV`, `JWT_SECRET`. The backend connects with these `POSTGRES_*` vars (not a single `DATABASE_URL`).
- **The app does NOT run DB migrations.** `main.go` only runs `internal/db/seed.sql` (which is INSERT-only). On a fresh database you must first apply the schema by running the files in `internal/db/migrations/` in numeric order (e.g. `for f in internal/db/migrations/0*.sql; do docker exec -i bookmarket_postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 < "$f"; done`). If tables are missing, seeding fails silently (logs a warning) and all API calls error out. Migrations and seed are idempotent.
- **Run the backend from the repo root.** `db.Seed` reads `internal/db/seed.sql` via a relative path, so the working directory must be the repository root.

### Testing / lint
- `go test ./...` runs unit + Testcontainers integration tests (needs Docker running). `go test -short ./...` skips the integration test.
- `npm run lint` (in `frontend/`) currently reports ~22 pre-existing lint errors in the app source; this is the repo's baseline, not an environment problem.

### Demo accounts
`admin@demo.com / admin123` (ADMIN), `seller@demo.com / seller123` (SELLER), `buyer@demo.com / buyer123` (BUYER). The signature feature is the admin "Sistemi Sıfırla" (Golden State) reset.
