package handler

import (
	"context"
	"errors"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) SubscribeToChannel(ctx context.Context, channelID string) error {
	parsedChannel, err := h.parser.Parse(channelID)
	if err != nil {
		return err
	}

	channel := models.Channel{
		ID:   channelID,
		Name: parsedChannel.Name,
	}

	err = h.db.SubscribeToChannel(ctx, channel)
	if err != nil && !errors.Is(err, db.ErrChannelExists) {
		return err
	}

	return nil
}

func (h *handler) UnsubscribeFromChannel(ctx context.Context, channelID string) error {
	return h.db.UnsubscribeFromChannel(ctx, channelID)
}
