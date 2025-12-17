package ytrssil_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/handler"
	"github.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
	mockFeedparser "github.com/TheEdgeOfRage/ytrssil-api/mocks/feedparser"
	mockYouTube "github.com/TheEdgeOfRage/ytrssil-api/mocks/youtube"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func init() {
	time.Local = time.UTC
}

type EndpointsTestSuite struct {
	suite.Suite
	cfg           config.Config
	schema        string
	db            db.DB
	dbConn        *pgx.Conn
	parser        *mockFeedparser.ParserMock
	youtubeClient *mockYouTube.ClientMock
	server        *http.Server
}

func (s *EndpointsTestSuite) SetupSuite() {
	var err error
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	s.cfg = config.TestConfig()

	s.schema = fmt.Sprintf("ytrssil_test_%s", ulid.Make().String())
	s.dbConn, err = pgx.Connect(context.Background(), s.cfg.DBURI)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	_, err = s.dbConn.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA %s", s.schema))
	if err != nil {
		panic(fmt.Sprintf("failed to create test schema: %v", err))
	}

	testDBURI := fmt.Sprintf("%s&search_path=%s", s.cfg.DBURI, s.schema)
	cmd := exec.Command("./bin/migrate", "-database", testDBURI, "-path", "migrations", "up")
	cmd.Dir = "../.."
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v\nOutput: %s", err, string(output)))
	}

	s.db, err = db.NewPostgresDB(l, testDBURI)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	s.parser = &mockFeedparser.ParserMock{
		ParseFunc: func(channelID string) (*feedparser.Channel, error) {
			publishTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
			return &feedparser.Channel{
				ID:   channelID,
				Name: fmt.Sprintf("Test Channel %s", channelID),
				Videos: []*feedparser.Video{
					{
						ID:        fmt.Sprintf("yt:video:%s-video1", channelID),
						Title:     "Test Video 1",
						Published: feedparser.Date(publishTime),
						IsShort:   false,
					},
				},
			}, nil
		},
	}

	s.youtubeClient = &mockYouTube.ClientMock{
		GetVideoDurationsFunc: func(ctx context.Context, videos map[string]*models.Video) error {
			for _, video := range videos {
				video.DurationSeconds = 300
			}
			return nil
		},
		GetVideoMetadataFunc: func(ctx context.Context, videoID string) (*models.Video, error) {
			return &models.Video{
				ID:              videoID,
				Title:           "Test Video",
				PublishedTime:   time.Now().Add(-24 * time.Hour),
				DurationSeconds: 300,
				IsShort:         false,
				ChannelID:       "test-channel",
				ChannelName:     "Test Channel",
			}, nil
		},
		GetChannelImageURLFunc: func(ctx context.Context, channelID string) (string, error) {
			return fmt.Sprintf("https://example.com/%s.jpg", channelID), nil
		},
	}

	h := handler.New(l, s.db, s.parser, s.youtubeClient)

	gin.SetMode(gin.TestMode)
	router, err := ytrssil.SetupGinRouter(l, s.cfg, h)
	if err != nil {
		panic(fmt.Sprintf("failed to setup gin router: %v", err))
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.Port),
		Handler: router,
	}
}

func (s *EndpointsTestSuite) TearDownSuite() {
	defer s.dbConn.Close(context.Background())
	_, err := s.dbConn.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA %s CASCADE", s.schema))
	if err != nil {
		panic(fmt.Sprintf("failed to drop test schema: %v", err))
	}
}

func (s *EndpointsTestSuite) SetupTest() {
	_, err := s.dbConn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s.videos", s.schema))
	if err != nil {
		panic(fmt.Sprintf("failed to drop test schema: %v", err))
	}
}
