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
	"github.com/stretchr/testify/assert"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
	"github.com/TheEdgeOfRage/ytrssil-api/handler"
	"github.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
	"github.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
)

var testConfig config.Config

func init() {
	// always use UTC
	time.Local = time.UTC
	testConfig = config.TestConfig()
}

func setupTestServer(t *testing.T) *http.Server {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))

	handler := handler.New(l, nil, nil)

	gin.SetMode(gin.TestMode)
	router, err := ytrssil.SetupGinRouter(
		l,
		handler,
		auth.APIAuthMiddleware(""),
		auth.PageAuthMiddleware(""),
	)
	assert.Nil(t, err)

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", testConfig.Port),
		Handler: router,
	}
}

func TestHealthz(t *testing.T) {
	server := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "healthy", w.Body.String())
}
