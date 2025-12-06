ALTER TABLE videos DROP COLUMN IF EXISTS watch_timestamp;

CREATE TABLE IF NOT EXISTS users (
	username text NOT NULL PRIMARY KEY
	, password text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_videos (
	username text NOT NULL REFERENCES users(username)
	, video_id text NOT NULL REFERENCES videos(id)
	, watch_timestamp timestamp with time zone
	, CONSTRAINT user_videos_pkey PRIMARY KEY (username, video_id)
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
	username text NOT NULL REFERENCES users(username)
	, channel_id text NOT NULL REFERENCES channels(id)
	, CONSTRAINT user_subscriptions_pkey PRIMARY KEY (channel_id, username)
);
