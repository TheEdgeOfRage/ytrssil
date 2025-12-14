package db

import (
	"context"
	"errors"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

var (
	ErrChannelExists     = errors.New("channel already exists")
	ErrChannelNotFound   = errors.New("no channel with that ID found")
	ErrAlreadySubscribed = errors.New("already subscribed to channel")
	ErrVideoExists       = errors.New("video already exists")
)

// DB represents a database layer for getting video and channel data
type DB interface {
	// ListChannels lists all channels from the database
	ListChannels(ctx context.Context) ([]models.Channel, error)
	// SubscribeToChannel will start fetching new videos from that channel
	SubscribeToChannel(ctx context.Context, channel models.Channel) error
	// UnsubscribeToChannel will stop fetching videos from that channel
	UnsubscribeFromChannel(ctx context.Context, channelID string) error

	// GetNewVideos returns a list of unwatched videos from all subscribed channels
	GetNewVideos(ctx context.Context, sortDesc bool) ([]models.Video, error)
	// GetWatchedVideos returns a list of all watched videos
	GetWatchedVideos(ctx context.Context, sortDesc bool, limit int, offset int) ([]models.Video, error)
	// HasVideo returns true if the video with the given ID exists in the DB
	HasVideo(ctx context.Context, videoID string) (bool, error)
	// AddVideo adds a newly published video to the database
	AddVideo(ctx context.Context, video models.Video, channelID string) error
	// SetVideoWatchTime sets or unsets the watch timestamp of a video
	SetVideoWatchTime(ctx context.Context, videoID string, watchTime *time.Time) error
	// SetVideoProgress sets or unsets the watch progress of a video
	SetVideoProgress(ctx context.Context, videoID string, progress int) (*models.Video, error)
}
