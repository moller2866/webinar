package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/webinar/backend/internal/model"
	"github.com/webinar/backend/internal/service"
)

// contextKey is an unexported type for context keys in this package.
type contextKey int

const userIDKey contextKey = iota

// UserIDFromContext retrieves the authenticated user ID from the request context.
// Returns 0 and false if not present.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

// AuthRequest is the shared DTO for register and login.
type AuthRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// AuthResponse is the JWT response returned after successful register/login.
type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	userService *service.UserService
	jwtSecret   []byte
}

func NewAuthHandler(userService *service.UserService, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{userService: userService, jwtSecret: jwtSecret}
}

// register handles POST /api/auth/register.
func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Register(req.DisplayName, req.Email, req.Password)
	if err != nil {
		var ve *model.ValidationError
		var ce *model.ConflictError
		switch {
		case errors.As(err, &ve):
			writeError(w, http.StatusBadRequest, ve.Message)
		case errors.As(err, &ce):
			writeError(w, http.StatusConflict, ce.Message)
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, AuthResponse{Token: token, User: user})
}

// login handles POST /api/auth/login.
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if user == nil {
		writeError(w, http.StatusUnauthorized, "email or password is incorrect")
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{Token: token, User: user})
}

// generateToken creates a signed HS256 JWT for the given user ID (7-day expiry).
func (h *AuthHandler) generateToken(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

// AuthMiddleware returns HTTP middleware that validates a Bearer JWT.
// On success it injects the user ID into the request context and calls next.
// On failure it writes a 401 and does not call next.
func AuthMiddleware(jwtSecret []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "authorization required")
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Subject == "" {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
