package handler

import (
	"context"
	"errors"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) SubscribeToChannel(ctx context.Context, username string, channelID string) error {
	parsedChannel, err := h.parser.Parse(channelID)
	if err != nil {
		return err
	}

	channel := models.Channel{
		ID:   channelID,
		Name: parsedChannel.Name,
	}

	err = h.db.CreateChannel(ctx, channel)
	if err != nil && !errors.Is(err, db.ErrChannelExists) {
		return err
	}

	return h.db.SubscribeUserToChannel(ctx, username, channelID)
}

func (h *handler) UnsubscribeFromChannel(ctx context.Context, username string, channelID string) error {
	return h.db.UnsubscribeUserFromChannel(ctx, username, channelID)
}
