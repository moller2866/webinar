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
	dbPath := "blog.db"
	if v := os.Getenv("DB_PATH"); v != "" {
		dbPath = v
	}

	db, err := repository.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	postRepo := repository.NewSQLitePostRepository(db)
	commentRepo := repository.NewSQLiteCommentRepository(db)

	postService := service.NewPostService(postRepo, commentRepo)

	h := handler.NewHandler(postService)
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
