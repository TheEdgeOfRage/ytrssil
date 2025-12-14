package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APIChannelListResponse struct {
	Items []struct {
		Snippet struct {
			Thumbnails struct {
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

func (c *youTubeClient) GetChannelImageURL(ctx context.Context, channelID string) (string, error) {
	query := url.Values{}
	query.Add("id", channelID)
	query.Add("part", "snippet")
	query.Add("fields", "items/snippet/thumbnails")
	query.Add("key", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to set up request: %w", err)
	}
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/youtube/v3/channels",
		RawQuery: query.Encode(),
	}

	c.log.Info("Making request to YouTube API for channel image", "channelID", channelID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch channel details: %w", err)
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
		return "", fmt.Errorf("got non-200 status from YouTube API [%d]: %v", resp.StatusCode, bodyStr)
	}

	decoder := json.NewDecoder(resp.Body)
	var respData APIChannelListResponse
	err = decoder.Decode(&respData)
	if err != nil {
		return "", fmt.Errorf("failed to decode channel details: %w", err)
	}

	if len(respData.Items) == 0 {
		return "", fmt.Errorf("channel not found [%s]", channelID)
	}

	channelData := respData.Items[0]
	return channelData.Snippet.Thumbnails.Medium.URL, nil
}
