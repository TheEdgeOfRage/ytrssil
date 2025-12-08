package youtube

import (
	"context"
	"log/slog"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

type Client interface {
	GetVideoDurations(ctx context.Context, videos map[string]*models.Video) error
	GetVideoMetadata(ctx context.Context, videoID string) (*models.Video, error)
}

type youTubeClient struct {
	log    *slog.Logger
	apiKey string
}

var _ Client = (*youTubeClient)(nil)

func NewYouTubeClient(log *slog.Logger, apiKey string) *youTubeClient {
	return &youTubeClient{
		log:    log,
		apiKey: apiKey,
	}
}
