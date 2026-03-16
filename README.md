<p align="center">
  <a href="https://github.com/TheEdgeOfRage/ytrssil">
    <picture>
      <img src="assets/ytrssil.svg" alt="ytrssil logo" width="200">
    </picture>
  </a>
</p>

<h1 align="center">ytrssil</h1>
<p align="center">YouTube subscription feed, but better</p>

---

## What is ytrssil?

ytrssil is a self-hosted YouTube subscription manager that gives you a clean, ad-free feed of your subscribed channels. It tracks watched videos, handles downloads, and gives you full control over your YouTube experience.

## Quick Start

### Docker (Recommended)

```bash
git clone https://github.com/TheEdgeOfRage/ytrssil
cd ytrssil
docker-compose up -d
```

Visit `http://localhost:8080` in your browser.

That's it. The container includes everything you need: PostgreSQL database, the Go application, and automatic migrations.

## Configuration

ytrssil is configured via environment variables. Create a `.env` file in the project directory:

```bash
# PostgreSQL connection (optional, defaults to Docker Compose service)
POSTGRES_URL=postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable

# Server port
SERVER_PORT=8080

# Where downloaded videos are saved
DOWNLOAD_PATH=/downloads

# How many days to keep downloaded videos (default: 2)
CLEANUP_DAYS=2

# How often to check for new videos (default: 5m)
FETCH_INTERVAL=5m

# How often to cleanup old downloads (default: 1h)
CLEANUP_INTERVAL=1h

# Disable background jobs (useful for development)
DISABLE_JOBS=false
```

### Docker Compose

The included `docker-compose.yml` sets up PostgreSQL and the app automatically. Just run:

```bash
docker-compose up -d
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

Downloaded videos are automatically cleaned up after the configured number of days (default: 2 days). This keeps your disk usage in check while still giving you access to your content.

## Features

- **Smart Subscriptions** - Add channels via URL or search
- **Watch Tracking** - Never lose track of what you've seen
- **Video Downloads** - Save videos locally with automatic cleanup
- **Shorts Filter** - Per-channel control over YouTube Shorts
- **Clean Interface** - No ads, no recommendations, just your feed
- **Auto Updates** - Checks for new videos every 5 minutes
- **Docker Ready** - One command to get everything running

## Project Structure

```
ytrssil/
├── cmd/              # Application entry point
├── handler/          # Business logic
├── httpserver/       # HTTP routes (HTML + API)
├── pages/            # UI templates
├── lib/              # External clients (YouTube, RSS, downloader)
├── db/               # Database operations
├── migrations/       # Database schema changes
├── assets/           # Static files
└── docker-compose.yml
```

## Known Limitations

- Background jobs are disabled in development mode (`DISABLE_JOBS=true`)
- Mobile responsiveness is still being improved
- User authentication is not yet implemented

## Roadmap

- [ ] User authentication and multi-user support
- [ ] Enhanced mobile UI
- [ ] Better error handling
- [ ] Advanced search and filtering
- [ ] Notification system for new videos
- [ ] API documentation

## Support

For development setup, code architecture, or contributing, see [`AGENTS.md`](AGENTS.md).

---

Built with ❤️ using Go, PostgreSQL, and Templ.
