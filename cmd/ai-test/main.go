package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not loaded")
	}

	classifier, err := ai.NewGroqClassifier()
	if err != nil {
		log.Fatal(err)
	}

	content := "Saya sedang mencari developer untuk membuat website bisnis saya. Kalau ada yang bisa bantu, boleh DM."

	result, err := classifier.Classify(
		context.Background(),
		content,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== GROQ CLASSIFICATION ===")
	fmt.Printf("Is Lead:         %v\n", result.IsLead)
	fmt.Printf("Category:        %s\n", result.Category)
	fmt.Printf("Intent:          %s\n", result.Intent)
	fmt.Printf("Lead Score:      %d\n", result.LeadScore)
	fmt.Printf("Has Budget:      %v\n", result.HasBudget)
	fmt.Printf("Has Urgency:     %v\n", result.HasUrgency)
	fmt.Printf("Needs Developer: %v\n", result.NeedsDeveloper)
	fmt.Printf("Summary:         %s\n", result.Summary)
}
