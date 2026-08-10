package scheduler

import (
	"context"
	"sharev2/internal/download"
	"time"
)

type ExpireCleanupScheduler struct {
	manager  *download.Manager
	interval time.Duration
}

func NewExprieCleanupScheduler(dm *download.Manager, interval time.Duration) *ExpireCleanupScheduler {
	return &ExpireCleanupScheduler{
		manager:  dm,
		interval: interval,
	}
}

func (s *ExpireCleanupScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.manager.Cleanup()
		case <-ctx.Done():
			return
		}
	}
}
