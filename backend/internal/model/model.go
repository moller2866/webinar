package model

import "time"

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"createdAt"`
	Comments  []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"createdAt"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }
