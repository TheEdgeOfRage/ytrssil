package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type VideoFormat struct {
	Height    int    `json:"height"`
	FormatID  string `json:"format_id"`
	Note      string `json:"note"`
	Extension string `json:"extension"`
}

const ytDlpFormat = "bestvideo[height<=1080][vcodec^=vp9]+bestaudio/bestvideo[height<=1080][vcodec^=av01]+bestaudio/bestvideo[height<=1080]+bestaudio/best[height<=1080]/best" // nolint:lll

type Downloader interface {
	Download(ctx context.Context, videoID string, title string, outputDir string) (string, error)
	DownloadWithFormat(ctx context.Context, videoID string, title string, outputDir string, format string) (string, error)
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
	return d.DownloadWithFormat(ctx, videoID, title, outputDir, ytDlpFormat)
}

func (d *ytdlpDownloader) DownloadWithFormat(
	ctx context.Context,
	videoID string,
	title string,
	outputDir string,
	format string,
) (string, error) {
	outputTemplate := filepath.Join(outputDir, fmt.Sprintf("%s.%%(ext)s", videoID))

	cmd := exec.CommandContext(ctx,
		"yt-dlp",
		"--format", format,
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

func ParseFormats(output []byte) []VideoFormat {
	var formats []VideoFormat

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
					if height, ok := fMap["height"].(float64); ok && height > 0 {
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

						formats = append(formats, VideoFormat{
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

	// Sort by height descending
	sort.Slice(formats, func(i, j int) bool {
		return formats[i].Height > formats[j].Height
	})

	return formats
}
