package threads

import (
	"context"
	"strings"
	"time"
)

type MockClient struct {
	Posts []Post
}

func NewMockClient() *MockClient {
	return &MockClient{
		Posts: []Post{
			{
				ID:       "mock-001",
				Username: "owner_umkm",
				URL:      "https://www.threads.net/",
				Content:  "Ada yang bisa bantu bikin website untuk bisnis saya? Saya butuh developer.",
				PublishedAt: time.Now().Add(
					-10 * time.Minute,
				),
			},
			{
				ID:       "mock-002",
				Username: "coffee_owner",
				URL:      "https://www.threads.net/",
				Content:  "Saya sedang mencari developer untuk aplikasi kasir dan inventory.",
				PublishedAt: time.Now().Add(
					-20 * time.Minute,
				),
			},
			{
				ID:       "mock-003",
				Username: "random_user",
				URL:      "https://www.threads.net/",
				Content:  "Ada rekomendasi laptop untuk coding?",
				PublishedAt: time.Now().Add(
					-30 * time.Minute,
				),
			},
			{
				ID:       "mock-004",
				Username: "business_owner",
				URL:      "https://www.threads.net/",
				Content:  "Saya mau bikin landing page untuk bisnis baru saya.",
				PublishedAt: time.Now().Add(
					-40 * time.Minute,
				),
			},
			{
				ID:       "mock-005",
				Username: "hr_owner",
				URL:      "https://www.threads.net/",
				Content:  "Ada yang bisa bantu bikin sistem HRIS untuk perusahaan kami?",
				PublishedAt: time.Now().Add(
					-50 * time.Minute,
				),
			},
			{
				ID:       "mock-006",
				Username: "automation_owner",
				URL:      "https://www.threads.net/",
				Content:  "Butuh bantuan automation Google Apps Script untuk laporan perusahaan.",
				PublishedAt: time.Now().Add(
					-60 * time.Minute,
				),
			},
		},
	}
}

func (m *MockClient) SearchPosts(
	ctx context.Context,
	keyword string,
) ([]Post, error) {
	keyword = strings.ToLower(
		strings.TrimSpace(keyword),
	)

	var results []Post

	for _, post := range m.Posts {
		content := strings.ToLower(post.Content)

		if strings.Contains(content, keyword) {
			results = append(results, post)
		}
	}

	return results, nil
}
