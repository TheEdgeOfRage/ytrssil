.PHONY: all lint govulncheck test air templ build gen-mocks migrate image-build

DB_URI ?= postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable

bin:
	mkdir bin
bin/moq: bin
	GOBIN=$(PWD)/bin go install github.com/matryer/moq@v0.7.1
bin/golangci-lint: bin
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
bin/migrate: bin
	GOBIN=$(PWD)/bin go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
bin/air: bin
	GOBIN=$(PWD)/bin go install github.com/air-verse/air@v1.65.1
bin/templ: bin
	GOBIN=$(PWD)/bin go install github.com/a-h/templ/cmd/templ@v0.3.1001
bin/govulncheck: bin
	GOBIN=$(PWD)/bin go install golang.org/x/vuln/cmd/govulncheck@v1.2.0

fmt:
	go mod tidy
	bin/golangci-lint fmt

lint: bin/golangci-lint
	go mod tidy -diff
	go vet ./...
	bin/golangci-lint run
	bin/golangci-lint fmt -d
	$(MAKE) govulncheck

govulncheck:  bin/govulncheck
	bin/govulncheck ./...

test:
	go test -timeout=30s -race ./...

air: bin/air
	@./bin/air -c .air.toml

templ: bin/templ
	bin/templ generate

build: templ
	go build -o dist/ytrssil ./cmd/main.go

gen-mocks: bin/moq
	./bin/moq -pkg db_mock -out ./mocks/db/db.go ./db DB
	./bin/moq -pkg parser_mock -out ./mocks/feedparser/feedparser.go ./feedparser Parser
	./bin/moq -pkg youtube_mock -out ./mocks/youtube/youtube.go ./lib/clients/youtube Client
	go fmt ./...

migrate: bin/migrate
	bin/migrate -database "$(DB_URI)" -path migrations up

image-build:
	docker buildx build --push -t theedgeofrage/ytrssil:api --target api .
