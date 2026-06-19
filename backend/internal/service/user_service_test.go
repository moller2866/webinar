package service_test

import (
	"testing"

	"github.com/webinar/backend/internal/model"
	"github.com/webinar/backend/internal/service"
)

// --- stub repository ---

type stubUserRepo struct {
	users map[string]*model.User
	nextID int64
}

func newStubUserRepo() *stubUserRepo {
	return &stubUserRepo{users: make(map[string]*model.User)}
}

func (r *stubUserRepo) Create(user *model.User) error {
	r.nextID++
	user.ID = r.nextID
	r.users[user.Email] = user
	return nil
}

func (r *stubUserRepo) GetByEmail(email string) (*model.User, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *stubUserRepo) GetByID(id int64) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

// --- tests ---

func TestRegister_Success(t *testing.T) {
	svc := service.NewUserService(newStubUserRepo())
	user, err := svc.Register("Alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", user.Email)
	}
	if user.DisplayName != "Alice" {
		t.Errorf("expected display name Alice, got %s", user.DisplayName)
	}
	if user.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}
	if user.PasswordHash == "password123" {
		t.Error("password must be hashed, not stored in plain text")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newStubUserRepo()
	svc := service.NewUserService(repo)

	_, err := svc.Register("Alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	_, err = svc.Register("Alice2", "alice@example.com", "password456")
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	var ce *model.ConflictError
	if !isConflictError(err, &ce) {
		t.Errorf("expected *model.ConflictError, got %T: %v", err, err)
	}
}

func TestRegister_ValidationErrors(t *testing.T) {
	svc := service.NewUserService(newStubUserRepo())

	cases := []struct {
		name        string
		displayName string
		email       string
		password    string
	}{
		{"empty display name", "", "a@b.com", "password123"},
		{"short display name", "A", "a@b.com", "password123"},
		{"empty email", "Alice", "", "password123"},
		{"short password", "Alice", "a@b.com", "short"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Register(tc.displayName, tc.email, tc.password)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			var ve *model.ValidationError
			if !isValidationError(err, &ve) {
				t.Errorf("expected *model.ValidationError, got %T: %v", err, err)
			}
		})
	}
}

func TestLogin_Success(t *testing.T) {
	svc := service.NewUserService(newStubUserRepo())

	_, err := svc.Register("Alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, err := svc.Login("alice@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", user.Email)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := service.NewUserService(newStubUserRepo())

	_, err := svc.Register("Alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, err := svc.Login("alice@example.com", "wrongpassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Error("expected nil user for wrong password, got user")
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	svc := service.NewUserService(newStubUserRepo())

	user, err := svc.Login("nobody@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Error("expected nil user for unknown email, got user")
	}
}

func TestLogin_EmailNormalized(t *testing.T) {
	svc := service.NewUserService(newStubUserRepo())

	_, err := svc.Register("Alice", "ALICE@EXAMPLE.COM", "password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, err := svc.Login("alice@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user after email normalization, got nil")
	}
}

// helpers to avoid importing errors package in tests

func isValidationError(err error, target **model.ValidationError) bool {
	if ve, ok := err.(*model.ValidationError); ok {
		*target = ve
		return true
	}
	return false
}

func isConflictError(err error, target **model.ConflictError) bool {
	if ce, ok := err.(*model.ConflictError); ok {
		*target = ce
		return true
	}
	return false
}
