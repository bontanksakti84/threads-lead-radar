package leads

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetPosts(ctx context.Context) ([]Post, error) {
	return s.repository.GetPosts(ctx)
}

func (s *Service) GetPostByID(ctx context.Context, id string) (*Post, error) {
	return s.repository.GetPostByID(ctx, id)
}

func (s *Service) GetPostByExternalID(
	ctx context.Context,
	externalID string,
) (*Post, error) {
	return s.repository.GetPostByExternalID(ctx, externalID)
}

func (s *Service) CreatePost(ctx context.Context, post Post) error {
	return s.repository.CreatePost(ctx, post)
}

func (s *Service) MergeMatchedKeywords(
	ctx context.Context,
	externalID string,
	keywords []string,
) error {
	return s.repository.MergeMatchedKeywords(
		ctx,
		externalID,
		keywords,
	)
}
