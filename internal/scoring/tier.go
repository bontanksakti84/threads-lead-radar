package scoring

type Tier string

const (
	TierHot       Tier = "hot"
	TierWarm      Tier = "warm"
	TierPotential Tier = "potential"
	TierLow       Tier = "low"
)

func GetTier(score int) Tier {
	switch {
	case score >= 80:
		return TierHot

	case score >= 60:
		return TierWarm

	case score >= 40:
		return TierPotential

	default:
		return TierLow
	}
}
