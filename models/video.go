package models

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

type Video struct {
	// YouTube ID of the video
	ID string `json:"video_id"`
	// Name of the channel the video belongs to
	ChannelName string `json:"channel_name"`
	// ID of the channel the video belongs to
	ChannelID string `json:"-"`
	// Title of the video
	Title string `json:"title"`
	// Video publish timestamp
	PublishedTime time.Time `json:"published_timestamp"`
	// Video watch timestamp
	WatchTime *time.Time `json:"watch_timestamp"`
	// ProgressSeconds is the total duration of the video
	DurationSeconds int `json:"duration"`
	// ProgressSeconds is the saved progress of the video
	ProgressSeconds int `json:"progress"`
	// IsShort indicates if a video is a YouTube short
	IsShort bool `json:"short"`
	// DownloadedAt is the timestamp when the video was downloaded to the server
	DownloadedAt *time.Time `json:"downloaded_at"`
	// FilePath is the path to the downloaded video file on the server
	FilePath *string `json:"-"`
	// DownloadStatus tracks the download state: "", "pending", "downloading", "completed", "failed"
	DownloadStatus *string `json:"download_status"`
	// DownloadError stores the error message if download failed
	DownloadError *string `json:"download_error"`
}

// ProgressPercentage returns the current progress of the video as an integer from 0-100
func (v Video) ProgressPercentage() int {
	return int(100 * float64(v.ProgressSeconds) / float64(v.DurationSeconds))
}

// WatchURL returns the formatted YouTube watch URL including the timestamp pointing to  the current progress
func (v Video) WatchURL() string {
	return fmt.Sprintf("https://youtube.com/watch?v=%s&t=%d", v.ID, v.ProgressSeconds)
}

// Duration returns the total duration of the video in the hh:mm:ss format
func (v Video) Duration() string {
	duration := time.Duration(v.DurationSeconds) * time.Second
	duration = duration.Round(time.Second)
	h := duration / time.Hour
	duration -= h * time.Hour
	m := duration / time.Minute
	duration -= m * time.Minute
	s := duration / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func (v Video) HumanizedPublishTime() string {
	return humanize.Time(v.PublishedTime)
}

func (v Video) IsDownloaded() bool {
	return v.DownloadStatus != nil && *v.DownloadStatus == "completed"
}

func (v Video) IsDownloading() bool {
	return v.DownloadStatus != nil && (*v.DownloadStatus == "pending" || *v.DownloadStatus == "downloading")
}

func (v Video) DownloadFailed() bool {
	return v.DownloadStatus != nil && *v.DownloadStatus == "failed"
}

type PaginatedVideos struct {
	Videos     []Video
	NextOffset int
}
