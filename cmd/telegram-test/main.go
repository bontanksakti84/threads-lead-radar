package main

import (
	"context"
	"log"
	"os"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
	"github.com/bontanksakti84/threads-lead-radar/internal/telegram"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not loaded")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	bot := telegram.NewBot(token, chatID)

	post := leads.Post{
		Username:  "demo_umkm",
		URL:       "https://www.threads.net/",
		Category:  "website",
		Intent:    "hire_developer",
		AISummary: "Membutuhkan developer untuk membuat website bisnis.",
	}

	classification := ai.Classification{
		NeedsDeveloper: true,
		HasBudget:      true,
		HasUrgency:     true,
	}

	message := telegram.FormatLeadAlert(
		post,
		classification,
		85,
		[]string{
			"buat website",
			"butuh developer",
		},
	)

	if err := bot.SendMessage(context.Background(), message); err != nil {
		log.Fatal(err)
	}

	log.Println("Lead alert sent successfully")
}
