ALTER TABLE videos ADD COLUMN IF NOT EXISTS watch_timestamp timestamp with time zone DEFAULT NULL;

UPDATE videos
  SET watch_timestamp = user_videos.watch_timestamp
  FROM user_videos
  WHERE videos.id = user_videos.video_id;

DROP TABLE IF EXISTS user_subscriptions;
DROP TABLE IF EXISTS user_videos;
DROP TABLE IF EXISTS users;
