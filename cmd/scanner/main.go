package main

import (
	"context"
	"log"
	"os"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
	"github.com/bontanksakti84/threads-lead-radar/internal/database"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
	"github.com/bontanksakti84/threads-lead-radar/internal/threads"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not loaded")
	}

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	db, err := database.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()

	repository := leads.NewRepository(db.Pool)
	leadService := leads.NewService(repository)

	keywordRepository := leads.NewKeywordRepository(db.Pool)
	scanRunRepository := leads.NewScanRunRepository(db.Pool)

	threadsClient := threads.NewMockClient()

	classifier, err := ai.NewGroqClassifier()
	if err != nil {
		log.Fatal(err)
	}

	scanner := threads.NewScanner(
		threadsClient,
		leadService,
		keywordRepository,
		classifier,
		scanRunRepository,
	)

	if err := scanner.ScanAll(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Scan completed successfully")
}
