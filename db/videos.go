package db

import (
	"context"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (d *postgresDB) GetNewVideos(ctx context.Context, sortDesc bool) ([]models.Video, error) {
	query := `
		SELECT
			videos.id
			, title
			, published_timestamp
			, is_short
			, duration
			, progress
			, channels.name
			, channels.id
		FROM videos
		LEFT JOIN channels ON videos.channel_id=channels.id
		WHERE watch_timestamp IS NULL
		ORDER BY published_timestamp
	`
	if sortDesc {
		query += " DESC"
	}

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		d.l.Error("Failed to query new videos", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	videos := make([]models.Video, 0)
	for rows.Next() {
		var video models.Video
		err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.PublishedTime,
			&video.IsShort,
			&video.DurationSeconds,
			&video.ProgressSeconds,
			&video.ChannelName,
			&video.ChannelID,
		)
		if err != nil {
			d.l.Error("Failed to scan rows for get new videos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (d *postgresDB) GetWatchedVideos(
	ctx context.Context, sortDesc bool, limit int, offset int,
) ([]models.Video, error) {
	query := `
		SELECT
			videos.id
			, title
			, published_timestamp
			, watch_timestamp
			, is_short
			, duration
			, progress
			, channels.name
			, channels.id
		FROM videos
		LEFT JOIN channels ON videos.channel_id=channels.id
		WHERE watch_timestamp IS NOT NULL
		ORDER BY watch_timestamp
	`
	if sortDesc {
		query += " DESC"
	}
	query += " LIMIT $1 OFFSET $2"

	rows, err := d.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		d.l.Error("Failed to query for watched videos", "call", "sql.QueryContext", "error", err)
		return nil, err
	}
	defer rows.Close()

	videos := make([]models.Video, 0)
	for rows.Next() {
		var video models.Video
		err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.PublishedTime,
			&video.WatchTime,
			&video.IsShort,
			&video.DurationSeconds,
			&video.ProgressSeconds,
			&video.ChannelName,
			&video.ChannelID,
		)
		if err != nil {
			d.l.Error("Failed to scan rows for watched videos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (d *postgresDB) HasVideo(ctx context.Context, videoID string) (bool, error) {
	query := `SELECT COUNT(1) FROM videos WHERE id = $1`
	row := d.db.QueryRowContext(ctx, query, videoID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		d.l.Error("Failed to query for video", "call", "sql.QueryRowContext", "error", err)
		return false, err
	}

	return count == 1, nil
}

func (d *postgresDB) AddVideo(ctx context.Context, video models.Video, channelID string) error {
	query := `
		INSERT INTO videos (
			id
			, title
			, published_timestamp
			, duration
			, is_short
			, channel_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`

	resp, err := d.db.ExecContext(
		ctx,
		query,
		video.ID,
		video.Title,
		video.PublishedTime,
		video.DurationSeconds,
		video.IsShort,
		channelID,
	)
	if err != nil {
		d.l.Error("Failed to add video", "call", "sql.Exec", "error", err)
		return err
	}
	if affected, _ := resp.RowsAffected(); affected == 0 {
		return ErrVideoExists
	}

	return nil
}

func (d *postgresDB) SetVideoWatchTime(
	ctx context.Context,
	videoID string,
	watchTime *time.Time,
) error {
	const query = `UPDATE videos SET watch_timestamp = $1 WHERE id = $2`
	_, err := d.db.ExecContext(ctx, query, watchTime, videoID)
	if err != nil {
		d.l.Error("", "error", err)
		return err
	}

	return nil
}

func (d *postgresDB) SetVideoProgress(ctx context.Context, videoID string, progress int) (*models.Video, error) {
	const query = `
		WITH updated AS (
			UPDATE videos SET progress = $1 WHERE id = $2 RETURNING *
		)
		SELECT
			updated.id
			, title
			, published_timestamp
			, is_short
			, duration
			, progress
			, channels.name
			, channels.id
		FROM updated
		LEFT JOIN channels ON updated.channel_id = channels.id
	`

	row := d.db.QueryRowContext(ctx, query, progress, videoID)
	var video models.Video
	err := row.Scan(
		&video.ID,
		&video.Title,
		&video.PublishedTime,
		&video.IsShort,
		&video.DurationSeconds,
		&video.ProgressSeconds,
		&video.ChannelName,
		&video.ChannelID,
	)
	if err != nil {
		d.l.Error("Failed to scan rows for get new videos", "call", "sql.Scan", "error", err)
		return nil, err
	}

	return &video, nil
}
