package handler

import (
	"context"
	"errors"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) SubscribeToChannel(ctx context.Context, channelID string) (*models.Channel, error) {
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
		ID:         channelID,
		Name:       parsedChannel.Name,
		Subscribed: true,
		ImageURL:   imageURL,
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
