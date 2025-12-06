package ytrssil

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/handler"
)

type server struct {
	log     *slog.Logger
	handler handler.Handler
}

func NewServer(log *slog.Logger, handler handler.Handler) (*server, error) {
	return &server{log: log, handler: handler}, nil
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
func SetupGinRouter(l *slog.Logger, handler handler.Handler, authMiddleware func(c *gin.Context)) (*gin.Engine, error) {
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

	srv, err := NewServer(l, handler)
	if err != nil {
		return nil, err
	}
	engine.GET("/healthz", srv.Healthz)
	engine.POST("/register", srv.CreateUserJSON)
	engine.POST("/fetch", srv.FetchVideosJSON)

	pages := engine.Group("")
	pages.Use(authMiddleware)

	pages.GET("/", srv.NewVideosPage)
	pages.POST("/videos/:video_id/watch", srv.MarkVideoAsWatchedPage)

	// all APIs go in this routing group and require authentication
	api := engine.Group("/api")
	api.Use(authMiddleware)
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
