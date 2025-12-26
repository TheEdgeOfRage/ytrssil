package handler

import (
	"context"
	"os"
	"time"
)

const (
	CleanupInterval = 1 * time.Hour
	CleanupAge      = 48 * time.Hour
)

func (h *handler) CleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	h.log.Info("Starting cleanup goroutine", "interval", CleanupInterval, "age", CleanupAge)

	for {
		select {
		case <-ctx.Done():
			h.log.Info("Cleanup context done, stopping cleanup goroutine")
			return
		case <-ticker.C:
			h.performCleanup(ctx)
		}
	}
}

func (h *handler) performCleanup(ctx context.Context) {
	videos, err := h.db.GetVideosForCleanup(ctx, CleanupAge)
	if err != nil {
		h.log.Error("Failed to get videos for cleanup", "error", err)
		return
	}

	for _, video := range videos {
		if video.FilePath == nil {
			continue
		}

		h.log.Info("Cleaning up video file", "video_id", video.ID, "path", *video.FilePath)

		if err := os.Remove(*video.FilePath); err != nil {
			h.log.Error("Failed to delete file", "path", *video.FilePath, "error", err)
			continue
		}

		if err := h.db.DeleteVideoFile(ctx, video.ID); err != nil {
			h.log.Error("Failed to update DB after cleanup", "video_id", video.ID, "error", err)
		}
	}

	if len(videos) > 0 {
		h.log.Info("Cleanup completed", "files_deleted", len(videos))
	}
}
