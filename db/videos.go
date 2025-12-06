package db

import (
	"context"
	"time"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (d *postgresDB) GetNewVideos(ctx context.Context, username string, sortDesc bool) ([]models.Video, error) {
	getNewVideosQuery := `
		SELECT
			videos.id
			, videos.title
			, videos.published_timestamp
			, videos.is_short
			, channels.name as channel_name
			, channels.id as channel_id
		FROM user_videos
		LEFT JOIN videos ON video_id=videos.id
		LEFT JOIN channels ON channel_id=channels.id
		WHERE
			1=1
			AND watch_timestamp IS NULL
			AND username=$1
		ORDER BY published_timestamp
	`
	if sortDesc {
		getNewVideosQuery += "DESC"
	}

	rows, err := d.db.QueryContext(ctx, getNewVideosQuery, username)
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

const getWatchedVideosQuery = `
	SELECT
		videos.id
		, videos.title
		, videos.published_timestamp
		, videos.watch_timestamp
		, videos.is_short
		, channels.name as channel_name
		, channels.id as channel_id
	FROM user_videos
	LEFT JOIN videos ON video_id=videos.id
	LEFT JOIN channels ON channel_id=channels.id
	WHERE
		1=1
		AND watch_timestamp IS NOT NULL
		AND username=$1
	ORDER BY watch_timestamp DESC
`

func (d *postgresDB) GetWatchedVideos(ctx context.Context, username string) ([]models.Video, error) {
	rows, err := d.db.QueryContext(ctx, getWatchedVideosQuery, username)
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

const addVideoQuery = `
INSERT INTO videos (
	id
	, title
	, published_timestamp
	, is_short
	, channel_id
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING
`

func (d *postgresDB) AddVideo(ctx context.Context, video models.Video, channelID string) error {
	resp, err := d.db.ExecContext(
		ctx,
		addVideoQuery,
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

const addVideoToUserQuery = `INSERT INTO user_videos (username, video_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

func (d *postgresDB) AddVideoToUser(ctx context.Context, username string, videoID string) error {
	_, err := d.db.ExecContext(ctx, addVideoToUserQuery, username, videoID)
	if err != nil {
		d.l.Error("Failed to add video to user", "error", err)
		return err
	}

	return nil
}

const setVideoWatchTimeQuery = `UPDATE user_videos SET watch_timestamp = $1 WHERE username = $2 AND video_id = $3`

func (d *postgresDB) SetVideoWatchTime(
	ctx context.Context, username string, videoID string, watchTime *time.Time,
) error {
	_, err := d.db.ExecContext(ctx, setVideoWatchTimeQuery, watchTime, username, videoID)
	if err != nil {
		d.l.Error("", "error", err)
		return err
	}

	return nil
}
