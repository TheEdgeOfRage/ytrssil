package handler

import (
	"context"
	"errors"
	"fmt"
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
			h.log.Error("Failed to check if video already exists", "call", "db.HasVideo", "err", err)
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

	if len(videos) == 0 {
		return
	}

	err = h.youTubeClient.GetVideoDurations(ctx, videos)
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

	video, err := h.db.SetVideoProgress(ctx, videoID, int(progress.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidProgress, err.Error())
	}

	return video, nil
}

func (h *handler) AddCustomVideo(ctx context.Context, videoID string) error {
	exists, err := h.db.HasVideo(ctx, videoID)
	if err != nil {
		h.log.Error("Failed to check if video already exists", "call", "db.HasVideo", "err", err)
		return err
	}

	if exists {
		h.log.Warn("Video already in db", "call", "db.HasVideo")
		return nil
	}

	video, err := h.youTubeClient.GetVideoMetadata(ctx, videoID)
	if err != nil {
		h.log.Error("Failed to get video metadata", "error", err)
		return err
	}

	channel := models.Channel{
		ID:         video.ChannelID,
		Name:       video.ChannelName,
		Subscribed: false,
	}
	err = h.db.SubscribeToChannel(ctx, channel)
	if err != nil {
		h.log.Error("Failed to insert channel", "error", err)
		return err
	}

	err = h.db.AddVideo(ctx, *video, video.ChannelID)
	if err != nil {
		if !errors.Is(err, db.ErrVideoExists) {
			h.log.Error("Failed to save video to db", "call", "db.AddVideo", "err", err)
			return err
		}
	}

	return nil
}
