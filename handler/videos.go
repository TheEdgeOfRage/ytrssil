package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

var ErrInvalidProgress = errors.New("invalid progress time")

func (h *handler) GetNewVideos(ctx context.Context, sortDesc bool) ([]models.Video, error) {
	return h.db.GetNewVideos(ctx, sortDesc)
}

func (h *handler) GetWatchedVideos(ctx context.Context, sortDesc bool) ([]models.Video, error) {
	return h.db.GetWatchedVideos(ctx, sortDesc)
}

type APIVideoListResponse struct {
	Items []APIVideo `json:"items"`
}

type APIVideo struct {
	ID             string            `json:"id"`
	ContentDetails APIContentDetails `json:"contentDetails"`
}

type APIContentDetails struct {
	Duration string `json:"duration"`
}

func (h *handler) getVideoDurations(ctx context.Context, videos map[string]*models.Video) error {
	ids := strings.Builder{}
	for id := range videos {
		ids.WriteString(id)
		ids.WriteString(",")
	}
	query := url.Values{}
	query.Add("id", ids.String())
	query.Add("part", "contentDetails")
	query.Add("key", h.youTubeAPIKey)

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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch video details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		var bodyStr string
		if err != nil {
			h.log.Error("failed to decode error body", "error", err)
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

func (h *handler) addVideosForChannel(ctx context.Context, parsedChannel *feedparser.Channel) {
	var err error
	videos := make(map[string]*models.Video, len(parsedChannel.Videos))
	for _, parsedVideo := range parsedChannel.Videos {
		date, err := parsedVideo.Published.Parse()
		if err != nil {
			h.log.Error("Failed to parse video information", "call", "feedparser.Parse", "err", err)
			continue
		}

		videoID := strings.Split(parsedVideo.ID, ":")[2]
		exists, err := h.db.HasVideo(ctx, videoID)
		if err != nil {
			h.log.Error("Failed to save video to db", "call", "db.AddVideo", "err", err)
			continue
		}
		if !exists {
			videos[videoID] = &models.Video{
				ID:            videoID,
				Title:         parsedVideo.Title,
				PublishedTime: date,
				IsShort:       parsedVideo.IsShort,
			}
		}
	}

	err = h.getVideoDurations(ctx, videos)
	if err != nil {
		h.log.Error("Failed to get video durations", "call", "handler.getVideoDurations", "err", err)
		return
	}

	for _, video := range videos {
		err = h.db.AddVideo(ctx, *video, parsedChannel.ID)
		if err != nil {
			if !errors.Is(err, db.ErrVideoExists) {
				h.log.Error("Failed to save video to db", "call", "db.AddVideo", "err", err)
			}
			continue
		}
	}
}

func (h *handler) FetchVideos(ctx context.Context) error {
	h.log.Info("Fetching new videos for all channels")

	channels, err := h.db.ListChannels(ctx)
	if err != nil {
		return err
	}
	parsedChannels := make(chan *feedparser.Channel, len(channels))
	errors := make(chan error, len(channels))
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, channel := range channels {
		wg.Add(1)
		go h.parser.ParseThreadSafe(channel.ID, parsedChannels, errors, &mu, &wg)
	}
	wg.Wait()

	for range channels {
		parsedChannel := <-parsedChannels
		err = <-errors
		if err != nil {
			continue
		}
		h.addVideosForChannel(ctx, parsedChannel)
	}

	return nil
}

func (h *handler) MarkVideoAsWatched(ctx context.Context, videoID string) error {
	watchTime := time.Now()
	return h.db.SetVideoWatchTime(ctx, videoID, &watchTime)
}

func (h *handler) MarkVideoAsUnwatched(ctx context.Context, videoID string) error {
	return h.db.SetVideoWatchTime(ctx, videoID, nil)
}

func (h *handler) SetVideoProgress(ctx context.Context, videoID string, progressTime string) (*models.Video, error) {
	progress, err := time.ParseDuration(progressTime)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidProgress, err.Error())
	}

	video, err := h.db.SetVideoProgress(ctx, videoID, int(math.Floor(progress.Seconds())))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidProgress, err.Error())
	}

	return video, nil
}
