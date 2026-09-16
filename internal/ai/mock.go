package ai

import (
	"context"
	"strings"
)

type MockClassifier struct{}

func NewMockClassifier() *MockClassifier {
	return &MockClassifier{}
}

func (m *MockClassifier) Classify(
	ctx context.Context,
	content string,
) (Classification, error) {

	text := strings.ToLower(content)

	result := Classification{
		IsLead:    false,
		Category:  "other",
		Intent:    "other",
		LeadScore: 10,
		Summary:   "Tidak teridentifikasi sebagai lead.",
	}

	// ---------------------------------
	// Developer / programming need
	// ---------------------------------

	if strings.Contains(text, "developer") ||
		strings.Contains(text, "programmer") {

		result.NeedsDeveloper = true
		result.Category = "custom_software"
		result.LeadScore = 70
		result.Summary = "Pengguna terlihat membutuhkan developer."
	}

	// ---------------------------------
	// Hire developer
	// ---------------------------------

	if strings.Contains(text, "butuh developer") ||
		strings.Contains(text, "cari developer") ||
		strings.Contains(text, "mencari developer") ||
		strings.Contains(text, "butuh programmer") ||
		strings.Contains(text, "cari programmer") {

		result.IsLead = true
		result.Intent = "hire_developer"
		result.LeadScore = 85
		result.Summary = "Pengguna terlihat sedang mencari developer."
	}

	// ---------------------------------
	// Project inquiry
	// ---------------------------------

	if strings.Contains(text, "bisa bantu") ||
		strings.Contains(text, "bantu bikin") ||
		strings.Contains(text, "bantu buat") ||
		strings.Contains(text, "ada yang bisa") {

		result.IsLead = true
		result.Intent = "project_inquiry"

		if result.LeadScore < 60 {
			result.LeadScore = 70
		}

		result.Summary = "Pengguna terlihat sedang mencari bantuan untuk sebuah project."
	}

	// ---------------------------------
	// Website
	// ---------------------------------

	if strings.Contains(text, "website") ||
		strings.Contains(text, "landing page") ||
		strings.Contains(text, "company profile") {

		result.Category = "website"
	}

	// ---------------------------------
	// POS
	// ---------------------------------

	if strings.Contains(text, "kasir") ||
		strings.Contains(text, "point of sale") {

		result.Category = "pos"
	}

	// ---------------------------------
	// Inventory
	// ---------------------------------

	if strings.Contains(text, "inventory") {
		result.Category = "custom_software"
	}

	// ---------------------------------
	// Urgency
	// ---------------------------------

	if strings.Contains(text, "segera") ||
		strings.Contains(text, "urgent") ||
		strings.Contains(text, "secepatnya") {

		result.HasUrgency = true
		result.LeadScore += 10
	}

	if result.LeadScore > 100 {
		result.LeadScore = 100
	}

	return result, nil
}
