package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

type APIVideoListResponse struct {
	Items []APIVideo `json:"items"`
}

type APIVideo struct {
	ID             string            `json:"id"`
	Snippet        APISnippet        `json:"snippet"`
	ContentDetails APIContentDetails `json:"contentDetails"`
	Player         APIPlayer         `json:"player"`
}

type APISnippet struct {
	PublishedAt string `json:"publishedAt"`
	ChannelID   string `json:"channelId"`
	Title       string `json:"title"`
	ChannelName string `json:"channelTitle"`
}

type APIPlayer struct {
	EmbedHTML string `json:"embedHtml"`
}

type APIContentDetails struct {
	Duration   string `json:"duration"`
	Definition string `json:"definition"`
}

func (c *youTubeClient) GetVideoDurations(ctx context.Context, videos map[string]*models.Video) error {
	ids := strings.Builder{}
	for id := range videos {
		ids.WriteString(id)
		ids.WriteString(",")
	}
	videoIDs := ids.String()
	query := url.Values{}
	query.Add("id", videoIDs)
	query.Add("part", "contentDetails")
	query.Add("key", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	if err != nil {
		return fmt.Errorf("failed to set up request: %w", err)
	}
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/youtube/v3/videos",
		RawQuery: query.Encode(),
	}

	c.log.Info("Making request to YouTube API for video durations", "videoIDs", videoIDs)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch video details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		var bodyStr string
		if err != nil {
			c.log.Error("failed to decode error body", "error", err)
			bodyStr = "failed to decode body"
		} else {
			bodyStr = string(body)
		}
		return fmt.Errorf("got non-200 status from YouTube API [%d]: %v", resp.StatusCode, bodyStr)
	}

	decoder := json.NewDecoder(resp.Body)
	var respData APIVideoListResponse
	err = decoder.Decode(&respData)
	if err != nil {
		return fmt.Errorf("failed to decode video details: %w", err)
	}

	for _, v := range respData.Items {
		durationStr := strings.TrimPrefix(v.ContentDetails.Duration, "PT")
		durationStr = strings.ToLower(durationStr)
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("failed to parse video duration [%v]: %w", durationStr, err)
		}
		videos[v.ID].DurationSeconds = int(math.Round(duration.Seconds()))
	}

	return nil
}
