package handler

import (
	"context"
	"log/slog"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/lib/clients/youtube"
	"github.com/TheEdgeOfRage/ytrssil-api/lib/downloader"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

type Handler interface {
	SubscribeToChannel(ctx context.Context, channelID string) (*models.Channel, error)
	UnsubscribeFromChannel(ctx context.Context, channelID string) error
	ListChannels(ctx context.Context) ([]models.Channel, error)
	GetNewVideos(ctx context.Context, sortDesc bool) ([]models.Video, error)
	GetWatchedVideos(ctx context.Context, sortDesc bool, page int) ([]models.Video, error)
	FetchVideos(ctx context.Context) error
	MarkVideoAsWatched(ctx context.Context, videoID string) error
	MarkVideoAsUnwatched(ctx context.Context, videoID string) error
	SetVideoProgress(ctx context.Context, videoID string, progressTime string) (*models.Video, error)
	AddCustomVideo(ctx context.Context, videoID string) error
	DownloadVideo(ctx context.Context, videoID string) error
	ServeVideoFile(ctx context.Context, videoID string) (filePath string, filename string, err error)
	CleanupRoutine(ctx context.Context)
}

type handler struct {
	log           *slog.Logger
	db            db.DB
	parser        feedparser.Parser
	youTubeClient youtube.Client
	downloader    downloader.Downloader
	downloadsDir  string
}

func New(
	log *slog.Logger,
	db db.DB,
	parser feedparser.Parser,
	youTubeClient youtube.Client,
	downloader downloader.Downloader,
	downloadsDir string,
) *handler {
	return &handler{
		log:           log,
		db:            db,
		parser:        parser,
		youTubeClient: youTubeClient,
		downloader:    downloader,
		downloadsDir:  downloadsDir,
	}
}
