.PHONY: all setup ytrssil-api build gen-mocks lint yamllint test migrate image-build image-push

DB_URI ?= postgres://ytrssil:ytrssil@localhost:5432/ytrssil?sslmode=disable

bin/moq:
	GOBIN=$(PWD)/bin go install github.com/matryer/moq@v0.6.0
bin/golangci-lint:
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
bin/migrate: bin
	GOBIN=$(PWD)/bin go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0

lint: bin/golangci-lint
	go mod tidy
	go vet ./...
	bin/golangci-lint -c .golangci.yml fmt ./...
	bin/golangci-lint -c .golangci.yml run ./...

test:
	go mod tidy
	go test -timeout=10s -race -benchmem ./...

gen-mocks: bin/moq
	./bin/moq -pkg db_mock -out ./mocks/db/db.go ./db DB
	./bin/moq -pkg parser_mock -out ./mocks/feedparser/feedparser.go ./feedparser Parser
	go fmt ./...

migrate: bin/migrate
	bin/migrate -database "$(DB_URI)" -path migrations up

image-build:
	docker buildx build --push -t theedgeofrage/ytrssil:api --target api .
	docker buildx build --push -t theedgeofrage/ytrssil:migrations --target migrations .
