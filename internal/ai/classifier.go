package ai

import (
	"context"
	"fmt"
	"strings"
)

type Classification struct {
	IsLead         bool   `json:"is_lead"`
	Category       string `json:"category"`
	Intent         string `json:"intent"`
	LeadScore      int    `json:"lead_score"`
	HasBudget      bool   `json:"has_budget"`
	HasUrgency     bool   `json:"has_urgency"`
	NeedsDeveloper bool   `json:"needs_developer"`
	Summary        string `json:"summary"`
}

type Classifier interface {
	Classify(ctx context.Context, content string) (Classification, error)
}

func ValidateClassification(result Classification) error {
	if result.LeadScore < 0 {
		return fmt.Errorf("lead_score cannot be below 0")
	}

	if result.LeadScore > 100 {
		return fmt.Errorf("lead_score cannot exceed 100")
	}

	allowedCategories := map[string]bool{
		"website":         true,
		"mobile_app":      true,
		"custom_software": true,
		"erp":             true,
		"hrm":             true,
		"pos":             true,
		"automation":      true,
		"ai":              true,
		"ecommerce":       true,
		"maintenance":     true,
		"other":           true,
	}

	if !allowedCategories[result.Category] {
		return fmt.Errorf("invalid category: %s", result.Category)
	}

	allowedIntents := map[string]bool{
		"hire_developer":   true,
		"looking_for_jasa": true,
		"project_inquiry":  true,
		"recommendation":   true,
		"product_research": true,
		"general_question": true,
		"other":            true,
	}

	if !allowedIntents[result.Intent] {
		return fmt.Errorf("invalid intent: %s", result.Intent)
	}

	result.Summary = strings.TrimSpace(result.Summary)

	return nil
}
