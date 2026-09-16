package leads

import (
	"context"
	"errors"

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

func (r *Repository) CreatePost(ctx context.Context, post Post) error {
	_, err := r.db.Exec(
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
	)

	return err
}

func (r *Repository) GetPosts(ctx context.Context) ([]Post, error) {
	rows, err := r.db.Query(
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
		order by created_at desc
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err := rows.Scan(
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

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
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
