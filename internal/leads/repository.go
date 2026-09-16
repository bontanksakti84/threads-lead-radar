package leads

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

type Post struct {
	ID              string   `json:"id"`
	ExternalID      string   `json:"external_id"`
	Username        string   `json:"username"`
	URL             string   `json:"post_url"`
	Content         string   `json:"content"`
	Category        string   `json:"category"`
	Intent          string   `json:"intent"`
	LeadScore       int      `json:"lead_score"`
	HasBudget       bool     `json:"has_budget"`
	HasUrgency      bool     `json:"has_urgency"`
	NeedsDeveloper  bool     `json:"needs_developer"`
	AISummary       string   `json:"ai_summary"`
	MatchedKeywords []string `json:"matched_keywords"`
}

type DashboardStats struct {
	TotalPosts     int `json:"total_posts"`
	Hot            int `json:"hot"`
	Warm           int `json:"warm"`
	Potential      int `json:"potential"`
	Low            int `json:"low"`
	HasBudget      int `json:"has_budget"`
	HasUrgency     int `json:"has_urgency"`
	NeedsDeveloper int `json:"needs_developer"`
}

type PostFilter struct {
	Page           int
	Limit          int
	Search         string
	Tier           string
	Category       string
	Intent         string
	HasBudget      *bool
	HasUrgency     *bool
	NeedsDeveloper *bool
	MinScore       int
}

type PaginatedPosts struct {
	Data       []Post `json:"data"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
}

type Lead struct {
	ID            string `json:"id"`
	PostID        string `json:"post_id"`
	Status        string `json:"status"`
	ContactStatus string `json:"contact_status"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type PostWithLead struct {
	ID              string   `json:"id"`
	ExternalID      string   `json:"external_id"`
	Username        string   `json:"username"`
	URL             string   `json:"post_url"`
	Content         string   `json:"content"`
	Category        string   `json:"category"`
	Intent          string   `json:"intent"`
	LeadScore       int      `json:"lead_score"`
	HasBudget       bool     `json:"has_budget"`
	HasUrgency      bool     `json:"has_urgency"`
	NeedsDeveloper  bool     `json:"needs_developer"`
	AISummary       string   `json:"ai_summary"`
	MatchedKeywords []string `json:"matched_keywords"`
	Lead            *Lead    `json:"lead"`
}

type PaginatedPostsWithLead struct {
	Data       []PostWithLead `json:"data"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

func (r *Repository) CreatePost(ctx context.Context, post Post) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var postID string

	err = tx.QueryRow(
		ctx,
		`
		insert into public.threads_posts
			(
				external_id,
				username,
				post_url,
				content,
				category,
				intent,
				lead_score,
				has_budget,
				has_urgency,
				needs_developer,
				ai_summary,
				matched_keywords
			)
		values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		on conflict (external_id)
		do update set
			username = excluded.username,
			post_url = excluded.post_url,
			content = excluded.content,
			category = excluded.category,
			intent = excluded.intent,
			lead_score = excluded.lead_score,
			has_budget = excluded.has_budget,
			has_urgency = excluded.has_urgency,
			needs_developer = excluded.needs_developer,
			ai_summary = excluded.ai_summary,
			matched_keywords = excluded.matched_keywords
		returning id
		`,
		post.ExternalID,
		post.Username,
		post.URL,
		post.Content,
		post.Category,
		post.Intent,
		post.LeadScore,
		post.HasBudget,
		post.HasUrgency,
		post.NeedsDeveloper,
		post.AISummary,
		post.MatchedKeywords,
	).Scan(&postID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		insert into public.leads (
			post_id,
			status,
			contact_status
		)
		values (
			$1,
			'new',
			'unknown'
		)
		on conflict (post_id) do nothing
		`,
		postID,
	)

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPosts(
	ctx context.Context,
	filter PostFilter,
) (PaginatedPosts, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 20
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit

	where := "where 1=1"
	args := make([]any, 0)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")

		where += fmt.Sprintf(
			` and (
				content ilike $%d
				or username ilike $%d
				or ai_summary ilike $%d
			)`,
			len(args),
			len(args),
			len(args),
		)
	}

	if filter.Tier != "" {
		switch filter.Tier {
		case "hot":
			where += " and lead_score >= 80"

		case "warm":
			where += " and lead_score >= 60 and lead_score < 80"

		case "potential":
			where += " and lead_score >= 40 and lead_score < 60"

		case "low":
			where += " and lead_score < 40"
		}
	}

	if filter.Category != "" {
		args = append(args, filter.Category)

		where += fmt.Sprintf(
			" and category = $%d",
			len(args),
		)
	}

	if filter.Intent != "" {
		args = append(args, filter.Intent)

		where += fmt.Sprintf(
			" and intent = $%d",
			len(args),
		)
	}

	if filter.HasBudget != nil {
		args = append(args, *filter.HasBudget)

		where += fmt.Sprintf(
			" and has_budget = $%d",
			len(args),
		)
	}

	if filter.HasUrgency != nil {
		args = append(args, *filter.HasUrgency)

		where += fmt.Sprintf(
			" and has_urgency = $%d",
			len(args),
		)
	}

	if filter.NeedsDeveloper != nil {
		args = append(args, *filter.NeedsDeveloper)

		where += fmt.Sprintf(
			" and needs_developer = $%d",
			len(args),
		)
	}

	if filter.MinScore > 0 {
		args = append(args, filter.MinScore)

		where += fmt.Sprintf(
			" and lead_score >= $%d",
			len(args),
		)
	}

	var total int

	countQuery := `
		select count(*)
		from public.threads_posts
		` + where

	if err := r.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return PaginatedPosts{}, err
	}

	limitArg := len(args) + 1
	offsetArg := len(args) + 2

	query := fmt.Sprintf(`
		select
			id,
			external_id,
			coalesce(username, ''),
			coalesce(post_url, ''),
			content,
			coalesce(category, ''),
			coalesce(intent, ''),
			coalesce(lead_score, 0),
			coalesce(has_budget, false),
			coalesce(has_urgency, false),
			coalesce(needs_developer, false),
			coalesce(ai_summary, ''),
			coalesce(matched_keywords, '{}')
		from public.threads_posts

		%s

		order by lead_score desc, created_at desc

		limit $%d
		offset $%d
	`, where, limitArg, offsetArg)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return PaginatedPosts{}, err
	}
	defer rows.Close()

	posts := make([]Post, 0)

	for rows.Next() {
		var post Post

		if err := rows.Scan(
			&post.ID,
			&post.ExternalID,
			&post.Username,
			&post.URL,
			&post.Content,
			&post.Category,
			&post.Intent,
			&post.LeadScore,
			&post.HasBudget,
			&post.HasUrgency,
			&post.NeedsDeveloper,
			&post.AISummary,
			&post.MatchedKeywords,
		); err != nil {
			return PaginatedPosts{}, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return PaginatedPosts{}, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + filter.Limit - 1) / filter.Limit
	}

	return PaginatedPosts{
		Data:       posts,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *Repository) GetPostsWithLead(
	ctx context.Context,
	filter PostFilter,
) (PaginatedPostsWithLead, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 20
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit

	where := "where 1=1"
	args := make([]any, 0)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")

		where += fmt.Sprintf(
			` and (
				p.content ilike $%d
				or p.username ilike $%d
				or p.ai_summary ilike $%d
			)`,
			len(args),
			len(args),
			len(args),
		)
	}

	if filter.Tier != "" {
		switch filter.Tier {
		case "hot":
			where += " and p.lead_score >= 80"
		case "warm":
			where += " and p.lead_score >= 60 and p.lead_score < 80"
		case "potential":
			where += " and p.lead_score >= 40 and p.lead_score < 60"
		case "low":
			where += " and p.lead_score < 40"
		}
	}

	if filter.Category != "" {
		args = append(args, filter.Category)
		where += fmt.Sprintf(
			" and p.category = $%d",
			len(args),
		)
	}

	if filter.Intent != "" {
		args = append(args, filter.Intent)
		where += fmt.Sprintf(
			" and p.intent = $%d",
			len(args),
		)
	}

	if filter.HasBudget != nil {
		args = append(args, *filter.HasBudget)
		where += fmt.Sprintf(
			" and p.has_budget = $%d",
			len(args),
		)
	}

	if filter.HasUrgency != nil {
		args = append(args, *filter.HasUrgency)
		where += fmt.Sprintf(
			" and p.has_urgency = $%d",
			len(args),
		)
	}

	if filter.NeedsDeveloper != nil {
		args = append(args, *filter.NeedsDeveloper)
		where += fmt.Sprintf(
			" and p.needs_developer = $%d",
			len(args),
		)
	}

	if filter.MinScore > 0 {
		args = append(args, filter.MinScore)

		where += fmt.Sprintf(
			" and p.lead_score >= $%d",
			len(args),
		)
	}

	var total int

	countQuery := `
		select count(*)
		from public.threads_posts p
		` + where

	if err := r.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return PaginatedPostsWithLead{}, err
	}

	limitArg := len(args) + 1
	offsetArg := len(args) + 2

	query := fmt.Sprintf(`
		select
			p.id,
			p.external_id,
			coalesce(p.username, ''),
			coalesce(p.post_url, ''),
			p.content,
			coalesce(p.category, ''),
			coalesce(p.intent, ''),
			coalesce(p.lead_score, 0),
			coalesce(p.has_budget, false),
			coalesce(p.has_urgency, false),
			coalesce(p.needs_developer, false),
			coalesce(p.ai_summary, ''),
			coalesce(p.matched_keywords, '{}'),

			l.id,
			l.post_id,
			coalesce(l.status, 'new'),
			coalesce(l.contact_status, 'unknown'),
			coalesce(l.notes, ''),
			l.created_at::text,
			l.updated_at::text

		from public.threads_posts p

		left join public.leads l
			on l.post_id = p.id

		%s

		order by p.lead_score desc, p.created_at desc

		limit $%d
		offset $%d
	`, where, limitArg, offsetArg)

	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return PaginatedPostsWithLead{}, err
	}

	defer rows.Close()

	posts := make([]PostWithLead, 0)

	for rows.Next() {
		var post PostWithLead
		var leadID *string
		var leadPostID *string
		var leadStatus *string
		var leadContactStatus *string
		var leadNotes *string
		var leadCreatedAt *string
		var leadUpdatedAt *string

		if err := rows.Scan(
			&post.ID,
			&post.ExternalID,
			&post.Username,
			&post.URL,
			&post.Content,
			&post.Category,
			&post.Intent,
			&post.LeadScore,
			&post.HasBudget,
			&post.HasUrgency,
			&post.NeedsDeveloper,
			&post.AISummary,
			&post.MatchedKeywords,

			&leadID,
			&leadPostID,
			&leadStatus,
			&leadContactStatus,
			&leadNotes,
			&leadCreatedAt,
			&leadUpdatedAt,
		); err != nil {
			return PaginatedPostsWithLead{}, err
		}

		if leadID != nil {
			post.Lead = &Lead{
				ID:            *leadID,
				PostID:        *leadPostID,
				Status:        *leadStatus,
				ContactStatus: *leadContactStatus,
				Notes:         *leadNotes,
				CreatedAt:     *leadCreatedAt,
				UpdatedAt:     *leadUpdatedAt,
			}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return PaginatedPostsWithLead{}, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + filter.Limit - 1) / filter.Limit
	}

	return PaginatedPostsWithLead{
		Data:       posts,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *Repository) GetPostByID(ctx context.Context, id string) (*Post, error) {
	var post Post

	err := r.db.QueryRow(
		ctx,
		`
		select
			id,
			external_id,
			coalesce(username, ''),
			coalesce(post_url, ''),
			content,
			coalesce(category, ''),
			coalesce(intent, ''),
			coalesce(lead_score, 0),
			coalesce(has_budget, false),
			coalesce(has_urgency, false),
			coalesce(needs_developer, false),
			coalesce(ai_summary, ''),
			coalesce(matched_keywords, '{}')
		from public.threads_posts
		where id = $1
		`,
		id,
	).Scan(
		&post.ID,
		&post.ExternalID,
		&post.Username,
		&post.URL,
		&post.Content,
		&post.Category,
		&post.Intent,
		&post.LeadScore,
		&post.HasBudget,
		&post.HasUrgency,
		&post.NeedsDeveloper,
		&post.AISummary,
		&post.MatchedKeywords,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) GetPostByExternalID(ctx context.Context, externalID string) (*Post, error) {
	var post Post

	err := r.db.QueryRow(
		ctx,
		`
		select
			id,
			external_id,
			coalesce(username, ''),
			coalesce(post_url, ''),
			content,
			coalesce(category, ''),
			coalesce(intent, ''),
			coalesce(lead_score, 0),
			coalesce(has_budget, false),
			coalesce(has_urgency, false),
			coalesce(needs_developer, false),
			coalesce(ai_summary, ''),
			coalesce(matched_keywords, '{}')
		from public.threads_posts
		where external_id = $1
		`,
		externalID,
	).Scan(
		&post.ID,
		&post.ExternalID,
		&post.Username,
		&post.URL,
		&post.Content,
		&post.Category,
		&post.Intent,
		&post.LeadScore,
		&post.HasBudget,
		&post.HasUrgency,
		&post.NeedsDeveloper,
		&post.AISummary,
		&post.MatchedKeywords,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) MergeMatchedKeywords(
	ctx context.Context,
	externalID string,
	keywords []string,
) error {
	if len(keywords) == 0 {
		return nil
	}

	_, err := r.db.Exec(
		ctx,
		`
		update public.threads_posts
		set matched_keywords = (
			select coalesce(array_agg(distinct keyword order by keyword), '{}')
			from unnest(
				coalesce(matched_keywords, '{}') ||
				$2::text[]
			) as keyword
		)
		where external_id = $1
		`,
		externalID,
		keywords,
	)

	return err
}

func (r *Repository) GetDashboardStats(
	ctx context.Context,
) (DashboardStats, error) {
	var stats DashboardStats

	err := r.db.QueryRow(
		ctx,
		`
		select
			count(*) as total_posts,

			count(*) filter (
				where lead_score >= 80
			) as hot,

			count(*) filter (
				where lead_score >= 60
				and lead_score < 80
			) as warm,

			count(*) filter (
				where lead_score >= 40
				and lead_score < 60
			) as potential,

			count(*) filter (
				where lead_score < 40
			) as low,

			count(*) filter (
				where has_budget = true
			) as has_budget,

			count(*) filter (
				where has_urgency = true
			) as has_urgency,

			count(*) filter (
				where needs_developer = true
			) as needs_developer

		from public.threads_posts
		`,
	).Scan(
		&stats.TotalPosts,
		&stats.Hot,
		&stats.Warm,
		&stats.Potential,
		&stats.Low,
		&stats.HasBudget,
		&stats.HasUrgency,
		&stats.NeedsDeveloper,
	)

	if err != nil {
		return DashboardStats{}, err
	}

	return stats, nil
}

func (r *Repository) GetLeadByPostID(
	ctx context.Context,
	postID string,
) (*Lead, error) {
	var lead Lead

	err := r.db.QueryRow(
		ctx,
		`
		select
			id,
			post_id,
			status,
			contact_status,
			coalesce(notes, ''),
			created_at::text,
			updated_at::text
		from public.leads
		where post_id = $1
		limit 1
		`,
		postID,
	).Scan(
		&lead.ID,
		&lead.PostID,
		&lead.Status,
		&lead.ContactStatus,
		&lead.Notes,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &lead, nil
}

func (r *Repository) UpdateLeadStatus(
	ctx context.Context,
	postID string,
	status string,
) error {
	result, err := r.db.Exec(
		ctx,
		`
		update public.leads
		set
			status = $1,
			updated_at = now()
		where post_id = $2
		`,
		status,
		postID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("lead not found")
	}

	return nil
}

func (r *Repository) UpdateLeadNotes(
	ctx context.Context,
	postID string,
	notes string,
) error {
	result, err := r.db.Exec(
		ctx,
		`
		update public.leads
		set
			notes = $1,
			updated_at = now()
		where post_id = $2
		`,
		notes,
		postID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("lead not found")
	}

	return nil
}
