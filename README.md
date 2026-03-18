<p align="center">
  <a href="https://github.com/TheEdgeOfRage/ytrssil">
    <picture>
      <img src="assets/ytrssil.svg" alt="ytrssil logo" width="100">
    </picture>
  </a>
</p>

<h1 align="center">ytrssil</h1>
<p align="center">YouTube subscription feed, but better</p>

---

## What is ytrssil?

ytrssil is a self-hosted YouTube subscription manager that gives you a clean feed of your subscribed channels and nothing else. It tracks watched videos, handles downloads, and gives you control over the videos you get served.

## Quick Start

### Docker (Recommended)

```bash
git clone https://github.com/TheEdgeOfRage/ytrssil
cd ytrssil
docker-compose up -d
```

Visit `http://localhost:8080` in your browser.

## Configuration

ytrssil is configured via environment variables. Edit the env vars in compose.yaml

```bash
# PostgreSQL connection (optional, defaults to Docker Compose service)
POSTGRES_URL=postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable

# Where downloaded videos are saved
DOWNLOADS_DIR=/var/lib/ytrssil/downloads

# How often to check for new videos (default: 5m)
FETCH_INTERVAL=5m

# How often to cleanup old downloads (default: 1h)
CLEANUP_INTERVAL=1h
```

## Usage

### Adding Channels

1. Click the **"Add Channel"** button
2. Paste a YouTube channel URL or search by name
3. The channel appears in your subscription list

### Managing Videos

- **Watch status**: Click a video to mark it as watched
- **Downloads**: Click the download button to save videos locally
- **Shorts filter**: Toggle the shorts switch on each channel to filter out YouTube Shorts
- **Progress**: The dashboard shows unwatched counts and recent activity

### Download Settings

Downloaded videos are automatically cleaned up 2 days (configurable) after marking it as watched.

## Features

- **Channel subscriptions** - Add channels via channel name or ID
- **Watch History** - Keep a list of what you've watched
- **Progress Tracking** - Keep track of your watch progress in videos
- **Video Downloads** - Save videos locally with automatic cleanup
- **Shorts Filter** - Per-channel control over YouTube Shorts
- **Clean Interface** - No ads, no recommendations, just your feed
- **Auto Updates** - Checks for new videos every 5 minutes
- **Docker Ready** - One command to get everything running

## Support

For development setup, code architecture, or contributing, see [`AGENTS.md`](AGENTS.md).
