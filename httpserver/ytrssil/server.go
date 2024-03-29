package ytrssil

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/handler"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
)

type server struct {
	log     log.Logger
	handler handler.Handler
}

func NewServer(log log.Logger, handler handler.Handler) (*server, error) {
	return &server{log: log, handler: handler}, nil
}

func (s *server) Healthz(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}

// SetupGinRouter sets up routes for all APIs on a Gin server (aka router)
func SetupGinRouter(l log.Logger, handler handler.Handler, authMiddleware func(c *gin.Context)) (*gin.Engine, error) {
	engine := gin.New()
	// Middlewares are executed top to bottom in a stack-like manner
	engine.Use(
		gin.LoggerWithFormatter(log.GinFormatterWithUTCAndBodySize),
		gin.Recovery(), // Recovery needs to go before other middlewares to catch panics
	)

	srv, err := NewServer(l, handler)
	if err != nil {
		return nil, err
	}
	engine.GET("/healthz", srv.Healthz)
	engine.POST("/register", srv.CreateUser)
	engine.POST("/fetch", srv.FetchVideos)

	// all APIs go in this routing group and require authentication
	api := engine.Group("/api")
	api.Use(authMiddleware)
	{
		api.POST("channels/:channel_id/subscribe", srv.SubscribeToChannel)
		api.POST("channels/:channel_id/unsubscribe", srv.UnsubscribeFromChannel)
		api.GET("videos/new", srv.GetNewVideos)
		api.GET("videos/watched", srv.GetWatchedVideos)
		api.POST("videos/:video_id/watch", srv.MarkVideoAsWatched)
		api.POST("videos/:video_id/unwatch", srv.MarkVideoAsUnwatched)
	}

	return engine, nil
}
