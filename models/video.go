package models

import (
	"math"
	"time"
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

func (v Video) ProgressPercentage() int {
	return int(math.Floor(100 * float64(v.ProgressSeconds) / float64(v.DurationSeconds)))
}

type PaginatedVideos struct {
	Videos     []Video
	NextOffset int
}
