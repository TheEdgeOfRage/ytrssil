package handler

import (
	"context"
	"log/slog"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

type Handler interface {
	CreateUser(ctx context.Context, user models.User) error
	SubscribeToChannel(ctx context.Context, username string, channelID string) error
	UnsubscribeFromChannel(ctx context.Context, username string, channelID string) error
	GetNewVideos(ctx context.Context, username string, sortDesc bool) ([]models.Video, error)
	GetWatchedVideos(ctx context.Context, username string) ([]models.Video, error)
	FetchVideos(ctx context.Context) error
	MarkVideoAsWatched(ctx context.Context, username string, videoID string) error
	MarkVideoAsUnwatched(ctx context.Context, username string, videoID string) error
}

type handler struct {
	log    *slog.Logger
	db     db.DB
	parser feedparser.Parser
}

func New(log *slog.Logger, db db.DB, parser feedparser.Parser) *handler {
	return &handler{
		log:    log,
		db:     db,
		parser: parser,
	}
}
