package ytrssil_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
	"github.com/TheEdgeOfRage/ytrssil-api/handler"
	"github.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
)

func init() {
	// always use UTC
	time.Local = time.UTC
}

type APITestSuite struct {
	suite.Suite

	cfg    config.Config
	server *http.Server
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	s.cfg = config.TestConfig()

	handler := handler.New(l, nil, nil, nil, nil, s.cfg.DownloadsDir)

	gin.SetMode(gin.TestMode)
	router, err := ytrssil.SetupGinRouter(l, s.cfg, handler)
	s.Require().NoError(err)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.Port),
		Handler: router,
	}
}

func (s *APITestSuite) TestHealthz() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Equal("healthy", w.Body.String())
}
