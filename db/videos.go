package db

import (
	"context"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (db *postgresDB) GetNewVideos(ctx context.Context, sortDesc bool) ([]models.Video, error) {
	query := `
		SELECT
			videos.id
			, title
			, published_timestamp
			, is_short
			, duration
			, progress
			, downloaded_at
			, file_path
			, download_status
			, download_error
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

	rows, err := db.db.Query(ctx, query)
	if err != nil {
		db.l.Error("Failed to query new videos", "call", "sql.QueryContext", "error", err)
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
			&video.DownloadedAt,
			&video.FilePath,
			&video.DownloadStatus,
			&video.DownloadError,
			&video.ChannelName,
			&video.ChannelID,
		)
		if err != nil {
			db.l.Error("Failed to scan rows for get new videos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *postgresDB) GetWatchedVideos(
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
			, downloaded_at
			, file_path
			, download_status
			, download_error
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

	rows, err := db.db.Query(ctx, query, limit, offset)
	if err != nil {
		db.l.Error("Failed to query for watched videos", "call", "sql.QueryContext", "error", err)
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
			&video.DownloadedAt,
			&video.FilePath,
			&video.DownloadStatus,
			&video.DownloadError,
			&video.ChannelName,
			&video.ChannelID,
		)
		if err != nil {
			db.l.Error("Failed to scan rows for watched videos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *postgresDB) HasVideo(ctx context.Context, videoID string) (bool, error) {
	query := `SELECT COUNT(1) FROM videos WHERE id = $1`
	row := db.db.QueryRow(ctx, query, videoID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		db.l.Error("Failed to query for video", "call", "sql.QueryRowContext", "error", err)
		return false, err
	}

	return count == 1, nil
}

func (db *postgresDB) AddVideo(ctx context.Context, video models.Video, channelID string) error {
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

	resp, err := db.db.Exec(
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
		db.l.Error("Failed to add video", "call", "sql.Exec", "error", err)
		return err
	}
	if resp.RowsAffected() == 0 {
		return ErrVideoExists
	}

	return nil
}

func (db *postgresDB) SetVideoWatchTime(
	ctx context.Context,
	videoID string,
	watchTime *time.Time,
) error {
	const query = `UPDATE videos SET watch_timestamp = $1 WHERE id = $2`
	_, err := db.db.Exec(ctx, query, watchTime, videoID)
	if err != nil {
		db.l.Error("", "error", err)
		return err
	}

	return nil
}

func (db *postgresDB) SetVideoProgress(ctx context.Context, videoID string, progress int) (*models.Video, error) {
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

	row := db.db.QueryRow(ctx, query, progress, videoID)
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
		db.l.Error("Failed to scan rows for get new videos", "call", "sql.Scan", "error", err)
		return nil, err
	}

	return &video, nil
}

func (db *postgresDB) GetVideo(ctx context.Context, videoID string) (*models.Video, error) {
	query := `
		SELECT
			videos.id
			, title
			, published_timestamp
			, is_short
			, duration
			, progress
			, watch_timestamp
			, downloaded_at
			, file_path
			, download_status
			, download_error
			, channels.name
			, channels.id
		FROM videos
		LEFT JOIN channels ON videos.channel_id = channels.id
		WHERE videos.id = $1
	`

	row := db.db.QueryRow(ctx, query, videoID)
	var video models.Video
	err := row.Scan(
		&video.ID,
		&video.Title,
		&video.PublishedTime,
		&video.IsShort,
		&video.DurationSeconds,
		&video.ProgressSeconds,
		&video.WatchTime,
		&video.DownloadedAt,
		&video.FilePath,
		&video.DownloadStatus,
		&video.DownloadError,
		&video.ChannelName,
		&video.ChannelID,
	)
	if err != nil {
		db.l.Error("Failed to query video", "call", "sql.QueryRow", "error", err)
		return nil, err
	}

	return &video, nil
}

func (db *postgresDB) SetVideoDownloadStatus(ctx context.Context, videoID string, status string) error {
	query := `UPDATE videos SET download_status = $1 WHERE id = $2`
	_, err := db.db.Exec(ctx, query, status, videoID)
	if err != nil {
		db.l.Error("Failed to set video download status", "error", err)
		return err
	}
	return nil
}

func (db *postgresDB) SetVideoDownloadCompleted(ctx context.Context, videoID string, filePath string) error {
	query := `
		UPDATE videos
		SET
			downloaded_at = $1,
			file_path = $2,
			download_status = $3,
			download_error = NULL
		WHERE id = $4
	`
	now := time.Now()
	_, err := db.db.Exec(ctx, query, now, filePath, "completed", videoID)
	if err != nil {
		db.l.Error("Failed to mark video as downloaded", "error", err)
		return err
	}
	return nil
}

func (db *postgresDB) SetVideoDownloadFailed(ctx context.Context, videoID string, errorMsg string) error {
	query := `UPDATE videos SET download_status = $1, download_error = $2 WHERE id = $3`
	_, err := db.db.Exec(ctx, query, "failed", errorMsg, videoID)
	if err != nil {
		db.l.Error("Failed to mark video download as failed", "error", err)
		return err
	}
	return nil
}

func (db *postgresDB) GetVideosForCleanup(ctx context.Context, olderThan time.Duration) ([]models.Video, error) {
	query := `
		SELECT
			id
			, file_path
			, downloaded_at
			, watch_timestamp
		FROM videos
		WHERE downloaded_at IS NOT NULL
			AND file_path IS NOT NULL
			AND download_status = 'completed'
			AND watch_timestamp IS NOT NULL
			AND watch_timestamp < $1
	`
	cutoffTime := time.Now().Add(-olderThan)

	rows, err := db.db.Query(ctx, query, cutoffTime)
	if err != nil {
		db.l.Error("Failed to query videos for cleanup", "error", err)
		return nil, err
	}
	defer rows.Close()

	videos := make([]models.Video, 0)
	for rows.Next() {
		var video models.Video
		err = rows.Scan(&video.ID, &video.FilePath, &video.DownloadedAt, &video.WatchTime)
		if err != nil {
			db.l.Error("Failed to scan cleanup video", "error", err)
			continue
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *postgresDB) DeleteVideoFile(ctx context.Context, videoID string) error {
	query := `
		UPDATE videos
		SET
			downloaded_at = NULL,
			file_path = NULL,
			download_status = NULL,
			download_error = NULL
		WHERE id = $1
	`
	_, err := db.db.Exec(ctx, query, videoID)
	if err != nil {
		db.l.Error("Failed to clear video download fields", "error", err)
		return err
	}
	return nil
}
