package repository

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"

	"github.com/webinar/backend/internal/model"
)

// NewSQLiteDB opens a SQLite database, configures it for concurrent use, and initializes the schema.
func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// SQLite concurrency: serialize writes, enable WAL for read concurrency, and
	// wait up to 5s for locks instead of failing immediately.
	db.SetMaxOpenConns(1)
	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, err
		}
	}

	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author TEXT NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id)
	);`

	_, err := db.Exec(schema)
	return err
}

// --- Post Repository ---

type SQLitePostRepository struct {
	db *sql.DB
}

func NewSQLitePostRepository(db *sql.DB) *SQLitePostRepository {
	return &SQLitePostRepository{db: db}
}

func (r *SQLitePostRepository) GetAll() ([]model.Post, error) {
	rows, err := r.db.Query(
		"SELECT id, title, content, author, likes, dislikes, created_at FROM posts ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var p model.Post
		var createdAt string
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Likes, &p.Dislikes, &createdAt); err != nil {
			return nil, err
		}
		if p.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *SQLitePostRepository) GetByID(id int64) (*model.Post, error) {
	row := r.db.QueryRow(
		"SELECT id, title, content, author, likes, dislikes, created_at FROM posts WHERE id = ?", id,
	)

	var p model.Post
	var createdAt string
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Likes, &p.Dislikes, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if p.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SQLitePostRepository) Create(post *model.Post) error {
	result, err := r.db.Exec(
		"INSERT INTO posts (title, content, author, likes, dislikes, created_at) VALUES (?, ?, ?, 0, 0, ?)",
		post.Title, post.Content, post.Author, post.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = id
	return nil
}

func (r *SQLitePostRepository) IncrementLikes(id int64) error {
	_, err := r.db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", id)
	return err
}

func (r *SQLitePostRepository) IncrementDislikes(id int64) error {
	_, err := r.db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", id)
	return err
}

// --- Comment Repository ---

type SQLiteCommentRepository struct {
	db *sql.DB
}

func NewSQLiteCommentRepository(db *sql.DB) *SQLiteCommentRepository {
	return &SQLiteCommentRepository{db: db}
}

func (r *SQLiteCommentRepository) GetByPostID(postID int64) ([]model.Comment, error) {
	rows, err := r.db.Query(
		"SELECT id, post_id, author, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at ASC",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		var createdAt string
		if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &createdAt); err != nil {
			return nil, err
		}
		if c.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *SQLiteCommentRepository) Create(comment *model.Comment) error {
	result, err := r.db.Exec(
		"INSERT INTO comments (post_id, author, content, created_at) VALUES (?, ?, ?, ?)",
		comment.PostID, comment.Author, comment.Content, comment.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = id
	return nil
}
