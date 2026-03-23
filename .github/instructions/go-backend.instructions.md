---
description: "Use when writing Go backend code — handlers, services, repositories, models, or main.go wiring. Covers layer responsibilities, error handling, routing, and repository pattern."
applyTo: "backend/**/*.go"
---

# Go Backend Conventions

## Layer Responsibilities

| Package | Owns | Must NOT |
|---------|------|----------|
| `model/` | Domain structs (`Post`, `Comment`, `ValidationError`) | Touch SQL or HTTP |
| `repository/` | Interfaces in `repository.go` + SQLite implementations in `sqlite.go` | Contain business logic |
| `service/` | Business logic, validation, orchestration | Touch `net/http` or SQL directly |
| `handler/` | HTTP request/response DTOs, parsing, routing | Contain business logic |

## Dependency Injection

All structs are constructed via `New*` functions receiving their dependencies as interfaces or concrete types:

```go
func NewPostService(posts repository.PostRepository, comments repository.CommentRepository) *PostService
func NewHandler(postService *service.PostService) *Handler
```

Never instantiate concrete repository types (e.g. `SQLitePostRepository`) outside of `main.go`.

## Error Handling

- User-facing validation errors: return `*model.ValidationError{Message: "..."}` from the service layer
- Handlers distinguish error types with `errors.As`:
  - `*model.ValidationError` → 400 Bad Request
  - Any other error → 500 Internal Server Error

```go
var ve *model.ValidationError
if errors.As(err, &ve) {
    writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: ve.Message})
    return
}
writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
```

## Routing

Use Go 1.22+ method+pattern routing — no external router:

```go
mux.HandleFunc("GET /api/posts/{id}", h.getPost)
```

Extract path values with `r.PathValue("id")`. Do not introduce chi, gin, or gorilla/mux.

## Repository Pattern

When adding a new data operation:
1. Add the method to the interface in `repository/repository.go` **first**
2. Implement it in `repository/sqlite.go`
3. Consume it in the service layer — never call SQL from outside `repository/`

## SQLite

Driver: `modernc.org/sqlite` — pure Go, no CGO required.
Use the `database/sql` standard interface. Never import the driver directly outside `repository/sqlite.go`.
