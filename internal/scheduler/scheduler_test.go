package scheduler_test

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/scheduler"
)

var silentLogger = log.New(nopWriter{}, "", 0)

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestRun_ExecutesJobImmediately(t *testing.T) {
	var calls int32
	job := func(_ context.Context) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	s := scheduler.New(10*time.Second, job, silentLogger)
	_ = s.Run(ctx)

	if atomic.LoadInt32(&calls) < 1 {
		t.Fatal("expected job to be called at least once immediately")
	}
}

func TestRun_RepeatsOnInterval(t *testing.T) {
	var calls int32
	job := func(_ context.Context) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	s := scheduler.New(25*time.Millisecond, job, silentLogger)
	_ = s.Run(ctx)

	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", atomic.LoadInt32(&calls))
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Run

	s := scheduler.New(time.Second, func(_ context.Context) error { return nil }, silentLogger)
	err := s.Run(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRun_JobErrorDoesNotStop(t *testing.T) {
	var calls int32
	job := func(_ context.Context) error {
		atomic.AddInt32(&calls, 1)
		return errors.New("transient error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	s := scheduler.New(25*time.Millisecond, job, silentLogger)
	_ = s.Run(ctx)

	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("expected scheduler to continue after error, got %d calls", atomic.LoadInt32(&calls))
	}
}
