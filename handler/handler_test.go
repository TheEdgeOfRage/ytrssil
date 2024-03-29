package handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
	db_mock "gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/mocks/db"
	parser_mock "gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/mocks/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var testConfig config.Config

func init() {
	testConfig = config.TestConfig()
}

func TestGetNewVideos(t *testing.T) {
	// Arrange
	l := log.NewNopLogger()

	handler := New(l, &db_mock.DBMock{
		GetNewVideosFunc: func(ctx context.Context, username string) ([]models.Video, error) {
			return []models.Video{
				{
					ID:            "test",
					ChannelName:   "test",
					Title:         "test",
					PublishedTime: time.Now(),
					WatchTime:     nil,
				},
			}, nil
		},
	}, &parser_mock.ParserMock{})

	// Act
	resp, err := handler.GetNewVideos(context.TODO(), "username")

	// Assert
	if assert.NoError(t, err) {
		if assert.NotNil(t, resp) {
			assert.Equal(t, resp[0].ID, "test")
			assert.Equal(t, resp[0].Title, "test")
		}
	}
}
