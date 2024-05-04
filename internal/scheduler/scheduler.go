package scheduler

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	interval time.Duration
	task     func(ctx context.Context)
}

func NewScheduler(interval time.Duration, task func(ctx context.Context)) *Scheduler {
	return &Scheduler{
		interval,
		task,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	var wg sync.WaitGroup

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.task(ctx)
			}()

		case <-ctx.Done():
			break loop
		}
	}

	wg.Wait()
}
