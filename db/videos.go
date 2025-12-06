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

func (d *postgresDB) GetWatchedVideos(ctx context.Context, sortDesc bool) ([]models.Video, error) {
	query := `
		SELECT
			videos.id
			, title
			, published_timestamp
			, watch_timestamp
			, is_short
			, channels.name
			, channels.id
		FROM videos
		LEFT JOIN channels ON videos.channel_id=channels.id
		WHERE watch_timestamp IS NOT NULL
		ORDER BY watch_timestamp DESC
	`
	if sortDesc {
		query += " DESC"
	}

	rows, err := d.db.QueryContext(ctx, query)
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
			&video.ChannelName,
		)
		if err != nil {
			d.l.Error("Failed to scan rows for watched videos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (d *postgresDB) AddVideo(ctx context.Context, video models.Video, channelID string) error {
	query := `
		INSERT INTO videos (
			id
			, title
			, published_timestamp
			, is_short
			, channel_id
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
	`

	resp, err := d.db.ExecContext(
		ctx,
		query,
		video.ID,
		video.Title,
		video.PublishedTime,
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

const setVideoWatchTimeQuery = `UPDATE videos SET watch_timestamp = $1 WHERE id = $2`

func (d *postgresDB) SetVideoWatchTime(
	ctx context.Context,
	videoID string,
	watchTime *time.Time,
) error {
	_, err := d.db.ExecContext(ctx, setVideoWatchTimeQuery, watchTime, videoID)
	if err != nil {
		d.l.Error("", "error", err)
		return err
	}

	return nil
}
