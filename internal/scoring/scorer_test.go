package scoring

import (
	"testing"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
)

func TestCalculateHireDeveloper(t *testing.T) {
	result := ai.Classification{
		NeedsDeveloper: true,
		Intent:         "hire_developer",
		HasBudget:      false,
		HasUrgency:     false,
	}

	score := Calculate(Input{
		Classification: result,
		KeywordMatches: 1,
	})

	expected := 60

	if score != expected {
		t.Fatalf(
			"expected score %d, got %d",
			expected,
			score,
		)
	}
}

func TestCalculateProjectInquiry(t *testing.T) {
	result := ai.Classification{
		NeedsDeveloper: true,
		Intent:         "project_inquiry",
		HasBudget:      true,
		HasUrgency:     true,
	}

	score := Calculate(Input{
		Classification: result,
		KeywordMatches: 1,
	})

	expected := 75

	if score != expected {
		t.Fatalf(
			"expected score %d, got %d",
			expected,
			score,
		)
	}
}

func TestCalculateNoLeadSignals(t *testing.T) {
	result := ai.Classification{
		NeedsDeveloper: false,
		Intent:         "general_question",
		HasBudget:      false,
		HasUrgency:     false,
	}

	score := Calculate(Input{
		Classification: result,
		KeywordMatches: 0,
	})

	expected := 0

	if score != expected {
		t.Fatalf(
			"expected score %d, got %d",
			expected,
			score,
		)
	}
}
