package scoring

import "testing"

func TestGetTier(t *testing.T) {
	tests := []struct {
		score    int
		expected Tier
	}{
		{100, TierHot},
		{80, TierHot},
		{79, TierWarm},
		{60, TierWarm},
		{59, TierPotential},
		{40, TierPotential},
		{39, TierLow},
		{0, TierLow},
	}

	for _, tt := range tests {
		result := GetTier(tt.score)

		if result != tt.expected {
			t.Fatalf(
				"score %d: expected %s, got %s",
				tt.score,
				tt.expected,
				result,
			)
		}
	}
}
