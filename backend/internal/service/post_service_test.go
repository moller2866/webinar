package service

import (
	"errors"
	"testing"
	"time"

	"github.com/webinar/backend/internal/model"
)

// --- Mock repositories ---

type mockPostRepository struct {
	posts   []model.Post
	nextID  int64
	forceErr error
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{nextID: 1}
}

func (m *mockPostRepository) GetAll() ([]model.Post, error) {
	if m.forceErr != nil {
		return nil, m.forceErr
	}
	return m.posts, nil
}

func (m *mockPostRepository) GetByID(id int64) (*model.Post, error) {
	if m.forceErr != nil {
		return nil, m.forceErr
	}
	for i, p := range m.posts {
		if p.ID == id {
			return &m.posts[i], nil
		}
	}
	return nil, nil
}

func (m *mockPostRepository) Create(post *model.Post) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	post.ID = m.nextID
	m.nextID++
	m.posts = append(m.posts, *post)
	return nil
}

type mockCommentRepository struct {
	comments []model.Comment
	nextID   int64
	forceErr  error
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{nextID: 1}
}

func (m *mockCommentRepository) GetByPostID(postID int64) ([]model.Comment, error) {
	if m.forceErr != nil {
		return nil, m.forceErr
	}
	result := []model.Comment{}
	for _, c := range m.comments {
		if c.PostID == postID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCommentRepository) GetByID(id int64) (*model.Comment, error) {
	if m.forceErr != nil {
		return nil, m.forceErr
	}
	for i, c := range m.comments {
		if c.ID == id {
			return &m.comments[i], nil
		}
	}
	return nil, nil
}

func (m *mockCommentRepository) Create(comment *model.Comment) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	comment.ID = m.nextID
	m.nextID++
	m.comments = append(m.comments, *comment)
	return nil
}

// --- Helpers ---

func newTestService() (*PostService, *mockPostRepository, *mockCommentRepository) {
	posts := newMockPostRepository()
	comments := newMockCommentRepository()
	svc := NewPostService(posts, comments)
	return svc, posts, comments
}

// --- CreatePost tests ---

func TestCreatePost_Valid(t *testing.T) {
	svc, postRepo, _ := newTestService()

	post := &model.Post{Title: "Hello", Content: "World"}
	if err := svc.CreatePost(post); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != 1 {
		t.Errorf("expected ID=1, got %d", post.ID)
	}
	if post.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if len(postRepo.posts) != 1 {
		t.Errorf("expected 1 post stored, got %d", len(postRepo.posts))
	}
}

func TestCreatePost_MissingTitle(t *testing.T) {
	svc, _, _ := newTestService()

	err := svc.CreatePost(&model.Post{Content: "no title"})
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if ve.Message != "title is required" {
		t.Errorf("unexpected message: %s", ve.Message)
	}
}

func TestCreatePost_MissingContent(t *testing.T) {
	svc, _, _ := newTestService()

	err := svc.CreatePost(&model.Post{Title: "no content"})
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if ve.Message != "content is required" {
		t.Errorf("unexpected message: %s", ve.Message)
	}
}

func TestCreatePost_NilImagesDefaultsToEmpty(t *testing.T) {
	svc, postRepo, _ := newTestService()

	post := &model.Post{Title: "T", Content: "C", Images: nil}
	if err := svc.CreatePost(post); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stored := postRepo.posts[0]
	if stored.Images == nil {
		t.Error("expected Images to be non-nil empty slice")
	}
	if len(stored.Images) != 0 {
		t.Errorf("expected 0 images, got %d", len(stored.Images))
	}
}

func TestCreatePost_RepositoryError(t *testing.T) {
	svc, postRepo, _ := newTestService()
	postRepo.forceErr = errors.New("db error")

	err := svc.CreatePost(&model.Post{Title: "T", Content: "C"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- GetPost tests ---

func TestGetPost_Found(t *testing.T) {
	svc, postRepo, commentRepo := newTestService()

	postRepo.posts = []model.Post{{ID: 1, Title: "T", Content: "C", CreatedAt: time.Now(), Images: []string{}}}
	commentRepo.comments = []model.Comment{
		{ID: 1, PostID: 1, Content: "nice", CreatedAt: time.Now()},
	}

	post, err := svc.GetPost(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != 1 {
		t.Errorf("expected ID=1, got %d", post.ID)
	}
	if len(post.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(post.Comments))
	}
}

func TestGetPost_NotFound(t *testing.T) {
	svc, _, _ := newTestService()

	_, err := svc.GetPost(99)
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

// --- AddComment tests ---

func TestAddComment_Valid(t *testing.T) {
	svc, postRepo, commentRepo := newTestService()

	postRepo.posts = []model.Post{{ID: 1, Title: "T", Content: "C", CreatedAt: time.Now(), Images: []string{}}}

	comment := &model.Comment{PostID: 1, Content: "great post"}
	if err := svc.AddComment(comment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 1 {
		t.Errorf("expected ID=1, got %d", comment.ID)
	}
	if comment.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if len(commentRepo.comments) != 1 {
		t.Errorf("expected 1 comment stored, got %d", len(commentRepo.comments))
	}
}

func TestAddComment_MissingContent(t *testing.T) {
	svc, postRepo, _ := newTestService()
	postRepo.posts = []model.Post{{ID: 1, Title: "T", Content: "C", Images: []string{}}}

	err := svc.AddComment(&model.Comment{PostID: 1, Content: ""})
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if ve.Message != "content is required" {
		t.Errorf("unexpected message: %s", ve.Message)
	}
}

func TestAddComment_PostNotFound(t *testing.T) {
	svc, _, _ := newTestService()

	err := svc.AddComment(&model.Comment{PostID: 99, Content: "hello"})
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestAddComment_WithParentID(t *testing.T) {
	svc, postRepo, commentRepo := newTestService()

	postRepo.posts = []model.Post{{ID: 1, Title: "T", Content: "C", CreatedAt: time.Now(), Images: []string{}}}

	parentID := int64(42)
	comment := &model.Comment{PostID: 1, Content: "reply", ParentID: &parentID}
	if err := svc.AddComment(comment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stored := commentRepo.comments[0]
	if stored.ParentID == nil || *stored.ParentID != 42 {
		t.Errorf("expected ParentID=42, got %v", stored.ParentID)
	}
}

// --- ListPosts tests ---

func TestListPosts(t *testing.T) {
	svc, postRepo, _ := newTestService()

	postRepo.posts = []model.Post{
		{ID: 1, Title: "A", Content: "a", Images: []string{}},
		{ID: 2, Title: "B", Content: "b", Images: []string{}},
	}

	posts, err := svc.ListPosts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestListPosts_RepositoryError(t *testing.T) {
	svc, postRepo, _ := newTestService()
	postRepo.forceErr = errors.New("db error")

	_, err := svc.ListPosts()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
