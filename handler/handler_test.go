package handler

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
	db_mock "github.com/TheEdgeOfRage/ytrssil-api/mocks/db"
	parser_mock "github.com/TheEdgeOfRage/ytrssil-api/mocks/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

var testConfig config.Config

func init() {
	testConfig = config.TestConfig()
}

func TestGetNewVideos(t *testing.T) {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := New(l, testConfig, &db_mock.DBMock{
		GetNewVideosFunc: func(ctx context.Context, _ bool) ([]models.Video, error) {
			return []models.Video{
				{
					ID:            "test",
					ChannelName:   "test",
					Title:         "test",
					PublishedTime: time.Now(),
				},
			}, nil
		},
	}, &parser_mock.ParserMock{})
	resp, err := handler.GetNewVideos(context.TODO(), false)

	if assert.NoError(t, err) {
		if assert.NotNil(t, resp) {
			assert.Equal(t, resp[0].ID, "test")
			assert.Equal(t, resp[0].Title, "test")
		}
	}
}
