package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
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
	return h.DownloadVideoWithFormat(ctx, videoID, "")
}

func (h *handler) DownloadVideoWithFormat(ctx context.Context, videoID string, format string) error {
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

	go h.performDownloadWithFormat(videoID, format)

	return nil
}

func (h *handler) performDownloadWithFormat(videoID string, format string) {
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

	h.log.Info("Starting video download", "video_id", videoID, "title", video.Title, "format", format)

	var filePath string
	if format != "" {
		filePath, err = h.downloader.DownloadWithFormat(ctx, videoID, video.Title, h.config.DownloadsDir, format)
	} else {
		filePath, err = h.downloader.Download(ctx, videoID, video.Title, h.config.DownloadsDir)
	}

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

func (h *handler) GetVideoFormats(ctx context.Context, videoID string) ([]models.VideoFormat, error) {
	cmd := exec.CommandContext(ctx,
		"yt-dlp",
		"--dump-json",
		"--no-playlist",
		"--no-warnings",
		"--quiet",
		"--format", "best",
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get video formats: %w", err)
	}

	var formats []models.VideoFormat

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var formatInfo map[string]interface{}
		if err := json.Unmarshal([]byte(line), &formatInfo); err != nil {
			continue
		}

		if formatArr, ok := formatInfo["formats"].([]interface{}); ok {
			for _, f := range formatArr {
				if fMap, ok := f.(map[string]interface{}); ok {
					if height, ok := fMap["height"].(float64); ok {
						if height > 0 {
							formatID := ""
							if id, ok := fMap["format_id"].(string); ok {
								formatID = id
							}

							note := ""
							if n, ok := fMap["note"].(string); ok {
								note = n
							}

							ext := ""
							if e, ok := fMap["ext"].(string); ok {
								ext = e
							}

							formats = append(formats, models.VideoFormat{
								Height:    int(height),
								FormatID:  formatID,
								Note:      note,
								Extension: ext,
							})
						}
					}
				}
			}
		}
	}

	return formats, nil
}
