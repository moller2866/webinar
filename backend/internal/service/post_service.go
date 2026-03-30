package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/webinar/backend/internal/model"
	"github.com/webinar/backend/internal/repository"
)

type PostService struct {
	posts    repository.PostRepository
	comments repository.CommentRepository
}

func NewPostService(posts repository.PostRepository, comments repository.CommentRepository) *PostService {
	return &PostService{posts: posts, comments: comments}
}

func (s *PostService) ListPosts() ([]model.Post, error) {
	return s.posts.GetAll()
}

func (s *PostService) GetPost(id int64) (*model.Post, error) {
	post, err := s.posts.GetByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, &model.ValidationError{Message: fmt.Sprintf("post %d not found", id)}
	}

	comments, err := s.comments.GetByPostID(id)
	if err != nil {
		return nil, err
	}
	post.Comments = comments
	return post, nil
}

func (s *PostService) CreatePost(post *model.Post) error {
	if post.Title == "" {
		return &model.ValidationError{Message: "title is required"}
	}
	if post.Content == "" {
		return &model.ValidationError{Message: "content is required"}
	}
	if post.Author == "" {
		return &model.ValidationError{Message: "author is required"}
	}
	post.Tags = sanitizeTags(post.Tags)
	post.CreatedAt = time.Now()
	return s.posts.Create(post)
}

func (s *PostService) LikePost(id int64) error {
	post, err := s.posts.GetByID(id)
	if err != nil {
		return err
	}
	if post == nil {
		return &model.ValidationError{Message: fmt.Sprintf("post %d not found", id)}
	}
	return s.posts.IncrementLikes(id)
}

func (s *PostService) DislikePost(id int64) error {
	post, err := s.posts.GetByID(id)
	if err != nil {
		return err
	}
	if post == nil {
		return &model.ValidationError{Message: fmt.Sprintf("post %d not found", id)}
	}
	return s.posts.IncrementDislikes(id)
}

func (s *PostService) AddComment(comment *model.Comment) error {
	if comment.Author == "" {
		return &model.ValidationError{Message: "author is required"}
	}
	if comment.Content == "" {
		return &model.ValidationError{Message: "content is required"}
	}

	post, err := s.posts.GetByID(comment.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		return &model.ValidationError{Message: fmt.Sprintf("post %d not found", comment.PostID)}
	}

	comment.CreatedAt = time.Now()
	return s.comments.Create(comment)
}

func (s *PostService) LikeComment(id int64) error {
	comment, err := s.comments.GetByID(id)
	if err != nil {
		return err
	}
	if comment == nil {
		return &model.ValidationError{Message: fmt.Sprintf("comment %d not found", id)}
	}
	return s.comments.IncrementLikes(id)
}

func (s *PostService) DislikeComment(id int64) error {
	comment, err := s.comments.GetByID(id)
	if err != nil {
		return err
	}
	if comment == nil {
		return &model.ValidationError{Message: fmt.Sprintf("comment %d not found", id)}
	}
	return s.comments.IncrementDislikes(id)
}

// sanitizeTags trims whitespace from each tag, removes empty strings, and deduplicates.
func sanitizeTags(in []string) []string {
	seen := make(map[string]bool)
	out := []string{}
	for _, t := range in {
		t = strings.TrimSpace(t)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}
