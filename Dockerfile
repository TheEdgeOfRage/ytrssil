FROM golang:1.25-trixie AS builder

WORKDIR /app/
ENV GOCACHE="/cache"
ENV PORT="80"

COPY go.mod go.sum /app/
RUN --mount=type=cache,target="/cache" go mod download

COPY . /app/
RUN --mount=type=cache,target="/cache" go build -o dist/ytrssil-api cmd/main.go

FROM debian:trixie-slim AS api

HEALTHCHECK --start-period=2s --start-interval=2s CMD exec curl -sf localhost:$PORT/healthz
WORKDIR /app/
ENTRYPOINT ["./ytrssil-api"]

RUN apt update \
	&& apt install -y ca-certificates curl \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

COPY ./assets/ ./assets/
RUN assets/load.sh && rm assets/load.sh

COPY --from=builder /app/dist/ ./

FROM migrate/migrate:4 AS migrations
COPY ./migrations/ /migrations/
