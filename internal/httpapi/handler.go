package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
	"net/http"
	"strconv"
	"strings"
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

func (h *Handler) GetPosts(
	w http.ResponseWriter,
	r *http.Request,
) {
	filter := leads.PostFilter{
		Page:     1,
		Limit:    20,
		Search:   r.URL.Query().Get("search"),
		Tier:     r.URL.Query().Get("tier"),
		Category: r.URL.Query().Get("category"),
		Intent:   r.URL.Query().Get("intent"),
	}

	if value := r.URL.Query().Get("page"); value != "" {
		page, err := strconv.Atoi(value)

		if err != nil || page < 1 {
			writeError(
				w,
				http.StatusBadRequest,
				"page must be a positive integer",
			)
			return
		}

		filter.Page = page
	}

	if value := r.URL.Query().Get("limit"); value != "" {
		limit, err := strconv.Atoi(value)

		if err != nil || limit < 1 {
			writeError(
				w,
				http.StatusBadRequest,
				"limit must be a positive integer",
			)
			return
		}

		if limit > 100 {
			writeError(
				w,
				http.StatusBadRequest,
				"limit cannot exceed 100",
			)
			return
		}

		filter.Limit = limit
	}

	if value := r.URL.Query().Get("min_score"); value != "" {
		minScore, err := strconv.Atoi(value)

		if err != nil || minScore < 0 || minScore > 100 {
			writeError(
				w,
				http.StatusBadRequest,
				"min_score must be an integer between 0 and 100",
			)
			return
		}

		filter.MinScore = minScore
	}

	if err := validatePostFilters(filter); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	hasBudget, err := parseOptionalBool(
		r,
		"has_budget",
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"has_budget must be true or false",
		)
		return
	}

	filter.HasBudget = hasBudget

	hasUrgency, err := parseOptionalBool(
		r,
		"has_urgency",
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"has_urgency must be true or false",
		)
		return
	}

	filter.HasUrgency = hasUrgency

	needsDeveloper, err := parseOptionalBool(
		r,
		"needs_developer",
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"needs_developer must be true or false",
		)
		return
	}

	filter.NeedsDeveloper = needsDeveloper

	ctx := r.Context()

	result, err := h.leads.GetPostsWithLead(
		ctx,
		filter,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		result,
	)
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

func (h *Handler) GetDashboardStats(
	w http.ResponseWriter,
	r *http.Request,
) {
	stats, err := h.leads.GetDashboardStats(r.Context())
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		stats,
	)
}

func parseOptionalBool(
	r *http.Request,
	name string,
) (*bool, error) {
	value := strings.TrimSpace(
		r.URL.Query().Get(name),
	)

	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func isAllowedValue(
	value string,
	allowed []string,
) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}

	return false
}

func validatePostFilters(
	filter leads.PostFilter,
) error {
	if filter.Tier != "" &&
		!isAllowedValue(
			filter.Tier,
			[]string{
				"hot",
				"warm",
				"potential",
				"low",
			},
		) {
		return fmt.Errorf(
			"invalid tier: %s",
			filter.Tier,
		)
	}

	if filter.Category != "" &&
		!isAllowedValue(
			filter.Category,
			[]string{
				"website",
				"mobile_app",
				"custom_software",
				"erp",
				"hrm",
				"pos",
				"crm",
				"automation",
				"ai",
				"ecommerce",
				"maintenance",
				"other",
			},
		) {
		return fmt.Errorf(
			"invalid category: %s",
			filter.Category,
		)
	}

	if filter.Intent != "" &&
		!isAllowedValue(
			filter.Intent,
			[]string{
				"hire_developer",
				"looking_for_jasa",
				"project_inquiry",
				"recommendation",
				"product_research",
				"general_question",
				"other",
			},
		) {
		return fmt.Errorf(
			"invalid intent: %s",
			filter.Intent,
		)
	}

	return nil
}

var allowedLeadStatuses = map[string]bool{
	"new":       true,
	"contacted": true,
	"replied":   true,
	"qualified": true,
	"won":       true,
	"lost":      true,
}

func isValidLeadStatus(status string) bool {
	return allowedLeadStatuses[status]
}

func (h *Handler) GetLeadByPostID(
	w http.ResponseWriter,
	r *http.Request,
) {
	postID := strings.TrimPrefix(
		r.URL.Path,
		"/api/posts/",
	)

	postID = strings.TrimSuffix(
		postID,
		"/lead",
	)

	if postID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"post id is required",
		)
		return
	}

	lead, err := h.leads.GetLeadByPostID(
		r.Context(),
		postID,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	if lead == nil {
		writeError(
			w,
			http.StatusNotFound,
			"lead not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		lead,
	)
}

type updateLeadStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateLeadStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	postID := strings.TrimPrefix(
		r.URL.Path,
		"/api/posts/",
	)

	postID = strings.TrimSuffix(
		postID,
		"/lead/status",
	)

	if postID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"post id is required",
		)
		return
	}

	var request updateLeadStatusRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON",
		)
		return
	}

	request.Status = strings.ToLower(
		strings.TrimSpace(request.Status),
	)

	if !isValidLeadStatus(request.Status) {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid lead status",
		)
		return
	}

	if err := h.leads.UpdateLeadStatus(
		r.Context(),
		postID,
		request.Status,
	); err != nil {
		if err.Error() == "lead not found" {
			writeError(
				w,
				http.StatusNotFound,
				"lead not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]any{
			"success": true,
			"status":  request.Status,
		},
	)
}

type updateLeadNotesRequest struct {
	Notes string `json:"notes"`
}

func (h *Handler) UpdateLeadNotes(
	w http.ResponseWriter,
	r *http.Request,
) {
	postID := strings.TrimPrefix(
		r.URL.Path,
		"/api/posts/",
	)

	postID = strings.TrimSuffix(
		postID,
		"/lead/notes",
	)

	if postID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"post id is required",
		)
		return
	}

	var request updateLeadNotesRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON",
		)
		return
	}

	if err := h.leads.UpdateLeadNotes(
		r.Context(),
		postID,
		request.Notes,
	); err != nil {
		if err.Error() == "lead not found" {
			writeError(
				w,
				http.StatusNotFound,
				"lead not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]any{
			"success": true,
			"notes":   request.Notes,
		},
	)
}
