package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/handler"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
)

const LevelFatal slog.Level = slog.LevelError + 4

func init() {
	// always use UTC
	time.Local = time.UTC
}

func fetcherRoutine(l *slog.Logger, h handler.Handler) {
	for {
		err := h.FetchVideos(context.Background())
		if err != nil {
			l.Error("Failed to fetch videos", "call", "handler.FetchVideos", "err", err)
		}
		time.Sleep(5 * time.Minute)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := config.Parse()
	if err != nil {
		logger.Log(context.Background(), LevelFatal, "Failed to parse config", "call", "config.Parse", "error", err)
		return
	}
	db, err := db.NewPostgresDB(logger, config.DB)
	if err != nil {
		logger.Log(
			context.Background(), LevelFatal,
			"Failed to create DB connection",
			"call", "db.NewPostgresDB",
			"error", err,
		)
		return
	}
	parser := feedparser.NewParser(logger)
	handler := handler.New(logger, db, parser)
	gin.SetMode(gin.ReleaseMode)
	router, err := ytrssil.SetupGinRouter(
		logger,
		handler,
		auth.AuthMiddleware(db),
	)
	if err != nil {
		logger.Log(
			context.Background(), LevelFatal,
			"Failed to set up gin server",
			"call", "ytrssil.SetupGinServer",
			"error", err,
		)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Gin.Port),
		Handler: router,
	}
	quit := make(chan os.Signal, 1)
	// handle Interrupt (ctrl-c) Term, used by `kill` et al, HUP which is commonly used to reload configs
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		s := <-quit
		logger.Info("Received signal, shutting down", "signal", s)
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error(
				"Failed to shutdown server",
				"call", "server.Shutdown",
				"error", err,
			)
		}
	}()

	// start periodic fetch videos routine
	go fetcherRoutine(logger, handler)

	logger.Info("ytrssil API is starting up", "port", config.Gin.Port)
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			logger.Error("Server crashed", "call", "server.ListenAndServe", "error", err)
		}
	}
	logger.Info("exit complete")
}
