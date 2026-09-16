package scoring

import (
	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
)

type Input struct {
	Classification ai.Classification
	KeywordMatches int
}

func Calculate(input Input) int {
	score := 0
	result := input.Classification

	if result.NeedsDeveloper {
		score += 25
	}

	switch result.Intent {
	case "hire_developer":
		score += 25

	case "looking_for_jasa":
		score += 20

	case "project_inquiry":
		score += 20
	}

	if result.HasBudget {
		score += 10
	}

	if result.HasUrgency {
		score += 10
	}

	// Keyword matching is a supporting signal,
	// not the primary determinant of lead quality.
	if input.KeywordMatches > 0 {
		score += 10
	}

	if score > 100 {
		score = 100
	}

	return score
}
