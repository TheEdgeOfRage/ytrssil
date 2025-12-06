package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) GetNewVideos(ctx context.Context, username string, sortDesc bool) ([]models.Video, error) {
	return h.db.GetNewVideos(ctx, username, sortDesc)
}

func (h *handler) GetWatchedVideos(ctx context.Context, username string) ([]models.Video, error) {
	return h.db.GetWatchedVideos(ctx, username)
}

func (h *handler) addVideoToAllSubscribers(ctx context.Context, channelID string, videoID string) error {
	subs, err := h.db.GetChannelSubscribers(ctx, channelID)
	if err != nil {
		h.log.Error("failed to get channel subscribers", "call", "db.GetChannelSubscribers", "err", err)
		return err
	}

	for _, sub := range subs {
		err = h.db.AddVideoToUser(ctx, sub, videoID)
		if err != nil {
			h.log.Error("Failed to add video to user", "call", "db.AddVideoToUser", "err", err)
			continue
		}
	}

	return nil
}

func (h *handler) isShort(ctx context.Context, videoID string) (bool, error) {
	url := fmt.Sprintf("https://www.youtube.com/shorts/%s", videoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		h.log.Error(
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
		h.log.Error("Failed to check if a video is a short", "videoID", videoID, "error", err)
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (h *handler) addVideosForChannel(ctx context.Context, parsedChannel *feedparser.Channel) {
	for _, parsedVideo := range parsedChannel.Videos {
		date, err := parsedVideo.Published.Parse()
		if err != nil {
			h.log.Warn("Failed to parse video information", "call", "feedparser.Parse", "err", err)
			continue
		}

		videoID := strings.Split(parsedVideo.ID, ":")[2]
		isShort, err := h.isShort(ctx, videoID)
		if err != nil {
			isShort = false
		}
		video := models.Video{
			ID:            videoID,
			Title:         parsedVideo.Title,
			PublishedTime: date,
			IsShort:       isShort,
		}
		err = h.db.AddVideo(ctx, video, parsedChannel.ID)
		if err != nil {
			if !errors.Is(err, db.ErrVideoExists) {
				h.log.Warn("Failed to save video to db", "call", "db.AddVideo", "err", err)
			}
			continue
		}
		err = h.addVideoToAllSubscribers(ctx, parsedChannel.ID, videoID)
		if err != nil {
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

func (h *handler) MarkVideoAsWatched(ctx context.Context, username string, videoID string) error {
	watchTime := time.Now()
	return h.db.SetVideoWatchTime(ctx, username, videoID, &watchTime)
}

func (h *handler) MarkVideoAsUnwatched(ctx context.Context, username string, videoID string) error {
	return h.db.SetVideoWatchTime(ctx, username, videoID, nil)
}
