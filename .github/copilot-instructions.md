# Webinar Blog — Project Guidelines

## Architecture

This is a monorepo with two top-level workspaces:

- `backend/` — Go API server, port 8080
- `frontend/` — React + TypeScript SPA (Vite), port 5173

**Go module path**: `github.com/webinar/backend`

**Layered architecture** (backend):
```
handler → service → repository → model
```
Each layer only imports the one directly below it. Never skip layers.

## API Endpoints

All endpoints are defined in `backend/internal/handler/handler.go` and wired in `backend/cmd/server/main.go`. Update the table below and the handler code when adding new endpoints.


| Method | Path                        | Description           |
|--------|-----------------------------|-----------------------|
| GET    | /api/posts                  | List all posts        |
| POST   | /api/posts                  | Create a post         |
| GET    | /api/posts/{id}             | Get post + comments   |
| POST   | /api/posts/{id}/comments    | Add comment to post   |
| POST   | /api/posts/{id}/like        | Increment likes       |
| POST   | /api/posts/{id}/dislike     | Increment dislikes    |
| POST   | /api/comments/{id}/like     | Like a comment        |
| POST   | /api/comments/{id}/dislike  | Dislike a comment     |

## Key Design Decisions

- **No auth** — `author` is a plain text field. This is intentional; auth would obscure the architectural patterns.
- **No external router** — Go 1.22+ method+pattern routing is used (`"GET /api/posts/{id}"`). Do not introduce chi, gin, or gorilla/mux.
- **PostgreSQL** — The database runs as a separate container (`postgres:17-alpine`). The Go backend connects via `jackc/pgx/v5` through the standard `database/sql` interface. There is no embedded database.
- **No frontend state library** — `useState`/`useEffect` is sufficient. Do not introduce Redux, Zustand, Jotai, etc.
- **No axios** — use the built-in `fetch` API via `frontend/src/api.ts`.

## Wiring

All dependency injection happens in `backend/cmd/server/main.go`:
```
NewPostgresDB → NewPostgresPostRepository + NewPostgresCommentRepository → NewPostService → NewHandler → mux
```
Add new wiring there only.

## Build & Run

```bash
docker compose up --build
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### Local Development (without Docker)

Requires a running PostgreSQL instance on `localhost:5432` with database `blog`.

```bash
# Backend
cd backend && go run ./cmd/server
# → http://localhost:8080

# Frontend
cd frontend && npm run dev
# → http://localhost:5173
```

## Environment Variables

| Variable | Default | Where | Purpose |
|----------|---------|-------|---------|
| `DATABASE_URL` | `postgres://blog:blog@localhost:5432/blog?sslmode=disable` | Backend | PostgreSQL connection string |
| `VITE_API_BASE` | `http://localhost:8080/api` | Frontend (build-time) | API base URL baked into the Vite build |

**⚠ Blank-page pitfall:** In Docker, the nginx config proxies `/api` to the backend, so `VITE_API_BASE` should be `/api` (relative). If it's set to `http://localhost:8080/api` inside the container, the browser can't reach it and the page renders blank with no visible error.

# Instructions

whenever you use a custom instruction in the `.github/instructions` directory, make sure to state which instruction(s) you are using and why. For example:

```md
I am working on the repository layer. I will use the instruction in `.github/instructions/repository.md` to guide my implementation because it provides the necessary steps and best practices for updating the repository layer.
```
