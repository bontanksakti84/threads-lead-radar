package threads

import (
	"testing"

	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
)

func TestMatchKeywords(t *testing.T) {
	keywords := []leads.Keyword{
		{Keyword: "mencari developer", Category: "developer", Active: true},
		{Keyword: "website", Category: "website", Active: true},
		{Keyword: "inventory", Category: "inventory", Active: true},
	}

	content := `
		Saya sedang mencari developer untuk
		membuat website dan sistem inventory.
	`

	matches := MatchKeywords(content, keywords)

	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
}
