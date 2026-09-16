package threads

import (
	"context"
	"testing"
)

func TestMockClientSearchPosts(t *testing.T) {
	client := NewMockClient()

	posts, err := client.SearchPosts(
		context.Background(),
		"website",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Fatalf(
			"expected 1 post, got %d",
			len(posts),
		)
	}

	if posts[0].ID != "mock-001" {
		t.Fatalf(
			"expected mock-001, got %s",
			posts[0].ID,
		)
	}
}

func TestMockClientSearchLandingPage(t *testing.T) {
	client := NewMockClient()

	posts, err := client.SearchPosts(
		context.Background(),
		"landing page",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Fatalf(
			"expected 1 post, got %d",
			len(posts),
		)
	}

	if posts[0].ID != "mock-004" {
		t.Fatalf(
			"expected mock-004, got %s",
			posts[0].ID,
		)
	}
}

func TestMockClientSearchNoResults(t *testing.T) {
	client := NewMockClient()

	posts, err := client.SearchPosts(
		context.Background(),
		"xyz-does-not-exist",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 0 {
		t.Fatalf(
			"expected 0 posts, got %d",
			len(posts),
		)
	}
}
