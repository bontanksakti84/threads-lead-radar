package ai

import (
	"context"
	"testing"
)

func TestMockClassifier(t *testing.T) {
	classifier := NewMockClassifier()

	content := "Saya sedang mencari developer untuk membuat website bisnis saya."

	result, err := classifier.Classify(
		context.Background(),
		content,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsLead {
		t.Fatal("expected lead")
	}

	if !result.NeedsDeveloper {
		t.Fatal("expected needs_developer=true")
	}

	if result.Category != "website" {
		t.Fatalf(
			"expected category website, got %s",
			result.Category,
		)
	}

	if result.Intent != "hire_developer" {
		t.Fatalf(
			"expected intent hire_developer, got %s",
			result.Intent,
		)
	}

	if result.LeadScore <= 0 {
		t.Fatal("expected positive lead score")
	}
}

func TestMockClassifierNonLead(t *testing.T) {
	classifier := NewMockClassifier()

	content := "Ada rekomendasi laptop bagus untuk coding?"

	result, err := classifier.Classify(
		context.Background(),
		content,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsLead {
		t.Fatal("expected non-lead")
	}

	if result.LeadScore <= 0 {
		t.Fatal("expected score")
	}
}

func TestMockClassifierProjectInquiry(t *testing.T) {
	classifier := NewMockClassifier()

	content := "Ada yang bisa bantu bikin website untuk bisnis saya?"

	result, err := classifier.Classify(
		context.Background(),
		content,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsLead {
		t.Fatal("expected lead")
	}

	if result.Category != "website" {
		t.Fatalf(
			"expected category website, got %s",
			result.Category,
		)
	}

	if result.Intent != "project_inquiry" {
		t.Fatalf(
			"expected intent project_inquiry, got %s",
			result.Intent,
		)
	}

	if result.LeadScore < 60 {
		t.Fatalf(
			"expected lead score >= 60, got %d",
			result.LeadScore,
		)
	}
}
