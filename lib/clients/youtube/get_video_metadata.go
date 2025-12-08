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

func (c *youTubeClient) GetVideoMetadata(ctx context.Context, videoID string) (*models.Video, error) {
	query := url.Values{}
	query.Add("id", videoID)
	query.Add("part", "snippet,contentDetails,player")
	query.Add("key", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to set up request: %w", err)
	}
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/youtube/v3/videos",
		RawQuery: query.Encode(),
	}

	c.log.Info("Making request to YouTube API for video metadata", "videoID", videoID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video details: %w", err)
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
		return nil, fmt.Errorf("got non-200 status from YouTube API [%d]: %v", resp.StatusCode, bodyStr)
	}

	decoder := json.NewDecoder(resp.Body)
	var respData APIVideoListResponse
	err = decoder.Decode(&respData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode video details: %w", err)
	}

	if len(respData.Items) == 0 {
		return nil, fmt.Errorf("video not found [%s]", videoID)
	}

	videoData := respData.Items[0]
	video := &models.Video{
		ID:          videoData.ID,
		Title:       videoData.Snippet.Title,
		ChannelName: videoData.Snippet.ChannelName,
		ChannelID:   videoData.Snippet.ChannelID,
	}

	publishedTime, err := time.Parse(time.RFC3339, videoData.Snippet.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published time [%v]: %w", videoData.Snippet.PublishedAt, err)
	}
	video.PublishedTime = publishedTime

	durationStr := strings.TrimPrefix(videoData.ContentDetails.Duration, "PT")
	durationStr = strings.ToLower(durationStr)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse video duration [%v]: %w", durationStr, err)
	}
	video.DurationSeconds = int(math.Round(duration.Seconds()))

	video.IsShort, err = c.isShort(ctx, videoID)
	if err != nil {
		return nil, err
	}

	return video, nil
}
