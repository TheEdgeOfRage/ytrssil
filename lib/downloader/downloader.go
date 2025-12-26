package downloader

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
)

const ytDlpFormat = "bestvideo[height<=1080][vcodec^=vp9]+bestaudio/bestvideo[height<=1080][vcodec^=av01]+bestaudio/bestvideo[height<=1080]+bestaudio/best[height<=1080]/best" // nolint:lll

type Downloader interface {
	Download(ctx context.Context, videoID string, title string, outputDir string) (string, error)
	ValidateInstallation() error
}

type ytdlpDownloader struct {
	log *slog.Logger
}

func NewYtdlpDownloader(log *slog.Logger) *ytdlpDownloader {
	return &ytdlpDownloader{log: log}
}

func (d *ytdlpDownloader) ValidateInstallation() error {
	cmd := exec.Command("yt-dlp", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp not found or not executable: %w", err)
	}
	return nil
}

func (d *ytdlpDownloader) Download(
	ctx context.Context,
	videoID string,
	title string,
	outputDir string,
) (string, error) {
	outputTemplate := filepath.Join(outputDir, fmt.Sprintf("%s.%%(ext)s", videoID))

	cmd := exec.CommandContext(ctx,
		"yt-dlp",
		"--format", ytDlpFormat,
		"--output", outputTemplate,
		"--no-playlist",
		"--no-warnings",
		"--quiet",
		"--progress",
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		d.log.Error("yt-dlp failed", "output", string(output), "error", err)
		return "", fmt.Errorf("download failed: %w", err)
	}

	matches, err := filepath.Glob(filepath.Join(outputDir, fmt.Sprintf("%s.*", videoID)))
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("downloaded file not found")
	}

	return matches[0], nil
}
