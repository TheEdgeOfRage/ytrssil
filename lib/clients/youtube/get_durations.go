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
	ID                   string                   `json:"id"`
	Snippet              APISnippet               `json:"snippet"`
	ContentDetails       APIContentDetails        `json:"contentDetails"`
	Player               APIPlayer                `json:"player"`
	LiveStreamingDetails *APILiveStreamingDetails `json:"liveStreamingDetails"`
}

type APISnippet struct {
	PublishedAt          string `json:"publishedAt"`
	ChannelID            string `json:"channelId"`
	Title                string `json:"title"`
	ChannelName          string `json:"channelTitle"`
	LiveBroadcastContent string `json:"liveBroadcastContent"`
}

type APILiveStreamingDetails struct {
	ActualEndTime string `json:"actualEndTime"`
}

type APIPlayer struct {
	EmbedHTML string `json:"embedHtml"`
}

type APIContentDetails struct {
	Duration   string `json:"duration"`
	Definition string `json:"definition"`
}

func parseISO8601Duration(s string) (time.Duration, error) {
	if s == "P0D" {
		return 0, nil
	}

	var total time.Duration
	s = strings.ToLower(strings.TrimPrefix(s, "P"))

	if i := strings.Index(s, "d"); i != -1 {
		days, err := time.ParseDuration(s[:i] + "h")
		if err != nil {
			return 0, fmt.Errorf("failed to parse days: %w", err)
		}
		total += days * 24
		s = s[i+1:]
	}

	s = strings.TrimPrefix(s, "t")
	if s == "" {
		return total, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}

	return total + d, nil
}

func (c *youTubeClient) GetVideoDurations(ctx context.Context, videos map[string]*models.Video) error {
	if len(videos) == 0 {
		return nil
	}

	ids := strings.Builder{}
	for id := range videos {
		ids.WriteString(id)
		ids.WriteString(",")
	}
	videoIDs := ids.String()
	query := url.Values{}
	query.Add("id", videoIDs)
	query.Add("part", "contentDetails,liveStreamingDetails")
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
		duration, err := parseISO8601Duration(v.ContentDetails.Duration)
		if err != nil {
			c.log.Error(
				"Failed to parse video duration",
				"error", err,
				"durationString", v.ContentDetails.Duration,
				"videoID", v.ID,
			)
			return fmt.Errorf("failed to parse video duration [%v]: %w", v.ContentDetails.Duration, err)
		}
		videos[v.ID].DurationSeconds = int(math.Round(duration.Seconds()))
		videos[v.ID].IsLive = v.LiveStreamingDetails != nil && v.LiveStreamingDetails.ActualEndTime == ""
	}

	return nil
}
