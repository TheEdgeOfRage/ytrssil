package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/handler"
	"github.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
	"github.com/TheEdgeOfRage/ytrssil-api/lib/clients/youtube"
)

func init() {
	// always use UTC
	time.Local = time.UTC
}

func fetcherRoutine(ctx context.Context, l *slog.Logger, h handler.Handler) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			l.Info("Fetcher context done, stopping fetcher")
			return
		case <-ticker.C:
			err := h.FetchVideos(context.Background())
			if err != nil {
				l.Error("Failed to fetch videos", "call", "handler.FetchVideos", "err", err)
			}
		}
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Parse()
	if err != nil {
		logger.Error("Failed to parse config", "call", "config.Parse", "error", err)
		return
	}
	db, err := db.NewPostgresDB(logger, cfg.DBURI)
	if err != nil {
		logger.Error(
			"Failed to create DB connection",
			"call", "db.NewPostgresDB",
			"error", err,
		)
		return
	}
	parser := feedparser.NewParser(logger)
	youTubeClient := youtube.NewYouTubeClient(logger, cfg.YouTubeAPIKey)
	handler := handler.New(logger, db, parser, youTubeClient)
	if cfg.Dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router, err := ytrssil.SetupGinRouter(logger, cfg, handler)
	if err != nil {
		logger.Error(
			"Failed to set up gin server",
			"call", "ytrssil.SetupGinServer",
			"error", err,
		)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	wg := &sync.WaitGroup{}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// start periodic fetch videos routine
	fetcherContext, cancelFetcher := context.WithCancel(context.Background())
	wg.Go(func() {
		fetcherRoutine(fetcherContext, logger, handler)
	})

	wg.Go(func() {
		logger.Info("ytrssil API is starting up", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Error("Server crashed", "call", "server.ListenAndServe", "error", err)
			}
		}
	})

	s := <-quit
	logger.Info("Received signal, shutting down", "signal", s)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error(
			"Failed to shutdown server",
			"call", "server.Shutdown",
			"error", err,
		)
	}
	cancelFetcher()

	wg.Wait()
	logger.Info("exit complete")
}
