package repository

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"

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
	CREATE TABLE IF NOT EXISTS users (
		id           BIGSERIAL PRIMARY KEY,
		email        TEXT UNIQUE NOT NULL,
		display_name TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		bio          TEXT NOT NULL DEFAULT '',
		avatar_url   TEXT NOT NULL DEFAULT '',
		created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS categories (
		id          BIGSERIAL PRIMARY KEY,
		name        TEXT UNIQUE NOT NULL,
		slug        TEXT UNIQUE NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		created_by  BIGINT REFERENCES users(id),
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS posts (
		id          BIGSERIAL PRIMARY KEY,
		title       TEXT NOT NULL,
		content     TEXT NOT NULL,
		user_id     BIGINT REFERENCES users(id),
		category_id BIGINT REFERENCES categories(id),
		images      TEXT[] NOT NULL DEFAULT '{}',
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS comments (
		id         BIGSERIAL PRIMARY KEY,
		post_id    BIGINT NOT NULL REFERENCES posts(id),
		user_id    BIGINT REFERENCES users(id),
		parent_id  BIGINT REFERENCES comments(id),
		content    TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS votes (
		id          BIGSERIAL PRIMARY KEY,
		user_id     BIGINT NOT NULL REFERENCES users(id),
		target_type TEXT NOT NULL CHECK (target_type IN ('post', 'comment')),
		target_id   BIGINT NOT NULL,
		vote_type   TEXT NOT NULL CHECK (vote_type IN ('highfive', 'meh')),
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE (user_id, target_type, target_id)
	);

	CREATE INDEX IF NOT EXISTS idx_comments_post_id   ON comments(post_id);
	CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
	CREATE INDEX IF NOT EXISTS idx_posts_user_id      ON posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_posts_category_id  ON posts(category_id);
	CREATE INDEX IF NOT EXISTS idx_votes_target       ON votes(target_type, target_id);
	CREATE INDEX IF NOT EXISTS idx_categories_slug    ON categories(slug);

	-- Migration: add new columns to existing posts table (safe re-runs)
	ALTER TABLE posts ADD COLUMN IF NOT EXISTS user_id     BIGINT REFERENCES users(id);
	ALTER TABLE posts ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES categories(id);
	ALTER TABLE posts ADD COLUMN IF NOT EXISTS images      TEXT[] NOT NULL DEFAULT '{}';

	-- Migration: drop deprecated posts columns
	ALTER TABLE posts DROP COLUMN IF EXISTS author;
	ALTER TABLE posts DROP COLUMN IF EXISTS tags;
	ALTER TABLE posts DROP COLUMN IF EXISTS likes;
	ALTER TABLE posts DROP COLUMN IF EXISTS dislikes;

	-- Migration: add new columns to existing comments table (safe re-runs)
	ALTER TABLE comments ADD COLUMN IF NOT EXISTS user_id   BIGINT REFERENCES users(id);
	ALTER TABLE comments ADD COLUMN IF NOT EXISTS parent_id BIGINT REFERENCES comments(id);

	-- Migration: drop deprecated comments columns
	ALTER TABLE comments DROP COLUMN IF EXISTS author;
	ALTER TABLE comments DROP COLUMN IF EXISTS likes;
	ALTER TABLE comments DROP COLUMN IF EXISTS dislikes;`

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
		"SELECT id, title, content, user_id, category_id, images, created_at FROM posts ORDER BY created_at DESC LIMIT 100",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var p model.Post
		var userID, categoryID sql.NullInt64
		var images pq.StringArray
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &userID, &categoryID, &images, &p.CreatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			p.UserID = &userID.Int64
		}
		if categoryID.Valid {
			p.CategoryID = &categoryID.Int64
		}
		p.Images = []string(images)
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *PostgresPostRepository) GetByID(id int64) (*model.Post, error) {
	row := r.db.QueryRow(
		"SELECT id, title, content, user_id, category_id, images, created_at FROM posts WHERE id = $1", id,
	)

	var p model.Post
	var userID, categoryID sql.NullInt64
	var images pq.StringArray
	err := row.Scan(&p.ID, &p.Title, &p.Content, &userID, &categoryID, &images, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		p.UserID = &userID.Int64
	}
	if categoryID.Valid {
		p.CategoryID = &categoryID.Int64
	}
	p.Images = []string(images)
	return &p, nil
}

func (r *PostgresPostRepository) Create(post *model.Post) error {
	return r.db.QueryRow(
		"INSERT INTO posts (title, content, user_id, category_id, images, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		post.Title, post.Content, post.UserID, post.CategoryID, pq.StringArray(post.Images), post.CreatedAt,
	).Scan(&post.ID)
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
		"SELECT id, post_id, user_id, parent_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		var userID, parentID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.PostID, &userID, &parentID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			c.UserID = &userID.Int64
		}
		if parentID.Valid {
			c.ParentID = &parentID.Int64
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *PostgresCommentRepository) GetByID(id int64) (*model.Comment, error) {
	row := r.db.QueryRow(
		"SELECT id, post_id, user_id, parent_id, content, created_at FROM comments WHERE id = $1", id,
	)
	var c model.Comment
	var userID, parentID sql.NullInt64
	err := row.Scan(&c.ID, &c.PostID, &userID, &parentID, &c.Content, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		c.UserID = &userID.Int64
	}
	if parentID.Valid {
		c.ParentID = &parentID.Int64
	}
	return &c, nil
}

func (r *PostgresCommentRepository) Create(comment *model.Comment) error {
	return r.db.QueryRow(
		"INSERT INTO comments (post_id, user_id, parent_id, content, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		comment.PostID, comment.UserID, comment.ParentID, comment.Content, comment.CreatedAt,
	).Scan(&comment.ID)
}

// --- User Repository ---

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(id int64) (*model.User, error) {
	row := r.db.QueryRow(
		"SELECT id, email, display_name, password_hash, bio, avatar_url, created_at FROM users WHERE id = $1", id,
	)
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*model.User, error) {
	row := r.db.QueryRow(
		"SELECT id, email, display_name, password_hash, bio, avatar_url, created_at FROM users WHERE email = $1", email,
	)
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.DisplayName, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Create(user *model.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, display_name, password_hash, bio, avatar_url, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Email, user.DisplayName, user.PasswordHash, user.Bio, user.AvatarURL, user.CreatedAt,
	).Scan(&user.ID)
}

func (r *PostgresUserRepository) Update(user *model.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET display_name = $1, bio = $2, avatar_url = $3 WHERE id = $4",
		user.DisplayName, user.Bio, user.AvatarURL, user.ID,
	)
	return err
}

// --- Category Repository ---

type PostgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) GetAll() ([]model.Category, error) {
	rows, err := r.db.Query(
		"SELECT id, name, slug, description, created_by, created_at FROM categories ORDER BY name ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []model.Category{}
	for rows.Next() {
		var c model.Category
		var createdBy sql.NullInt64
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &createdBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		if createdBy.Valid {
			c.CreatedBy = &createdBy.Int64
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *PostgresCategoryRepository) GetByID(id int64) (*model.Category, error) {
	row := r.db.QueryRow(
		"SELECT id, name, slug, description, created_by, created_at FROM categories WHERE id = $1", id,
	)
	var c model.Category
	var createdBy sql.NullInt64
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &createdBy, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if createdBy.Valid {
		c.CreatedBy = &createdBy.Int64
	}
	return &c, nil
}

func (r *PostgresCategoryRepository) GetBySlug(slug string) (*model.Category, error) {
	row := r.db.QueryRow(
		"SELECT id, name, slug, description, created_by, created_at FROM categories WHERE slug = $1", slug,
	)
	var c model.Category
	var createdBy sql.NullInt64
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &createdBy, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if createdBy.Valid {
		c.CreatedBy = &createdBy.Int64
	}
	return &c, nil
}

func (r *PostgresCategoryRepository) Create(category *model.Category) error {
	return r.db.QueryRow(
		"INSERT INTO categories (name, slug, description, created_by, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		category.Name, category.Slug, category.Description, category.CreatedBy, category.CreatedAt,
	).Scan(&category.ID)
}

// --- Vote Repository ---

type PostgresVoteRepository struct {
	db *sql.DB
}

func NewPostgresVoteRepository(db *sql.DB) *PostgresVoteRepository {
	return &PostgresVoteRepository{db: db}
}

func (r *PostgresVoteRepository) Upsert(vote *model.Vote) error {
	return r.db.QueryRow(
		`INSERT INTO votes (user_id, target_type, target_id, vote_type, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, target_type, target_id)
		 DO UPDATE SET vote_type = EXCLUDED.vote_type, created_at = EXCLUDED.created_at
		 RETURNING id`,
		vote.UserID, vote.TargetType, vote.TargetID, vote.VoteType, vote.CreatedAt,
	).Scan(&vote.ID)
}

func (r *PostgresVoteRepository) Delete(userID int64, targetType string, targetID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM votes WHERE user_id = $1 AND target_type = $2 AND target_id = $3",
		userID, targetType, targetID,
	)
	return err
}

func (r *PostgresVoteRepository) GetByTarget(targetType string, targetID int64) ([]model.Vote, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, target_type, target_id, vote_type, created_at FROM votes WHERE target_type = $1 AND target_id = $2",
		targetType, targetID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := []model.Vote{}
	for rows.Next() {
		var v model.Vote
		if err := rows.Scan(&v.ID, &v.UserID, &v.TargetType, &v.TargetID, &v.VoteType, &v.CreatedAt); err != nil {
			return nil, err
		}
		votes = append(votes, v)
	}
	return votes, rows.Err()
}

func (r *PostgresVoteRepository) GetByUser(userID int64, targetType string, targetID int64) (*model.Vote, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, target_type, target_id, vote_type, created_at FROM votes WHERE user_id = $1 AND target_type = $2 AND target_id = $3",
		userID, targetType, targetID,
	)
	var v model.Vote
	err := row.Scan(&v.ID, &v.UserID, &v.TargetType, &v.TargetID, &v.VoteType, &v.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}
