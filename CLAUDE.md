- Don't use serena for templ files, it only supports normal go code.
- Don't attempt to edit the generated templ code
- Don't invoke ANY go commands other than `go vet` and `go mod tidy`

## Architecture

**Ytrssil** is a YouTube RSS feed aggregator with video download capabilities built with Go, Gin, PostgreSQL, and Templ for server-side rendering.

**Tech Stack**: Go 1.25, Gin (HTTP), Templ (SSR), PostgreSQL (pgx), yt-dlp (downloads), HTMX (frontend interactivity), Bootstrap 5 (UI)

**Layers**:
- **HTTP**: Gin router with dual interfaces (HTML pages + JSON API), authentication middleware, HTMX-enabled templates
- **Handler**: Business logic layer orchestrating channels, videos, downloads, and cleanup routines
- **Database**: PostgreSQL with pgx connection pooling, 7 migrations tracking schema evolution
- **External**: YouTube API client, RSS feed parser, yt-dlp downloader wrapper

**Key Features**: Subscribe to YouTube channels via RSS, track watched/unwatched videos with progress, async video downloads with status polling, automatic file cleanup (2 days after watched), YouTube Shorts detection

**Background Jobs**: Video fetcher (5min intervals), cleanup routine (1hr intervals) - both disabled in dev mode

**Directory Structure**: `cmd/` (main), `config/`, `db/`, `models/`, `handler/`, `httpserver/`, `pages/` (templ), `lib/clients/`, `lib/downloader/`, `feedparser/`, `migrations/`, `assets/`
