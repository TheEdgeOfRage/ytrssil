## IMPORTANT

- Don't edit the generated templ code or mock files
- Don't invoke ANY go commands, the development workflow is orchestrated with make targets

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
└── assets/           # Static files
```

## Development Workflow

All Go tooling is managed via `make` targets. Never run `go` commands directly.

**Make Targets**:

- `make lint` — runs `go mod tidy`, `go vet`, and golangci-lint (lint + fmt). Use to check code quality.
- `make test` — runs `go test -timeout=30s -race ./...`
- `make templ` — generates Go code from `.templ` files via `templ generate`. Required before build.
- `make gen-mocks` — regenerates mock files in `mocks/` using moq. Run after changing interfaces in `db/`, `feedparser/`, or `lib/clients/youtube/`.
- `make migrate` — applies PostgreSQL migrations from `migrations/` dir. Configurable via `DB_URI` env var (default: `postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable`)

**Tool binaries** are installed into `bin/` locally (golangci-lint, moq, migrate) and auto-installed by their respective targets.
