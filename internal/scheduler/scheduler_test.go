package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"
)

type mockScanner struct {
	mu    sync.Mutex
	calls int
}

func (m *mockScanner) ScanAll(ctx context.Context) error {
	m.mu.Lock()
	m.calls++
	m.mu.Unlock()

	return nil
}

func (m *mockScanner) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.calls
}

func TestSchedulerRunsImmediately(t *testing.T) {
	mock := &mockScanner{}

	s := New(
		mock,
		time.Hour,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Run(ctx)

	time.Sleep(50 * time.Millisecond)
	cancel()

	if mock.Calls() != 1 {
		t.Fatalf(
			"expected scanner to run once, got %d",
			mock.Calls(),
		)
	}
}

type blockingScanner struct {
	started chan struct{}
	release chan struct{}

	mu    sync.Mutex
	calls int
}

func newBlockingScanner() *blockingScanner {
	return &blockingScanner{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
}

func (m *blockingScanner) ScanAll(ctx context.Context) error {
	m.mu.Lock()
	m.calls++
	m.mu.Unlock()

	select {
	case <-m.started:
	default:
		close(m.started)
	}

	<-m.release

	return nil
}

func (m *blockingScanner) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.calls
}

func TestSchedulerPreventsOverlappingScans(t *testing.T) {
	mock := newBlockingScanner()

	s := New(
		mock,
		time.Hour,
	)

	ctx := context.Background()

	firstScanDone := make(chan struct{})

	go func() {
		s.runScan(ctx)
		close(firstScanDone)
	}()

	select {
	case <-mock.started:
	case <-time.After(time.Second):
		t.Fatal("first scan did not start")
	}

	// Try to start another scan while the first one is still running.
	secondScanDone := make(chan struct{})

	go func() {
		s.runScan(ctx)
		close(secondScanDone)
	}()

	select {
	case <-secondScanDone:
	case <-time.After(time.Second):
		t.Fatal("second scan did not return")
	}

	if calls := mock.Calls(); calls != 1 {
		t.Fatalf(
			"expected only one scan while first scan is running, got %d",
			calls,
		)
	}

	// Allow the first scan to finish.
	close(mock.release)

	select {
	case <-firstScanDone:
	case <-time.After(time.Second):
		t.Fatal("first scan did not finish")
	}

	// A new scan should now be allowed.
	mock2 := newBlockingScanner()

	s2 := New(
		mock2,
		time.Hour,
	)

	done := make(chan struct{})

	go func() {
		s2.runScan(ctx)
		close(done)
	}()

	select {
	case <-mock2.started:
	case <-time.After(time.Second):
		t.Fatal("new scan did not start after previous scan finished")
	}

	if calls := mock2.Calls(); calls != 1 {
		t.Fatalf(
			"expected new scan to start, got %d",
			calls,
		)
	}

	close(mock2.release)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("new scan did not finish")
	}
}
