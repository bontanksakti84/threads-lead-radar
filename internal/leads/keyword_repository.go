package leads

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Keyword struct {
	ID       string
	Keyword  string
	Category string
	Active   bool
}

type KeywordRepository struct {
	db *pgxpool.Pool
}

func NewKeywordRepository(db *pgxpool.Pool) *KeywordRepository {
	return &KeywordRepository{
		db: db,
	}
}

func (r *KeywordRepository) GetActiveKeywords(
	ctx context.Context,
) ([]Keyword, error) {
	rows, err := r.db.Query(
		ctx,
		`
		select
			id,
			keyword,
			coalesce(category, ''),
			is_active
		from public.keywords
		where is_active = true
		order by keyword
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var keywords []Keyword

	for rows.Next() {
		var keyword Keyword

		if err := rows.Scan(
			&keyword.ID,
			&keyword.Keyword,
			&keyword.Category,
			&keyword.Active,
		); err != nil {
			return nil, err
		}

		keywords = append(keywords, keyword)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}