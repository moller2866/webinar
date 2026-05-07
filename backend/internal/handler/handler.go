package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/webinar/backend/internal/model"
	"github.com/webinar/backend/internal/service"
)

// Request/response DTOs

type CreatePostRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	CategoryID *int64   `json:"categoryId"`
	Images     []string `json:"images"`
}

type CreateCommentRequest struct {
	Content  string `json:"content"`
	ParentID *int64 `json:"parentId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Handler wires HTTP routes to the PostService.
type Handler struct {
	postService *service.PostService
}

func NewHandler(postService *service.PostService) *Handler {
	return &Handler{postService: postService}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/posts", h.listPosts)
	mux.HandleFunc("POST /api/posts", h.createPost)
	mux.HandleFunc("GET /api/posts/{id}", h.getPost)
	mux.HandleFunc("POST /api/posts/{id}/comments", h.addComment)
}

// --- Handlers ---

func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.ListPosts()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list posts")
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post := &model.Post{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		Images:     req.Images,
	}

	if err := h.postService.CreatePost(post); err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.postService.GetPost(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) addComment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comment := &model.Comment{
		PostID:   id,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := h.postService.AddComment(comment); err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

// --- Helpers ---

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// handleServiceError maps service-layer errors to HTTP status codes.
func handleServiceError(w http.ResponseWriter, err error) {
	var ve *model.ValidationError
	if errors.As(err, &ve) {
		if strings.Contains(ve.Message, "not found") {
			writeError(w, http.StatusNotFound, ve.Message)
		} else {
			writeError(w, http.StatusBadRequest, ve.Message)
		}
		return
	}
	writeError(w, http.StatusInternalServerError, "internal server error")
}
