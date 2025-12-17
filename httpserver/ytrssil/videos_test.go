package ytrssil_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

type VideosTestSuite struct {
	EndpointsTestSuite
}

func TestVideosTestSuite(t *testing.T) {
	suite.Run(t, new(VideosTestSuite))
}

func (s *VideosTestSuite) TestGetNewVideosJSON() {
	ctx := context.Background()
	channelID := "test-channel-789"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video123",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/videos/new", nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string][]models.Video
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Require().Contains(response, "videos")
	s.Require().Len(response["videos"], 1)
	s.Equal("video123", response["videos"][0].ID)
}

func (s *VideosTestSuite) TestGetWatchedVideosJSON() {
	ctx := context.Background()
	channelID := "test-channel-101"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video456",
		Title:           "Watched Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	watchTime := time.Now()
	err = s.db.SetVideoWatchTime(ctx, "video456", &watchTime)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/videos/watched", nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string][]models.Video
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Require().Contains(response, "videos")
	s.Require().Len(response["videos"], 1)
	s.Equal("video456", response["videos"][0].ID)
}

func (s *VideosTestSuite) TestMarkVideoAsWatchedJSON() {
	ctx := context.Background()
	channelID := "test-channel-202"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video789",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/videos/video789/watch", nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Require().Contains(response, "msg")
	s.Equal("marked video as watched", response["msg"])
}

func (s *VideosTestSuite) TestMarkVideoAsUnwatchedJSON() {
	ctx := context.Background()
	channelID := "test-channel-303"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video101",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	watchTime := time.Now()
	err = s.db.SetVideoWatchTime(ctx, "video101", &watchTime)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/videos/video101/unwatch", nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Require().Contains(response, "msg")
	s.Equal("cleared video from watch history", response["msg"])
}

func (s *VideosTestSuite) TestNewVideosPage() {
	ctx := context.Background()
	channelID := "test-channel-404"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video202",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Test Video")
}

func (s *VideosTestSuite) TestWatchedVideosPage() {
	ctx := context.Background()
	channelID := "test-channel-505"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video303",
		Title:           "Watched Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	watchTime := time.Now()
	err = s.db.SetVideoWatchTime(ctx, "video303", &watchTime)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/watched", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Watched Video")
}

func (s *VideosTestSuite) TestAddVideoPage() {
	form := url.Values{}
	form.Add("video_id", "custom-video-123")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/videos", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusAccepted, w.Code)
}

func (s *VideosTestSuite) TestMarkVideoAsWatchedPage() {
	ctx := context.Background()
	channelID := "test-channel-909"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video404",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/videos/video404/watch", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *VideosTestSuite) TestMarkVideoAsUnwatchedPage() {
	ctx := context.Background()
	channelID := "test-channel-1010"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video505",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	watchTime := time.Now()
	err = s.db.SetVideoWatchTime(ctx, "video505", &watchTime)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/videos/video505/unwatch", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *VideosTestSuite) TestSetVideoProgressPage() {
	ctx := context.Background()
	channelID := "test-channel-1111"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	err = s.db.AddVideo(ctx, models.Video{
		ID:              "video606",
		Title:           "Test Video",
		PublishedTime:   time.Now().Add(-1 * time.Hour),
		DurationSeconds: 300,
		IsShort:         false,
	}, channelID)
	s.Require().NoError(err)

	form := url.Values{}
	form.Add("progress", "2:30")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/videos/video606/progress", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *VideosTestSuite) TestFetchVideosJSON() {
	ctx := context.Background()

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         "channel-1212",
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/fetch", nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Require().Contains(response, "msg")
	s.Equal("videos fetched successfully", response["msg"])
}

func (s *VideosTestSuite) TestFetchEndpointNoAuthRequired() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/fetch", nil)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *VideosTestSuite) TestVideosRequireAuth() {
	testCases := []struct {
		name   string
		method string
		path   string
		is401  bool
	}{
		{"API new videos", "GET", "/api/videos/new", true},
		{"API watched videos", "GET", "/api/videos/watched", true},
		{"API mark watched", "POST", "/api/videos/test/watch", true},
		{"API mark unwatched", "POST", "/api/videos/test/unwatch", true},
		{"Page home", "GET", "/", false},
		{"Page watched", "GET", "/watched", false},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			s.server.Handler.ServeHTTP(w, req)

			if tc.is401 {
				s.Equal(http.StatusUnauthorized, w.Code)
			} else {
				s.Equal(http.StatusFound, w.Code)
			}
		})
	}
}
