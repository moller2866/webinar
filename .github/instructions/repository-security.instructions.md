---
description: "Use when writing or modifying repository package code — postgres.go, repository.go, or any database query. Covers SQL injection prevention, query safety, schema conventions, and PostgreSQL-specific best practices."
applyTo: "backend/internal/repository/**/*.go"
---

# Repository & Database Security Conventions

## SQL Injection — Always Use Parameterised Queries

Never interpolate values into SQL strings. Always use numbered `$1`, `$2`, ... placeholders and pass values as arguments.

```go
// WRONG — SQL injection risk
db.Query("SELECT * FROM posts WHERE author = '" + author + "'")

// CORRECT
db.Query("SELECT id, title FROM posts WHERE author = $1", author)
```

This applies to every call: `db.Query`, `db.QueryRow`, `db.Exec`.

## Column Allowlists for Dynamic Queries

If you ever need a dynamic `ORDER BY` or `WHERE` column, validate against an explicit allowlist — never interpolate user input:

```go
allowed := map[string]bool{"title": true, "created_at": true, "likes": true}
if !allowed[col] {
    return nil, errors.New("invalid sort column")
}
query := "SELECT ... FROM posts ORDER BY " + col // safe — col is validated
```

## Select Only What You Need

Never use `SELECT *`. List columns explicitly. This prevents accidentally exposing new columns added to the schema and makes scan order deterministic.

```go
// WRONG
db.Query("SELECT * FROM posts WHERE id = $1", id)

// CORRECT
db.QueryRow("SELECT id, title, content, author, likes, dislikes, created_at FROM posts WHERE id = $1", id)
```

## Close Rows Immediately with defer

Always `defer rows.Close()` immediately after a successful `db.Query` call. Check `rows.Err()` after the iteration loop.

```go
rows, err := r.db.Query("SELECT ...", args)
if err != nil {
    return nil, err
}
defer rows.Close()

for rows.Next() { ... }
return results, rows.Err()
```

## Use Foreign Keys

PostgreSQL enforces `FOREIGN KEY` constraints natively. All new tables with FK relationships must declare them in the schema — do not rely on application-level enforcement alone.

## Schema Changes

- Always use `CREATE TABLE IF NOT EXISTS` — schema init runs on every startup.
- Never `DROP` or `ALTER` columns in `initSchema` without a migration strategy. For this project, add new columns with `DEFAULT` values or create new tables.
- Add indexes for columns used in `WHERE` or `ORDER BY` clauses in frequent queries:

```sql
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
```

## Atomic Counter Updates

Use SQL arithmetic for counters — never read-then-write:

```go
// WRONG — race condition
post.Likes++
db.Exec("UPDATE posts SET likes = $1 WHERE id = $2", post.Likes, id)

// CORRECT — atomic
db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = $1", id)
```

## Limit Result Sets

For list queries, add a reasonable `LIMIT` if the table can grow unbounded. This prevents a single query from returning the entire table to the caller:

```go
db.Query("SELECT ... FROM posts ORDER BY created_at DESC LIMIT 100")
```

## No SQL Outside This Package

SQL strings must only appear in `repository/postgres.go`. No other package may import `database/sql` or construct queries. The interfaces in `repository.go` are the only contract exposed to the rest of the application.
