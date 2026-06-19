package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"displayName"`
	PasswordHash string    `json:"-"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatarUrl"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedBy   *int64    `json:"createdBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Vote struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	TargetType string    `json:"targetType"`
	TargetID   int64     `json:"targetId"`
	VoteType   string    `json:"voteType"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Post struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	UserID     *int64    `json:"userId,omitempty"`
	CategoryID *int64    `json:"categoryId,omitempty"`
	Images     []string  `json:"images"`
	CreatedAt  time.Time `json:"createdAt"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	UserID    *int64    `json:"userId,omitempty"`
	ParentID  *int64    `json:"parentId,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
