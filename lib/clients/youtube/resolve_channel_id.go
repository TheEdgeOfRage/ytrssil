package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type apiChannelIDResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

func (c *youTubeClient) ResolveChannelID(ctx context.Context, handle string) (string, error) {
	query := url.Values{}
	query.Add("forHandle", handle)
	query.Add("part", "id")
	query.Add("fields", "items/id")
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

	c.log.Info("Resolving channel handle to ID", "handle", handle)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch channel: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		var bodyStr string
		if err != nil {
			bodyStr = "failed to decode body"
		} else {
			bodyStr = string(body)
		}
		return "", fmt.Errorf("got non-200 status from YouTube API [%d]: %v", resp.StatusCode, bodyStr)
	}

	var respData apiChannelIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(respData.Items) == 0 {
		return "", fmt.Errorf("channel not found for handle %q", handle)
	}

	return respData.Items[0].ID, nil
}
