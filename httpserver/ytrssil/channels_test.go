package ytrssil_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

type ChannelsTestSuite struct {
	EndpointsTestSuite
}

func TestChannelsTestSuite(t *testing.T) {
	suite.Run(t, new(ChannelsTestSuite))
}

func (s *ChannelsTestSuite) TestSubscribeToChannelJSON() {
	channelID := "test-channel-123"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/channels/%s/subscribe", channelID), nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response models.Channel
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Equal(channelID, response.ID)
	s.Equal(fmt.Sprintf("Test Channel %s", channelID), response.Name)
	s.True(response.Subscribed)
}

func (s *ChannelsTestSuite) TestUnsubscribeFromChannelJSON() {
	channelID := "test-channel-456"

	ctx := context.Background()
	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/channels/%s/unsubscribe", channelID), nil)
	req.Header.Set("Authorization", s.cfg.AuthToken)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Equal("unsubscribed from channel successfully", response["msg"])
}

func (s *ChannelsTestSuite) TestChannelsPage() {
	ctx := context.Background()

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         "channel-606",
		Name:       "Test Channel Page",
		Subscribed: true,
		ImageURL:   "https://example.com/image.jpg",
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/channels", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Test Channel Page")
}

func (s *ChannelsTestSuite) TestSubscribeToChannelPage() {
	form := url.Values{}
	form.Add("channel_id", "channel-707")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "channel-707")
}

func (s *ChannelsTestSuite) TestUnsubscribeFromChannelPage() {
	ctx := context.Background()
	channelID := "channel-808"

	err := s.db.SubscribeToChannel(ctx, models.Channel{
		ID:         channelID,
		Name:       "Test Channel",
		Subscribed: true,
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/channels/%s/unsubscribe", channelID), nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *ChannelsTestSuite) TestSubscribeRequiresAuth() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/channels/test/subscribe", nil)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *ChannelsTestSuite) TestUnsubscribeRequiresAuth() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/channels/test/unsubscribe", nil)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}
