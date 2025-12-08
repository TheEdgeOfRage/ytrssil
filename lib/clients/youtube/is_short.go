package youtube

import (
	"context"
	"fmt"
	"net/http"
)

func (c *youTubeClient) isShort(ctx context.Context, videoID string) (bool, error) {
	url := fmt.Sprintf("https://www.youtube.com/shorts/%s", videoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		c.log.Error(
			"Failed to create request to check if a video is a short",
			"videoID", videoID,
			"error", err,
		)
		return false, err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		c.log.Error("Failed to check if a video is a short", "videoID", videoID, "error", err)
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}
