package db

import (
	"context"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (d *postgresDB) SubscribeToChannel(ctx context.Context, channel models.Channel) error {
	const query = `INSERT INTO channels (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	resp, err := d.db.ExecContext(ctx, query, channel.ID, channel.Name)
	if err != nil {
		d.l.Error("Failed to create channel", "call", "sql.ExecContext", "error", err)
		return err
	}
	if affected, _ := resp.RowsAffected(); affected == 0 {
		return ErrChannelExists
	}

	return nil
}

func (d *postgresDB) ListChannels(ctx context.Context) ([]models.Channel, error) {
	const query = `SELECT id, name FROM channels`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		d.l.Error("Failed to list channels", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	channels := make([]models.Channel, 0)
	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.Name)
		if err != nil {
			d.l.Error("Failed to scan rows to list channels", "call", "sql.Scan", "error", err)
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

func (d *postgresDB) UnsubscribeFromChannel(ctx context.Context, channelID string) error {
	const query = `DELETE FROM channels WHERE id = $1`
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
