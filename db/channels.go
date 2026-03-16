package db

import (
	"context"
	"fmt"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (db *postgresDB) SubscribeToChannel(ctx context.Context, channel models.Channel) error {
	const query = `
		INSERT INTO channels (id, name, subscribed, image_url, enable_shorts) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET subscribed = $3, image_url = $4, enable_shorts = $5
	`
	resp, err := db.db.Exec(ctx, query, channel.ID, channel.Name, channel.Subscribed,
		channel.ImageURL, channel.EnableShorts)
	if err != nil {
		db.l.Error("Failed to subscribe to channel", "call", "sql.ExecContext", "error", err)
		return err
	}

	if resp.RowsAffected() == 0 {
		db.l.Error("Failed to subscribe to channel, no rows affected", "call", "sql.RowsAffected")
		return fmt.Errorf("failed to subscribe to channel")
	}

	return nil
}

func (db *postgresDB) ListChannels(ctx context.Context) ([]models.Channel, error) {
	const query = `
		SELECT
			channels.id,
			channels.name,
			channels.subscribed,
			COALESCE(channels.image_url, '') as image_url,
			COALESCE(channels.enable_shorts, true) as enable_shorts,
			COUNT(videos.id) FILTER (WHERE videos.watch_timestamp IS NULL AND videos.is_discarded = false) as unwatched_count
		FROM channels
		LEFT JOIN videos ON channels.id = videos.channel_id
		WHERE channels.subscribed = true
		GROUP BY channels.id, channels.name, channels.subscribed, channels.image_url, channels.enable_shorts
		ORDER BY channels.name
	`
	rows, err := db.db.Query(ctx, query)
	if err != nil {
		db.l.Error("Failed to list channels", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	channels := make([]models.Channel, 0)
	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.Name, &channel.Subscribed, &channel.ImageURL,
			&channel.EnableShorts, &channel.UnwatchedCount)
		if err != nil {
			db.l.Error("Failed to scan rows to list channels", "call", "sql.Scan", "error", err)
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

func (db *postgresDB) UnsubscribeFromChannel(ctx context.Context, channelID string) error {
	const query = `UPDATE channels SET subscribed = false WHERE id = $1`
	resp, err := db.db.Exec(ctx, query, channelID)
	if err != nil {
		db.l.Error("Failed to unsubscribe from channel", "call", "sql.ExecContext", "error", err)
		return err
	}

	if resp.RowsAffected() != 1 {
		return ErrChannelNotFound
	}

	return nil
}

func (db *postgresDB) ToggleChannelShorts(ctx context.Context, channelID string, enableShorts bool) error {
	const query = `UPDATE channels SET enable_shorts = $1 WHERE id = $2`
	resp, err := db.db.Exec(ctx, query, enableShorts, channelID)
	if err != nil {
		db.l.Error("Failed to toggle channel shorts", "call", "sql.ExecContext", "error", err)
		return err
	}

	if resp.RowsAffected() != 1 {
		return ErrChannelNotFound
	}

	return nil
}

func (db *postgresDB) GetChannelByID(ctx context.Context, channelID string) (*models.Channel, error) {
	const query = `
		SELECT
			id,
			name,
			subscribed,
			COALESCE(image_url, '') as image_url,
			COALESCE(enable_shorts, true) as enable_shorts
		FROM channels
		WHERE id = $1
	`
	row := db.db.QueryRow(ctx, query, channelID)
	var channel models.Channel
	err := row.Scan(&channel.ID, &channel.Name, &channel.Subscribed, &channel.ImageURL, &channel.EnableShorts)
	if err != nil {
		db.l.Error("Failed to query channel by ID", "call", "sql.QueryRowContext", "error", err)
		return nil, err
	}

	return &channel, nil
}
