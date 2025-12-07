package ytrssil

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
	"github.com/TheEdgeOfRage/ytrssil-api/handler"
	"github.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
)

type server struct {
	log     *slog.Logger
	cfg     config.Config
	handler handler.Handler
}

func NewServer(log *slog.Logger, cfg config.Config, handler handler.Handler) (*server, error) {
	return &server{log: log, cfg: cfg, handler: handler}, nil
}

func ginLogFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("time=%v method=%s path=%#v status=%d ip=%s size=%v duration=%v\n",
		param.TimeStamp.UTC().Format(time.RFC3339Nano),
		param.Method,
		param.Path,
		param.StatusCode,
		param.ClientIP,
		param.BodySize,
		param.Latency,
	)
}

func (srv *server) Healthz(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}

// SetupGinRouter sets up routes for all APIs on a Gin server (aka router)
func SetupGinRouter(
	l *slog.Logger,
	cfg config.Config,
	handler handler.Handler,
) (*gin.Engine, error) {
	engine := gin.New()
	// Middlewares are executed top to bottom in a stack-like manner
	engine.Use(
		gin.LoggerWithFormatter(ginLogFormatter),
		gin.Recovery(), // Recovery needs to go before other middlewares to catch panics
	)
	engine.HandleMethodNotAllowed = true
	engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "URL not found"})
	})
	engine.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "HTTP method not allowed"})
	})

	srv, err := NewServer(l, cfg, handler)
	if err != nil {
		return nil, err
	}
	engine.GET("/healthz", srv.Healthz)
	engine.POST("/fetch", srv.FetchVideosJSON)
	engine.Static("/assets", "./assets")

	engine.GET("/auth", srv.AuthPage)
	engine.POST("/auth", srv.HandleAuth)

	pages := engine.Group("")
	pages.Use(auth.PageAuthMiddleware(cfg.AuthToken))
	{
		pages.GET("/", srv.NewVideosPage)
		pages.POST("/subscribe", srv.SubscribeToChannelPage)
		pages.PATCH("/videos/:video_id/watch", srv.MarkVideoAsWatchedPage)
		pages.PATCH("/videos/:video_id/progress", srv.SetVideoProgressPage)
	}

	api := engine.Group("/api")
	api.Use(auth.APIAuthMiddleware(cfg.AuthToken))
	{
		api.POST("channels/:channel_id/subscribe", srv.SubscribeToChannelJSON)
		api.POST("channels/:channel_id/unsubscribe", srv.UnsubscribeFromChannelJSON)
		api.GET("videos/new", srv.GetNewVideosJSON)
		api.GET("videos/watched", srv.GetWatchedVideosJSON)
		api.POST("videos/:video_id/watch", srv.MarkVideoAsWatchedJSON)
		api.POST("videos/:video_id/unwatch", srv.MarkVideoAsUnwatchedJSON)
	}

	return engine, nil
}
