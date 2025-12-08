ALTER TABLE channels ADD COLUMN IF NOT EXISTS subscribed boolean NOT NULL DEFAULT true;
