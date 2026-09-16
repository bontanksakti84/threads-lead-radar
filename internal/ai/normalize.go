package ai

func NormalizeClassification(result Classification) Classification {
	switch result.Intent {
	case "hire_developer":
		result.IsLead = true
		result.NeedsDeveloper = true

	case "looking_for_jasa":
		result.IsLead = true

	case "project_inquiry":
		result.IsLead = true
	}

	if result.LeadScore < 0 {
		result.LeadScore = 0
	}

	if result.LeadScore > 100 {
		result.LeadScore = 100
	}

	return result
}
