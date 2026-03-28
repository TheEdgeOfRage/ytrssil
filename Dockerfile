FROM golang:1.25-trixie@sha256:84ad9c55a9914b5861e04c6391151c2bb8931f4115acc9ab7658b01cc85e2854 AS builder

WORKDIR /app/
ENV GOCACHE="/cache"
ENV PORT="80"

COPY go.mod go.sum /app/
RUN --mount=type=cache,target="/cache" go mod download;

COPY . /app/
RUN --mount=type=cache,target="/cache" go build -o dist/ytrssil-api cmd/main.go;

FROM debian:trixie-slim@sha256:26f98ccd92fd0a44d6928ce8ff8f4921b4d2f535bfa07555ee5d18f61429cf0c AS api

HEALTHCHECK --start-period=2s --start-interval=2s CMD exec curl -sf localhost:$PORT/healthz;
WORKDIR /app/
ENTRYPOINT ["./ytrssil-api"]
VOLUME /var/lib/ytrssil/downloads

RUN apt update \
	&& apt install -y ca-certificates curl yt-dlp \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*;

COPY ./assets/ ./assets/
RUN assets/load.sh && rm assets/load.sh;

COPY --from=builder /app/dist/ ./

FROM migrate/migrate:4 AS migrations
COPY ./migrations/ /migrations/
