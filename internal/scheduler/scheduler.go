package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type Scanner interface {
	ScanAll(ctx context.Context) error
}

type Scheduler struct {
	scanner  Scanner
	interval time.Duration

	mu      sync.Mutex
	running bool
}

func New(scanner Scanner, interval time.Duration) *Scheduler {
	return &Scheduler{
		scanner:  scanner,
		interval: interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	log.Printf(
		"Scheduler started | interval=%s",
		s.interval,
	)

	// Run immediately when the process starts.
	s.runScan(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return

		case <-ticker.C:
			s.runScan(ctx)
		}
	}
}

func (s *Scheduler) runScan(ctx context.Context) {
	if !s.startScan() {
		log.Println("Scheduled scan skipped | previous scan still running")
		return
	}

	defer s.finishScan()

	startedAt := time.Now()

	log.Println("Scheduled scan started")

	if err := s.scanner.ScanAll(ctx); err != nil {
		log.Printf(
			"Scheduled scan failed: %v",
			err,
		)
		return
	}

	log.Printf(
		"Scheduled scan completed | duration=%s",
		time.Since(startedAt),
	)
}

func (s *Scheduler) startScan() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return false
	}

	s.running = true
	return true
}

func (s *Scheduler) finishScan() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.running = false
}
