package scheduler

import (
	"context"
	"log"
	"time"
)

// Scheduler runs a job function on a fixed interval until the context is cancelled.
type Scheduler struct {
	interval time.Duration
	job      func(ctx context.Context) error
	logger   *log.Logger
}

// New creates a Scheduler with the given interval and job function.
func New(interval time.Duration, job func(ctx context.Context) error, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}
}

// Run executes the job immediately, then repeats on each tick until ctx is done.
func (s *Scheduler) Run(ctx context.Context) error {
	if err := s.tick(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scheduler: context cancelled, stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := s.tick(ctx); err != nil {
				s.logger.Printf("scheduler: job error: %v", err)
			}
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) error {
	s.logger.Printf("scheduler: running job at %s", time.Now().Format(time.RFC3339))
	return s.job(ctx)
}
