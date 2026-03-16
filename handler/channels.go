package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

// isChannelID reports whether s looks like a raw YouTube channel ID (UCxxxxxxxx…).
func isChannelID(s string) bool {
	return strings.HasPrefix(s, "UC") && len(s) == 24
}

func (h *handler) SubscribeToChannel(ctx context.Context, channelID string) (*models.Channel, error) {
	if !isChannelID(channelID) {
		// Treat input as a handle; normalise to @handle form for the API.
		handle := channelID
		if !strings.HasPrefix(handle, "@") {
			handle = "@" + handle
		}
		resolved, err := h.youTubeClient.ResolveChannelID(ctx, handle)
		if err != nil {
			return nil, err
		}
		channelID = resolved
	}

	parsedChannel, err := h.parser.Parse(channelID)
	if err != nil {
		return nil, err
	}

	imageURL, err := h.youTubeClient.GetChannelImageURL(ctx, channelID)
	if err != nil {
		h.log.Warn("Failed to fetch channel image URL", "channelID", channelID, "error", err)
		imageURL = ""
	}

	channel := models.Channel{
		ID:           channelID,
		Name:         parsedChannel.Name,
		Subscribed:   true,
		ImageURL:     imageURL,
		EnableShorts: true,
	}

	err = h.db.SubscribeToChannel(ctx, channel)
	if err != nil && !errors.Is(err, db.ErrChannelExists) {
		return nil, err
	}

	return &channel, nil
}

func (h *handler) UnsubscribeFromChannel(ctx context.Context, channelID string) error {
	return h.db.UnsubscribeFromChannel(ctx, channelID)
}

func (h *handler) ListChannels(ctx context.Context) ([]models.Channel, error) {
	return h.db.ListChannels(ctx)
}

func (h *handler) GetChannelByID(ctx context.Context, channelID string) (*models.Channel, error) {
	return h.db.GetChannelByID(ctx, channelID)
}

func (h *handler) ToggleChannelShorts(ctx context.Context, channelID string, enableShorts bool) error {
	return h.db.ToggleChannelShorts(ctx, channelID, enableShorts)
}
