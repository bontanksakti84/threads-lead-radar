package leads

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScanRun struct {
	ID              string
	StartedAt       time.Time
	FinishedAt      *time.Time
	KeywordCount    int
	PostsFound      int
	CandidatesFound int
	NewLeads        int
	ExistingLeads   int
	GroqCalls       int
	DurationMs      int64
	Status          string
	ErrorMessage    string
}

type ScanRunRepository struct {
	db *pgxpool.Pool
}

func NewScanRunRepository(db *pgxpool.Pool) *ScanRunRepository {
	return &ScanRunRepository{db: db}
}

func (r *ScanRunRepository) Start(
	ctx context.Context,
	keywordCount int,
) (string, error) {
	var id string

	err := r.db.QueryRow(
		ctx,
		`
		insert into public.scan_runs (
			keyword_count,
			status
		)
		values ($1, 'running')
		returning id
		`,
		keywordCount,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *ScanRunRepository) Complete(
	ctx context.Context,
	id string,
	postsFound int,
	candidatesFound int,
	newLeads int,
	existingLeads int,
	groqCalls int,
	durationMs int64,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		update public.scan_runs
		set
			finished_at = now(),
			posts_found = $2,
			leads_found = $3,
			candidates_found = $4,
			new_leads = $5,
			existing_leads = $6,
			groq_calls = $7,
			duration_ms = $8,
			status = 'completed'
		where id = $1
		`,
		id,
		postsFound,
		newLeads,
		candidatesFound,
		newLeads,
		existingLeads,
		groqCalls,
		durationMs,
	)

	return err
}

func (r *ScanRunRepository) Fail(
	ctx context.Context,
	id string,
	errMessage string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		update public.scan_runs
		set
			finished_at = now(),
			status = 'failed',
			error_message = $2
		where id = $1
		`,
		id,
		errMessage,
	)

	return err
}
