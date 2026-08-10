package scheduler

import (
	"context"
	"sharev2/internal/upload"
	"time"
)

type UploadCleanupScheduler struct {
	manager  *upload.Manager
	interval time.Duration
}

func NewUploadCleanupScheduler(manager *upload.Manager, interval time.Duration) *UploadCleanupScheduler {
	return &UploadCleanupScheduler{
		manager:  manager,
		interval: interval,
	}
}

func (s *UploadCleanupScheduler) Run(ctx context.Context) {
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
