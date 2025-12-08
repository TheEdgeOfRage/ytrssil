package db

import (
	"context"
	"fmt"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (d *postgresDB) SubscribeToChannel(ctx context.Context, channel models.Channel) error {
	const query = `
		INSERT INTO channels (id, name, subscribed) VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET subscribed = $3
	`
	resp, err := d.db.ExecContext(ctx, query, channel.ID, channel.Name, channel.Subscribed)
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
	const query = `SELECT id, name, subscribed FROM channels WHERE subscribed = true`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		d.l.Error("Failed to list channels", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	channels := make([]models.Channel, 0)
	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.Name, &channel.Subscribed)
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
