package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/bontanksakti84/threads-lead-radar/internal/database"
	"github.com/bontanksakti84/threads-lead-radar/internal/httpapi"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	db, err := database.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Pool.Close()

	log.Println("Database connected successfully")

	repository := leads.NewRepository(db.Pool)
	leadService := leads.NewService(repository)
	handler := httpapi.NewHandler(leadService)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.Health)

	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetPosts(w, r)

		case http.MethodPost:
			handler.CreatePost(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handler.GetPostByID(w, r)
	})

	port := ":8080"

	log.Printf("Threads Lead Radar running on http://localhost%s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
