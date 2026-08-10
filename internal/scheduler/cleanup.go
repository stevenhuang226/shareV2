package scheduler

import (
	"context"
	"sharev2/internal/upload"
	"time"
)

type CleanupScheduler struct {
	manager  *upload.Manager
	interval time.Duration
}

func NewCleanupScheduler(manager *upload.Manager, interval time.Duration) *CleanupScheduler {
	return &CleanupScheduler{
		manager:  manager,
		interval: interval,
	}
}

func (s *CleanupScheduler) Run(ctx context.Context) {
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
