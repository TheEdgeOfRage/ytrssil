package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

const WatchedVideosPageSize = 100

func (h *handler) GetWatchedVideos(ctx context.Context, sortDesc bool, page int) ([]models.Video, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * WatchedVideosPageSize
	return h.db.GetWatchedVideos(ctx, sortDesc, WatchedVideosPageSize, offset)
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

type parseResult struct {
	channel *feedparser.Channel
	err     error
}

func (h *handler) FetchVideos(ctx context.Context) error {
	h.log.Info("Fetching new videos for all channels")

	channels, err := h.db.ListChannels(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	results := make(chan parseResult, 1)
	for _, channel := range channels {
		wg.Go(func() {
			parsedChannel, err := h.parser.Parse(channel.ID)
			results <- parseResult{channel: parsedChannel, err: err}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			h.log.Error("failed to parse channel feed", "error", result.err)
			continue
		}
		h.addVideosForChannel(ctx, result.channel)
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

// parseTimeProgress parses a duration string in  the Go duration format, hh:mm:ss, and mm:ss to a time.Duration
func parseTimeProgress(progressTime string) (time.Duration, error) {
	// Try Go duration format first
	duration, err := time.ParseDuration(progressTime)
	if err == nil {
		return duration, nil
	}

	// Try hh:mm:ss or mm:ss format
	parts := strings.Split(progressTime, ":")
	if len(parts) == 2 {
		// mm:ss format
		minutes, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid mm:ss format")
		}
		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid mm:ss format")
		}
		return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
	} else if len(parts) == 3 {
		// hh:mm:ss format
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid hh:mm:ss format")
		}
		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid hh:mm:ss format")
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, fmt.Errorf("invalid hh:mm:ss format")
		}
		return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
	}

	return 0, fmt.Errorf("unsupported time format: expected Go duration, hh:mm:ss, or mm:ss")
}

func (h *handler) SetVideoProgress(ctx context.Context, videoID string, progressTime string) (*models.Video, error) {
	progress, err := parseTimeProgress(progressTime)
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
