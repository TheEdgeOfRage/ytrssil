package models

import (
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
	// IsShort indicates if a video is a YouTube short
	IsShort bool `json:"short"`
}

type PaginatedVideos struct {
	Videos     []Video
	NextOffset int
}
