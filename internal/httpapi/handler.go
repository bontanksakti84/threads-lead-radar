package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
)

type Handler struct {
	leads *leads.Service
}

func NewHandler(leadsService *leads.Service) *Handler {
	return &Handler{
		leads: leadsService,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.leads.GetPosts(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/posts/")

	if id == "" {
		writeError(w, http.StatusBadRequest, "post id is required")
		return
	}

	post, err := h.leads.GetPostByID(r.Context(), id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if post == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post leads.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if post.ExternalID == "" {
		writeError(w, http.StatusBadRequest, "external_id is required")
		return
	}

	if post.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	if err := h.leads.CreatePost(r.Context(), post); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
