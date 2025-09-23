package handler

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) GetNewVideos(ctx context.Context, username string) ([]models.Video, error) {
	return h.db.GetNewVideos(ctx, username)
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

func (h *handler) addVideosForChannel(ctx context.Context, parsedChannel *feedparser.Channel) {
	for _, parsedVideo := range parsedChannel.Videos {
		date, err := parsedVideo.Published.Parse()
		if err != nil {
			h.log.Warn("Failed to parse video information", "call", "feedparser.Parse", "err", err)
			continue
		}

		videoID := strings.Split(parsedVideo.ID, ":")[2]
		video := models.Video{
			ID:            videoID,
			Title:         parsedVideo.Title,
			PublishedTime: date,
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
