package threads

import (
	"context"
	"time"
)

type Post struct {
	ID          string
	Username    string
	URL         string
	Content     string
	PublishedAt time.Time
}

type Client interface {
	SearchPosts(ctx context.Context, keyword string) ([]Post, error)
}
