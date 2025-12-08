package models

import (
	"fmt"
	"math"
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
}

// ProgressPercentage returns the current progress of the video as an integer from 0-100
func (v Video) ProgressPercentage() int {
	return int(math.Floor(100 * float64(v.ProgressSeconds) / float64(v.DurationSeconds)))
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

type PaginatedVideos struct {
	Videos     []Video
	NextOffset int
}
