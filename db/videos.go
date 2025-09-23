package db

import (
	"context"
	"time"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var getNewVideosQuery = `
	SELECT
		video_id
		, title
		, published_timestamp
		, watch_timestamp
		, name as channel_name
	FROM user_videos
	LEFT JOIN videos ON video_id=videos.id
	LEFT JOIN channels ON channel_id=channels.id
	WHERE
		1=1
		AND watch_timestamp IS NULL
		AND username=$1
	ORDER BY published_timestamp
`

func (d *postgresDB) GetNewVideos(ctx context.Context, username string) ([]models.Video, error) {
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
			&video.WatchTime,
			&video.ChannelName,
		)
		if err != nil {
			d.l.Error("Failed to scan rows for get new videos", "call", "sql.Scan", "error", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

var getWatchedVideosQuery = `
	SELECT
		video_id
		, title
		, published_timestamp
		, watch_timestamp
		, name as channel_name
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

var addVideoQuery = `
INSERT INTO videos (
	id
	, title
	, published_timestamp
	, channel_id
) VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
`

func (d *postgresDB) AddVideo(ctx context.Context, video models.Video, channelID string) error {
	resp, err := d.db.ExecContext(ctx, addVideoQuery, video.ID, video.Title, video.PublishedTime, channelID)
	if err != nil {
		d.l.Error("Failed to add video", "call", "sql.Exec", "error", err)
		return err
	}
	if affected, _ := resp.RowsAffected(); affected == 0 {
		return ErrVideoExists
	}

	return nil
}

var addVideoToUserQuery = `INSERT INTO user_videos (username, video_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

func (d *postgresDB) AddVideoToUser(ctx context.Context, username string, videoID string) error {
	_, err := d.db.ExecContext(ctx, addVideoToUserQuery, username, videoID)
	if err != nil {
		d.l.Error("Failed to add video to user", "error", err)
		return err
	}

	return nil
}

var setVideoWatchTimeQuery = `UPDATE user_videos SET watch_timestamp = $1 WHERE username = $2 AND video_id = $3`

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
