# Ytrssil API - Project Overview

## Purpose
Ytrssil API is a YouTube RSS-to-interface layer application. It provides an HTTP API for managing YouTube channel subscriptions and video tracking with watch progress functionality.

## Tech Stack
- **Language**: Go 1.25
- **Web Framework**: Gin (github.com/gin-gonic/gin v1.11.0)
- **Templating**: templ (github.com/a-h/templ v0.3.960) for HTML generation
- **Database**: PostgreSQL (via lib/pq driver)
- **Testing**: testify (github.com/stretchr/testify v1.11.1)
- **Hot Reload**: air (for development)
- **Containerization**: Docker with buildx

## Project Structure
```
.
├── cmd/               # Entry point (main.go)
├── httpserver/        # HTTP server implementations
│   ├── ytrssil/      # Main API routes
│   └── auth/         # Authentication routes
├── handler/          # Business logic handlers
│   ├── channels.go
│   ├── videos.go
│   └── handler_test.go
├── models/           # Data models
│   ├── channel.go
│   └── video.go
├── db/               # Database layer
│   ├── db.go
│   ├── psql.go
│   ├── channels.go
│   └── videos.go
├── feedparser/       # RSS feed parsing logic
├── pages/            # templ templates for HTML pages
├── lib/              # Shared libraries
├── mocks/            # Generated mock interfaces for testing
├── migrations/       # Database migrations
├── config/           # Configuration files
└── assets/           # Static assets
```

## Dependencies
Database runs on PostgreSQL 18 via Docker Compose.