package main

import (
	"log"
	"net/http"
	"os"

	"github.com/webinar/backend/internal/handler"
	"github.com/webinar/backend/internal/repository"
	"github.com/webinar/backend/internal/service"
)

func main() {
	databaseURL := "postgres://blog:blog@localhost:5432/blog?sslmode=disable"
	if v := os.Getenv("DATABASE_URL"); v != "" {
		databaseURL = v
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	db, err := repository.NewPostgresDB(databaseURL)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	postRepo := repository.NewPostgresPostRepository(db)
	commentRepo := repository.NewPostgresCommentRepository(db)
	userRepo := repository.NewPostgresUserRepository(db)

	postService := service.NewPostService(postRepo, commentRepo)
	userService := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(userService, []byte(jwtSecret))
	h := handler.NewHandler(postService, authHandler, []byte(jwtSecret))
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	corsHandler := corsMiddleware(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
