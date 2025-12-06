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
	ErrUserExists        = errors.New("user already exists")
)

// DB represents a database layer for getting video and channel data
type DB interface {
	// AuthenticateUser verifies a user's password against a hashed value
	AuthenticateUser(ctx context.Context, user models.User) (bool, error)
	// CreateUser registers a new user in the database
	CreateUser(ctx context.Context, user models.User) error
	// DeleteUser registers a new user in the database
	DeleteUser(ctx context.Context, username string) error

	// CreateChannel starts tracking a new channel and fetch new videos for it
	CreateChannel(ctx context.Context, channel models.Channel) error
	// ListChannels lists all channels from the database
	ListChannels(ctx context.Context) ([]models.Channel, error)
	// GetChannelSubscribers lists all channels from the database
	GetChannelSubscribers(ctx context.Context, channelID string) ([]string, error)
	// SubscribeUserToChannel will start adding new videos from that channel to the user
	SubscribeUserToChannel(ctx context.Context, username string, channelID string) error
	// SubscribeUserToChannel will stop adding videos from that channel to the user
	UnsubscribeUserFromChannel(ctx context.Context, username string, channelID string) error

	// GetNewVideos returns a list of unwatched videos from all subscribed channels
	GetNewVideos(ctx context.Context, username string, sortDesc bool) ([]models.Video, error)
	// GetWatchedVideos returns a list of all watched videos for a user
	GetWatchedVideos(ctx context.Context, username string) ([]models.Video, error)
	// AddVideo adds a newly published video to the database
	AddVideo(ctx context.Context, video models.Video, channelID string) error
	// AddVideoToUser will list the video in the users feed
	AddVideoToUser(ctx context.Context, username string, videoID string) error
	// SetVideoWatchTime sets or unsets the watch timestamp of a user's video
	SetVideoWatchTime(ctx context.Context, username string, videoID string, watchTime *time.Time) error
}
