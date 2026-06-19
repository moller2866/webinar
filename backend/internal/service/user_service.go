package service

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/webinar/backend/internal/model"
	"github.com/webinar/backend/internal/repository"
)

const bcryptCost = 12

type UserService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
}

// Register creates a new user account, hashing the password with bcrypt.
// Returns the created user or an error. Duplicate email yields *model.ConflictError.
func (s *UserService) Register(displayName, email, password string) (*model.User, error) {
	displayName = strings.TrimSpace(displayName)
	email = strings.ToLower(strings.TrimSpace(email))

	if displayName == "" {
		return nil, &model.ValidationError{Message: "display name is required"}
	}
	if len(displayName) < 2 || len(displayName) > 30 {
		return nil, &model.ValidationError{Message: "display name must be between 2 and 30 characters"}
	}
	if email == "" {
		return nil, &model.ValidationError{Message: "email is required"}
	}
	if len(password) < 8 {
		return nil, &model.ValidationError{Message: "password must be at least 8 characters"}
	}

	existing, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &model.ConflictError{Message: "an account with that email already exists"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login validates credentials and returns the user on success.
// Returns nil, nil when the email is not found or the password is wrong
// (callers should treat this as a generic 401 — no enumeration).
func (s *UserService) Login(email, password string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// Perform a dummy bcrypt comparison to avoid timing-based email enumeration.
		bcrypt.CompareHashAndPassword([]byte("$2a$12$dummy"), []byte(password)) //nolint:errcheck
		return nil, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil
	}
	return user, nil
}
