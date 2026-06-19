package repository

import "github.com/webinar/backend/internal/model"

type PostRepository interface {
	GetAll() ([]model.Post, error)
	GetByID(id int64) (*model.Post, error)
	Create(post *model.Post) error
}

type CommentRepository interface {
	GetByPostID(postID int64) ([]model.Comment, error)
	GetByID(id int64) (*model.Comment, error)
	Create(comment *model.Comment) error
}

type UserRepository interface {
	GetByID(id int64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

type CategoryRepository interface {
	GetAll() ([]model.Category, error)
	GetByID(id int64) (*model.Category, error)
	GetBySlug(slug string) (*model.Category, error)
	Create(category *model.Category) error
}

type VoteRepository interface {
	Upsert(vote *model.Vote) error
	Delete(userID int64, targetType string, targetID int64) error
	GetByTarget(targetType string, targetID int64) ([]model.Vote, error)
	GetByUser(userID int64, targetType string, targetID int64) (*model.Vote, error)
}
