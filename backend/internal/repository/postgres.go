package repository

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/webinar/backend/internal/model"
)

// NewPostgresDB opens a PostgreSQL connection pool and initializes the schema.
func NewPostgresDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
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
		id BIGSERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author TEXT NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS comments (
		id BIGSERIAL PRIMARY KEY,
		post_id BIGINT NOT NULL REFERENCES posts(id),
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);

	ALTER TABLE comments ADD COLUMN IF NOT EXISTS likes INTEGER NOT NULL DEFAULT 0;
	ALTER TABLE comments ADD COLUMN IF NOT EXISTS dislikes INTEGER NOT NULL DEFAULT 0;`

	_, err := db.Exec(schema)
	return err
}

// --- Post Repository ---

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) *PostgresPostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) GetAll() ([]model.Post, error) {
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
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Likes, &p.Dislikes, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *PostgresPostRepository) GetByID(id int64) (*model.Post, error) {
	row := r.db.QueryRow(
		"SELECT id, title, content, author, likes, dislikes, created_at FROM posts WHERE id = $1", id,
	)

	var p model.Post
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Likes, &p.Dislikes, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresPostRepository) Create(post *model.Post) error {
	return r.db.QueryRow(
		"INSERT INTO posts (title, content, author, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		post.Title, post.Content, post.Author, post.CreatedAt,
	).Scan(&post.ID)
}

func (r *PostgresPostRepository) IncrementLikes(id int64) error {
	_, err := r.db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = $1", id)
	return err
}

func (r *PostgresPostRepository) IncrementDislikes(id int64) error {
	_, err := r.db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = $1", id)
	return err
}

// --- Comment Repository ---

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) GetByPostID(postID int64) ([]model.Comment, error) {
	rows, err := r.db.Query(
		"SELECT id, post_id, author, content, likes, dislikes, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *PostgresCommentRepository) GetByID(id int64) (*model.Comment, error) {
	row := r.db.QueryRow(
		"SELECT id, post_id, author, content, likes, dislikes, created_at FROM comments WHERE id = $1", id,
	)
	var c model.Comment
	err := row.Scan(&c.ID, &c.PostID, &c.Author, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCommentRepository) IncrementLikes(id int64) error {
	_, err := r.db.Exec("UPDATE comments SET likes = likes + 1 WHERE id = $1", id)
	return err
}

func (r *PostgresCommentRepository) IncrementDislikes(id int64) error {
	_, err := r.db.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = $1", id)
	return err
}

func (r *PostgresCommentRepository) Create(comment *model.Comment) error {
	return r.db.QueryRow(
		"INSERT INTO comments (post_id, author, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		comment.PostID, comment.Author, comment.Content, comment.CreatedAt,
	).Scan(&comment.ID)
}
