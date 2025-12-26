package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func sanitizeFilename(title string) string {
	title = strings.ReplaceAll(title, " ", "_")

	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-\.]`)
	title = reg.ReplaceAllString(title, "")

	if len(title) > 200 {
		title = title[:200]
	}

	return title
}

func (h *handler) DownloadVideo(ctx context.Context, videoID string) error {
	exists, err := h.db.HasVideo(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to check video existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("video not found")
	}

	if err := h.db.SetVideoDownloadStatus(ctx, videoID, "pending"); err != nil {
		return fmt.Errorf("failed to set download status: %w", err)
	}

	go h.performDownload(videoID)

	return nil
}

func (h *handler) performDownload(videoID string) {
	ctx := context.Background()

	video, err := h.db.GetVideo(ctx, videoID)
	if err != nil {
		h.log.Error("Failed to get video for download", "video_id", videoID, "error", err)
		if dbErr := h.db.SetVideoDownloadFailed(ctx, videoID, "Failed to get video info"); dbErr != nil {
			h.log.Error("Failed to update download status to failed", "video_id", videoID, "error", dbErr)
		}
		return
	}

	if err := h.db.SetVideoDownloadStatus(ctx, videoID, "downloading"); err != nil {
		h.log.Error("Failed to update download status to downloading", "video_id", videoID, "error", err)
		return
	}

	h.log.Info("Starting video download", "video_id", videoID, "title", video.Title)

	filePath, err := h.downloader.Download(ctx, videoID, video.Title, h.downloadsDir)
	if err != nil {
		h.log.Error("Video download failed", "video_id", videoID, "error", err)
		if dbErr := h.db.SetVideoDownloadFailed(ctx, videoID, err.Error()); dbErr != nil {
			h.log.Error("Failed to update download status to failed", "video_id", videoID, "error", dbErr)
		}
		return
	}

	if err := h.db.SetVideoDownloadCompleted(ctx, videoID, filePath); err != nil {
		h.log.Error("Failed to mark video as downloaded", "video_id", videoID, "error", err)
		os.Remove(filePath)
		if dbErr := h.db.SetVideoDownloadFailed(ctx, videoID, "Failed to update database"); dbErr != nil {
			h.log.Error("Failed to update download status to failed", "video_id", videoID, "error", dbErr)
		}
		return
	}

	h.log.Info("Video downloaded successfully", "video_id", videoID, "path", filePath)
}

func (h *handler) ServeVideoFile(ctx context.Context, videoID string) (filePath string, filename string, err error) {
	video, err := h.db.GetVideo(ctx, videoID)
	if err != nil {
		return "", "", fmt.Errorf("video not found: %w", err)
	}

	if video.FilePath == nil {
		return "", "", fmt.Errorf("video not downloaded")
	}

	if _, err := os.Stat(*video.FilePath); os.IsNotExist(err) {
		h.db.DeleteVideoFile(ctx, videoID)
		return "", "", fmt.Errorf("file not found on disk")
	}

	sanitizedTitle := sanitizeFilename(video.Title)
	if sanitizedTitle == "" {
		sanitizedTitle = videoID
	}
	ext := filepath.Ext(*video.FilePath)
	filename = sanitizedTitle + ext

	return *video.FilePath, filename, nil
}
