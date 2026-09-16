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

func (s *Service) GetPosts(
	ctx context.Context,
	filter PostFilter,
) (PaginatedPosts, error) {
	return s.repository.GetPosts(ctx, filter)
}

func (s *Service) GetPostsWithLead(
	ctx context.Context,
	filter PostFilter,
) (PaginatedPostsWithLead, error) {
	return s.repository.GetPostsWithLead(ctx, filter)
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

func (s *Service) GetDashboardStats(
	ctx context.Context,
) (DashboardStats, error) {
	return s.repository.GetDashboardStats(ctx)
}

func (s *Service) GetLeadByPostID(
	ctx context.Context,
	postID string,
) (*Lead, error) {
	return s.repository.GetLeadByPostID(ctx, postID)
}

func (s *Service) UpdateLeadStatus(
	ctx context.Context,
	postID string,
	status string,
) error {
	return s.repository.UpdateLeadStatus(
		ctx,
		postID,
		status,
	)
}

func (s *Service) UpdateLeadNotes(
	ctx context.Context,
	postID string,
	notes string,
) error {
	return s.repository.UpdateLeadNotes(
		ctx,
		postID,
		notes,
	)
}
