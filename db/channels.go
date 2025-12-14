package db

import (
	"context"
	"fmt"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (d *postgresDB) SubscribeToChannel(ctx context.Context, channel models.Channel) error {
	const query = `
		INSERT INTO channels (id, name, subscribed, image_url) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET subscribed = $3, image_url = $4
	`
	resp, err := d.db.ExecContext(ctx, query, channel.ID, channel.Name, channel.Subscribed, channel.ImageURL)
	if err != nil {
		d.l.Error("Failed to subscribe to channel", "call", "sql.ExecContext", "error", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		d.l.Error("Failed to subscribe to channel", "call", "sql.RowsAffected", "error", err)
		return err
	}
	if rows == 0 {
		d.l.Error("Failed to subscribe to channel, no rows affected", "call", "sql.RowsAffected")
		return fmt.Errorf("failed to subscribe to channel")
	}

	return nil
}

func (d *postgresDB) ListChannels(ctx context.Context) ([]models.Channel, error) {
	const query = `
		SELECT
			channels.id,
			channels.name,
			channels.subscribed,
			COALESCE(channels.image_url, '') as image_url,
			COUNT(videos.id) FILTER (WHERE videos.watch_timestamp IS NULL) as unwatched_count
		FROM channels
		LEFT JOIN videos ON channels.id = videos.channel_id
		WHERE channels.subscribed = true
		GROUP BY channels.id, channels.name, channels.subscribed, channels.image_url
		ORDER BY channels.name
	`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		d.l.Error("Failed to list channels", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	channels := make([]models.Channel, 0)
	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.Name, &channel.Subscribed, &channel.ImageURL, &channel.UnwatchedCount)
		if err != nil {
			d.l.Error("Failed to scan rows to list channels", "call", "sql.Scan", "error", err)
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

func (d *postgresDB) UnsubscribeFromChannel(ctx context.Context, channelID string) error {
	const query = `UPDATE channels SET subscribed = false WHERE id = $1`
	resp, err := d.db.ExecContext(ctx, query, channelID)
	if err != nil {
		d.l.Error("Failed to unsubscribe from channel", "call", "sql.ExecContext", "error", err)
		return err
	}

	if affected, err := resp.RowsAffected(); err != nil || affected != 1 {
		return ErrChannelNotFound
	}

	return nil
}
