package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

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
	mux.HandleFunc("/api/dashboard/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handler.GetDashboardStats(w, r)
	})

	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// /api/posts/ = list posts
		if path == "/api/posts/" {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			handler.GetPosts(w, r)
			return
		}

		switch {
		case strings.HasSuffix(path, "/lead"):
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			handler.GetLeadByPostID(w, r)
			return

		case strings.HasSuffix(path, "/lead/status"):
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			handler.UpdateLeadStatus(w, r)
			return

		case strings.HasSuffix(path, "/lead/notes"):
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			handler.UpdateLeadNotes(w, r)
			return

		default:
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			handler.GetPostByID(w, r)
		}
	})

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler.GetDashboardStats(w, r)
	})

	port := ":8080"

	log.Printf("Threads Lead Radar running on http://localhost%s", port)

	if err := http.ListenAndServe(port, enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
