package repository

import "github.com/webinar/backend/internal/model"

type PostRepository interface {
	GetAll() ([]model.Post, error)
	GetByID(id int64) (*model.Post, error)
	Create(post *model.Post) error
	IncrementLikes(id int64) error
	IncrementDislikes(id int64) error
}

type CommentRepository interface {
	GetByPostID(postID int64) ([]model.Comment, error)
	GetByID(id int64) (*model.Comment, error)
	Create(comment *model.Comment) error
	IncrementLikes(id int64) error
	IncrementDislikes(id int64) error
}
