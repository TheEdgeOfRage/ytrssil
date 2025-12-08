# Suggested Commands

## Development Commands

### Building
```bash
make build          # Generate templates and build binary to dist/ytrssil
templ generate      # Generate Go code from templ templates
```

### Testing
```bash
make test           # Run tests with race detection and benchmarks (10s timeout)
go test -timeout=10s -race -benchmem ./...
```

### Linting and Formatting
```bash
make lint           # Run full linting pipeline:
                    # - go mod tidy
                    # - go vet
                    # - golangci-lint fmt (format)
                    # - golangci-lint run (lint)
```

### Development Server
```bash
make air            # Start hot-reload development server using air
```

### Mock Generation
```bash
make gen-mocks      # Generate mock implementations for testing:
                    # - db/DB interface
                    # - feedparser/Parser interface
                    # - lib/clients/youtube/Client interface
```

### Database
```bash
make migrate        # Run database migrations
                    # Default: postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable
                    # Override with: make migrate DB_URI=<your-uri>
```

### Docker
```bash
make image-build    # Build and push Docker images:
                    # - theedgeofrage/ytrssil:api
                    # - theedgeofrage/ytrssil:migrations

docker compose up   # Start full stack (postgres + api + migrations)
```

## Binary Installation (in bin/)
The Makefile installs tools locally to avoid global pollution:
- `bin/moq`: Mock generator
- `bin/golangci-lint`: Linter
- `bin/migrate`: Database migration tool
- `bin/air`: Hot reload development server