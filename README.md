# Blog Application — Three-Layer Monorepo

A blog application demonstrating clean three-layer architecture:
**Frontend** (React + TypeScript + MUI) → **Backend** (Go) → **Persistence** (SQLite)

## Architecture

```
frontend/          → React 19 + Vite + MUI v6
backend/
  cmd/server/      → Entry point
  internal/
    handler/       → HTTP handlers (transport layer)
    service/       → Business logic
    repository/    → Data access (SQLite, repository pattern)
    model/         → Domain types
```

## Quick Start

### Docker (recommended)
```bash
docker compose up --build
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### Local Development

#### Backend
```bash
cd backend
go run ./cmd/server
# Server starts on :8080
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
# Dev server starts on :5173
```

## API Endpoints

| Method | Path                        | Description           |
|--------|-----------------------------|-----------------------|
| GET    | /api/posts                  | List all posts        |
| POST   | /api/posts                  | Create a post         |
| GET    | /api/posts/{id}             | Get post + comments   |
| POST   | /api/posts/{id}/comments    | Add comment to post   |
| POST   | /api/posts/{id}/like        | Increment likes       |
| POST   | /api/posts/{id}/dislike     | Increment dislikes    |
