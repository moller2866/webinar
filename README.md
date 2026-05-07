# Blog Application — Three-Layer Monorepo

A blog application demonstrating clean three-layer architecture:
**Frontend** (React + TypeScript + MUI) → **Backend** (Go) → **Persistence** (PostgreSQL)

## Architecture

```
frontend/          → React 19 + Vite + MUI v6
backend/
  cmd/server/      → Entry point
  internal/
    handler/       → HTTP handlers (transport layer)
    service/       → Business logic
    repository/    → Data access (PostgreSQL, repository pattern)
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
JWT_SECRET=your-secret go run ./cmd/server
# Server starts on :8080
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
# Dev server starts on :5173
```

## Environment Variables

| Variable       | Default                                              | Where   | Purpose                                    |
|----------------|------------------------------------------------------|---------|--------------------------------------------|
| `DATABASE_URL` | `postgres://blog:blog@localhost:5432/blog?sslmode=disable` | Backend | PostgreSQL connection string          |
| `JWT_SECRET`   | *(required)*                                         | Backend | HS256 signing key for JWTs; server refuses to start if empty |
| `VITE_API_BASE`| `http://localhost:8080/api`                          | Frontend (build-time) | API base URL baked into the Vite build |

## API Endpoints

| Method | Path                        | Auth     | Description                |
|--------|-----------------------------|----------|----------------------------|
| POST   | /api/auth/register          | No       | Create account, return JWT |
| POST   | /api/auth/login             | No       | Validate credentials, return JWT |
| GET    | /api/posts                  | No       | List all posts             |
| POST   | /api/posts                  | Required | Create a post              |
| GET    | /api/posts/{id}             | No       | Get post + comments        |
| POST   | /api/posts/{id}/comments    | Required | Add comment to post        |
| POST   | /api/posts/{id}/like        | Required | Increment likes            |
| POST   | /api/posts/{id}/dislike     | Required | Increment dislikes         |
| POST   | /api/comments/{id}/like     | Required | Like a comment             |
| POST   | /api/comments/{id}/dislike  | Required | Dislike a comment          |

