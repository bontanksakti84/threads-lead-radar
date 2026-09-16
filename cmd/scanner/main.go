package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
	"github.com/bontanksakti84/threads-lead-radar/internal/database"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
	"github.com/bontanksakti84/threads-lead-radar/internal/scheduler"
	"github.com/bontanksakti84/threads-lead-radar/internal/telegram"
	"github.com/bontanksakti84/threads-lead-radar/internal/threads"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not loaded")
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

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

	telegramBot := telegram.NewBot(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		os.Getenv("TELEGRAM_CHAT_ID"),
	)

	scanner := threads.NewScanner(
		threadsClient,
		leadService,
		keywordRepository,
		classifier,
		scanRunRepository,
		telegramBot,
		60,
	)

	intervalMinutes := 10

	if value := os.Getenv("SCAN_INTERVAL_MINUTES"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			log.Printf(
				"Invalid SCAN_INTERVAL_MINUTES=%q, using default 10 minutes",
				value,
			)
		} else {
			intervalMinutes = parsed
		}
	}

	scanScheduler := scheduler.New(
		scanner,
		time.Duration(intervalMinutes)*time.Minute,
	)

	scanScheduler.Run(ctx)
}
